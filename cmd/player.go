/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"errors"
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
		logger.Info(fmt.Sprintf("Using %s store", storeKind))

		// set stores, one for players and one for history
		st, cleanup, err := openStores(context.Background(), storeKind, logger)
		if err != nil {
			return err
		}
		defer cleanup()

		playerStore := st.players

		asJSON, err := cmd.Flags().GetBool("json")
		if err != nil {
			return err
		}

		// No argument: list every player.
		if len(args) == 0 {
			return printAllPlayers(playerStore, asJSON)
		}

		// Argument: look up a single player by ID if it parses as a
		// UUID, otherwise by name.
		arg := args[0]

		player, err := lookupPlayer(playerStore, arg)
		if errors.Is(err, types.ErrPlayerNotFound) {
			return fmt.Errorf("no player %q in %s store", arg, storeKind)
		}
		if err != nil {
			return fmt.Errorf("getting player %q from %s store: %w", arg, storeKind, err)
		}

		if asJSON {
			return printJSON(player)
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
	gameCmd.AddCommand(playerCmd)
}

// MARK:Helpers

func listPlayers(players types.Players) {
	fmt.Println("Players:")
	for _, player := range players {
		fmt.Printf("%v\t%v\n", player.ID, player.Name)
	}
}

func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}
	fmt.Println(string(b))
	return nil
}

func printAllPlayers(ps types.PlayerStore, asJSON bool) error {
	players, err := ps.GetAllPlayers()
	if err != nil {
		return fmt.Errorf("getting players from store: %w", err)
	}
	if asJSON {
		return printJSON(players)
	}
	listPlayers(players)
	return nil
}
