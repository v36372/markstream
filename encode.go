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
    // "log"
)

const (
	// normScale = float64(math.MaxInt16)
  freqThes = 0.0001
  offset = 1
  spf = 80
)

func main() {
    file , _ := os.Open("test.wav")
    reader , _ := wav.New(file)
    // fmt.Println(reader.Samples)

    // l,err := reader.ReadSamples(reader.Samples)
    l,err := reader.ReadFloatsScale(reader.Samples)
    if err!= nil {
      fmt.Println(err)
    }

    // fmt.Println(max)

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

    // fmt.Println(reader.Samples)
    // for i :=0 ;i<10;i++ {
    //   fmt.Print(l[i], " ")
    // }
    mag := make([]float64, 0)
    phs := make([]float64, 0)
    // mag := make([]float64, 10)
    // phs := make([]float64, 10)

    // fmt.Println(mag[0])
    // fourier := fft.FFTReal32(l)
    // fmt.Println(fourier[1])
    // fmt.Println(fft.IFFT(fourier))
    var max float64
    max = 0
    var i = spf-1
    var j = 0
    for i<len(l) {
      max = 0
      submag := make([]float64, spf)
      subphs := make([]float64, spf)
      // fmt.Println(j, " ", i+1)
      var subl = l[j:i+1]
      // fmt.Println(len(subl))
      subfourier :=  fft.FFTReal32(subl)
      // fmt.Println(subfourier)
      for k,x :=range subfourier {
        submag[k],subphs[k] = cmplx.Polar(x)
        if submag[k] > max {
          max = submag[k]
        }
      }
      mag = append(mag, submag...)
      phs = append(phs, subphs...)
      j=i+1
      i+=spf
      // fmt.Println(i)
      if len(l)-i>0&&len(l)-i<spf{
        i=len(l)-1
      }
    }

    // fmt.Println(max)
    // for i=0;i<200;i++{
    //   fmt.Println(mag[i], " ", phs[i])
    // }
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
    step := [5]float64{pi/10,pi/8,pi/6,pi/4,pi/2}

    var info = "Nguyen Trong Tin"
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
      // fmt.Println(substr)
      stringbit += substr
    }

    fmt.Println(stringbit)
    var k=offset
    var count=0
    i =0
    for i<len(stringbit){
      if math.Abs(mag[k]) < freqThes {
        // fmt.Println("ooops ",mag[k])
        k++
        continue
      }
      fmt.Print(phs[k], " ")
      var stepsize = findStep(mag[k],max)
      // if stringbit[i] == '0'{
      //   phs[k] = math.Floor(phs[k]/step[stepsize] + 0.5)*step[stepsize]
      // }
      // if stringbit[i] == '1' {
      //   phs[k] = math.Floor(phs[k]/step[stepsize])*step[stepsize] + step[stepsize]/2
      // }

      var bit float64
      // bit =0;
      if stringbit[i] == '0'{
        bit=0
      }else {
        bit=1
      }
      phs[k]=step[stepsize]*Round(phs[k]/step[stepsize] - bit/2,0.5,0) + step[stepsize]*bit/2
      fmt.Println(phs[k], " " , stepsize)
      count++
      if count==10{
        count=0
        i++
      }
      k++
    }
    // fmt.Println(k)
    // fmt.Println(phs[1]," ",phs[2]," ",phs[3]," ",phs[4])
    cmplxArray := make([]complex128, reader.Samples)
    // cmplxArray := make([]complex128, 10)
    for i,_ := range mag {
      cmplxArray[i] = cmplx.Rect(mag[i],phs[i])
    }
    // fmt.Println(cmplxArray[1])
    i = spf-1
    j=0
    var newWav = make([]float64, 0)
    for i<len(l) {
      // max = 0
      // submag := make([]float64, spf)
      // subphs := make([]float64, spf)
      // fmt.Println(j, " ", i+1)
      var subcmplx = cmplxArray[j:i+1]
      // fmt.Println(len(subcmplx))
      subIFFT :=  fft.IFFTRealOutput(subcmplx)
      newWav = append(newWav, subIFFT...)

      // fmt.Println(subfourier)
      // mag = append(mag, submag...)
      // phs = append(phs, subphs...)
      j=i+1
      i+=spf
      // fmt.Println(i)
      if len(l)-i>=0&&len(l)-i<spf{
        // fmt.Println("aaaaaaaaaaaa")
        i=len(l)-1
      }
    }
    // var wm_frame = fft.IFFT(cmplxArray)
    // for i :=0 ;i<10;i++ {
    //   fmt.Print(newWav[i], " ")
    // }

    // var newWav = make([]float64, reader.Samples)
    // for i,_ := range wm_frame {
    //   newWav[i] = real(wm_frame[i])
    // }

    // outfile, _ := os.Create("a_wm.wav")
    // var wr := wav.NewWriter(outfile,reader.WavData.Size,)
    // outfile, _ := os.Create("a_wm.wav")
    // reader := wav.NewWriter(file)
    // format, _ := reader.Format()
    // fmt.Println(reader.Header.NumChannels)
    // fmt.Println(reader.Header.BitsPerSample)
    // fmt.Println(reader.Header.SampleRate)
    // fmt.Println(reader.Header.AudioFormat)

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
    // fmt.Println("out put lenght : ", len(newWav))
  	// var freq float64
  	// freq = 0.0001
    // b := make([]byte, 2)
    // fmt.Println(writer.SamplesWritten)
  	for n := 0; n < reader.Samples; n += 1 {
      integer := int16(newWav[n]*math.MaxInt16)
      // toNumber := uint16(newWav[n]   * normScale) // Inverse the read scaling
  		// binary.LittleEndian.PutUint16(b, uint16(toNumber))
  		// writer.WriteSample(b)
      // fmt.Print(integer, " ")
  		err = writer.WriteInt16(integer)
  		checkErr(err)
  	}
    // fmt.Println(writer.SamplesWritten)
  	// fmt.Printf("Simulation Done. Took:%v\n", time.Since(start))


}

func findStep(mag float64,max float64) int32{
  // var k = int32(math.Ceil(math.Abs(max)/mag))
  // if k > 4{
  //   k = 4
  // }
  // // fmt.Println(mag, " ",k)
  // return k
  var sMag = mag/(0.005)
  var group = math.Ceil(sMag/0.2)
  // fmt.Println(group)
  if group==0{
    group=0
  }
  if group>4{
    group=4
  }
  return int32(group)
}

func Round(val float64, roundOn float64, places int ) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
