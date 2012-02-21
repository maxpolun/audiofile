package audiofile

import (
	"bytes"
	"testing"
)

func Test_WavefileShouldBeAbleToBeUsedAsAnAudiofile(t *testing.T) {
	// This test is mostly a static assertion
	var af AudioFile = &Wavefile{}
	switch af.(type) {
	default:
		t.Errorf("Wavefile should be able to be converted back and forth with AudioFile inteface")
	case *Wavefile:
	}
}

func Test_EmptyReaderShouldReturnAnError(t *testing.T) {
	empty := bytes.NewBuffer([]byte{})
	af := &Wavefile{}
	err := af.Load(empty)
	if err == nil {
		t.Errorf("expected to get an error when loading an empty file, got nil instead")
	}

}

func Test_ValidWaveHeaderShouldNotReturnAnError(t *testing.T) {
	af := &Wavefile{}
	err := af.Load(validWave)
	if err != nil {
		t.Errorf("expected no errors with valid wave header, got %v", err)
	}
}
func badId() *bytes.Buffer {
	newbuf := make([]byte, 44)
	copy(newbuf, validWave.Bytes())
	newbuf[0] = 0
	return bytes.NewBuffer(newbuf)
}
func badFormat() *bytes.Buffer {
	newbuf := make([]byte, 44)
	copy(newbuf, validWave.Bytes())
	newbuf[8] = 0
	return bytes.NewBuffer(newbuf)
}
func badSubchunkId() *bytes.Buffer {
	newbuf := make([]byte, 44)
	copy(newbuf, validWave.Bytes())
	newbuf[12] = 0
	return bytes.NewBuffer(newbuf)
}
func badSubchunkId2() *bytes.Buffer {
	newbuf := make([]byte, 44)
	copy(newbuf, validWave.Bytes())
	newbuf[36] = 0
	return bytes.NewBuffer(newbuf)
}
func Test_ShouldCheckForFileCorruption(t *testing.T) {
	buffuncs := []func() *bytes.Buffer{
		badId,
		badFormat,
		badSubchunkId,
		badSubchunkId2}
	msgs := []string{
		"Riff header ID",
		"Wave header ID",
		"Wave subchunk 1 ID",
		"Wave data chunk header ID"}

	for i := range buffuncs {
		buf := buffuncs[i]()
		af := Wavefile{}
		err := af.Load(buf)
		if err == nil {
			t.Errorf("Expected an error on corrupted %v, got no error", msgs[i])
		}
	}
}
func Test_InitShouldProduceAValidWave(t *testing.T) {
	af := &Wavefile{}
	af.Init()
	if err := validate(af.header); err != nil {
		t.Errorf("Expected Init to produce a valid wavefile, got error %v", err)
		t.Errorf("header has chunkID = %v, wave headerID = %v, wave subchunk1ID = %v, Wave  data chunk header = %v", af.header.chunkID, af.header.format, af.header.subchunk1ID, af.header.subchunk2ID)
	}
}
