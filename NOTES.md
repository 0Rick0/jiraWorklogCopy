# Notes
## Date formating error
I had to revert this commit in the go-jira library to get this application to work.
https://github.com/andygrunwald/go-jira/commit/8c77107df3757c4ec5eae6e9d7c018618e708bfa

Given the documentation, this commit seems to be correct.
But the api call kept failing with this format.

This should be investigated later on to at least know why, and to propose a fix to the upstream project.

To build the project locally, replace the date format in `vendor/github.com/andygrunwald/go-jira/jira.go` `func (t Time) MarshalJSON() ([]byte, error){}`