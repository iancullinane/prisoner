/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/iancullinane/prisoner/api"
	"github.com/iancullinane/prisoner/internal/logging"
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
		st, cleanup, err := openStores(context.Background(), storeKind)
		if err != nil {
			return err
		}
		defer cleanup()

		fmt.Printf("starting server (%s store)...\n", storeKind)
		logger := logging.NewLogger()
		logger.With(
			slog.String("service", "prisoner"),
		)
		logger.Info("starting server")

		server := api.NewPlayerServer(logger, st.players, st.history)
		logger.Info("listening on :5001")
		return http.ListenAndServe(":5001", server)
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
