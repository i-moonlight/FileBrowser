package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/users"
)

var ctx = context.Background()

type EncryptedCredentials struct {
	Iv            string `json:"iv"`
	EncryptedData string `json:"encryptedData"`
}

type DecryptedCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
	OU       string `json:"OU"`
	Hostname string `json:"hostname"`
}

type userInfo struct {
	Locale               string               `json:"locale"`
	ViewMode             users.ViewMode       `json:"viewMode"`
	SingleClick          bool                 `json:"singleClick"`
	Perm                 users.Permissions    `json:"perm"`
	Commands             []string             `json:"commands"`
	LockPassword         bool                 `json:"lockPassword"`
	HideDotfiles         bool                 `json:"hideDotfiles"`
	DateFormat           bool                 `json:"dateFormat"`
	Scope                string               `json:"scope"`
	EncryptedCredentials EncryptedCredentials `json:"credentials"`
}

type authToken struct {
	User userInfo `json:"user"`
	jwt.RegisteredClaims
}

type RedisTokenInfo struct {
	Payload  authToken `json:"authToken"`
	IsActive bool      `json:"isActive"`
}

const (
	ListViewMode   ViewMode = "list"
	MosaicViewMode ViewMode = "mosaic"
)

type Permissions struct {
	Admin    bool `json:"admin"`
	Execute  bool `json:"execute"`
	Create   bool `json:"create"`
	Rename   bool `json:"rename"`
	Modify   bool `json:"modify"`
	Delete   bool `json:"delete"`
	Share    bool `json:"share"`
	Download bool `json:"download"`
}

type tokenStruct struct {
	Scope                string               `json:"scope"`
	Locale               string               `json:"locale"`
	ViewMode             ViewMode             `json:"viewMode"`
	Perm                 Permissions          `json:"perm"`
	Fs                   afero.Fs             `json:"-" yaml:"-"`
	HideDotfiles         bool                 `json:"hideDotfiles"`
	EncryptedCredentials EncryptedCredentials `json:"credentiald"`
}

type extractor []string

func (e extractor) ExtractToken(r *http.Request) (string, error) {
	token, _ := request.HeaderExtractor{"X-Auth"}.ExtractToken(r)

	// Checks if the token isn't empty and if it contains two dots.
	// The former prevents incompatibility with URLs that previously
	// used basic auth.
	if token != "" && strings.Count(token, ".") == 2 {
		return token, nil
	}

	auth := r.URL.Query().Get("auth")
	if auth != "" && strings.Count(auth, ".") == 2 {
		return auth, nil
	}

	if r.Method == http.MethodGet {
		cookie, _ := r.Cookie("auth")
		if cookie != nil && strings.Count(cookie.Value, ".") == 2 {
			return cookie.Value, nil
		}
	}

	return "", request.ErrNoTokenInRequest
}

func withUser(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return d.settings.Key, nil
		}

		var tk authToken
		token, err := request.ParseFromRequest(r, &extractor{}, keyFunc, request.WithClaims(&tk))

		// Check is token valid
		if err != nil || !token.Valid {
			return http.StatusUnauthorized, nil
		}

		// Check token expiration
		expired := !tk.VerifyExpiresAt(time.Now(), true)

		if expired {
			return http.StatusUnauthorized, nil
		}

		// Check session in Redis
		val, err := d.redis.Get(ctx, token.Raw).Result()
		if err != nil {
			return http.StatusUnauthorized, nil
		}

		var rTokenInfo RedisTokenInfo
		json.Unmarshal([]byte(val), &rTokenInfo)

		if !rTokenInfo.IsActive {
			return http.StatusUnauthorized, nil
		}

		scope := filepath.Join(d.server.Root, filepath.Join("/", tk.User.Scope)) //nolint:gocritic
		fs := afero.NewBasePathFs(afero.NewOsFs(), scope)

		tokenPayload := &tokenStruct{
			Scope:                tk.User.Scope,
			Locale:               tk.User.Locale,
			ViewMode:             ViewMode(tk.User.ViewMode),
			Perm:                 Permissions(tk.User.Perm),
			Fs:                   fs,
			HideDotfiles:         tk.User.HideDotfiles,
			EncryptedCredentials: tk.User.EncryptedCredentials,
		}

		d.token = tokenPayload
		return fn(w, r, d)
	}
}

var checkTokenHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	return http.StatusOK, nil
})

var mountHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	// Decrypt credentials data
	decryptedCredentials, err := decryptData(d.token.EncryptedCredentials.EncryptedData, d.server.TokenCredentialsSecret, d.token.EncryptedCredentials.Iv)
	if err != nil {
		return http.StatusUnauthorized, nil
	}

	jsonString := string(decryptedCredentials)

	var credentials DecryptedCredentials
	json.Unmarshal([]byte(jsonString), &credentials)

	fmt.Println("Decrypted Credentials:", credentials)

	e := executeScript("./script.sh", credentials.Username, credentials.Password, credentials.Type, credentials.OU, credentials.Hostname)
	if e != nil {
		fmt.Println("Error executing script:", e)
		return http.StatusBadRequest, e
	}

	fmt.Println("Script executed successfully.")
	return http.StatusOK, nil
})
