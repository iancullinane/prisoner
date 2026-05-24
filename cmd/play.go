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

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
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

		fmt.Printf("Using %s store (%d players in league)\n", storeKind, len(store.GetLeague()))

		roundCount, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		player1 := types.NewPlayer("")
		player2 := types.NewPlayer("")

		rounds := make([]types.Round, 0, roundCount)
		for i := 0; i < roundCount; i++ {
			rounds = append(rounds,
				types.NewRound(
					player1,
					player2,
					prisoner.Cooperate,
					prisoner.Cooperate))
		}

		fmt.Printf("Playing %d rounds\n", len(rounds))

		for _, round := range rounds {
			fmt.Println(round.PrintScore(payoff))
		}

		// TODO: Store the results of the game for later, only applies to file or database storage
		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

	playCmd.Flags().String("scoring", "positive", "scoring system to use, accepts, classic or positive")
	_ = viper.BindPFlag("scoring", playCmd.Flags().Lookup("scoring"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
