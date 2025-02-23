package template

import (
	"bytes"
	"html/template"
)

// TemplateManager 定义模板管理接口
type TemplateManager interface {
	// ExecuteTemplate 执行指定模板并返回渲染结果
	ExecuteTemplate(name string, data interface{}) (string, error)
	// AddTemplate 添加新模板
	AddTemplate(name, content string) error
}

// templateManager 实现模板管理接口
type templateManager struct {
	templates *template.Template
}

// NewTemplateManager 创建模板管理器实例
func NewTemplateManager() (TemplateManager, error) {
	tm := &templateManager{
		templates: template.New(""),
	}
	if err := tm.InitDefaultTemplates(); err != nil {
		return nil, err
	}
	return tm, nil
}

// ExecuteTemplate 执行指定模板
func (tm *templateManager) ExecuteTemplate(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// AddTemplate 添加新模板
func (tm *templateManager) AddTemplate(name, content string) error {
	_, err := tm.templates.New(name).Parse(content)
	return err
}
