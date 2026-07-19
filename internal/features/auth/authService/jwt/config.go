package jwt

import (
	"fmt"
	"time"
)

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func newJWTConfig() (JWTConfig, error) { //ПЕРЕПИСАТЬ
	secret := "dsadasdsa-aasd2d3a_sdad123-5142#4q24-as" // Вынести в .env
	if secret == "" {
		return JWTConfig{}, fmt.Errorf("JWT_SECRET is required")
	}

	return JWTConfig{
		Secret:          secret,
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}, nil
}
