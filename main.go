package main

import (
	"log"
	"time"

	vlc "github.com/adrg/libvlc-go/v3"
	"github.com/briandowns/spinner"
	"github.com/common-nighthawk/go-figure"
)

func main() {

	myFigure := figure.NewColorFigure("Radio Paradise.", "", "green", true)
	myFigure.Print()
	s := spinner.New(spinner.CharSets[39], 1*time.Second)
	s.Prefix = "\nstream-uk1.radioparadise.com: "
	s.Start()

	if err := vlc.Init("--no-video", "--quiet"); err != nil {
		log.Fatal(err)
	}
	defer vlc.Release()

	// Create a new player.
	player, err := vlc.NewPlayer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		player.Stop()
		player.Release()
	}()

	media, err := player.LoadMediaFromURL("http://stream-uk1.radioparadise.com/mp3-32")
	if err != nil {
		log.Fatal(err)
	}
	defer media.Release()

	manager, err := player.EventManager()
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan struct{})
	eventCallback := func(event vlc.Event, userData interface{}) {
		s.Stop()
		close(quit)

	}

	eventID, err := manager.Attach(vlc.MediaPlayerEndReached, eventCallback, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer manager.Detach(eventID)

	// Start playing the media.
	err = player.Play()
	if err != nil {
		log.Fatal(err)
	}

	<-quit
}
