package apimaker

func ErrorFunc(err error)  {
	if err != nil {	
		panic(err)
	}
}

