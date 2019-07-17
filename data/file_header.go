package data

import (
	"encoding/binary"
	"io"
)

type FileHeader struct {
	Compressed   uint8
	HeaderSize   uint64
	OriginalSize uint64
}

func (h *FileHeader) ReadFrom(r io.Reader) (nn int64, err error) {
	err = binary.Read(r, KiriKiriEndian, h)
	if err != nil {
		return 0, err
	}
	return 17, nil
}
