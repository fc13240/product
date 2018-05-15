package text

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadLine(filename string, line func(line string)) {
	file, err := os.Open(filename)

	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	var emi = 0
	buff := bufio.NewReader(file) // 读入缓冲区
	for {

		str, err := buff.ReadString('\n')
		if err != nil {
			if !(err.Error() == "EOF") {
				fmt.Println(err)
				os.Exit(1)
				return
			}
		} else {
			str = strings.Trim(str, "\n")
			str = strings.Trim(str, "\r")
			str = strings.Trim(str, " ")
			line(str)
		}

		if str == "" {
			emi++
			if emi > 10 {
				return
			}
		}
	}
}

func Log(filename, content string) {
	if f, err := os.OpenFile(filename, os.O_WRONLY, 0664); err == nil {
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
}
