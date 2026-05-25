/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/iancullinane/prisoner/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// playCmd represents the play command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start an http server for the prisoner's dilemma",
	Long:  `Server will start a long running process serving HTTP requests`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storeKind := viper.GetString("store")
		store, cleanup, err := openPlayerStore(context.Background(), storeKind)
		if err != nil {
			return err
		}
		if cleanup != nil {
			defer cleanup()
		}

		historyStoreKind := viper.GetString("history-store")
		historyStore, cleanup, err := openHistoryStore(context.Background(), historyStoreKind)
		if err != nil {
			return err
		}
		if cleanup != nil {
			defer cleanup()
		}

		fmt.Printf("starting server (%s store)...\n", storeKind)
		server := api.NewPlayerServer(store, historyStore)
		log.Println("listening on :5001")
		log.Fatal(http.ListenAndServe(":5001", server))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
