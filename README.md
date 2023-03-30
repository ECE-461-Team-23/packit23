# **Group-24-ECE461**

Group 24's implementation of a package rater written in Golang for ECE 461.

Group:

- Tony Ni @ ni86@purdue.edu
- Yigitkan Balci @ ybalci@purdue.edu
- Varun Dengi @ vdengi@purdue.edu

# **Components**

# CMD

- This is the directory where main.go is located.
- This is where LOG_FILE will default to if there is not set value.
- You can run the code manually in this directory using: `go run main.go args`

# API

# Internal

- The internal directory contains all the business logic that is unique to this project, it should not be used outside of this project.

## CLI

## Config

- For configuration files we used the cobra package.
- All environmental variables are accessed through a cofig struct that get initialized using the NewConfig() function.
- In main.go there is a check to see if GITHUB_TOKEN is set if it not set the program will terminate with an error.

## Error

- Our project consists of 2 implementations error structs for handling errors: Request Errors and General Errors.

- Request Error will only appear in the API when sending out HTTP request to the package host.

- General Error will appear through out the code base whenever an error is thrown.

### Request Error

- Request Error = {
  StatusCode: int,
  Message: string,
  RequestType: string
  }
- The Status Code represents the returned HTTP status code.
- The Message contains the Error message.
- The Request Type is what kind of API request type i.e. (RESTful or GraphQL).

### General Error

- General Error = {
  Function: string,
  Message: string
  }
- The Function represents when function caused the error.
- The message represents the Error message.

## Logger

- We used the Zap library to build our logger.
- Zap provides fast peformance while also giving the option for log levels.
- The logger will be initialized in main.go and a single logger object will be passed around.

### LOG_LEVEL

- The default LOG_LEVEL is set to silent, but you can change this by setting LOG_LEVEL in bashrc.
- LOG_LEVEL=1 means debug.
- LOG_LEVEL=2 means info.
- default is silent.

### LOG_FILE

- the LOG_FILE enviromental variable determines where the log file will be output to.
- Without setting a value for LOG_FILE, logger will create a "mylog.log" in the cmd/ directory.

## Models

- The models package is used to setup querys for HTTP requests.
- This also contains additional information (listed below) that will be used in the scorer package.

### licenses

- This contains a list of licenses thare are compatible with Gnu GPL.

- The list is sourced from: [Gnu GPL List](https://www.gnu.org/licenses/license-list.en.html).

### query

- This sets up the GraphQL query used in api/api.go

### repositroy

- The repositroy sets up a Respositroy struct for each package given in the input.
- This struct will then be passed to the scorer package.

## Scorer

- the scorer package caculates a score given a Repository struct using our algorithm.

# Test
