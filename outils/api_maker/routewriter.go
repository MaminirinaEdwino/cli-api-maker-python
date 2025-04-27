package apimaker

import (
	"fmt"
	"os"

	configoutils "github.com/MaminirinaEdwino/cli_api_maker/outils/config_outils"
)

func WriteGetAll(model_name string) string {
	return fmt.Sprintf(`
@%s_router.get("/")
async def get_all_%s(db: Session = Depends(get_db)):
	db_%s = db.query(%s).all()	
	return db_%s`, model_name, model_name, model_name, model_name, model_name)
}

func WriteGetById(model_name string) string {
	return fmt.Sprintf(`
@%s_router.get("/{id}")
async def get_%s_by_id(id: int, db: Session = Depends(get_db)):
	db_%s = db.query(%s).filter(%s.id == id).first()
	if not db_%s:
		raise HTTPException(status_code=404, detail="%s not found")
	return db_%s`, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name)
}

func WritePost(model_name string) string {
	return fmt.Sprintf(`
@%s_router.post("/")
async def create_%s(%s_post: %s_create, db: Session = Depends(get_db)):
	db_%s = %s(**%s_post.model_dump())
	db.add(db_%s)
	db.commit()
	db.refresh(db_%s)
	return db_%s`, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name)

}
func WritePut(model_name string) string {
	return fmt.Sprintf(`
@%s_router.put("/{id}")
async def update_%s(id: int, %s_put: %s_update, db: Session = Depends(get_db)):
	db_%s = db.query(%s).filter(%s.id == id).first()
	if not db_%s:
		raise HTTPException(status_code=404, detail="%s not found")
	for key, value in %s_put.model_dump().items():
		if value is not None:
			setattr(db_%s, key, value)
	db.commit()
	db.refresh(db_%s)
	return db_%s`, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name)
}

func WriteDelete(model_name string) string {
	return fmt.Sprintf(`
@%s_router.delete("/{id}")
async def delete_%s(id: int, db: Session = Depends(get_db)):
	db_%s = db.query(%s).filter(%s.id == id).first()
	if not %s:
		raise HTTPException(status_code=404, detail="%s not found")
	db.delete(db_%s)
	db.commit()
	return {"message": "%s deleted successfully"}`, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name, model_name)
}
func WriteGetByField(model_name string, field_name string) string {
	return fmt.Sprintf(`
@%s_router.get("/{%s}")
async def get_%s_by_%s(%s: str, db: Session = Depends(get_db)):
	db_%s = db.query(%s).filter(%s.%s == %s).all()	
	if not db_%s:
		raise HTTPException(status_code=404, detail="%s not found")
	return db_%s`, model_name, field_name, model_name, field_name, field_name, model_name, model_name, model_name, field_name, field_name, model_name, model_name, model_name)
}
func RouteWriter(c *configoutils.Models_config, api_dir string) {
	var code []string

	DEFAULT_IMPORT := "from fastapi import APIRouter,Depends, HTTPException \nfrom security import *\n"
	DEFAULT_IMPORT += fmt.Sprintf("from %s.model import %s, %s_create, %s_update \n", c.Model_name, c.Model_name, c.Model_name, c.Model_name)
	DEFAULT_IMPORT += fmt.Sprintln("from db import get_db")
	DEFAULT_IMPORT += fmt.Sprintln("from requests import Session")

	code = append(code, fmt.Sprintf("%s_router = APIRouter(prefix=\"/%s\", tags=['%s'], dependencies=[Depends(get_current_active_user)])\n", c.Model_name, c.Model_name, c.Model_name))
	code = append(code, WriteGetAll(c.Model_name))
	code = append(code, WriteGetById(c.Model_name))
	code = append(code, WritePost(c.Model_name))
	code = append(code, WritePut(c.Model_name))
	code = append(code, WriteDelete(c.Model_name))
	for _, field := range c.Attributs {
		code = append(code, WriteGetByField(c.Model_name, field.Attribut_name))
	}
	file, err := os.OpenFile(api_dir+c.Model_name+"/"+"route.py", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	ErrorFunc(err)
	FromStringWriter(DEFAULT_IMPORT, file)
	FromArrayWriter(code, file)
	defer file.Close()
}
