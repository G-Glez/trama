package cli

import (
	"database/sql"

	"github.com/spf13/cobra"
)

type CLI struct {
	DB *sql.DB
}

func New(db *sql.DB) *CLI {
	return &CLI{
		DB: db,
	}
}

func (c *CLI) Start() error {
	rootCmd := &cobra.Command{
		Use:   "trama-cli",
		Short: "TRAMA CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	return rootCmd.Execute()
}
