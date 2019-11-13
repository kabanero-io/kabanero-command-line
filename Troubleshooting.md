Troubleshooting guide

Logs/config file are located in `root/.kabanero`


### All commands 
Running any commands output `[Error] EOF`

This is almost always a problem where the microservice has changed what it sends in the JSON response and the CLI has not been updated to parse those changes. 


### Login
`[Error] The url: <url> is not a valid kabanero url`
 
 User gets the kab url from <ask dave where they get it>. The login command is flexible so it takes it with or without https

`Unable to validate user: `
The microservice was not able to issue a jwt to be used for all subsequent commands due to invalid login credentials from the user. 

