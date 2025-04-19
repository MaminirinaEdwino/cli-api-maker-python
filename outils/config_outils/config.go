package configoutils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Db_config struct {
	Db_name string
	Sgbd    string
	Async   bool
}

type Attributs struct {
	Attribut_name  string
	Attributs_type string
}

type Models_config struct {
	Model_name string
	Attributs  []Attributs
}


type Config_file struct {
	Api_name          string
	Api_version       string
	Api_description   string
	Api_reload        bool
	Host              string
	Port              string
	Db                Db_config
	Models            []Models_config
	Template_dirs     string
	Static_files_dirs string
}

func (c *Config_file) ConfigGenerator() error {
	fmt.Println("generate")
	c.Api_name = "My api"
	c.Api_description = c.Api_name + " api description"
	c.Api_version = "0.0.0"
	c.Host = "0.0.0.0"
	c.Port = "8000"
	c.Api_reload = true
	c.Db.Db_name = c.Api_name + "_db"
	c.Db.Sgbd = "sqlite"
	c.Db.Async = false

	var Attributs_list []Attributs
	Attributs_list = append(Attributs_list, Attributs{
		Attribut_name:  "name",
		Attributs_type: "default_type",
	})

	var model_default []Models_config
	model_default = append(model_default, Models_config{
		Model_name: "model",
		Attributs:  Attributs_list,
	})

	c.Models = append(c.Models, model_default...)
	c.Template_dirs = "template"
	c.Static_files_dirs = "static_file"

	config, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", config, 0644)
}

func (c *Config_file) ReadConfig() error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}
