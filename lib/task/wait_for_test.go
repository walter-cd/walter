package task

import "testing"

func TestValidateWaitFor(t *testing.T) {
	w := &WaitFor{}
	w.Port = 80
	w.File = "/tmp"
	err := w.validate()
	if err == nil {
		t.Fatalf("Error should be returned: %#v", w)
	}

	w = &WaitFor{}
	w.Port = 80
	w.Delay = 10
	err = w.validate()
	if err == nil {
		t.Fatalf("Error should be returned: %#v", w)
	}
}
