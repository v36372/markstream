package main
import (
    // "flag"
    "fmt"
    // "github.com/youpy/go-wav"
    "github.com/mjibson/go-dsp/fft"
    // "github.com/mjibson/go-dsp/wav"
    // "io"
    // "os"
    // "math"
    // "math/cmplx"
    // "strconv"
)


func main(){
  var a = fft.FFTReal([]float64 {1, 2, 3})
  fmt.Println(a)
  var b = fft.IFFT(a)
  fmt.Println(b)
}
