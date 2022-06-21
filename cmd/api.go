package cmd

import (
	"github.com/leilei3167/chat/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "开启api服务,处理用户请求",
	Run: func(cmd *cobra.Command, args []string) {
		api.New().Run()
	},
}
