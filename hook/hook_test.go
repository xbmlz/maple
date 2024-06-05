package hook

import "testing"

func TestHookAddAndPreAdd(t *testing.T) {
	h := Hook[int]{}

	if total := len(h.handlers); total != 0 {
		t.Fatalf("Expected no handlers, found %d", total)
	}

	triggerSequence := ""

	f1 := func(data int) error { triggerSequence += "f1"; return nil }
	f2 := func(data int) error { triggerSequence += "f2"; return nil }
	f3 := func(data int) error { triggerSequence += "f3"; return nil }
	f4 := func(data int) error { triggerSequence += "f4"; return nil }

	h.Add(f1)
	h.Add(f2)
	h.PreAdd(f3)
	h.PreAdd(f4)
	h.Trigger(1)

	if total := len(h.handlers); total != 4 {
		t.Fatalf("Expected %d handlers, found %d", 4, total)
	}

	expectedTriggerSequence := "f4f3f1f2"

	if triggerSequence != expectedTriggerSequence {
		t.Fatalf("Expected trigger sequence %s, got %s", expectedTriggerSequence, triggerSequence)
	}
}
