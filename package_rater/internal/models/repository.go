package models

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Repository struct {
	Name                      string  `json:"name"`
	Owner                     string  `json:"owner"`
	Url                       string  `json:"Url"`
	StarsCount                int     `json:"starsCount"`       // Used in the calculation of correctness score
	OpenIssues                int     `json:"openIssues"`       // Used in the calculation of responsiveness score
	OpenPRs                   int     `json:"openPullRequests"` // Used in the calculation of responsiveness score
	TopContributions          int     `json:"topContributions"` // Used in the calculation of bus factor score
	Commits                   int     `json:"commits"`          // Used in the calculation of bus factor score
	DependencyCount           int     `json:"dependencyCount"`  //
	PinnedVersions            int     `json:"pinnedVersions"`   //
	License                   string  `json:"license"`          // Used in the calculation of license compatibility score
	Readme                    string  `json:"readme"`           // Used in the calculation of ramp up time score
	RampUpTimeScore           float64 `json:"rampUpTimeScore"`
	CorrectnessScore          float64 `json:"correctnessScore"`
	BusFactorScore            float64 `json:"busFactorScore"`
	ResponsivenessScore       float64 `json:"responsivenessScore"`
	LicenseCompatibilityScore float64 `json:"licenseCompatibilityScore"`
	VersionScore              float64 `json:"versionScore"`
	NetScore                  float64 `json:"netScore"`
	NetPercentage             float64 `json:"netPercentage"`
}

type Package struct {
	Url                       string  `json:"URL"`
	NetPercentage             float64 `json:"NET_SCORE"`
	RampUpTimeScore           float64 `json:"RAMP_UP_SCORE"`
	CorrectnessScore          float64 `json:"CORRECTNESS_SCORE"`
	BusFactorScore            float64 `json:"BUS_FACTOR_SCORE"`
	ResponsivenessScore       float64 `json:"RESPONSIVENESS_MAINTAINER_SCORE"`
	LicenseCompatibilityScore float64 `json:"LICENSE_SCORE"`
	VersionScore              float64 `json:"VERSION_SCORE"`
}

func NewRepository() *Repository { // Initialize empty *Repository object
	return &Repository{}
}

var Weights = map[string]float64{
	"Ramp Up Time":          2.0,
	"Correctness":           2.0,
	"Bus Factor":            3.0,
	"Responsiveness":        2.0,
	"License Compatibility": 2.0,
}

func SortRepositories(repos []*Repository) []*Repository {
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].NetScore > repos[j].NetScore
	})

	return repos
}

func ShowResults(repos []*Repository) {
	var data []Package
	repos = SortRepositories(repos)

	for _, repo := range repos {
		m := Package{
			Url:                       repo.Url,
			NetPercentage:             float64(int(repo.NetPercentage)) / 100,
			RampUpTimeScore:           float64(int(repo.RampUpTimeScore*100)) / 100,
			CorrectnessScore:          float64(int(repo.CorrectnessScore*100)) / 100,
			BusFactorScore:            float64(int(repo.BusFactorScore*100)) / 100,
			ResponsivenessScore:       float64(int(repo.ResponsivenessScore*100)) / 100,
			LicenseCompatibilityScore: float64(int(repo.LicenseCompatibilityScore*100)) / 100,
			VersionScore:              float64(int(repo.VersionScore*100)) / 100,
		}

		data = append(data, m)
	}

	for _, obj := range data {
		b, err := json.Marshal(obj)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(string(b))
	}
}
