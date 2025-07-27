package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	baseURL = "http://localhost:8080"
)

func TestMain(m *testing.M) {
	// Give the Docker Compose services some time to start up
	time.Sleep(10 * time.Second)
	m.Run()
}

func TestUserRegistrationAndLogin(t *testing.T) {
	t.Skip("Skipping integration test")
	username := "testuser"
	password := "testpassword"

	// Register user
	registerPayload := map[string]string{
		"username": username,
		"password": password,
	}
	registerBody, _ := json.Marshal(registerPayload)
	resp, err := http.Post(fmt.Sprintf("%s/api/users/register", baseURL), "application/json", bytes.NewBuffer(registerBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Login user
	loginPayload := map[string]string{
		"username": username,
		"password": password,
	}
	loginBody, _ := json.Marshal(loginPayload)
	resp, err = http.Post(fmt.Sprintf("%s/api/users/login", baseURL), "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResponse struct {
		Status string `json:"status"`
		Data   struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	assert.NoError(t, err)
	assert.Equal(t, "success", loginResponse.Status)
	assert.NotEmpty(t, loginResponse.Data.Token)
}

func TestURLShortening(t *testing.T) {
	t.Skip("Skipping integration test")
	// First, register and login a user to get a token
	username := "urltestuser"
	password := "urltestpassword"
	registerPayload := map[string]string{
		"username": username,
		"password": password,
	}
	registerBody, _ := json.Marshal(registerPayload)
	_, err := http.Post(fmt.Sprintf("%s/api/users/register", baseURL), "application/json", bytes.NewBuffer(registerBody))
	assert.NoError(t, err)

	loginPayload := map[string]string{
		"username": username,
		"password": password,
	}
	loginBody, _ := json.Marshal(loginPayload)
	resp, err := http.Post(fmt.Sprintf("%s/api/users/login", baseURL), "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	var loginResponse struct {
		Status string `json:"status"`
		Data   struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	assert.NoError(t, err)
	token := loginResponse.Data.Token

	// Create a short URL
	originalURL := "https://www.google.com"
	createURLPayload := map[string]string{
		"original_url": originalURL,
	}
	createURLBody, _ := json.Marshal(createURLPayload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/urls", baseURL), bytes.NewBuffer(createURLBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createURLResponse struct {
		Status string `json:"status"`
		Data   struct {
			ShortCode string `json:"short_code"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&createURLResponse)
	assert.NoError(t, err)
	assert.Equal(t, "success", createURLResponse.Status)
	assert.NotEmpty(t, createURLResponse.Data.ShortCode)

	shortCode := createURLResponse.Data.ShortCode

	// Test redirection
	resp, err = http.Get(fmt.Sprintf("%s/%s", baseURL, shortCode))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
	assert.Equal(t, originalURL, resp.Header.Get("Location"))

	// Test QR code generation
	resp, err = http.Get(fmt.Sprintf("%s/%s/qr", baseURL, shortCode))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))
	qrCodeBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, qrCodeBytes)
}