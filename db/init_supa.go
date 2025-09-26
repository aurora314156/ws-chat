package db

import (
	"fmt"
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
)

func NewSupabaseClient() (*supabase.Client, error) {
	SUPABASEURL := os.Getenv("SUPABASEURL")
	SUPABASEKEY := os.Getenv("SUPABASEKEY")
	if SUPABASEURL == "" {
		return nil, fmt.Errorf("Supabase URL cannot be empty")
	}
	if SUPABASEKEY == "" {
		return nil, fmt.Errorf("Supabase API Key cannot be empty")
	}
	var err error
	supaClient, err := supabase.NewClient(SUPABASEURL, SUPABASEKEY, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}

	return supaClient, nil
}
