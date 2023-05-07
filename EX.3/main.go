package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Key used to sign the cookie value
const secretKey = "secret"

// Creates a signed cookie
func createSignedCookie(name, value string, w http.ResponseWriter) {
	// Create the cookie
	cookie := http.Cookie{
		Name:     name,
		Value:    signValue(value),
		Expires:  time.Now().Add(time.Hour * 24 * 30), // Expires in 30 days
		HttpOnly: true,
	}

	// Set the cookie header
	http.SetCookie(w, &cookie)
}

// Signs the cookie value using HMAC-SHA256
func signValue(value string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(value))
	signature := mac.Sum(nil)

	// Encode the signature in base64
	encodedSig := base64.RawURLEncoding.EncodeToString(signature)

	// Combine the value and signature
	return fmt.Sprintf("%s|%s", value, encodedSig)
}

// Verifies the signature of a signed cookie
func verifySignedCookie(name string, r *http.Request) (string, bool) {
	// Get the cookie
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", false
	}

	// Split the value and signature
	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 2 {
		return "", false
	}
	value := parts[0]
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", false
	}

	// Verify the signature
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(value))
	expectedSig := mac.Sum(nil)
	if !hmac.Equal(sig, expectedSig) {
		return "", false
	}

	return value, true
}

func main() {
	// Example usage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set the signed cookie
		createSignedCookie("mycookie", "Hello, World!", w)

		// Verify the signed cookie
		value, ok := verifySignedCookie("mycookie", r)
		if ok {
			fmt.Fprintf(w, "Verified cookie value: %s", value)
		} else {
			fmt.Fprint(w, "Invalid cookie")
		}
	})
	fmt.Println("Listening on port 4000...")
	http.ListenAndServe(":4000", nil)
}
