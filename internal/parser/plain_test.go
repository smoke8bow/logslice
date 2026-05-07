package parser

import (
	"testing"
	"time"
)

func TestParsePlain_NoTimestamp(t *testing.T) {
	ll := parsePlain("this is a plain log message")
	if ll.Message != "this is a plain log message" {
		t.Errorf("unexpected message: %q", ll.Message)
	}
	if !ll.Timestamp.IsZero() {
		t.Errorf("expected zero timestamp, got %v", ll.Timestamp)
	}
	if ll.Fields["msg"] != ll.Message {
		t.Errorf("fields[msg] mismatch")
	}
}

func TestParsePlain_RFC3339Timestamp(t *testing.T) {
	line := "2024-03-15T12:00:00Z application started successfully"
	ll := parsePlain(line)

	want, _ := time.Parse(time.RFC3339, "2024-03-15T12:00:00Z")
	if !ll.Timestamp.Equal(want) {
		t.Errorf("expected %v, got %v", want, ll.Timestamp)
	}
	if ll.Message != "application started successfully" {
		t.Errorf("unexpected message: %q", ll.Message)
	}
}

func TestParsePlain_SpaceSeparatedTimestamp(t *testing.T) {
	line := "2024-03-15 12:00:05 some event occurred"
	ll := parsePlain(line)

	want, _ := time.Parse("2006-01-02 15:04:05", "2024-03-15 12:00:05")
	if !ll.Timestamp.Equal(want) {
		t.Errorf("expected %v, got %v", want, ll.Timestamp)
	}
	if ll.Message != "some event occurred" {
		t.Errorf("unexpected message: %q", ll.Message)
	}
}

func TestParsePlain_EmptyLine(t *testing.T) {
	ll := parsePlain("")
	if ll.Message != "" {
		t.Errorf("expected empty message, got %q", ll.Message)
	}
	if !ll.Timestamp.IsZero() {
		t.Errorf("expected zero timestamp")
	}
}

func TestParsePlain_RawPreserved(t *testing.T) {
	original := "  2024-03-15T08:00:00Z   padded message  "
	ll := parsePlain(original)
	if ll.Raw != original {
		t.Errorf("Raw field modified: got %q", ll.Raw)
	}
}

func TestTryParseTime_ValidLayouts(t *testing.T) {
	cases := []string{
		"2024-01-02T15:04:05Z",
		"2024-01-02T15:04:05.999Z",
		"2024-01-02T15:04:05",
		"2024-01-02 15:04:05",
	}
	for _, c := range cases {
		if _, ok := tryParseTime(c); !ok {
			t.Errorf("expected %q to parse successfully", c)
		}
	}
}

func TestTryParseTime_Invalid(t *testing.T) {
	if _, ok := tryParseTime("not-a-timestamp"); ok {
		t.Error("expected parse failure for non-timestamp string")
	}
}
