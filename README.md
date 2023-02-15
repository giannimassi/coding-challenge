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

The program is a command-line application that allows to check a list of urls provided as a list of lines in a file passed to the application via the environment variable

## Comments
Understanding who is the user of the application and what are the usecases is quite important in choosing the right interface to such program. In this case I made the least possible number of assumptions without overengineering this simple program, choosing to read the urls from a file passed via environment variable and print the output to stdout which is compatible with automated use (daemon) and command-line on-demand use (e.g. power user utility). An obvious improvement in case it used with automation is to make the output machine-readable.

About the testing, I have chosen to implement a unit test for the `checkResponseCode` function and for the `main`function, providing two levels of testing, one for the most important function in the code base, the other one, to make sure the whole code path is tested. Given the simplicity of the task I didn't deem necessary to implement unit tests for all the functions individually, althought that would be preferred, providing more usefuleness that the current approach when it comes to maintaining the code.
