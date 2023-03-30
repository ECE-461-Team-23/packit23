package scorer

import (
	"regexp"

	"github.com/packit461/packit23/package_rater/internal/models"
)

//TODO:
//identify all the endpoints needed:
//
//Ramp-up Time: parse or use regex on readme to get installation, quickstart, docs and examples ;; using REST API
//Correctness: number of stargazers, found in Repository ;; using GraphQL API
//Bus Factor: Total commits, and top 5 contributors ;; using REST API
//Responsiveness: pull requests in the last week to open issues in the last week ;; using REST API
//License Compatibility: 1 If license hs lesser general public, 0 otherwise (will use regex if need to search for a specific license) ;; using GraphQL API

func CalculatePackageScore(repo *models.Repository) {
	CalculateRampUpTime(repo)
	CalculateCorrectness(repo)
	CalculateBusFactor(repo)
	CalculateResponsiveness(repo)
	CalculateLicenseCompatibility(repo)
	CalculateNetScore(repo)
}

func CalculateNetScore(repo *models.Repository) {
	netScore := (repo.RampUpTimeScore * models.Weights["Ramp Up Time"]) + (repo.CorrectnessScore * models.Weights["Correctness"]) +
		(repo.BusFactorScore * models.Weights["Bus Factor"] * -1.0) + (repo.ResponsivenessScore * models.Weights["Responsiveness"]) +
		(repo.LicenseCompatibilityScore * models.Weights["License Compatibility"])
	repo.NetScore = netScore
	repo.NetPercentage = (netScore / float64(8)) * float64(100)
}

func CalculateRampUpTime(repo *models.Repository) {
	rampUpTime := 0.0
	res, _ := regexp.MatchString(`(?i)docs\b`, repo.Readme)
	if res {
		rampUpTime = rampUpTime + 0.25
	}

	res, _ = regexp.MatchString(`(?i)quick start\b`, repo.Readme)
	if res {
		rampUpTime = rampUpTime + 0.25
	}

	res, _ = regexp.MatchString(`(?i)installation\b`, repo.Readme)
	if res {
		rampUpTime = rampUpTime + 0.25
	}

	res, _ = regexp.MatchString(`(?i)example\b`, repo.Readme)
	if res {
		rampUpTime = rampUpTime + 0.25
	}

	repo.RampUpTimeScore = rampUpTime
}

func CalculateCorrectness(repo *models.Repository) {
	score := float64(repo.StarsCount) / (float64(100)) * 0.01
	if score > 1 {
		score = 1.0
	}

	repo.CorrectnessScore = score
}

func CalculateBusFactor(repo *models.Repository) {
	score := float64(repo.TopContributions) / float64(repo.Commits)
	if score > 1 || repo.Commits == 0.0 {
		score = 1.0
	}

	repo.BusFactorScore = score
}

func CalculateResponsiveness(repo *models.Repository) {
	score := float64(repo.OpenPRs) / float64(repo.OpenIssues)
	if score > 1 || repo.OpenIssues == 0.0 {
		score = 1.0
	}

	repo.ResponsivenessScore = score
}

func CalculateLicenseCompatibility(repo *models.Repository) {
	var score float64 = 0.0

	for _, l := range models.Licenses {
		re := regexp.MustCompile(`\b` + l + `\b`)
		if re.MatchString(repo.Readme) {
			score = 1.0
			break
		}
	}

	repo.LicenseCompatibilityScore = score
}
