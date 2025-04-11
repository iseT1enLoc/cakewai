package handlers

import (
	"bytes"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// type GeminiHandler struct {
// 	ApiKey string
// 	Model  string
// }
// type Prompt struct {
// 	Text string `json:"user_input"`
// }

// func (gen *GeminiHandler) GenerateFineGrainPrompt() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		// Declare a variable to hold the input
// 		var userInput Prompt

// 		if err := ctx.ShouldBindBodyWithJSON(&userInput); err != nil {
// 			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
// 				Code:    400,
// 				Message: "Can not parsing the json file",
// 				Error:   err.Error(),
// 			})
// 		}

// 		instruction := fmt.Sprintf("Role: you're an experienced translator. "+
// 			"Task: translate the user input to English. "+
// 			"User input %s"+
// 			"Response format:only the translation text", userInput.Text)
// 		fmt.Print(userInput.Text)

// 		gemini_output, err := service.GeminiReqClientPrompt(ctx, instruction)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
// 				Code:    400,
// 				Message: "Error happened while request to gemini",
// 				Error:   err.Error(),
// 			})
// 		}
// 		finegrain_prompt := fmt.Sprintf(
// 			"Role: You are an experienced cake maker with years of expertise in creating beautifully customized cakes tailored to clients' requests. "+
// 				"Task: Generate a high-resolution image of a customized cake based on the following user input: '%s'. "+
// 				"The cake should have intricate details such as realistic textures, smooth icing, and vibrant colors. "+
// 				"Consider incorporating creative design elements such as edible flowers, custom decorations, and personalized messages. "+
// 				"The style should reflect an elegant and modern aesthetic, ensuring the cake looks both delicious and visually stunning. "+
// 				"Ensure the lighting and shadows in the image highlight the cake's details, making it appear as realistic as possible. "+
// 				"The background should be simple and clean to keep the focus on the cake, with soft lighting to emphasize its texture and color.",
// 			gemini_output)

// 		// Define API URL and token
// 		apiURL := "https://api.claid.ai/v1-beta1/image/generate"
// 		//apiToken := "fe71fd94cad54f09917ce22b20d6d0ff" // Replace with your token

// 		if err := ctx.ShouldBindJSON(&userInput); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 			return
// 		}
// 		// Prepare input and options
// 		requestBody := map[string]interface{}{
// 			"input": finegrain_prompt, // Replace with your input
// 			"options": map[string]interface{}{
// 				"number_of_images": 1,
// 				"guidance_scale":   5.0,
// 			},
// 		}

// 		// Convert request body to JSON
// 		jsonBody, err := json.Marshal(requestBody)
// 		if err != nil {
// 			fmt.Println("Error marshalling JSON:", err)
// 			return
// 		}

// 		// Create a POST request
// 		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API request"})
// 			return
// 		}

// 		// Send the request using an HTTP client
// 		client := &http.Client{}
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println("Error sending request:", err)
// 			return
// 		}
// 		defer resp.Body.Close()

//			// Read the response
//			var result map[string]interface{}
//			err = json.NewDecoder(resp.Body).Decode(&result)
//			if err != nil {
//				fmt.Println("Error decoding response:", err)
//				return
//			}
//			ctx.JSON(http.StatusOK, response.Success{
//				ResponseFormat: response.ResponseFormat{
//					Code:    200,
//					Message: "Successfully translate user input",
//				},
//				Data: resp,
//			})
//		}
//	}
type GeminiHandler struct {
	ApiKey string
	Model  string
}

type Prompt struct {
	Text string `json:"user_input" binding:"required"`
}

func (gen *GeminiHandler) GenerateFineGrainPrompt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse the user input
		var userInput Prompt
		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid input: failed to parse JSON",
				Error:   err.Error(),
			})
			return
		}

		// Create the instruction for translation
		instruction := fmt.Sprintf(
			"Role: you're an experienced translator. "+
				"Task: translate the user input to English. "+
				"User input: '%s'. "+
				"Response format: only the translation text.", userInput.Text,
		)

		// Call the Gemini translation service
		geminiOutput, err := service.GeminiReqClientPrompt(ctx, instruction)
		fmt.Println(geminiOutput)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error occurred while requesting Gemini translation",
				Error:   err.Error(),
			})
			return
		}

		// Generate the fine-grain prompt
		// fineGrainPrompt := fmt.Sprintf(
		// 	"Role: You are an experienced cake maker with years of expertise in creating beautifully customized cakes tailored to clients' requests. "+
		// 		"Task: Generate a high-resolution image of a customized cake based on the following user input: '%s'. "+
		// 		"The cake should have intricate details such as realistic textures, smooth icing, and vibrant colors. "+
		// 		"Consider incorporating creative design elements such as edible flowers, custom decorations, and personalized messages. "+
		// 		"The style should reflect an elegant and modern aesthetic, ensuring the cake looks both delicious and visually stunning. "+
		// 		"Ensure the lighting and shadows in the image highlight the cake's details, making it appear as realistic as possible. "+
		// 		"The background should be simple and clean to keep the focus on the cake, with soft lighting to emphasize its texture and color.",
		// 	geminiOutput,
		// )
		// Generate the fine-grain prompt
		// Generate the fine-grain prompt
		// fineGrainPrompt := fmt.Sprintf(
		// 	"Role: You are a skilled cake artist with experience in creating detailed, customized cakes. " +
		// 		"Task: Generate a high-resolution image of a cake based on this user input: '%s'. " +
		// 		"The style should be modern and elegant,creative ,with a clean background that highlights the cake's details. " +
		// 		geminiOutput,
		// )
		fineGrainPrompt := fmt.Sprintf(
			"Role: You are a skilled cakes shop owner. " +
				"Please generate an image of a single cake based on this user input: " + geminiOutput +
				"The cake should be modern, elegant, and creative, with a clean background that highlights the cake's details, " +
				"ensuring it is visually appealing and ready for the user to choose.")
		println(fineGrainPrompt)
		// Prepare the request payload for the image generation API
		requestBody := map[string]interface{}{
			"input": fineGrainPrompt,
			"options": map[string]interface{}{
				"number_of_images": 1,
				"guidance_scale":   5.0,
			},
		}

		// Convert the request body to JSON
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to marshal request body to JSON",
				Error:   err.Error(),
			})
			return
		}

		// Define the API URL and set headers
		//apiURL := "https://api.claid.ai/v1-beta1/image/generate"
		apiURL := "http://localhost:5050/api/v1/generate-cake-image"
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create HTTP request",
				Error:   err.Error(),
			})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLAIDAI"))) // Use ApiKey from handler struct

		// Send the request using an HTTP client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error occurred while sending request to image generation API",
				Error:   err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		// Check if the response status is not successful
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			ctx.JSON(resp.StatusCode, response.FailedResponse{
				Code:    resp.StatusCode,
				Message: "Image generation API returned an error",
				Error:   string(bodyBytes),
			})
			return
		}

		// Decode the response body
		var apiResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to decode API response",
				Error:   err.Error(),
			})
			return
		}

		// Return only the "output" from the response
		if output, ok := apiResponse["data"].(map[string]interface{})["output"]; ok {
			ctx.JSON(http.StatusOK, response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "Successfully processed user input and generated image",
				},
				Data: output,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Missing 'output' field in API response",
				Error:   "Output data not found",
			})
		}
	}
}
func (gen *GeminiHandler) GeminiHandlerNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data Prompt
		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "nothing interesting"})
			return
		}
		res, err := service.GenImageRequest(ctx, data.Text)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "nothing interesting"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": res})
	}
}
