package cli

import (
	"encoding/json"
	"testing"
)

func TestDatasetValue_UnmarshalJSON_SimpleString(t *testing.T) {
	jsonData := []byte(`"test-value"`)
	var dv DatasetValue

	err := json.Unmarshal(jsonData, &dv)
	if err != nil {
		t.Fatalf("Failed to unmarshal string value: %v", err)
	}

	if dv.JqQuery != "test-value" {
		t.Errorf("Expected JqQuery to be 'test-value', got '%s'", dv.JqQuery)
	}
}

func TestDatasetValue_UnmarshalJSON_SimpleNumber(t *testing.T) {
	jsonData := []byte(`2`)
	var dv DatasetValue

	err := json.Unmarshal(jsonData, &dv)
	if err != nil {
		t.Fatalf("Failed to unmarshal number value: %v", err)
	}

	if dv.JqQuery != "2" {
		t.Errorf("Expected JqQuery to be '2', got '%s'", dv.JqQuery)
	}
}

func TestDatasetValue_UnmarshalJSON_SimpleFloat(t *testing.T) {
	jsonData := []byte(`3.14`)
	var dv DatasetValue

	err := json.Unmarshal(jsonData, &dv)
	if err != nil {
		t.Fatalf("Failed to unmarshal float value: %v", err)
	}

	if dv.JqQuery != "3.14" {
		t.Errorf("Expected JqQuery to be '3.14', got '%s'", dv.JqQuery)
	}
}

func TestDatasetValue_UnmarshalJSON_SimpleBoolTrue(t *testing.T) {
	jsonData := []byte(`true`)
	var dv DatasetValue

	err := json.Unmarshal(jsonData, &dv)
	if err != nil {
		t.Fatalf("Failed to unmarshal bool value: %v", err)
	}

	if dv.JqQuery != "true" {
		t.Errorf("Expected JqQuery to be 'true', got '%s'", dv.JqQuery)
	}
}

func TestDatasetValue_UnmarshalJSON_SimpleBoolFalse(t *testing.T) {
	jsonData := []byte(`false`)
	var dv DatasetValue

	err := json.Unmarshal(jsonData, &dv)
	if err != nil {
		t.Fatalf("Failed to unmarshal bool value: %v", err)
	}

	if dv.JqQuery != "false" {
		t.Errorf("Expected JqQuery to be 'false', got '%s'", dv.JqQuery)
	}
}

func TestDatasetValue_UnmarshalJSON_JqQueryObject(t *testing.T) {
	jsonData := []byte(`{"jqQuery": ".user.email"}`)
	var dv DatasetValue

	err := json.Unmarshal(jsonData, &dv)
	if err != nil {
		t.Fatalf("Failed to unmarshal jqQuery object: %v", err)
	}

	if dv.JqQuery != ".user.email" {
		t.Errorf("Expected JqQuery to be '.user.email', got '%s'", dv.JqQuery)
	}
}

func TestDatasetValue_UnmarshalJSON_Null(t *testing.T) {
	jsonData := []byte(`null`)
	var dv DatasetValue

	err := json.Unmarshal(jsonData, &dv)
	if err != nil {
		t.Fatalf("Failed to unmarshal null value: %v", err)
	}

	if dv.JqQuery != "" {
		t.Errorf("Expected JqQuery to be empty, got '%s'", dv.JqQuery)
	}
}

func TestDatasetRule_UnmarshalJSON_WithSimpleValue(t *testing.T) {
	jsonData := []byte(`{
		"property": "cpuLimit",
		"operator": ">",
		"value": 2
	}`)

	var rule DatasetRule
	err := json.Unmarshal(jsonData, &rule)
	if err != nil {
		t.Fatalf("Failed to unmarshal dataset rule: %v", err)
	}

	if rule.Property == nil || *rule.Property != "cpuLimit" {
		t.Errorf("Expected property to be 'cpuLimit'")
	}

	if rule.Operator != ">" {
		t.Errorf("Expected operator to be '>', got '%s'", rule.Operator)
	}

	if rule.Value == nil {
		t.Fatal("Expected value to be non-nil")
	}

	if rule.Value.JqQuery != "2" {
		t.Errorf("Expected value.JqQuery to be '2', got '%s'", rule.Value.JqQuery)
	}
}

