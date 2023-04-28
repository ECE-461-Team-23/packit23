package models

var License struct { // graphql query structure to get License body of a repository
	Repository struct {
		LicenseInfo struct {
			Body string `json:"body"`
		}
	} `graphql:"repository(owner: $owner, name: $name)"` // passing in owner and name for the repository to look for
}

var Stars struct { // graphql query structure to get stargazers count of a repository
	Repository struct {
		StargazerCount int `json:"stars"`
	} `graphql:"repository(owner: $owner, name: $name)"` // passing in owner and name for the repository to look for
}

var Dependency struct {
	Repository struct {
		DependencyGraphManifests struct {
			TotalCount int `json:"totalCount"`
			Nodes      []struct {
				Filename string `json:"filename"`
			}
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
			Edges []struct {
				Node struct {
					Dependencies DependenciesConnection `graphql:"dependencies(first: $first, after: $after)"`
				}
			}
		}
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type DependenciesConnection struct {
	TotalCount int `json:"dependencyCount"`
	Nodes      []struct {
		PackageName     string `json:"packageName"`
		Requirements    string `json:"requirements"`
		HasDependencies bool   `json:"hasDependencies"`
		PackageManager  string `json:"packageManager"`
	}
	PageInfo struct {
		EndCursor   string `json:"endCursor"`
		HasNextPage bool   `json:"hasNextPage"`
	} `json:"pageInfo"`
}
