package models

import (
	"encoding/json"
	"testing"
)

func TestUserPublicJSONUsesFullName(t *testing.T) {
	user := UserPublic{
		FullName: "Test User",
		Email:    "test@example.com",
		Role:     RoleRegistrar,
		Language: LanguageAm,
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if _, ok := parsed["name"]; ok {
		t.Fatal("expected no name field in JSON")
	}
	if parsed["full_name"] != "Test User" {
		t.Fatalf("expected full_name, got %v", parsed["full_name"])
	}
}
