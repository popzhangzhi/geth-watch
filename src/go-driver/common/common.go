package common

import "log"

func Output(Ftype int8, arg ...interface{}) {
	switch Ftype {
	case 1:
		log.Println(arg)
	case 2:
		log.Fatalln(arg)
	}

}
