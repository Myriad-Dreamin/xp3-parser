package main

import (
	"os"
	"time"

	"github.com/Myriad-Dreamin/xp3-parser/decoder"
	logger "github.com/Myriad-Dreamin/xp3-parser/log"
)

var path = "do not existing path?"

func init() {
	path = os.Args[1]
}

func main() {
	decoder.Decode(path)
	time.Sleep(logger.FlushTime)
	time.Sleep(logger.FlushTime)
	logger.Infof("success\n")
}
