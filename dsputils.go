package markstream

import (
	"encoding/binary"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"time"
)

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 18, 64)
}

func Float64bytes(float float64) []byte {
	bits := math.Float32bits(float32(float))
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func Int16bytes(integer int16) []byte {
	bytes := make([]byte, 2)
	// big edian
	bytes[1] = uint8(integer >> 8)
	bytes[0] = uint8(integer & 0xff)
	return bytes
}

func Int16ArrayByte(f []int16) []byte {
	bytes := make([]byte, 0)
	for _, val := range f {
		bytes = append(bytes, Int16bytes(val)...)
	}
	return bytes
}

func FloatArrayByte(f []float64) []byte {
	bytes := make([]byte, 0)
	for _, val := range f {
		bytes = append(bytes, Float64bytes(val)...)
	}
	return bytes
}

func (ms *MarkStream) Read(filename string) []float64 {
	file, _ := os.Open(filename)
	reader, _ := wav.New(file)

	l, _ := reader.ReadFloatsScale(reader.Samples)

	// config.Header = reader.Header
	return l
}

func PrepareString(info string) string {
	var stringbit = ""
	byteArray := []byte(info)
	for _, char := range byteArray {
		n := int64(char)
		substr := strconv.FormatInt(n, 2)
		if len(substr) < 8 {
			length := len(substr)
			for j := 1; j <= 8-length; j++ {
				substr = "0" + substr
			}
		}
		stringbit += substr
	}

	return stringbit
}

func (ms *MarkStream) Embedding(l []float64) {
	var i = SAMPLE_PER_FRAME - 1
	var j = 0

	for i < len(l) {
		select {
		case watermark := <-ms.userInputChan:
			var pos = 0
			submag := make([]float64, i+1-j)
			subphs := make([]float64, i+1-j)
			var stringbit = PrepareString("1" + watermark + "\n")
			for pos < len(stringbit) {
				var subl = l[j : i+1]
				subfourier := fft.FFTReal64(subl)
				var count = 0
				var bitrepeat = 0

				for k, x := range subfourier {
					submag[k], subphs[k] = cmplx.Polar(x)
					if submag[k] < MAG_THRES {
						continue
					}
					if pos < len(stringbit) && count < BIN_PER_FRAME {
						subphs[k] = QIMEncode(submag[k], subphs[k], int(stringbit[pos]))
						count++
						bitrepeat++
					}
					if bitrepeat == BIT_REPEAT {
						bitrepeat = 0
						pos++
					}
				}

				cmplxArray := make([]complex128, len(subl))
				for i, _ := range cmplxArray {
					cmplxArray[i] = cmplx.Rect(submag[i], subphs[i])
				}
				newWav := fft.IFFTRealOutput(cmplxArray)
				Wav16bit := Scale(newWav)
				ms.ConnManager.audioDataChan <- Wav16bit
				j = i + 1
				i += SAMPLE_PER_FRAME
				if len(l)-i > 0 && len(l)-i < SAMPLE_PER_FRAME {
					i = len(l) - 1
				}
				time.Sleep(400 * time.Millisecond)
			}
		default:
			var pos = 0
			var subl = l[j : i+1]
			subfourier := fft.FFTReal64(subl)
			submag := make([]float64, i+1-j)
			subphs := make([]float64, i+1-j)
			var stringbit = PrepareString("0")
			for pos < len(stringbit) {
				var count = 0
				var bitrepeat = 0

				for k, x := range subfourier {
					submag[k], subphs[k] = cmplx.Polar(x)
					if submag[k] < MAG_THRES {
						continue
					}
					if pos < len(stringbit) && count < BIN_PER_FRAME {
						subphs[k] = QIMEncode(submag[k], subphs[k], int(stringbit[pos]))
						count++
						bitrepeat++
					}
					if bitrepeat == BIT_REPEAT {
						bitrepeat = 0
						pos++
					}
				}

			}
			cmplxArray := make([]complex128, len(subl))
			for i, _ := range cmplxArray {
				cmplxArray[i] = cmplx.Rect(submag[i], subphs[i])
			}
			newWav := fft.IFFTRealOutput(cmplxArray)
			Wav16bit := Scale(newWav)
			ms.ConnManager.audioDataChan <- Wav16bit
			j = i + 1
			i += SAMPLE_PER_FRAME
			if len(l)-i > 0 && len(l)-i < SAMPLE_PER_FRAME {
				i = len(l) - 1
			}
			time.Sleep(400 * time.Millisecond)
		}
	}
}

func Scale(wav []float64) []int16 {
	var newWav = make([]int16, len(wav))
	for i, x := range wav {
		integer := int16(x * math.MaxInt16)
		newWav[i] = integer
	}

	return newWav
}

func QIMEncode(mag float64, phs float64, bit int) float64 {
	step := [5]float64{PI / 18, PI / 14, PI / 10, PI / 6, PI / 2}
	var stepsize = findStep(mag)
	if bit == 48 {
		return math.Floor(phs/step[stepsize]+0.5) * step[stepsize]
	} else {
		return math.Floor(phs/step[stepsize])*step[stepsize] + step[stepsize]/2
	}
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
