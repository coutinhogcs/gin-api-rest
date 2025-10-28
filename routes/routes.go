package routes

import (
	"github.com/coutinhogcs/api-go-gin/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequest() {
	r := gin.Default()
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	r.GET("/alunos/:id", controllers.ProcuraAlunoPorId)
	r.GET("alunos/cpf/:cpf", controllers.BuscaPorCpf)
	r.POST("/alunos", controllers.CriaNovoAluno)
	r.DELETE("/alunos/:id", controllers.DeletaAlunoPorId)
	r.PATCH("alunos/:id", controllers.AtualizaAluno)
	r.Run(":9000")
}
