package apimaker

import (
	"os"
)

func FromStringWriter(code string,file *os.File)  {
	_,err := file.WriteString(code+"\n")  
	ErrorFunc(err)
}