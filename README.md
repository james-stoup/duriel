# duriel
A tool for interpreting Golang's test coverage output

## What does this do and why do I care?
Go includes helpful tools for viewing the test coverage of your test. Using 'go test' and 'go tool cover' you can generate output that looks like this:

```
$ go test  -coverprofile=scanner_coverage.out
$ go tool cover -func=scanner_coverage.out
github.com/rogpeppe/godef/go/scanner/errors.go:38:	Reset			0.0%
github.com/rogpeppe/godef/go/scanner/errors.go:41:	ErrorCount		100.0%
github.com/rogpeppe/godef/go/scanner/errors.go:52:	Error			0.0%
github.com/rogpeppe/godef/go/scanner/errors.go:65:	Len			100.0%
github.com/rogpeppe/godef/go/scanner/errors.go:66:	Swap			100.0%
...
```

And while this is very useful, when dealing with lots of files (and hundreds of functions) it can be difficult to tell at a glance what functions need to have their test coverage increased. Since the cover tool only reports percentages, it is up to the user to figure out which functions need more coverage. When dealing with large projects this can get a bit difficult because it is possible to have a 200 line function at 90% coverage and a 10 line function at 60% coverage. In this case, the larger function has more uncovered lines even though at first glance the second function would seem to be a greater priority. Duriel solves this be giving you a more comprehensive rundown of your testing status like so:

```
$ duriel cover_output.txt
                                                                   FILE|                 FUNCTION|    SIZE|   COVERAGE|   REMAINING|
    /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/errors.go |                   Error |      1 |      100% |          0 |
   /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/scanner.go |              scanNumber |     61 |      100% |          0 |
   /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/scanner.go |              scanEscape |     39 |     92.6% |          3 |
   /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/scanner.go |                   error |      4 |      100% |          0 |
   /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/scanner.go |             findLineEnd |     38 |     95.8% |          2 |
   /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/scanner.go |                    Scan |    151 |      100% |          0 |
   /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/scanner.go |           scanRawString |     11 |      100% |          0 |
    /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/errors.go |                   Reset |      1 |        0% |          0 |
   /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/scanner.go |                digitVal |      9 |      100% |          0 |
   /home/breezy/go/src/github.com/rogpeppe/godef/go/scanner/scanner.go |              scanString |     14 |       90% |          2 |
```

## Installation
Simply pull down this repo, cd into it and run the install command to put the binary in the bin directory in your GOPATH.

```
$ cd ~/go/src/github.com/james-stoup/duriel
$ go install
```

## How to use Duriel
Using this tool is a three step process. First, you have to run the test code and generate a coverage profile of all the functions. Second, you have to use the cover tool to convert that output to something Durel can use. Third, run Duriel to see your nicely formatted results. Here are the steps in more detail. For additional resources as well as a very good write up of Go's testing functionality, take a look at [the cover story](https://blog.golang.org/cover) at the Go Blog.

### Generating the coverage profile
To generate the coverage profile run your test like you normally would, just append the coverage flag and specify the output file.

```
$ go test -coverprofile=my_lib_coverage.out
```

This .out file can be used by other Go tools for other types of output, but we don't care about that in this case.

### Converting the coverage file
Go's cover tool can take the coverage file and generate the helpful coverage percentages that we saw earlier. Duriel takes this as input when it does calculations.

```
$ go tool cover -func=my_lib_coverage.out > my_lib_coverage.txt
```

### Run Duriel
Now we can pass this newly converted file to Duriel and watch it do its magic.

```
$ duriel my_lib_coverage.txt
```

