package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type trackingResponse struct {
	Status string `json:"status"`
}

func TestTrackingAPI_Success(t *testing.T) {

	jsonData := []byte(`{
		"shipment_id":1,
		"tracking_number":"LOG-1-123456",
		"status":"IN_TRANSIT",
		"location":"Jakarta",
		"note":"Paket sedang dikirim"
	}`)

	resp, err := http.Post(
		"http://localhost:8087/tracking",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result trackingResponse

	json.NewDecoder(resp.Body).Decode(&result)

	if result.Status != "IN_TRANSIT" {
		t.Errorf("Expected IN_TRANSIT, got %s", result.Status)
	}
}