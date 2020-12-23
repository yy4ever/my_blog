package utils

import (
	"bufio"
	"os"
	"strings"
)

func Confs2Map (fileObj *os.File) *map[string]string {
	confs := make(map[string]string)
	scanner := bufio.NewScanner(fileObj)
	for scanner.Scan() {
		ln := scanner.Text()
		if ln == "" {
			continue
		}
		kv := strings.SplitN(ln, "=", 2)
		k, v := kv[0], kv[1]
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		confs[k] = v
	}
	return &confs
}