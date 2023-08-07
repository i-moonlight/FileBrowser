package cmd

import (
	"encoding/json"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/filebrowser/filebrowser/v2/auth"
	"github.com/filebrowser/filebrowser/v2/settings"
)

func init() {
	configCmd.AddCommand(configImportCmd)
}

type settingsFile struct {
	Settings *settings.Settings `json:"settings"`
	Server   *settings.Server   `json:"server"`
	Auther   interface{}        `json:"auther"`
}

var configImportCmd = &cobra.Command{
	Use:   "import <path>",
	Short: "Import a configuration file",
	Long: `Import a configuration file. This will replace all the existing
configuration. Can be used with or without unexisting databases.

If used with a nonexisting database, a key will be generated
automatically. Otherwise the key will be kept the same as in the
database.

The path must be for a json or yaml file.`,
	Args: jsonYamlArg,
	Run: python(func(cmd *cobra.Command, args []string, d pythonData) {
		var key []byte
		if d.hadDB {
			settings, err := d.store.Settings.Get()
			checkErr(err)
			key = settings.Key
		} else {
			key = generateKey()
		}

		file := settingsFile{}
		err := unmarshal(args[0], &file)
		checkErr(err)

		file.Settings.Key = key
		err = d.store.Settings.Save(file.Settings)
		checkErr(err)

		err = d.store.Settings.SaveServer(file.Server)
		checkErr(err)
	}, pythonConfig{allowNoDB: true}),
}

func getAuther(sample auth.Auther, data interface{}) interface{} {
	authType := reflect.TypeOf(sample)
	auther := reflect.New(authType).Interface()
	bytes, err := json.Marshal(data)
	checkErr(err)
	err = json.Unmarshal(bytes, &auther)
	checkErr(err)
	return auther
}
