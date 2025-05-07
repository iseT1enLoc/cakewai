package handlers

import (
	"bytes"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/service"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (gen *GeminiHandler) GenerateFineGrainPromptWithNebius() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userInput Prompt
		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid input: failed to parse JSON",
				Error:   err.Error(),
			})
			return
		}

		instruction := fmt.Sprintf(
			"Role: you're an experienced translator. "+
				"Task: translate the user input to English. "+
				"User input: '%s'. "+
				"Response format: only the translation text.", userInput.Text,
		)

		geminiOutput, err := service.GeminiReqClientPrompt(ctx, instruction)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error occurred while requesting Gemini translation",
				Error:   err.Error(),
			})
			return
		}

		fineGrainPrompt := fmt.Sprintf(
			"Role: You are a skilled cakes shop owner. " +
				"Please generate an image of a single cake based on this user input: " + geminiOutput +
				"The cake should be modern, elegant, and creative, with a clean background that highlights the cake's details, " +
				"ensuring it is visually appealing and ready for the user to choose.")

		requestBody := map[string]string{
			"prompt": fineGrainPrompt,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to marshal request body to JSON",
				Error:   err.Error(),
			})
			return
		}
		baseimaegURL := os.Getenv("BASE_IMAGE_GENERATION")
		apiURL := baseimaegURL + "/api/v1/generate-cake-image"
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

		// ✅ Check that response is an image
		contentType := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			bodyBytes, _ := io.ReadAll(resp.Body)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "The image generation API did not return a valid image",
				Error:   fmt.Sprintf("Unexpected content type: %s\nBody: %s", contentType, string(bodyBytes)),
			})
			return
		}

		// ✅ Read the image bytes
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to read image response",
				Error:   err.Error(),
			})
			return
		}

		// Step 2: Prepare multipart request to /upload-cloud
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		fw, err := w.CreateFormFile("file", "cake.png")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create form file",
				Error:   err.Error(),
			})
			return
		}
		_, err = fw.Write(respBytes)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to write file to form",
				Error:   err.Error(),
			})
			return
		}
		w.Close()
		base_backend := os.Getenv("BASE_BACKEND")
		// Step 3: POST multipart form to upload API
		uploadReq, err := http.NewRequest("POST", base_backend+"/api/upload", &b)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create upload request",
				Error:   err.Error(),
			})
			return
		}
		uploadReq.Header.Set("Content-Type", w.FormDataContentType())

		uploadResp, err := client.Do(uploadReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Cloudinary upload failed",
				Error:   err.Error(),
			})
			return
		}
		defer uploadResp.Body.Close()

		if uploadResp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(uploadResp.Body)
			ctx.JSON(uploadResp.StatusCode, response.FailedResponse{
				Code:    uploadResp.StatusCode,
				Message: "Cloudinary upload returned error",
				Error:   string(bodyBytes),
			})
			return
		}

		var result map[string]interface{}
		err = json.NewDecoder(uploadResp.Body).Decode(&result)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to decode Cloudinary response",
				Error:   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Image generated and uploaded successfully",
			},
			Data: result["data"], // Cloudinary secure URL
		})
	}
}
