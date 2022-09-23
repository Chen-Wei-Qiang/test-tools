package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/axgle/mahonia"
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
		// "val",
		val = val[1:len(val)]

		res = append(res, translateLine{tag: tag, value: val})
	}
	return res
}

func saveAsNewFile(langType, fileName string, data []translateLine) {
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
		enc := mahonia.NewDecoder("UTF-8")
		resData := enc.ConvertString(string(wl))
		bw.WriteString(resData)
		//bw.WriteString(wl)
	}
	bw.WriteString("}")
	bw.Flush()
	//fmt.Printf("generated file:%s line:%d\n ", fileName, len(data))
}

var sourceFile string
var diffFile string
var targetFile string
var zhCheckFile string
var langType string
var prefixPath string

var file_one string
var file_two string

func init() {
	flag.StringVar(&file_one, "file_one", "", "对比文件1路径")
	flag.StringVar(&file_two, "file_two", "", "对比文件2路径")
	flag.Parse()
}

func main() {

	file_one := &file_one

	file_two := &file_two
	currentData := readFile(*file_one)
	//fmt.Println("当前所有翻译项：", currentData[0])
	//fmt.Println("当前所有翻译项数量：", len(currentData))
	Data := readFile(*file_two)
	//fmt.Println("当前所有翻译项：", Data[0])
	//fmt.Println("当前所有翻译项数量：", len(Data))

	fmt.Printf("{")

	newData := make([]translateLine, 0)

	var biaoji bool
	for _, n := range currentData {
		biaoji = false
		for _, n1 := range Data {
			if n.tag == n1.tag {
				biaoji = true
				continue
			}
		}
		if biaoji == false {
			newData1 := new(translateLine)
			fmt.Printf("\"%v\"", n.tag)
			fmt.Printf(":\"%v\",", n.value)
			newData1.tag = n.tag
			newData1.value = n.value
			newData = append(newData, *newData1)
		}
	}

	fmt.Printf("}")

	saveAsNewFile("", "diff.json", newData)
}
