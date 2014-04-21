package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/triggit/MessageDroid/common"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/shurcooL/go/time_util"
	"github.com/sourcegraph/apiproxy"
	"github.com/sourcegraph/apiproxy/service/github"
)

type worker struct {
	client  *github.Client
	orgName string

	closedIssues uint64 // Using uint64, so that there's lots of room to do great things. ;)
}

func NewWorker(token, orgName string) common.ServiceWorker {
	// GitHub authentication.
	authTransport := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}

	// Memory caching.
	memoryCacheTransport := httpcache.NewMemoryCacheTransport()
	memoryCacheTransport.Transport = authTransport

	transport := &apiproxy.RevalidationTransport{
		Transport: memoryCacheTransport,
		Check: (&githubproxy.MaxAge{
			User:         time.Hour * 24,
			Repository:   time.Hour * 24,
			Repositories: time.Hour * 24,
			Activity:     time.Hour * 12,
		}).Validator(),
	}

	httpClient := &http.Client{Transport: transport}

	return &worker{client: github.NewClient(httpClient), orgName: orgName}
}

func (this *worker) update() {
	// Reset the counter to 0 and count all issues.
	this.closedIssues = 0

	// Setup a filter to get all closed issues since the beginning of the current week, in local time.
	startOfWeek := time_util.StartOfWeek(time.Now())
	opt := &github.IssueListOptions{Filter: "all", State: "closed", Since: startOfWeek}

	for {
		issues, resp, err := this.client.Issues.ListByOrg(this.orgName, opt)
		if err != nil {
			log.Panicln("github.ListByOrg:", err)
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
}

func (this *worker) GetServiceUpdate() string {
	this.update()

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
