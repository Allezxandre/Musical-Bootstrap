package main

import (
	"fmt"
	"github.com/pkg/browser"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

// RedirectURI The callback URI Spotify will use
// once authorization is granted.
const RedirectURI = "http://localhost:19331/callback"

// AutocloseHTML An HTML snippet that
// simply closes its browser tab.
const AutocloseHTML = `
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

var (
	auth                   = spotify.NewAuthenticator(RedirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState)
	ch                     = make(chan *spotify.Client)
	state                  = "749772f2dbd39a78e4d84c80066a7060" // Really just a random string
	client *spotify.Client = nil
)

// launchServer Launches the callback listening server and returns.
func launchServer() {
	http.HandleFunc("/callback", completeAuth)
	go http.ListenAndServe(":19331", nil)
}

// startAuth Opens a web-browser to request authorization from user.
func startAuth() {
	url := auth.AuthURL(state)
	fmt.Println("Opening browser for log-in...")
	browser.OpenURL(url)

	// wait for auth to complete
	client = <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are now logged in as:", user.ID)
}

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
	fmt.Fprint(w, AutocloseHTML)
	ch <- &client
}

func searchForArtist(n string, r_ch chan<- *spotify.FullArtist) {
	fmt.Println("Looking for", n)
	results, err := spotify.Search(n, spotify.SearchTypeArtist)
	if err != nil {
		fmt.Println("Cannot find artist:", err)
		return
	}

	if results.Artists != nil {
		// "Return" the first artist if we have one
		if len(results.Artists.Artists) > 0 {
			r_ch <- &results.Artists.Artists[0]
			fmt.Println(n, "found.")
		} else {
			r_ch <- nil
			fmt.Println(n, "not found.")
		}
	}
}

func artistsTopTracks(a spotify.FullArtist) []spotify.FullTrack {
	tracks, err := spotify.GetArtistsTopTracks(a.ID, spotify.CountryUnitedKingdom)
	if err != nil {
		fmt.Println("An error happened when trying to fetch top tracks:", err)
		return []spotify.FullTrack{}
	}
	return tracks
}

func artistsAlbumTracks(a spotify.FullArtist) (trackIDs []spotify.ID) {
	albumsPage, err := spotify.GetArtistAlbums(a.ID) // TODO: use options to avoid fetching too much
	if err != nil {
		fmt.Println("An error happened when trying to fetch albums:", err)
		return
	}
	for _, album := range albumsPage.Albums {
		tracksPage, err := spotify.GetAlbumTracks(album.ID)
		if err != nil {
			fmt.Println("An error happened when trying to fetch tracks:", err)
			continue
		}
		for _, t := range tracksPage.Tracks {
			trackIDs = append(trackIDs, t.ID)
		}
	}
	return
}

// audioFeaturesForTracks Returns
func audioFeaturesForTracks(t []spotify.FullTrack) []*spotify.AudioFeatures {
	// TODO: sort tracks first? (By release date for instance)
	IDs := IDsForTracks(t)
	return audioFeaturesForTrackIDs(IDs...)
}

func audioFeaturesForTrackIDs(tIDs ...spotify.ID) []*spotify.AudioFeatures {
	features, err := client.GetAudioFeatures(tIDs...)
	if err != nil {
		fmt.Println("An error happened:", err)
		return []*spotify.AudioFeatures{}
	}
	return features
}

// IDsForTracks Converts and returns an array of tracks
// into an array of IDs.
func IDsForTracks(t []spotify.FullTrack) (IDs []spotify.ID) {
	for _, track := range t {
		IDs = append(IDs, track.ID)
	}
	return
}
