package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func exe_cmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println(cmd)
	parts := strings.Fields(cmd)
	out, err := exec.Command(parts[0], parts[1]).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done()
}

func main() {
	wg := new(sync.WaitGroup)
	commands := []string{"echo newline >> foo.o", "echo newline >> f1.o", "echo newline >> f2.o"}
	for _, str := range commands {
		wg.Add(1)
		go exe_cmd(str, wg)
	}
	wg.Wait()
}
