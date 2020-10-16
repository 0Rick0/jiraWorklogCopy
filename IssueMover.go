package jiraWorklogCopy

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andygrunwald/go-jira"
)

type JiraMoveContext struct {
	sourceClient *jira.Client
	sourceSelf   *jira.User
	targetClient *jira.Client
	targetSelf   *jira.User
	movedTotal   int
	dryRun       bool
}

type MoveRule struct {
	From	string
	To		string
}

func CreateContext() JiraMoveContext {
	return JiraMoveContext{
		movedTotal: 0,
		dryRun: false,
	}
}

func (ctx *JiraMoveContext) SetDryRun() {
	ctx.dryRun = true
}

func (ctx *JiraMoveContext) LoginToSourceJira(baseUrl string, username string, password string) error {
	jiraClient, err := loginToJira(baseUrl, username, password)
	if err != nil {
		return err
	}
	ctx.sourceClient = jiraClient
	user, _, err := ctx.sourceClient.User.GetSelf()
	if err != nil {
		panic(err)
		return err
	}
	ctx.sourceSelf = user
	return nil
}

func (ctx *JiraMoveContext) LoginToTargetJira(baseUrl string, username string, password string) error {
	jiraClient, err := loginToJira(baseUrl, username, password)
	if err != nil {
		return err
	}
	ctx.targetClient = jiraClient
	user, _, err := ctx.targetClient.User.GetSelf()
	if err != nil {
		return err
	}
	ctx.targetSelf = user
	return nil
}

func loginToJira(baseUrl string, username string, password string) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	return jira.NewClient(tp.Client(), baseUrl)
}

func (ctx *JiraMoveContext) Move(rules []MoveRule, start time.Time, end time.Time) error {
	for _, rule := range rules {
		err := ctx.moveOne(rule, start, end)
		if err != nil {
			return err
		}
	}
	fmt.Printf("moved %d hours\n", ctx.movedTotal / 60 / 60)
	return nil
}

func (ctx *JiraMoveContext) moveOne(rule MoveRule, start time.Time, end time.Time) error {
	log, _, err := ctx.sourceClient.Issue.GetWorklogs(rule.From, func(request *http.Request) error {
		values := request.URL.Query()
		values.Add("startAfter", fmt.Sprintf("%d", start.Unix()))
		request.URL.RawQuery = values.Encode()
		return nil
	})
	if err != nil {
		return err
	}
	//fmt.Printf("%+v\n%+v\n", ctx.sourceSelf, log.Worklogs[0].Author)
	for _, worklog := range log.Worklogs {
		if worklog.Author.Key != ctx.sourceSelf.Key {
			continue
		}
		if start.After(time.Time(*worklog.Started)) {
			continue
		}
		if end.Before(time.Time(*worklog.Started)) {
			continue
		}
		ctx.movedTotal += worklog.TimeSpentSeconds

		fmt.Printf("Moving from %s at %s by %s to %s amount %s\n", rule.From, time.Time(*worklog.Started).Format("2006-01-02"), worklog.Author.Name, rule.To, worklog.TimeSpent)
		if !ctx.dryRun {
			_, _, err := ctx.targetClient.Issue.AddWorklogRecord(rule.To, &jira.WorklogRecord{
				Comment:   worklog.Comment,
				Started:   worklog.Started,
				TimeSpentSeconds: worklog.TimeSpentSeconds,
			})
			if err != nil {
				return err
			}

		}
	}
	return nil
}
