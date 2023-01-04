package main

import (
	"fmt"
	"testing"
	"time"
)

func TestEmptyStats(t *testing.T) {
	var data RadarData
	data.Timestamps = make([]string, 0)
	data.Values = make([]string, 0)

	res, err := calculateStats(data)

	if err != nil {
		t.Error("calculateStats returned an error on empty data set")
	}

	if res.Mean == "" {
		t.Error("Empty Mean")
	}
	if res.Median == "" {
		t.Error("Empty Median")
	}
	if res.Max == "" {
		t.Error("Empty Max")
	}
	if res.Min == "" {
		t.Error("Empty Min")
	}

}

func TestBasicStats(t *testing.T) {
	var data RadarData
	data.Timestamps = make([]string, 0)
	data.Values = make([]string, 0)
	data.Timestamps = append(data.Timestamps, fmt.Sprintf("%d", time.Now().Unix()))
	data.Values = append(data.Values, "1.5")

	res, err := calculateStats(data)

	if err != nil {
		t.Error("calculateStats returned an error on 1 item data set")
	}

	if res.Mean != "1.500" {
		t.Errorf("Invalid Mean - Expected 1.5, Got %s", res.Mean)
	}

	if res.Median != "1.500" {
		t.Errorf("Invalid Median - Expected 1.5, Got %s", res.Median)
	}

	if res.Min != "1.500" {
		t.Errorf("Invalid Min - Expected 1.5, Got %s", res.Min)
	}

	if res.Max != "1.500" {
		t.Errorf("Invalid Max - Expected 1.5, Got %s", res.Max)
	}

}

func TestPrepareRequest(t *testing.T) {

	timestamp := time.Now().Unix()
	_, err := prepareRadarRequest(fmt.Sprintf("%d", timestamp))
	if err != nil {
		t.Errorf("Expected no error in prepareRequest, got: %s", err)
	}

	timestamp = time.Now().Unix() + 10
	_, err = prepareRadarRequest(fmt.Sprintf("%d", timestamp))
	if err == nil {
		t.Errorf("Expected error in prepareRequest, got none")
	}
}

func TestGetRadarData(t *testing.T) {
	radarEndpoint := "https://cfisysapi.developers.workers.dev"
	timestamp := time.Now().Unix()
	url := fmt.Sprintf("%s/stats?timestamp=%d", radarEndpoint, timestamp)

	stats, err := getRadarData(url)

	if err != nil {
		t.Errorf("Unexpected error in getRadarData: %s", err)
	}

	if stats.Mean == "" {
		t.Errorf("getRadarData mean returned empty string")
	}

	if stats.Median == "" {
		t.Errorf("getRadarData Median returned empty string")
	}

	if stats.Max == "" {
		t.Errorf("getRadarData Max returned empty string")
	}

	if stats.Min == "" {
		t.Errorf("getRadarData Min returned empty string")
	}
}
