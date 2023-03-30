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
