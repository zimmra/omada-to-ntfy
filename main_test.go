package main_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	main "github.com/zimmra/omada-to-ntfy"
)

func TestInitMain(t *testing.T) {

	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)

	t.Run("NTFY_URL is required", func(t *testing.T) {
		buf.Reset()
		_, _, _, err := main.InitMain(logger)
		if err.Error() != "NTFY_URL environment variable is required" {
			logger.Fatalf("Failed test whether NTFY_URL is required; log is `%v`", buf.String())
		}
	})

	os.Setenv("NTFY_URL", "https://ntfy.sh")

	t.Run("NTFY_TOPIC is required", func(t *testing.T) {
		buf.Reset()
		_, _, _, err := main.InitMain(logger)
		if err.Error() != "NTFY_TOPIC environment variable is required" {
			logger.Fatalf("Failed test whether NTFY_TOPIC is required; log is `%v`", buf.String())
		}
	})

	os.Setenv("NTFY_TOPIC", "my_omada_alerts")

	t.Run("OMADA_SHARED_SECRET is required", func(t *testing.T) {
		buf.Reset()
		_, _, _, err := main.InitMain(logger)
		if err.Error() != "OMADA_SHARED_SECRET environment variable is required" {
			logger.Fatalf("Failed test whether OMADA_SHARED_SECRET is required; log is `%v`", buf.String())
		}
	})

	os.Setenv("OMADA_SHARED_SECRET", "foo")

	t.Run("Can initialise after environment variables are set", func(t *testing.T) {
		buf.Reset()

		ntfyClient, server, port, err := main.InitMain(logger)

		if err != nil {
			logger.Fatalf("Still failed to initialize main; log is %v", buf.String())
		}

		if ntfyClient.NtfyURL != "https://ntfy.sh" {
			logger.Fatalf("Failed to initialize ntfy client properly; NtfyURL is `%v`", ntfyClient.NtfyURL)
		}

		if ntfyClient.Topic != "my_omada_alerts" {
			logger.Fatalf("Failed to initialize ntfy client properly; Topic is `%v`", ntfyClient.Topic)
		}

		if port != "8080" {
			logger.Fatalf("Failed to initialize server port; PORT is `%v`", port)
		}

		if server == nil {
			logger.Fatal("The server wasn't created by the init call")
		}
	})
}
