package cmd

import (
	"github.com/leilei3167/chat/internal/logic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logicCmd)
}

var logicCmd = &cobra.Command{
	Use:   "logic",
	Short: "开启logic",
	Run: func(cmd *cobra.Command, args []string) {
		logic.New().Run()
	},
}
