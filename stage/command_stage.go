package stage

import (
	"bytes"
	"io"
	"log"
	"os/exec"
)

type CommandStage struct {
	command   string
	arguments []string
	outResult string
}

func (self *CommandStage) GetStdoutResult() string {
	return self.outResult
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
	self.outResult = copyStream(out)
	err = cmd.Wait()
	if err != nil {
		return false
	}
	return true
}

func copyStream(reader io.Reader) string {
	var err error
	var n int
	var buffer bytes.Buffer
	tmpBuf := make([]byte, 1024)
	for {
		if n, err = reader.Read(tmpBuf); err != nil {
			break
		}
		buffer.Write(tmpBuf[0:n])
	}
	if err == io.EOF {
		err = nil
	} else {
		log.Println("ERROR: " + err.Error())
	}
	return buffer.String()
}

func (self *CommandStage) AddCommand(command string, arguments ...string) {
	self.command = command
	self.arguments = arguments
}

func NewCommandStage() *CommandStage {
	return &CommandStage{}
}
