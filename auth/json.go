package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/filebrowser/filebrowser/v2/settings"
)

// MethodJSONAuth is used to identify json auth.
const MethodJSONAuth settings.AuthMethod = "json"

// JSONAuth is a json implementation of an Auther.
type JSONAuth struct {
	ReCaptcha *ReCaptcha `json:"recaptcha" yaml:"recaptcha"`
}

// LoginPage tells that json auth doesn't require a login page.
func (a JSONAuth) LoginPage() bool {
	return true
}

const reCaptchaAPI = "/recaptcha/api/siteverify"

// ReCaptcha identifies a recaptcha connection.
type ReCaptcha struct {
	Host   string `json:"host"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

// Ok checks if a reCaptcha responde is correct.
func (r *ReCaptcha) Ok(response string) (bool, error) {
	body := url.Values{}
	body.Set("secret", r.Secret)
	body.Add("response", response)

	client := &http.Client{}

	resp, err := client.Post(
		r.Host+reCaptchaAPI,
		"application/x-www-form-urlencoded",
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var data struct {
		Success bool `json:"success"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false, err
	}

	return data.Success, nil
}
