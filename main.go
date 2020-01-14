package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var data []byte

	if len(os.Args) < 2 {
		content, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}
		data = content
	} else {
		fileContent, readFileErr := ioutil.ReadFile(os.Args[1])
		if readFileErr != nil {
			fmt.Println(readFileErr)
			return
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
		for j > 0 && bytes.Compare(dataSplit[j], dataSplit[j-1]) < 0 {
			dataSplit[j], dataSplit[j-1] = dataSplit[j-1], dataSplit[j]
			j -= 1
		}
	}
	dataJoin := bytes.Join(dataSplit, []byte("\n"))
	fmt.Printf("%s\n", strings.TrimPrefix(string(dataJoin), "\n"))
}
