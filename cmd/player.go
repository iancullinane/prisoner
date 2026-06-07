/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var (
// 	player1, _ = uuid.FromBytes([]byte("player1"))
// 	player2, _ = uuid.FromBytes([]byte("player2"))
// )

// playerCmd represents the play command
var playerCmd = &cobra.Command{
	Use:   "player",
	Short: "Commands related to players",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		storeKind := viper.GetString("store")

		// set stores, one for players and one for history
		st, cleanup, err := openStores(context.Background(), storeKind, logger)
		if err != nil {
			return err
		}
		defer cleanup()

		playerStore := st.players

		playerName := args[0]
		fmt.Printf("Using %s store for player %q\n", storeKind, playerName)

		player, err := playerStore.GetPlayerByName(playerName)
		if err != nil {
			return fmt.Errorf("getting player from %s store: %w", storeKind, err)
		}
		fmt.Printf("Player: %+v\n", player)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(playerCmd)

	// playerCmd.Flags().String("scoring", "positive", "scoring system to use, accepts, classic or positive")
	// _ = viper.BindPFlag("scoring", playerCmd.Flags().Lookup("scoring"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
