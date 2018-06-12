package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"
)

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

	log.Printf("crc token found: %v", crcToken)

	mac := hmac.New(sha256.New, []byte(twitterConsumerSecret))
	mac.Write([]byte(crcToken))

	w.Write([]byte(base64.StdEncoding.EncodeToString(mac.Sum(nil))))
}

func main() {
	http.HandleFunc("/"+webhookPath, crcCheck)
	appengine.Main()
}
