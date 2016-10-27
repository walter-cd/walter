package pipeline

import (
	"io/ioutil"
	"sync"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/context"

	"github.com/go-yaml/yaml"
	"github.com/walter-cd/walter/lib/notify"
	"github.com/walter-cd/walter/lib/task"
)

type Pipeline struct {
	Build     Build
	Deploy    Deploy
	Notifiers []notify.Notifier
}

type Build struct {
	Tasks   Tasks
	Cleanup Tasks
}

type Deploy struct {
	Tasks   Tasks
	Cleanup Tasks
}

type Tasks []*task.Task

func Load(b []byte) (Pipeline, error) {
	p := Pipeline{}
	err := yaml.Unmarshal(b, &p)
	if err != nil {
		return p, err
	}

	p.Notifiers, err = notify.NewNotifiers(b)

	return p, err
}

func LoadFromFile(file string) (Pipeline, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return Pipeline{}, err
	}
	return Load(data)
}

func (p *Pipeline) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	p.runTasks(ctx, cancel, p.Build.Tasks)

	ctx, cancel = context.WithCancel(context.Background())
	p.runTasks(ctx, cancel, p.Build.Cleanup)
}

func (p *Pipeline) runTasks(ctx context.Context, cancel context.CancelFunc, tasks Tasks) {
	failed := false
	for i, t := range tasks {
		if len(t.Parallel) > 0 {
			p.runParallel(ctx, cancel, t.Parallel)
			t.Status = task.Succeeded
			for _, child := range t.Parallel {
				if child.Status == task.Failed {
					t.Status = task.Failed
				}
			}
		}

		if len(t.Serial) > 0 {
			p.runTasks(ctx, cancel, t.Serial)
			t.Status = task.Succeeded
			for _, child := range t.Serial {
				if child.Status == task.Failed {
					t.Status = task.Failed
				}
			}
		}

		if failed || (i > 0 && tasks[i-1].Status == task.Failed) {
			t.Status = task.Skipped
			failed = true
			log.Infof("[%s] Task skipped because previous task failed", t.Name)
			continue
		}

		log.Infof("[%s] Start task", t.Name)
		err := t.Run(ctx, cancel)
		if err != nil {
			log.Errorf("[%s] %s", t.Name, err)
		}

		if t.Status == task.Succeeded {
			log.Infof("[%s] End task", t.Name)
		}

		for _, n := range p.Notifiers {
			n.Notify(t)
		}
	}
}

func (p *Pipeline) runParallel(ctx context.Context, cancel context.CancelFunc, tasks Tasks) {
	var wg sync.WaitGroup
	for _, t := range tasks {
		wg.Add(1)
		go func(t *task.Task) {
			defer wg.Done()
			log.Infof("[%s] Start task", t.Name)
			t.Run(ctx, cancel)
			if t.Status == task.Succeeded {
				log.Infof("[%s] End task", t.Name)
			}

			for _, n := range p.Notifiers {
				n.Notify(t)
			}
		}(t)
	}
	wg.Wait()
}
