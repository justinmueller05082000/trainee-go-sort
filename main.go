package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	numericParameter := flag.Bool("n", false, "Sort numbers from lowest to highest")
	writingParameter := flag.String("o", "", "Save result into a file")
	flag.Parse()

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

	dataSplit := bytes.Split(data, []byte("\n"))
	lengthOfData := len(dataSplit)

	for i := 0; i < lengthOfData; i++ {
		j := i
		for j > 0 {
			if !*numericParameter {
				if bytes.Compare(dataSplit[j], dataSplit[j-1]) < 0 {
					dataSplit[j], dataSplit[j-1] = dataSplit[j-1], dataSplit[j]
					j -= 1
				} else {
					break
				}
			} else {
				comparer1, convertErr1 := strconv.Atoi(string(dataSplit[j]))
				if convertErr1 != nil {
					fmt.Fprintln(os.Stderr, convertErr1)
					os.Exit(1)
				}
				comparer2, convertErr2 := strconv.Atoi(string(dataSplit[j-1]))
				if convertErr2 != nil {
					fmt.Fprintln(os.Stderr, convertErr2)
					os.Exit(1)
				}
				if comparer1 < comparer2 {
					dataSplit[j], dataSplit[j-1] = dataSplit[j-1], dataSplit[j]
					j -= 1
				} else {
					break
				}
			}
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
