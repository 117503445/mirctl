package utils

import (
	"bytes"
	"fmt"
	"text/template"
)

func RenderTemplate(tmpl string, data any) (string, error) {
	// 创建一个新的模板并解析模板字符串。
	t := template.New("")
	t, err := t.Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	// 创建一个缓冲区用于存储渲染后的输出。
	var buf bytes.Buffer

	// 使用给定的数据执行模板，并将结果写入缓冲区。
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	// 返回缓冲区中的内容作为字符串。
	return buf.String(), nil
}
