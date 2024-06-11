package header

import (
	"encoding/binary"
	"os"
	"unsafe"
	
	"github.com/The-EpaG/DataToVideo/internal/errors"
)

type Header struct {
	HeaderSize uint64
	Size       uint64
	Filename   string
}

func (header *Header) bytesSize() (int) {
	size := int(unsafe.Sizeof(uint64(0)))
	size *= 2
    size += len(header.Filename)

	return size
}

func (header *Header) ToBytes() []byte {
	var headerBytes []byte = make([]byte, 0, header.bytesSize())

	var tmp []byte = make([]byte, unsafe.Sizeof(uint64(0)))

	binary.BigEndian.PutUint64(tmp, header.HeaderSize)
	headerBytes = append(headerBytes, tmp...)

	binary.BigEndian.PutUint64(tmp, header.Size)
	headerBytes = append(headerBytes, tmp...)

	headerBytes = append(headerBytes, []byte(header.Filename)...)

	return headerBytes
}

func FromBytes(buffer []byte) (*Header, error) {
	if len(buffer) == 0 {
		return nil, &errors.EmptyBufferError{}
	}

	var header Header = Header{}
	uint64Size := int(unsafe.Sizeof(uint64(0)))

	if len(buffer) < uint64Size*2 {
		return nil, &errors.BufferTooSmallError{Size: len(buffer), MinSize: uint64Size*2}
	}

	header.HeaderSize = binary.BigEndian.Uint64(buffer[:uint64Size])
	header.Size = binary.BigEndian.Uint64(buffer[uint64Size:uint64Size*2])
	header.Filename = string(buffer[uint64Size*2:header.HeaderSize])

	return &header, nil
}

func New(file *os.File, filename string) (*Header, error) {
	var header Header = Header{0, 0, filename}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := fi.Size()
	header.Size = uint64(size)

	header.HeaderSize = uint64(header.bytesSize())

	return &header, nil
}