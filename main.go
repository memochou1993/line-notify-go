package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/memochou1993/line-notify/app"
	"os"

	"html/template"
	"log"
	"net/http"
	"net/url"
)

var (
	clientID     string
	clientSecret string
	callbackURL  string
	token        string
)

func main() {
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	callbackURL = os.Getenv("CALLBACK_URL")

	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/notify", notifyHandler)
	http.HandleFunc("/auth", authHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), nil))
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
	}

	code := r.Form.Get("code")

	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", callbackURL)
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)

	payload, err := app.Call("POST", "https://notify-bot.line.me/oauth/token", data, "")

	if err != nil {
		log.Println(err.Error())
	}

	res := app.Parse(payload)

	token = res.AccessToken

	if _, err := w.Write(payload); err != nil {
		log.Println(err.Error())
	}
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
	}

	msg := r.Form.Get("msg")

	data := url.Values{}
	data.Add("message", msg)

	payload, err := app.Call("POST", "https://notify-bot.line.me/api/notify", data, token)

	if err != nil {
		log.Println(err.Error())
	}

	res := app.Parse(payload)

	token = res.AccessToken

	if _, err := w.Write(payload); err != nil {
		log.Println(err.Error())
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/auth.html"))

	err := tmpl.Execute(w, struct {
		ClientID    string
		CallbackURL string
	}{
		ClientID:    clientID,
		CallbackURL: callbackURL,
	})

	if err != nil {
		log.Println(err.Error())
	}
}
