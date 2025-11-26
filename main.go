package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type CookieData struct {
	Name     string
	Value    string
	Expires  string
	HttpOnly bool
	Secure   bool
	SameSite string
	Path     string
}

func main() {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set cookies if they don't exist
		setCookieIfNotExists(w, "regular", "regular-value", false, false, http.SameSiteLaxMode, "/")
		setCookieIfNotExists(w, "httpOnly", "http-only-value", true, false, http.SameSiteLaxMode, "/")
		setCookieIfNotExists(w, "secure", "secure-value", false, true, http.SameSiteLaxMode, "/")
		setCookieIfNotExists(w, "samesite-lax", "samesite-lax-value", false, false, http.SameSiteLaxMode, "/")
		setCookieIfNotExists(w, "samesite-strict", "samesite-strict-value", false, false, http.SameSiteStrictMode, "/")
		setCookieIfNotExists(w, "samesite-none", "samesite-none-value", false, true, http.SameSiteNoneMode, "/") // Requires Secure

		// Get all cookies
		cookies := r.Cookies()
		var cookieData []CookieData

		for _, cookie := range cookies {
			cookieData = append(cookieData, CookieData{
				Name:     cookie.Name,
				Value:    cookie.Value,
				Expires:  cookie.Expires.Format(time.RFC1123),
				HttpOnly: cookie.HttpOnly,
				Secure:   cookie.Secure,
				SameSite: sameSiteString(cookie.SameSite),
				Path:     cookie.Path,
			})
		}

		tmpl.Execute(w, cookieData)
	})

	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func setCookieIfNotExists(w http.ResponseWriter, name, value string, httpOnly, secure bool, sameSite http.SameSite, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: sameSite,
		Path:     path,
	})
}

func sameSiteString(s http.SameSite) string {
	switch s {
	case http.SameSiteNoneMode:
		return "None"
	case http.SameSiteLaxMode:
		return "Lax"
	case http.SameSiteStrictMode:
		return "Strict"
	default:
		return "Default"
	}
}
