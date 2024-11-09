package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cakewai/cakewai.com/component/response"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func UploadCloud() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse multipart form with a max upload size of 10 MB
		err := c.Request.ParseMultipartForm(10 << 20)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.BasicResponse{
				Code:    http.StatusBadRequest,
				Message: "Failed to parse multipart form",
				Error:   err.Error(),
			})
			return
		}

		// Retrieve the uploaded file
		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, response.BasicResponse{
				Code:    http.StatusBadRequest,
				Message: "Error retrieving the file",
				Error:   err.Error(),
			})
			return
		}
		defer file.Close()

		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		// Load environment variables from .env file
		errEnv := godotenv.Load()
		if errEnv != nil {
			log.Fatalf("Error loading .env file")
		}

		// Initialize Cloudinary configuration
		cld, err := cloudinary.NewFromParams(
			os.Getenv("CLOUD_NAME"),
			os.Getenv("CLOUD_API_KEY"),
			os.Getenv("CLOUD_API_SECRET"),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.BasicResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error initializing Cloudinary configuration",
				Error:   err.Error(),
			})
			return
		}

		// Get the file name without extension for the public ID
		var publicID = handler.Filename[:len(handler.Filename)-4]

		// Upload file to Cloudinary
		uploadApi := cld.Upload
		result, err := uploadApi.Upload(context.Background(), file, uploader.UploadParams{PublicID: publicID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.BasicResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error uploading file to Cloudinary",
				Error:   err.Error(),
			})
			return
		}

		fmt.Println("File Successfully Uploaded to Cloudinary")
		fmt.Println("Public ID:", result.PublicID)
		fmt.Println("URL:", result.SecureURL)

		// Respond to the client
		c.JSON(http.StatusOK, response.BasicResponse{
			Code:    http.StatusOK,
			Message: "File successfully uploaded to Cloudinary",
			Error:   "no error",
			Data:    result.SecureURL,
		})
	}

}
