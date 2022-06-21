package main

import (
	"github.com/leilei3167/chat/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Run(); err != nil {
		logrus.Fatal(err)
	}
}
