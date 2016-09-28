package task

import (
	"bufio"
	"os/exec"
	"sync"

	log "github.com/Sirupsen/logrus"
)

const (
	Running = iota
	Succeeded
	Failed
	Skipped
	Aborted
)

type Task struct {
	Name           string
	Command        string
	Parallel       Parallel
	Serial         Tasks
	Stdout         []string
	Stderr         []string
	CombinedOutput []string
	Status         int
}

type Tasks []*Task
type Parallel []*Task

func (tasks Tasks) Run() {
	failed := false
	for i, t := range tasks {
		if failed || (i > 0 && tasks[i-1].Status == Failed) {
			t.Status = Skipped
			failed = true
			log.Infof("[%s] Task skipped because previous task failed", t.Name)
			continue
		}
		log.Infof("[%s] Start task", t.Name)
		t.Run()
		log.Infof("[%s] End task", t.Name)
	}
}

func (t *Task) Run() error {
	if len(t.Parallel) > 0 {
		t.Parallel.Run()
		t.Status = Succeeded
		for _, task := range t.Parallel {
			if task.Status == Failed {
				t.Status = Failed
			}
		}
	}

	if len(t.Serial) > 0 {
		t.Serial.Run()
		t.Status = Succeeded
		for _, task := range t.Serial {
			if task.Status == Failed {
				t.Status = Failed
			}
		}
	}

	if t.Command == "" {
		return nil
	}

	cmd := exec.Command("sh", "-c", t.Command)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	t.Status = Running

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			text := scanner.Text()
			t.Stdout = append(t.Stdout, text)
			t.CombinedOutput = append(t.CombinedOutput, text)
			log.Infof("[%s] %s", t.Name, text)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			text := scanner.Text()
			t.Stderr = append(t.Stderr, text)
			t.CombinedOutput = append(t.CombinedOutput, text)
			log.Infof("[%s] %s", t.Name, text)
		}
	}()

	wg.Wait()

	cmd.Wait()

	if cmd.ProcessState.Success() {
		t.Status = Succeeded
	} else {
		log.Errorf("[%s] Task failed", t.Name)
		t.Status = Failed
	}

	return nil
}

func (tasks Parallel) Run() {
	var wg sync.WaitGroup
	for _, t := range tasks {
		wg.Add(1)
		go func(t *Task) {
			defer wg.Done()
			log.Infof("[%s] Start task", t.Name)
			t.Run()
			log.Infof("[%s] End task", t.Name)
		}(t)
	}
	wg.Wait()
}
