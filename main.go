package main

import (
	"flag"
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
	port := flag.Int("port", 8080, "Port for the web server")
	flag.Parse()

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

	addr := fmt.Sprintf(":%d", *port)

	fmt.Printf("Server starting on http://localhost%s\n", addr)
	http.ListenAndServe(addr, nil)
}

func setCookieIfNotExists(w http.ResponseWriter, name, value string, httpOnly, secure bool, sameSite http.SameSite, path string) {

	expireTime := time.Now().Add(365 * 24 * time.Hour)

	readableExpireTime := expireTime.Format("2006-01-02 15:04:05")

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value + " | expires: " + readableExpireTime,
		Expires:  expireTime,
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
