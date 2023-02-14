package main

import (

  "testing"

  github "app/github"


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

