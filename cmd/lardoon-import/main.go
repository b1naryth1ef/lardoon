package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/b1naryth1ef/lardoon"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <target-directory>\n", os.Args[0])
		return
	}

	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".acmi") {
			err := lardoon.ImportFile(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error: %v", err)
	}

	// err := lardoon.ImportFile(os.Args[1])
	// if err != nil {
	// 	fmt.Printf("error: %v", err)
	// }
}
