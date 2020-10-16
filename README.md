# Jira Worklog Copy
This tool copies the worklog from one Jira instance to another given a couple of rules.

## Usage
1. Download the appropriate binary and the example config from the releases page.
    1. There is a Windows, MacOS(Darwin) and a Linux build available.
1. Place them in the same folder
1. Rename the example configuration to `config.yaml`
1. Modify the `config.yaml` so it has your credentials and jira info.
1. Setup some copying rules.
1. Execute the program (add `-dry` to not save the new worklog, only report what would happen)
```shell script
./jiraWorklogCopy -startDate 2020-10-12 -endDate 2020-10-16 -dry
```

## Copying rules
The copying rules are very simple, you define a from ticket number, and a to ticket number.
All worklog items in a given date range will be copied from the source to the target. 

## State of the project
The project is currently in a minimal viable state.
If there is enough demand for it, the project might evolve.
