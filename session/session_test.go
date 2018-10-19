package session

import (
	"os"
	"regexp"
	"testing"
)

func TestSessionID(t *testing.T) {
	s := NewSession()
	if len(s.ID) != 32 {
		t.Errorf("Session.ID '%s' length %d != 32\n", s.ID, len(s.ID))
	}

	re := regexp.MustCompile("[0-9a-f]{32}")
	if !re.MatchString(s.ID) {
		t.Errorf("Session.ID '%s' contains non-hex characters\n", s.ID)
	}
}

func TestSaveLoad(t *testing.T) {
	sessions := []Session{
		Session{"123", 10},
		Session{"456", -1},
		Session{"abc", -1},
		Session{"abc", 999999},
	}

	for _, s := range sessions {
		s.Save()
		filename := getFilename(s.ID)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Expected '%s' to exist", filename)
		}
		os.Remove(filename)
	}
}
