package naming

import "testing"

func TestNormalizeIdentifier(t *testing.T) {
	cases := map[string]string{
		"blink":     "blink",
		"my-task":   "my_task",
		"my_task":   "my_task",
		"my task":   "my_task",
		"1task":     "_1task",
		"Task.Name": "Task_Name",
	}
	for input, want := range cases {
		if got := NormalizeIdentifier(input); got != want {
			t.Fatalf("NormalizeIdentifier(%q) = %q, want %q", input, got, want)
		}
	}
}
