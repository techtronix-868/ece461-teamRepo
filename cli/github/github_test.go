package github

import (

	"testing"

  )

func TestSumNumbersInList(t *testing.T) {
	expectedOutput := 415
	output := Get_com("nullivex","nodist")

	if expectedOutput != output {
		t.Errorf("Failed ! got %c want %v", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}
	
}


func TestReleases(t *testing.T) {
	expectedOutput := 21

	output := Get_releases("nullivex","nodist")

	if expectedOutput != output {
		t.Errorf("Failed ! got %v want %c", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}
}

func TestScoreResponsivness(t *testing.T){
	expectedOutput := "Between 0 and 1"

	output := ScoreResponsiveness("nullivex","nodist")

	if output > 1 && output < 0 {
		t.Errorf("Failed ! got %v want %s", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}

}

func TestScoreBusFactor(t *testing.T){
	expectedOutput := "Between 0 and 1"

	output := ScoreBusFactor("nullivex","nodist")

	if output > 1 && output < 0  {
		t.Errorf("Failed ! got %v want %s", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}

}

func TestScore(t *testing.T){
	expectedScores := 0.35
	output:=Score("https://github.com/nullivex/nodist")

	if  output.Overall_score == expectedScores{
		t.Logf("Success !")
	}else {
		t.Errorf("Failed ! ")
	}

}