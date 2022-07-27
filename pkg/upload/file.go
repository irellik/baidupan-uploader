package upload

import (
	"fmt"
	"os"
)

func HandleBigFile(fileName string, size int32, handle func(int, []byte)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	s := make([]byte, size)
	index := 0
	for {
		switch nr, err := f.Read(s); true {
		case nr < 0:
			fmt.Fprintf(os.Stderr, "cat: error reading: %s\n", err.Error())
			os.Exit(1)
		case nr == 0: // EOF
			return nil
		case nr > 0:
			handle(index, s[0:nr])
			index += 1
		}
	}
}
