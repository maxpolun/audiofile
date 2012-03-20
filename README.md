audiofile: a simple and well tested interface to audio files for go.

For now it only supports wav files, but other lossless formats should be simple
enough to add. Lossy formats should be able to be read, but writing to them
would be significantly more work. This is why the interface is split -- a format can 
support only reading, only writing, or both.

an example of usage is below:

    package main
    import (
    	"os"
    	"github.com/maxpolun/audiofile"
    )

    func main() {
    	infile := os.Open("src.wav")
    	af := &WaveFile{}
    	af.Load(infile)
    	bytes := af.GetBytes()

    	// process the data in some way

    	outfile := os.Create("dest.wav")
    	af.Save(outfile)
    }

The extensive unit tests should also contain quite a few examples of usage.


Copyright (c) 2012, Max Polun
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.