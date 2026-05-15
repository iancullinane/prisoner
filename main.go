/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/iancullinane/prisoner/cmd"
	"github.com/iancullinane/prisoner/cmd/api"
)

func main() {

	runServer := flag.Bool("server", false, "set to run the server, otherwise execute cli")

	flag.Parse()

	if *runServer {

		server := api.NewPlayerServer(api.NewInMemoryPlayerStore())
		log.Fatal(http.ListenAndServe(":5001", server))
	}

	cmd.Execute()
}
