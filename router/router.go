package router

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

const Endpoint = "https://api.openai.com/v1/completions"

func Get() *gin.Engine {
	router := gin.New()

	gpt := router.Group("/api/gpt")
	gpt.POST("/prompt", postGPTPrompt)

	gpt.Use(setAccessControl())
	return router
}

func setAccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://aimusic.ca")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST")
	}
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


