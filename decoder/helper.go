package decoder

import (
	"io"
	"os"
	"path/filepath"

	logger "github.com/Myriad-Dreamin/xp3-parser/log"
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

func checkAvailable(fileName string) bool {
	for _, charf := range fileName {
		if charf == '$' {
			logger.Warnln("escaping", fileName)
			return false
		}
	}
	return true
}

func EscapeSpace(fileName []byte) []byte {
	// var buffer = bytes.NewBuffer(fileName)
	// var wb = bytes.NewBuffer(make([]byte, 512))
	// wb.Reset()
	// var err error
	// var b byte
	// for {
	// 	b, err = buffer.ReadByte()
	// 	if err != nil {
	// 		break
	// 	}
	// 	if b != ' ' {
	// 		// wb.WriteByte('\\')
	// 		wb.WriteByte(b)
	// 	}
	// }
	// if err == io.EOF {
	// 	return wb.Bytes()
	// }
	// logger.Fatal(err)
	return fileName
}

func createAndWrite(fileName []byte, buffer io.Reader) {

	var decodedFileName = string(EscapeSpace(UTF16ToUTF8(fileName)))

	if !checkAvailable(decodedFileName) {
		return
	}
	var wantedFileName = "./decode/" + decodedFileName
	var wantedDir = filepath.Dir(wantedFileName)
	if _, err := os.Stat(wantedDir); os.IsNotExist(err) {
		os.MkdirAll(wantedDir, 0755)
	} else if err != nil {
		logger.Fatal(err)
		return
	}
	if f, err := os.Create(wantedFileName); err != nil {
		logger.Fatal(err)
		return
	} else {
		defer f.Close()
		_, err = io.TeeReader(buffer, f).Read(buf)
		if err != nil {
			logger.Fatal(err)
			return
		}
		logger.Infof("extracted file: %v...\n", wantedFileName)
	}
}
