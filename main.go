package main

import (
	"github.com/ahmedYasserM/qo/cmd"
	"github.com/ahmedYasserM/qo/pkg/logger"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logger.Error(err)
	}

}
