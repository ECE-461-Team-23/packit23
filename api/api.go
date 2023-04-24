package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/packit461/packit23/package_rater/internal/config"
	"github.com/packit461/packit23/package_rater/internal/error"
	"github.com/packit461/packit23/package_rater/internal/helper"
	"github.com/packit461/packit23/package_rater/internal/models"

	"github.com/google/go-github/github"
	"github.com/patrickmn/go-cache"
	"github.com/shurcooL/githubv4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// TODO:
// identify all the endpoints needed:
//
// Ramp-up Time: parse or use regex on readme to get installation, quickstart, docs and url ;; using REST API
// Correctness: number of stargazers, found in Repository ;; using GraphQL API
// Bus Factor: Total commits, and top 5 contributors ;; using REST API
// Responsiveness: pull requests in the last week to open issues in the last week ;; using REST API
// License Compatibility: 1 If license, 0 otherwise (will use regex if need to search for a specific license) ;; using GraphQL API
var flag int = 0

// type DependenciesFactors struct {
// 	Data struct {
// 		Repository struct
// 	}
// }

func SendRequests(client *github.Client, graphqlClient *githubv4.Client, ctx context.Context, graphqlCtx context.Context, repo *models.Repository, logger *zap.Logger, cache *cache.Cache) (f int) {
	flag = 0

	GetDependencyQuery(graphqlClient, ctx, repo, logger, cache)
	if flag == 1 {
		return flag
	}
	GetStars(graphqlClient, ctx, repo, logger, cache)
	if flag == 1 {
		return flag
	}
	GetReadme(client, ctx, repo, logger, cache)
	if flag == 1 {
		return flag
	}
	GetPullRequests(client, ctx, repo, logger, cache)
	if flag == 1 {
		return flag
	}
	GetIssues(client, ctx, repo, logger, cache)
	if flag == 1 {
		return flag
	}
	GetContributors(client, ctx, repo, logger, cache)
	if flag == 1 {
		return flag
	}
	GetCommits(client, ctx, repo, logger, cache)
	if flag == 1 {
		return flag
	}

	return flag
}

func GetRepoOwnerFromNPM(pack string) string {
	url := "https://registry.npmjs.org/" + pack
	res, err := http.Get(url)
	if err != nil {

		fmt.Println(err)
		return ""
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		//log err
		return ""
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		//log error
		return ""
	}

	repo := data["repository"]
	value, ok := repo.(map[string]interface{})
	var repo_url string
	if !ok {
		fmt.Println("error reading values from response")
	}
	repo_url = value["url"].(string)
	parts := strings.Split(repo_url, ".com/")
	parts = strings.Split(parts[1], "/")
	return parts[0]
}

func CreateRESTClient() (*github.Client, context.Context) { // function to create github REST api client
	ctx := context.Background() // create empty context
	cfg := config.NewConfig()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.GithubToken}) // configure auth header for the client
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, ctx // returns the github rest api client and the empty context
}

func CreateGQLClient() (*githubv4.Client, context.Context) { // function to creategithub GraphQL api client
	// ctx := context.Background() // create empty context
	// cfg := config.NewConfig()
	// ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.GithubToken}) // configure auth header for the client
	// tc := oauth2.NewClient(ctx, ts)                                             // create new http client
	// graphqlClient := githubv4.NewClient(tc)                                     // create new github graphql api client from the http client template

	ctx := context.Background() // create empty context
	cfg := config.NewConfig()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.GithubToken}) // configure auth header for the client
	tc := oauth2.NewClient(ctx, ts)                                             // create new http client
	headers := map[string]string{
		"Accept": "application/vnd.github.hawkgirl-preview+json", // changed headers to enable dependency graph preview
	}
	tc = NewHTTPClientWithHeaders(tc.Transport, headers)
	graphqlClient := githubv4.NewClient(tc) // create new github graphql api client from the http client template

	return graphqlClient, ctx // returns the github graphql api client and the empty context
}

