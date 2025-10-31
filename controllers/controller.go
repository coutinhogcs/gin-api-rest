package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/coutinhogcs/api-go-gin/config"
	"github.com/coutinhogcs/api-go-gin/database"
	"github.com/coutinhogcs/api-go-gin/models"
	"github.com/coutinhogcs/api-go-gin/utils" // <-- Garanta que "utils" está importado
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const oauthStateCookieName = "oauthstate"

// HandleGoogleLogin (Não muda)
func HandleGoogleLogin(c *gin.Context) {
	state, err := utils.GenerateOauthStateString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar o 'state'"})
		return
	}
	c.SetCookie(oauthStateCookieName, state, 3600, "/", "localhost", false, true)
	url := config.AppConfig.GoogleOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleCallback (Atualizado para GERAR o JWT)
func HandleGoogleCallback(c *gin.Context) {
	// 1. Validar 'state'
	stateCookie, err := c.Cookie(oauthStateCookieName)
	if err != nil || c.Query("state") != stateCookie {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State inválido ou cookie não encontrado"})
		return
	}
	c.SetCookie(oauthStateCookieName, "", -1, "/", "localhost", false, true)

	// 2. Trocar código por token
	code := c.Query("code")
	token, err := config.AppConfig.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao trocar o código: " + err.Error()})
		return
	}

	// 3. Buscar dados do usuário no Google
	client := config.AppConfig.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar dados do usuário: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	userData, _ := io.ReadAll(resp.Body)
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	json.Unmarshal(userData, &googleUser)

	// 4. LÓGICA DE BANCO DE DADOS: Encontrar ou Criar o Usuário
	var user models.User
	err = database.DB.Where("provider = ? AND provider_id = ?", "google", googleUser.ID).First(&user).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		// --- Usuário NÃO existe: Criar ---
		log.Println("Usuário não encontrado. Criando novo usuário...")
		user = models.User{
			Email:      googleUser.Email,
			Name:       googleUser.Name,
			AvatarURL:  googleUser.Picture,
			Provider:   "google",
			ProviderID: googleUser.ID,
		}
		if err = database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar usuário no banco"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro no banco de dados ao buscar usuário"})
		return
	}

	// 5. Usuário (encontrado ou criado) está na variável 'user'
	log.Printf("Usuário %s logado com sucesso (ID: %s)", user.Email, user.ID)

	// --- 6. GERAR O TOKEN (A MUDANÇA) ---
	tokenString, err := utils.GenerateJWT(user.ID) // <-- CHAMA A FUNÇÃO QUE CRIAMOS
	if err != nil {
		log.Printf("Erro ao gerar JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar token de autenticação"})
		return
	}

	// --- 7. REDIRECIONAR PARA O FRONTEND (A CORREÇÃO) ---
	// <-- CORREÇÃO: Removemos o bloco c.JSON daqui.

	frontendURL := "http://localhost:3000"

	// Anexa o "crachá" como um parâmetro na URL
	redirectURL := frontendURL + "?token=" + tokenString

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GetMyProfile (A NOVA FUNÇÃO que seu routes.go precisa)
func GetMyProfile(c *gin.Context) {
	// 1. Pegar o ID do usuário que o middleware (RequireAuth) anexou
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao obter ID do usuário do contexto"})
		return
	}

	// <-- CORREÇÃO: Removemos a geração de token e o redirecionamento מכאן.

	// 2. Buscar o usuário no banco de dados
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// 3. Retornar os dados do usuário em JSON (A CORREÇÃO)
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}
