package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"goapp/clusters"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	// main context
	mainCtx context.Context
	// Version contains the current version.
	Version = "dev"
	// BuildDate contains a string with the build date.
	BuildDate = "unknown"
	// GitCommit git commit SHA
	GitCommit = "dirty"
	// GitBranch git branch
	GitBranch = "master"
)

func mainCmd(c *cli.Context) error {
	log.Printf("running main command with %s", c.FlagNames())
	origClusInfo:= clusters.ReadCluster(c)
	clusters.CreateClusters(c, origClusInfo)
	return nil
}

func init() {
	// handle termination signal
	mainCtx = handleSignals()
}

func handleSignals() context.Context {
	// Graceful shut-down on SIGINT/SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// create cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		sid := <-sig
		log.Printf("received signal: %d\n", sid)
		log.Println("canceling main command ...")
	}()

	return ctx
}

func main() {
	log.Print("Starting")
	dir, _ := os.Getwd()
	os.Stderr.WriteString("..........................\n")
	os.Stderr.WriteString(dir)
	os.Stderr.WriteString("\n")
	files, _ := ioutil.ReadDir("./")

	for _, f := range files {
		os.Stderr.WriteString(f.Name())
	}

     os.Stderr.WriteString("..........................\n")
	files0, _ := ioutil.ReadDir("./goapp")

	for _, f2 := range files0 {
		os.Stderr.WriteString(f2.Name())
	}

//	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS","gcp-credentials-for-docker.json")
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project",
				Usage:    "GCP project",
				Required: true, //todo use current GCP default
			},
			&cli.StringFlag{
				Name:        "location",
				Usage:       "GCP zone",

			},
		},
		Name:    "goapp",
		Usage:   "goapp CLI",
		Action:  mainCmd,
		Version: Version,
	}
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("goapp %s\n", Version)
		fmt.Printf("  Build date: %s\n", BuildDate)
		fmt.Printf("  Git commit: %s\n", GitCommit)
		fmt.Printf("  Git branch: %s\n", GitBranch)
		fmt.Printf("  Built with: %s\n", runtime.Version())
	}

	error := app.Run(os.Args)
	if error != nil {
		log.Fatal(error)
	}
}
