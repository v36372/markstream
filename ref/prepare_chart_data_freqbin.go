package main

import (
	"fmt"
	"bufio"
	"os"
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

func main(){
  f, error := os.Create("ref/test_result_summary_biterr_freq.txt")
  if error != nil {
    fmt.Println(error)
  }

  lines, _ := readLines("testfiles.txt")

	var err1,err2,err3,err4,err5 [800]int

	for i:=0;i<800;i++{
		err1[i] = err2[i] = err3[i] = err4[i] = err5[i] = 0
	}

  numErr1 := "";
  numErr2 := "";
  numErr3 := "";
  numErr4 := "";
  numErr5 := "";

	for _, x := range lines {
    lines1, _ := readLines("test_result/" + x + "_1.txt")
    lines2, _ := readLines("test_result/" + x + "_2.txt")
    lines3, _ := readLines("test_result/" + x + "_3.txt")
    lines4, _ := readLines("test_result/" + x + "_4.txt")
    lines5, _ := readLines("test_result/" + x + "_5.txt")


    if strings.Split(lines1[len(lines1)-1]," ")[0] != "0"{
      numErr1 += strings.Split(lines1[len(lines1)-1]," ")[0] + ", "
    }
    if strings.Split(lines2[len(lines2)-1]," ")[0] != "0"{
      numErr2 += strings.Split(lines2[len(lines2)-1]," ")[0] + ", "
    }
    if strings.Split(lines3[len(lines3)-1]," ")[0] != "0"{
      numErr3 += strings.Split(lines3[len(lines3)-1]," ")[0] + ", "
    }
    if strings.Split(lines4[len(lines4)-1]," ")[0] != "0"{
      numErr4 += strings.Split(lines4[len(lines4)-1]," ")[0] + ", "
    }
    if strings.Split(lines5[len(lines5)-1]," ")[0] != "0"{
      numErr5 += strings.Split(lines5[len(lines5)-1]," ")[0] + ", "
    }
  }

  f.WriteString("\n" + numErr1)
  f.WriteString("\n" + numErr2)
  f.WriteString("\n" + numErr3)
  f.WriteString("\n" + numErr4)
  f.WriteString("\n" + numErr5)
}
