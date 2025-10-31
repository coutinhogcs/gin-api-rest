package utils

import (
	"time"

	"github.com/coutinhogcs/api-go-gin/config" // Importa o config pelo segredo
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims é o "molde" do nosso "crachá" (JWT).
// É a struct que está faltando no seu middleware.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT cria um novo token JWT para um usuário
func GenerateJWT(userID uuid.UUID) (string, error) {
	// O token expira em 7 dias
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		UserID: userID, // Nosso dado customizado
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "serb-api",
		},
	}

	// Assina o token com o nosso SESSION_SECRET do .env
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.SessionSecret))

	return tokenString, err
}
