package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	cookieName  = "my_cookie"
	cookieKey   = "my_secret_key_123"
	cookieValue = "my_cookie_value"
)

func main() {
	http.HandleFunc("/", setCookieHandler)
	fmt.Println("Listening on port 4000...")
	http.ListenAndServe(":4000", nil)
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := createCookie(cookieName, cookieValue, cookieKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	fmt.Fprintln(w, "Cookie set")
}

func createCookie(name, value, key string) (*http.Cookie, error) {
	// Encrypt the cookie value
	encryptedValue, err := encrypt(value, key)
	if err != nil {
		return nil, err
	}

	// Hash the encrypted value to create a signature
	signature := hash(encryptedValue, key)

	// Combine the encrypted value and signature into a single string
	cookieValue := base64.StdEncoding.EncodeToString(encryptedValue) + "|" + base64.StdEncoding.EncodeToString(signature)

	// Create the cookie
	cookie := &http.Cookie{
		Name:     name,
		Value:    cookieValue,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Hour * 24),
	}

	return cookie, nil
}

func encrypt(text string, key string) ([]byte, error) {
	// Convert the key to a 32-byte slice using SHA-256
	keyBytes := hash([]byte(key), "")[:32]

	// Generate a random 16-byte initialization vector
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Create a new AES cipher using the key and initialization vector
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	// Add padding to the plaintext to ensure it's a multiple of the block size
	paddedText := pad([]byte(text), aes.BlockSize)

	// Create a new cipher.StreamWriter to encrypt the plaintext
	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(paddedText))
	stream.XORKeyStream(ciphertext, paddedText)

	// Prepend the initialization vector to the ciphertext
	ciphertextWithIV := append(iv, ciphertext...)

	return ciphertextWithIV, nil
}

func hash(data []byte, key string) []byte {
	// Combine the data and key into a single byte slice
	message := append(data, []byte(key)...)

	// Hash the message using SHA-256
	hash := sha256.Sum256(message)

	return hash[:]
}

func pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	padding := strings.Repeat(string(rune(padLen)), padLen)

	return append(data, []byte(padding)...)
}
