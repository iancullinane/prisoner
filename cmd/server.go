/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/iancullinane/prisoner/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// playCmd represents the play command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start an http server for the prisoner's dilemma",
	Long:  `Server will start a long running process serving HTTP requests. Takes the same store flags as using the CLI app. Supports memory, file, and postgres stores.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storeKind := viper.GetString("store")

		log := logger.With(slog.String("service", "prisoner"))

		st, cleanup, err := openStores(context.Background(), storeKind, log)
		if err != nil {
			return err
		}
		defer cleanup()

		log.Info("starting server", slog.String("store", storeKind))

		server := api.NewPlayerServer(log, st.players, st.history)
		log.Info("listening on :5001")
		return http.ListenAndServe(":5001", server)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
