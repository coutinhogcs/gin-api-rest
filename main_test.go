package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coutinhogcs/api-go-gin/database"
	"github.com/coutinhogcs/api-go-gin/models"
	"github.com/coutinhogcs/api-go-gin/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	R *gin.Engine
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	database.DbConection()
	R = gin.Default()
	routes.SetupRoutes(R)

	m.Run()
}

func TestCriaAlunoMock(t *testing.T) {
	aluno := models.Aluno{Nome: "Nome do Aluno Teste", CPF: "12345678911", RG: "123456789"}
	alunoJson, _ := json.Marshal(aluno)
	req, _ := http.NewRequest("POST", "/alunos", bytes.NewBuffer(alunoJson))
	resp := httptest.NewRecorder()
	R.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)
	fmt.Print(resp.Code)
}

func TestDeletaAlunoMock(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/alunos/15", nil)
	resp := httptest.NewRecorder()
	R.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestBuscaAlunoPorID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/alunos/2", nil)
	res := httptest.NewRecorder()
	R.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestExibeTodosAlunos(t *testing.T) {
	req, _ := http.NewRequest("GET", "/alunos", nil)
	resp := httptest.NewRecorder()
	R.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestAtualizaAlunoPorID(t *testing.T) {
	body := models.Aluno{Nome: "Nome Atualizado", CPF: "00099988876", RG: "987654321"}
	bodyJson, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", "/alunos/16", bytes.NewBuffer(bodyJson))
	resp := httptest.NewRecorder()
	R.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
