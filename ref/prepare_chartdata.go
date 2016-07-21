package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	f, error := os.Create("ref/test_result_summary_charerr.txt")
	if error != nil {
		fmt.Println(error)
	}

	lines, _ := readLines("testfiles.txt")

	// numErr1 := ""
	numErr2 := ""
	numErr3 := ""
	numErr4 := ""
	// numErr5 := ""
	testnumber := ""

	total2 := 0
	total3 := 0
	total4 := 0
	for i, x := range lines {
		// lines1, _ := readLines("test_result/" + x + "_1.txt")
		lines2, _ := readLines("test_result/" + x + "_2.txt")
		lines3, _ := readLines("test_result/" + x + "_3.txt")
		lines4, _ := readLines("test_result/" + x + "_4.txt")
		// lines5, _ := readLines("test_result/" + x + "_5.txt")

		// if strings.Split(lines1[len(lines1)-1]," ")[0] != "0"{
		//   numErr1 += strings.Split(lines1[len(lines1)-1]," ")[0] + ", "
		// }
		temp, _ := strconv.Atoi(strings.Split(lines2[len(lines2)-1], " ")[0])
		total2 += temp
		temp, _ = strconv.Atoi(strings.Split(lines3[len(lines3)-1], " ")[0])
		total3 += temp
		temp, _ = strconv.Atoi(strings.Split(lines4[len(lines4)-1], " ")[0])
		total4 += temp

		if strings.Split(lines2[len(lines2)-1], " ")[1] != "0" {
			testnumber += strconv.Itoa(i+1) + ", "
			numErr2 += strings.Split(lines2[len(lines2)-1], " ")[1] + ", "
			// }
			// if strings.Split(lines3[len(lines3)-1], " ")[0] != "0" {
			numErr3 += strings.Split(lines3[len(lines3)-1], " ")[1] + ", "
			// }
			// if strings.Split(lines4[len(lines4)-1], " ")[0] != "0" {
			numErr4 += strings.Split(lines4[len(lines4)-1], " ")[1] + ", "
		}
		// if strings.Split(lines5[len(lines5)-1]," ")[0] != "0"{
		//   numErr5 += strings.Split(lines5[len(lines5)-1]," ")[0] + ", "
		// }
	}

	f.WriteString(testnumber)
	f.WriteString("\n" + numErr2)
	f.WriteString("\n" + numErr3)
	f.WriteString("\n" + numErr4 + "\n")
	f.WriteString(strconv.Itoa(total2) + " " + strconv.Itoa(total3) + " " + strconv.Itoa(total4))
}
