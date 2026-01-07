package webhook_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zimmra/omada-to-ntfy/ntfy"
	"github.com/zimmra/omada-to-ntfy/omada"
	"github.com/zimmra/omada-to-ntfy/webhook"
)

// Mock implementation for ntfy client to avoid making real HTTP requests

type NtfyClientMock struct {
	Calls       int
	returnError error
}

func (mock *NtfyClientMock) Send(payload *omada.OmadaMessage) error {
	mock.Calls += 1
	return mock.returnError
}

// Here's an io.Reader implementation that returns an error on reading

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

// Now for the tests

func TestWebhookServer(t *testing.T) {

	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)

	const sharedSecret = "vewySecwet"

	ntfyClient := ntfy.NtfyClient{
		NtfyURL:  "https://ntfy.sh",
		Topic:    "test_topic",
		Username: "",
		Password: "",
		Logger:   logger,
	}

	server := &webhook.WebhookServer{
		NtfyClient:   ntfyClient,
		SharedSecret: sharedSecret,
		Logger:       logger,
	}

	notAuthorizedTests := []struct {
		name  string
		token string
	}{
		{
			name:  "Not authenticated with an empty access token header",
			token: "",
		},
		{
			name:  "Not authenticated with a wrong access token header",
			token: server.SharedSecret + "_",
		},
		{
			name:  "Not authenticated with an all uppercase token header",
			token: strings.ToUpper(server.SharedSecret),
		},
		{
			name:  "Not authenticated with an all lowercase token header",
			token: strings.ToLower(server.SharedSecret),
		},
	}

	for _, tt := range notAuthorizedTests {
		t.Run(tt.name, func(t *testing.T) {
			json := []byte(`{"Site":"Some site","description":"This is a webhook message from Omada Controller","shardSecret":"fef97b18-e440-45bc-8826-be957e4dc8f6","text":["[2.5G WAN1] of [gateway:98-03-8E-3A-8D-53] is down.\r","[gateway:98-03-8E-3A-8D-53]: The online detection result of [2.5G WAN1] was offline.\r"],"Controller":"Omada Controller_347044","timestamp":1758852904877}`)
			body := bytes.NewReader(json)

			request, _ := http.NewRequest(http.MethodPost, "/", body)

			t.Logf("Using `%v` for the access token for this test", tt.token)
			request.Header.Set("Access_token", tt.token)

			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			got := response.Result().Status
			want := "403 Forbidden"

			if got != want {
				t.Errorf("Expected status code to be `%s`, but got `%s`", want, got)
			} else {
				t.Logf("Got the expected %v status code", got)
			}

			got = response.Body.String()
			want = "Not authorized\n"

			if got != want {
				t.Errorf("Got %q, want %q", got, want)
			}
		})
	}

	t.Run("Authenticated but incorrect JSON input", func(t *testing.T) {
		json := []byte(`{"Site":"Some site",`)
		body := bytes.NewReader(json)

		request, _ := http.NewRequest(http.MethodPost, "/", body)

		request.Header.Set("Access_token", server.SharedSecret) // CORRECT

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Result().Status
		want := "500 Internal Server Error"

		if got != want {
			t.Errorf("Expected status code to be `%s`, but got `%s`", want, got)
		}

		got = response.Body.String()
		want = "Internal message parsing error\n"

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})

	t.Run("Error on reading the data from the request", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", errReader(0))

		request.Header.Set("Access_token", server.SharedSecret) // CORRECT

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Result().Status
		want := "400 Bad Request"

		if got != want {
			t.Errorf("Expected status code to be `%s`, but got `%s`", want, got)
		}

		got = response.Body.String()
		want = "Bad Request\n"

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})
}
