package ntfy

import (
	"testing"

	"github.com/zimmra/omada-to-ntfy/omada"
)

func TestMapPriority(t *testing.T) {
	tests := []struct {
		name          string
		omadaPriority int
		expected      int
	}{
		{"Max priority", 10, 5},
		{"High priority", 7, 4},
		{"Default priority", 4, 3},
		{"Low priority", 0, 2},
		{"Unknown priority", 5, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapPriority(tt.omadaPriority)
			if result != tt.expected {
				t.Errorf("MapPriority(%d) = %d; want %d", tt.omadaPriority, result, tt.expected)
			}
		})
	}
}

func TestGetTagsForMessageType(t *testing.T) {
	tests := []struct {
		name        string
		messageType omada.OmadaMessageType
		expected    []string
	}{
		{"Offline message", omada.OmadaOfflineMessage, []string{"rotating_light"}},
		{"Online message", omada.OmadaOnlineMessage, []string{"white_check_mark"}},
		{"Test message", omada.OmadaTestMessage, []string{"test_tube"}},
		{"Unrecognised message", omada.UnrecognisedMessage, []string{"warning"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTagsForMessageType(tt.messageType)
			if len(result) != len(tt.expected) {
				t.Errorf("GetTagsForMessageType(%v) returned %d tags; want %d", tt.messageType, len(result), len(tt.expected))
				return
			}
			for i, tag := range result {
				if tag != tt.expected[i] {
					t.Errorf("GetTagsForMessageType(%v)[%d] = %s; want %s", tt.messageType, i, tag, tt.expected[i])
				}
			}
		})
	}
}

// EOF
