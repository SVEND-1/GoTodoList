package cookies

import (
	"net/http"
	"time"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

type Options struct {
	Secure   bool
	Domain   string
	SameSite http.SameSite
}

func SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken string, accessTTL, refreshTTL time.Duration, opts Options) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    accessToken,
		Path:     "/",
		Domain:   opts.Domain,
		Expires:  time.Now().Add(accessTTL),
		HttpOnly: true,
		Secure:   opts.Secure,
		SameSite: opts.SameSite,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    refreshToken,
		Path:     "/api/v1/auth/refresh",
		Domain:   opts.Domain,
		Expires:  time.Now().Add(refreshTTL),
		HttpOnly: true,
		Secure:   opts.Secure,
		SameSite: opts.SameSite,
	})
}

func ClearAuthCookies(w http.ResponseWriter, opts Options) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Path:     "/",
		Domain:   opts.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   opts.Secure,
		SameSite: opts.SameSite,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		Domain:   opts.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   opts.Secure,
		SameSite: opts.SameSite,
	})
}

func GetTokenFromCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
