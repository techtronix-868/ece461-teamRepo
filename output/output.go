package output

import (
	"encoding/json"
	"fmt"
	"strings"
	"sort"
)

type NdJson struct {
	URL                   string
	Overall_score         float64
	Ramp_up_score         float64
	Bus_factor            float64
	Responsiveness        float64
	License_compatability float64
	Correctness           float64
}

type ByOverallScore []*NdJson

func (a ByOverallScore) Len() int {
	return len(a)
}

func (a ByOverallScore) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByOverallScore) Less(i, j int) bool {
	// Change the sign to < for ascending.
	return a[i].Overall_score > a[j].Overall_score
}

func (nd *NdJson) DataToNd(_package string, Overall_score float64, _Ramp_up_score float64, _Bus_factor float64, _Responsiveness float64, _Correctness float64, _License_compatability float64) *NdJson {
	nd.URL = _package
	nd.Overall_score = Overall_score
	nd.Ramp_up_score = _Ramp_up_score
	nd.Bus_factor = _Bus_factor
	nd.Responsiveness = _Responsiveness
	nd.License_compatability = _License_compatability
	nd.Correctness = _Correctness

	return nd
}

func FormattedOutput(data []*NdJson) string {
	var jsonStrings []string

	// Sorting in the order for the best score at the top
	sort.Sort(ByOverallScore(data))

	// Loop over the slice of structs and convert each to JSON
	for _, record := range data {
		recordJSON, err := json.Marshal(record)
		if err != nil {
			fmt.Println("Error marshaling struct:", err)
			continue
		}
		jsonStrings = append(jsonStrings, string(recordJSON))
	}

	// Join the JSON strings with a newline character to produce the NDJSON output
	ndjson := strings.Join(jsonStrings, "\n")
	return ndjson
}
