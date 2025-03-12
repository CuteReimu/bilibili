package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func snakeToCamel(snakeStr string) string {
	parts := strings.Split(snakeStr, "_")
	for i, word := range parts {
		parts[i] = strings.Title(word) // nolint: staticcheck
	}
	return strings.Join(parts, "")
}

var headMapping = map[string]string{
	"参数名": "name",
	"字段":  "name",
	"字段名": "name",
	"项":   "name",
	"类型":  "type",
	"内容":  "content",
	"必要性": "notnull",
	"备注":  "comment",
}

func main() {
	for {
		var head []string
		var lines []map[string]string
		fmt.Println("请输入Markdown表格，在最后一行之后输入ok表示结束：（退出请输入exit）")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			if line == "exit" {
				os.Exit(0)
			}
			if line == "ok" {
				break
			}
			line = strings.ReplaceAll(line, "`", "")
			if len(head) == 0 {
				head = strings.Split(line, "|")
				for i, v := range head {
					head[i] = strings.TrimSpace(v)
					head[i] = headMapping[head[i]]
				}
			} else {
				l := strings.Split(line, "|")
				if len(l) != len(head) {
					fmt.Println("列数不一致，无法解析")
					os.Exit(0)
				}
				m := make(map[string]string)
				for i, v := range l {
					v = strings.TrimSpace(v)
					if strings.HasPrefix(v, "---") {
						m = nil
						break
					}
					m[head[i]] = v
				}
				if m != nil {
					lines = append(lines, m)
				}
			}
		}
		fmt.Println("type T struct {")
		for _, m := range lines {
			name := snakeToCamel(m["name"])
			switch m["type"] {
			case "num":
				m["type"] = "int"
			case "str":
				m["type"] = "string"
			case "bool":
			case "array", "Array", "list", "List", "array(obj)", "Array(obj)", "list(obj)", "List(obj)":
				m["type"] = "[]"
				if !strings.HasSuffix(name, "s") {
					m["type"] += name
				} else {
					m["type"] += name[:len(name)-1]
				}
			case "array(num)", "Array(num)", "list(num)", "List(num)":
				m["type"] = "[]int"
			case "array(str)", "Array(str)", "list(str)", "List(str)":
				m["type"] = "[]string"
			case "array(bool)", "Array(bool)", "list(bool)", "List(bool)":
				m["type"] = "[]bool"
			case "obj":
				m["type"] = name
			}
			notNullValue := `"`
			if m["notnull"] != "必要" && m["notnull"] != "必须" && m["notnull"] != "必填" && m["notnull"] != "√" {
				notNullValue = `,omitempty" request:"query,omitempty"`
			}
			content := m["content"]
			comment := m["comment"]
			sep := ""
			if content != "" && comment != "" {
				sep = "。"
			}
			commentText := fmt.Sprintf("// %s%s%s", content, sep, comment)
			commentText = strings.ReplaceAll(commentText, "<br>", "。")
			commentText = strings.ReplaceAll(commentText, "<br/>", "。")
			commentText = strings.ReplaceAll(commentText, "<br />", "。")
			commentText = strings.ReplaceAll(commentText, `\`, "")
			fmt.Printf("\t%s %s `json:\"%s%s`%s\n", name, m["type"], m["name"], notNullValue, commentText)
		}
		fmt.Println("}")
	}
}
