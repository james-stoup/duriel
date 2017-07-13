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

// map of values would look like this:
// path:ling:function -> {coverage, line_count, lines_untested}
// example: fmt/print.go:1133:doPrintln -> {8, 100%, 0}
//

// should pass in the file with the list of functions and their percentages, and the location of the file/files to read
// start with one file though

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
	filePath := args[0]
	file, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("ERROR! Can't open file %v: %v", args[1], err)
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	curLine := ""
	inFunc := false
	funcSize := 0
	index1 := 0
	curFuncName := ""

	// map to hold func names and line counts
	funcLineMap := make(map[string]int)

	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		curLine = scanner.Text()

		if inFunc {
			strings.TrimSpace(curLine)

			if len(curLine) > 0 {

				// If we are in a function, first check to see if this is the last line of the function
				if curLine[0:1] == "}" {
					inFunc = false
				} else {
					// Don't count comments
					if curLine[0:2] != "//" {
						funcSize = funcSize + 1
						funcLineMap[filePath+"/"+curFuncName] = funcSize

					}
				}
			}
		} else {
			if len(curLine) > 4 {
				if curLine[0:4] == "func" {
					inFunc = true
					funcSize = 0
					curFuncName = ""
					funcNameStr := ""

					// need to check if this is an interface method
					index1 = strings.Index(curLine, "(")
					if index1 == 5 {
						tempLine := curLine[6:]
						index2 := strings.Index(tempLine, ")")
						tempLine2 := tempLine[index2+2:]
						index3 := strings.Index(tempLine2, "(")
						funcNameStr = tempLine2[:index3]
						curFuncName = funcNameStr
					} else {

						curFuncName = curLine[5:index1]
					}

					funcLineMap[filePath+"/"+curFuncName] = 0
				}
			}
		}
	}

	for key, val := range funcLineMap {
		fmt.Printf("%v \t %v\n", key, val)
	}
}
