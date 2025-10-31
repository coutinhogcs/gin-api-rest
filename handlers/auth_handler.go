package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/coutinhogcs/api-go-gin/config"
	"github.com/coutinhogcs/api-go-gin/utils"
	"github.com/gin-gonic/gin"
)

const oauthStateCookieName = "oauthstate"

// HandleGoogleLogin inicia o fluxo de login do Google.
// Redireciona o usuário para a tela de consentimento do Google.
func HandleGoogleLogin(c *gin.Context) {
	state, err := utils.GenerateOauthStateString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar o 'state'"})
		return
	}

	// Armazena o 'state' em um cookie seguro
	c.SetCookie(oauthStateCookieName, state, 3600, "/", "localhost", false, true)

	// Gera a URL de redirecionamento
	url := config.AppConfig.GoogleOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleCallback é a URL de callback que o Google chama após o consentimento.
func HandleGoogleCallback(c *gin.Context) {
	// 1. Validar o 'state'
	stateCookie, err := c.Cookie(oauthStateCookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cookie de 'state' não encontrado"})
		return
	}

	if c.Query("state") != stateCookie {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incompatibilidade de 'state'"})
		return
	}

	// 2. Limpar o cookie do 'state'
	c.SetCookie(oauthStateCookieName, "", -1, "/", "localhost", false, true)

	// 3. Trocar o código por um token
	code := c.Query("code")
	token, err := config.AppConfig.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao trocar o código pelo token: " + err.Error()})
		return
	}

	// 4. Buscar informações do usuário no Google
	client := config.AppConfig.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar dados do usuário: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao ler os dados do usuário"})
		return
	}

	// Struct temporária para os dados do Google
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.Unmarshal(userData, &googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao processar dados do usuário"})
		return
	}

	/*
		--- PRÓXIMOS PASSOS: LÓGICA DE BANCO DE DADOS E JWT ---

		1. Chamar um `userService.FindOrCreateUser(...)`
			- Tentar encontrar um usuário no seu DB com `provider="google"` e `providerID=googleUser.ID`
			- Se não encontrar, criar um novo usuário com `models.NewUser(...)`

		2. Gerar um Token JWT
			- Criar um JWT (JSON Web Token) contendo o `user.ID`
			- Este token será usado para autenticar o usuário nas rotas privadas da sua API

		3. Enviar o JWT para o Frontend
			- Redirecionar o usuário de volta para o seu site (Frontend)
			- Ex: c.Redirect(http.StatusTemporaryRedirect, "http://meu-frontend.com/dashboard?token=SEU_JWT")
			- Ou enviar o JWT em um cookie HttpOnly
	*/

	// Por enquanto, vamos apenas retornar os dados do usuário como JSON
	log.Printf("Usuário logado: %s (%s)", googleUser.Name, googleUser.Email)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Login bem-sucedido!",
		"user_info": googleUser,
		"token":     "AQUI_VAI_O_SEU_JWT_FUTURO",
	})
}