func TestDataset_UnmarshalJSON_CompleteExample(t *testing.T) {
	jsonData := []byte(`{
		"combinator": "and",
		"rules": [
			{
				"property": "cpuLimit",
				"operator": ">",
				"value": 2
			},
			{
				"property": "status",
				"operator": "=",
				"value": "active"
			},
			{
				"property": "$team",
				"operator": "containsAny",
				"value": {
					"jqQuery": ".user.team"
				}
			}
		]
	}`)

	var dataset Dataset
	err := json.Unmarshal(jsonData, &dataset)
	if err != nil {
		t.Fatalf("Failed to unmarshal dataset: %v", err)
	}

	if dataset.Combinator != "and" {
		t.Errorf("Expected combinator to be 'and', got '%s'", dataset.Combinator)
	}

	if len(dataset.Rules) != 3 {
		t.Fatalf("Expected 3 rules, got %d", len(dataset.Rules))
	}

	// Check first rule (number value)
	if dataset.Rules[0].Value.JqQuery != "2" {
		t.Errorf("Expected first rule value to be '2', got '%s'", dataset.Rules[0].Value.JqQuery)
	}

	// Check second rule (string value)
	if dataset.Rules[1].Value.JqQuery != "active" {
		t.Errorf("Expected second rule value to be 'active', got '%s'", dataset.Rules[1].Value.JqQuery)
	}

	// Check third rule (jqQuery object)
	if dataset.Rules[2].Value.JqQuery != ".user.team" {
		t.Errorf("Expected third rule value to be '.user.team', got '%s'", dataset.Rules[2].Value.JqQuery)
	}
}

func TestDatasetValue_MarshalJSON_Number(t *testing.T) {
	dv := DatasetValue{JqQuery: "2"}

	jsonData, err := json.Marshal(dv)
	if err != nil {
		t.Fatalf("Failed to marshal number value: %v", err)
	}

	expected := `2`
	if string(jsonData) != expected {
		t.Errorf("Expected marshaled value to be '%s', got '%s'", expected, string(jsonData))
	}
}

func TestDatasetValue_MarshalJSON_String(t *testing.T) {
	dv := DatasetValue{JqQuery: "active"}

	jsonData, err := json.Marshal(dv)
	if err != nil {
		t.Fatalf("Failed to marshal string value: %v", err)
	}

	expected := `"active"`
	if string(jsonData) != expected {
		t.Errorf("Expected marshaled value to be '%s', got '%s'", expected, string(jsonData))
	}
}

func TestDatasetValue_MarshalJSON_Bool(t *testing.T) {
	dv := DatasetValue{JqQuery: "true"}

	jsonData, err := json.Marshal(dv)
	if err != nil {
		t.Fatalf("Failed to marshal bool value: %v", err)
	}

	expected := `true`
	if string(jsonData) != expected {
		t.Errorf("Expected marshaled value to be '%s', got '%s'", expected, string(jsonData))
	}
}

func TestDatasetValue_MarshalJSON_Null(t *testing.T) {
	dv := DatasetValue{JqQuery: ""}

	jsonData, err := json.Marshal(dv)
	if err != nil {
		t.Fatalf("Failed to marshal null value: %v", err)
	}

	expected := `null`
	if string(jsonData) != expected {
		t.Errorf("Expected marshaled value to be '%s', got '%s'", expected, string(jsonData))
	}
}

func TestDataset_RoundTrip(t *testing.T) {
	original := []byte(`{
		"combinator": "and",
		"rules": [
			{
				"property": "cpuLimit",
				"operator": ">",
				"value": 2
			}
		]
	}`)

	var dataset Dataset
	err := json.Unmarshal(original, &dataset)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	marshaled, err := json.Marshal(dataset)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var dataset2 Dataset
	err = json.Unmarshal(marshaled, &dataset2)
	if err != nil {
		t.Fatalf("Failed to unmarshal round-trip: %v", err)
	}

	if dataset2.Rules[0].Value.JqQuery != "2" {
		t.Errorf("Round-trip failed: expected value '2', got '%s'", dataset2.Rules[0].Value.JqQuery)
	}
}
