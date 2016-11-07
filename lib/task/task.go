package task

import (
	"bytes"
	"io"
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
	Stdout         *bytes.Buffer
	Stderr         *bytes.Buffer
	CombinedOutput *bytes.Buffer
	Status         int
	Cmd            *exec.Cmd
	Include        string
}

type outputHandler struct {
	writer io.Writer
	copy   io.Writer
	mu     *sync.Mutex
}

func (t *Task) Run(ctx context.Context, cancel context.CancelFunc, prevTask *Task) error {
	if t.Command == "" {
		return nil
	}

	log.Infof("[%s] Start task", t.Name)

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

	if prevTask != nil && prevTask.Stdout != nil {
		t.Cmd.Stdin = bytes.NewBuffer(prevTask.Stdout.Bytes())
	}

	t.Stdout = new(bytes.Buffer)
	t.Stderr = new(bytes.Buffer)
	t.CombinedOutput = new(bytes.Buffer)

	var mu sync.Mutex
	t.Cmd.Stdout = &outputHandler{t.Stdout, t.CombinedOutput, &mu}
	t.Cmd.Stderr = &outputHandler{t.Stderr, t.CombinedOutput, &mu}

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

	if t.Status == Succeeded {
		log.Infof("[%s] End task", t.Name)
	}

	return nil
}

func (o *outputHandler) Write(b []byte) (int, error) {
	log.Info(strings.TrimSuffix(string(b), "\n"))

	o.mu.Lock()
	defer o.mu.Unlock()
	o.writer.Write(b)
	o.copy.Write(b)

	return len(b), nil
}
