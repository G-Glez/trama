package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"trama/internal/migrate"
)

func migrateCmd(c *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := migrate.Run(c.DB); err != nil {
				return fmt.Errorf("migrate: %w", err)
			}
			log.Println("Migrations applied successfully")
			return nil
		},
	}

	return cmd
}
