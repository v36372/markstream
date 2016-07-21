package main

import (
	"bufio"
	"fmt"
	"os"
	// "strconv"
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
	f, error := os.Create("markstream_charerr.chart")
	if error != nil {
		fmt.Println(error)
	}

	f.WriteString("ChartType = column\n")
	f.WriteString("Title = Markstream Test Result\n")
	f.WriteString("SubTitle = Character Detection Error\n")
	// f.WriteString("ValueSuffix = °C\n")
	f.WriteString("XAxisNumbers = ")

	lines, _ := readLines("ref/test_result_summary_charerr.txt")

	// for i := 1; i <= 11; i++ {
	f.WriteString(lines[0])
	// }

	f.WriteString("\nYAxisText = Character errors\n")

	f.WriteString("Data|Settings #1 = " + lines[1] + "\n")
	f.WriteString("Data|Settings #2 = " + lines[2] + "\n")
	f.WriteString("Data|Settings #3 = " + lines[3] + "\n")
	// f.WriteString("Data|Settings #4 = " + lines[3] + "\n")
	// f.WriteString("Data|Settings #5 = " + lines[4])
}
