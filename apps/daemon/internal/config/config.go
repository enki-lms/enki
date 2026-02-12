package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	OIDC     OIDCConfig     `mapstructure:"oidc"`
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Sandbox  SandboxConfig  `mapstructure:"sandbox"`
	AI       AIConfig       `mapstructure:"ai"`
}

// OIDCConfig holds OIDC provider configuration
type OIDCConfig struct {
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	IssuerURL    string   `mapstructure:"issuer_url"`
	RedirectURL  string   `mapstructure:"redirect_url"`
	Scopes       []string `mapstructure:"scopes"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	URL      string `mapstructure:"url"`
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port        int    `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	FrontendURL string `mapstructure:"frontend_url"`
	UploadsDir  string `mapstructure:"uploads_dir"`
}

// AuthConfig holds authentication/authorization settings
type AuthConfig struct {
	TeacherRoleName string   `mapstructure:"teacher_role_name"`
	AdminEmails     []string `mapstructure:"admin_emails"`
	JWTSecret       string   `mapstructure:"jwt_secret"`
}

// SandboxConfig holds NSJail sandbox settings
type SandboxConfig struct {
	NSJailPath    string `mapstructure:"nsjail_path"`
	PythonPath    string `mapstructure:"python_path"`
	LowerDir      string `mapstructure:"lower_dir"`
	TempDir       string `mapstructure:"temp_dir"`
	TimeoutMs     int    `mapstructure:"timeout_ms"`
	MemoryLimitMB int    `mapstructure:"memory_limit_mb"`
	Unsafe        bool   `mapstructure:"unsafe"`
}

// AIConfig holds AI assistant settings (OAI-compatible API)
type AIConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	BaseURL      string `mapstructure:"base_url"`
	APIKey       string `mapstructure:"api_key"`
	Model        string `mapstructure:"model"`
	SystemPrompt string `mapstructure:"system_prompt"`
}

// Load reads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/enki/")

	// Environment variable support
	viper.SetEnvPrefix("ENKI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.frontend_url", "http://localhost:5173")
	viper.SetDefault("server.uploads_dir", "./uploads")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "enki_user")
	viper.SetDefault("database.password", "enki_password")
	viper.SetDefault("database.dbname", "enki_db")
	viper.SetDefault("auth.teacher_role_name", "teacher")
	viper.SetDefault("oidc.scopes", []string{"openid", "profile", "email"})
	// Sandbox defaults
	viper.SetDefault("sandbox.nsjail_path", "/usr/bin/nsjail")
	viper.SetDefault("sandbox.python_path", "/usr/bin/python3")
	viper.SetDefault("sandbox.lower_dir", "/var/lib/enki/rootfs")
	viper.SetDefault("sandbox.temp_dir", "/tmp/enki-sandbox")
	viper.SetDefault("sandbox.timeout_ms", 5000)
	viper.SetDefault("sandbox.memory_limit_mb", 128)
	// AI defaults
	viper.SetDefault("ai.enabled", false)
	viper.SetDefault("ai.base_url", "https://api.openai.com/v1")
	viper.SetDefault("ai.model", "gpt-4o-mini")
	viper.SetDefault("ai.system_prompt", "You are a helpful programming tutor assisting students with computer science problems. Provide clear explanations and guide students toward understanding concepts rather than just giving answers.")

	// Read config file (optional - env vars can override)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK - we'll use env vars
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that all required configuration is present
func (c *Config) Validate() error {
	if c.OIDC.ClientID == "" {
		return fmt.Errorf("oidc.client_id is required")
	}
	if c.OIDC.ClientSecret == "" {
		return fmt.Errorf("oidc.client_secret is required")
	}
	if c.OIDC.IssuerURL == "" {
		return fmt.Errorf("oidc.issuer_url is required")
	}
	if c.OIDC.RedirectURL == "" {
		return fmt.Errorf("oidc.redirect_url is required")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("auth.jwt_secret is required")
	}
	return nil
}

// DatabaseURL returns the database connection URL
func (c *Config) DatabaseURL() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
	)
}

// ServerAddress returns the server listen address
func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
