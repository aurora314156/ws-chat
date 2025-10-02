package db_test

import (
	"testing"
	"ws-chat/db"

	"github.com/joho/godotenv"
)

func TestSupabaseConnection(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log("Warning: Could not load .env file. Assuming environment variables are already set.")
	}

	supaClient, err := db.NewSupabaseClient()

	if err != nil || supaClient == nil {
		t.Fatalf("Failed to initialize Supabase client: %v", err)
	}

	_, _, err = supaClient.From("users").Select("*", "exact", false).Limit(1, "").Execute()

	if err != nil {
		t.Errorf("Supabase connection check failed! Error: %v", err)
	}

	t.Log("✅ Supabase connection success!")
}
