package processor

import (
	"strings"
	"testing"
	"time"
)

func TestNewStats_InitializesStartTime(t *testing.T) {
	before := time.Now()
	s := NewStats()
	after := time.Now()

	if s.StartTime.Before(before) || s.StartTime.After(after) {
		t.Errorf("StartTime %v not in expected range [%v, %v]", s.StartTime, before, after)
	}
}

func TestStats_Counters(t *testing.T) {
	s := NewStats()
	s.LinesRead.Add(100)
	s.LinesMatched.Add(60)
	s.LinesDropped.Add(40)
	s.BytesRead.Add(4096)

	if s.LinesRead.Load() != 100 {
		t.Errorf("expected LinesRead=100, got %d", s.LinesRead.Load())
	}
	if s.LinesMatched.Load() != 60 {
		t.Errorf("expected LinesMatched=60, got %d", s.LinesMatched.Load())
	}
	if s.LinesDropped.Load() != 40 {
		t.Errorf("expected LinesDropped=40, got %d", s.LinesDropped.Load())
	}
	if s.BytesRead.Load() != 4096 {
		t.Errorf("expected BytesRead=4096, got %d", s.BytesRead.Load())
	}
}

func TestStats_Duration_BeforeFinish(t *testing.T) {
	s := NewStats()
	time.Sleep(10 * time.Millisecond)
	d := s.Duration()
	if d < 10*time.Millisecond {
		t.Errorf("expected duration >= 10ms, got %v", d)
	}
}

func TestStats_Duration_AfterFinish(t *testing.T) {
	s := NewStats()
	time.Sleep(10 * time.Millisecond)
	s.Finish()
	d1 := s.Duration()
	time.Sleep(20 * time.Millisecond)
	d2 := s.Duration()
	if d1 != d2 {
		t.Errorf("duration should be fixed after Finish: d1=%v d2=%v", d1, d2)
	}
}

func TestStats_Summary_ContainsFields(t *testing.T) {
	s := NewStats()
	s.LinesRead.Add(10)
	s.LinesMatched.Add(7)
	s.LinesDropped.Add(3)
	s.BytesRead.Add(512)
	s.Finish()

	summary := s.Summary()
	for _, want := range []string{"lines_read=10", "lines_matched=7", "lines_dropped=3", "bytes_read=512", "duration="} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary missing %q: %s", want, summary)
		}
	}
}
