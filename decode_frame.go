package main
import (
    "fmt"
    "github.com/mjibson/go-dsp/fft"
    "github.com/mjibson/go-dsp/wav"
    "os"
    "math"
    "math/cmplx"
)

const (
  FREQ_THRES = 0.0001
  BIT_OFFSET = 1
  SAMPLE_PER_FRAME=3000
  BIT_REPEAT=5
  PI=math.Pi
)

func main(){
  file , _ := os.Open("test_wm_frame.wav")
  reader , _ := wav.New(file)
  l,err := reader.ReadFloatsScale(reader.Samples-8)
  if err!= nil {
    fmt.Println(err)
  }

  fmt.Println(l[1]," ",l[2]," ",l[3]," ",l[4]," ",l[5]," ",l[6])
  //----------------------fft the whole file-------------------
  // mag := make([]float64, reader.Samples)
  // phs := make([]float64, reader.Samples)
  //
  // fmt.Println(mag[0])
  // fourier := fft.FFTReal(l)
  //
  // for i,a:= range fourier{
  //     mag[i], phs[i] = cmplx.Polar(a)
  // }
  //----------------------fft the whole file-------------------

  //---------------------divide into frames---------------------
  mag := make([]float64, 0)
  phs := make([]float64, 0)

  var max float64
  max = 0
  var i = SAMPLE_PER_FRAME-1
  var j = 0
  for i<len(l) {
    max = 0
    submag := make([]float64, i+1-j)
    subphs := make([]float64, i+1-j)
    var subl = l[j:i+1]
    subfourier :=  fft.FFTReal32(subl)
    for k,x :=range subfourier {
      submag[k],subphs[k] = cmplx.Polar(x)
      if submag[k] > max {
        max = submag[k]
      }
    }
    mag = append(mag, submag...)
    phs = append(phs, subphs...)
    j=i+1
    i+=SAMPLE_PER_FRAME
    if len(l)-i>=0&&len(l)-i<SAMPLE_PER_FRAME{
      i=len(l)-1
    }
  }

  fmt.Println(phs[BIT_OFFSET], " ",phs[BIT_OFFSET+1]," ", phs[BIT_OFFSET+2]," ",phs[BIT_OFFSET+3])
  //---------------------divide into frames---------------------

  var str = ""
  step := [5]float64{PI/20,PI/16,PI/12,PI/8,PI/4}
  var k=BIT_OFFSET
  var countzero=0
  var countone=1
  var res=0
  for res<16*8{
    if math.Abs(mag[k]) < FREQ_THRES{
      k++
      continue
    }
    var stepsize = findStep(mag[k])
    integer := int64(math.Floor(phs[k]/(step[stepsize]/2)))
    r := phs[k]/(step[stepsize]/2) - math.Floor(phs[k]/(step[stepsize]/2))
    // fmt.Println(phs[k]," ",r)
    if r < 0.5 {
      if integer % 2 == 0 {
        countzero++
      } else {
        countone++
      }
    }
    if r >= 0.5 {
      if integer % 2 == 0 {
        countone++
      } else {
        countzero++
      }
    }
    if countzero+countone==BIT_REPEAT{
      fmt.Println(countzero, " ", countone, " ", res)
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
  fmt.Println(len(str))
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

func findStep(mag float64) int32{
  var sMag = mag/(0.005)
  var group = math.Ceil(sMag/0.2)
  if group==0{
    group=0
  }
  if group>4{
    group=4
  }
  return int32(group)
}