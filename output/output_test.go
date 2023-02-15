package output

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFormattedOutput(t *testing.T) {
	// Create test data
	data := []*NdJson{
		{
			URL:                   "github.com/Davies",
			Overall_score:         93.45,
			Ramp_up_score:         95.32,
			Bus_factor:            87.23,
			Responsiveness:        96.75,
			License_compatability: 92.17,
			Correctness:           90.11,
		},
		{
			URL:                   "github.com/Kahaan",
			Overall_score:         88.23,
			Ramp_up_score:         90.34,
			Bus_factor:            82.17,
			Responsiveness:        91.12,
			License_compatability: 89.13,
			Correctness:           85.65,
		},
		{
			URL:                   "github.com/Ishaan",
			Overall_score:         95.67,
			Ramp_up_score:         97.42,
			Bus_factor:            93.12,
			Responsiveness:        98.13,
			License_compatability: 94.76,
			Correctness:           92.89,
		},
	}

	expectedOutput := "{\"URL\":\"github.com/Ishaan\",\"Overall_score\":95.67,\"Ramp_up_score\":97.42,\"Bus_factor\":93.12,\"Responsiveness\":98.13,\"License_compatability\":94.76,\"Correctness\":92.89}\n{\"URL\":\"github.com/Davies\",\"Overall_score\":93.45,\"Ramp_up_score\":95.32,\"Bus_factor\":87.23,\"Responsiveness\":96.75,\"License_compatability\":92.17,\"Correctness\":90.11}\n{\"URL\":\"github.com/Kahaan\",\"Overall_score\":88.23,\"Ramp_up_score\":90.34,\"Bus_factor\":82.17,\"Responsiveness\":91.12,\"License_compatability\":89.13,\"Correctness\":85.65}"

	// Call the function being tested
	output := FormattedOutput(data)

	// Check the result
	assert.Equal(t, expectedOutput, output)
}

func TestNdJson_DataToNd(t *testing.T) {
    // Create an instance of the NdJson struct
    nd := &NdJson{}

    // Call the DataToNd method
    nd.DataToNd("test-package", 1.2345, 2.3456, 3.4567, 4.5678, 5.6789, 6.7890)

    // Check if the struct was populated with the correct values
    if nd.URL != "test-package" {
        t.Errorf("Expected URL to be 'test-package', got %v", nd.URL)
    }
    if nd.Overall_score != 1.23 {
        t.Errorf("Expected Overall_score to be 1.23, got %v", nd.Overall_score)
    }
    if nd.Ramp_up_score != 2.35 {
        t.Errorf("Expected Ramp_up_score to be 2.35, got %v", nd.Ramp_up_score)
    }
    if nd.Bus_factor != 3.46 {
        t.Errorf("Expected Bus_factor to be 3.46, got %v", nd.Bus_factor)
    }
    if nd.Responsiveness != 4.57 {
        t.Errorf("Expected Responsiveness to be 4.57, got %v", nd.Responsiveness)
    }
    if nd.License_compatability != 6.79 {
        t.Errorf("Expected License_compatability to be 6.79, got %v", nd.License_compatability)
    }
    if nd.Correctness != 5.68 {
        t.Errorf("Expected Correctness to be 5.68, got %v", nd.Correctness)
    }
}
