package apimaker

import (
	"os"
)

func FromArrayWriter(array []string, file *os.File)  {
	for i := 0; i < len(array); i++ {
		_,err := file.WriteString(array[i]+"\n")  
		ErrorFunc(err)
	}
}