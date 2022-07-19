package main_test

import (
	"os"
	"testing"
	"h4kor/kiss-social"
)

func TestCanCreateANewUser(t *testing.T) {
	// Create a new user
	main.CreateNewUser("/tmp/test", "testuser")
	if _, err := os.Stat("/tmp/test/testuser"); err != nil {
		t.Error("User directory not created")
	}
}

func TestCreateUserAddsVersionFile(t *testing.T) {
	// Create a new user
	main.CreateNewUser("/tmp/test", "testuser")
	if _, err := os.Stat("/tmp/test/testuser/meta/VERSION"); err != nil {
		t.Error("Version file not created")
	}
}

func TestCreateUserAddsBaseHtmlFile(t *testing.T) {
	// Create a new user
	main.CreateNewUser("/tmp/test", "testuser")
	if _, err := os.Stat("/tmp/test/testuser/meta/base.html"); err != nil {
		t.Error("Base html file not created")
	}
}