package cmd

import (
	"os"
	"urlshortener/cmd/urlshortenerapi"
)

func Execute() error {
	// There is one command application
	return urlshortenerapi.Execute(os.Args)
}
