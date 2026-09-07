/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/iancullinane/prisoner/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// logger is the process-wide root logger, built from the resolved config in
// PersistentPreRunE before any subcommand runs. Commands derive component
// loggers from it when constructing stores and servers.
var logger *slog.Logger

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prisoner",
	Short: "Command line playing of prionser's dilemna",
	Long: `A small laboratory for the Prisoner's Dilemma.

Two players choose to cooperate or defect. The payoff matrix rewards
temptation and punishment, but mutual cooperation beats mutual defection—
so trust is fragile and interesting.

Scores are recorded via their algebraic symbols as follows:
- Temptation: T
- Reward: R
- Punish: P
- Sucker: S

Scoring mechanisms convert these symbols based on the payoff matrix used.

Use this CLI to play rounds from the terminal and to experiment with
strategies and simulations as the tool grows.`,
	// PersistentPreRunE runs for every subcommand after config is loaded, so the
	// logger reflects the resolved --log-level / --log-format (flag, env, config,
	// then default). Subcommands read the package-level logger var.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Args/flag validation runs before this hook, so anything that gets
		// here has already passed usage checks. Silence usage so runtime
		// errors (store failures, not-found, etc.) don't dump --help too.
		cmd.SilenceUsage = true

		level := logging.ParseLevel(viper.GetString("log-level"))
		logger = logging.New(os.Stdout, level, viper.GetString("log-format"))
		logger.Debug("logger initialised",
			slog.String("log_level", level.String()),
			slog.String("log_format", viper.GetString("log-format")),
		)
		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.prisoner)")

	rootCmd.PersistentFlags().String("store", StoreMemory, "player store backend: memory, file, or postgres")
	_ = viper.BindPFlag("store", rootCmd.PersistentFlags().Lookup("store"))

	rootCmd.PersistentFlags().String("log-level", "debug", "log level: debug, info, warn, or error")
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	rootCmd.PersistentFlags().String("log-format", "text", "log output format: text or json")
	_ = viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		if os.Getenv("DEBUG") != "" {
			home = "/Users/iancullinane/dev/iancullinane/prisoner"
		}

		// Search config in home directory with name ".prisoner" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".prisoner")
	}

	viper.SetEnvPrefix("prisoner")
	viper.SetDefault("store", StorePostgres)
	viper.SetDefault("scoring", "classic")
	viper.SetDefault("log-level", "debug")
	viper.SetDefault("log-format", "text")

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
