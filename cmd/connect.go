package cmd

import (
	"github.com/leilei3167/chat/internal/connect"
	"github.com/spf13/cobra"
	"log"
)

var mode = "wb"

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVarP(&mode, "mode", "m", mode, "指定连接模式为wb(websocket)或者tcp")
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "开启connect层服务",

	Run: func(cmd *cobra.Command, args []string) {
		switch mode {
		case "wb":
			connect.New().Run()
		case "tcp":
			fallthrough
		default:
			log.Fatal("unsupported!")
		}
	},
}
