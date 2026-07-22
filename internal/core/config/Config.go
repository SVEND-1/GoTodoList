package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	TimeZone *time.Location
	Email    *EmailConfig
}

type EmailConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
}

func NewConfig() (*Config, error) {
	tz := os.Getenv("TIME_ZONE")
	if tz == "" {
		tz = "UTC"
	}

	zone, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("failed to load time zone: %w", err)
	}

	//emailConfig, err := newEmailConfig()
	//if err != nil {
	//	return nil, fmt.Errorf("failed to load email config: %w", err)
	//}

	return &Config{
		TimeZone: zone,
		Email:    nil,
	}, nil
}

//func newEmailConfig() (EmailConfig, error) {
//	host := os.Getenv("SMTP_HOST")
//	if host == "" {
//		host = "smtp.gmail.com"
//	}
//
//	portStr := os.Getenv("SMTP_PORT")
//	if portStr == "" {
//		portStr = "587"
//	}
//	port, err := strconv.Atoi(portStr)
//	if err != nil {
//		return EmailConfig{}, fmt.Errorf("failed to parse SMTP_PORT: %w", err)
//	}
//
//	username := os.Getenv("SMTP_USERNAME")
//	if username == "" {
//		return EmailConfig{}, fmt.Errorf("SMTP_USERNAME is required")
//	}
//
//	password := os.Getenv("SMTP_PASSWORD")
//	if password == "" {
//		return EmailConfig{}, fmt.Errorf("SMTP_PASSWORD is required")
//	}
//
//	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
//	if fromEmail == "" {
//		fromEmail = username
//	}
//
//	return EmailConfig{
//		Host:      host,
//		Port:      port,
//		Username:  username,
//		Password:  password,
//		FromEmail: fromEmail,
//	}, nil
//}

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}
	return config
}
