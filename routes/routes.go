package routes

import (
	"log"
	"time" // <-- Pode ser necessário importar "time"

	"github.com/coutinhogcs/api-go-gin/controllers"
	"github.com/coutinhogcs/api-go-gin/middleware"
	"github.com/gin-contrib/cors" // <-- 1. IMPORTAR O CORS
	"github.com/gin-gonic/gin"
)

func HandleRequest() {
	router := gin.Default()

	// --- 2. ADICIONAR O MIDDLEWARE DE CORS AQUI ---
	router.Use(cors.New(cors.Config{
		// Permite que o seu frontend (em localhost:3000) acesse esta API
		AllowOrigins: []string{"http://localhost:3000"},

		// Métodos que o frontend pode usar
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		// Cabeçalhos que o frontend pode enviar
		// (IMPORTANTE: "Authorization" é necessário para o seu "crachá" JWT)
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- Rotas Públicas (Login) ---
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/google/login", controllers.HandleGoogleLogin)
		authGroup.GET("/google/callback", controllers.HandleGoogleCallback)
	}

	// --- Rotas Privadas (Exigem o "crachá" JWT) ---
	api := router.Group("/api/v1")
	api.Use(middleware.RequireAuth)
	{
		api.GET("/me", controllers.GetMyProfile)
	}

	log.Println("Servidor Gin rodando na porta 8080")
	router.Run(":8080")
}
