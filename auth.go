package main

import (
	"fmt"
	"net/http"

	"github.com/strava/go.strava"
)

// authHandler provides an auth url and HandlerFunc to handle its redirect
func authHandler(port string) (string, string, http.HandlerFunc, error) {

	// Application id and secret can be found at https://www.strava.com/settings/api
	// define a strava.OAuthAuthenticator to hold state.
	// The callback url is used to generate an AuthorizationURL.
	authenticator := &strava.OAuthAuthenticator{
		CallbackURL: fmt.Sprintf("http://localhost:%s/exchange_token", port),
	}

	callbackPath, err := authenticator.CallbackPath()
	if err != nil {
		return "", callbackPath, nil, err
	}

	authURL := authenticator.AuthorizationURL(
		"state1", strava.Permissions.WriteViewPrivate, true)

	handler := authenticator.HandlerFunc(success, failure)

	return authURL, callbackPath, handler, nil
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
