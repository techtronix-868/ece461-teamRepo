package main

import (
	"testing"
	nd "app/output"
	log "app/lg"
	"os"
)

func TestSeperateLinks(t *testing.T) {

	log.Init(os.Getenv("LOG_FILE"))
	links := []string{"https://www.npmjs.com/package/browserify", "https://www.npmjs.com/package/express", "https://api.npms.io/v2/package/ws"}
	
	expectedScores := []*nd.NdJson{&nd.NdJson{}, &nd.NdJson{}, &nd.NdJson{}}
	
	scores := seperateLinks(links)

	if len(scores) != len(expectedScores) {
		t.Errorf("Expected length %d, got %d", len(expectedScores), len(scores))
	}
	for i := range scores {
		if scores[i] == nil {
			t.Errorf("Expected non-nil value at index %d, got nil", i)
		}
	}
}
