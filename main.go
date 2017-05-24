package main

import (
	"flag"
	"github.com/zmb3/spotify"
)

// UI constants

// DefaultSearchText The default placeholder of the search box.

var artistParameter = flag.String("artist", "", "The artist to look for.")

func main() {
	flag.Parse()

	// Ready server
	launchServer()
	startAuth()

	res_ch := make(chan *spotify.FullArtist)
	if artistParameter != nil && len(*artistParameter) > 0 {
		go searchForArtist(*artistParameter, res_ch)
	}

	artist := <-res_ch
	if artist != nil {
		performAnalysisOnArtist(*artist)
	}
}
