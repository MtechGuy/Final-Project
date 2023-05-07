// Filename: main.go
package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/set", setCookieHandler)
	mux.HandleFunc("/get", getCookieHandler)

	log.Print("Starting server on :3000")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "pepto",
		Value:    "howdy!",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("the cookie has been set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("pepto")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "a cookie was not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	w.Write([]byte(cookie.Value))
}

func Write(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))
	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}
	http.SetCookie(w, &cookie)
	return nil
}

func Read(r *http.Request, name string) (string, error) {
	// read the cookie
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err

	}
	value, err := base64.URLEncoding.DecodeString(cookie.Value)

	if err != nil {
		return "", ErrInvalidValue
	}
	return string(value), nil
}
