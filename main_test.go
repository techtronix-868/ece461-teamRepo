package main

import (

  "testing"

  github "app/github"
  lg "app/lg"
  "os"
  "os/exec"
  nd "app/output"
  "io/ioutil"
  //"github.com/stretchr/testify/require"
  //"fmt"



)

func TestSumNumbersInList(t *testing.T) {
	expectedOutput := 415

	output := github.Get_com("nullivex","nodist")

	if expectedOutput != output {
		t.Errorf("Failed ! got %v want %c", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}
}


func TestReleases(t *testing.T) {
	expectedOutput := 21

	output := github.Get_releases("nullivex","nodist")

	if expectedOutput != output {
		t.Errorf("Failed ! got %v want %c", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}
}

func TestScoreResponsivness(t *testing.T){
	expectedOutput := "Between 0 and 1"

	output := github.ScoreResponsiveness("nullivex","nodist")

	if output > 1 && output < 0 {
		t.Errorf("Failed ! got %v want %s", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}

}

func TestScoreBusFactor(t *testing.T){
	expectedOutput := "Between 0 and 1"

	output := github.ScoreBusFactor("nullivex","nodist")

	if output > 1 && output < 0  {
		t.Errorf("Failed ! got %v want %s", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}

}

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

func TestMainFunction(t *testing.T) {
	// Set up the input file

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	inputFilePath := "test_input.txt"
	inputFile, err := os.Create(inputFilePath)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}
	defer os.Remove(inputFilePath)
	if _, err := inputFile.WriteString("https://www.npmjs.com/package/express\n"); err != nil {
		t.Fatalf("Failed to write to input file: %v", err)
	}
	if err := inputFile.Close(); err != nil {
		t.Fatalf("Failed to close input file: %v", err)
	}

	args := []string{"cmd", inputFile.Name()}
	main()
    // if err := run(args); err != nil {
    //     t.Fatalf("failed to run main function: %v", err)
    // }
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	t.Fatalf("Failed to run command: %v", err)
	// }

	// Check the output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.String() != "You provided the argument:\n" {
		t.Errorf("Unexpected output: %s", buf.String())
	}
}


