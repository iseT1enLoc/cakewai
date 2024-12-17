package handlers

import (
	"fmt"
	"net/http"

	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/service"

	"github.com/gin-gonic/gin"
)

type GeminiHandler struct {
	ApiKey string
	Model  string
}
type Prompt struct {
	Text string `json:"user_input"`
}

func (gen *GeminiHandler) GenerateFineGrainPrompt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Declare a variable to hold the input
		var userInput Prompt

		if err := ctx.ShouldBindBodyWithJSON(&userInput); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    400,
				Message: "Can not parsing the json file",
				Error:   err.Error(),
			})
		}

		instruction := fmt.Sprintf("Role: you're an experienced translator. "+
			"Task: translate the user input to English. "+
			"User input %s"+
			"Response format:only the translation text", userInput.Text)
		fmt.Print(userInput.Text)

		gemini_output, err := service.GeminiReqClientPrompt(ctx, instruction)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    400,
				Message: "Error happened while request to gemini",
				Error:   err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    0,
				Message: "Successfully translate user input",
			},
			Data: gemini_output,
		})
	}
}
