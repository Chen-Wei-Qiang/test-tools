//修改替换json文件
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const LINE_RDWR_SIZE = 4096 * 1024

type translateLine struct {
	tag   string
	value string
}

var originalEntry string
var newEntry string
var fileAbsolutePath string
var filename string

func init() {
	flag.StringVar(&originalEntry, "originalEntry", "", "原词条")
	flag.StringVar(&newEntry, "newEntry", "", "新词条")
	flag.StringVar(&fileAbsolutePath, "fileAbsolutePath", "", "词条文件目录（绝对路径）")
	flag.StringVar(&filename, "filename", "", "文件名称")

	flag.Parse()
}

//读取文件词条
func readFile(fileName string) []translateLine {
	res := make([]translateLine, 0)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Open file Failed", err)
		return res
	}
	defer file.Close()

	br := bufio.NewReaderSize(file, LINE_RDWR_SIZE)
	lineNum := 0
	for {
		byteLine, _, c := br.ReadLine()
		lineNum++
		if c == io.EOF {
			break
		}
		strLine := string(byteLine)
		if len(strLine) <= 0 {
			continue
		}
		//fmt.Println(lineNum, len(strLine), strLine)
		if strLine[0] == '{' || strLine[0] == '}' {
			continue
		}
		if strLine[len(strLine)-1] == ',' {
			strLine = strLine[:len(strLine)-1]
		}
		splitIdx := strings.Index(strLine, ": ")
		if splitIdx == -1 {
			fmt.Printf("%s split error", strLine)
			continue
		}
		tag := strings.TrimSpace(strLine[:splitIdx])
		val := strings.TrimSpace(strLine[splitIdx+2 : len(strLine)-1])

		// "tag"
		tag = tag[1 : len(tag)-1]
		// "val"
		val = val[1:len(val)]

		res = append(res, translateLine{tag: tag, value: val})
	}
	return res
}

//存入文件词条
func saveAsNewFile(fileName string, data []translateLine) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("create file failed", err)
		return
	}
	defer file.Close()
	bw := bufio.NewWriterSize(file, LINE_RDWR_SIZE)
	bw.Write([]byte("{\n"))
	lastIndex := len(data) - 1
	for i, lineData := range data {
		suffix := "\",\n"
		if i == lastIndex {
			suffix = "\"\n"
		}
		wl := "  \"" + lineData.tag + "\": \"" + lineData.value + suffix
		bw.WriteString(wl)
	}
	bw.WriteString("}")
	bw.Flush()
}

func main() {

	/**
	originalEntry  原词条
	newEntry 新词条
	fileAbsolutePath 词条文件目录
	filename 修改后生成的文件名称
	*/

	originalEntry := originalEntry
	newEntry := newEntry
	fileAbsolutePath := fileAbsolutePath
	filename := filename
	a1 := readFile(fileAbsolutePath)

	res := make([]translateLine, 0)

	count := 0
	key := ""
	for _, n := range a1 {
		if n.value == originalEntry {
			count++
			key = n.tag
		}
	}
	if count == 1 {
		for _, n := range a1 {
			if n.tag == key {
				n.value = newEntry
			}
			res = append(res, translateLine{tag: n.tag, value: n.value})
		}
		fmt.Println("value值唯一，词条替换完成")
	} else {
		fmt.Println("value值不唯一，无法替换，请手动排查冲突value")
	}
	saveAsNewFile(filename, res)
	fmt.Println(res)
}
