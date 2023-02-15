package main

import (

  "testing"

  lg "app/lg"
  "os"
//   "os/exec"
  nd "app/output"
  "io/ioutil"
//   "io"
  //"github.com/stretchr/testify/require"
//   "bytes"



)



func TestSeperateLinks1(t *testing.T) {

	lg.Init(os.Getenv("LOG_FILE"))
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

func TestSeperateLinks2(t *testing.T) {

	lg.Init(os.Getenv("LOG_FILE"))
	links := []string{"https://github.com/nullivex/nodist", "https://github.com/lodash/lodash", "https://github.com/cloudinary/cloudinary_npm"}
	

	
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

func TestReadInput1(t *testing.T) {
    // Create a temporary input file for testing
    inputContent := "Line 1\nLine 2\nLine 3"
    tmpfile, err := ioutil.TempFile("", "test_input")
    if err != nil {
        t.Errorf("Failed to create temporary input file: %s", err)
        return
    }
    defer os.Remove(tmpfile.Name()) // clean up

    if _, err := tmpfile.Write([]byte(inputContent)); err != nil {
        t.Errorf("Failed to write to temporary input file: %s", err)
        return
    }

    expectedLines := []string{"Line 1", "Line 2", "Line 3"}
    resultLines := readInput(tmpfile.Name())

    if len(resultLines) != len(expectedLines) {
        t.Errorf("Expected %d lines, but got %d", len(expectedLines), len(resultLines))
        return
    }

    for i, line := range resultLines {
        if line != expectedLines[i] {
            t.Errorf("Line %d does not match. Expected '%s', but got '%s'", i+1, expectedLines[i], line)
            return
        }
    }
}


func TestReadInput2(t *testing.T) {
    // Create a temporary input file for testing

    resultLines := readInput("xxxx")

    if resultLines != nil {
        t.Errorf("Expected nil , got %s", resultLines)
        return
    }

}

// func TestMainFunction(t *testing.T) {
// 	inputFile, err := ioutil.TempFile("", "input")
// 	if err != nil {
// 		t.Fatalf("Failed to create input file: %v", err)
// 	}
// 	defer os.Remove(inputFile.Name())


// 	inputData := "https://www.npmjs.com/package/express\n"
// 	if _, err := inputFile.WriteString(inputData); err != nil {
// 		t.Fatalf("Failed to write to input file: %v", err)
// 	}
// 	if err := inputFile.Close(); err != nil {
// 		t.Fatalf("Failed to close input file: %v", err)
// 	}

// 	old := os.Stdout
// 	r, w, _ := os.Pipe()
// 	os.Stdout = w
// 	defer func() {
// 		os.Stdout = old
// 	}()
// 	os.Args = []string{"cmd", inputFile.Name()}
// 	main()

// 	// Compare the output with the expected value
// 	expectedOutput := "https://example.com\t1\n"
// 	w.Close()
// 	var buf bytes.Buffer
// 	io.Copy(&buf, r)
// 	actualOutput := buf.String()
// 	if actualOutput != expectedOutput {
// 		t.Errorf("Output does not match expected value. Got %s, expected %s", actualOutput, expectedOutput)
// 	}
// }


