package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/axgle/mahonia"
	"github.com/liangzibo/go-excel/lzbExcel"
	"github.com/mitchellh/mapstructure"
	"log"
	"os"
	"strconv"
	"time"
)

const LINE_RDWR_SIZE = 4096 * 1024

type translateLine struct {
	tag   string
	value string
}

type ExcelTest struct {
	key string `json:"key" name:"key" index:"0"`
	zh  string `json:"zh" name:"zh" index:"1"`
	en  string `json:"en" name:"en" index:"2"`
	ja  string `json:"ja" name:"ja" index:"3"`
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

var basePath string

func init() {
	flag.StringVar(&basePath, "basePath", "", "词条目录--规定目录下的文件名必需是(en.json、ja.json、en.json)")

	flag.Parse()
}

func main() {

	BasePath := &basePath

	xlsx, err := excelize.OpenFile(*BasePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Get all the rows in a sheet.
	rows := xlsx.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//结果在  arr 中
	var arr []ExcelTest
	err = lzbExcel.NewExcelStructDefault().SetPointerStruct(&ExcelTest{}).RowsAllProcess(rows, func(maps map[string]interface{}) error {
		var ptr ExcelTest
		// map 转 结构体
		if err2 := mapstructure.Decode(maps, &ptr); err2 != nil {
			return err2
		}
		arr = append(arr, ptr)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(arr)

	//结果在  arr 中
	var arr2 []ExcelTest
	//StartRow 开始行,索引从 0开始
	//IndexMax  索引最大行,如果 结构体中的 index 大于配置的,那么使用结构体中的
	err = lzbExcel.NewExcelStruct(1, 10).SetPointerStruct(&ExcelTest{}).RowsAllProcess(rows, func(maps map[string]interface{}) error {
		var ptr ExcelTest
		// map 转 结构体
		if err2 := mapstructure.Decode(maps, &ptr); err2 != nil {
			return err2
		}
		arr2 = append(arr2, ptr)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(arr2)

	var arr3 []map[string]interface{}
	excel := lzbExcel.NewExcelStruct(1, 10).SetPointerStruct(&ExcelTest{})
	for i, row := range rows {
		//If the index is less than the set start row, skip
		//如果 索引 小于 已设置的 开始行,那么跳过
		if i < excel.StartRow {
			continue
		}
		//单行处理
		m, err3 := excel.Row(row)
		if err3 != nil {
			fmt.Println(err3)
		}

		arr3 = append(arr3, m)
	}
	fmt.Println("arr3")
	fmt.Println(arr3)

	mapa := make(map[string]interface{})
	mapb := make(map[string]interface{})
	mapc := make(map[string]interface{})

	mSlicea := make([]map[string]interface{}, 0)
	mSliceb := make([]map[string]interface{}, 0)
	mSlicec := make([]map[string]interface{}, 0)

	newDataa := make([]translateLine, 0)
	newDatab := make([]translateLine, 0)
	newDatac := make([]translateLine, 0)

	for i, j := range arr3 {
		fmt.Println(i, j)
		for s1, s := range j {
			fmt.Println(s1, s)
			if s1 == "zh" {
				newData1 := new(translateLine)
				mapa[s1] = s
				mapa["key"] = j["key"]
				mSlicea = append(mSlicea, mapa)
				newData1.tag = j["key"].(string)
				newData1.value = s.(string)
				newDataa = append(newDataa, *newData1)
			} else if s1 == "en" {
				newData1 := new(translateLine)
				mapb[s1] = s
				mapb["key"] = j["key"]
				mSliceb = append(mSliceb, mapb)
				newData1.tag = j["key"].(string)
				newData1.value = s.(string)
				newDatab = append(newDatab, *newData1)
			} else if s1 == "ja" {
				newData1 := new(translateLine)
				mapc[s1] = s
				mapc["key"] = j["key"]
				mSlicec = append(mSlicec, mapc)
				newData1.tag = j["key"].(string)
				newData1.value = s.(string)
				newDatac = append(newDatac, *newData1)
			}
		}
	}

	fmt.Println(mSlicea)
	fmt.Println(mSliceb)
	fmt.Println(mSlicec)

	//func saveAsNewFile(langType, fileName string, data []translateLine)

	t1 := strconv.Itoa(time.Now().Year())       //年
	t2 := strconv.Itoa(int(time.Now().Month())) //月
	t3 := strconv.Itoa(time.Now().Day())        // 日

	wenjian_name := "./" + t1 + "_" + t2 + "_" + t3

	saveAsNewFile("", wenjian_name+"/zh.json", newDataa)
	saveAsNewFile("", wenjian_name+"/en.json", newDatab)
	saveAsNewFile("", wenjian_name+"/ja.json", newDatac)
}
