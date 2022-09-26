package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/xuri/excelize/v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const LINE_RDWR_SIZE = 4096 * 1024

type translateLine struct {
	tag   string
	value string
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

// Convert map json string
func MapToJson(m map[string]string) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return "", nil
	}

	return string(jsonByte), nil
}

var basePath string

func init() {
	flag.StringVar(&basePath, "basePath", "", "词条目录--规定目录下的文件名必需是(en.json、ja.json、en.json)")

	flag.Parse()
}

func main() {
	BasePath := &basePath
	files, _ := ioutil.ReadDir(*BasePath)
	fmt.Println(len(files))

	var s1 [3]string
	for i, f := range files {
		zh_diff := filepath.Join(*BasePath, f.Name())
		fmt.Println(zh_diff)
		s1[i] = zh_diff
	}
	fmt.Println(s1)

	countryCapitalMap := make(map[string][]translateLine)
	for i := range s1 {
		split := strings.Split(s1[i], "/")
		currentData := readFile(s1[i])
		sq := split[len(split)-1]
		acx := sq[:2]
		countryCapitalMap[acx] = currentData
	}
	fmt.Println(countryCapitalMap)
	countMap := make(map[string]int)

	//对不同的key值中的不同的取并集  凡事缺少一个翻译就直接提取出来
	for countryCapitalMap_key, countryCapitalMap_value := range countryCapitalMap {
		for _, i := range countryCapitalMap_value {
			countMap[i.tag]++
		}
		fmt.Println(countryCapitalMap_key, countryCapitalMap_value)
	}

	DataMap := make(map[string]map[string]string)

	//如果值不为文件夹中文件的总数时，则输出相应的数值操作
	for key, value := range countMap {
		fmt.Println(key, value)
		if value != 3 {
			numbers := make(map[string]string)
			for countryCapitalMap_key, countryCapitalMap_value := range countryCapitalMap {
				numbers[countryCapitalMap_key] = ""
				for _, i := range countryCapitalMap_value {
					if key == i.tag {
						numbers[countryCapitalMap_key] = i.value
					}
				}
				DataMap[key] = numbers
			}
		}
	}

	DataMap1 := make(map[string]string, 2)
	jsonRes, err := MapToJson(DataMap1)
	var identifier []string
	for key, DataMap_value := range DataMap {
		DataMap1["key"] = key
		for DataMap_key, DataMap_value_value := range DataMap_value {
			DataMap1[DataMap_key] = DataMap_value_value
		}
		jsonRes, err = MapToJson(DataMap1)
		identifier = append(identifier, jsonRes)
	}
	fmt.Println(identifier)

	jsonStr := strings.Replace(strings.Trim(fmt.Sprint(identifier), ""), "} {", "},{", -1)
	fmt.Println(jsonStr)
	mSlice := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(jsonStr), &mSlice)
	if err != nil {
		fmt.Println("反序列化失败")
	} else {
		fmt.Println(mSlice)
	}

	//DataTitle := [4]string{
	//	"key",
	//	"zh",
	//	"en",
	//	"ja",
	//}

	//f := excelize.NewFile()

	// 创建一个工作表
	//index := f.NewSheet("Sheet1")
	//设置表头
	//for i, title := range DataTitle {
	//	cellNum := fmt.Sprintf("%c1", 65+i)
	//	_ = f.SetCellValue("Sheet1", cellNum, title)
	//}
	var data []map[string]interface{}

	for key, _ := range mSlice {
		data = append(data, mSlice[key])
	}

	fmt.Println(data)

	//HeadKey := [4]string{
	//	"key",
	//	"zh",
	//	"en",
	//	"ja",
	//}

	//for i, rows := range data { //一行
	//	for num, key := range HeadKey { //每列
	//		line := fmt.Sprintf("%c%d", 65+num, i+2)
	//		//正常值
	//		enc := mahonia.NewDecoder("utf-8")
	//		resData := enc.ConvertString(rows[key].(string))
	//		_ = f.SetCellValue("Sheet1", line, resData)
	//	}
	//}
	//// 设置工作簿的默认工作表
	//f.SetActiveSheet(index)

	t1 := strconv.Itoa(time.Now().Year())       //年
	t2 := strconv.Itoa(int(time.Now().Month())) //月
	t3 := strconv.Itoa(time.Now().Day())        // 日

	wenjian_name := "./" + t1 + "_" + t2 + "_" + t3

	err = os.Mkdir(wenjian_name, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	//wenjian_name1 := t1 + "_" + t2 + "_" + t3

	//bea := wenjian_name1 + "/" + "Book2.xlsx"

	// 根据指定路径保存文件
	//if err := f.SaveAs(bea); err != nil {
	//	fmt.Println(err)
	//}

	//进行提取  key  ready_to_translate（待翻译的文案）   ja
	//如果需要中文  中文没有补上英文
	var ready_to_translate string
	var new_map map[string]string

	mSlice2 := make([]map[string]string, 0)
	for i, j := range data {
		ready_to_translate = " "
		fmt.Println(i, j)
		new_map = make(map[string]string)
		if j["ja"].(string) == "" {
			ready_to_translate_zh := j["zh"].(string)
			if ready_to_translate_zh == "" {
				ready_to_translate_en := j["en"].(string)
				if ready_to_translate_en == "" {
					ready_to_translate = ""
				} else {
					ready_to_translate = j["en"].(string)
				}
			} else {
				ready_to_translate = j["zh"].(string)
			}
			new_map["key"] = j["key"].(string)
			new_map["ready_to_translate"] = ready_to_translate
			new_map["ja"] = " "
			mSlice2 = append(mSlice2, new_map)
		}
	}
	fmt.Println(mSlice2)
	//中文英文 都没有 则为空

	DataTitle1 := [4]string{
		"key",
		"ready_to_translate",
		"ja",
	}

	f1 := excelize.NewFile()

	// 创建一个工作表
	index1 := f1.NewSheet("Sheet1")
	//设置表头
	for i, title := range DataTitle1 {
		cellNum := fmt.Sprintf("%c1", 65+i)
		_ = f1.SetCellValue("Sheet1", cellNum, title)
	}
	var data1 []map[string]string

	for key, _ := range mSlice2 {
		data1 = append(data1, mSlice2[key])
	}

	fmt.Println(56757121)
	fmt.Println(data1)

	HeadKey1 := [4]string{
		"key",
		"ready_to_translate",
		"ja",
	}

	for i, rows := range data1 { //一行
		for num, key := range HeadKey1 { //每列
			line := fmt.Sprintf("%c%d", 65+num, i+2)
			//正常值
			enc := mahonia.NewDecoder("utf-8")
			resData := enc.ConvertString(rows[key])
			_ = f1.SetCellValue("Sheet1", line, resData)
		}
	}
	// 设置工作簿的默认工作表
	f1.SetActiveSheet(index1)

	wenjian_name2 := t1 + "_" + t2 + "_" + t3

	bea1 := wenjian_name2 + "/" + "execl_one.xlsx"

	// 根据指定路径保存文件
	if err := f1.SaveAs(bea1); err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Printf("Convert json to map failed with error: %+v\n", err)
	}
}
