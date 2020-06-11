package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	whatphone "samhofi.us/x/whatphone/pkg/api"
)

const (
	// Current version string
	version = "0.2.0"

	// Exit code on failure
	exitFail = 1
)

type configFunc func() (*whatphone.API, error)

type configReader struct {
	reader configFunc
}

func newConfigReader(f configFunc) configReader {
	return configReader{
		reader: f,
	}
}

func main() {
	if err := run(os.Args, os.Stdout, newConfigReader(readConfig)); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitFail)
	}
}

func run(args []string, stdout io.Writer, cr configReader) error {
	app := cli.App{
		Name:                   "WhatPhone",
		HelpName:               "whatphone",
		Usage:                  "Phone number lookup via EveryoneAPI",
		UseShortOptionHandling: true,
		Writer:                 stdout,
		Version:                version,
		Metadata:               map[string]interface{}{"configReader": cr},

		Commands: []*cli.Command{
			{
				Name:      "lookup",
				Usage:     "Perform a phone number lookup",
				Action:    cmdLookup,
				ArgsUsage: "<phone number>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						// TODO: implement this
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "Output JSON data",
						Hidden:  true,
					},
					&cli.BoolFlag{
						Name:    "pricing-breakdown",
						Aliases: []string{"b"},
						Usage:   "Include pricing breakdown of request",
					},
					&cli.BoolFlag{
						Name:  "all",
						Usage: "Request all data points",
					},
					&cli.BoolFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Request name data",
					},
					&cli.BoolFlag{
						Name:    "profile",
						Aliases: []string{"p"},
						Usage:   "Request profile data",
					},
					&cli.BoolFlag{
						Name:    "cnam",
						Aliases: []string{"i"},
						Usage:   "Request CNAM data",
					},
					&cli.BoolFlag{
						Name:    "gender",
						Aliases: []string{"g"},
						Usage:   "Request gender data",
					},
					&cli.BoolFlag{
						Name:    "image",
						Aliases: []string{"m"},
						Usage:   "Request image data",
					},
					&cli.BoolFlag{
						Name:    "address",
						Aliases: []string{"a"},
						Usage:   "Request address data",
					},
					&cli.BoolFlag{
						Name:    "location",
						Aliases: []string{"l"},
						Usage:   "Request location data",
					},
					&cli.BoolFlag{
						Name:    "line-provider",
						Aliases: []string{"r"},
						Usage:   "Request line provider data",
					},
					&cli.BoolFlag{
						Name:    "carrier",
						Aliases: []string{"c"},
						Usage:   "Request carrier data",
					},
					&cli.BoolFlag{
						Name:    "original-carrier",
						Aliases: []string{"o"},
						Usage:   "Request original carrier data",
					},
					&cli.BoolFlag{
						Name:    "linetype",
						Aliases: []string{"t"},
						Usage:   "Request linetype data",
					},
				},
			},
			{
				Name:   "init",
				Usage:  "Initialize the app with your EveryoneAPI credentials",
				Action: cmdInit,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "accountsid",
						Aliases:  []string{"s"},
						Usage:    "EveryoneAPI Account SID",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "authtoken",
						Aliases:  []string{"t"},
						Usage:    "EveryoneAPI Auth Token",
						Required: true,
					},
				},
			},
		},
	}

	err := app.Run(args)
	if err != nil {
		return err
	}

	return nil
}

func cmdInit(c *cli.Context) error {
	var configFile string
	var err error

	if configFile, err = getConfigFile(); err != nil {
		return err
	}

	f, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	api := whatphone.New(c.String("accountsid"), c.String("authtoken"))
	err = json.NewEncoder(f).Encode(api)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.App.Writer, "Config successfully written to %s\n", configFile)
	return nil
}

