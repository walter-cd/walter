package config

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	configData := *ReadConfig("../tests/fixtures/pipeline.yml")
	actual := configData["pipeline"].(map[interface{}]interface{})["command_stage_1"].(map[interface{}]interface{})["command"]

	expected := "echo \"hello, world\""
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
