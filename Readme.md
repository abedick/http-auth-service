Note: This is more of a demonstration of a simple Go microservice rather than a production authentication server.

# Simple HTTP Authentication Service
To Run locally, use the `startLocal.sh` script, otherwise scripts are provided in the `/scripts` directory to build the docker image for the service.

## Setup
Upon first run, the server will be in the not-running state as it will not have been configured. To configure it, use curl or something like [Postman](https://www.getpostman.com/) and set these parameters by POSTing them as JSON to `127.0.0.1:3895/update` for local or `127.0.0.1:3000/update` for Docker service.


```json
{
	"server_name":"name",
	"debug": false,
	"require_auth": false,
	"auth_key": "onlyUsedIfrequireAuthisTrue",
	"token_expiration":"6h",
	"access_modes":[""]
}
```

The server should respond with a JSON object that has "status" set as "success" if the config has been updated. The current config will follow, again in JSON format.

**Note: If require_auth is set to true, on each further HTTP request, you must attach a field called "auth" with the value you pass as "auth_key" or your requests will be denied.**

## Status
To see if the server is running, not configured, or degraded use Postman and send a get request to `/status`. You can also POST to `/status` and get the config.

## Issue
To issue a [JWT token](https://jwt.io) POST the following to `/issue`

```json
{
	"access_mode": "superuser",
	"claims": {
		"id": "string",
		"iss": "string",
		"custom_claims": {
			"username": "string"		
		}
	}
}
```

The server will respond with a JSON object with "status" as "success" and a token in the "token" field. Use [JWT.io](https://jwt.io) to look at this token and the values it holds.

## Validation
To validate a JWT Token, send a POST request to `/validate`

```json
{ "token": "token_string" }
```

The response will be a JSON object with the claims as set above or otherwise an error message.
