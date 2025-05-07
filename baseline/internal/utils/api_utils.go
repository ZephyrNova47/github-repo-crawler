package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

var baseURL = "https://github.com"

func GetRepoURL(repo string) string {
	return baseURL + "repos/" + repo
}

func GetNumRelease(repoOwner string, repoName string) int {
	repoURL := baseURL + "/" + repoOwner + "/" + repoName

	c := colly.NewCollector()

	numRelease := 0

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("a.Link--primary.no-underline.Link", func(e *colly.HTMLElement) {
		text := e.Text
		if strings.Contains(text, "Releases") {
			// fmt.Println("Text:", text)
			re := regexp.MustCompile(`\d+`)
			match := re.FindString(text)
			numRelease, _ = strconv.Atoi(match)
			// fmt.Println("Number of releases:", numRelease)
		}
	})

	err := c.Visit(repoURL)
	if err != nil {
		fmt.Println("Error visiting URL:", err)
	}

	return numRelease
}

func GetReleaseTags(owner string, repo string, numRelease int) []string {
	log := logrus.New()
	releaseURL := baseURL + "/" + owner + "/" + repo + "/releases"

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
	})
	tags := make([]string, 0, numRelease)

	c.OnHTML("a.Link--primary.Link", func(e *colly.HTMLElement) {
		tagHref := strings.Split(e.Attr("href"), "/")
		tag := tagHref[len(tagHref)-1]
		tags = append(tags, tag)
		// fmt.Println(tag)
	})

	currentPage := 1
	for true {
		if len(tags) >= numRelease {
			break
		}
		visitURL := releaseURL + "?page=" + strconv.Itoa(currentPage)
		if err := c.Visit(visitURL); err != nil {
			log.WithError(err).Errorf("Error visiting %s: %v", visitURL, err)
			break

		}
		currentPage++
	}

	return tags
}

func GetReleaseURLs(repo string, tags []string) []string {
	releaseURLs := make([]string, len(tags))
	for i, tag := range tags {
		releaseURLs[i] = baseURL + repo + "/releases/tag/" + tag
	}
	return releaseURLs
}

func GetCommitURLs(repo string, sha string) string {
	return baseURL + "repos/" + repo + "/commits/" + sha
}

func GetNumCommitRelease(releaseURL string) int {
	log := logrus.New()
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {})

	numCommits := 0

	c.OnHTML("div.d-flex.flex-row.flex-wrap.color-fg-muted.flex-items-end", func(e *colly.HTMLElement) {

		text := e.Text
		re := regexp.MustCompile(`(\d+)\s+commits`)
		match := re.FindStringSubmatch(text)
		numCommits, _ = strconv.Atoi(match[1])
		// fmt.Println("Number of commits:", numCommits)
	})

	if err := c.Visit(releaseURL); err != nil {
		log.WithError(err).Errorf("Error visiting %s: %v", releaseURL, err)
	}
	return numCommits
}

// func main() {
// 	repo := "/opencv/opencv"
// 	fmt.Print(GetTags(repo))
// }
