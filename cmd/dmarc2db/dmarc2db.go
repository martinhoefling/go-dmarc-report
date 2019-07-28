package main

import (
	"fmt"
	"github.com/martinhoefling/go-dmarc-report/acquire"
	"log"
	"os"
)

func getEnvOrDefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}

func getEnvOrDie(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s not found", key))
	}
	return value
}

func main() {
	login := getEnvOrDie("IMAP_LOGIN")
	password := getEnvOrDie("IMAP_PASSWORD")
	server := getEnvOrDie("IMAP_SERVER")
	mailbox := getEnvOrDefault("IMAP_MAILBOX", "INBOX")
	dburl := getEnvOrDie("DATABASE_URL")
	dbpass := getEnvOrDie("DATABASE_PASSWORD")

	err := acquire.IngestMissingAttachments(server, login, password, mailbox, dburl, dbpass)

	if err != nil {
		log.Fatal(err)
	}
}
