package task

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
)

const (
	Init = iota
	Running
	Succeeded
	Failed
	Skipped
	Aborted
)

type Task struct {
	Name           string
	Command        string
	Directory      string
	Parallel       []*Task
	Serial         []*Task
	Stdout         *bytes.Buffer
	Stderr         *bytes.Buffer
	CombinedOutput *bytes.Buffer
	Status         int
	Cmd            *exec.Cmd
}

type copyWriter struct {
	original io.Writer
	copy     io.Writer
}

func (t *Task) Run(ctx context.Context, cancel context.CancelFunc) error {
	if t.Command == "" {
		return nil
	}

	t.Cmd = exec.Command("sh", "-c", t.Command)
	t.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if t.Directory != "" {
		re := regexp.MustCompile(`\$[A-Z1-9\-_]+`)
		matches := re.FindAllString(t.Directory, -1)
		for _, m := range matches {
			env := os.Getenv(strings.TrimPrefix(m, "$"))
			t.Directory = strings.Replace(t.Directory, m, env, -1)
		}
		t.Cmd.Dir = t.Directory
	}

	t.Stdout = new(bytes.Buffer)
	t.Stderr = new(bytes.Buffer)
	t.CombinedOutput = new(bytes.Buffer)

	t.Cmd.Stdout = &copyWriter{t.Stdout, t.CombinedOutput}
	t.Cmd.Stderr = &copyWriter{t.Stderr, t.CombinedOutput}

	if err := t.Cmd.Start(); err != nil {
		t.Status = Failed
		return err
	}

	t.Status = Running

	go func(t *Task) {
		for {
			select {
			case <-ctx.Done():
				if t.Status == Running {
					t.Status = Aborted
					t.Cmd.Process.Kill()
					pgid, err := syscall.Getpgid(t.Cmd.Process.Pid)
					if err == nil {
						syscall.Kill(-pgid, syscall.SIGTERM)
					}
					log.Warnf("[%s] aborted", t.Name)
				}
				return
			}
		}
	}(t)

	t.Cmd.Wait()

	if t.Cmd.ProcessState.Success() {
		t.Status = Succeeded
	} else if t.Status == Running {
		log.Errorf("[%s] Task failed", t.Name)
		t.Status = Failed
		cancel()
	}

	return nil
}

func (c *copyWriter) Write(b []byte) (int, error) {
	log.Info(strings.TrimSuffix(string(b), "\n"))

	c.original.Write(b)
	c.copy.Write(b)

	return len(b), nil
}
