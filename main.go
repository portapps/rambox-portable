//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/utl"
)

var (
	app *App
)

func init() {
	var err error

	// Init app
	if app, err = New("rambox-portable", "Rambox"); err != nil {
		Log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "Rambox.exe")

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
