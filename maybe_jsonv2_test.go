//go:build goexperiment.jsonv2

package maybe

import (
	"encoding/json/v2"
	"testing"
)

// Test struct with scalar Maybe fields
type TestConfigScalar struct {
	Name    string      `json:"name"`
	Timeout Maybe[int]  `json:"timeout,omitzero"`
	Enabled Maybe[bool] `json:"enabled,omitzero"`
}

// Test struct with nested object Maybe field
type TestAddress struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

type TestConfigObject struct {
	Name    string             `json:"name"`
	Address Maybe[TestAddress] `json:"address,omitzero"`
}

func TestUnmarshalOmittedScalar(t *testing.T) {
	// JSON with timeout field missing
	jsonData := []byte(`{"name": "test"}`)

	var cfg TestConfigScalar
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if cfg.Timeout.IsSome() {
		t.Error("Expected Timeout to be None when field is omitted")
	}
	if cfg.Enabled.IsSome() {
		t.Error("Expected Enabled to be None when field is omitted")
	}
}

func TestUnmarshalOmittedObject(t *testing.T) {
	// JSON with address field missing
	jsonData := []byte(`{"name": "test"}`)

	var cfg TestConfigObject
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if cfg.Address.IsSome() {
		t.Error("Expected Address to be None when field is omitted")
	}
}

func TestUnmarshalPresentScalar(t *testing.T) {
	jsonData := []byte(`{"name": "test", "timeout": 30, "enabled": true}`)

	var cfg TestConfigScalar
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if cfg.Timeout.IsNone() {
		t.Error("Expected Timeout to be Some when field is present")
	}
	if val, _ := cfg.Timeout.Unwrap(); val != 30 {
		t.Errorf("Expected Timeout value 30, got %d", val)
	}

	if cfg.Enabled.IsNone() {
		t.Error("Expected Enabled to be Some when field is present")
	}
	if val, _ := cfg.Enabled.Unwrap(); val != true {
		t.Errorf("Expected Enabled value true, got %v", val)
	}
}

func TestUnmarshalPresentObject(t *testing.T) {
	jsonData := []byte(`{"name": "test", "address": {"street": "123 Main", "city": "NYC"}}`)

	var cfg TestConfigObject
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if cfg.Address.IsNone() {
		t.Error("Expected Address to be Some when field is present")
	}
	addr, _ := cfg.Address.Unwrap()
	if addr.Street != "123 Main" || addr.City != "NYC" {
		t.Errorf("Expected address {123 Main, NYC}, got {%s, %s}", addr.Street, addr.City)
	}
}

func TestMarshalSomeValue(t *testing.T) {
	cfg := TestConfigScalar{
		Name:    "test",
		Timeout: Some(42),
		Enabled: Some(true),
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := `{"name":"test","timeout":42,"enabled":true}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestMarshalNoneWithOmitzero(t *testing.T) {
	cfg := TestConfigScalar{
		Name:    "test",
		Timeout: None[int](),
		Enabled: None[bool](),
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// With omitzero, None fields should be omitted
	expected := `{"name":"test"}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestMarshalNoneObject(t *testing.T) {
	cfg := TestConfigObject{
		Name:    "test",
		Address: None[TestAddress](),
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := `{"name":"test"}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestUnmarshalNullValue(t *testing.T) {
	// JSON with explicit null value
	jsonData := []byte(`{"name": "test", "timeout": null}`)

	var cfg TestConfigScalar
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Field was present (as null), so it should be Some with zero value
	if cfg.Timeout.IsNone() {
		t.Error("Expected Timeout to be Some when field is present (even as null)")
	}
}
