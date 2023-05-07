package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type User struct {
	Name    string
	Age     int
	Email   string
	Address string
}

func main() {
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
    fmt.Println("Listening on port 4000...")
	http.ListenAndServe(":4000", nil)
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new user.
	user := User{
		Name:    "Alex Peraza",
		Age:     20,
		Email:   "alex.peraza@example.com",
		Address: "2 Miles George Price Highway, Belize City",
	}

	// Convert the user to a JSON string.
	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// URL-encode the JSON string and create a new cookie with the encoded value.
	encodedValue := url.QueryEscape(string(data))
	cookie := http.Cookie{
		Name:    "user",
		Value:   encodedValue,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}

	// Set the cookie.
	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "User cookie set successfully!")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user cookie.
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// URL-decode the cookie value and convert it to a user struct.
	decodedValue, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var user User
	err = json.Unmarshal([]byte(decodedValue), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Output the user's information.
	fmt.Fprintf(w, "Name: %s\n", user.Name)
	fmt.Fprintf(w, "Age: %d\n", user.Age)
	fmt.Fprintf(w, "Email: %s\n", user.Email)
	fmt.Fprintf(w, "Address: %s\n", user.Address)
}
