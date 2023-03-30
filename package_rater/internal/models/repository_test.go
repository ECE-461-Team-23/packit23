package models

import "testing"

func TestNewRepository(t *testing.T) {
	newRepo := NewRepository()
	if newRepo == nil {
		t.Error("Error creating repository!")
	}
}
