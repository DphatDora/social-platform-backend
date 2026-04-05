package logger

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LogFileEntry describes a .log file in the log directory (active file + rotated backups).
type LogFileEntry struct {
	Name    string `json:"name"`
	Current bool   `json:"current"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"modTimeUnix"`
}

// LogDirectory returns the absolute directory containing the active log file.
func LogDirectory() string {
	p := LogFilePath()
	if p == "" {
		return ""
	}
	return filepath.Dir(p)
}

// ListLogFiles returns app.log and any rotated *.log files in the same directory, newest first after the current file.
func ListLogFiles() ([]LogFileEntry, error) {
	active := LogFilePath()
	if active == "" {
		return nil, errors.New("log file path unavailable")
	}
	dir := filepath.Dir(active)
	currentBase := filepath.Base(active)
	des, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []LogFileEntry
	for _, de := range des {
		if de.IsDir() {
			continue
		}
		name := de.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".log") {
			continue
		}
		info, err := de.Info()
		if err != nil {
			continue
		}
		out = append(out, LogFileEntry{
			Name:    name,
			Current: name == currentBase,
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Current != out[j].Current {
			return out[i].Current
		}
		return out[i].ModTime > out[j].ModTime
	})
	return out, nil
}

// ResolveLogFile returns an absolute path to a log file under the log directory.
// fileName must be a base name only (e.g. app.log), ending with .log.
func ResolveLogFile(fileName string) (string, error) {
	active := LogFilePath()
	if active == "" {
		return "", errors.New("log file path unavailable")
	}
	dir := filepath.Dir(active)
	base := filepath.Base(strings.TrimSpace(fileName))
	if base == "." || base == "" {
		return "", errors.New("invalid file name")
	}
	if strings.ContainsAny(base, `/\`) {
		return "", errors.New("invalid file name")
	}
	if !strings.HasSuffix(strings.ToLower(base), ".log") {
		return "", errors.New("not a log file")
	}
	full := filepath.Join(dir, base)
	cleanDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	cleanFull, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(cleanDir, cleanFull)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", errors.New("invalid path")
	}
	return cleanFull, nil
}
