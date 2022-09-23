package bitreader

import (
	"fmt"
	"io"
)

const (
	mask1 = uint8(0b10000000)
	mask2 = uint8(0b11000000)
	mask3 = uint8(0b11100000)
	mask4 = uint8(0b11110000)
	mask5 = uint8(0b11111000)
	mask6 = uint8(0b11111100)
	mask7 = uint8(0b11111110)
	mask8 = uint8(0b11111111)
)

type Reader struct {
	reader    io.Reader
	buffer    []byte
	offset    int64
	bitOffset int64
}

func New(reader io.Reader) (*Reader, error) {
	r := &Reader{reader: reader, buffer: nil, bitOffset: 8, offset: -1}
	if err := r.nextByte(); err != nil {
		return nil, fmt.Errorf("read.New: %w", err)
	}
	return r, nil
}

func (r *Reader) readMore(byteSize int) error {
	temp := make([]byte, byteSize)
	_, err := r.reader.Read(temp)
	if err != nil {
		return err
	}
	r.buffer = append(r.buffer, temp...)
	return nil
}

func (r *Reader) nextByte() error {
	r.bitOffset = 8
	r.offset++
	if r.offset >= int64(len(r.buffer)) {
		if err := r.readMore(1); err != nil {
			return fmt.Errorf("Reader.nextByte: %w", err)
		}
	}
	return nil
}

func (r *Reader) readNBits(bitSize int64) (uint8, uint8, int64) {
	if r.bitOffset == 0 {
		return 0, 0, bitSize
	}
	min := bitSize
	if min > r.bitOffset {
		min = r.bitOffset
	}
	mask := uint8(0)
	switch min {
	case 1:
		mask = mask1
	case 2:
		mask = mask2
	case 3:
		mask = mask3
	case 4:
		mask = mask4
	case 5:
		mask = mask5
	case 6:
		mask = mask6
	case 7:
		mask = mask7
	case 8:
		mask = mask8
	}
	value := (r.buffer[r.offset] << byte(8-r.bitOffset)) & mask
	r.bitOffset -= min
	return value, uint8(min), bitSize - min
}

func (r *Reader) Read(bitSize int64) ([]uint8, error) {
	result := []uint8(nil)
	value := uint8(0)
	remain := bitSize
	read := uint8(0)
	prevRead := uint8(0)
	for remain != 0 {
		value, read, remain = r.readNBits(remain)
		if read == 0 {
			if err := r.nextByte(); err != nil {
				return result, fmt.Errorf("Reader.Read: %w", err)
			}
			continue
		}
		if prevRead > 0 && prevRead < 8 {
			result[len(result)-1] |= value >> prevRead
			read -= 8 - prevRead
			value <<= 8 - prevRead
		}
		if read > 0 && read < 8 {
			result = append(result, value)
		}
		prevRead = read
	}
	return result, nil
}

func (r *Reader) Reset() {
	r.bitOffset = 8
	if len(r.buffer) == 0 {
		r.offset = -1
		return
	}
	r.offset = 0
}

func (r *Reader) MoveTo(bitPos int64) error {
	bitOffset := bitPos % 8
	offset := bitPos / 8
	if offset > int64(len(r.buffer)) {
		return fmt.Errorf("Reader.Seek: offset %d is out of range", offset)
	}
	r.offset = offset
	r.bitOffset = 8 - bitOffset
	return nil
}

func (r *Reader) Peek(bitSize int64) ([]uint8, error) {
	offset := r.offset
	bitOffset := r.bitOffset
	defer func() {
		r.offset = offset
		r.bitOffset = bitOffset
	}()
	result, err := r.Read(bitSize)
	if err != nil {
		return nil, fmt.Errorf("Reader.Peek: %w", err)
	}
	return result, nil
}
