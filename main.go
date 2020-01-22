package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func main() {
	numericParameter := flag.Bool("n", false, "Sort numbers from lowest to highest")
	writingParameter := flag.String("o", "", "Save result into a file")
	uniqueParameter := flag.Bool("u", false, "Output only the first of an equal run")
	leadingBlanksParameter := flag.Bool("b", false, "Ignore leading blanks")
	ignoreCaseParameter := flag.Bool("f", false, "Fold lower case to upper case characters")
	randomSortingParameter := flag.Bool("R", false, "shuffle, but group identical keys.")
	flag.Parse()

	var comparator func(a, b []byte) int

	var data []byte

	if len(flag.Args()) < 1 {
		content, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		data = content
	} else {
		fileContent, readFileErr := ioutil.ReadFile(flag.Args()[0])
		if readFileErr != nil {
			fmt.Fprintln(os.Stderr, readFileErr)
			os.Exit(1)
		}
		data = fileContent
	}

	if string(data) == "" {
		return
	}

	if *ignoreCaseParameter {
		data = []byte(strings.ToUpper(string(data)))
	}

	dataSplit := bytes.Split(data, []byte("\n"))

	if *leadingBlanksParameter {
		for v := range dataSplit {
			dataSplit[v] = bytes.TrimLeftFunc(dataSplit[v], unicode.IsSpace)
		}
	}

	if *uniqueParameter {
		duplicates := make(map[string]struct{})
		for _, v := range dataSplit {
			duplicates[string(v)] = struct{}{}
		}
		dataSplit = nil
		for b := range duplicates {
			dataSplit = append(dataSplit, []byte(b))
		}
	}

	if *randomSortingParameter {
		comparator = randomComparator
		rand.Seed(time.Now().UnixNano())
	} else if *numericParameter {
		comparator = numberComparator
	} else {
		comparator = bytes.Compare
	}

	for i := 0; i < len(dataSplit); i++ {
		j := i
		for j > 0 && comparator(dataSplit[j], dataSplit[j-1]) < 0 {
			dataSplit[j], dataSplit[j-1] = dataSplit[j-1], dataSplit[j]
			j -= 1
		}
	}
	dataJoin := bytes.Join(dataSplit, []byte("\n"))

	if *writingParameter != "" {
		writeErr := ioutil.WriteFile(*writingParameter, []byte(strings.TrimPrefix(string(dataJoin), "\n")+"\n"), 0644)
		if writeErr != nil {
			fmt.Fprintln(os.Stderr, writeErr)
			os.Exit(1)
		}
	} else {
		fmt.Printf("%s\n", strings.TrimPrefix(string(dataJoin), "\n"))
	}
}

func numberComparator(a, b []byte) int {
	comparer1, convertErr1 := strconv.Atoi(string(a))
	if convertErr1 != nil {
		fmt.Fprintln(os.Stderr, convertErr1)
		os.Exit(1)
	}
	comparer2, convertErr2 := strconv.Atoi(string(b))
	if convertErr2 != nil {
		fmt.Fprintln(os.Stderr, convertErr2)
		os.Exit(1)
	}
	return comparer1 - comparer2
}

func randomComparator(a, b []byte) int {
	return rand.Intn(3) - 1
}
