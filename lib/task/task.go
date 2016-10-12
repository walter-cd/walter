package task

import (
	"bufio"
	"os/exec"
	"sync"

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
	Parallel       Parallel
	Serial         Tasks
	Stdout         []string
	Stderr         []string
	CombinedOutput []string
	Status         int
	Cmd            *exec.Cmd
}

type Tasks []*Task
type Parallel []*Task

func (tasks Tasks) Run(ctx context.Context, cancel context.CancelFunc) {
	failed := false
	for i, t := range tasks {
		if failed || (i > 0 && tasks[i-1].Status == Failed) {
			t.Status = Skipped
			failed = true
			log.Infof("[%s] Task skipped because previous task failed", t.Name)
			continue
		}
		log.Infof("[%s] Start task", t.Name)

		err := t.Run(ctx, cancel)
		if err != nil {
			log.Errorf("[%s] %s", t.Name, err)
		}

		if t.Status == Succeeded {
			log.Infof("[%s] End task", t.Name)
		}
	}
}

func (t *Task) Run(ctx context.Context, cancel context.CancelFunc) error {
	if len(t.Parallel) > 0 {
		t.Parallel.Run(ctx, cancel)
		t.Status = Succeeded
		// Set Failed to parent task if one of parallel tasks failed
		for _, task := range t.Parallel {
			if task.Status == Failed {
				t.Status = Failed
			}
		}
	}

	if len(t.Serial) > 0 {
		t.Serial.Run(ctx, cancel)
		t.Status = Succeeded
		// Set Failed to parent task if one of serial tasks failed
		for _, task := range t.Serial {
			if task.Status == Failed {
				t.Status = Failed
			}
		}
	}

	if t.Command == "" {
		return nil
	}

	t.Cmd = exec.Command("sh", "-c", t.Command)

	if t.Directory != "" {
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
					stdoutPipe.Close()
					stderrPipe.Close()
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

func (tasks Parallel) Run(ctx context.Context, cancel context.CancelFunc) {
	var wg sync.WaitGroup
	for _, t := range tasks {
		wg.Add(1)
		go func(t *Task) {
			defer wg.Done()
			log.Infof("[%s] Start task", t.Name)
			t.Run(ctx, cancel)
			if t.Status == Succeeded {
				log.Infof("[%s] End task", t.Name)
			}
		}(t)
	}
	wg.Wait()
}
