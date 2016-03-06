package main
import (
    // "flag"
    "fmt"
    // "github.com/youpy/go-wav"
    "github.com/mjibson/go-dsp/fft"
    "github.com/mjibson/go-dsp/wav"
    // "io"
    "os"
    "math"
    "math/cmplx"
    // "strconv"
)


func main(){
  file , _ := os.Open("test_wm.wav")
  reader , _ := wav.New(file)
  l,err := reader.ReadFloats(reader.Samples-8)
  if err!= nil {
    fmt.Println(err)
  }

  // file, _ := os.Open("cosine_wm.wav")
  // reader := wav.NewReader(file)
  // defer file.Close()
  // // fmt.Println(reader.WavData.Size)
  // // FrameSyns(reader,l)
  // var l []float64
  // for {
  //     samples, err := reader.ReadSamples()
  //     // l := make([]float64, reader.WavData.Size)
  //     if err == io.EOF {
  //         break
  //     }
  //     l = make([]float64, reader.WavData.Size)
  //     // r := make([]int, reader.WavData.Size)
  //     for i, sample := range samples {
  //         // fmt.Printf("L/R: %d/%d\n", reader.IntValue(sample, 0), reader.IntValue(sample, 1))
  //         l[i] = float64(reader.IntValue(sample,0))
  //         // fmt.Println(l[i])
  //         // r[i] = reader.IntValue(sample,1)
  //     }
  // }

  fmt.Println(reader.Samples)
  // for i:=0;i<reader.Samples-8;i++ {
  //   fmt.Print(l.([]int16)[i], " ")
  //   if i%10 == 0 {
  //     fmt.Println("")
  //   }
  // }

  mag := make([]float64, reader.Samples)
  phs := make([]float64, reader.Samples)

  fmt.Println(mag[0])
  fourier := fft.FFTReal32(l)

  for i,a:= range fourier{
      mag[i], phs[i] = cmplx.Polar(a)
  }

  var str = ""
  var pi = math.Pi
  var step =  pi/6;

  for i:=10;i<34;i++ {
    integer := int64(math.Floor(phs[i]/(step/2)))
    r := phs[i]/(step/2) - math.Floor(phs[i]/(step/2))
    if r < 0.5 {
      if integer % 2 == 0 {
        str += "0"
      } else {
        str += "1"
      }
    }
    if r >= 0.5 {
      if integer % 2 == 0 {
        str += "1"
      } else {
        str += "0"
      }
    }
  }
  fmt.Println(str)
  var sum byte
  var last = 0
  for i,_ := range str {
    sum <<= 1;
    sum += str[i] - '0';
    if (i-last+1)%8==0{
      fmt.Print(string(sum))
      sum = 0
      last=i+1
    }
  }
}
