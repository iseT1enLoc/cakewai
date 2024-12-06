package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Define a struct for the main response body
type ApiResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    TokenData `json:"data"`
}

func getJWTToken() (string, error) {
	// Assuming your login API returns a JSON response with the token.
	url := "http://localhost:8080/api/public/login" // Adjust to your login endpoint
	payload := `{"email": "25dw0f@gm.uit.edu.vn", "password": "0123456789"}`

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
	fmt.Println(resp)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", resp.Status)
	}

	var result ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}
	fmt.Println(result.Data.AccessToken)
	return result.Data.AccessToken, nil
}

func sendRequest(wg *sync.WaitGroup, url string, payload []byte, token string) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
	}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add JWT token to Authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed request: %s, Status Code: %d\n", url, resp.StatusCode)
		return
	}

	// If you want to log the successful requests
	fmt.Println("Request successful, Status Code:", resp.StatusCode)
}

func main() {
	// Obtain JWT token
	token, err := getJWTToken()
	//fmt.Println(token)
	if err != nil {
		fmt.Println("Failed to get JWT token:", err)
		return
	}

	// Define the backend API endpoint and payload for the POST request
	url := "http://localhost:8080/api/protected/gg" // Adjust to your protected API endpoint
	// Construct the payload with the access token
	payload := []byte("hihi")

	// Define the number of concurrent requests to simulate
	concurrentRequests := 100000 // You can increase this number for higher traffic

	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	// Send concurrent requests
	for i := 0; i < concurrentRequests; i++ {
		go sendRequest(&wg, url, payload, token)
	}

	// Wait for all requests to complete
	wg.Wait()
	fmt.Println("Load test completed")
}
