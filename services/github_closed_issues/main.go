package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dinedal/message_droid/common"

	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/shurcooL/go-goon"
	"github.com/shurcooL/go/time_util"
	"golang.org/x/oauth2"
)

type worker struct {
	client  *github.Client
	orgName string

	closedIssues uint64 // Using uint64, so that there's lots of room to do great things. ;)
}

func NewWorker(token, orgName string) common.ServiceWorker {
	var transport http.RoundTripper

	// GitHub API authentication.
	transport = &oauth2.Transport{
		Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	}

	// Memory caching.
	transport = &httpcache.Transport{
		Transport:           transport,
		Cache:               httpcache.NewMemoryCache(),
		MarkCachedResponses: true,
	}

	httpClient := &http.Client{Transport: transport}

	return &worker{client: github.NewClient(httpClient), orgName: orgName}
}

func (this *worker) update() error {
	// Reset the counter to 0 and count all issues.
	this.closedIssues = 0

	// Setup a filter to get all closed issues since the beginning of the current week, in local time.
	startOfWeek := time_util.StartOfWeek(time.Now())
	opt := &github.IssueListOptions{Filter: "all", State: "closed", Since: startOfWeek}

	for {
		issues, resp, err := this.client.Issues.ListByOrg(this.orgName, opt)
		if err != nil {
			goon.Dump(err)
			log.Println("github.ListByOrg:", err)
			return err
		}

		// Despite the Since filter, issues with other activity will show up.
		// So we check the ClosedAt time and only count those closed after target time.
		for _, issue := range issues {
			if issue.ClosedAt.After(startOfWeek) {
				this.closedIssues++
			}
		}

		// Iterate over all paginated results.
		if resp.NextPage != 0 {
			opt.ListOptions.Page = resp.NextPage
			continue
		}

		break
	}

	return nil
}

func (this *worker) GetServiceUpdate() string {
	if err := this.update(); err != nil {
		return fmt.Sprintf("Error happened while calculating closed issues.")
	}

	return fmt.Sprintf("Issues Closed this week: %v", this.closedIssues)
}

func main() {
	orgNameFlag := flag.String("org-name", "", "Name of GitHub organization to get closed issues for (required).")
	flag.Parse()

	// Check for required flag value.
	if *orgNameFlag == "" {
		flag.Usage()
		os.Exit(2)
	}

	fmt.Println("Enter a GitHub token:")
	token, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Panicln("ioutil.ReadAll:", err)
	}

	fmt.Println("\nStarting.")

	common.ServiceMainLoop(NewWorker(string(token), *orgNameFlag), "github_closed_issues", 30*time.Second)
}
