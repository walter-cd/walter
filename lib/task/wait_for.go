package task

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

type WaitFor struct {
	Host  string
	Port  int
	File  string
	State string
	Delay float64
}

func (t *Task) wait() error {
	err := t.WaitFor.validate()
	if err != nil {
		return err
	}
	t.WaitFor.wait(t)
	return nil
}

func (w *WaitFor) wait(t *Task) {
	switch {
	case w.Delay > 0.0:
		w.waitForDelay(t)
	case w.Port > 0:
		w.waitForPort(t)
	case w.File != "":
		w.waitForFile(t)
	}
}

func (w *WaitFor) waitForDelay(t *Task) {
	log.Infof("[%s] wait_for: %f seconds delay", t.Name, w.Delay)
	time.Sleep(time.Duration(w.Delay) * time.Second)
	return
}

func (w *WaitFor) waitForPort(t *Task) {
	log.Infof("[%s] wait_for: %s:%d %s", t.Name, w.Host, w.Port, w.State)
	if w.State == "ready" || w.State == "present" {
		for {
			if connected(w.Host, w.Port) {
				return
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	if w.State == "unready" || w.State == "absent" {
		for {
			if !connected(w.Host, w.Port) {
				return
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func connected(host string, port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, strconv.Itoa(port)))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (w *WaitFor) waitForFile(t *Task) {
}

func (w *WaitFor) validate() error {
	switch {
	case w.Port != 0 && w.File != "":
		return errors.New("wait_for: cannot use port and file at the same time")
	case w.Port != 0 && w.Delay > 0.0:
		return errors.New("wait_for: cannot use port and delay at the same time")
	case w.File != "" && w.Delay > 0.0:
		return errors.New("wait_for: cannot use file and delay at the same time")
	case w.Delay < 0:
		return errors.New("wait_for: delay must be positive")
	case w.Port < 0:
		return errors.New("wait_for: port must be positive")
	case w.Port > 0 && w.Host == "":
		return errors.New("wait_for: cannot use port without host")
	case w.Port == 0 && w.Host != "":
		return errors.New("wait_for: cannot use host without port")
	case w.Delay == 0 && !includes([]string{"ready", "unready", "present", "absent"}, w.State):
		return errors.New("wait_for: state does not support " + w.State)
	case w.Port > 0 && w.State == "":
		return errors.New("wait_for: cannot use port without state")
	case w.File != "" && w.State == "":
		return errors.New("wait_for: cannot use file without state")
	}

	return nil
}

func includes(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
