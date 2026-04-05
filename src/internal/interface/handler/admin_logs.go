package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"social-platform-backend/package/logger"
	"social-platform-backend/web/admin"

	"github.com/gin-gonic/gin"
)

// ServeAdminLogsDashboard serves the static HTML UI.
func ServeAdminLogsDashboard(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", admin.LogsDashboardHTML)
}

// GetAdminLogFiles lists log files in the log directory (active + rotated backups).
func GetAdminLogFiles(c *gin.Context) {
	files, err := logger.ListLogFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"dir":   logger.LogDirectory(),
		"files": files,
	})
}

// GetAdminLogs returns the last N lines from a chosen log file as JSON.
// Query params:
//   - file  : base name of the log file (default: active log file)
//   - lines : max lines to return (default 200, max 2000)
//   - search: optional search value
//   - search_type: one of "user_id", "token", "ip" (default: "")
//   - level: one of "DEBUG", "INFO", "WARN", "ERROR" (default: "")
func GetAdminLogs(c *gin.Context) {
	active := logger.LogFilePath()
	if active == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "log file path unavailable"})
		return
	}

	fileParam := c.Query("file")
	if fileParam == "" {
		fileParam = filepath.Base(active)
	}

	path, err := logger.ResolveLogFile(fileParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	n, _ := strconv.Atoi(c.DefaultQuery("lines", "200"))
	if n < 1 {
		n = 200
	}
	if n > 2000 {
		n = 2000
	}

	search := strings.TrimSpace(c.Query("search"))
	searchType := strings.ToLower(strings.TrimSpace(c.Query("search_type")))
	level := strings.ToUpper(strings.TrimSpace(c.Query("level")))

	const maxRead = 2 * 1024 * 1024

	// When search is applied, we read more lines (up to maxRead) to increase chances of matching records,
	// then filter down to N results in memory.
	readN := n
	if search != "" {
		readN = 10000
	}

	lines, err := logger.ReadTailLines(path, readN, maxRead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reverseLines(lines)

	if search != "" || level != "" {
		lines = filterLogLines(lines, searchType, search, level)
		// Lines are already newest-first, keep only top N newest records.
		if len(lines) > n {
			lines = lines[:n]
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"path":        path,
		"lines":       lines,
		"search":      search,
		"search_type": searchType,
		"level":       level,
	})
}

func GetAdminMetrics(c *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	requestMetrics := logger.SnapshotRequestMetrics()

	c.JSON(http.StatusOK, gin.H{
		"time": time.Now().UTC().Format(time.RFC3339),
		"system": gin.H{
			"num_goroutine": runtime.NumGoroutine(),
			"memory": gin.H{
				"alloc_bytes":       mem.Alloc,
				"total_alloc_bytes": mem.TotalAlloc,
				"heap_alloc_bytes":  mem.HeapAlloc,
				"heap_inuse_bytes":  mem.HeapInuse,
				"heap_sys_bytes":    mem.HeapSys,
				"stack_inuse_bytes": mem.StackInuse,
				"stack_sys_bytes":   mem.StackSys,
				"num_gc":            mem.NumGC,
			},
		},
		"request": requestMetrics,
	})
}

// reverseLines flips the order so the newest record appears first.
func reverseLines(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
}

// filterLogLines filters log lines by search type. If searchType is empty or unknown,
// it falls back to case-insensitive substring matching on raw JSON line.
func filterLogLines(lines []string, searchType string, query string, level string) []string {
	q := strings.TrimSpace(query)
	lv := strings.ToUpper(strings.TrimSpace(level))
	if q == "" && lv == "" {
		return lines
	}
	qLower := strings.ToLower(q)
	qTokenHash := sha256Hex(q)
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if !matchLogLevel(l, lv) {
			continue
		}
		if q == "" || matchLogLine(l, searchType, q, qLower, qTokenHash) {
			out = append(out, l)
		}
	}
	return out
}

func matchLogLevel(raw string, level string) bool {
	if level == "" {
		return true
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return strings.Contains(strings.ToUpper(raw), `"LEVEL":"`+level+`"`)
	}

	v, ok := obj["level"]
	if !ok {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(anyToString(v)), level)
}

func matchLogLine(raw string, searchType string, query string, queryLower string, queryTokenHash string) bool {
	if searchType == "" {
		return strings.Contains(strings.ToLower(raw), queryLower)
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return strings.Contains(strings.ToLower(raw), queryLower)
	}

	switch searchType {
	case "user_id":
		v, ok := obj["user_id"]
		if !ok {
			return false
		}
		return strings.EqualFold(strings.TrimSpace(anyToString(v)), query)
	case "ip":
		v, ok := obj["client_ip"]
		if !ok {
			return false
		}
		return strings.Contains(strings.ToLower(anyToString(v)), queryLower)
	case "token":
		if hint, ok := obj["token_hint"]; ok {
			if strings.Contains(strings.ToLower(anyToString(hint)), queryLower) {
				return true
			}
		}
		if h, ok := obj["token_hash"]; ok {
			return strings.EqualFold(strings.TrimSpace(anyToString(h)), queryTokenHash)
		}
		return false
	default:
		return strings.Contains(strings.ToLower(raw), queryLower)
	}
}

func anyToString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
