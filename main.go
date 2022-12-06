package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/memochou1993/line-notify-go/app"

	"html/template"
	"log"
	"math/rand"
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

	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", r.Form.Get("code"))
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

	data := url.Values{}
	data.Add("message", r.Form.Get("message"))

	payload, err := app.Call("POST", "https://notify-api.line.me/api/notify", data, token)

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
		State       int
	}{
		ClientID:    clientID,
		CallbackURL: callbackURL,
		State:       rand.Int(),
	})

	if err != nil {
		log.Println(err.Error())
	}
}
