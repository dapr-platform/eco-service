package client

import (
	"testing"
)

func TestGetBoxesHourStats(t *testing.T) {
	// Test case 1: Valid parameters with hour
	params := map[string]string{
		"projectCode": "P00000001",
		"mac":         "98CC4D000000",
		"year":        "2024",
		"month":       "01",
		"day":         "01",
		"hour":        "12",
	}

	stats, err := GetBoxesHourStats(params)
	if err != nil {
		t.Errorf("GetBoxesHourStats failed: %v", err)
	}
	if len(stats) == 0 {
		t.Error("Expected non-empty stats response")
	}
	projectCode, err := GetBoxProjectCode("98CC4D150A00")
	if err != nil {
		t.Errorf("GetBoxProjectCode failed: %v", err)
	}
	// Test case 2: Valid parameters without hour (returns full day data)
	params = map[string]string{
		"projectCode": projectCode,
		"mac":         "98CC4D150A00",
		"year":        "2024",
		"month":       "11",
		"day":         "01",
	}

	stats, err = GetBoxesHourStats(params)
	if err != nil {
		t.Errorf("GetBoxesHourStats failed: %v", err)
	}
	if len(stats) == 0 {
		t.Error("Expected non-empty stats response")
	}
}
