package http

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"

	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	libErrors "github.com/filebrowser/filebrowser/v2/errors"
)

func renderJSON(w http.ResponseWriter, _ *http.Request, data interface{}) (int, error) {
	marsh, err := json.Marshal(data)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(marsh); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func errToStatus(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case os.IsPermission(err):
		return http.StatusForbidden
	case os.IsNotExist(err), err == libErrors.ErrNotExist:
		return http.StatusNotFound
	case os.IsExist(err), err == libErrors.ErrExist:
		return http.StatusConflict
	case errors.Is(err, libErrors.ErrPermissionDenied):
		return http.StatusForbidden
	case errors.Is(err, libErrors.ErrInvalidRequestParams):
		return http.StatusBadRequest
	case errors.Is(err, libErrors.ErrRootUserDeletion):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// This is an addaptation if http.StripPrefix in which we don't
// return 404 if the page doesn't have the needed prefix.
func stripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" || prefix == "/" {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		h.ServeHTTP(w, r2)
	})
}

func generateKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func decryptData(encryptedData string, key string, iv string) ([]byte, error) {
	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	ciphertext, err := hex.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	byteKey := generateKey(key)

	block, err := aes.NewCipher(byteKey)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	ciphertext, err = unpadPKCS7(ciphertext)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// unpadPKCS7 removes PKCS#7 padding from the decrypted data.
func unpadPKCS7(data []byte) ([]byte, error) {
	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}

func executeScript(scriptPath string, args ...string) error {
	cmd := exec.Command(scriptPath, args...)

	// Set the working directory if the script requires it
	// cmd.Dir = "/path/to/working/directory"

	// Set environment variables if the script requires it
	// cmd.Env = append(os.Environ(), "KEY=VALUE")

	// Redirect the standard output and error to the current process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	return err
}
