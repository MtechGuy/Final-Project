package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type User struct {
    Name     string
    Age      int
    Email    string
    Password string
}

func main() {
    http.HandleFunc("/", setCookie)
    http.HandleFunc("/get", getCookie)
	fmt.Println("Listening on port 4000...")
    http.ListenAndServe(":4000", nil)
}

func setCookie(w http.ResponseWriter, r *http.Request) {
    user := User{
        Name:     "Alex",
        Age:      30,
        Email:    "alex@example.com",
        Password: "password123",
    }
    cookieValue, err := json.Marshal(user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    cookie := &http.Cookie{
        Name:     "user",
        Value:    string(cookieValue),
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
    }
    http.SetCookie(w, cookie)
    fmt.Fprintf(w, "Cookie is set")
}

func getCookie(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("user")
    if err != nil {
        if err == http.ErrNoCookie {
            fmt.Fprint(w, "Cookie not found")
            return
        }
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    var user User
    err = json.Unmarshal([]byte(cookie.Value), &user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Name: %s\nAge: %d\nEmail: %s\nPassword: %s", user.Name, user.Age, user.Email, user.Password)
}
