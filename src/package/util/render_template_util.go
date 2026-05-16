package util

import (
	"bytes"
	"html/template"
	"strings"

	tmplfs "social-platform-backend/package/template"
)

func RenderTemplate(filePath string, data interface{}) (string, error) {
	// Strip "package/template/" prefix to get the path relative to the embedded FS
	fsPath := strings.TrimPrefix(filePath, "package/template/")

	t, err := template.ParseFS(tmplfs.FS, fsPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
