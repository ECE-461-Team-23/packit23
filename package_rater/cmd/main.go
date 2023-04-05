package main

import (
	"fmt"
	"io"
	"os"
	"net/http"

	"go.uber.org/zap"

	"github.com/packit461/packit23/package_rater/internal/cli"
	"github.com/packit461/packit23/package_rater/internal/config"
	"github.com/packit461/packit23/package_rater/internal/error"
	"github.com/packit461/packit23/package_rater/internal/logger"
	"github.com/packit461/packit23/package_rater/github_apis"
)


func return_error_packet(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 - Internal error"))
}

// Handle a request
func handle_request(w http.ResponseWriter, r *http.Request, logger *zap.Logger) {
	// Log request
	logger.Info(fmt.Sprintf("Received %s request", r.Method))
	headers := "Headers:\n"
    for key, value := range r.Header {
        headers += fmt.Sprintf("%s=%s\n", key, value)
    }
	logger.Info(headers)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Info("\nError reading body of request\n")
		return_error_packet(w, r)
		return
	}
	logger.Info(fmt.Sprintf("Body:\n%s\n", body))

	// Validate request
	if r.Method != "POST" {
		logger.Warn(fmt.Sprintf("Rejecting packet, not a POST method. It is a: %s\n", r.Method))
		return_error_packet(w, r)
		return
	}

	// TODO: Auth

	// Validate body
	url := string(body[:])
	owner, name := github_apis.ParseUrl(url)
	if owner == "" || name == "" {
		logger.Warn(fmt.Sprintf("Error parsing body (Invalid Url type): %s", url))
		return_error_packet(w, r)
		return
	}

	// Process
	score := cli.ScoreSingle(url, logger)
	if score == "" {
		logger.Warn(fmt.Sprintf("Error getting score for %s", url))
		return_error_packet(w, r)
		return
	}

	// Return results
	fmt.Fprintf(w, "%s", score) // Implicit WriteHeader(http.StatusOK)
}

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
	fmt.Println("Starting Application")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handle_request(w, r, logger)
	})

	logger.Fatal(http.ListenAndServe(":8080", nil).Error())
	os.Exit(1)
}
