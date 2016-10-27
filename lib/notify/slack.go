package notify

import (
	"github.com/drone/drone/shared/build/log"
	"github.com/walter-cd/walter/lib/task"
)

type Slack struct {
	Channel  string
	URL      string
	IconURL  string
	UserName string
}

func NewSlack(m map[string]string) *Slack {
	s := &Slack{}
	s.Channel = m["channel"]
	s.URL = m["url"]
	s.IconURL = m["icon_url"]
	s.UserName = m["username"]
	return s
}

func (s *Slack) Notify(t *task.Task) error {
	if s.Channel[0] != '#' {
		s.Channel = "#" + s.Channel
	}

	log.Info("Notify!!!!")
	//resp, err := http.PostForm(s.url)
	return nil
}
