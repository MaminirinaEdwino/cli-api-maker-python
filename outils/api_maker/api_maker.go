package apimaker

import configoutils "github.com/MaminirinaEdwino/cli_api_maker/outils/config_outils"

func Api_maker_simple_structure()  {
	conf := configoutils.Config_file{}
	conf.ConfigGenerator()
	conf.ReadConfig()
	CodeWriter(&conf)
}