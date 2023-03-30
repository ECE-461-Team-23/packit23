package scorer

import (
	"fmt"
	"github.com/packit461/packit23/package_rater/internal/models"
	"testing"
)

func TestCalculateNetScore(t *testing.T) {
	repo := models.NewRepository()
	repo.RampUpTimeScore = 0.6
	repo.CorrectnessScore = 1
	repo.BusFactorScore = 1
	repo.ResponsivenessScore = 1
	repo.LicenseCompatibilityScore = 0

	CalculateNetScore(repo)
	t.Log()
	if repo.NetScore < -3 || repo.NetScore > 8 {
		t.Error("NetScore is outside determined bounds")
	}
	fmt.Println(repo.NetScore)
}

func TestCalculateRampUpTime(t *testing.T) {
	repo := models.NewRepository()
	repo.Readme = "DoCs, quickStart, InstalLation, Example"
	CalculateRampUpTime(repo)
	if repo.RampUpTimeScore > 1.0 {
		t.Error("Ramp Up Time score can't  be greater than 1")
	}

	if repo.RampUpTimeScore < 0.0 {
		t.Error("Ramp Up Time score can't be less than 0")
	}
}

func TestCalculateCorrectness(t *testing.T) {
	repo := models.NewRepository()
	repo.StarsCount = 100
	CalculateCorrectness(repo)
	if repo.CorrectnessScore > 1.0 {
		t.Error("Correctness score can't  be greater than 1")
	}

	if repo.CorrectnessScore < 0.0 {
		t.Error("Correctness score can't be less than 0")
	}
	repo.StarsCount = 30000
	CalculateCorrectness(repo)
	if repo.CorrectnessScore > 1.0 {
		t.Error("Correctness score can't  be greater than 1")
	}

	if repo.CorrectnessScore < 0.0 {
		t.Error("Correctness score can't be less than 0")
	}

}

func TestCalculateBusFactor(t *testing.T) {
	repo := models.NewRepository()
	repo.TopContributions = 3000
	repo.Commits = 10000
	CalculateBusFactor(repo)
	if repo.BusFactorScore > 1.0 {
		t.Error("Bus Factor score can't  be greater than 1")
	}

	if repo.BusFactorScore < 0.0 {
		t.Error("Bus Factor score can't be less than 0")
	}

	repo.Commits = 0
	CalculateBusFactor(repo)
	if repo.BusFactorScore > 1.0 {
		t.Error("Bus Factor score can't  be greater than 1")
	}

	if repo.BusFactorScore < 0.0 {
		t.Error("Bus Factor score can't be less than 0")
	}
}

func TestCalculateResponsiveness(t *testing.T) {
	repo := models.NewRepository()
	repo.OpenPRs = 10
	repo.OpenIssues = 40
	CalculateResponsiveness(repo)
	if repo.ResponsivenessScore > 1.0 {
		t.Error("Responsiveness score can't  be greater than 1")
	}

	if repo.ResponsivenessScore < 0.0 {
		t.Error("Responsiveness score can't be less than 0")
	}

	repo.OpenIssues = 0
	CalculateResponsiveness(repo)
	if repo.ResponsivenessScore > 1.0 {
		t.Error("Responsiveness score can't  be greater than 1")
	}

	if repo.ResponsivenessScore < 0.0 {
		t.Error("Responsiveness score can't be less than 0")
	}
}

func TestCalculateLicenseCompatibility(t *testing.T) {
	repo := models.NewRepository()
	repo.License = "jkahsgkjah wekut qbnhwoiuhalksndg oah o9iahsdg oihl  lesser general public askjdghaksj nqwerknka"
	CalculateLicenseCompatibility(repo)
	if repo.LicenseCompatibilityScore > 1.0 {
		t.Error("License Compatibility score can't  be greater than 1")
	}

	if repo.LicenseCompatibilityScore < 0.0 {
		t.Error("License Compatibility score can't be less than 0")
	}

	repo.License = ""
	CalculateLicenseCompatibility(repo)
	if repo.LicenseCompatibilityScore > 1.0 {
		t.Error("License Compatibility score can't  be greater than 1")
	}

	if repo.LicenseCompatibilityScore < 0.0 {
		t.Error("License Compatibility score can't be less than 0")
	}
}
