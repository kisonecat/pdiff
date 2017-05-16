package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	prefixed "github.com/kisonecat/logrus-prefixed-formatter"
	"github.com/urfave/cli"
	"net/url"
	"os"
	"sort"
)

var log = logrus.New()
var repository string
var keyFingerprint string
var ximeraUrl *url.URL

func init() {
	formatter := new(prefixed.TextFormatter)
	formatter.DisableTimestamp = true
	formatter.DisableUppercase = true
	log.Formatter = formatter
}

func main() {
	app := cli.NewApp()

	app.Name = "pdiff"
	app.Usage = "a diff tool for PDFs"
	app.UsageText = "pdiff [options] a.pdf b.pdf"
	app.Version = "0.0.2"

	fmt.Printf("This is pdiff, Version " + app.Version + "\n\n")

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}

	// BADBAD: This should produce nicer error outputs
	w := log.Writer()
	defer w.Close()
	log.WriterLevel(logrus.ErrorLevel)
	cli.ErrWriter = w

	repository, _ = os.Getwd()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v, debug, d",
			Usage: "Display additional debugging information",
		},
		cli.BoolFlag{
			Name:  "no-color, C",
			Usage: "Disable color",
		},
	}

	app.Action = func(c *cli.Context) error {
		lefthand := c.Args().Get(0)
		righthand := c.Args().Get(1)

		var err error

		if len(c.Args()) != 2 {
			cli.ShowAppHelp(c)
			log.Error("pdiff expects two PDF filenames.")
			os.Exit(2)
			return nil
		}

		if _, err = os.Stat(lefthand); os.IsNotExist(err) {
			log.Error(lefthand + " does not exist.")
			os.Exit(2)
			return nil
		}

		if _, err = os.Stat(righthand); os.IsNotExist(err) {
			log.Error(righthand + " does not exist.")
			os.Exit(2)
			return nil
		}

		err = Compare(lefthand, righthand)

		if err != nil {
			log.Error(err)
			os.Exit(1)
			return err
		}

		log.Info("These two PDFs look to be the same.")
		return err
	}

	app.HideHelp = true

	app.Before = func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.Level = logrus.DebugLevel
		}

		if c.Bool("no-color") {
			color.NoColor = true
			plainLogs := new(prefixed.TextFormatter)
			plainLogs.DisableColors = true
			plainLogs.DisableTimestamp = true
			log.Formatter = plainLogs
		}

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	app.Run(os.Args)
}