func cmdLookup(c *cli.Context) error {
	cr := c.App.Metadata["configReader"].(configReader)
	reader := cr.reader
	config, err := reader()

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("unable to read config; you may need to run the init command")
		}
		return err
	}

	if config.AccountSID == "" || config.AuthToken == "" {
		return fmt.Errorf("authentication strings not set")
	}

	if c.NArg() < 1 {
		return fmt.Errorf("missing phone number")
	}

	opts := make([]whatphone.Option, 0)
	if c.Bool("name") {
		opts = append(opts, whatphone.WithName())
	}
	if c.Bool("profile") {
		opts = append(opts, whatphone.WithProfile())
	}
	if c.Bool("cnam") {
		opts = append(opts, whatphone.WithCNAM())
	}
	if c.Bool("gender") {
		opts = append(opts, whatphone.WithGender())
	}
	if c.Bool("image") {
		opts = append(opts, whatphone.WithImage())
	}
	if c.Bool("address") {
		opts = append(opts, whatphone.WithAddress())
	}
	if c.Bool("location") {
		opts = append(opts, whatphone.WithLocation())
	}
	if c.Bool("line-provider") {
		opts = append(opts, whatphone.WithLineProvider())
	}
	if c.Bool("carrier") {
		opts = append(opts, whatphone.WithCarrier())
	}
	if c.Bool("original-carrier") {
		opts = append(opts, whatphone.WithOriginalCarrier())
	}
	if c.Bool("linetype") {
		opts = append(opts, whatphone.WithLineType())
	}

	if len(opts) == 0 && !c.Bool("all") {
		return fmt.Errorf("no data points selected; use --all to request all data points")
	}

	phonenumber := c.Args().Get(0)
	result, err := config.Lookup(phonenumber, opts...)
	if err != nil {
		return err
	}

	if result.Data.Name != nil {
		fmt.Fprintf(c.App.Writer, "Name: %s\n", *result.Data.Name)
	}
	if result.Data.Profile != nil {
		profile := *result.Data.Profile
		fmt.Fprintf(c.App.Writer, "Profile:\n")
		fmt.Fprintf(c.App.Writer, "  Edu: %s\n  Job: %s\n  Relationship: %s\n", profile.Edu, profile.Job, profile.Relationship)
	}
	if result.Data.Cnam != nil {
		fmt.Fprintf(c.App.Writer, "CNAM: %s\n", *result.Data.Cnam)
	}
	if result.Data.Gender != nil {
		fmt.Fprintf(c.App.Writer, "Gender: %s\n", *result.Data.Gender)
	}
	if result.Data.Image != nil {
		image := *result.Data.Image
		fmt.Fprintf(c.App.Writer, "Image:\n")
		fmt.Fprintf(c.App.Writer, "  Cover: %s\n  Small: %s\n  Medium: %s\n  Large: %s\n", image.Cover, image.Small, image.Med, image.Large)
	}
	if result.Data.Address != nil {
		fmt.Fprintf(c.App.Writer, "Address: %s\n", *result.Data.Address)
	}
	if result.Data.Location != nil {
		location := *result.Data.Location
		fmt.Fprintf(c.App.Writer, "Location:\n")
		fmt.Fprintf(c.App.Writer, "  City, State, Zip: %s, %s, %s\n", location.City, location.State, location.Zip)
		fmt.Fprintf(c.App.Writer, "  Lat, Long: %s, %s\n", location.Geo.Latitude, location.Geo.Longitude)
	}
	if result.Data.LineProvider != nil {
		lineprovider := *result.Data.LineProvider
		fmt.Fprintf(c.App.Writer, "Line Provider:\n")
		fmt.Fprintf(c.App.Writer, "  ID: %s\n  Name: %s\n  MMS E-mail: %s\n  SMS E-mail: %s\n", lineprovider.ID, lineprovider.Name, lineprovider.MmsEmail, lineprovider.SmsEmail)
	}
	if result.Data.Carrier != nil {
		carrier := *result.Data.Carrier
		fmt.Fprintf(c.App.Writer, "Carrier:\n")
		fmt.Fprintf(c.App.Writer, "  ID: %s\n  Name: %s\n", carrier.ID, carrier.Name)
	}
	if result.Data.CarrierO != nil {
		carriero := *result.Data.CarrierO
		fmt.Fprintf(c.App.Writer, "Original Carrier:\n")
		fmt.Fprintf(c.App.Writer, "  ID: %s\n  Name: %s\n", carriero.ID, carriero.Name)
	}
	if result.Data.Linetype != nil {
		fmt.Fprintf(c.App.Writer, "Linetype: %s\n", *result.Data.Linetype)
	}
	if result.Note != "" {
		fmt.Fprintf(c.App.Writer, "Note: %s\n", result.Note)
	}
	fmt.Fprintf(c.App.Writer, "Price Total: %.4f\n", result.Pricing.Total)
	if c.Bool("pricing-breakdown") {
		fmt.Fprintf(c.App.Writer, "  Name: %.4f\n", result.Pricing.Breakdown.Name)
		fmt.Fprintf(c.App.Writer, "  Profile: %.4f\n", result.Pricing.Breakdown.Profile)
		fmt.Fprintf(c.App.Writer, "  CNAM: %.4f\n", result.Pricing.Breakdown.Cnam)
		fmt.Fprintf(c.App.Writer, "  Gender: %.4f\n", result.Pricing.Breakdown.Gender)
		fmt.Fprintf(c.App.Writer, "  Image: %.4f\n", result.Pricing.Breakdown.Image)
		fmt.Fprintf(c.App.Writer, "  Address: %.4f\n", result.Pricing.Breakdown.Address)
		fmt.Fprintf(c.App.Writer, "  Location: %.4f\n", result.Pricing.Breakdown.Location)
		fmt.Fprintf(c.App.Writer, "  Line Provider: %.4f\n", result.Pricing.Breakdown.LineProvider)
		fmt.Fprintf(c.App.Writer, "  Carrier: %.4f\n", result.Pricing.Breakdown.Carrier)
		fmt.Fprintf(c.App.Writer, "  Original Carrier: %.4f\n", result.Pricing.Breakdown.Carrier0)
		fmt.Fprintf(c.App.Writer, "  Linetype: %.4f\n", result.Pricing.Breakdown.Linetype)
	}

	if len(result.Missed) > 0 {
		fmt.Fprintf(c.App.Writer, "\nMissed: %s\n", strings.Join(result.Missed, ", "))
	}

	return nil
}

// getconfigfile determines the appropriate path to read and write the config file
func getConfigFile() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := configDir + "/whatphone"
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err = os.Mkdir(appDir, 0744)
		if err != nil {
			return "", err
		}
	}

	return appDir + "/config.json", nil
}

// loadConfig loads a config from a reader and returns an api object
func loadConfig(r io.Reader) (*whatphone.API, error) {
	var config whatphone.API
	var err error

	err = json.NewDecoder(r).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// readConfig gets the config location, opens it, and returns an api object
func readConfig() (*whatphone.API, error) {
	configFile, err := getConfigFile()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return loadConfig(f)
}
