package dedupe

import (
	"testing"
	"time"
)

func TestIsDuplicate_FirstOccurrenceNotDuplicate(t *testing.T) {
	d := New(5 * time.Second)
	e := Entry{"message": "hello", "level": "info"}
	if d.IsDuplicate(e) {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_SecondOccurrenceIsDuplicate(t *testing.T) {
	d := New(5 * time.Second)
	e := Entry{"message": "hello", "level": "info"}
	d.IsDuplicate(e)
	if !d.IsDuplicate(e) {
		t.Fatal("second occurrence within window should be a duplicate")
	}
}

func TestIsDuplicate_DifferentEntriesNotDuplicate(t *testing.T) {
	d := New(5 * time.Second)
	e1 := Entry{"message": "hello"}
	e2 := Entry{"message": "world"}
	d.IsDuplicate(e1)
	if d.IsDuplicate(e2) {
		t.Fatal("different entries should not be duplicates of each other")
	}
}

func TestIsDuplicate_AfterWindowExpires(t *testing.T) {
	now := time.Unix(1000, 0)
	d := New(2 * time.Second)
	d.now = func() time.Time { return now }

	e := Entry{"message": "expiring"}
	d.IsDuplicate(e)

	// Advance time beyond window
	d.now = func() time.Time { return now.Add(3 * time.Second) }
	if d.IsDuplicate(e) {
		t.Fatal("entry seen after window expiry should not be a duplicate")
	}
}

func TestIsDuplicate_EmptyEntry(t *testing.T) {
	d := New(5 * time.Second)
	e := Entry{}
	d.IsDuplicate(e)
	if !d.IsDuplicate(e) {
		t.Fatal("identical empty entries within window should be duplicates")
	}
}

func TestIsDuplicate_IndependentWindows(t *testing.T) {
	d := New(5 * time.Second)
	e1 := Entry{"src": "a", "msg": "x"}
	e2 := Entry{"src": "b", "msg": "x"}
	d.IsDuplicate(e1)
	if d.IsDuplicate(e2) {
		t.Fatal("entries from different sources should be independent")
	}
}
