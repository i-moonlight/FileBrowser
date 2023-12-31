package settings

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/filebrowser/filebrowser/v2/rules"
)

const DefaultUsersHomeBasePath = "/users"

// Settings contain the main settings of the application.
type Settings struct {
	Key              []byte              `json:"key"`
	CreateUserDir    bool                `json:"createUserDir"`
	UserHomeBasePath string              `json:"userHomeBasePath"`
	Defaults         UserDefaults        `json:"defaults"`
	Branding         Branding            `json:"branding"`
	Commands         map[string][]string `json:"commands"`
	Shell            []string            `json:"shell"`
	Rules            []rules.Rule        `json:"rules"`
}

// GetRules implements rules.Provider.
func (s *Settings) GetRules() []rules.Rule {
	return s.Rules
}

// Server specific settings.
type Server struct {
	Root                   string `json:"root"`
	BaseURL                string `json:"baseURL"`
	Socket                 string `json:"socket"`
	TLSKey                 string `json:"tlsKey"`
	TLSCert                string `json:"tlsCert"`
	Port                   string `json:"port"`
	Address                string `json:"address"`
	Log                    string `json:"log"`
	EnableThumbnails       bool   `json:"enableThumbnails"`
	ResizePreview          bool   `json:"resizePreview"`
	EnableExec             bool   `json:"enableExec"`
	TypeDetectionByHeader  bool   `json:"typeDetectionByHeader"`
	AuthHook               string `json:"authHook"`
	RedisUrl               string `json:"redisUrl"`
	TokenSecret            string `json:"tokenSecret"`
	TokenCredentialsSecret string `json:"tokenCredentialsSecret"`
	MountScriptPath        string `json:"mountScriptPath"`
}

// Clean cleans any variables that might need cleaning.
func (s *Server) Clean() {
	s.BaseURL = strings.TrimSuffix(s.BaseURL, "/")
}

// GenerateKey generates a key of 512 bits.
func GenerateKey() ([]byte, error) {
	b := make([]byte, 64) //nolint:gomnd
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func IsValidKey(keyString string) ([]byte, bool) {
	decodedKey, err := base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		return decodedKey, false
	}

	return decodedKey, len(decodedKey) == 64
}
