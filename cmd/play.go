/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var (
// 	player1, _ = uuid.FromBytes([]byte("player1"))
// 	player2, _ = uuid.FromBytes([]byte("player2"))
// )

// simulateCmd represents the play command
var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Play prisoner's dilemna",
	Long:  `Provide a move and get a result`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		storeKind := viper.GetString("store")
		store, cleanup, err := openPlayerStore(context.Background(), storeKind)
		if err != nil {
			return err
		}
		if cleanup != nil {
			defer cleanup()
		}

		scoring := viper.GetString("scoring")

		var payoff prisoner.Payoff[int]
		switch scoring {
		case "positive", "":
			payoff = prisoner.Positive
		case "classic":
			payoff = prisoner.Classic
		default:
			return fmt.Errorf("unknown scoring %q: use classic or positive", scoring)
		}

		league, err := store.GetLeague()
		if err != nil {
			return fmt.Errorf("loading league from %s store: %w", storeKind, err)
		}
		fmt.Printf("Using %s store (%d players in league)\n", storeKind, len(league))

		roundCount, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		player1 := types.NewPlayer("")
		player2 := types.NewPlayer("")

		//------------------------/------------------------/------------------------

		historyStore, cleanup, err := openHistoryStore(context.Background(), storeKind)
		if err != nil {
			return err
		}
		if cleanup != nil {
			defer cleanup()
		}

		history, err := historyStore.GetHistory()
		if err != nil {
			return fmt.Errorf("loading history from %s store: %w", storeKind, err)
		}
		fmt.Printf("Using %s store (%d interactions in history)\n", storeKind, len(history))

		//------------------------/------------------------/------------------------
		fmt.Printf("Playing %d rounds\n", roundCount)

		// rounds := make([]types.Interaction, 0, roundCount)
		for range roundCount {
			interaction := types.NewInteraction(
				player1.ID,
				player2.ID,
				prisoner.Cooperate,
				prisoner.Cooperate)
			err = historyStore.RecordInteraction(interaction)
			if err != nil {
				return fmt.Errorf("recording interaction in %s store: %w", storeKind, err)
			}
		}

		history, err = historyStore.GetHistory()
		if err != nil {
			return fmt.Errorf("loading history from %s store: %w", storeKind, err)
		}

		for _, interaction := range history {
			fmt.Println(interaction.PrintScore(payoff))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(simulateCmd)

	simulateCmd.Flags().String("scoring", "positive", "scoring system to use, accepts, classic or positive")
	_ = viper.BindPFlag("scoring", simulateCmd.Flags().Lookup("scoring"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
