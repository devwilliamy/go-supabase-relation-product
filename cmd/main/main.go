package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devwilliamy/go-supabase-relation-product/internal/supabase"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	supabaseURL := os.Getenv("NEXT_PUBLIC_SUPABASE_URL")
    supabaseKey := os.Getenv("NEXT_PUBLIC_SUPABASE_KEY")
    supabaseAdminURL := os.Getenv("NEXT_PUBLIC_SUPABASE_ADMIN_NODE_URL")
    supabaseAdminKey := os.Getenv("NEXT_PUBLIC_SUPABASE_ADMIN_NODE_KEY")

	client, err := supabase.NewClient(supabaseURL, supabaseKey)
	if err != nil {
		log.Fatalf("Unable to create Supabase client %v", err)
	}

	adminClient, err := supabase.NewClient(supabaseAdminURL, supabaseAdminKey)
	if err != nil {
		log.Fatalf("Unable to create Supabase admin client: %v", err)
	}

	fmt.Println("Start Supabase relation project")
	if err := supabase.CreateRealationalTable(client, adminClient); err!= nil {
		log.Fatalf("Error creating relational table %v", err)
	}
}