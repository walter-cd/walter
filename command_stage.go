package plumber

import (
        "log"
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
                log.Fatal("failed to retrieve pipe. %s", err)
                return false
        }
 
	err = cmd.Start()
        if err != nil {
                log.Fatal("failed to execute external command. %s", err)
                return false
        }
	return true;
}

func NewCommandStage() *CommandStage {
	return &CommandStage{}
}
