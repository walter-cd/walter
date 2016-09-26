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
	Stdout   string
}

func (t *Task) Run() error {
	if len(t.Parallel) > 0 {
		runParallel(t.Parallel)
	}

	if len(t.Serial) > 0 {
		runSerial(t.Serial)
	}

	if t.Command != "" {
		cmd := exec.Command("sh", "-c", t.Command)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		cmd.Start()

		in := bufio.NewScanner(stdout)
		for in.Scan() {
			t.Stdout += in.Text()
		}

		cmd.Wait()
	}

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
