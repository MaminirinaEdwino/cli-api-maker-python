package apimaker

import (
	"os"

	configoutils "github.com/MaminirinaEdwino/cli_api_maker/outils/config_outils"
)

func Db_Writer(c *configoutils.Db_config, db_url string, file *os.File) {
	DEFAULT_IMPORT := `
from fastapi import Depends
from sqlalchemy import create_engine, Column, Integer, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
	`
	DATABASE_URL := "DATABASE_URL = '" + db_url + "'\n"
	code := `
engine = create_engine(DATABASE_URL, connect_args={"check_same_thread": False})
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()

def get_db():
   db = SessionLocal()
   try:
       yield db
   finally:
       db.close()
       
db_dependency = Depends(get_db)
	`
	FromStringWriter(DEFAULT_IMPORT, file)
	FromStringWriter(DATABASE_URL, file)
	FromStringWriter(code, file)
}



func DataBaseWriter(c *configoutils.Db_config, file *os.File, api_dir string) {
	var DATABASE_URL string
	if c.Sgbd == "sqlite" {
		// code := fmt.Sprintf("DATABASE_URL = 'sqlite:///%s.db'", c.Db_name)
		DATABASE_URL= "sqlite:///" + c.Db_name + ".db"
		os.Create(api_dir+c.Db_name + ".db")
		Db_Writer(c, DATABASE_URL, file)
	} else if c.Sgbd == "mysql" {
		DATABASE_URL = "mysql+pymysql://postgres:root@localhost:5432" + c.Db_name + ".db"
		Db_Writer(c, DATABASE_URL, file)
	} else if c.Sgbd == "postgresql" {
		DATABASE_URL = "postgresql://postgres:root@localhost:5432" + c.Db_name + ".db"
		Db_Writer(c, DATABASE_URL, file)
	}
}
