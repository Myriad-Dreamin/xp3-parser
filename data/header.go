package data

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	// 2.8 and 3.0
	h1 = []byte{0x58, 0x50, 0x33, 0x0d, 0x0a, 0x20, 0x0a, 0x1a}
	h2 = []byte{0x8b, 0x67, 0x01}

	// 3.0
	cushion_index        = []byte{0x17, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	header_minor_version = []byte{0x01, 0x00, 0x00, 0x00}
	cushion_header       = []byte{0x80}
	index_size           = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

type XP3Header interface {
	GetFileHeaderOffset() uint64
}

type XP3Header28 struct {
	Header1          [8]byte
	Header2          [3]byte
	FileHeaderOffset uint64
}

type XP3Header30 struct {
	Header1            [8]byte
	Header2            [3]byte
	CushionIndex       [8]byte
	HeaderMinorVersion [4]byte
	CushionHeader      [1]byte
	IndexSize          [8]byte
	FileHeaderOffset   uint64
}

func (h *XP3Header28) GetFileHeaderOffset() uint64 { return h.FileHeaderOffset }
func (h *XP3Header30) GetFileHeaderOffset() uint64 { return h.FileHeaderOffset }

func (h *XP3Header28) ReadFrom(r io.Reader) (nn int64, err error) {
	err = binary.Read(r, KiriKiriEndian, h)
	if err != nil {
		return 0, err
	}
	if !bytes.Equal(h1, h.Header1[:]) {
		return 0, errors.New("not equal of magic header 1(0-7 bytes)(kirikiri 2.28)")
	}
	if !bytes.Equal(h2, h.Header2[:]) {
		return 8, errors.New("not equal of magic header 2(8-10 bytes)(kirikiri 2.28)")
	}

	return 19, nil
}

func (h *XP3Header30) ReadFrom(r io.Reader) (nn int64, err error) {
	err = binary.Read(r, KiriKiriEndian, &h.Header1)
	if err != nil {
		return 0, err
	}
	if !bytes.Equal(h1, h.Header1[:]) {
		return 0, errors.New("not equal of magic header 1(0-7 bytes)(kirikiri 2.30)")
	}
	err = binary.Read(r, KiriKiriEndian, &h.Header2)
	if err != nil {
		return 0, err
	}
	if !bytes.Equal(h2, h.Header2[:]) {
		return 8, errors.New("not equal of magic header 2(8-10 bytes)(kirikiri 2.30)")
	}
	err = binary.Read(r, KiriKiriEndian, &h.CushionIndex)
	if err != nil {
		return 8, err
	}
	if !bytes.Equal(cushion_index, h.CushionIndex[:]) {
		return 11, errors.New("not equal of cushion index (11-18 bytes)(kirikiri 2.30)")
	}
	err = binary.Read(r, KiriKiriEndian, &h.HeaderMinorVersion)
	if err != nil {
		return 11, err
	}
	if !bytes.Equal(header_minor_version, h.HeaderMinorVersion[:]) {
		return 19, errors.New("not equal of header minor version (19-22 bytes)(kirikiri 2.30)")
	}
	err = binary.Read(r, KiriKiriEndian, &h.CushionHeader)
	if err != nil {
		return 23, err
	}
	if !bytes.Equal(cushion_header, h.CushionHeader[:]) {
		return 23, errors.New("not equal of cushion header (23-23 bytes)(kirikiri 2.30)")
	}
	err = binary.Read(r, KiriKiriEndian, &h.IndexSize)
	if err != nil {
		return 24, err
	}
	if !bytes.Equal(index_size, h.IndexSize[:]) {
		return 24, errors.New("not equal of index size (24-31 bytes)(kirikiri 2.30)")
	}
	return 40, nil
}

func (h *XP3Header28) ReReadFrom(r io.Reader) (nn int64, err error) {
	// not need to set header
	err = binary.Read(r, KiriKiriEndian, &h.FileHeaderOffset)
	if err != nil {
		return 11, err
	}
	return 19, nil
}

func ParseHeader(r *os.File) XP3Header {
	var header = new(XP3Header30)
	n, err := header.ReadFrom(r)
	if err != nil && n != 11 {
		fmt.Println(n, err)
		log.Fatal(err)
		return nil
	} else if err != nil && n == 11 {
		var header2 = new(XP3Header28)
		header2.Header1 = header.Header1
		header2.Header2 = header.Header2
		header2.FileHeaderOffset = KiriKiriEndian.Uint64(header.CushionIndex[:])
		return header2
	}
	return header
}
