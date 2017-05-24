package main

import (
	"flag"
	"fmt"
	"github.com/zmb3/spotify"
)

// UI constants

// DefaultSearchText The default placeholder of the search box.

// artistParameter An optional parameter that will start the analysis with no further user interaction.
var artistParameter = flag.String("artist", "", "The artist to look for.")

func main() {
	flag.Parse()

	// Ready server
	launchServer()
	startAuth()

	// Quit if there's no artist
	if artistParameter == nil || len(*artistParameter) == 0 {
		fmt.Println("Choose an artist with the `-artist` flag")
		return
	}

	res_ch := make(chan *spotify.FullArtist)

	go searchForArtist(*artistParameter, res_ch)
	artist := <-res_ch

	if artist != nil {
		performAnalysisOnArtist(*artist)
	}
}
