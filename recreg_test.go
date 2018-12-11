package recreg

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}
}

func TestListActions(t *testing.T) {
	c := NewClient(nil)
	date, err := time.Parse(ISO8601Date, FirstDate)
	if err != nil {
		t.Error(err)
	}
	actions, err := c.ListActions(date)
	if err != nil {
		t.Error(err)
	}
	if got, want := len(actions), 4; got != want {
		t.Errorf("len(ListActions(%q)) is %d, want %d", FirstDate, got, want)
	}
}