func GetPullRequests(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger, c *cache.Cache) { // function to make get request for pull requests
	cachedResp, found := c.Get(fmt.Sprintf("%s-prs", repo.Url))
	if found {
		// use the cached response
		repo.OpenPRs = cachedResp.(int)
		return
	}
	since := helper.GetLastWeek()                                                              // get date to filter results by
	s := fmt.Sprintf("org:%s repo:%s created:>%s is:pr is:open", repo.Owner, repo.Name, since) // create query string
	prs, response, err := client.Search.Issues(ctx, s, &github.SearchOptions{})                // make the request
	if err != nil {
		flag = 1
		newError := error.NewRequestError("REST", err.Error(), response.StatusCode)
		fmt.Println(newError.Error())
		logger.Info(newError.Error())
		return
	}
	c.Set(fmt.Sprintf("%s-prs", repo.Url), *prs.Total, cache.DefaultExpiration)
	logger.Debug(fmt.Sprintf("Get Pull Request: %s", response.Status))
	repo.OpenPRs = *prs.Total // populate reposigotory field

}

func GetDependencyQuery(client *githubv4.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger, c *cache.Cache) {

	variables := map[string]interface{}{ // variables to dynamically populate the graphql query structure
		"owner": githubv4.String(repo.Owner),
		"name":  githubv4.String(repo.Name),
	}
	err := client.Query(ctx, &models.Dependency, variables)

	if err != nil {
		flag = 1
		newError := error.NewRequestError("GraphQL", err.Error(), 400)
		fmt.Println(newError.Error())
		logger.Info(newError.Error())
		return
	}
	// numNodes := models.Dependency.Repository.DependencyGraphManifests.TotalCount
	numNodes := len(models.Dependency.Repository.DependencyGraphManifests.Nodes)

	edges := models.Dependency.Repository.DependencyGraphManifests.Edges
	re, _ := regexp.Compile(">|>=|<|<=|`^`")
	// logger.Info(fmt.Sprintf("Dependency count: %d", repo.DependencyCount))
	for i := 0; i < numNodes; i++ {
		dependency := edges[i].Node.Dependencies
		repo.DependencyCount += dependency.TotalCount
		logger.Info(fmt.Sprintf("Dependency count: %d", repo.DependencyCount))
		for j := 0; j < dependency.TotalCount; j++ {
			requirement := dependency.Nodes[j].Requirements
			matched := re.MatchString(requirement) // cases where major + minor version is NOT pinned
			if !matched {
				repo.PinnedVersions += 1
			}
		}
	}

	// logger.Info(fmt.Sprintf("Dependency count: %d", repo.DependencyCount))
	// logger.Info(fmt.Sprintf("Version score: %f", repo.VersionScore))
	// logger.Info(fmt.Sprintf("numNodes: %d", numNodes))

	// logger.Info(fmt.Sprintf("%+v\n", models.Dependency.Repository))

}

func GetIssues(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger, c *cache.Cache) { // function to make get requests for issues
	cachedResp, found := c.Get(fmt.Sprintf("%s-issues", repo.Url))
	if found {
		// use the cached response
		repo.OpenIssues = cachedResp.(int)
		return
	}

	since := helper.GetLastWeek()                                                              // get date to filter results by
	s := fmt.Sprintf("org:%s repo:%s created:>%s is:pr is:open", repo.Owner, repo.Name, since) // create query string
	issues, response, err := client.Search.Issues(ctx, s, &github.SearchOptions{})             // make the request
	if err != nil {
		flag = 1
		newError := error.NewRequestError("REST", err.Error(), response.StatusCode)
		fmt.Println(newError.Error())
		logger.Info(newError.Error())
		return
	}
	c.Set(fmt.Sprintf("%s-issues", repo.Url), *issues.Total, cache.DefaultExpiration)
	logger.Debug(fmt.Sprintf("Get Issues: %s", response.Status))
	repo.OpenIssues = *issues.Total // populate repository field

}

