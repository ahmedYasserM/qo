package main

import (
	"os"

	"github.com/ahmedYasserM/qo/cmd"
	"github.com/ahmedYasserM/qo/pkg/logger"
	"github.com/ahmedYasserM/qo/pkg/sandbox"
)

func main() {

	if len(os.Args) == 1 && os.Args[0] == "init" {
		if err := sandbox.StartSandBox(); err != nil {
			logger.Error(err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	if err := cmd.Execute(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}

}
