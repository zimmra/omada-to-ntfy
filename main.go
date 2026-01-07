package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/zimmra/omada-to-ntfy/ntfy"
	"github.com/zimmra/omada-to-ntfy/webhook"
)

var version = "development"

func main() {
	logger := log.Default()

	_, server, port, err := InitMain(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Printf("omada-to-ntfy %s server starting on port %s ...", version, port)

	logger.Fatal(http.ListenAndServe(":"+port, server))
}

func InitMain(logger *log.Logger) (nc ntfy.NtfyClient, s *webhook.WebhookServer, p string, err error) {
	ntfyURL := os.Getenv("NTFY_URL")
	if ntfyURL == "" {
		return ntfy.NtfyClient{}, nil, "", errors.New("NTFY_URL environment variable is required")
	}

	ntfyTopic := os.Getenv("NTFY_TOPIC")
	if ntfyTopic == "" {
		return ntfy.NtfyClient{}, nil, "", errors.New("NTFY_TOPIC environment variable is required")
	}

	// Username and password are optional for ntfy (some instances may not require auth)
	ntfyUser := os.Getenv("NTFY_USER")
	ntfyPassword := os.Getenv("NTFY_PASSWORD")

	sharedSecret := os.Getenv("OMADA_SHARED_SECRET")
	if sharedSecret == "" {
		return ntfy.NtfyClient{}, nil, "", errors.New("OMADA_SHARED_SECRET environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ntfyClient := ntfy.NtfyClient{
		NtfyURL:  ntfyURL,
		Topic:    ntfyTopic,
		Username: ntfyUser,
		Password: ntfyPassword,
		Logger:   logger,
	}

	server := &webhook.WebhookServer{
		NtfyClient:   ntfyClient,
		SharedSecret: sharedSecret,
		Logger:       logger,
	}

	return ntfyClient, server, port, nil
}

// EOF
