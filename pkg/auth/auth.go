package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/strava/go.strava"
)

var s *http.Server

// Start runs an http server and captures an OAuth response
func Start(port string) error {

	// Application id and secret can be found at https://www.strava.com/settings/api
	// define a strava.OAuthAuthenticator to hold state.
	// The callback url is used to generate an AuthorizationURL.
	authenticator := &strava.OAuthAuthenticator{
		CallbackURL: fmt.Sprintf("http://localhost:%s/exchange_token", port),
	}

	path, err := authenticator.CallbackPath()
	if err != nil {
		return err
	}

	m := http.NewServeMux()
	s = &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: m}

	m.HandleFunc(path, authenticator.HandlerFunc(oAuthSuccess, oAuthFailure))

	fmt.Printf("-------------------------------\n")
	fmt.Printf(
		"Open this URL in a browser window:\n%s\n",
		authenticator.AuthorizationURL(
			"state1", strava.Permissions.WriteViewPrivate, true))
	fmt.Printf("-------------------------------\n")

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func oAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Access Token: %s\n", auth.AccessToken)
	fmt.Printf("-------------------------------\n")

	fmt.Fprintf(w, "Complete. Close this window\n")
	go func() {
		if err := s.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
}

func oAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authorization Failure:\n")
	if err == strava.OAuthAuthorizationDeniedErr {
		fmt.Fprint(w, "The user clicked the 'Do not Authorize' button on the previous page.\n")
		fmt.Fprint(w, "This is the main error your application should handle.")
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
