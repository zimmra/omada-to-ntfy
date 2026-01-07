package ntfy

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zimmra/omada-to-ntfy/omada"
)

type NtfyClient struct {
	NtfyURL  string
	Topic    string
	Username string
	Password string
	Logger   *log.Logger
}

// MapPriority maps Omada priorities (0-10) to ntfy priorities (1-5)
// 10 -> 5 (Max/Urgent)
// 7  -> 4 (High)
// 4  -> 3 (Default)
// 0  -> 2 (Low)
func MapPriority(omadaPriority int) int {
	switch omadaPriority {
	case 10:
		return 5 // Max/Urgent
	case 7:
		return 4 // High
	case 4:
		return 3 // Default
	case 0:
		return 2 // Low
	default:
		return 3 // Default for unrecognized priorities
	}
}

// GetTagsForMessageType returns emoji tags for different message types
// These tags will be converted to emojis by ntfy
func GetTagsForMessageType(messageType omada.OmadaMessageType) []string {
	switch messageType {
	case omada.OmadaOfflineMessage:
		return []string{"rotating_light"} // 🚨
	case omada.OmadaOnlineMessage:
		return []string{"white_check_mark"} // ✅
	case omada.OmadaTestMessage:
		return []string{"test_tube"} // 🧪
	case omada.UnrecognisedMessage:
		return []string{"warning"} // ⚠️
	default:
		return []string{"information_source"} // ℹ️
	}
}

// Send sends a message to ntfy using the provided payload
func (nc *NtfyClient) Send(payload *omada.OmadaMessage) error {
	// Construct the full URL
	url := fmt.Sprintf("%s/%s", strings.TrimSuffix(nc.NtfyURL, "/"), nc.Topic)

	// Create the request body
	body := []byte(payload.Body())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		nc.Logger.Printf("Could not create ntfy request: %v", err)
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("Title", payload.Title())
	req.Header.Set("Priority", fmt.Sprintf("%d", MapPriority(payload.Priority())))

	// Add tags based on message type
	tags := GetTagsForMessageType(payload.Type())
	if len(tags) > 0 {
		req.Header.Set("Tags", strings.Join(tags, ","))
	}

	// Add authentication if provided
	if nc.Username != "" && nc.Password != "" {
		req.SetBasicAuth(nc.Username, nc.Password)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		nc.Logger.Printf("Could not send message to ntfy: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		nc.Logger.Printf("ntfy returned non-success status code: %d", resp.StatusCode)
		return fmt.Errorf("ntfy returned status code %d", resp.StatusCode)
	}

	nc.Logger.Println("Message sent to ntfy")
	return nil
}

// EOF
