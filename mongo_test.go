package main

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestNewClient(t *testing.T) {
	// Call the newClient function to create a new client instance
	client := newClient()

	// List the available databases to check if the connection is working
	databases, err := client.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		t.Errorf("Failed to list databases: %v", err)
		return
	}

	// Verify that at least one database is returned
	if len(databases) == 0 {
		t.Errorf("Expected at least one database, but got none")
		return
	}

	// Check if the target string "testing" is present in the databases array
	found := false
	for _, db := range databases {
		if db == "testing" {
			found = true
			break
		}
	}

	// Verify that the target string was found
	if !found {
		t.Errorf("Expected to find the 'testing' database, but it was not present. Available databases: %v", databases)
		return
	}

	// Test passed
	t.Logf("Successfully connected to the database. Available databases: %v", databases)
}