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
		finegrain_prompt := fmt.Sprintf(
			"Role: You are an experienced cake maker with years of expertise in creating beautifully customized cakes tailored to clients' requests. "+
				"Task: Generate a high-resolution image of a customized cake based on the following user input: '%s'. "+
				"The cake should have intricate details such as realistic textures, smooth icing, and vibrant colors. "+
				"Consider incorporating creative design elements such as edible flowers, custom decorations, and personalized messages. "+
				"The style should reflect an elegant and modern aesthetic, ensuring the cake looks both delicious and visually stunning. "+
				"Ensure the lighting and shadows in the image highlight the cake's details, making it appear as realistic as possible. "+
				"The background should be simple and clean to keep the focus on the cake, with soft lighting to emphasize its texture and color.",
			gemini_output)

		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    0,
				Message: "Successfully translate user input",
			},
			Data: finegrain_prompt,
		})
	}
}
