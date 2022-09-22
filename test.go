package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

type wcSendcontent struct {
	Content string `json:"content"`
}

type WcSendMsg struct {
	MsgType string        `json:"msgtype"`
	Text    wcSendcontent `json:"text"`
}

// ExecShell ...
func ExecShell(command string, arg ...string) (out string, err error) {
	var Stdout []byte
	cmd := exec.Command(command, arg...)
	Stdout, err = cmd.CombinedOutput()
	out = string(Stdout)
	return
}

// Repo ...
// git config core.hooksPath .githooks  解决方式1
//                                      解决方式2
func Repo() (repo string, err error) {
	var (
		out string
	)
	if out, err = ExecShell("/bin/sh", "-c", "git remote -v"); err != nil {
		return
	}
	if repo = out[strings.Index(out, ":")+1 : strings.Index(out, ".git")]; repo == "" {
		err = fmt.Errorf("not found, %s", out)
		return
	}
	return
}

// Branch ...
func Branch() (branch string, err error) {
	var (
		out string
	)
	if out, err = ExecShell("/bin/sh", "-c", "git branch"); err != nil {
		return
	}
	list := strings.Split(out, "\n")
	for _, v := range list {
		if strings.HasPrefix(v, "*") {
			branch = v[strings.Index(v, "*")+2:]
			return
		}
	}
	err = fmt.Errorf("not found, %s", out)
	return
}

func SplitString(s string, myStrings []rune) []string {
	Split := func(r rune) bool {
		for _, v := range myStrings {
			if v == r {
				return true
			}
		}
		return false
	}
	a := strings.FieldsFunc(s, Split)
	return a
}

//git diff --name-only HEAD~ HEAD
func main() {

	cmd1, _ := ExecShell("git", "diff", "--name-only", "HEAD~", "HEAD")

	fmt.Println(cmd1)

	cmd, _ := ExecShell("git", "log", "--oneline")
	fmt.Println(cmd)
	myStrings := SplitString(cmd, []rune{' ', '\n'})

	now_commid := myStrings[len(myStrings)-2]
	master_commid := myStrings[0]
	//fmt.Println(myStrings[0], myStrings[len(myStrings)-2])
	fmt.Println()
	fmt.Println()
	fmt.Println()

	var data []string
	diff, _ := ExecShell("git", "diff", now_commid, master_commid, "test1.txt")
	if diff != "" {
		diffs := strings.Split(diff, "\n")
		for i, diff_name := range diffs {
			fmt.Println(i, diff_name)
			reg := regexp.MustCompile(`^[\+\-]{1}[^\+|\-].*`)
			s1 := reg.FindAllString(diff_name, -1)
			if s1 != nil {
				fmt.Println(s1)
				data = append(data, s1[0])
			}
		}
		fmt.Println(data)
	}

	var str_shan []string
	var str_add []string

	for _, j := range data {
		if j[0] == '-' {
			str_shan = append(str_shan, j[1:])
		}
		if j[0] == '+' {
			str_add = append(str_add, j[1:])
		}
	}

	fmt.Println(str_shan)
	fmt.Println(str_add)

	str_shans := strings.Join(str_shan, ",")
	str_adds := strings.Join(str_add, ",")

	//删除
	str_shans_connect := fmt.Sprintf("删除词条%s \n 添加词条%s", str_shans, str_adds)

	SendCardMsg(cmd1, str_shans_connect)
}

//git diff --name-only HEAD~ HEAD git比较

//企业微信应用消息提醒方法如下
func SendCardMsg(fliename string, str_shans_connect string) (WcSendMsg, error) {
	//查看是那些词条文件发生了改变
	flie_diff := fmt.Sprintf("%s 词条文件发生改变\n %s", fliename, str_shans_connect)

	req := WcSendMsg{MsgType: "text", Text: wcSendcontent{Content: flie_diff}}

	sendurl := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=cd75d3a3-0899-4c63-a1a7-fe578784b9e2"
	data, err := httpPostJson(sendurl, req)
	if err != nil {
		log.Println(err)
		return WcSendMsg{MsgType: "", Text: wcSendcontent{Content: ""}}, err
	}
	return data, nil
}

func httpPostJson(url string, data WcSendMsg) (WcSendMsg, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return WcSendMsg{MsgType: "", Text: wcSendcontent{Content: ""}}, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(res))
	if err != nil {
		return WcSendMsg{MsgType: "", Text: wcSendcontent{Content: ""}}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return WcSendMsg{MsgType: "", Text: wcSendcontent{Content: ""}}, err
	}
	return data, nil
}
