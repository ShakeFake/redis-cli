package utils

import (
	"bufio"
	"fmt"
	"os"
)

func ReadRedisKeyFromFile(path string) []string {
	fileH, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bufioReader := bufio.NewReader(fileH)

	var result []string
	for {
		line, _, err := bufioReader.ReadLine()
		if err != nil {
			fmt.Println(err)
			break
		}
		result = append(result, string(line))
	}

	return result
}
