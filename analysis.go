package main

import (
	"fmt"
	"github.com/go-music-theory/music-theory/key"
	"github.com/go-music-theory/music-theory/note"
	"github.com/zmb3/spotify"
)

// MinimumTracks The least amount of tracks that would be sufficient for the analysis.
const MinimumTracks = 15

// MinimumTracks The maximum amount of tracks to use for the analysis.
const MaximumTracks = 30

type artistStyle struct {
	BPM int
	Key key.Key
}

func performAnalysisOnArtist(a spotify.FullArtist) {
	t := artistsTopTracks(a)
	tIDs := IDsForTracks(t)
	if len(tIDs) < MinimumTracks {
		// Add albums from the artist
		tIDs = append(tIDs, artistsAlbumTracks(a)...)
	}
	if len(tIDs) > MaximumTracks {
		tIDs = tIDs[0:MaximumTracks]
	}
	features := audioFeaturesForTrackIDs(tIDs...)
	tempos, keys := extractRelevantFeatures(features)
	fav_k := mostFrequentKey(keys)
	fav_t := mostFrequentTempo(tempos)
	if fav_k != nil && fav_t > 0 {
		fmt.Println("Most frequent Key:", *fav_k)
		fmt.Println("Most frequent Tempo:", fav_t)
	}
}

func extractRelevantFeatures(f []*spotify.AudioFeatures) (tempos map[int]int, keys map[int]int) {
	// Initialize dicts, they hold the counts of each
	tempos = make(map[int]int)
	keys = make(map[int]int)
	for _, feature := range f {
		if feature == nil {
			continue
		}
		// Determine Tempo
		tempo := int(feature.Tempo)
		// Update counts
		tempos[tempo] += 1
		keys[spotifyKeyToInt(feature.Key, feature.Mode)] += 1
	}
	return
}

// Array Utilities

// mostFrequentKey Returns the most frequent Key
// from a dictionary of counts
func mostFrequentKey(m map[int]int) (k_max *key.Key) {
	k_max = nil
	max := 0
	for k, c := range m {
		if c > max {
			k_max = intToKey(k)
			max = c
		}
	}
	return
}

// mostFrequentTempo Returns the most frequent Tempo
// from a dictionary of counts.
func mostFrequentTempo(m map[int]int) (t_max int) {
	t_max = -1
	max := 0
	for t, c := range m {
		if c > max {
			t_max = t
			max = c
		}
	}
	return
}

// Note Utilities

// classFromInteger Converts spotify pitch integer to a class.
// From `https://en.wikipedia.org/wiki/Pitch_class#Other_ways_to_label_pitch_classes`
func classFromInteger(c int) note.Class {
	switch spotify.Key(c) {
	case spotify.C:
		return note.C
	case spotify.CSharp:
		return note.Cs
	case spotify.D:
		return note.D
	case spotify.DSharp:
		return note.Ds
	case spotify.E:
		return note.E
	case spotify.F:
		return note.F
	case spotify.FSharp:
		return note.Fs
	case spotify.G:
		return note.G
	case spotify.GSharp:
		return note.Gs
	case spotify.A:
		return note.A
	case spotify.ASharp:
		return note.As
	case spotify.B:
		return note.B
	default:
		fmt.Println("Unexpected integer for PitchClass:", c)
		return note.Nil
	}
}

func spotifyKeyToInt(key int, mode int) int {
	c_int := key * 100
	m_int := int(mode)
	return c_int + m_int
}

func intToKey(i int) *key.Key {
	// Find Class
	c_int := i / 100
	class := classFromInteger(c_int)
	// Find mode
	m_int := i - c_int*100
	var mode key.Mode
	switch spotify.Mode(m_int) {
	case spotify.Minor:
		mode = key.Minor
	case spotify.Major:
		mode = key.Major
	default:
		fmt.Println("Unexpected integer for Mode:", i, "because Mode would be", m_int)
		mode = key.Nil
	}
	return &key.Key{Root: class, AdjSymbol: note.No, Mode: mode}
}
