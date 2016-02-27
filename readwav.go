package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mjibson/go-dsp/wav"
)

func main() {
	// infilepath := flag.String("infile", "", "test.wav")
	// flag.Parse()

	file, _ := os.Open("test.wav")
	reader, _ := wav.New(file)

	defer file.Close()

	for {
		samples, err := reader.ReadSamples(reader.Samples)
		if err == io.EOF {
			break
		}

		fmt.Printf(samples)
		// for _, sample := range samples {
		// 	fmt.Printf("L/R: %d/%d\n", reader.IntValue(sample, 0), reader.IntValue(sample, 1))
		// }
	}
	// fmt.Println(fft.FFTReal([]float64{1, 2, 3}))
}
