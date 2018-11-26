package controller

import "go-driver/blockDriver"
import log "go-driver/common"

var (
	systemNum = 3
	userNum   = 10
)

func init() {}

func Generate() {

	for i := 0; i < systemNum+userNum; i++ {
		addr, priKey, err := blockDriver.DoCreate()
		if err != nil {
			log.Output(2, "number "+(string(i)), err)
			break
		}
		log.Output(1, addr, priKey)
	}

}
