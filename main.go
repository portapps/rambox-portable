//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/utl"
)

type config struct {
	Cleanup bool `yaml:"cleanup" mapstructure:"cleanup"`
}

var (
	app *App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Cleanup: false,
	}

	// Init app
	if app, err = NewWithCfg("rambox-portable", "Rambox", cfg); err != nil {
		Log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "Rambox.exe")
	app.Args = []string{
		"--user-data-dir=" + app.DataPath,
		"--without-update",
	}

	// Cleanup on exit
	if cfg.Cleanup {
		defer func() {
			utl.Cleanup([]string{
				path.Join(os.Getenv("APPDATA"), "Rambox"),
			})
		}()
	}

	configPath := path.Join(app.DataPath, "config.json")
	if _, err := os.Stat(configPath); err == nil {
		rawSettings, err := ioutil.ReadFile(configPath)
		if err == nil {
			jsonMapSettings := make(map[string]interface{})
			json.Unmarshal(rawSettings, &jsonMapSettings)
			Log.Info().Msgf("Current config: %s", jsonMapSettings)

			jsonMapSettings["auto_launch"] = false
			Log.Info().Msgf("New config: %s", jsonMapSettings)

			jsonSettings, err := json.Marshal(jsonMapSettings)
			if err != nil {
				Log.Error().Err(err).Msg("Update config marshal")
			}
			err = ioutil.WriteFile(configPath, jsonSettings, 0644)
			if err != nil {
				Log.Error().Err(err).Msg("Write config")
			}
		}
	}

	app.Launch(os.Args[1:])
}
