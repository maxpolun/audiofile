/*
	Package audiofile is a pure Go generic audiofile library, it reads from a 
	standard io.Reader and writes to an io.Writer

	currently it supports wav files, with other lossless formats planned
*/
package audiofile

import (
	"encoding/binary"
	"errors"
	//	"fmt"
	"io"
)

// The AudioReader interface is a file that can be read -- it can't 
// necessarily be written out in the case of a format like mp3 that
// requires complex encoding
type AudioReader interface {
	// Load the audio file from an io.Reader It performs some basic 
	// sanity checks and returns an error if there is either an io error
	// or if the file appears to be corrupted
	Load(io.Reader) error
	// GetBytes gets the raw bytes of audio data from the audio file. It performs 
	// no conversions and does not return any of the metadata
	GetBytes() []byte
}

// The AudioWriter interface is an audiofle that can be written out
type AudioWriter interface {
	// Writes a valid audio file out to the given io.Writer. It 
	// should only return an error if there is an io error
	Save(io.Writer) error
	// Init sets up the audio file metadata for the 
	// audiofile with the default values
	Init()
	// SetBytes sets the audio data to the given raw bytes. Like 
	// GetBytes, it performs no conversions.
	SetBytes([]byte)
}

// The AudioFile interface is a file that can be read or written. Most of the 
// supported file types should implement this interface
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

// wave structure from https://ccrma.stanford.edu/courses/422/projects/WaveFormat/
// byte arrays for strings, uints for numbers 
type Waveheader struct {
	ChunkID       [4]byte // BigEndian
	ChunkSize     uint32  // LittleEndian
	Format        [4]byte // BigEndian
	Subchunk1ID   [4]byte // BigEndian
	Subchunk1Size uint32  // LittleEndian
	AudioFormat   uint16  // LittleEndian
	NumChannels   uint16  // LittleEndian
	SampleRate    uint32  // LittleEndian
	ByteRate      uint32  // LittleEndian
	BlockAlign    uint16  // LittleEndian
	BitsPerSample uint16  // LittleEndian
	Subchunk2ID   [4]byte // BigEndian
	Subchunk2Size uint32  // LittleEndian
}

// Wavefile format -- *Wavefile implements AudioFile
type Wavefile struct {
	Header Waveheader
	Data   []byte
}

// BadFile indicates that the audiofile is corrupt
var BadFile = errors.New("File is corrupt or not the proper format")

func (w *Wavefile) Load(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &w.Header.ChunkID); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.Header.ChunkSize); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.BigEndian, &w.Header.Format); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.BigEndian, &w.Header.Subchunk1ID); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.Header.Subchunk1Size); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.Header.AudioFormat); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.Header.NumChannels); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.Header.SampleRate); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.Header.ByteRate); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.Header.BlockAlign); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.Header.BitsPerSample); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.BigEndian, &w.Header.Subchunk2ID); err != nil {
		return BadFile
	}
	if err := binary.Read(r, binary.LittleEndian, &w.Header.Subchunk2Size); err != nil {
		return BadFile
	}
	w.Data = make([]byte, w.Header.Subchunk2Size)
	binary.Read(r, binary.LittleEndian, &w.Data)

	return validate(w.Header)
}
func validate(h Waveheader) error {
	if h.ChunkID != [4]byte{'R', 'I', 'F', 'F'} {
		return BadFile
	}
	if h.Format != [4]byte{'W', 'A', 'V', 'E'} {
		return BadFile
	}
	if h.Subchunk1ID != [4]byte{'f', 'm', 't', ' '} {
		return BadFile
	}
	if h.Subchunk2ID != [4]byte{'d', 'a', 't', 'a'} {
		return BadFile
	}
	return nil
}

func (w *Wavefile) Save(writer io.Writer) error {
	if err := binary.Write(writer, binary.BigEndian, w.Header.ChunkID); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Header.ChunkSize); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.BigEndian, w.Header.Format); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.BigEndian, w.Header.Subchunk1ID); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Header.Subchunk1Size); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Header.AudioFormat); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Header.NumChannels); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Header.SampleRate); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Header.ByteRate); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Header.BlockAlign); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Header.BitsPerSample); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.BigEndian, w.Header.Subchunk2ID); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Header.Subchunk2Size); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, w.Data); err != nil {
		return err
	}
	return nil
}

func (w *Wavefile) Init() {
	w.Header.ChunkID = [4]byte{'R', 'I', 'F', 'F'}
	w.Header.ChunkSize = 36
	w.Header.Format = [4]byte{'W', 'A', 'V', 'E'}
	w.Header.Subchunk1ID = [4]byte{'f', 'm', 't', ' '}
	w.Header.Subchunk1Size = 16
	w.Header.AudioFormat = 1
	w.Header.NumChannels = 1
	w.Header.SampleRate = 44100
	w.Header.ByteRate = 44100 * 2
	w.Header.BlockAlign = 2
	w.Header.BitsPerSample = 16
	w.Header.Subchunk2ID = [4]byte{'d', 'a', 't', 'a'}
	w.Header.Subchunk2Size = 0
}
func (w *Wavefile) GetBytes() []byte {
	return w.Data
}
func (w *Wavefile) SetBytes(b []byte) {
	w.Data = b
	w.Header.Subchunk2Size = uint32(len(b))
}

// GetPCM is a utility function -- it calls GetBytes and interperets the data
// as 16bit signed LPCM data
func GetPCM(areader AudioReader) []int16 {
	bytes := areader.GetBytes()
	out := make([]int16, len(bytes)/2)
	for i := range out {
		out[i] = BytesToSigned16(bytes[i*2], bytes[i*2+1])
	}
	return out
}

// SetPCM is a utility function -- it saves the given 16bit signed LPCM data
// in the given AudioWriter
func SetPCM(awriter AudioWriter, pcm []int16) {
	bytes := make([]byte, len(pcm)*2)
	for i := range pcm {
		low, high := Signed16ToBytes(pcm[i])
		bytes[i*2] = low
		bytes[i*2+1] = high
	}
	awriter.SetBytes(bytes)
}

// BytesToSigned16 is a bit by bit copy of 2 bytes into a signed 16bit value
func BytesToSigned16(low, high byte) (out int16) {
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

// Signed16ToBytes is a bit by bit copy of a signed 16bit int into 2 bytes
func Signed16ToBytes(in int16) (low, high byte) {
	if in == -1 {
		return 0, 128
	}
	if in < 0 {
		in = (^in) + 1
		low = byte(in)
		high = byte(in >> 8)
		high |= 128
	} else {
		low = byte(in)
		high = byte(in >> 8)
	}
	return
}
