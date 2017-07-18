// durial

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
)

type funcStat struct {
	size      int     // size of the function
	covered   float64 // percentage covered
	remaining float64 // lines remaining to be tested
}

// a map to hold the func path and its stats
type statMap map[string]funcStat

// a map of to be used for function names to line counts
type flmap map[string]int

// usage simply prints out the usage for the program
func usage() {
	fmt.Printf("USAGE: durial <coverage.out>\n")
}

// countFunctionLines counts all the lines of each function in the passed
// in file and returns a map of function names to line count
func countFunctionLines(filePath string, funcLineMap statMap) error {

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

	curStat := funcStat{}

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
						curStat = funcLineMap[filePath+":"+curFuncName]
						curStat.size = funcSize
						uncovered := float64(1 - (curStat.covered / 100.0))
						curStat.remaining = math.Ceil(float64(curStat.size) * uncovered)
						if curStat.covered == 100.0 {
							curStat.remaining = 0
						}
						funcLineMap[filePath+":"+curFuncName] = curStat
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

					curStat = funcLineMap[filePath+":"+curFuncName]
					curStat.size = funcSize
					curStat.remaining = math.Floor(float64(curStat.size)*1.0 - (curStat.covered / 100.0))
					if curStat.covered == 100.0 {
						curStat.remaining = 0
					}
					funcLineMap[filePath+":"+curFuncName] = curStat
				}
			}
		}
	}

	return nil
}

// used to get just the filepath
func extractFilePath(curLine string, uniqueFilePaths flmap) {
	index := strings.Index(curLine, ":")
	filePathStr := curLine[0:index]
	joinPath := os.Getenv("GOPATH") + "/src"
	newPath := path.Join(joinPath, filePathStr)
	uniqueFilePaths[newPath] = 1
}

// populateFuncStats gets the ball rolling
func populateFuncStats(curLine string, funcStats statMap) {
	// get filepath first
	index := strings.Index(curLine, ":")
	filePathStr := curLine[0:index]

	// trim off the :XX: part
	tempStr := curLine[index+1:]
	index = strings.Index(tempStr, ":")
	tempStr = tempStr[index+1:]

	// now get the function name
	index = strings.Index(tempStr, "%")
	funcName := strings.TrimSpace(tempStr[:index-5])

	// now grab the percentage
	pIndex := strings.Index(curLine, "%")
	percentage := strings.TrimSpace(curLine[pIndex-5 : pIndex])

	joinPath := os.Getenv("GOPATH") + "/src"
	fullPath := path.Join(joinPath, filePathStr)

	pVal, err := strconv.ParseFloat(percentage, 64)
	if err != nil {
		log.Printf("ERROR - Can't convert percentage to float")
		pVal = -1.0
	}

	key := fullPath + ":" + funcName

	funcStats[key] = funcStat{
		covered: pVal,
	}

}

// this pulls out the list of files as well as the function coverage percentages
func parseFunctionList(outFile string) (flmap, statMap) {

	file, err := os.Open(outFile)
	if err != nil {
		fmt.Printf("ERROR! Can't open file %v: %v", outFile, err)
		log.Fatal(err)
		return flmap{}, statMap{}
	}

	defer file.Close()

	uniqueFilePaths := make(flmap)
	funcStats := make(statMap)

	scanner := bufio.NewScanner(file)
	curLine := ""

	for scanner.Scan() {
		curLine = scanner.Text()

		if curLine[0:6] == "total:" {
			continue
		}

		// pull out only the list of files we need later
		extractFilePath(curLine, uniqueFilePaths)

		// start populating the function map with the coverage values
		populateFuncStats(curLine, funcStats)
	}

	return uniqueFilePaths, funcStats
}

// main
func main() {

	// check that the file passed in is a .out file
	// then make sure that only one made it in

	// grab the file passed in
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Printf("WARNING! Incorrect usage!\n")
		usage()
		os.Exit(1)
	}

	// try to open it (should be of type .out)
	filePath := args[0]

	fileList, funcStats := parseFunctionList(filePath)

	for key, _ := range fileList {
		err := countFunctionLines(key, funcStats)

		if err != nil {
			log.Printf("ERROR - Can't get line count for file %v: %v", key, err)
		}

	}

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "NAME\tSIZE\tCOVERAGE\tREMAINING\t")
	for key, val := range funcStats {
		//fmt.Fprintln(w, "\t\t\t\t", key, val.size, val.covered, val.remaining)
		fmt.Fprintln(w, key, "\t", val.size, "\t", val.covered, "\t", val.remaining, "\t")
		//fmt.Printf("%v\t\t\t%v lines\t%v%%\t(%v lines uncovered)\n", key, val.size, val.covered, val.remaining)
	}
	w.Flush()
}
