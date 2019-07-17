package decoder

import (
	"io"
	"log"
	"os"
)

type callBackFunction func(*os.File)

func GetFileReader(path string, cb callBackFunction) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	cb(f)
	return nil
}

func createAndWrite(fileName []byte, buffer io.Reader) {
	if f, err := os.Create("./decode/" + string(UTF16ToUTF8(fileName))); err != nil {
		log.Fatal(err)
		return
	} else {
		defer f.Close()
		_, err = io.TeeReader(buffer, f).Read(buf)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
