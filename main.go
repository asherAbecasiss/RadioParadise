package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	vlc "github.com/adrg/libvlc-go/v3"
	"github.com/briandowns/spinner"
	"github.com/common-nighthawk/go-figure"
)


const version = "v0.1.0"

type numaa [][]string

func NewNumaa(n int) numaa {

	switch n {
	case 1:
		return numaa{{"1"}}
	case 2:
		return numaa{{"2"}}
	case 3:
		return numaa{{"3"}}
	case 4:
		return numaa{{"4"}}
	case 5:
		return numaa{{"5"}}
	case 6:
		return numaa{{"6"}}
	case 7:
		return numaa{{"7"}}
	case 8:
		return numaa{{"8"}}
	case 9:
		return numaa{{"9"}}
	case 0:
		return numaa{{"0"}}
	default: // colon
		return numaa{{":"}}
	}
}

func (na numaa) join() string {
	var s string
	for i := 0; i < len(na); i++ {
		r := strings.Join(na[i], "")
		s += r + "\n"
	}
	return s
}

func (na numaa) merge(na1 numaa) numaa {
	for i := 0; i < len(na); i++ {
		na[i] = append(na[i], na1[i]...)
	}
	return na
}

func Play() {

}

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
	quitClock := make(chan bool)
	go func() {
		for {
			select {
			case <-quitClock:
				return
			default:
				v := flag.Bool("version", false, "Show version")
				ns := flag.Bool("nosec", false, "Not print second on the clock")
				t := flag.String("timezone", "", "Specify timezone \nEx) Asia/Tokyo")
				flag.Parse()

				if *v {
					fmt.Printf("%s %s\n", os.Args[0], version)
					os.Exit(0)
				}

				if *t != "" {
					tz, err := time.LoadLocation(*t)
					if err != nil {
						log.Fatal(err)
					}
					time.Local = tz
				}

			
				for {
					now := time.Now()
					h1 := NewNumaa(now.Hour() / 10)
					h2 := NewNumaa(now.Hour() % 10)
					m1 := NewNumaa(now.Minute() / 10)
					m2 := NewNumaa(now.Minute() % 10)
					colon := NewNumaa(10) // colon
					h1.merge(h2).merge(colon).merge(m1).merge(m2)
					if !*ns {
						s1 := NewNumaa(now.Second() / 10)
						s2 := NewNumaa(now.Second() % 10)
						h1.merge(colon).merge(s1).merge(s2)
					}
					fmt.Printf("%s\033[%dA", h1.join(), len(h1))

					time.Sleep(time.Second)


				}
			}
		}
	}()


	quit := make(chan struct{})
	eventCallback := func(event vlc.Event, userData interface{}) {
		s.Stop()
		close(quit)
		quitClock <- true
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
