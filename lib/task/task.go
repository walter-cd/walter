package task

import (
	"bufio"
	"os/exec"
	"sync"
)

type Task struct {
	Name     string
	Command  string
	Parallel []Task
	Serial   []Task
	Stdout   []string
	Stderr   []string
}

func (t *Task) Run() error {
	if len(t.Parallel) > 0 {
		runParallel(t.Parallel)
	}

	if len(t.Serial) > 0 {
		runSerial(t.Serial)
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

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			t.Stdout = append(t.Stdout, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			t.Stderr = append(t.Stderr, scanner.Text())
		}
	}()

	cmd.Wait()

	return nil
}

func runParallel(tasks []Task) {
	var wg sync.WaitGroup
	for _, t := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			t.Run()
		}(t)
	}
	wg.Wait()
}

func runSerial(tasks []Task) {
	for _, t := range tasks {
		t.Run()
	}
}
