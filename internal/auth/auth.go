package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kaputi/nikaro/internal/res"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Authorization struct {
}

func NewAthorization() *Authorization {
	return &Authorization{}
}

func NewContext(ctx context.Context, token *jwt.Token, claims *CustomClaims) context.Context {
	ctx = context.WithValue(ctx, "token", token)
	ctx = context.WithValue(ctx, "claims", claims)
	return ctx
}

func (a *Authorization) GenerateToken(id, role string, exp time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *Authorization) ParseToken(tokenString string) (*jwt.Token, *CustomClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			// TODO:
			return nil, nil, errors.New("token expired")
		}
		return nil, nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, fmt.Errorf("invalid token")
}

func (a *Authorization) SetTokenToCookie(w http.ResponseWriter, name, token, path string, exp time.Duration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    token,
		Expires:  time.Now().Add(exp),
		Secure:   true,
		HttpOnly: true,
		Path:     path,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func (a *Authorization) ClearTokenCookie(w http.ResponseWriter, r *http.Request, name string) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return
	}

	cookie.Expires = time.Now()
	cookie.Value = ""

	http.SetCookie(w, cookie)
}

func (a *Authorization) GetTokenFromCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (a *Authorization) VerifyToken(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := a.GetTokenFromCookie(r, name)
			if tokenString == "" {
				res.Fail(w, "No token", http.StatusUnauthorized)
				return
			}
			token, claims, err := a.ParseToken(tokenString)
			if err != nil {
				res.Fail(w, err.Error(), http.StatusUnauthorized)
				// token has invalid claims: token is expired
				return
			}
			r = r.WithContext(NewContext(r.Context(), token, claims))
			next.ServeHTTP(w, r)
		})
	}
}

func (a *Authorization) AuthorizeAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, _ := r.Context().Value("claims").(*CustomClaims)

			if claims.Role != "admin" {
				res.Fail(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)

		})
	}
}

func (a *Authorization) GetTokenFromContext(ctx context.Context) (*jwt.Token, *CustomClaims) {
	token, _ := ctx.Value("token").(*jwt.Token)
	claims := ctx.Value("claims").(*CustomClaims)
	return token, claims
}
