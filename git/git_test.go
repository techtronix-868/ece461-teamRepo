package git

import (
	"testing"
)

func TestClone(t *testing.T) {
	url := "https://github.com/rattle99/QtNotepad"
	if !Clone(url) {
		t.Errorf("failed to find a license in %s", url)
	}

}
