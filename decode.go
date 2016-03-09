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

const (
  freqThes = 0.0001
  offset = 1
  spf = 80
)

func main(){
  file , _ := os.Open("test_wm.wav")
  reader , _ := wav.New(file)
  l,err := reader.ReadFloatsScale(reader.Samples-8)
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
  // for i :=0 ;i<10;i++ {
  //   fmt.Print(l[i], " ")
  // }
  // fmt.Println(reader.Samples)
  // for i:=0;i<reader.Samples-8;i++ {
  //   fmt.Print(l[i], " ")
  //   if i%10 == 0 {
  //     fmt.Println("")
  //   }
  // }

  mag := make([]float64, 0)
  phs := make([]float64, 0)

  // fmt.Println(mag[0])
  // fourier := fft.FFTReal32(l)
  // var max float64
  //
  // for i,a:= range fourier{
  //     mag[i], phs[i] = cmplx.Polar(a)
  //     if mag[i] > max {
  //       max = mag[i]
  //     }
  // }

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
    // fmt.Println(max)
    if len(l)-i>=0&&len(l)-i<spf{
      i=len(l)-1
    }
  }
  fmt.Println(phs[1]," ",phs[2]," ",phs[3]," ",phs[4])
  var str = ""
  var pi = math.Pi
  step := [5]float64{pi/10,pi/8,pi/6,pi/4,pi/2}
  fmt.Println(max)
  var k=offset
  var countzero=0
  var countone=1
  var res=0
  for res<16*10 {
    if math.Abs(mag[k]) < freqThes{
      fmt.Println("oopps")
      k++
      continue
    }
    var stepsize = findStep(mag[k],max)
    fmt.Print(phs[k]," ",stepsize," ")
    // integer := int64(math.Floor(phs[k]/(step[stepsize]/2)))
    // r := phs[k]/(step[stepsize]/2) - math.Floor(phs[k]/(step[stepsize]/2))
    var q1 = step[stepsize]*math.Floor(phs[k]/step[stepsize]+0.5)
	  var q2 = step[stepsize]*math.Floor(phs[k]/step[stepsize]) + step[stepsize]/2
    fmt.Println(q1," ",q2)
    if q1<q2{
      countzero++
    }else{
      countone++
    }
    // if r < 0.25||r>0.75 {
    //   countzero++
    // }
    // if r >= 0.25 && r<=0.75 {
    //   countone++
    // }
    if countzero+countone==10{
      if countzero>countone{
        str+= "0"
      }else{
        str+="1"
      }
      countzero=0
      countone=0
      res++
    }
    k++
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
