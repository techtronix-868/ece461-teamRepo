package github

import (
	"testing"
)

func TestSumNumbersInList(t *testing.T) {
	expectedOutput := 415
	output := Get_com("nullivex", "nodist")

	if expectedOutput != output {
		t.Errorf("Failed ! got %c want %d", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}

}

func TestReleases(t *testing.T) {
	expectedOutput := 21

	output := Get_releases("nullivex", "nodist")

	if expectedOutput != output {
		t.Errorf("Failed ! got %v want %c", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}
}

func TestScoreResponsivness(t *testing.T) {
	expectedOutput := "Between 0 and 1"

	output := scoreResponsiveness("nullivex", "nodist")

	if output > 1 && output < 0 {
		t.Errorf("Failed ! got %v want %s", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}

}

func TestScoreBusFactor(t *testing.T) {
	expectedOutput := "Between 0 and 1"

	output := scoreBusFactor("nullivex", "nodist")

	if output > 1 && output < 0 {
		t.Errorf("Failed ! got %v want %s", output, expectedOutput)
	} else {
		t.Logf("Success !")
	}

}

func TestScoreVersionPinning(t *testing.T) {
	depList := []string{"1.2.3", "~1.2.3", "^0.2.3", "~1", "^1.2.3"} //first 3 pinned, last 2 not
	output := scoreVersionPinning(depList)
	expectedOutput := .60
	if output == expectedOutput {
		t.Logf("ScoreVersionPinning Passed")
	} else {
		t.Logf("Failed! got %f want %f", output, expectedOutput)
	}
}

func TestScore(t *testing.T) {
	expectedScores := 0.35
	output := Score("https://github.com/nullivex/nodist")

	if output.NET_SCORE == expectedScores {
		t.Logf("Success !")
	} else {
		t.Errorf("Failed ! ")
	}

}
