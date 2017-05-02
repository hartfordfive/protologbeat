package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/hartfordfive/protologbeat/beater"
)

func main() {
	err := beat.Run("protologbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
