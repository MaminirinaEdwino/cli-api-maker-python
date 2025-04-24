package apimaker

import (
	"fmt"
	"os"

	configoutils "github.com/MaminirinaEdwino/cli_api_maker/outils/config_outils"
)

func ImportRouter(model_name string) string {
	return fmt.Sprintf("from %s.route import %s_router", model_name, model_name)
}

func IncludeRouter(model_name string) string {
	return fmt.Sprintf("app.include_router(%s_router, tags=['%s'])", model_name, model_name)
}



func CodeWriter(c* configoutils.Config_file)  {
	var code []string

	DEFAULT_IMPORT := []string{
		"from fastapi import FastAPI, Depends, HTTPException, status",
		"import uvicorn",
		"from fastapi.middleware.cors import CORSMiddleware",
		"from db import Base, engine, SessionLocal, get_db",
		"from fastapi.security import  OAuth2PasswordRequestForm",
		"from datetime import timedelta",
		"from user import UserDB, UserCreate, User, Token",
		"from security import *",
	}

	user_code := `
from db import Base
from pydantic import BaseModel
from sqlalchemy import Column, Integer, String, Boolean
from typing import Optional

class UserDB(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String, unique=True, index=True)
    email = Column(String, unique=True, index=True)
    full_name = Column(String)
    role = Column(String, default="user") 
    is_active = Column(Boolean, default=True)
    hashed_password = Column(String)


# --- Data Models (Pydantic) ---
class User(BaseModel):
    username: str
    email: str
    full_name: Optional[str] = None
    is_active: bool = True
    role: str = "user"  # Default role is 'user'
    class Config:
        from_attributes = True

class UserCreate(BaseModel):
    username: str
    password: str
    email: str
    full_name: Optional[str] = None
    role: str = "user"  # Default role is 'user'
    class Config:
        from_attributes = True

class Token(BaseModel):
    access_token: str
    token_type: str
	`
	user_route_code := `
@app.post("/users/", response_model=User, status_code=201)
async def create_user(user: UserCreate, db: SessionLocal = Depends(get_db)):
    db_user = await get_user_by_username(db, username=user.username)
    if db_user:
        raise HTTPException(status_code=400, detail="Username already registered")
    return await create_user_db(db, user)

@app.post("/token", response_model=Token)
async def login_for_access_token(form_data: OAuth2PasswordRequestForm = Depends(), db: SessionLocal = Depends(get_db)):
    user = await get_user_by_username(db, username=form_data.username)
    if not user:
        raise HTTPException(status_code=400, detail="Incorrect username or password")
    if not verify_password(form_data.password, user.hashed_password):
        raise HTTPException(status_code=400, detail="Incorrect username or password")
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(data={"sub": user.username}, expires_delta=access_token_expires)
    return {"access_token": access_token, "token_type": "bearer"}

@app.get("/users/me/", response_model=User)
async def read_users_me(current_user: UserDB = Depends(get_current_active_user)):
    return current_user

@app.get("/protected/", dependencies=[Depends(get_current_active_user)])
async def protected_route():
    return {"message": "This route is protected!"}

@app.get("/public/")
async def public_route():
    return {"message": "This route is public."}
	`

	security_code := `
from fastapi import Depends, HTTPException, status
from jose import JWTError, jwt
from passlib.context import CryptContext
from fastapi.security import OAuth2PasswordBearer
from datetime import datetime, timedelta
from typing import Optional
from streamlit import status
from db import SessionLocal, get_db
from user import * 
import os


SECRET_KEY = os.environ.get("SECRET_KEY", "super-secret-key")
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

# --- Security Utilities ---
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

def verify_password(plain_password: str, hashed_password: str) -> bool:
    return pwd_context.verify(plain_password, hashed_password)

def get_password_hash(password: str) -> str:
    return pwd_context.hash(password)

def create_access_token(data: dict, expires_delta: Optional[timedelta] = None):
    to_encode = data.copy()
    expire = datetime.utcnow() + (expires_delta or timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES))
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt

async def get_user_by_username(db: SessionLocal, username: str):
    return db.query(UserDB).filter(UserDB.username == username).first()

async def get_user(db: SessionLocal, user_id: int):
    return db.query(UserDB).filter(UserDB.id == user_id).first()

async def create_user_db(db: SessionLocal, user: "UserCreate"):
    hashed_password = get_password_hash(user.password)
    db_user = UserDB(
        username=user.username,
        email=user.email,
        full_name=user.full_name,
        hashed_password=hashed_password,
    )
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    return db_user

async def get_current_user(db: SessionLocal = Depends(get_db), token: str = Depends(oauth2_scheme)):
    credentials_exception = HTTPException(
        status_code=401,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        username: str = payload.get("sub")
        if username is None:
            raise credentials_exception
    except JWTError:
        raise credentials_exception
    user = await get_user_by_username(db, username=username)
    if user is None:
        raise credentials_exception
    return user

async def get_current_active_user(current_user: UserDB = Depends(get_current_user)):
    if not current_user.is_active:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Inactive user")
    return current_user
async def get_current_active_admin(current_user: UserDB = Depends(get_current_user)):
    if not current_user.is_active or current_user.role != "admin":
        raise HTTPException(status_code=400, detail="Inactive user or not an admin")
    return current_user
async def get_current_active_user_or_admin(current_user: UserDB = Depends(get_current_user)):
    if not current_user.is_active or current_user.role not in ["admin", "user"]:
        raise HTTPException(status_code=400, detail="Inactive user or not an admin")
    return current_user

	`
	api_declaration := fmt.Sprintf("app = FastAPI(title='%s', description='%s', version='%s') \n", c.Api_name, c.Api_description, c.Api_version)
	api_declaration += fmt.Sprintf("app.add_middleware(CORSMiddleware, allow_origins=['*'], allow_credentials=True, allow_methods=['*'], allow_headers=['*'])\nBase.metadata.create_all(bind=engine)\n")
	api_launcher := fmt.Sprintf(`if __name__ == '__main__':
	uvicorn.run('main:app', host='%s', port=%s, reload=True)`, c.Host, c.Port)

	for _, model := range c.Models {
		code = append(code, ImportRouter(model.Model_name))
	}
	code = append(code, api_declaration)

	code = append(code, user_route_code)

	for _, model := range c.Models {
		code = append(code, IncludeRouter(model.Model_name))
	}

	code = append(code, api_launcher)

	os.Mkdir(c.Api_name, os.ModePerm)
	file, err := os.OpenFile(c.Api_name+"/main.py", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	ErrorFunc(err)
	defer file.Close()

	user_dile, err_user_file := os.OpenFile(c.Api_name+"/user.py", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	ErrorFunc(err_user_file)
	defer user_dile.Close()

	security_file, err_security := os.OpenFile(c.Api_name+"/security.py", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	ErrorFunc(err_security)
	defer security_file.Close()

	db_file, db_err := os.OpenFile(c.Api_name+"/db.py", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	ErrorFunc(db_err)
	defer db_file.Close()

	//writing the db_file
	DataBaseWriter(&c.Db, db_file, c.Api_name+"/")
	WriteAllModel(c.Models, c.Api_name+"/")
	// Write the code to the file
	_, err = file.WriteString("#!/usr/bin/env python3\n")
	ErrorFunc(err)
	_, err = file.WriteString("# -*- coding: utf-8 -*-\n")
	ErrorFunc(err)

	//ecriture de l'import et des codes Base de l'api
	FromArrayWriter(DEFAULT_IMPORT, file)
	FromStringWriter(user_code, user_dile)
	FromStringWriter(security_code, security_file)
	FromArrayWriter(code, file)
}