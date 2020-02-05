package main

import (
	"fmt"
	"net/http"

	strava "github.com/strava/go.strava"
)

// AuthHandler provides an url to direct the user to as well as
// an http.HandlerFunc to handle the redirect from the remote host.
func AuthHandler(port string) (authURL string, localPath string, handler http.HandlerFunc) {
	authenticator := &strava.OAuthAuthenticator{
		CallbackURL: fmt.Sprintf("http://localhost:%s/exchange_token", port),
	}

	// the path that our server should listen on
	localPath = "/exchange_token"

	authURL = authenticator.AuthorizationURL("state1", "activity:write", true)

	handler = authenticator.HandlerFunc(success, failure)

	return
}

func success(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Access Token: %s\n", auth.AccessToken)
	fmt.Printf("-------------------------------\n")
}

func failure(err error, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authorization Failure:\n")
	if err == strava.OAuthAuthorizationDeniedErr {
		fmt.Fprint(w, "You clicked the 'Do not Authorize' button on the previous page.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		fmt.Fprint(w, "You provided an incorrect client_id or client_secret.")
	} else if err == strava.OAuthInvalidCodeErr {
		fmt.Fprint(w, "The temporary token was not recognized.")
	} else if err == strava.OAuthServerErr {
		fmt.Fprint(w, "There was a server error, please try again")
	} else {
		fmt.Fprint(w, err)
	}
}
