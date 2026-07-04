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

		asJSON, err := cmd.Flags().GetBool("json")
		if err != nil {
			return err
		}

		// set stores, one for players and one for history
		st, cleanup, err := openStores(context.Background(), storeKind, logger)
		if err != nil {
			return err
		}
		defer cleanup()
		historyStore := st.history
		playerStore := st.players

		roundCount, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Playing %d rounds\n", roundCount)

		// empty strings generate random names
		player1, player2, err := getTwoRandomPlayers(playerStore)
		if err != nil {
			return fmt.Errorf("getting two random players: %w", err)
		}

		for range roundCount {
			player1Move := prisoner.RandomMove()
			player2Move := prisoner.RandomMove()
			interaction := types.NewInteraction(
				player1.ID,
				player2.ID,
				player1Move,
				player2Move)
			err = historyStore.RecordInteraction(interaction)
			fmt.Println(interaction.PrintInteraction(payoff, asJSON))
			if err != nil {
				return fmt.Errorf("recording interaction in %s store: %w", storeKind, err)
			}
		}

		return nil
	},
}

func init() {
	gameCmd.AddCommand(simulateCmd)

	simulateCmd.Flags().String("scoring", "positive", "scoring system to use, accepts, classic or positive")
	_ = viper.BindPFlag("scoring", simulateCmd.Flags().Lookup("scoring"))
}

func getTwoRandomPlayers(ps types.PlayerStore) (types.Player, types.Player, error) {
	player1, err := ps.GetRandomPlayer()
	if err != nil {
		return types.Player{}, types.Player{}, fmt.Errorf("getting random player: %w", err)
	}
	player2, err := ps.GetRandomPlayerExcept(player1.ID)
	if err != nil {
		return types.Player{}, types.Player{}, fmt.Errorf("getting random player except %v: %w", player1.ID, err)
	}
	return player1, player2, nil
}
