package pipeline

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	yaml := `
build:
  tasks:
    - name: command_stage_1
      command: echo "hello, world"
    - name: command_stage_2
      command: echo "hello, world, command_stage_2"
    - name: command_stage_3
      command: echo "hello, world, command_stage_3"
`
	p, err := Load([]byte(yaml))

	fmt.Printf("%v", p)
	if err != nil {
		t.Fatal(err)
	}

}
