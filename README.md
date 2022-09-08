# Simple API Server
The goal of this assignment is to build a simple HTTP service (with the appropriate unit tests) with two endpoints:

### POST /auth

Accepts JSON input in the format:

`{"username": "<user name>", "password": "<user password>"}`

and returns JWT OAUTH 2/OIDC token with the username as a subject. The username and the password don't have to be verified, but should not accept empty strings. The JWT token should expire in one hour.

It should return appropriate error status code if the JSON payload is not valid, or the username and password are not valid (are empty)


### POST /sum

Protected with a valid JWT token, generated by the **/auth** endpoint, provided as a Bearer Authorization header.

Accepts arbitrary JSON document as payload, which can contain a variety of things: arrays **[1,2,3,4]**, objects **{"a":1, "b":2, "c":3}**, **numbers**, and **strings**. The endpoint should find all of the numbers throughout the document and add them together.

For example:

- **[1,2,3,4]** and **{"a":6,"b":4}** both have a sum of **10**.
- **[[[2]]]** and **{"a":{"b":4},"c":-2}** both have a sum of **2**.
- **{"a":[-1,1,"dark"]}** and **[-1,{"a":1, "b":"light"}]** both have a sum of **0**.
- **[]** and **{}** both have a sum of **0**.

The response should be the **SHA256 hash of the sum of all numbers in the document**. It should return the appropriate error status code if the JWT token or the JSON payload are not valid.

## Technical details

## Requirements
The only strict requirement we pose is that the app is written in **Go** using Go modules. We recommend the use of Go patterns wherever meaningful of course, and be prepared to show and explain them when presenting the solution. 

### Project Setup
Please setup the application to be build in a standard Go way or else provide a `start.sh` script.

## Bonus (optional)
A Dockerfile to build and run the service

### External Libraries
We don't mind if you use whatever external libraries you like as long as they are with a non-restrictive for commercial use license.

### Estimation
You will see an open issue "Call for estimation". Please estimate by writing a comment when you think the task will be ready before you start. We don't set any hard deadlines.

### Notes
How to run:
- Tun the command go run main.go
- The port is 8080
- The API localhost:8080/sumapi/v1/auth generates a token
- The API localhost:8080/sumapi/v1/sum takes in a bearer token with a json body and finds the sum of the numbers. Use this as a test for the json: body:{
    "data1": [1,2,3,4],
    "data2": {"a":6,"b":4},
    "data3": [[[2]]],
    "data4": {"a":{"b":4},"c":-2},
    "data5": {"a":[-1,1,"dark"]},
    "data6": [-1,{"a":1, "b":"light"}],
    "data7": [],
    "data8": {}
}

Time spent: 
- 2 days. I spent a little longer on this test to cover best practices, flexible mocking, high test coverage on logic etc and also to remind myself of how I would do things...Hopefully not too over engineered

Packages:

1. sumapi: sum API and routes
2. tokenhelper: generates and verifies tokens
3. jsonprovider: takes in unmarshalled json as a map[string]interface{}, finds all floats and then populates the float64 slice pointer
4. golib: leverages interfaces for 3rd party APIs which can be mocked out(look at mock.go). There maybe a better way to manage this by putting each library in their own package or some other way. Also not every 3rd party API needs to be mocked out, achieving 100% test coverage may not be necessary and it can add a little complexity but I have done some 3rd party API mocking as an example
5. common: API error handling and typed errors
6. constant: viper names and some default config values

Points:

- Used context value dependency injection to pass around services, check inject.go in corresponding packages
- With the use of dependency injection and leveraging of interfaces I am able to write my own mocks for my libraries and 3rd party libraries where I can potentially get 100% coverage. Most if not all paths are covered except for the error paths which may not be worth the hassle but I have tested a few error paths using my mocks. Note: I prefer to write my own mocks than to use a 3rd party library like gomock or mock gen as I can make it more flexible and also it helps to better understand the code.
- packages golib, tokenhelper and jsonprovider have mocks check mock.go in their corresponding packages
- Avoid sentinel errors, used type errors. If I spent more time I probably would use error AS/IS error matching to improve errors. Errors should also be propagated up in a format like service1: service2: token error: the error
- Prefer to return generic 500 error for some errors and log the error internally so it does not give any information away for a potential hacker
- All input should be verified, can use regular expression to prevent hacks like sql injection
- Input json body maybe should have a length so that I can specify slice capacity which improves perfomance
- Named jsonprovider as a provider to maybe in future ultilise the factory pattern, i.e. what if the user wants to use xml instead?
- Avoid inits() they are deterministic but can be error prone if not careful
- jsons uses floats for numbers when unmarshaled but other formats tend to use int, maybe can use reflection to convert to int automatically if we know it should only be int
- used route versioning i.e. using v1 at the moment and can add v2 but still have v1 remaining if a consumer is not ready to use v2
- use golangci-lint to enforce go standards


if you have any improvements or something doesn't work please let me know!
