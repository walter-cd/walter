package notify

import (
	"github.com/go-yaml/yaml"
	"github.com/walter-cd/walter/lib/task"
)

type Notifier interface {
	Notify(*task.Task) error
}

type Notify struct {
	Notify []map[string]string
}

type Default struct{}

func NewNotifiers(b []byte) ([]Notifier, error) {
	notify := Notify{}
	err := yaml.Unmarshal(b, &notify)

	var notifiers []Notifier
	for _, n := range notify.Notify {
		switch n["type"] {
		case "slack":
			notifiers = append(notifiers, NewSlack(n))
		default:
			notifiers = append(notifiers, &Default{})
		}
	}

	return notifiers, err
}

func (d *Default) Notify(t *task.Task) error {
	return nil
}
