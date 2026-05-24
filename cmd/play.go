/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/iancullinane/prisoner/pkg/prisoner"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play prisoner's dilemna",
	Long:  `Provide a move and get a result`,
	RunE: func(cmd *cobra.Command, args []string) error {
		scoring, err := cmd.Flags().GetString("scoring")
		if err != nil {
			return err
		}

		payoff := prisoner.Positive

		fmt.Printf("Using %s scoring\n", scoring)
		r1, r2 := prisoner.Play(prisoner.Cooperate, prisoner.Cooperate)
		fmt.Println(payoff.Compute(r1), payoff.Compute(r2))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

	playCmd.Flags().String("scoring", "positive", "scoring system to use, accepts, classic or positive")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
