package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User define o modelo do usuário no banco de dados
type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;"`
	Email      string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Name       string    `gorm:"type:varchar(100)"`
	AvatarURL  string
	Provider   string `gorm:"type:varchar(50);not null"`
	ProviderID string `gorm:"type:varchar(100);not null;uniqueIndex:idx_provider_providerid"` // Índice único para (Provider + ProviderID)
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// BeforeCreate é um "hook" do GORM que gera um UUID antes de salvar
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}

// package models

// import (
// 	"gopkg.in/validator.v2"
// 	"gorm.io/gorm"
// )

// type Aluno struct {
// 	gorm.Model
// 	Nome string `json:"nome" validate:"nonzero"`
// 	CPF  string `json:"cpf" validate:"len=11, regexp=^[0-9]*$"`
// 	RG   string `json:"rg" validate:"len=9, regexp=^[0-9]*$"`
// }

// func ValidaDadosDeAluno(aluno *Aluno) error {
// 	if err := validator.Validate(aluno); err != nil {
// 		return err
// 	}
// 	return nil
// }

// var Alunos []Aluno
