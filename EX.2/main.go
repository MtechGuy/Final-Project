package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func main() {
	// Set up the cookie to be encoded
	cookieName := "mycookie"
	cookieValue := "Hello, World! This is a special character: ü"
	maxCookieLength := 256 // Maximum length in bytes

	// Encode the cookie value
	encodedValue := url.QueryEscape(cookieValue)

	// Check if the encoded cookie value is too long
	if len(encodedValue) > maxCookieLength {
		fmt.Printf("Encoded cookie value is too long (max: %d bytes)\n", maxCookieLength)
		return
	}

	// Set up the cookie with the encoded value
	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   encodedValue,
		Expires: time.Now().Add(24 * time.Hour),
	}

	// Set the cookie on a response
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, cookie)
		fmt.Fprintln(w, "Cookie set!")
	})

	// Start the server
	fmt.Println("Listening on port 4000...")
	http.ListenAndServe(":4000", nil)
}
