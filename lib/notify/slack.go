package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/walter-cd/walter/lib/task"
)

type Slack struct {
	Channel  string `json:"channel"`
	URL      string `json:"-"`
	IconURL  string `json:"icon_url"`
	UserName string `json:"username"`
	Text     string `json:"text"`
}

type payload struct {
	Slack
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}

func NewSlack(m map[string]string) *Slack {
	s := &Slack{}
	s.Channel = m["channel"]
	s.URL = m["url"]
	s.IconURL = m["icon_url"]
	s.UserName = m["username"]
	return s
}

func (s Slack) Notify(t *task.Task) error {
	if s.Channel[0] != '#' {
		s.Channel = "#" + s.Channel
	}

	var message string
	var color string

	switch t.Status {
	case task.Succeeded:
		message = fmt.Sprintf("[%s] Succeeded", t.Name)
		color = "good"
	case task.Failed:
		message = fmt.Sprintf("[%s] Failed", t.Name)
		color = "danger"
	case task.Skipped:
		message = fmt.Sprintf("[%s] Skipped", t.Name)
		color = "warning"
	case task.Aborted:
		message = fmt.Sprintf("[%s] Aborted", t.Name)
		color = "warning"
	}

	a := attachment{
		Text:  message,
		Color: color,
	}

	p := payload{Slack: s, Attachments: []attachment{a}}
	j, _ := json.Marshal(p)
	buf := bytes.NewBuffer(j)

	log.Infof("[%s] Notify to Slack", t.Name)
	resp, err := http.Post(s.URL, "application/json", buf)
	if err != nil {
		e := fmt.Sprintf("[%s] Failed to notify to Slack: %s", t.Name, message)
		log.Errorf(e)
		return errors.New(e)
	}

	defer resp.Body.Close()

	return nil
}
