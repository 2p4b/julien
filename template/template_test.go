package template

import (
	"testing"
)

func TestTemplate_Get(t *testing.T) {
	tmpl := &Template{
		Meta: map[string]interface{}{
			"key": "value",
		},
	}

	// Test with existing key
	value := tmpl.Get("key")
	if value != "value" {
		t.Errorf("Expected value 'value', got '%v'", value)
	}

	// Test with non-existing key
	value = tmpl.Get("nonExistingKey")
	if value != nil {
		t.Errorf("Expected value nil, got '%v'", value)
	}

	// Test with default value
	value = tmpl.Get("nonExistingKey", "defaultValue")
	if value != "defaultValue" {
		t.Errorf("Expected value 'defaultValue', got '%v'", value)
	}
}

func TestTemplate_GetString(t *testing.T) {
	tmpl := &Template{
		Meta: map[string]interface{}{
			"key": "value",
		},
	}

	// Test with existing key
	value := tmpl.GetString("key")
	if value != "value" {
		t.Errorf("Expected value 'value', got '%v'", value)
	}

	// Test with non-existing key
	value = tmpl.GetString("nonExistingKey")
	if value != "" {
		t.Errorf("Expected value '', got '%v'", value)
	}

	// Test with default value
	value = tmpl.GetString("nonExistingKey", "defaultValue")
	if value != "defaultValue" {
		t.Errorf("Expected value 'defaultValue', got '%v'", value)
	}
}

func TestTemplate_GetInt(t *testing.T) {
	tmpl := &Template{
		Meta: map[string]interface{}{
			"key": 123,
		},
	}

	// Test with existing key
	value := tmpl.GetInt("key")
	if value != 123 {
		t.Errorf("Expected value 123, got %v", value)
	}

	// Test with non-existing key
	value = tmpl.GetInt("nonExistingKey")
	if value != 0 {
		t.Errorf("Expected value 0, got %v", value)
	}

	// Test with default value
	value = tmpl.GetInt("nonExistingKey", 456)
	if value != 456 {
		t.Errorf("Expected value 456, got %v", value)
	}
}

func TestTemplate_GetInt64(t *testing.T) {
	tmpl := &Template{
		Meta: map[string]interface{}{
			"key": int64(123),
		},
	}

	// Test with existing key
	value := tmpl.GetInt64("key")
	if value != 123 {
		t.Errorf("Expected value 123, got %v", value)
	}

	// Test with non-existing key
	value = tmpl.GetInt64("nonExistingKey")
	if value != 0 {
		t.Errorf("Expected value 0, got %v", value)
	}

	// Test with default value
	value = tmpl.GetInt64("nonExistingKey", 456)
	if value != 456 {
		t.Errorf("Expected value 456, got %v", value)
	}
}
