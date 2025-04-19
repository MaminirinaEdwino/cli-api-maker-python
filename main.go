package main

import (
	"flag"
	"fmt"

	apimaker "github.com/MaminirinaEdwino/cli_api_maker/outils/api_maker"
	configoutils "github.com/MaminirinaEdwino/cli_api_maker/outils/config_outils"
)

func main() {
	fmt.Println("Api maker")

	generate_config := flag.Bool("generate_config", false, "Generate A base Config file")
	make_api := flag.Bool("make_api", false, "Make an api from the config file")
	flag.Parse()
	if *generate_config {
		fmt.Println("Generating config file")	
		conf := configoutils.Config_file{}
		conf.ConfigGenerator()
	}
	if *make_api {
		fmt.Println("Generating api from config file")
		conf := configoutils.Config_file{}
		conf.ReadConfig()
		apimaker.CodeWriter(&conf)
		
	}
	fmt.Println("By Edwino")
	// conf := configoutils.Config_file{}
	// fmt.Println("print config")
	// conf.ReadConfig()
	// fmt.Println(conf)
	// apimaker.CodeWriter(&conf)
}