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
	"github.com/filebrowser/filebrowser/v2/utils"
)

var ctx = context.Background()

type RedisTokenInfo struct {
	Locale    string `json:"locale"`
	Scope     string `json:"scope"`
	IsActive  bool   `json:"isActive"`
	SessionId string `json:"sessionId"`
	UA        string `json:"ua"`
	IP        string `json:"ip"`
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

	return "", request.ErrNoTokenInRequest
}

func extractSessionId(r *http.Request) string {
	sessionId := r.Header.Get("X-Session-Id")

	if sessionId == "" {
		sessionId = r.URL.Query().Get("sid")
	}

	return sessionId
}

func getTokenInfoFromRedis(d *data, token string) (RedisTokenInfo, error) {
	val, err := d.redis.Get(ctx, token).Result()
	if err != nil {
		return RedisTokenInfo{}, err
	}

	var rTokenInfo RedisTokenInfo
	json.Unmarshal([]byte(val), &rTokenInfo)
	return rTokenInfo, nil
}

func extractIPAddress(r *http.Request) string {
	// Attempt to get the client's IP address from the X-Real-Ip header
	ipAddress := r.Header.Get("X-Real-Ip")
	// ipAddress := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"

	if ipAddress == "" {
		// If X-Real-Ip is not set, fall back to getting the X-Forwarded-For
		// address from the header.
		ipAddress = r.Header.Get("X-Forwarded-For")
	}

	if ipAddress == "" {
		// If X-Forwarded-For is not set, fall back to getting the remote
		// address from the request.
		ipAddress = r.RemoteAddr
	}

	// Map default ip if user is localhost (dev mode)
	isLocal := strings.Contains(ipAddress, "127.0.0.1") || strings.Contains(ipAddress, "::1")

	if isLocal {
		ipAddress = "localhost"
	}

	return ipAddress
}

func withUser(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return d.settings.Key, nil
		}

		var tk users.AuthToken
		token, err := request.ParseFromRequest(r, &extractor{}, keyFunc, request.WithClaims(&tk))
		sessionId := extractSessionId(r)
		userAgent := r.Header.Get("User-Agent")
		ipAddress := extractIPAddress((r))

		// Check if sessionId is not empty
		if sessionId == "" {
			return http.StatusUnauthorized, nil
		}

		// Check is token valid
		if err != nil || !token.Valid {
			return http.StatusUnauthorized, nil
		}

		// Check token expiration
		expired := !tk.VerifyExpiresAt(time.Now(), true)

		if expired {
			return http.StatusUnauthorized, nil
		}

		rTokenInfo, err := getTokenInfoFromRedis(d, token.Raw)
		if err != nil {
			return http.StatusUnauthorized, nil
		}

		// Set new Session Id into Redis if it is empty
		if rTokenInfo.SessionId == "" {
			rTokenInfo.SessionId = sessionId
			rTokenInfo.UA = userAgent
			rTokenInfo.IP = ipAddress
			jsonBytes, _ := json.Marshal(rTokenInfo)

			err := d.redis.Set(ctx, token.Raw, jsonBytes, -1).Err()
			if err != nil {
				fmt.Println("Error while updating session Id in Redis")
				fmt.Println(err)
				return http.StatusUnauthorized, nil
			}

			rTokenInfo, _ = getTokenInfoFromRedis(d, token.Raw)
		}

		// Compare sessionId from redis and request header
		if rTokenInfo.SessionId != sessionId {
			return http.StatusUnauthorized, nil
		}

		// Compare IP address from redis and users request
		if rTokenInfo.IP != ipAddress {
			return http.StatusUnauthorized, nil
		}

		// Compare User Agent from redis and request header
		if rTokenInfo.UA != userAgent {
			return http.StatusUnauthorized, nil
		}

		scope := filepath.Join(d.server.Root, filepath.Join("/", tk.User.Scope)) //nolint:gocritic
		fs := afero.NewBasePathFs(afero.NewOsFs(), scope)

		tokenPayload := &users.TokenStruct{
			Scope:                tk.User.Scope,
			Locale:               tk.User.Locale,
			ViewMode:             users.ViewMode(tk.User.ViewMode),
			Perm:                 users.Permissions(tk.User.Perm),
			Fs:                   fs,
			HideDotfiles:         tk.User.HideDotfiles,
			EncryptedCredentials: tk.User.EncryptedCredentials,
			Raw:                  token.Raw,
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
	decryptedCredentials, err := utils.DecryptData(d.token.EncryptedCredentials.EncryptedData, d.server.TokenCredentialsSecret, d.token.EncryptedCredentials.Iv)
	if err != nil {
		return http.StatusUnauthorized, nil
	}

	jsonString := string(decryptedCredentials)

	var credentials users.DecryptedCredentials
	json.Unmarshal([]byte(jsonString), &credentials)

	fmt.Println("Script Path:", d.server.MountScriptPath)

	e := utils.ExecuteScript(d.server.MountScriptPath, credentials.Username, credentials.Password, credentials.Type, "1", credentials.Hostname)
	if e != nil {
		fmt.Println("Error executing script:", e)
		return http.StatusBadRequest, e
	}

	fmt.Println("Script executed successfully.")
	return http.StatusOK, nil
})

var logoutHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	// Decrypt credentials data
	decryptedCredentials, err := utils.DecryptData(d.token.EncryptedCredentials.EncryptedData, d.server.TokenCredentialsSecret, d.token.EncryptedCredentials.Iv)
	if err != nil {
		return http.StatusUnauthorized, nil
	}

	jsonString := string(decryptedCredentials)

	var credentials users.DecryptedCredentials
	json.Unmarshal([]byte(jsonString), &credentials)

	fmt.Println("Script Path:", d.server.MountScriptPath)

	e := utils.ExecuteScript(d.server.MountScriptPath, credentials.Username, credentials.Password, credentials.Type, "0", credentials.Hostname)
	if e != nil {
		fmt.Println("Error executing script:", e)
		return http.StatusBadRequest, e
	}

	fmt.Println("Script executed successfully (unmount).")

	d.redis.Del(ctx, d.token.Raw)
	return http.StatusOK, nil
})
