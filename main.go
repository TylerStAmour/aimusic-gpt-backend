package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	server := &http.Server{
		Addr: "127.0.0.1:8080",
		Handler: Get(),
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

const Endpoint = "https://api.openai.com/v1/completions"

func Get() *gin.Engine {
	router := gin.New()

	gpt := router.Group("/api/gpt")
	gpt.POST("/prompt", postGPTPrompt)

	return router
}

func postGPTPrompt(ctx *gin.Context) {
	prompt := ctx.Query("prompt")

	payload := fmt.Sprintf("\"model\": \"text-babbage-001\", \"prompt\": \"%s\", \"temperature\": 0, \"max_tokens\": 512}", prompt)
	request, err := http.NewRequest("POST", Endpoint, bytes.NewBufferString(payload))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer " + os.Getenv("OPENAI_ACCESS_KEY"))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, data)
}