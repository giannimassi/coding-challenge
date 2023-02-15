# 4Securitas - Go Coding Challenge

## Brief

Write a program in GO that open a file where a list of real-world http urls is separated by a newline (each url in a single line).
The program must take all the urls and start a goroutine for each of these, taking the http response and checking that the http response code is valid.
At least a unit test must be developed, where the http response must be mocked to ensure the requirements of the code.
The unit tests must cover at least the 99% of the code, a code coverage summary is required.

### Nice to have

- [] publish the program on github and enable a github action to start unit tests and code coverage on every commit on all branches.
- [] write the results of each http connection, with the timestamp of the http response code, into a file or a sqlite3 database
- [] use a logging framework
- [] dockerize the solution
- [x] provide a readme on how to test, build and deploy the solution

## Test

## Build and Deploy

Build the program with:

```bash
GOOS=darwin GOARCH=amd64 go build main.go -o ./challenge
```

Remember to ustomize GOOS and GOARCH according to your requirements (see go docs for more details).

The program is a command-line application that allows to check a list of urls provided as a list of lines in a file passed to the application via the environment variable.

It can be run with:

```bash
env FILE_PATH=<path to the urls file> go run main.go
```

If you have already built the program you can also use the executable:

```bash
env FILE_PATH=<path to the urls file> challenge
````

The output will be a list of urls, each followed by the result of the check (OK only if the response code is 200).

## Comments

Understanding who is the user of the application and what are the usecases is quite important in choosing the right interface to such program. In this case I made the least possible number of assumptions without overengineering this simple program, choosing to read the urls from a file passed via environment variable and print the output to stdout which is compatible with automated use (daemon) and command-line on-demand use (e.g. power user utility). An obvious improvement in case it used with automation is to make the output machine-readable.

When it comes to testing I would love to know why the brief required a coverage of 99%, since I suspect it is either higher than needed or coming from some certification requirements. In my experience 99% coverage is not always a good goal to have from an engineer perspective, since having a very high coverage does not necessarily mean the code is well tested. The threshold that in my experience makes more sense for a very high quality software development team is 90%. What I prefer to focus on is what code is tested. As a code reviewer I need to make sure that the relevant edge cases are tested, not just the code path.

In line with these thoughts I have decided to implement the unit tests that I felt were important to have.