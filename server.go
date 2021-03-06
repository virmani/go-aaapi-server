package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"
)

type crcResponse struct {
	ResponseToken string `json:"response_token"`
}

const crcTokenParam = "crc_token"
const webhookPath = "webhooks"

var twitterConsumerKey = os.Getenv("TWITTER_CONSUMER_KEY")
var twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")

func crcCheck(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	crcToken := params.Get(crcTokenParam)

	if crcToken == "" {
		log.Printf("No %v found in the query params", crcTokenParam)
		http.Error(w, "crc_token is not specified in the request query params",
			http.StatusBadRequest)
		return
	}

	mac := hmac.New(sha256.New, []byte(twitterConsumerSecret))
	mac.Write([]byte(crcToken))

	resp := crcResponse{"sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))}

	js, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	http.HandleFunc("/"+webhookPath, crcCheck)
	appengine.Main()
}
