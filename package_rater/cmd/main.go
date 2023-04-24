package main

import (
	"fmt"
	"os"

	"github.com/packit461/packit23/package_rater/internal/config"
	"github.com/packit461/packit23/package_rater/internal/error"
	"github.com/packit461/packit23/package_rater/internal/logger"
	"github.com/packit461/packit23/package_rater/internal/parser"
)

func main() {

	cfg := config.NewConfig()
	if err := cfg.CheckToken(); err != nil {
		fmt.Println(error.NewGeneralError("cfg.CheckToken", err.Error()).Error())
		os.Exit(1)
	}

	logger, err := logger.InitLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Application")
	// fmt.Println("Starting Application")
	argsWithOutProg := os.Args[1:]
	parser.ParseArguments(argsWithOutProg, logger)
	os.Exit(0)
}
