## Troubleshooting guide

Logs/config file are located in `root/.kabanero`

[Login](#Login)

[All Commands](#All commands)

--- 

### Login
`[Error] The url: <url> is not a valid kabanero url`
 
 User gets the kab url from their kabanero landing page. The login command is flexible so it takes it with or without https


`Unable to validate user: `

The microservice was not able to issue a jwt to be used for all subsequent commands due to invalid login credentials from the user. User needs to pass in their github username and github password/PAT. 

---

### All commands 
`[Error] EOF`

This is almost always a problem where the microservice has changed what it sends in the JSON response and the CLI has not been updated to parse those changes. 

`Session expired, login to your kabanero instance again`

User's jwt has expired. Login again for the microservice to issue a new jwt. If the user is just logging in with the same url, they do not have to pass it into the command again, just username and password will do since the kabanero url is stored in the config.yaml


