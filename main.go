// sort | (c) 2020 NETWAYS GmbH | GPLv2+

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
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
	reversePrintingParameter := flag.Bool("r", false, "Print result in reverse order")
	quickSortParameter := flag.Bool("qsort", false, "Use quick sort")
	mergeSortParameter := flag.Bool("mergesort", false, "Use merge sort")
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

	if len(data) <= 0 {
		return
	}

	if *ignoreCaseParameter {
		data = bytes.ToUpper(data)
	}

	dataSplit := bytes.Split(data, []byte{'\n'})

	if len(dataSplit[len(dataSplit)-1]) <= 0 {
		dataSplit = dataSplit[:len(dataSplit)-1]
	}

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

	if *reversePrintingParameter {
		comparator = wrapReverse(comparator)
	}

	if *quickSortParameter {
		quickSort(dataSplit, comparator)
	} else if *mergeSortParameter {
		dataSplit = mergeSort(dataSplit, comparator)
	} else {
		bubbleSort(dataSplit, comparator)
	}

	dataJoin := bytes.Join(dataSplit, []byte{'\n'})

	if *writingParameter != "" {
		writeErr := ioutil.WriteFile(*writingParameter, append(dataJoin,'\n'), 0644)
		if writeErr != nil {
			fmt.Fprintln(os.Stderr, writeErr)
			os.Exit(1)
		}
	} else {
		fmt.Printf("%s\n", string(dataJoin))
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

func wrapReverse(f func(a, b []byte) int) func(a, b []byte) int {
	return func(a, b []byte) int {
		return -f(a, b)
	}
}

func quickSort(input [][]byte, comp func(a, b []byte) int) [][]byte {
	if len(input) < 2 {
		return input
	}

	left, right := 0, len(input)-1
	pivot := rand.Int() % len(input)
	input[pivot], input[right] = input[right], input[pivot]

	for i := range input {
		if comp(input[i], input[right]) < 0 {
			input[left], input[i] = input[i], input[left]
			left++
		}
	}

	input[left], input[right] = input[right], input[left]

	quickSort(input[:left], comp)
	quickSort(input[left+1:], comp)

	return input
}

func bubbleSort(input [][]byte, comp func(a, b []byte) int) [][]byte {
	for i := 0; i < len(input); i++ {
		j := i
		for j > 0 && comp(input[j], input[j-1]) < 0 {
			input[j], input[j-1] = input[j-1], input[j]
			j -= 1
		}
	}

	return input
}

func mergeSort(input [][]byte, comp func(a, b []byte) int) [][]byte {
	if len(input) < 2 {
		return input
	}

	middle := len(input) / 2

	return merge(mergeSort(input[:middle], comp), mergeSort(input[middle:], comp), comp)
}

func merge(left, right [][]byte, comp func(a, b []byte) int) (result [][]byte) {
	result = make([][]byte, len(left)+len(right))

	i := 0
	for len(left) > 0 && len(right) > 0 {
		if comp(left[0], right[0]) < 0 {
			result[i] = left[0]
			left = left[1:]
		} else {
			result[i] = right[0]
			right = right[1:]
		}
		i++
	}

	for j := 0; j < len(left); j++ {
		result[i] = left[j]
		i++
	}
	for j := 0; j < len(right); j++ {
		result[i] = right[j]
		i++
	}

	return
}
