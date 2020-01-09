## Troubleshooting guide

Logs/config file are located in `root/.kabanero`

[Login](#Login)

[All Commands](#All commands)

--- 

### Login
* `[Error] The url: <url> is not a valid kabanero url`
 
 User gets the kab url from their kabanero landing page. The login command is flexible so it takes it with or without https
 

* `Unable to validate user: `

The microservice was not able to issue a jwt to be used for all subsequent commands due to invalid login credentials from the user. User needs to pass in their github username and github password/PAT. 

* `An error occurred. An error occurred during authentication for user [user]. The Github configuration is not complete: Github teams or organization have not been defined`

The user has not set up the Github authorization for the CLI. 

* `<url> is unreachable`

The user has passed an illegitimate, unreachable url. User gets the kab url from their kabanero landing page.

* `"Your session may have expired or the credentials entered may be invalid`

The user should already have been configured but the credentials passed may be incorrect. Another possibility is that their jwt expired. They can configure the lifetime of their jwt under `cliServices: sessionExpirationSeconds` more info can be found on the kabanero docs: https://kabanero.io/docs/ref/general/kabanero-cr-config.html



---

### List

Empty Lists

If the user's cli is correctly connected with the cli service, but running the list command results in an empty list, have the user run `./kabanero list` and check if the versions for the cli and service match. If they do not, direct the user to the cli release page to download the corresponding cli for the service version that they are running. 


---

### All commands 
* `[Error] EOF`

This is almost always a problem where the microservice has changed what it sends in the JSON response and the CLI has not been updated to parse those changes. User should update to a more recent version of the CLI and if still recievinf EOF errors, open up a bug report. 

* `Session expired, login to your kabanero instance again`

User's jwt has expired. Login again for the microservice to issue a new jwt. If the user is just logging in with the same url, they do not have to pass it into the command again, just username and password will do since the kabanero url is stored in the config.yaml



