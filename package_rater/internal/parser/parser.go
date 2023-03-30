package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/packit461/packit23/package_rater/internal/cli"

	"go.uber.org/zap"
)

func ParseArguments(argsWithOutProg []string, logger *zap.Logger) (e error) {
	arg := argsWithOutProg[0]
	if strings.Contains(arg, "/") {
		file, err := os.Open(arg)
		if err != nil {
			logger.Info("Error opening URL file:")
			fmt.Println("Error opening URL file:", err)
			return err
		}
		defer file.Close()
		var lines []string
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := strings.TrimRight(scanner.Text(), "\n")
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			logger.Info("Error reading URL file:")
			fmt.Println("Error reading URL file:", err)
			return err
		}
		cli.Score(lines, logger)
		return nil
	}

	if argsWithOutProg[0] == "install" {
		cli.Install(logger)
	}

	return nil
}
