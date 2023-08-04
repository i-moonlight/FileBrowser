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

// const (
// 	TokenExpirationTime = time.Hour * 2
// )

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
			Locale:               tk.User.Scope,
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

// func withAdmin(fn handleFunc) handleFunc {
// 	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
// 		if !d.user.Perm.Admin {
// 			return http.StatusForbidden, nil
// 		}

// 		return fn(w, r, d)
// 	})
// }

// var loginHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
// 	auther, err := d.store.Auth.Get(d.settings.AuthMethod)
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	user, err := auther.Auth(r, d.store.Users, d.settings, d.server)
// 	if err == os.ErrPermission {
// 		return http.StatusForbidden, nil
// 	} else if err != nil {
// 		return http.StatusInternalServerError, err
// 	} else {
// 		return printToken(w, r, d, user)
// 	}
// }

// type signupBody struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// var signupHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
// 	if !d.settings.Signup {
// 		return http.StatusMethodNotAllowed, nil
// 	}

// 	if r.Body == nil {
// 		return http.StatusBadRequest, nil
// 	}

// 	info := &signupBody{}
// 	err := json.NewDecoder(r.Body).Decode(info)
// 	if err != nil {
// 		return http.StatusBadRequest, err
// 	}

// 	if info.Password == "" || info.Username == "" {
// 		return http.StatusBadRequest, nil
// 	}

// 	user := &users.User{
// 		Username: info.Username,
// 	}

// 	d.settings.Defaults.Apply(user)

// 	pwd, err := users.HashPwd(info.Password)
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	user.Password = pwd

// 	userHome, err := d.settings.MakeUserDir(user.Username, user.Scope, d.server.Root)
// 	if err != nil {
// 		log.Printf("create user: failed to mkdir user home dir: [%s]", userHome)
// 		return http.StatusInternalServerError, err
// 	}
// 	user.Scope = userHome
// 	log.Printf("new user: %s, home dir: [%s].", user.Username, userHome)

// 	err = d.store.Users.Save(user)
// 	if err == errors.ErrExist {
// 		return http.StatusConflict, err
// 	} else if err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	return http.StatusOK, nil
// }

// var renewHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
// 	return printToken(w, r, d, d.user)
// })

// func printToken(w http.ResponseWriter, _ *http.Request, d *data, user *users.User) (int, error) {
// 	claims := &authToken{
// 		User: userInfo{
// 			ID:           user.ID,
// 			Locale:       user.Locale,
// 			ViewMode:     user.ViewMode,
// 			SingleClick:  user.SingleClick,
// 			Perm:         user.Perm,
// 			LockPassword: user.LockPassword,
// 			Commands:     user.Commands,
// 			HideDotfiles: user.HideDotfiles,
// 			DateFormat:   user.DateFormat,
// 			Scope:        user.Scope,
// 		},
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpirationTime)),
// 			Issuer:    "File Browser",
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signed, err := token.SignedString(d.settings.Key)
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	w.Header().Set("Content-Type", "text/plain")
// 	if _, err := w.Write([]byte(signed)); err != nil {
// 		return http.StatusInternalServerError, err
// 	}
// 	return 0, nil
// }
