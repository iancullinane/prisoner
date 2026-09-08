/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iancullinane/prisoner/api"
	"github.com/iancullinane/prisoner/internal/lambdahttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aws/aws-lambda-go/lambda"
)

// lambdaCmd runs the server in a AWS lambda
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Serve API behind AWS lambda",
	RunE: func(cmd *cobra.Command, args []string) error {
		storeKind := viper.GetString("store")
		if err := checkLambdaStore(storeKind); err != nil {
			return err
		}

		log := logger.With(
			slog.String("service", "prisoner"),
			slog.String("runtime", "lambda"),
		)

		st, cleanup, err := openStores(context.Background(), storeKind, log)
		if err != nil {
			return err
		}
		defer cleanup()

		log.Info("cold start complete", slog.String("store", storeKind))

		server := api.NewPlayerServer(log, st.players, st.history)
		lambda.StartHandlerFunc(lambdahttp.Handler(server))
		return nil
	},
}

// checkLambdaStore rejects the backends that cannot survive Lambda's execution
// model: memory holds state per instance, and file needs a writable, durable
// working directory. Only postgres shares state across concurrent instances.
func checkLambdaStore(kind string) error {
	if kind != StorePostgres {
		return fmt.Errorf("lambda requires the %q store, got %q", StorePostgres, kind)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(lambdaCmd)
}
