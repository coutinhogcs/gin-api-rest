package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/coutinhogcs/api-go-gin/config" // Importa a config (pelo SessionSecret)
	"github.com/coutinhogcs/api-go-gin/utils"  // Importa o utils (pela struct Claims)
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth é o nosso "segurança"
func RequireAuth(c *gin.Context) {
	// 1. Pegar o "crachá" (token) do cabeçalho
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Println("Cabeçalho de autorização vazio")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado: token não fornecido"})
		return
	}

	// 2. O token deve vir no formato "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado: formato de token inválido"})
		return
	}

	// 3. Validar o "crachá"
	claims := &utils.Claims{} // Usa a struct Claims que definimos em utils/jwt_utils.go

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura é o que esperamos (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, log.New(log.Writer(), "método de assinatura inesperado", log.LstdFlags).Output(2, "")
		}
		// Retorna o nosso segredo (do .env) para validação
		return []byte(config.AppConfig.SessionSecret), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Token inválido: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado: token inválido"})
		return
	}

	// 4. O "crachá" é válido!
	// Anexamos o ID do usuário à requisição para que o próximo controlador saiba quem ele é.
	c.Set("userID", claims.UserID)

	// 5. Deixar a requisição passar para o próximo controlador
	c.Next()
}
