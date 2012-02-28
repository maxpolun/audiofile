package audiofile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type AudioReader interface {
	Load(io.Reader) error
	GetData() []byte
}
type AudioWriter interface {
	Save(io.Writer) error
	Init() // set up the initial headers, etc.
	SetData([]byte)
}
type AudioFile interface {
	AudioReader
	AudioWriter
}

const (
	MIN_8_BIT  = 0
	MID_8_BIT  = 128
	MAX_8_BIT  = 255
	MIN_16_BIT = -0x8000
	MID_16_BIT = 0
	MAX_16_BIT = 0x7fff
)

type waveheader struct {
	// wave structure from https://ccrma.stanford.edu/courses/422/projects/WaveFormat/
	// byte arrays for strings, uints for numbers 
	chunkID       [4]byte // BigEndian
	chunkSize     uint32  // LittleEndian
	format        [4]byte // BigEndian
	subchunk1ID   [4]byte // BigEndian
	subchunk1Size uint32  // LittleEndian
	audioFormat   uint16  // LittleEndian
	numChannels   uint16  // LittleEndian
	sampleRate    uint32  // LittleEndian
	byteRate      uint32  // LittleEndian
	blockAlign    uint16  // LittleEndian
	bitsPerSample uint16  // LittleEndian
	subchunk2ID   [4]byte // BigEndian
	subchunk2Size uint32  // LittleEndian
}

type Wavefile struct {
	header waveheader
	Data   []byte
}

var BadFile = errors.New("File is corrupt or not the proper format")

func (w *Wavefile) Load(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &w.header.chunkID); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.header.chunkSize); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.BigEndian, &w.header.format); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.BigEndian, &w.header.subchunk1ID); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.header.subchunk1Size); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.header.audioFormat); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.header.numChannels); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.header.sampleRate); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.header.byteRate); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.header.blockAlign); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.header.bitsPerSample); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.BigEndian, &w.header.subchunk2ID); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.header.subchunk2Size); err != nil {
		return BadFile
	}
	w.Data = make([]byte, w.header.subchunk2Size)
	binary.Read(r, binary.LittleEndian, &w.Data)

	return validate(w.header)
}
func validate(h waveheader) error {
	if h.chunkID != [4]byte{'R', 'I', 'F', 'F'} {
		return BadFile
	}
	if h.format != [4]byte{'W', 'A', 'V', 'E'} {
		return BadFile
	}
	if h.subchunk1ID != [4]byte{'f', 'm', 't', ' '} {
		return BadFile
	}
	if h.subchunk2ID != [4]byte{'d', 'a', 't', 'a'} {
		return BadFile
	}
	return nil
}

func (w *Wavefile) Save(writer io.Writer) error {
	if err := binary.Write(writer, binary.BigEndian, w.header.chunkID); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.header.chunkSize); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.BigEndian, w.header.format); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.BigEndian, w.header.subchunk1ID); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.header.subchunk1Size); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.header.audioFormat); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.header.numChannels); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.header.sampleRate); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.header.byteRate); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.header.blockAlign); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.header.bitsPerSample); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.BigEndian, w.header.subchunk2ID); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.header.subchunk2Size); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Data); err != nil {
		return err
	}
	return nil
}

func (w *Wavefile) Init() {
	w.header.chunkID = [4]byte{'R', 'I', 'F', 'F'}
	w.header.chunkSize = 36
	w.header.format = [4]byte{'W', 'A', 'V', 'E'}
	w.header.subchunk1ID = [4]byte{'f', 'm', 't', ' '}
	w.header.subchunk1Size = 16
	w.header.audioFormat = 1
	w.header.numChannels = 1
	w.header.sampleRate = 44100
	w.header.byteRate = 44100 * 2
	w.header.blockAlign = 2
	w.header.bitsPerSample = 16
	w.header.subchunk2ID = [4]byte{'d', 'a', 't', 'a'}
	w.header.subchunk2Size = 0
}
func (w *Wavefile) GetData() []byte {
	return w.Data
}
func (w *Wavefile) SetData(b []byte) {
	w.Data = b
	w.header.subchunk2Size = uint32(len(b))
}
func BytesToSigned16(high, low byte) (out int16) {
	if high == 128 && low == 0 {
		return MIN_16_BIT
	}
	highi16 := int16(high & 127)
	highshifted := highi16 << 8
	out = highshifted + int16(low)

	if neg := high & 128; neg != 0 {
		out *= -1
	}

	return out
}
func Signed16ToBytes(in int16) (high, low byte) {
	if in < 0 {
		in = (^in) + 1
		low = byte(in)
		high = byte(in >> 8)
		high |= 128
	} else {
		low = byte(in)
		high = byte(in >> 8)
	}
	return high, low
}
