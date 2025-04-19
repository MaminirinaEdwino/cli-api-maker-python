package apimaker

import (
	"fmt"
	"os"

	configoutils "github.com/MaminirinaEdwino/cli_api_maker/outils/config_outils"
)

func ReturnTheRightType(type_ string) string {
	switch type_ {
	case "int":
		return "Integer"
	case "str":
		return "String(255)"
	case "json":
		return "JSON"
	case "date":	
		return "String(255)"
	case "bool": 
		return "Boolean"
	default:
		return type_
	}
}

func ModelWriter(c* configoutils.Models_config, api_dir string) {
	var code []string
	DEFAULT_IMPORT :=`
from pydantic import BaseModel
from typing import Optional
from db import Base
from sqlalchemy import Column, ForeignKey, Integer, String, JSON, Boolean
	`
	//code = append(code, fmt.Sprintf("class %s(BaseModel):\n", c.Name))
	code = append(code, fmt.Sprintf("class %s_create(BaseModel):\n", c.Model_name))
	for _, field := range c.Attributs {
		code = append(code, fmt.Sprintf("\t%s: %s", field.Attribut_name, field.Attributs_type))
	}
	code = append(code, "\n")

	code = append(code, fmt.Sprintf("class %s_update(BaseModel):\n", c.Model_name))
	for _, field := range c.Attributs {
		code = append(code, fmt.Sprintf("\t%s: Optional[%s] = None", field.Attribut_name, field.Attributs_type))
	}
	code = append(code, "\n")


	os.Mkdir(fmt.Sprintf("%s%s", api_dir, c.Model_name), os.ModePerm)
	file, err := os.OpenFile(api_dir+c.Model_name+"/"+
	"model.py", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	ErrorFunc(err)
	defer file.Close()

	code = append(code, fmt.Sprintf("class %s(Base):\n", c.Model_name))
	code = append(code, fmt.Sprintf("\t__tablename__ = '%s'\n", c.Model_name))
	code = append(code, fmt.Sprintf("\tid = Column(Integer, primary_key=True, index=True)\n"))

	for _, field := range c.Attributs {
		code = append(code, fmt.Sprintf("\t%s= Column(%s, nullable= False)", ReturnTheRightType(field.Attribut_name), ReturnTheRightType(field.Attributs_type)))
	}


	_, err = file.WriteString("#!/usr/bin/env python3\n")
	ErrorFunc(err)
	_, err = file.WriteString("# -*- coding: utf-8 -*-\n")
	ErrorFunc(err)

	FromStringWriter(DEFAULT_IMPORT, file)
	FromArrayWriter(code, file)
}

func WriteAllModel(models []configoutils.Models_config, api_dir string) {
	for _, model := range models {
		ModelWriter(&model, api_dir)
		RouteWriter(&model, api_dir)
	}
}