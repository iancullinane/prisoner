/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var (
// 	player1, _ = uuid.FromBytes([]byte("player1"))
// 	player2, _ = uuid.FromBytes([]byte("player2"))
// )

// playerCmd represents the play command
var playerCmd = &cobra.Command{
	Use:   "player [name|id]",
	Short: "Commands related to players",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		storeKind := viper.GetString("store")

		// set stores, one for players and one for history
		st, cleanup, err := openStores(context.Background(), storeKind, logger)
		if err != nil {
			return err
		}
		defer cleanup()

		playerStore := st.players

		// No argument: list every player.
		if len(args) == 0 {
			fmt.Printf("Using %s store, listing all players\n", storeKind)

			players, err := playerStore.GetAllPlayers()
			if err != nil {
				return fmt.Errorf("getting players from %s store: %w", storeKind, err)
			}
			fmt.Printf("Players: %+v\n", players)

			return nil
		}

		// Argument: look up a single player by ID if it parses as a
		// UUID, otherwise by name.
		arg := args[0]
		fmt.Printf("Using %s store for player %q\n", storeKind, arg)

		player, err := lookupPlayer(playerStore, arg)
		if err != nil {
			return fmt.Errorf("getting player %q from %s store: %w", arg, storeKind, err)
		}
		fmt.Printf("Player: %+v\n", player)

		return nil
	},
}

// lookupPlayer resolves a single player by ID when arg is a valid UUID,
// falling back to a name lookup otherwise.
func lookupPlayer(store types.PlayerStore, arg string) (types.Player, error) {
	if id, err := uuid.Parse(arg); err == nil {
		return store.GetPlayerByID(id)
	}
	return store.GetPlayerByName(arg)
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
