//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"encoding/json"
	"os"
	"path"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
)

type config struct {
	Cleanup bool `yaml:"cleanup" mapstructure:"cleanup"`
}

var (
	app *portapps.App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Cleanup: false,
	}

	// Init app
	if app, err = portapps.NewWithCfg("rambox-portable", "Rambox", cfg); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
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
		rawSettings, err := os.ReadFile(configPath)
		if err == nil {
			jsonMapSettings := make(map[string]interface{})
			json.Unmarshal(rawSettings, &jsonMapSettings)
			log.Info().Msgf("Current config: %s", jsonMapSettings)

			jsonMapSettings["auto_launch"] = false
			log.Info().Msgf("New config: %s", jsonMapSettings)

			jsonSettings, err := json.Marshal(jsonMapSettings)
			if err != nil {
				log.Error().Err(err).Msg("Update config marshal")
			}
			err = os.WriteFile(configPath, jsonSettings, 0644)
			if err != nil {
				log.Error().Err(err).Msg("Write config")
			}
		}
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
