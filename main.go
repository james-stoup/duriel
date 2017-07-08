// durial

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// setup...formatter, importer, autocomplete

// TODO
// 1. run test and generate *.out file
// 2. run cover tool to generate func stats
// 3. read in go files and start parsing functions
//   3a. read function name, store in map like so package:funcName
//   3b. skip blank lines and comments
//   3c. count total function lines
//   3d. store in map
// 4. when all files have been read, iterate through list of functions from go cover
// 5. calculate number of lines left in untested state and save in map
// 6. iterate through final map, pretty printing the output

// usage simply prints out the usage for the program
func usage() {
	fmt.Printf("USAGE: durial <coverage.out>\n")
}

// main
func main() {

	// grab the file passed in
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Printf("WARNING! Incorrect usage!\n")
		usage()
		os.Exit(1)
	}

	// try to open it (should be of type .out)
	file, err := os.Open(args[0])

	if err != nil {
		fmt.Printf("ERROR! Can't open file %v: %v", args[1], err)
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	curLine := ""
	inFunc := false
	funcSize := 0
	index := 0

	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		curLine = scanner.Text()

		if inFunc {
			strings.TrimSpace(curLine)

			if len(curLine) > 0 {

				// If we are in a function, first check to see if this is the last line of the function
				if curLine[0:1] == "}" {
					inFunc = false
					fmt.Printf("[%v lines]\n", funcSize)
				} else {
					// Don't count comments
					if curLine[0:2] != "//" {
						funcSize++
					}
				}
			}
		} else {
			if len(curLine) > 4 {
				if curLine[0:4] == "func" {
					inFunc = true
					index = strings.Index(curLine, "{")
					fmt.Printf("> %v ", curLine[:index])
				}
			}
		}
	}

}