func GetCommits(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger, c *cache.Cache) { // function to make get requests for commits
	cachedResp, found := c.Get(fmt.Sprintf("%s-issues", repo.Url))
	if found {
		// use the cached response
		repo.Commits = cachedResp.(int)
		return
	}

	opt := &github.CommitsListOptions{ // define options structure to indicate results per page
		ListOptions: github.ListOptions{PerPage: 30}, // here, it's 30 results per page, as default.
	}

	_, response, err := client.Repositories.ListCommits(ctx, repo.Owner, repo.Name, opt) // make the request
	if err != nil {
		flag = 1
		newError := error.NewRequestError("REST", err.Error(), response.StatusCode)
		fmt.Println(newError.Error())
		logger.Info(newError.Error())
		return
	}

	c.Set(fmt.Sprintf("%s-commits", repo.Url), response.LastPage*30, cache.DefaultExpiration)
	logger.Debug(fmt.Sprintf("Get Commits: %s", response.Status))
	repo.Commits = response.LastPage * 30 // populate repository field

}

func GetContributors(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger, c *cache.Cache) { // function to make get requests for contributors
	cachedResp, found := c.Get(fmt.Sprintf("%s-contributors", repo.Url))
	if found {
		// use the cached response
		repo.TopContributions = cachedResp.(int)
		return
	}

	contr, response, err := client.Repositories.ListContributors(ctx, repo.Owner, repo.Name, nil) // make the requests
	if err != nil {
		flag = 1
		newError := error.NewRequestError("REST", err.Error(), response.StatusCode)
		fmt.Println(newError.Error())
		logger.Info(newError.Error())
		return
	}
	topFive := helper.GetTopFiveContributions(contr)
	c.Set(fmt.Sprintf("%s-contributors", repo.Url), topFive, cache.DefaultExpiration)
	logger.Debug(fmt.Sprintf("Get Contributors: %s", response.Status))
	repo.TopContributions = topFive // populate repository field with the total contributions of top 5 contributors

}

func GetStars(client *githubv4.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger, c *cache.Cache) { // function to make get requests for stargazers
	cachedResp, found := c.Get(fmt.Sprintf("%s-stars", repo.Url))
	if found {
		// use the cached response
		repo.StarsCount = cachedResp.(int)
		return
	}

	variables := map[string]interface{}{ // variables to dynamically populate the graphql query structure
		"owner": githubv4.String(repo.Owner),
		"name":  githubv4.String(repo.Name),
	}
	err := client.Query(ctx, &models.Stars, variables) // make the graphql request
	if err != nil {
		flag = 1
		newError := error.NewRequestError("GraphQL", err.Error(), 400)
		fmt.Println(newError.Error())
		logger.Info(newError.Error())
		return
	}
	c.Set(fmt.Sprintf("%s-stars", repo.Url), models.Stars.Repository.StargazerCount, cache.DefaultExpiration)
	logger.Debug(fmt.Sprintf("Get Stars: %s", "200 OK"))
	repo.StarsCount = models.Stars.Repository.StargazerCount // populate repository field

}

func GetReadme(client *github.Client, ctx context.Context, repo *models.Repository, logger *zap.Logger, c *cache.Cache) { // function to make get requests for readme
	cachedResp, found := c.Get(fmt.Sprintf("%s-readme", repo.Url))
	if found {
		// use the cached response
		repo.Readme = cachedResp.(string)
		return
	}

	readme, response, err := client.Repositories.GetReadme(ctx, repo.Owner, repo.Name, &github.RepositoryContentGetOptions{}) // make the requests
	if err != nil {
		flag = 1
		newError := error.NewRequestError("REST", err.Error(), response.StatusCode)
		fmt.Println(newError.Error())
		logger.Info(newError.Error())
		return
	}
	rm := helper.Base64Decode(*readme.Content)
	c.Set(fmt.Sprintf("%s-readme", repo.Url), rm, cache.DefaultExpiration)
	logger.Debug(fmt.Sprintf("Get ReadMe: %s", response.Status))
	repo.Readme = rm // populate repository field

}

func ParseUrl(url string) (owner string, name string) {
	res, _ := regexp.MatchString(`(?i)github\b`, url)
	if res {
		owner, name = helper.GetOwnerAndName(url)
	}
	res, _ = regexp.MatchString(`(?i)npmjs\b`, url)
	if res {
		name = helper.GetPackageName(url)
		if name != "" {
			owner = GetRepoOwnerFromNPM(name)
		} else {
			owner = ""
		}
	}

	return owner, name
}
