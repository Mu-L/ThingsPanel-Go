package main

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"
)

func TestSampleMatchesUSRPT170Protocol(t *testing.T) {
	startedAt := time.Unix(0, 0)
	random := rand.New(rand.NewSource(1))
	limits := alarmLimits{temperatureMin: 0, temperatureMax: 35, humidityMin: 20, humidityMax: 85, illuminanceMax: 50000}

	got := sample(startedAt, startedAt, defaultSN, limits, random)
	if got.LocalTime != startedAt.Format("2006-01-02,15:04:05") || got.SN != defaultSN {
		t.Fatalf("unexpected identity fields: %+v", got)
	}
	payload, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]interface{}
	if err := json.Unmarshal(payload, &fields); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"temperature", "humidity", "illuminance"} {
		if _, ok := fields[key].(float64); !ok {
			t.Errorf("field %s must be a JSON number: %s", key, payload)
		}
	}
	if _, exists := fields["Alarm"]; exists {
		t.Errorf("normal telemetry must omit Alarm: %s", payload)
	}
}

func TestAlarmIsSetWhenAnyValueExceedsLimit(t *testing.T) {
	limits := alarmLimits{temperatureMin: 10, temperatureMax: 20, humidityMin: 20, humidityMax: 85, illuminanceMax: 50000}
	if !isAlarm(25, 55, 10000, limits) {
		t.Fatal("expected temperature alarm")
	}
	if isAlarm(15, 55, 10000, limits) {
		t.Fatal("did not expect alarm for values within limits")
	}
	got := sample(time.Unix(0, 0), time.Unix(0, 0), defaultSN, limits, rand.New(rand.NewSource(1)))
	if got.Alarm != 1 {
		t.Fatalf("Alarm must be numeric 1, got %v", got.Alarm)
	}
}

func TestNormalizeBroker(t *testing.T) {
	tests := map[string]string{
		"47.92.253.145:1883":          "tcp://47.92.253.145:1883",
		"tcp://localhost:1883":        "tcp://localhost:1883",
		"ssl://mqtt.example.com:8883": "ssl://mqtt.example.com:8883",
	}
	for input, want := range tests {
		if got := normalizeBroker(input); got != want {
			t.Errorf("normalizeBroker(%q) = %q, want %q", input, got, want)
		}
	}
}
