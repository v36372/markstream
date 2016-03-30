package main

import (
	"fmt"
	"os/exec"
	// "strings"
	// "sync"
	"bufio"
	"os"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	var (
		cmdOut []byte
		err    error
	)
	var str = "qSYiTF0HVcsmX7DejGD97PAP1p67jvQ854Z3NcJSmEMmASIk7gvSWYKtc3uGmUa8PZLLlSP0b4fsoaVOECIp6wJG3jbVE03zWka92lqkemNzcNLs88gfvePcktPOu2He3khAEVNO9aUPyYqFMCqjI9uz6i90WFPmZsTM5X2Q6HDhFONRcnIU7mPe4y2iGQwPBc44ADBJ3ikSVRcVLAHwXImAQXFtXa5JS226glGjY1T2nkhXIUfXR7RsinobNlU7T5nm1tOZV8M21W8msNs05zMael6tZls6oPSiGXEzce9ANAmZrybmWtFr59eyJ9YBacoTY0Q3B0ajNV1vyhOrRY96bRp92ga1y1nKoe0Z6RGl0BS61JsMwzExf2Xn79nYHWNT9ykrEVHBhpalP5k9hEFhgwnCrM9IKGD9BfgVrpP6597G5K4UW90Y0KWXp6LnpSp8yGpp2AUm76toeJVMF8p8GTRxuNJ1yJEo5m6VoHKHY7SW82Kw"
	var basedir = "RWC_60s/"
	lines, _ := readLines("testfiles.txt")
	for i, x := range lines {
		fmt.Print(i, " ")
		cmdName := "./encode_frame"

		filename := basedir + x
		cmdArgs := []string{filename, str}
		if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
			// fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
			// os.Exit(1)
			fmt.Println(err)
		}

		cmdName = "./decode_frame"
		cmdArgs = []string{filename + "_wm"}
		if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
			// fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
			// os.Exit(1)
			fmt.Println(err)
		}
		sha := string(cmdOut)
		fmt.Println(sha)
		// cmdName = "wine"
		// cmdArgs = []string{"ResampAudio.exe", "-s", "48000", filename + ".wav", filename + "_sample.wav"}
		// if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		// 	// fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
		// 	// os.Exit(1)
		// 	fmt.Println(err)
		// }

		// cmdName = "wine"
		// cmdArgs = []string{"ResampAudio.exe", "-s", "48000", filename + "_wm.wav", filename + "_wm_sample.wav"}
		// if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		// 	// fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
		// 	// os.Exit(1)
		// 	fmt.Println(err)
		// }
		// cmdName = "wine"
		// cmdArgs = []string{"PQevalAudio.exe", filename + "_sample.wav", filename + "_wm_sample.wav"}
		// if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		// 	// fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
		// 	// os.Exit(1)
		// 	fmt.Println(err)
		// }
		// var score = string(cmdOut)
		// fmt.Print(" ", score[len(score)-8:])
	}
}
