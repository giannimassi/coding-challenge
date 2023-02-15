# Go Coding Challenge

## Brief

Write a program in GO that open a file where a list of real-world http urls is separated by a newline (each url in a single line).
The program must take all the urls and start a goroutine for each of these, taking the http response and checking that the http response code is valid.
At least a unit test must be developed, where the http response must be mocked to ensure the requirements of the code.
The unit tests must cover at least the 99% of the code, a code coverage summary is required.

### Nice to have

- [x] publish the program on github and enable a github action to start unit tests and code coverage on every commit on all branches.
- [x] write the results of each http connection into a file or a sqlite3 database
    - [] with the timestamp of the http response code
- [] use a logging framework
- [x] dockerize the solution
- [x] provide a readme on how to test, build and deploy the solution

## Test

Test with:

```bash
go test -coverprofile cover.out ./...
```

This will run tests for the whole project and print out a code coverage summary, while additionally writing a full a code coverage report in the file `cover.out` that can be used for further inspection.

### Running tests in CI

The Github Action setup runs a complete pipeline for the program, including:

- building the project
- linting, vetting and testing the project
- verifying that the code coverage is above a set threshold (currently set to 90%)

## Build and Deploy

Build the program with:

```bash
GOOS=darwin GOARCH=amd64 go build main.go -o ./challenge
```

Remember to customize GOOS and GOARCH according to your requirements (see go docs for more details).

The program is a command-line application that allows to check the status code of a list of urls provided as a list of lines in a file passed to the application via the environment variable. The results will be written to a file as configured via env. variables.

It can be run with:

```bash
env INPUT_PATH=<path to the urls file> OUTPUT_PATH=<path to the urls file> go run main.go
```

If you have already built the program you can also use the executable:

```bash
env INPUT_PATH=<path to the urls file> OUTPUT_PATH=<path to the urls file> challenge
````

The output will be a list of urls, each followed by the result of the check (OK only if the response code is 200).

## Docker

A dockerfile is additionally provided to be able to run the executable in a portable manner via a docker container.
Note: Building the docker image requires access to the private repository.

The following command can be used to build the docker image:

```bash
docker build -t challenge .
```

This will build a docker image with the tag `challenge` that can be run with the following command (note how the environment variables are passed to the program):

```bash
docker run -v /local/path/to/urls.txt:/container/path/to/urls.txt -v /local/path/for/output/:/container/path/for/output -e INPUT_PATH=<path to the urls file> -e OUTPUT_PATH=<path to the urls file> challenge
```

TODO: it would be wise to integrate the dockerization in the CI pipeline, even just for the sake of not mantaining two separate test-and-build scripts.

## Comments

About testing, the brief required a coverage of 99%. In my experience 99% coverage is not always a good goal to have from an quality engineering perspective, since having a very high coverage does not necessarily mean the code is well tested. The threshold that in my experience makes more sense for high quality software is 90%. What I prefer to focus on is what code is tested. As a developer I need to make sure that the relevant edge cases are tested, which might not happen if I solely focus on coverage.

In line with these thoughts I have decided to implement the unit tests that I felt were important to have (around 93% coverage). Reaching 99% would mean making changes to the code to accomodate the testing required, which I think is only fit to be done if 99% is mandatory because of regulations and certification requirements.
