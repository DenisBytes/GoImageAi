package sb

import (
	"os"

	"github.com/nedpals/supabase-go"
)

var Client *supabase.Client

// not same as built-in init func
func Init() error {
	sbHost := os.Getenv("SUPABASE_URL")
	sbSecret := os.Getenv("SUPABASE_SECRET")
	Client = supabase.CreateClient(sbHost, sbSecret)

	return nil
}
