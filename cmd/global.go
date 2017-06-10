package cmd

import (
	"log"
	"os"

	"github.com/rumyantseva/go-velobike/velobike"
)

func mustClientFromEnv() *velobike.Client {
	id := os.Getenv("VELOBIKE_ID")
	if id == "" {
		log.Fatal("VELOBIKE_ID env var is not set")
	}
	pass := os.Getenv("VELOBIKE_PASS")
	if pass == "" {
		log.Fatal("VELOBIKE_PASS env var is not set")
	}

	tp := velobike.BasicAuthTransport{
		Username: id,
		Password: pass,
	}
	client := velobike.NewClient(tp.Client())
	auth, _, err := client.Authorization.Authorize()
	if err != nil {
		log.Fatalf("failed to authorize with velobike: %v", err)
	}
	client.SessionId = auth.SessionId

	return client
}
