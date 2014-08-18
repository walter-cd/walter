package config

import (
	"testing"

	"github.com/takahi-i/plumber/stages"
)

func TestParse(t *testing.T) {
	configData := ReadConfig("../tests/fixtures/pipeline.yml")
	actual := (*Parse(configData)).Stages.Front().Value.(*stages.CommandStage).Command

	expected := "echo \"hello, world\""
	t.Logf("got %v\nwant %v", actual, expected)
	//if expected != actual {
	//t.Errorf("got %v\nwant %v", actual, expected)
	//}
}
