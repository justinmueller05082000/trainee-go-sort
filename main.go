package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		return
	}

	if string(data) == "" {
		return
	}

	dataSplit := bytes.Split(data, []byte("\n"))
	lengthOfData := len(dataSplit)

	for i := 0; i < lengthOfData; i++ {
		j := i
		for j > 0 && bytes.Compare(dataSplit[j], dataSplit[j-1]) < 0 {
			dataSplit[j], dataSplit[j-1] = dataSplit[j-1], dataSplit[j]
			j -= 1
		}
	}
	dataJoin := bytes.Join(dataSplit, []byte("\n"))
	fmt.Printf("%s\n", strings.TrimPrefix(string(dataJoin), "\n"))
}
