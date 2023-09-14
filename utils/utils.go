package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/filebrowser/filebrowser/v2/users"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func generateKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func DecryptData(encryptedData string, key string, iv string) ([]byte, error) {
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

func ExecuteScript(scriptPath string, args ...string) error {
	fmt.Println(args)
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

func SubscribeRedisEvent(rdb *redis.Client, tokenCredentialsSecret string, tokenSecret string, mountScriptPath string) {
	pubsub := rdb.Subscribe(ctx, "__keyevent@0__:expired")
	// defer pubsub.Close()
	// Wait for the subscription to become ready
	_, err := pubsub.Receive(ctx)
	if err != nil {
		fmt.Println("Failed to subscribe to key expiration events:", err)
	}

	fmt.Println("Subscribed to key expiration events.")

	messageListener(pubsub, tokenSecret, tokenCredentialsSecret, mountScriptPath)
}

func messageListener(pubsub *redis.PubSub, tokenSecret string, tokenCredentialsSecret string, mountScriptPath string) {
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println("Error receiving message:", err)
		}

		fmt.Println(msg.Payload)

		tokenClaims := parseToken(msg.Payload, tokenSecret, tokenCredentialsSecret)
		decryptedCredentials := parseCredentials(tokenClaims.User.EncryptedCredentials.EncryptedData, tokenClaims.User.EncryptedCredentials.Iv, tokenCredentialsSecret)
		e := ExecuteScript(mountScriptPath, decryptedCredentials.Username, decryptedCredentials.Password, decryptedCredentials.Type, "0", decryptedCredentials.Hostname)
		if e != nil {
			fmt.Println("Error executing script:", e)
		}
		fmt.Println("Script executed successfully (unmount).")
	}
}

func parseToken(tokenString string, tokenSecret string, tokenCredentialsSecret string) *users.AuthToken {
	token, _ := jwt.ParseWithClaims(tokenString, &users.AuthToken{}, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})

	return token.Claims.(*users.AuthToken)
}

func parseCredentials(encData string, encIv string, tokenCredentialsSecret string) users.DecryptedCredentials {
	decryptedCredentials, _ := DecryptData(encData, tokenCredentialsSecret, encIv)
	jsonString := string(decryptedCredentials)

	var credentials users.DecryptedCredentials
	json.Unmarshal([]byte(jsonString), &credentials)

	return credentials
}
