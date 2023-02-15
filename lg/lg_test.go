package lg

import (

  	"testing"
  	"os"
  	"io/ioutil"
	"log"
	"strings"

)

func TestLogger(t *testing.T) {
	f, err := ioutil.TempFile("", "testlog")
	if err != nil {
		t.Fatalf("Failed to create temporary log file: %v", err)
	}
	defer os.Remove(f.Name())

	Init(f.Name())

	InfoLogger.Print("test info log")
	WarningLogger.Print("test warning log")
	ErrorLogger.Print("test error log")


	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(b), "test info log") {
		t.Errorf("Expected 'test info log' in log file, got: %s", string(b))
	}
	if !strings.Contains(string(b), "test warning log") {
		t.Errorf("Expected 'test warning log' in log file, got: %s", string(b))
	}
	if !strings.Contains(string(b), "test error log") {
		t.Errorf("Expected 'test error log' in log file, got: %s", string(b))
	}

	// Close the logger
	log.SetOutput(os.Stderr) // reset the log output
	InfoLogger = log.New(os.Stderr, "", 0)
	WarningLogger = log.New(os.Stderr, "", 0)
	ErrorLogger = log.New(os.Stderr, "", 0)
}

func TestInitErr(t *testing.T) {
	err := os.Remove("test.log")
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove test.log: %v", err)
	}
	err = Init("test.log")
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}
