package main

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
	"math"
	"math/cmplx"
	"os"
	"strconv"
)

var config struct {
	Header        wav.Header
	MAG_THRES     float64
	BIN_PER_FRAME int
	STEP_SIZE     float64
}

const (
	SAMPLE_PER_FRAME = 22050
	BIT_REPEAT       = 5
	PI               = math.Pi
)

func main() {
	var filename = string(os.Args[1]) + ".wav"
	config.MAG_THRES, _ = strconv.ParseFloat(os.Args[2], 64)
	config.BIN_PER_FRAME, _ = strconv.Atoi(os.Args[3])
	temp, _ := strconv.Atoi(os.Args[4])
	config.STEP_SIZE = float64(temp)

	var l []float64
	l = Read(filename)

	var str string
	var biterr int
	var charerr int
	// fmt.Print("\n......\n")
	str, biterr, charerr = Decode(l)
	fmt.Print("\n")
	Bit2Char(str)
	fmt.Print("\n")
	fmt.Print(biterr, " ", charerr)
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
	step := [5]float64{PI / float64(8+config.STEP_SIZE), PI / float64(6+config.STEP_SIZE), PI / float64(4+config.STEP_SIZE), PI / float64(2+config.STEP_SIZE), PI / (config.STEP_SIZE)}
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

func Decode(l []float64) (string, int, int) {
	// mag := make([]float64, 0)
	// phs := make([]float64, 0)

	var i = SAMPLE_PER_FRAME - 1
	var j = 0
	var pos = 0
	var str = ""
	var astr = "0100111001100111011101010111100101100101011011100010000001010100011100100110111101101110011001110010000001010100011010010110111000100000001011010010000001000001010100000100001101010011001100010011001000100000001011010010000001001000010000110100110101010101010100110010000000101101001000000100011101110010011000010110010001110101011000010111010001101001011011110110111000100000010101000110100001100101011100110110100101110011"
	var watermark = 424
	var biterr = 0
	var charerr = 0
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
			if submag[k] < config.MAG_THRES || k == 0 {
				continue
			}
			if count >= config.BIN_PER_FRAME {
				break
			}
			if pos < watermark && count < config.BIN_PER_FRAME {
				var bit = QIMDecode(submag[k], subphs[k])
				count++
				if bit == 1 {
					if astr[pos] != '1' {
						fmt.Println(k, " ", submag[k])
						biterr++
					}
					countone++
				} else {
					if astr[pos] != '0' {
						fmt.Println(k, " ", submag[k])
						biterr++
					}
					countzero++
				}
			}
			if countzero+countone == BIT_REPEAT {
				if countzero > countone {
					if astr[pos] != '0' {
						charerr++
					}
					str += "0"
				} else {
					if astr[pos] != '1' {
						charerr++
					}
					str += "1"
				}
				countzero = 0
				countone = 0
				pos++
			}
			if pos >= watermark {
				// break Loop
				// fmt.Println(Bit2Char(str))
				str = ""
				pos = 0
			}
		}
		if pos >= watermark {
			// break Loop
			// fmt.Println(Bit2Char(str))
			str = ""
			pos = 0
		}
		// mag = append(mag, submag...)
		// phs = append(phs, subphs...)
		j = i + 1
		i += SAMPLE_PER_FRAME
		if len(l)-i > 0 && len(l)-i < SAMPLE_PER_FRAME {
			i = len(l) - 1
		}
	}
	// fmt.Println(Bit2Char(str))
	return str, biterr, charerr
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
