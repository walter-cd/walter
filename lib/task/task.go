package task

import (
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
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
	Stdout         []string
	Stderr         []string
	CombinedOutput []string
	Status         int
	Cmd            *exec.Cmd
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

	stdoutPipe, err := t.Cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := t.Cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := t.Cmd.Start(); err != nil {
		t.Status = Failed
		return err
	}

	t.Status = Running

	var wg sync.WaitGroup

	wg.Add(1)
	go func(t *Task) {
		defer wg.Done()
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			text := scanner.Text()
			t.Stdout = append(t.Stdout, text)
			t.CombinedOutput = append(t.CombinedOutput, text)
			log.Infof("[%s] %s", t.Name, text)
		}
	}(t)

	wg.Add(1)
	go func(t *Task) {
		defer wg.Done()
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			text := scanner.Text()
			t.Stderr = append(t.Stderr, text)
			t.CombinedOutput = append(t.CombinedOutput, text)
			log.Infof("[%s] %s", t.Name, text)
		}
	}(t)

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

	wg.Wait()

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
