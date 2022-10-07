[![codecov](https://codecov.io/gh/wutzi15/knocken/branch/main/graph/badge.svg?token=MYj8xbirav)](https://codecov.io/gh/wutzi15/knocken)
![CI Status](https://github.com/wutzi15/knocken/actions/workflows/main.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


# knocken
***kocken*** is a tool that perioadically checks a website for changes and notifies you if there are any.

## Motivation
We encountered some cases of websites beeing hacked and their content being changed of redirected to another website. We wanted to be able to check the website for changes and notify the user if there are any. ***knocken*** is a tool that does this for you. ***Knocken*** is designed to be integrated with [prometheus](https://prometheus.io/), but it will also save data to a database if required. The [Blackbox exporter](https://github.com/prometheus/blackbox_exporter) provided by prometheus works in tandem with ***knocken***. While ***knocken*** is responsible for checking the website for changes, the Blackbox exporter is responsible for checking if the website is up and running. ***Knocken*** can also check, if a given website contains a given string. This can be used to check if a website is still the same, but also if a website is still up and running. Wordpress websites can be checked for the number of published posts in a given time interval.


## Available Checks
All checks will report their status on the ```/metrics``` endpoint in a prometheus compatible format. The results can also be stored in a database. The following checks are available:
### Difference Check
This check perioadically gets the HTML content of a website and compares it to the previuos HTML. The comparison is done either via the [levensthein](https://en.wikipedia.org/wiki/Levenshtein_distance) distance or the [Jaro-Winkler](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance) distance. When considering the distance algorithm please keep in mind that (citation needed :rofl: , please correct me if I'm wrong):
- The Jaro-Winkler distance is faster, but less accurate
- The levensthein distance is more expensice, but more accurate
- The Jaro-Winkler distance tends to be more accurate for short strings and is more sensitive to changes .
The resulting value will be between 0..1.0 representing the percentage of changed HTML code.
### Contains Check
This check perioadically gets the HTML content of a website and checks if it contains a given string. The resulting value will be either 0 or 1.0. 0 if the string is not found and 1.0 if the string is found.
### Wordpress Check
This check perioadically gets the number of published posts of a wordpress website.
## Installation
### First steps
- create a ```targets.yml``` file with your domains that should be checked. Take a look at the provided ```targets.sample.yml``` file for an example or the configuration section of the README.
- if required create a  ```ignore.yml``` file with your domains that should not be checked. Take a look at the provided ```ignore.sample.yml``` file for an example or the configuration section of the README. This can be useful, when ***knocken*** is used with prometheus and you want to ignore some domains, as the ```targets.yml``` can be used for blackbox_exporter and ***knocken***.
- if required create a ```.env``` file with your configuration. Take a look at the provided ```env.sample``` file for an example or the configuration section of the README.
### Docker
Knocken can be run locally without docker, however it is recommended to run it in a docker container. The docker container is available on [dockerhub](https://hub.docker.com/r/wutzi/knocken). The docker container can be run with the following command:
```bash
docker run -p 9101:9101 -v $(pwd)/targets.yml:/app/targets.yml -v $(pwd)/.env:/app/.env wutzi/knocken
```
The docker container exposes port 9101. This port is used by the prometheus to scrape the metrics. The targets.yml file is used to configure the targets that should be checked. The targets.yml file is described in the next section.
Though it is entirely possible to run ***knocken*** directly from the command line, mounting all of the required files can be tedieous. Therefore it is recommended to use the docker compose. A sample docker compose file is provided in the repository. The docker compose file can be run with the following command:
```bash
docker compose up
```
for older versions of docker, the command is:
```bash
docker-compose up
```
An example docker compose file looks like this:
```yaml
version: "3"

services:
  knocken:
    image: wutzi/knocken
    ports:
      - 9101:9101
    volumes:
      - ./targets.yml:/app/targets.yml
      - ./ignore.yml:/app/ignore.yml
      - ./.env:/app/.env
    restart: always
```

## runnning natively
***knocken*** can also be run natively. To do so, you need to have go installed. The following command will install ***knocken***:
```bash
go get -d -v ./...
go install -v ./...
go build knocken.go

or

go run knocken.go
```
Altough running in docker is the prefered method of running ***knocken***, it is handy to run natively for debugging purposes.
## Configuration
There are various configuration options for ***knocken***, which will be described in the following section.
### Targets
The configuration of ***knocken*** targets is done via a yaml file. This file has the same structure as the targets file used by blackbox_exporter, so it can be used for both programs. The file is called `targets.yml` and has the following structure:
```yaml
targets:
  - google.com
  - escsoftware.de
```
There is also a ```targets.sample.yml``` provided.
### Ignore
The ignore file can be used to exclude targets from the check. This can be useful, when ***knocken*** is used with prometheus and you want to ignore some domains, as the ```targets.yml``` can be used for blackbox_exporter and ***knocken***. The file is called `ignore.yml` and has the following structure:
```yaml
targets:
  - google.com
```
There is also a ```ignore.sample.yml``` provided.

### Contains
The contains file can be used to check if a given string is contained in the HTML of a website. The file is called `containstargets.yml` and has the following structure:
```yaml
targets:
- domain: google.com
  contain: google
- domain: https://www.escsoftware.de/impressum/
  contain: Geschäftsführer
```
There is also a ```conatinsTargets.sample.yml``` provided.
### Environment variables
Global options for ***knocken*** can be set via environment variables. The following environment variables are available:
- `CONTAINSTARGETS` sets the name of the targets for the "contains" check (default: containstargets.yml)
- `FASTDIFF` use the Jaro Winkler difference if set to frue (default: false)
- `IGNORE` sets the name of the targets to ignore (default: ignore.yml)
- `RUNCONTAIN` run the contains check (default: true)
- `RUNDIFF` run the diff check (default: true)
- `SAVECONFIG` save the config file to .env (default: true)
- `SAVEDIFF` save the diff value - used mostly for debugging and testing (default: false)
- `TARGETS` sets the name of the targets to check (default: targets.yml)
- `VERBOSE` run with more output (default: false)
- `WAITTIME` time to wait between to checks. Can be in seconds eg. 7s, minutes eg. 5m or any golang compatible time string (default: 5m)

There is a end.sample file provided.

# Sponsors
![esc software](https://www.escsoftware.de/wp-content/uploads/software-entwicklung-hamburg.svg)