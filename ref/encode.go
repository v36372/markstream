package main

import (
	"fmt"
	wavwriter "github.com/cryptix/wav"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
	// "log"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"time"
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
	var watermark = string(os.Args[2])
	config.MAG_THRES, _ = strconv.ParseFloat(os.Args[3], 64)
	config.BIN_PER_FRAME, _ = strconv.Atoi(os.Args[4])
	temp, _ := strconv.Atoi(os.Args[5])
	config.STEP_SIZE = float64(temp)
	var outputfile = string(os.Args[1]) + "_" + os.Args[5] + "_wm.wav"

	var l []float64
	l = Read(filename)

	var mag []float64
	var phs []float64

	// log.Println(watermark)
	start := time.Now()
	mag, phs = Embedding(l, watermark)
	elapsed := time.Since(start)
	fmt.Printf("%s ", elapsed)
	var newWav []float64
	// fmt.Println(currentpos)
	start = time.Now()
	newWav = Reconstruct(mag, phs)
	// newWav = Reconstruct(mag, phs)
	elapsed = time.Since(start)
	fmt.Printf("  %s ", elapsed)
	start = time.Now()
	Write(newWav, outputfile)
	elapsed = time.Since(start)
	fmt.Printf("  %s ", elapsed)
}

func Read(filename string) []float64 {
	file, _ := os.Open(filename)
	reader, _ := wav.New(file)

	l, _ := reader.ReadFloatsScale(reader.Samples)
	// fmt.Println(len(l))
	config.Header = reader.Header
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
	// log.Println(stringbit)
	return stringbit
}

func Embedding(l []float64, watermark string) ([]float64, []float64) {
	mag := make([]float64, 0)
	phs := make([]float64, 0)

	var i = SAMPLE_PER_FRAME - 1
	var j = 0
	// fmt.Println(watermark)
	var stringbit = PrepareString(watermark)
	var bitrepeat = 0
	var pos = 0
	for i < len(l) {
		submag := make([]float64, i+1-j)
		subphs := make([]float64, i+1-j)
		var subl = l[j : i+1]
		// fmt.Println(subl)
		subfourier := fft.FFTReal64(subl)
		var count = 0
		// var bitrepeat = 0

		for k, x := range subfourier {
			submag[k], subphs[k] = cmplx.Polar(x)
			if submag[k] < config.MAG_THRES || k == 0 {
				continue
			}
			if count >= config.BIN_PER_FRAME {
				break
			}
			if pos < len(stringbit) && count < config.BIN_PER_FRAME {
				subphs[k] = QIMEncode(submag[k], subphs[k], int(stringbit[pos]))
				count++
				bitrepeat++
			}
			if bitrepeat == BIT_REPEAT {
				bitrepeat = 0
				pos++
			}
			if pos >= len(stringbit) {
				// break
				pos = 0
			}
		}
		// fmt.Println(bitrepeat)
		mag = append(mag, submag...)
		phs = append(phs, subphs...)
		if pos >= len(stringbit) {
			// break
			pos = 0
		}
		j = i + 1
		i += SAMPLE_PER_FRAME
		// fmt.Println(i)
		if len(l)-i > 0 && len(l)-i < SAMPLE_PER_FRAME {
			i = len(l) - 1
		}
	}

	return mag, phs
}

func Reconstruct(mag []float64, phs []float64) []float64 {
	cmplxArray := make([]complex128, len(mag))
	for i, _ := range mag {
		cmplxArray[i] = cmplx.Rect(mag[i], phs[i])
	}

	var i = SAMPLE_PER_FRAME - 1
	var j = 0
	var newWav = make([]float64, 0)
	for i < len(mag) {
		var subcmplx = cmplxArray[j : i+1]
		subIFFT := fft.IFFTRealOutput(subcmplx)
		newWav = append(newWav, subIFFT...)
		j = i + 1
		i += SAMPLE_PER_FRAME
		if len(mag)-i >= 0 && len(mag)-i < SAMPLE_PER_FRAME {
			i = len(mag) - 1
		}
	}

	// newWav = append(newWav, original...)
	return newWav
}

func Write(newWav []float64, outputfile string) {
	wavOut, _ := os.Create(outputfile)
	defer wavOut.Close()

	meta := wavwriter.File{
		Channels:        1,
		SampleRate:      config.Header.SampleRate,
		SignificantBits: config.Header.BitsPerSample,
	}

	writer, _ := meta.NewWriter(wavOut)
	defer writer.Close()

	for n := 0; n < len(newWav); n += 1 {
		integer := int16(newWav[n] * math.MaxInt16)
		writer.WriteInt16(integer)
	}
}

func QIMEncode(mag float64, phs float64, bit int) float64 {
	step := [5]float64{PI / float64(8+config.STEP_SIZE), PI / float64(6+config.STEP_SIZE), PI / float64(4+config.STEP_SIZE), PI / float64(2+config.STEP_SIZE), PI / (config.STEP_SIZE)}
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
