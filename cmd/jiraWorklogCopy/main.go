package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"

	"github.com/0Rick0/jiraWorklogCopy"
)

const dateFormat = "2006-01-02T15:04:05"

type Login struct {
	BaseUrl	string
	Username string
	Password string
}

type Config struct {
	From	Login
	To		Login
	Rules	[]jiraWorklogCopy.MoveRule
}

var startDateString string
var endDateString string
var dryRun bool

var startDate time.Time
var endDate time.Time

func init() {
	flag.StringVar(&startDateString, "startDate", "", "Set the start date in YYYY-MM-DD format, e.g. 2020-01-01")
	flag.StringVar(&endDateString, "endDate", "", "Set the end date in YYYY-MM-DD format, e.g. 2020-02-01")
	flag.BoolVar(&dryRun, "dry", false, "If set, do not save hours")
}

func validateFlags() {
	if len(startDateString) == 0 || len(endDateString) == 0{
		print("Both start date and end date must be set!\n")
		os.Exit(1)
	}
	err := parseDates()
	if err != nil {
		print("One of the dates is not in a valid format\n")
		panic(err)
	}
}

func parseDates() (err error) {
	startDate, err = time.Parse(dateFormat, fmt.Sprintf("%sT00:00:00", startDateString))
	if err != nil {
		return err
	}
	endDate, err = time.Parse(dateFormat, fmt.Sprintf("%sT23:59:59", endDateString))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	validateFlags()

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	context := jiraWorklogCopy.CreateContext()
	if dryRun {
		context.SetDryRun()
	}

	err = context.LoginToSourceJira(config.From.BaseUrl, config.From.Username, config.From.Password)
	if err != nil {
		panic(err)
	}
	err = context.LoginToTargetJira(config.To.BaseUrl, config.To.Username, config.To.Password)
	if err != nil {
		panic(err)
	}
	err = context.Move(
		config.Rules,
		startDate,
		endDate,
	)
	if err != nil {
		panic(err)
	}
}