package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     os.Args[0] + " [cmd]",
	Short:   "开启指定的组件",
	Example: os.Args[0] + " api",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("请输入指定的子命令以开启服务!")
		os.Exit(0)
	},
}

func Run() error {
	return rootCmd.Execute()
}
