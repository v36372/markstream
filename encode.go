package main

import (
    // "flag"
    "fmt"
    // "encoding/binary"
    // "github.com/youpy/go-wav"
    "github.com/mjibson/go-dsp/wav"
    "github.com/mjibson/go-dsp/fft"
    wavwriter "github.com/cryptix/wav"
    // "io"
    "os"
    "math"
    "math/cmplx"
    "strconv"
)

const (
	normScale = float64(math.MaxInt16)
)

func main() {
    file , _ := os.Open("test.wav")
    reader , _ := wav.New(file)
    fmt.Println(reader.Samples)

    // l,err := reader.ReadSamples(reader.Samples)
    l,err := reader.ReadFloats(reader.Samples)
    if err!= nil {
      fmt.Println(err)
    }

    // fmt.Println(l)

    // file, _ := os.Open("a.wav")
    // reader := wav.NewReader(file)
    // defer file.Close()
    // // fmt.Println(reader.WavData.Size)
    // // FrameSyns(reader,l)
    // var l []float64
    // var r []float64
    // for {
    //     samples, err := reader.ReadSamples()
    //     // l := make([]float64, reader.WavData.Size)
    //     if err == io.EOF {
    //         break
    //     }
    //     l = make([]float64, reader.WavData.Size)
    //     r = make([]float64, reader.WavData.Size)
    //     for i, sample := range samples {
    //         // fmt.Printf("L/R: %d/%d\n", reader.IntValue(sample, 0), reader.IntValue(sample, 1))
    //         l[i] = float64(reader.IntValue(sample,0))
    //         // fmt.Println(l[i])
    //         r[i] = float64(reader.IntValue(sample,1))
    //     }
    // }
    //
    //
    // for i:=0;i<reader.Samples;i++ {
    //   fmt.Print(l.([]int16)[i], " ")
    //   if i%10 == 0 {
    //     fmt.Println("")
    //   }
    // }
    // for i:=0;i<reader.Samples;i++ {
    //   fmt.Print(l[i], " ")
    //   if i%10 == 0 {
    //     fmt.Println("")
    //   }
    // }

    // l := make([]float64,10)
    // l[0] = 1
    // l[1] = 2
    // l[2] = 3
    // l[3] = 4
    // l[4] = 5
    // l[5] = 6
    // l[6] = 7
    // l[7] = 8
    // l[8] = 9
    // l[9] = 10

    //f, _ := os.Create("cosine_test.txt")

    fmt.Println(reader.Samples)
    for i :=0 ;i<10;i++ {
      fmt.Print(l[i], " ")
    }
    mag := make([]float64, reader.Samples)
    phs := make([]float64, reader.Samples)
    // mag := make([]float64, 10)
    // phs := make([]float64, 10)

    // fmt.Println(mag[0])
    fourier := fft.FFTReal32(l)
    fmt.Println(fourier[0])
    // fmt.Println(fft.IFFT(fourier))
    for i,a:= range fourier{
        mag[i], phs[i] = cmplx.Polar(a)
    }
    // fmt.Println(phs[0],phs[1],phs[2],phs[3],phs[4],phs[5],phs[6],phs[7],phs[8])

    // fmt.Println(mag)
    // fmt.Println(phs)
    // var maxM,minM, maxP, minP int
    // maxM = minM = mag[0]
    // maxP = minP = phs[0]
    // for i,a := range mag{
    //     if maxM < mag[i]
    //       maxM = mag[i]
    //     if minM > mag[i]
    //       minM = mag[i]
    //     if maxP < phs[i]
    //       maxP = phs[i]
    //     if minP > phs[i]
    //       minP = phs[i]
    // }
    var pi = math.Pi
    var step =  pi/6;

    var info = "abc"
    var stringbit = ""
    byteArray := []byte(info)
    for _, char := range byteArray{
      n := int64(char)
      substr := strconv.FormatInt(n, 2)// 111001
      if len(substr) < 8{
        length := len(substr)
        for j:=1;j<=8-length;j++{
          substr = "0" + substr
        }
      }
      fmt.Println(substr)
      stringbit += substr
    }
    fmt.Println(stringbit)
    for i, bit := range stringbit{
      if bit == '0'{
        phs[i+10] = math.Floor(phs[i+10]/step + 0.5)*step
      }
      if bit == '1' {
        phs[i+10] = math.Floor(phs[i+10]/step)*step + step/2
      }
    }
    // fmt.Println(phs[0],phs[1],phs[2],phs[3],phs[4],phs[5],phs[6],phs[7],phs[8])
    cmplxArray := make([]complex128, reader.Samples)
    // cmplxArray := make([]complex128, 10)
    for i,_ := range mag {
      cmplxArray[i] = cmplx.Rect(mag[i],phs[i])
    }
    fmt.Println(cmplxArray[0])
    var wm_frame = fft.IFFT(cmplxArray)
    for i :=0 ;i<10;i++ {
      fmt.Print(real(wm_frame[i]), " ")
    }

    var newWav = make([]float64, reader.Samples)
    for i,_ := range wm_frame {
      newWav[i] = real(wm_frame[i])
    }

    // outfile, _ := os.Create("a_wm.wav")
    // var wr := wav.NewWriter(outfile,reader.WavData.Size,)
    // outfile, _ := os.Create("a_wm.wav")
    // reader := wav.NewWriter(file)
    // format, _ := reader.Format()
    fmt.Println(reader.Header.NumChannels)
    fmt.Println(reader.Header.BitsPerSample)
    fmt.Println(reader.Header.SampleRate)
    fmt.Println(reader.Header.AudioFormat)

    wavOut, err := os.Create("test_wm.wav")
  	checkErr(err)
  	defer wavOut.Close()

  	meta := wavwriter.File{
  		Channels:        1,
  		SampleRate:      reader.Header.SampleRate,
  		SignificantBits: reader.Header.BitsPerSample,
  	}

  	writer, err := meta.NewWriter(wavOut)
  	checkErr(err)
  	defer writer.Close()

  	// start := time.Now()

  	// var freq float64
  	// freq = 0.0001
    // b := make([]byte, 2)
    fmt.Println(writer.SamplesWritten)
  	for n := 0; n < reader.Samples; n += 1 {
      integer := int16(newWav[n]*65535 - 32378)
      // toNumber := uint16(newWav[n]   * normScale) // Inverse the read scaling
  		// binary.LittleEndian.PutUint16(b, uint16(toNumber))
  		// writer.WriteSample(b)
  		err = writer.WriteInt16(integer)
  		checkErr(err)
  	}
    fmt.Println(writer.SamplesWritten)
  	// fmt.Printf("Simulation Done. Took:%v\n", time.Since(start))


}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
