// durial

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
)

type funcStat struct {
	name      string  // name of the function
	path      string  // path to the file
	size      int     // size of the function
	covered   float64 // percentage covered
	remaining float64 // lines remaining to be tested
}

// A map to hold the func path and its stats
type statMap map[string]funcStat

// A map used for matching function names to line counts
type flmap map[string]int

// Usage simply prints out the usage for the program
func usage() {
	fmt.Printf("USAGE: durial <coverage.out>\n")
}

// Removes any comments that might be tacked on at the end of the line
func cleanComment(line string) string {
	cppcmt := regexp.MustCompile(`//.*`)
	cleanedBytes := cppcmt.ReplaceAll([]byte(line), []byte(""))
	return string(cleanedBytes[:])
}

// Handle calculating the coverage
func calcStats(funcSize int, curStat *funcStat) {
	uncovered := float64(1 - (curStat.covered / 100.0))
	curStat.remaining = math.Ceil(float64(curStat.size) * uncovered)
	if curStat.covered == 100.0 {
		curStat.remaining = 0
	}
}

// countFunctionLines counts the lines of each function in the passed
// in file and returns a map of function names to line counts. It ignores
// comments and empty lines when counting.
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
	key := ""

	// Iterate over the entire file
	for scanner.Scan() {
		curLine = scanner.Text()

		if inFunc {

			if len(curLine) > 0 {

				// First check to see if this is the last line of the function
				if curLine[0:1] == "}" {
					inFunc = false
				} else {
					// Don't count comments
					if curLine[0:2] != "//" {
						funcSize = funcSize + 1
						curStat = funcLineMap[key]
						curStat.size = funcSize
						calcStats(funcSize, &curStat)
						funcLineMap[key] = curStat
					}
				}
			}

		} else {

			if len(curLine) > 5 && curLine[0:5] == "func " {

				// remove any trailing comments
				curLine = cleanComment(curLine)

				inFunc = true
				funcSize = 0

				index1 = strings.Index(curLine, "(")
				curFuncName = curLine[5:index1]

				// need to check if this is an interface method
				if index1 == 5 {
					tempLine := curLine[6:]
					index2 := strings.Index(tempLine, ")")
					tempLine2 := tempLine[index2+2:]
					index3 := strings.Index(tempLine2, "(")
					curFuncName = tempLine2[:index3]
				}

				// set the value for the map
				key = filePath + ":" + curFuncName
				curStat = funcLineMap[key]
				curStat.name = curFuncName
				curStat.path = filePath

				// check if this is a one line function
				if strings.Contains(curLine, "{") && strings.Contains(curLine, "}") {
					curStat.size = 1
					inFunc = false
				}

				funcLineMap[key] = curStat
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
	fmt.Fprintln(w, "FILE\tFUNCTION\tSIZE\tCOVERAGE\tREMAINING\t")
	for _, val := range funcStats {
		fmt.Fprintln(w, val.path, "\t", val.name, "\t", val.size, "\t", fmt.Sprintf("%v%%", val.covered), "\t", val.remaining, "\t")
	}
	w.Flush()
}
