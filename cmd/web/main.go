package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dmathieu/sabayon/heroku"
	"github.com/joho/godotenv"
)

var (
	token   string
	appName string
)

func status(w http.ResponseWriter, r *http.Request) {
	herokuClient := heroku.NewClient(nil, token)
	certificates, err := herokuClient.GetSSLCertificates(appName)
	if err != nil {
		log.Fatal(err)
	}

	if len(certificates) > 1 {
		msg := fmt.Sprintf("Found %d certificate. Can only update one. Nothing done.", len(certificates))
		io.WriteString(w, msg)
		log.Println(msg)
	}

	if len(certificates) != 0 {
		certExpiration, err := time.Parse(time.RFC3339, certificates[0].SslCert.ExpiresAt)
		if err != nil {
			log.Fatal(err)
		}
		now := time.Now()
		renew := certExpiration.AddDate(0, -1, 0)

		if now.Before(renew) {
			msg := fmt.Sprintf("expires_at=\"%s\" renew_at=\"%s\"", certExpiration, renew)
			io.WriteString(w, msg)
			log.Println(msg)
		}
	}
}

func main() {
	godotenv.Load()

	token = os.Getenv("HEROKU_TOKEN")
	appName = os.Getenv("ACME_APP_NAME")

	http.HandleFunc("/", status)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Listening on port %s", port)
	http.ListenAndServe(":"+port, nil)
}
