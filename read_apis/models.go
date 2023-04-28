package main

import "time"

type ModelError struct {
	Code int32 `json:"code"`

	Message string `json:"message"`
}

type AuthenticationRequest struct {
	User *User `json:"User"`

	Secret *UserAuthenticationInfo `json:"Secret"`
}

// This is a \"union\" type. - On package upload, either Content or URL should be set. - On package update, exactly one field should be set. - On download, the Content field should be set.
type PackageData struct {
	// Package contents. This is the zip file uploaded by the user. (Encoded as text using a Base64 encoding).  This will be a zipped version of an npm package's GitHub repository, minus the \".git/\" directory.\" It will, for example, include the \"package.json\" file that can be used to retrieve the project homepage.  See https://docs.npmjs.com/cli/v7/configuring-npm/package-json#homepage.
	Content string `json:"Content,omitempty"`
	// Package URL (for use in public ingest).
	URL string `json:"URL,omitempty"`
	// A JavaScript program (for use with sensitive modules).
	JSProgram string `json:"JSProgram,omitempty"`
}

// One entry of the history of this package.
type PackageHistoryEntry struct {
	User *User `json:"User"`
	// Date of activity using ISO-8601 Datetime standard in UTC format.
	Date time.Time `json:"Date"`

	PackageMetadata *PackageMetadata `json:"PackageMetadata"`

	Action string `json:"Action"`
}

// The \"Name\" and \"Version\" are used as a unique identifier pair when uploading a package.  The \"ID\" is used as an internal identifier for interacting with existing packages.
type PackageMetadata struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	// Package version
	Version string `json:"Version"`
}

type PackageQuery struct {
	Name    string `json:"Name"`
	Version string `json:"Version,omitempty"`
}

type PackagesBody struct {
	Items []PackageQuery
}

type PackageRegExBody struct {
	RegEx string `json:"RegEx"`
}

type RegExReturn struct {
	Version string `json:"Version"`
	Name    string `json:"Name"`
}

// Package rating (cf. Project 1).  If the Project 1 that you inherited does not support one or more of the original properties, denote this with the value \"-1\".
type PackageRating struct {
	BusFactor float64 `json:"BusFactor"`

	Correctness float64 `json:"Correctness"`

	RampUp float64 `json:"RampUp"`

	ResponsiveMaintainer float64 `json:"ResponsiveMaintainer"`

	LicenseScore float64 `json:"LicenseScore"`
	// The fraction of its dependencies that are pinned to at least a specific major+minor version, e.g. version 2.3.X of a package. (If there are zero dependencies, they should receive a 1.0 rating. If there are two dependencies, one pinned to this degree, then they should receive a Â½ = 0.5 rating).
	GoodPinningPractice float64 `json:"GoodPinningPractice"`
	// The fraction of project code that was introduced through pull requests with a code review.
	PullRequest float64 `json:"PullRequest"`
	// From Part 1
	NetScore float64 `json:"NetScore"`
}

type PackageModel struct {
	Metadata *PackageMetadata `json:"metadata"`

	Data *PackageData `json:"data"`
}

type Package struct {
	Metadata *PackageMetadata `json:"metadata"`
	Data     *PackageData     `json:"data"`
}

// Authentication info for a user
type UserAuthenticationInfo struct {
	// Password for a user. Per the spec, this should be a \"strong\" password.
	Password string `json:"password"`
}

type User struct {
	Name string `json:"name"`
	// Is this user an admin?
	IsAdmin bool `json:"isAdmin"`
}
