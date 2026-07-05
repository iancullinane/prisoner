/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// gameCmd groups the player-facing gameplay commands (player, simulate).
// It carries the shared --json flag so every subcommand inherits it instead
// of redefining it per-command.
var gameCmd = &cobra.Command{
	Use:   "game",
	Short: "Commands for playing and inspecting the prisoner's dilemma",
}

func init() {
	rootCmd.AddCommand(gameCmd)

	gameCmd.PersistentFlags().Bool("json", false, "print output as JSON")
}
