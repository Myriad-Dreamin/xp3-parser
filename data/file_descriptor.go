package data

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	wcharSize = 2
)

type FileSection struct {
	Header [4]byte
	Size   uint64
}

type InfoSection struct {
	Header         [4]byte
	Size           uint64
	Flag           uint32
	OriginalSize   uint64
	CompressedSize uint64
	FileLength     uint16
	FileName       []byte
}

type SegmentSection struct {
	Header   [4]byte
	Size     uint64
	Segments []Segment
}

type Segment struct {
	Flag           uint32
	Offset         uint64
	OriginalSize   uint64
	CompressedSize uint64
}

type Adler32Section struct {
	Header  [4]byte
	Size    uint64
	Padding [4]byte
}

type FileDescriptor struct {
	FileSection    FileSection
	InfoSection    InfoSection
	SegmentSection SegmentSection
	Adler32Section Adler32Section
}

func (fs *FileSection) ReadFrom(r io.Reader) (nn int64, err error) {
	err = binary.Read(r, KiriKiriEndian, fs)
	if err != nil {
		return 0, err
	}
	return 12, nil
}

func (is *InfoSection) ReadFrom(r io.Reader) (nn int64, err error) {
	is.FileName = make([]byte, 0, 0)
	err = binary.Read(r, KiriKiriEndian, &is.Header)
	if err != nil {
		return 0, err
	}
	err = binary.Read(r, KiriKiriEndian, &is.Size)
	if err != nil {
		return 0, err
	}
	err = binary.Read(r, KiriKiriEndian, &is.Flag)
	if err != nil {
		return 0, err
	}
	err = binary.Read(r, KiriKiriEndian, &is.OriginalSize)
	if err != nil {
		return 0, err
	}
	err = binary.Read(r, KiriKiriEndian, &is.CompressedSize)
	if err != nil {
		return 0, err
	}
	err = binary.Read(r, KiriKiriEndian, &is.FileLength)
	if err != nil {
		return 0, err
	}
	is.FileLength *= wcharSize
	is.FileName = make([]byte, is.FileLength, is.FileLength)
	err = binary.Read(r, KiriKiriEndian, &is.FileName)
	if err != nil {
		return 0, err
	}
	return 34 + int64(is.FileLength), nil
}

func (ss *SegmentSection) ReadFrom(r io.Reader) (nn int64, err error) {
	err = binary.Read(r, KiriKiriEndian, &ss.Header)
	if err != nil {
		return 0, err
	}
	err = binary.Read(r, KiriKiriEndian, &ss.Size)
	if err != nil {
		return 0, err
	}
	if ss.Size%28 != 0 {
		return 0, errors.New("the size of segment section is not the mutiple of 28")
	}
	ss.Segments = make([]Segment, ss.Size/28, ss.Size/28)
	err = binary.Read(r, KiriKiriEndian, &ss.Segments)
	if err != nil {
		return 0, err
	}
	return 12 + int64(ss.Size), nil
}

func (as *Adler32Section) ReadFrom(r io.Reader) (nn int64, err error) {
	err = binary.Read(r, KiriKiriEndian, as)
	if err != nil {
		return 0, err
	}
	return 16, nil
}

func (d *FileDescriptor) ReadFrom(r io.Reader) (nn int64, err error) {
	var n int64
	n, err = d.FileSection.ReadFrom(r)
	if err != nil {
		return 0, err
	}
	nn += n
	n, err = d.InfoSection.ReadFrom(r)
	if err != nil {
		return 0, err
	}
	nn += n
	n, err = d.SegmentSection.ReadFrom(r)
	if err != nil {
		return 0, err
	}
	nn += n
	n, err = d.Adler32Section.ReadFrom(r)
	if err != nil {
		return 0, err
	}
	return nn + n, nil
}
