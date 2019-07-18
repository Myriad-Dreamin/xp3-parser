package decoder

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"

	data "github.com/Myriad-Dreamin/xp3-parser/data"
	logger "github.com/Myriad-Dreamin/xp3-parser/log"
)

const (
	maxBufferSize = 1024 * 1024 * 64
)

var buf, buf2 []byte
var fileDescriptors []*data.FileDescriptor

func init() {
	buf = make([]byte, maxBufferSize)
	buf2 = make([]byte, maxBufferSize)
}

func Decode(toDecode string) {
	var Parse = func(r *os.File) {
		var n int64
		var header = data.ParseHeader(r)
		_, err := r.ReadAt(buf, int64(header.GetFileHeaderOffset()))
		if err != nil && err != io.EOF {
			logger.Fatal(err)
			return
		}
		var fileheader = new(data.FileHeader)
		n, err = fileheader.ReadFrom(bytes.NewBuffer(buf))
		if err != nil && err != io.EOF {
			logger.Fatal(err)
			return
		}
		if fileheader.HeaderSize+uint64(n) > uint64(len(buf)) {
			buf = append(buf, make([]byte, fileheader.HeaderSize+uint64(n)-uint64(len(buf)))...)
		}
		_, err = r.ReadAt(buf, int64(header.GetFileHeaderOffset())+int64(n))
		if err != nil && err != io.EOF {
			logger.Fatal(err)
			return
		}

		var rr io.Reader
		if fileheader.Compressed != 0 {
			rr, err = zlib.NewReader(bytes.NewBuffer(buf[:fileheader.HeaderSize]))
			if err != nil && err != io.EOF {
				logger.Fatal(err)
				return
			}
		} else {
			rr = bytes.NewBuffer(buf)
		}

		var success = true
		for success {
			var fileDescriptor = new(data.FileDescriptor)

			_, err = fileDescriptor.ReadFrom(rr)
			if err != nil {
				if err != io.EOF {
					logger.Fatal(err)
				}
				return
			}

			fileDescriptors = append(fileDescriptors, fileDescriptor)
		}
	}

	var writeToFiles = func(r *os.File) {
		var rr io.Reader
		var err error
		for _, descriptor := range fileDescriptors {
			var mbuf = bytes.NewBuffer(buf2)
			mbuf.Reset()
			for _, segment := range descriptor.SegmentSection.Segments {
				if segment.CompressedSize > uint64(len(buf)) {
					buf = append(buf, make([]byte, segment.CompressedSize-uint64(len(buf)))...)
				}
				_, err = r.ReadAt(buf[:segment.CompressedSize], int64(segment.Offset))
				if err != nil && err != io.EOF {
					logger.Fatal(err)
					return
				}
				if segment.Flag != 0 {
					rr, err = zlib.NewReader(bytes.NewBuffer(buf[:segment.CompressedSize]))
					if err != nil && err != io.EOF {
						logger.Fatal(err)
						return
					}
				} else {
					rr = bytes.NewBuffer(buf[:segment.CompressedSize])
				}

				mbuf.ReadFrom(rr)
			}
			createAndWrite(descriptor.InfoSection.FileName, mbuf)
		}
	}

	if err := GetFileReader(toDecode, func(r *os.File) {
		Parse(r)
		writeToFiles(r)
	}); err != nil {
		logger.Fatal(err)
		return
	}
}
