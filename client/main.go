package main

import (
	"client/shell"
	"flag"
	"fmt"
	"log"
)

var (
	host string
	port int
	name string
)

func init() {
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.IntVar(&port, "p", 9999, "port")
	flag.StringVar(&name, "n", "", "name")
	flag.Parse()
}

func main() {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Fatal(shell.New(addr, name).Start())
}

/*func main() {
	// 定义要填充的数据
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
*/
