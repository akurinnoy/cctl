package main

import (
	"os"
	"testing"
)

func TestGetSlug_Basename(t *testing.T) {
	slug := GetSlug("/tmp/my-project")
	if slug != "my-project" {
		t.Errorf("expected 'my-project', got '%s'", slug)
	}
}

func TestGetSlug_Home(t *testing.T) {
	home, _ := os.UserHomeDir()
	slug := GetSlug(home)
	if slug != "[home]" {
		t.Errorf("expected '[home]', got '%s'", slug)
	}
}

func TestGetSlug_Root(t *testing.T) {
	slug := GetSlug("/")
	if slug != "[root]" {
		t.Errorf("expected '[root]', got '%s'", slug)
	}
}
