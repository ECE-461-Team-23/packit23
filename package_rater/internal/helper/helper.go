package helper

import (
	"encoding/base64"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/packit461/packit23/package_rater/internal/logger"
)

func Base64Decode(str string) string { // function to decode a base64 encoded string
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		logger, _ := logger.InitLogger()
		logger.Error("Error Decoding base 64")
		return ""
	}
	return string(data)
}

func GetTopFiveContributions(contr []*github.Contributor) int { // function to get the top 5 contributors of a list of contributors on a github repository and returns the sum of contributions from all 5 contributors.
	sort.Slice(contr, func(i, j int) bool {
		return *contr[i].Contributions > *contr[j].Contributions
	})
	var topFive []*github.Contributor
	if len(contr) > 5 {
		topFive = contr[:5]
	} else {
		topFive = contr
	}

	var sum int
	for _, c := range topFive {
		sum += *c.Contributions
	}
	return sum
}

func GetLastWeek() string { // function to get the day 1 week before today. Returns a string in the format YYYY-MM-DD.
	date := time.Now().Add(time.Duration(-24*7) * time.Hour)
	YYYYMMDD := "2006-01-02"
	return date.Format(YYYYMMDD)
}

func GetOwnerAndName(url string) (owner string, name string) {
	parts := strings.Split(url, ".com/")
	if len(parts) == 1 {
		return "", ""
	}
	parts = strings.Split(parts[1], "/")
	owner = parts[0]
	name = parts[1]

	return owner, name
}

func GetPackageName(url string) (name string) {
	parts := strings.Split(url, "package/")
	if len(parts) == 1 {
		return ""
	}
	name = parts[1]
	return name
}
