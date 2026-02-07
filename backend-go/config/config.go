package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// Massive API
	MassiveAPIKey  string
	MassiveBaseURL string

	// Supabase
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string

	// Server
	Port    string
	GinMode string

	// Database connection string (constructed from Supabase credentials)
	DatabaseURL string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	viper.SetConfigFile("../.env")
	viper.AutomaticEnv()

	// Try to read .env file (optional, env vars take precedence)
	_ = viper.ReadInConfig()

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("MASSIVE_BASE_URL", "https://api.massive.com/v3")

	config := &Config{
		MassiveAPIKey:      viper.GetString("MASSIVE_API_KEY"),
		MassiveBaseURL:     viper.GetString("MASSIVE_BASE_URL"),
		SupabaseURL:        viper.GetString("SUPABASE_URL"),
		SupabaseAnonKey:    viper.GetString("SUPABASE_ANON_KEY"),
		SupabaseServiceKey: viper.GetString("SUPABASE_SERVICE_KEY"),
		Port:               viper.GetString("PORT"),
		GinMode:            viper.GetString("GIN_MODE"),
	}

	// Validate required fields
	if config.MassiveAPIKey == "" {
		return nil, fmt.Errorf("MASSIVE_API_KEY is required")
	}
	if config.SupabaseURL == "" {
		return nil, fmt.Errorf("SUPABASE_URL is required")
	}
	if config.SupabaseServiceKey == "" {
		return nil, fmt.Errorf("SUPABASE_SERVICE_KEY is required")
	}

	// Construct database URL from Supabase URL
	// Format: postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres
	// We'll use the service key for server-side database access
	config.DatabaseURL = constructDatabaseURL(config.SupabaseURL, config.SupabaseServiceKey)

	return config, nil
}

// constructDatabaseURL builds the PostgreSQL connection string from Supabase credentials
// Note: For production, you should use the actual database password, not the service key
// This is a placeholder - you'll get the real connection string from Supabase dashboard
func constructDatabaseURL(supabaseURL, serviceKey string) string {
	// For now, return empty - we'll set this up properly when connecting to database
	// The actual connection string should come from Supabase dashboard: Settings -> Database
	return ""
}
