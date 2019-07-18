package decoder

import (
	"bytes"
	"io/ioutil"
	"strings"

	logger "github.com/Myriad-Dreamin/xp3-parser/log"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func trimZeros(fileName string) string {
	return strings.Replace(fileName, "\x00", "", -1)
}

func UTF16ToUTF8(b []byte) []byte {
	win16be := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())
	unicodeReader := transform.NewReader(bytes.NewReader(b), utf16bom)
	decoded, err := ioutil.ReadAll(unicodeReader)
	if err != nil {
		logger.Fatal(err)
		return nil
	}
	return decoded
}
