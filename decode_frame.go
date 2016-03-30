package main

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
	"math"
	"math/cmplx"
	"os"
)

const (
	MAG_THRES        = 0.0001
	SAMPLE_PER_FRAME = 22050
	BIN_PER_FRAME    = 800
	BIT_REPEAT       = 5
	PI               = math.Pi
)

func main() {
	var filename = string(os.Args[1]) + ".wav"

	var l []float64
	l = Read(filename)

	var str string
	var biterr int
	str, biterr = Decode(l)

	Bit2Char(str)
	fmt.Print(biterr)
}

func Bit2Char(str string) string {
	var msg string
	var sum byte
	msg = ""
	var last = 0
	for i, _ := range str {
		sum <<= 1
		sum += str[i] - '0'
		if (i-last+1)%8 == 0 {
			msg += string(sum)
			sum = 0
			last = i + 1
		}
	}

	return msg
}

func Read(filename string) []float64 {
	file, _ := os.Open(filename)
	reader, _ := wav.New(file)

	l, _ := reader.ReadFloatsScale(reader.Samples - 8)

	return l
}

func QIMDecode(mag float64, phs float64) int {
	step := [5]float64{PI / 18, PI / 14, PI / 10, PI / 6, PI / 2}
	var stepsize = findStep(mag)
	integer := int64(math.Floor(phs / (step[stepsize] / 2)))
	r := phs/(step[stepsize]/2) - math.Floor(phs/(step[stepsize]/2))
	if r < 0.5 {
		if integer%2 == 0 {
			return 0
		} else {
			return 1
		}
	}
	if r >= 0.5 {
		if integer%2 == 0 {
			return 1
		} else {
			return 0
		}
	}
	return 0
}

func Decode(l []float64) (string, int) {
	mag := make([]float64, 0)
	phs := make([]float64, 0)

	var i = SAMPLE_PER_FRAME - 1
	var j = 0
	var pos = 0
	var str = ""
	var watermark = 500 * 8
	var biterr = 0
Loop:
	for i < len(l) {
		submag := make([]float64, i+1-j)
		subphs := make([]float64, i+1-j)
		var subl = l[j : i+1]
		subfourier := fft.FFTReal64(subl)
		var countone = 0
		var countzero = 0
		var count = 0
		for k, x := range subfourier {
			submag[k], subphs[k] = cmplx.Polar(x)
			if submag[k] < MAG_THRES || k == 0 {
				continue
			}
			if pos < watermark && count < BIN_PER_FRAME {
				var bit = QIMDecode(submag[k], subphs[k])
				count++
				if bit == 1 {
					countone++
				} else {
					countzero++
				}
			} else if pos >= watermark {
				break Loop
			}
			if countzero+countone == BIT_REPEAT {
				if countzero > countone {
					if astr[pos] != '0' {
						biterr++
					}
					str += "0"
				} else {
					if astr[pos] != '1' {
						biterr++
					}
					str += "1"
				}
				countzero = 0
				countone = 0
				pos++
			}
		}
		mag = append(mag, submag...)
		phs = append(phs, subphs...)
		j = i + 1
		i += SAMPLE_PER_FRAME
		if len(l)-i > 0 && len(l)-i < SAMPLE_PER_FRAME {
			i = len(l) - 1
		}
	}
	return str, biterr
}

func findStep(mag float64) int32 {
	var sMag = mag / (0.005)
	var group = math.Ceil(sMag / 0.2)
	if group == 0 {
		group = 0
	}
	if group > 4 {
		group = 4
	}
	return int32(group)
}
