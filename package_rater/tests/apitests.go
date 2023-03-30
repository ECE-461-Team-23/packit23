package tests

// import (
// 	"fmt"
// 	"group-24-ECE461/api"
// 	"group-24-ECE461/internal/models"
// )

// var testCount int
// var passed int
// var OWNER string
// var NAME string

// func RunAPITests() {
// 	passed = 0    // initialize passed tests
// 	testCount = 0 // initialize all tests
// 	fmt.Println("Running Tests on API...")
// 	TestCreateRESTClient()                                                                                  // run test on CreateRESTClient function
// 	TestCreasteGQLClient()                                                                                  // run test on CreateGQLClient function
// 	TestGetPullRequests()                                                                                   // run test on GetPullRequests function
// 	TestGetIssues()                                                                                         // run test on GetIssues function
// 	TestGetCommits()                                                                                        // run test on GetCommits function
// 	TestGetContributors()                                                                                   // run test on GetContributors function
// 	TestGetLicense()                                                                                        // run test on GetLicense function
// 	TestGetStars()                                                                                          // run test on GetStars function
// 	TestGetReadme()                                                                                         // run test on GetReadme function
// 	fmt.Printf("Passed %d of %d tests. Completion: %%/%d\n\n\n", passed, testCount, (passed/testCount)*100) // print test results
// }

// func TestCreateRESTClient() { // test CreateRESTClient function in the api package
// 	testCount++
// 	client, ctx := api.CreateRESTClient()

// 	if client == nil || ctx == nil {
// 		fmt.Println("Error creating REST client!")
// 	} else {
// 		fmt.Println("Success creating REST client")
// 		passed++
// 	}
// }

// func TestCreasteGQLClient() { // test CreateGQLClient function in the api package
// 	testCount++
// 	client, ctx := api.CreateGQLClient()

// 	if client == nil || ctx == nil {
// 		fmt.Println("Error creating GraphQL client!")
// 	} else {
// 		fmt.Println("Success creating GraphQL client")
// 		passed++
// 	}
// }

// func TestGetPullRequests() { // test GetPullRequests function in the api package
// 	testCount++
// 	client, ctx := api.CreateRESTClient()
// 	repo := models.NewRepository()
// 	repo.Owner = OWNER
// 	repo.Name = NAME
// 	api.GetPullRequests(client, ctx, repo)
// 	val := interface{}(repo.OpenPRs)
// 	if _, ok := val.(int); !ok {
// 		fmt.Println("Error getting pull requests")
// 	} else {
// 		fmt.Println("Success getting pull requests")
// 		passed++
// 	}
// }

// func TestGetIssues() { // test GetIssues function in the api package
// 	testCount++
// 	client, ctx := api.CreateRESTClient()
// 	repo := models.NewRepository()
// 	repo.Owner = OWNER
// 	repo.Name = NAME
// 	api.GetIssues(client, ctx, repo)
// 	val := interface{}(repo.OpenIssues)
// 	if _, ok := val.(int); !ok {
// 		fmt.Println("Error getting issues")
// 	} else {
// 		fmt.Println("Success getting issues")
// 		passed++
// 	}
// }

// func TestGetCommits() { // test GetCommits function in the api package
// 	testCount++
// 	client, ctx := api.CreateRESTClient()
// 	repo := models.NewRepository()
// 	repo.Owner = OWNER
// 	repo.Name = NAME
// 	api.GetCommits(client, ctx, repo)
// 	val := interface{}(repo.Commits)
// 	if _, ok := val.(int); !ok {
// 		fmt.Println("Error getting commits")
// 	} else {
// 		fmt.Println("Success getting commits")
// 		passed++
// 	}
// }

// func TestGetContributors() { // test GetContributors function in the api package
// 	testCount++
// 	client, ctx := api.CreateRESTClient()
// 	repo := models.NewRepository()
// 	repo.Owner = OWNER
// 	repo.Name = NAME
// 	api.GetContributors(client, ctx, repo)
// 	val := interface{}(repo.TopContributions)
// 	if _, ok := val.(int); !ok {
// 		fmt.Println("Error getting top 5 contributions")
// 	} else {
// 		fmt.Println("Success getting top 5 contributions")
// 		passed++
// 	}
// }

// func TestGetLicense() { // test GetLicense function in the api package
// 	testCount++
// 	client, ctx := api.CreateGQLClient()
// 	repo := models.NewRepository()
// 	repo.Owner = OWNER
// 	repo.Name = NAME
// 	api.GetLicense(client, ctx, repo)
// 	val := interface{}(repo.License)
// 	if _, ok := val.(string); !ok {
// 		fmt.Println("Error getting license")
// 	} else {
// 		fmt.Println("Success getting license")
// 		passed++
// 	}
// }

// func TestGetStars() { // test GetStars function in the api package
// 	testCount++
// 	client, ctx := api.CreateGQLClient()
// 	repo := models.NewRepository()
// 	repo.Owner = OWNER
// 	repo.Name = NAME
// 	api.GetStars(client, ctx, repo)
// 	val := interface{}(repo.StarsCount)
// 	if _, ok := val.(int); !ok {
// 		fmt.Println("Error getting stargazers count")
// 	} else {
// 		fmt.Println("Success getting stargazers count")
// 		passed++
// 	}
// }

// func TestGetReadme() { // test GetReadme function in the api package
// 	testCount++
// 	client, ctx := api.CreateRESTClient()
// 	repo := models.NewRepository()
// 	repo.Owner = OWNER
// 	repo.Name = NAME
// 	api.GetReadme(client, ctx, repo)
// 	val := interface{}(repo.Readme)
// 	if _, ok := val.(string); !ok {
// 		fmt.Println("Error getting readme")
// 	} else {
// 		fmt.Println("Success getting readme")
// 		passed++
// 	}
// }
