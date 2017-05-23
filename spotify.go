package main

import (
	"fmt"
	"github.com/pkg/browser"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

const redirectURI = "http://localhost:19331/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	ch    = make(chan *spotify.Client)
	state = "749772f2dbd39a78e4d84c80066a7060" // Really just a random string
)

func launchServer() {
	http.HandleFunc("/callback", completeAuth)
	http.ListenAndServe(":19331", nil)
}

func startAuth() {
	url := auth.AuthURL(state)
	fmt.Println("Opening browser for log-in...")
	browser.OpenURL(url)

	// wait for auth to complete
	var client *spotify.Client = <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are now logged in as:", user.ID)
}

var autoclose_html = `
<html>
<head>
<script type= "text/javascript">
function closeWin() {
    window.close();
    return true;
}
</script>

</head>
<body onload="return closeWin();">
You may now close this window...
</body>
</html>
`

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, autoclose_html)
	ch <- &client
}
