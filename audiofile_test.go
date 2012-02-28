package audiofile

import (
	"bytes"
	"testing"
)

var validWaveBuf []byte = []byte{
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
	0, 0, 0, 0} //44 Subchunk 2 Size

var validWave *bytes.Buffer = bytes.NewBuffer(validWaveBuf)

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
	copy(newbuf, validWaveBuf)
	newbuf[0] = 0
	return bytes.NewBuffer(newbuf)
}
func badFormat() *bytes.Buffer {
	newbuf := make([]byte, 44)
	copy(newbuf, validWaveBuf)
	newbuf[8] = 0
	return bytes.NewBuffer(newbuf)
}
func badSubchunkId() *bytes.Buffer {
	newbuf := make([]byte, 44)
	copy(newbuf, validWaveBuf)
	newbuf[12] = 0
	return bytes.NewBuffer(newbuf)
}
func badSubchunkId2() *bytes.Buffer {
	newbuf := make([]byte, 44)
	copy(newbuf, validWaveBuf)
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
func Test_GetDataShouldReturnTheBytesOfDataFromTheFile(t *testing.T) {
	newData := []byte{
		0, 0,
		127, 255, // aka 0111 1111 1111 1111 aka max 16bit integer
		0, 0,
		128, 0, // aka 1000 0000 0000 0000 aka min 16bit integer
		0, 0}

	newBuf := bytes.Join([][]byte{validWaveBuf, newData}, nil)
	newBuf[40] = byte(len(newData))

	af := &Wavefile{}
	err := af.Load(bytes.NewBuffer(newBuf))
	if err != nil {
		t.Errorf("expected no errors with valid wave file with data, got %v on file %v", err, newBuf)
	}
	data := af.GetData()
	if bytes.Compare(newData, data) != 0 {
		t.Errorf("expected data returned from wavefile.GetData() to be the same as the input data.\n\nexpected %v, got %v",
			newData, data)
		t.Errorf("subchunk2Size = %v", af.header.subchunk2Size)
	}
}
func Test_SetDataShouldUpdateTheInternalBuffer(t *testing.T) {
	newData := []byte{
		0, 0,
		127, 255, // aka 0111 1111 1111 1111 aka max 16bit integer
		0, 0,
		128, 0, // aka 1000 0000 0000 0000 aka min 16bit integer
		0, 0}
	af := &Wavefile{}
	af.Init()
	af.SetData(newData)
	if bytes.Compare(af.Data, newData) != 0 {
		t.Errorf("expected to get %v, got %v in SetData()", newData, af.Data)
	}
}
func Test_SetDataShouldUpdateTheSizeCount(t *testing.T) {
	newData := []byte{
		0, 0,
		127, 255, // aka 0111 1111 1111 1111 aka max 16bit integer
		0, 0,
		128, 0, // aka 1000 0000 0000 0000 aka min 16bit integer
		0, 0}
	af := &Wavefile{}
	af.Init()
	af.SetData(newData)
	if af.header.subchunk2Size != uint32(len(newData)) {
		t.Errorf("expected data length to be %v, found %v", len(newData), af.header.subchunk2Size)
	}
}

func Test_BytesToSigned16(t *testing.T) {
	inputs := [][]byte{
		{0, 0},
		{127, 255},
		{128, 0},
		{0, 1},
		{128, 1}}
	expected := []int16{
		0,
		MAX_16_BIT,
		MIN_16_BIT,
		1,
		-1}
	for i := range inputs {
		if result := BytesToSigned16(inputs[i][0], inputs[i][1]); result != expected[i] {
			t.Errorf("expected %v, got %v with inputs %v", expected[i], result, inputs[i])
		}
	}
}
func Test_Signed16ToBytes(t *testing.T) {
	expected := [][]byte{
		{0, 0},
		{127, 255},
		{128, 0},
		{0, 1},
		{128, 1}}
	inputs := []int16{
		0,
		MAX_16_BIT,
		MIN_16_BIT,
		1,
		-1}
	for i := range inputs {
		if result1, result2 := Signed16ToBytes(inputs[i]); result1 != expected[i][0] || result2 != expected[i][1] {
			t.Errorf("expected %v, got [%v, %v] with input %v", expected[i], result1, result2, inputs[i])
		}
	}
}
func Test_SaveShouldNotHaveAnError(t *testing.T) {
	af := &Wavefile{}
	af.Load(bytes.NewBuffer(validWaveBuf))
	outbuf := make([]byte, 1024)
	err := af.Save(bytes.NewBuffer(outbuf))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func Test_SaveShouldGiveTheSameBytesAsInputPreviously(t *testing.T) {
	af := &Wavefile{}
	if err := af.Load(bytes.NewBuffer(validWaveBuf)); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	t.Logf("af = %v", af)
	outbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	err := af.Save(outbuf)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	b := outbuf.Bytes()
	if bytes.Compare(b, validWaveBuf) != 0 {
		t.Errorf("expected %v (len %v), got %v (len %v)", validWaveBuf, len(validWaveBuf), b, len(b))
	}
}
