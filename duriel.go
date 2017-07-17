// durial

package duriel

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
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

// a map of to be used for function names to line counts
type flmap map[string]int

// usage simply prints out the usage for the program
func usage() {
	fmt.Printf("USAGE: durial <coverage.out>\n")
}

// countFunctionLines counts all the lines of each function in the passed
// in file and returns a map of function names to line count
func countFunctionLines(filePath string, funcLineMap flmap) (flmap, error) {

	file, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("ERROR! Can't open file %v: %v", filePath, err)
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
	//funcLineMap := make(map[string]int)

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

	return funcLineMap, nil
}

func parseFilePaths(outFile string) flmap {

	file, err := os.Open(outFile)
	if err != nil {
		fmt.Printf("ERROR! Can't open file %v: %v", outFile, err)
		log.Fatal(err)
		return flmap{}
	}

	defer file.Close()

	uniqueFilePaths := make(map[string]int)
	scanner := bufio.NewScanner(file)
	filePathStr := ""
	curLine := ""
	newPath := ""

	for scanner.Scan() {
		curLine = scanner.Text()

		if curLine[0:6] == "total:" {
			continue
		}

		index := strings.Index(curLine, ":")
		filePathStr = curLine[0:index]
		newPath = path.Join(os.Getenv("GOPATH"), filePathStr)
		uniqueFilePaths[newPath] = 1

	}

	return uniqueFilePaths
}

// main
func main() {

	// check that the file passed in is a .out file
	// then make sure that only one made it in
	// then scan through that file and pickout the filenames
	// iterate over that list of filenames and call the counting function
	// once all the scanning is done (consider multithreading) then work the math
	// need to just do an estimate on number of remaining uncovered lines

	// grab the file passed in
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Printf("WARNING! Incorrect usage!\n")
		usage()
		os.Exit(1)
	}

	// try to open it (should be of type .out)
	filePath := args[0]

	// the map that will hold all the files
	newMap := make(flmap)

	//filePathList := parseFilePaths(curDir, filePath)
	filePathList := parseFilePaths(filePath)

	for key, val := range filePathList {
		log.Printf("%v - %v", key, val)
	}

	debug := false
	if debug {
		finalMap, err := countFunctionLines(filePath, newMap)

		if err != nil {
			log.Printf("ERROR - Can't get line count for file %v: %v", filePath, err)
		}

		testFP := "/usr/local/go/src/bufio/bufio.go"
		finalMap, err = countFunctionLines(testFP, finalMap)

		if err != nil {
			log.Printf("ERROR - Can't get line count for file %v: %v", filePath, err)
		}

		for key, val := range finalMap {
			fmt.Printf("[%v - %v]\n", key, val)
		}
	}
}
