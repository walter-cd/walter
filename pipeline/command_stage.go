package plumber

import (
        "os/exec"
)

type CommandStage struct {
	command string
	arguments []string
}

func (self *CommandStage) Run() bool {
	cmd := exec.Command(self.command)
        cmd.Args = self.arguments
        cmd.Dir = "."

        _, err := cmd.StdoutPipe()
        if err != nil {
                return false
        }
 
	err = cmd.Start()
        if err != nil {
                return false
        }
	return true;
}

func (self *CommandStage) AddCommand(command string, arguments []string)  {
	self.command = command
	self.arguments = arguments
}

func NewCommandStage() *CommandStage {
	return &CommandStage{}
}
