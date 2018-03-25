package main

import (
	"bufio"
	"canihazmidi/src/prolink2midi/midi"
	"errors"
	"fmt"
	"go.evanpurkhiser.com/prolink"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"strconv"
	"strings"
	syssync "sync"
	"time"
)

func initApp() *cli.App {
	return &cli.App{
		Name:        "prolink2midi",
		Description: "Converts prolink sync to midi",
		Commands: []*cli.Command{
			{
				Name:        "sync",
				Description: "Start sync",
				ArgsUsage:   "MIDI_OUT",
				Before: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return errors.New("missing arg MIDI_OUT")
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "midi-test",
						Usage:   "Start midi clock without network in interactive mode",
						Aliases: []string{"t"},
					},
				},
			},
		},
	}
}

func main() {
	app := initApp()
	app.Command("sync").Action = cli.ActionFunc(func(c *cli.Context) error {
		var midiOut *os.File
		var network *prolink.Network
		var err error

		if midiOut, err = openMidiOut(c.Args().First()); err != nil {
			return err
		}

		if c.IsSet("midi-test") {
			midiTest(midiOut)
			return nil
		}

		if network, err = initNetwork(); err != nil {
			return err
		}

		newSynchronizer(network.CDJStatusMonitor(), midi.NewClock(midiOut))

		time.Sleep(200 * time.Hour)
		return nil
	})

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func openMidiOut(devicePath string) (*os.File, error) {
	return os.OpenFile(devicePath, os.O_WRONLY, 0664)
}

func initNetwork() (*prolink.Network, error) {
	net, err := prolink.Connect()
	if err != nil {
		return nil, err
	}

	err = net.AutoConfigure(5 * time.Second)

	return net, err

}

func midiTest(midiOut *os.File) {
	var wg syssync.WaitGroup

	clk := midi.NewClock(midiOut)
	clk.Stop()
	fmt.Println("Available Commands: start, stop or a bpm value to be set (ex. 120.5)")

	// Endless
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer midiOut.Close()

		reader := bufio.NewReader(os.Stdin)
		var cmd string
		for {
			fmt.Print("> ")
			cmd, _ = reader.ReadString('\n')
			cmd = strings.TrimSuffix(cmd, "\n")
			switch {
			case cmd == "start":
				clk.Start()
			case cmd == "stop":
				clk.Stop()
			default:
				if bpm, err := strconv.ParseFloat(cmd, 32); err == nil {
					clk.SetBpm(float32(bpm))
				} else {
					fmt.Println("Invalid tempo value")
				}
			}
		}
	}()

	wg.Wait()
}
