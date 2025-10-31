package database

import (
	"log" // Vamos usar log.Fatal para parar a app se o DB falhar

	"github.com/coutinhogcs/api-go-gin/config" // 1. Importar o config
	"github.com/coutinhogcs/api-go-gin/models" // 2. Importar os models
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func DbConection() {
	// 3. Carregar as configurações ANTES de tentar conectar
	// Isso vai ler o .env e popular o config.AppConfig
	config.LoadConfig()

	// 4. Pegar a string de conexão do config (que leu o .env)
	strConn := config.AppConfig.DB_DSN
	if strConn == "" {
		log.Fatal("DB_DSN não está definido no arquivo .env")
	}

	DB, err = gorm.Open(postgres.Open(strConn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	log.Println("Conexão com banco de dados bem-sucedida.")

	// 5. Migrar AMBOS os modelos
	// O GORM é inteligente e só vai criar/atualizar o que for necessário.
	err = DB.AutoMigrate(
		&models.User{}, // Nosso novo modelo de usuário para OAuth
	)

	if err != nil {
		log.Fatalf("Erro ao rodar AutoMigrate: %v", err)
	}

	log.Println("Migração do banco de dados concluída.")
}
