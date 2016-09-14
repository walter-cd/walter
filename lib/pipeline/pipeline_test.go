package pipeline

import "testing"

func TestLoad(t *testing.T) {
	yaml := `
pipeline:
  - name: command_stage_1
    command: echo "hello, world"
  - name: command_stage_2
    command: echo "hello, world, command_stage_2"
  - name: command_stage_3
    command: echo "hello, world, command_stage_3"
`
	_, err := Load(yaml)

	if err != nil {
		t.Fatal(err)
	}

}
