package main

import (
	"flag"
	"log"
	"os"

	"github.com/howeyc/gopass"
	"github.com/martinhoefling/go-dmarc-report/acquire"
)

func getPassword(username, server string) (password string) {
	password = os.Getenv("IMAP_PASSWORD")

	if password == "" {
		log.Printf("Enter IMAP Password for %v on %v: ", username, server)
		passwordBytes, err := gopass.GetPasswd()
		if err != nil {
			panic(err)
		}
		password = string(passwordBytes)
	}
	return
}

func main() {
	var server, username, mailbox string
	var err error
	flag.StringVar(&server, "server", "", "Mail server to use")
	flag.StringVar(&username, "username", "", "Username for logging into the mail server")
	flag.StringVar(&mailbox, "mailbox", "", "Mailbox to read messages from")
	flag.Parse()

	password := getPassword(username, server)

	err = acquire.DownloadMissingAttachments(server, username, password, mailbox)

	if err != nil {
		log.Fatal(err)
	}
}
