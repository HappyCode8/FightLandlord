package ui

import (
	"fmt"
	"os"
	"text/template"
)

type UI struct {
	LastCards []string
	Cards     []string
	Chat      []string
}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) CreateOutput() {
	data := struct {
		Cards []string
		Chat  []string
	}{
		Cards: []string{"牌1", "牌2", "牌3"},
		Chat:  []string{"消息1", "消息2", "消息3"},
	}

	// 读取模板文件
	fileContent, err := os.ReadFile("./ui/ui.prob")
	if err != nil {
		fmt.Println("读取文件时出错:", err)
		return
	}

	// 解析模板
	tmpl, err := template.New("ui.prob").Parse(string(fileContent))
	if err != nil {
		fmt.Println("解析模板时出错:", err)
		return
	}

	// 执行模板并输出结果
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println("执行模板时出错:", err)
		return
	}
}
