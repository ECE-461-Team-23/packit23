package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	api "github.com/packit461/packit23/package_rater/github_apis"
	"github.com/packit461/packit23/package_rater/internal/models"
	"github.com/packit461/packit23/package_rater/internal/scorer"

	"github.com/patrickmn/go-cache"

	"go.uber.org/zap"
)

func Install(logger *zap.Logger) {
	out, err := exec.Command("go", "list", "-m", "all").Output()
	if err != nil {
		fmt.Println("Error executing 'go list -m all':", err)
		return
	}
	dependencies := strings.Split(string(out), "\n")
	fmt.Printf("Number of dependencies: %d\n", len(dependencies)-1)
}

func Score(urlLinks []string, logger *zap.Logger) {
	if len(urlLinks) == 0 {
		logger.Info("Input file is empty, No Url to read!")
		fmt.Println("Input file is empty, No Url to read!")
		return
	}
	// fmt.Println("Scoring.....")
	var repos []*models.Repository
	client, ctx := api.CreateRESTClient()
	graphqlClient, graphqlCtx := api.CreateGQLClient()
	c := cache.New(5*time.Minute, 10*time.Minute)
	//workingDir, _ := os.Getwd()
	//cacheDir := filepath.Join(workingDir, "program-cache")
	if _, err := os.Stat("cache.txt"); os.IsNotExist(err) {
		// Create input.txt if it does not exist
		file, err := os.Create("cache.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		logger.Debug("Created cache.txt")
	} else if err != nil {
		logger.Debug("Error checking for the cache file.")
	} else {
		logger.Debug("cache.txt exists")
		err := c.LoadFile("cache.txt")
		if err != nil {
			logger.Debug("Error loading cache from file!")
		}

	}
	for _, url := range urlLinks {
		owner, name := api.ParseUrl(url)
		if owner == "" || name == "" {
			fmt.Println("Error parsing Url string (Invalid Url type)")
		} else {
			repo := models.NewRepository()
			repo.Name = name
			repo.Owner = owner
			repo.Url = url
			repo.DependencyCount = 0
			repo.PinnedVersions = 0
			flag := api.SendRequests(client, graphqlClient, ctx, graphqlCtx, repo, logger, c)
			if flag != 0 {
				continue
			}

			scorer.CalculatePackageScore(repo)
			repos = append(repos, repo)

		}
	}
	c.SaveFile("cache.txt")
	models.ShowResults(repos)
}

func ScoreSingle(url string, logger *zap.Logger) (string) {
	if url == "" {
		logger.Info("No Url to read!")
		fmt.Println("No Url to read!")
		return ""
	}

	logger.Info("Scoring.....")
	// var repos []*models.Repository
	client, ctx := api.CreateRESTClient()
	graphqlClient, graphqlCtx := api.CreateGQLClient()
	c := cache.New(5*time.Minute, 10*time.Minute)
	//workingDir, _ := os.Getwd()
	//cacheDir := filepath.Join(workingDir, "program-cache")
	if _, err := os.Stat("cache.txt"); os.IsNotExist(err) {
		// Create input.txt if it does not exist
		file, err := os.Create("cache.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		logger.Debug("Created cache.txt")
	} else if err != nil {
		logger.Debug("Error checking for the cache file.")
	} else {
		logger.Debug("cache.txt exists")
		err := c.LoadFile("cache.txt")
		if err != nil {
			logger.Debug("Error loading cache from file!")
		}

	}

	owner, name := api.ParseUrl(url)
	if owner == "" || name == "" {
		fmt.Println("Error parsing Url string (Invalid Url type)")
		return ""
	}

	repo := models.NewRepository()
	repo.Name = name
	repo.Owner = owner
	repo.Url = url
	repo.DependencyCount = 0
	repo.PinnedVersions = 0
	flag := api.SendRequests(client, graphqlClient, ctx, graphqlCtx, repo, logger, c)
	if flag != 0 {
		return ""
	}

	scorer.CalculatePackageScore(repo)
	c.SaveFile("cache.txt")
	return models.ReturnResult(repo)
}

// func Test(logger *zap.Logger) {
// 	fmt.Println("Testing.....")
// 	cmd := exec.Command("go", "test", "./...", "-cover")
// 	fmt.Println(cmd)
// 	err := cmd.Run()
// 	fmt.Println(err)
// 	if err != nil {
// 		logger.Debug("Error running tests!")
// 		fmt.Println("Error running tests:", err)
// 		return
// 	}
// }
