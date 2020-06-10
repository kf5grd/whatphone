package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
	whatphone "samhofi.us/x/whatphone/pkg/api"
)

const (
	// Current version string
	version = "0.1.0"

	// Exit code on failure
	exitFail = 1
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitFail)
	}
}

func run(args []string, stdout io.Writer) error {
	app := cli.App{
		Name:                   "WhatPhone",
		HelpName:               "whatphone",
		Usage:                  "Phone number lookup via EveryoneAPI",
		UseShortOptionHandling: true,
		Writer:                 stdout,
		Version:                version,

		Commands: []*cli.Command{
			{
				Name:      "lookup",
				Usage:     "Perform a phone number lookup",
				Action:    cmdLookup,
				ArgsUsage: "<phone number>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "Output JSON data",
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
	config, err := readConfig()
	if err != nil {
		return err
	}

	if config.AccountSID == "" || config.AuthToken == "" {
		return fmt.Errorf("authentication strings not set")
	}

	if c.NArg() < 1 {
		return fmt.Errorf("missing phone number")
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
