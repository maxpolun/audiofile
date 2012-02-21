package audiofile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type AudioReader interface {
	Load(io.Reader) error
}
type AudioWriter interface {
	Save(io.Writer) error
	Init() // set up the initial headers, etc.
}
type AudioFile interface {
	AudioReader
	AudioWriter
}

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

func (w *Wavefile) Save(r io.Writer) error {
	return errors.New("Not Implemented")
}

var validWave *bytes.Buffer = bytes.NewBuffer([]byte{
	'R', 'I', 'F', 'F', //4 ChunkID
	36, 0, 0, 0, //8 ChunkSize
	'W', 'A', 'V', 'E', //12 Format
	'f', 'm', 't', ' ', //16 Subchunk1ID
	16, 0, 0, 0, //20 Subchunk1Size
	1, 0, //22 AudioFormat
	1, 0, //24 NumChannels
	0x44, 0xAC, 0, 0, //28 SampleRate 44100 
	0x88, 0x58, 0x01, 0, //32 ByteRate 
	2, 0, // 34 block Align
	16, 0, // 36 BitsPerSample
	'd', 'a', 't', 'a', // 40 Subchunk2ID
	0, 0, 0, 0}) //44 Subchunk 2 Size

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
