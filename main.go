package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/zmb3/spotify"
)

// UI constants

// DefaultSearchText The default placeholder of the search box.
const DefaultSearchText = "\x1b[38;5;244mSearch"

var artistParameter = flag.String("artist", "", "The artist to look for.")

func main() {
	flag.Parse()

	// Ready server
	launchServer()
	startAuth()

	/* FIXME: UI disabled
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	*/

	// Proof of concept
	res_ch := make(chan *spotify.FullArtist)
	if artistParameter != nil && len(*artistParameter) > 0 {
		go searchForArtist(*artistParameter, res_ch)
	}

	artist := <-res_ch
	if artist != nil {
		performAnalysisOnArtist(*artist)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("search-box", 2, maxY/2-1, maxX-2, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true

		if len(searchString) == 0 {
			fmt.Fprintln(v, DefaultSearchText)
			fmt.Print("WTF")
		} else {
			fmt.Print("Length not nil")
		}

		if _, err := g.SetCurrentView(v.Name()); err != nil {
			return err
		}
		v.Editor = gocui.EditorFunc(searchEditor)
		g.SetViewOnTop(v.Name())
	}
	return nil
}

var searchString string = ""

func searchEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	// Interpret
	switch {
	case ch != 0 && mod == 0:
		if len(searchString) == 0 {
			v.Clear()
		}
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		if len(searchString) == 0 {
			v.Clear()
		}
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		// TODO: Enter pressed
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
	//gocui.DefaultEditor.Edit(v, key, ch, mod)
	line, err := v.Line(0)
	if err != nil {
		fmt.Print("An error happened:", err)
		searchString = ""
		return
	}
	searchString = line
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
