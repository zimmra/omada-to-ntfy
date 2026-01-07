package webhook

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/zimmra/omada-to-ntfy/ntfy"
	"github.com/zimmra/omada-to-ntfy/omada"
)

type WebhookServer struct {
	NtfyClient   ntfy.NtfyClient
	SharedSecret string
	Logger       *log.Logger
}

func (ws *WebhookServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if r.Header["Access_token"] == nil || r.Header["Access_token"][0] != ws.SharedSecret {
		http.Error(w, "Not authorized", http.StatusForbidden)
		return
	}

	omadaMessage, err := omada.ParseOmadaMessage(ws.Logger, body)
	if err != nil || omadaMessage == nil {
		ws.Logger.Printf("Error parsing Omada notification message: %v", err)
		http.Error(w, "Internal message parsing error", http.StatusInternalServerError)
		return
	}

	// Send the message to ntfy
	err = ws.NtfyClient.Send(omadaMessage)

	if err != nil {
		ws.Logger.Printf("Error sending message to ntfy: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "") // or something like: "Webhook forwarded successfully" (Omada doesn't care though)
}

// EOF
