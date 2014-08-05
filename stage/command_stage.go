package stage

import (
	"io"
	"os"
	"os/exec"
)

type CommandStage struct {
	command   string
	arguments []string
}

func (self *CommandStage) Run() bool {
	cmd := exec.Command(self.command)
	cmd.Args = append([]string{self.command}, self.arguments...)
	cmd.Dir = "."

	out, err := cmd.StdoutPipe()
	if err != nil {
		return false
	}

	err = cmd.Start()
	if err != nil {
		return false
	}
	cmd.Stdout = os.Stdout
	io.Copy(os.Stdout, out)
	err = cmd.Wait()
	if err != nil {
		return false
	}
	return true
}

func (self *CommandStage) AddCommand(command string, arguments ...string) {
	self.command = command
	self.arguments = arguments
}

func NewCommandStage() *CommandStage {

	return &CommandStage{}
}
