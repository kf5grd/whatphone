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
	version = "0.1.3"

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

	fmt.Printf("Config successfully written to %s\n", configFile)
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
		fmt.Printf("Name: %s\n", *result.Data.Name)
	}
	if result.Data.Profile != nil {
		profile := *result.Data.Profile
		fmt.Printf("Profile:\n")
		fmt.Printf("  Edu: %s\n  Job: %s\n  Relationship: %s\n", profile.Edu, profile.Job, profile.Relationship)
	}
	if result.Data.Cnam != nil {
		fmt.Printf("CNAM: %s\n", *result.Data.Cnam)
	}
	if result.Data.Gender != nil {
		fmt.Printf("Gender: %s\n", *result.Data.Gender)
	}
	if result.Data.Image != nil {
		image := *result.Data.Image
		fmt.Printf("Image:\n")
		fmt.Printf("  Cover: %s\n  Small: %s\n  Medium: %s\n  Large: %s\n", image.Cover, image.Small, image.Med, image.Large)
	}
	if result.Data.Address != nil {
		fmt.Printf("Address: %s\n", *result.Data.Address)
	}
	if result.Data.Location != nil {
		location := *result.Data.Location
		fmt.Printf("Location:\n")
		fmt.Printf("  City, State, Zip: %s, %s, %s\n", location.City, location.State, location.Zip)
		fmt.Printf("  Lat, Long: %s, %s\n", location.Geo.Latitude, location.Geo.Longitude)
	}
	if result.Data.LineProvider != nil {
		lineprovider := *result.Data.LineProvider
		fmt.Printf("Line Provider:\n")
		fmt.Printf("  ID: %s\n  Name: %s\n  MMS E-mail: %s\n  SMS E-mail: %s\n", lineprovider.ID, lineprovider.Name, lineprovider.MmsEmail, lineprovider.SmsEmail)
	}
	if result.Data.Carrier != nil {
		carrier := *result.Data.Carrier
		fmt.Printf("Carrier:\n")
		fmt.Printf("  ID: %s\n  Name: %s\n", carrier.ID, carrier.Name)
	}
	if result.Data.CarrierO != nil {
		carriero := *result.Data.CarrierO
		fmt.Printf("Original Carrier:\n")
		fmt.Printf("  ID: %s\n  Name: %s\n", carriero.ID, carriero.Name)
	}
	if result.Data.Linetype != nil {
		fmt.Printf("Linetype: %s\n", *result.Data.Linetype)
	}
	if result.Note != "" {
		fmt.Printf("Note: %s\n", result.Note)
	}
	fmt.Printf("Price Total: %.4f\n", result.Pricing.Total)
	if c.Bool("pricing-breakdown") {
		fmt.Printf("  Name: %.4f\n", result.Pricing.Breakdown.Name)
		fmt.Printf("  Profile: %.4f\n", result.Pricing.Breakdown.Profile)
		fmt.Printf("  CNAM: %.4f\n", result.Pricing.Breakdown.Cnam)
		fmt.Printf("  Gender: %.4f\n", result.Pricing.Breakdown.Gender)
		fmt.Printf("  Image: %.4f\n", result.Pricing.Breakdown.Image)
		fmt.Printf("  Address: %.4f\n", result.Pricing.Breakdown.Address)
		fmt.Printf("  Location: %.4f\n", result.Pricing.Breakdown.Location)
		fmt.Printf("  Line Provider: %.4f\n", result.Pricing.Breakdown.LineProvider)
		fmt.Printf("  Carrier: %.4f\n", result.Pricing.Breakdown.Carrier)
		fmt.Printf("  Original Carrier: %.4f\n", result.Pricing.Breakdown.Carrier0)
		fmt.Printf("  Linetype: %.4f\n", result.Pricing.Breakdown.Linetype)
	}

	if len(result.Missed) > 0 {
		fmt.Printf("\nMissed: %s\n", strings.Join(result.Missed, ", "))
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
