package routes

import (
	"github.com/coutinhogcs/api-go-gin/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	r.GET("/alunos/:id", controllers.ProcuraAlunoPorId)
	r.GET("alunos/cpf/:cpf", controllers.BuscaPorCpf)
	r.POST("/alunos", controllers.CriaNovoAluno)
	r.DELETE("/alunos/:id", controllers.DeletaAlunoPorId)
	r.PATCH("alunos/:id", controllers.AtualizaAluno)
}

func HandleRequest() {
	r := gin.Default()

	// Chama a função de setup
	SetupRoutes(r)

	// Inicia o servidor
	r.Run(":9000")
}
