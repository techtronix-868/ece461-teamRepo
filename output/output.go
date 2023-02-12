package output

import (
	"encoding/json"
	"fmt"
	"strings"
)

type NdJson struct {
	URL               string
	Overall_score         float64
	Ramp_up_score         float64
	Bus_factor            float64
	Responsiveness        float64
	License_compatability float64
	Correctness           float64
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
