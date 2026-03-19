package lcu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadLockfile() (*LockfileData, error) {
	data, err := os.ReadFile(filepath.Join("C:\\", "Riot Games", "League of Legends", "lockfile"))
	if err != nil {
		return nil, fmt.Errorf("failed to read lockfile: %w", err)
	}

	parts := strings.Split(string(data), ":")

	return &LockfileData{
		ProcessName: parts[0],
		PID:         parts[1],
		Port:        parts[2],
		Password:    parts[3],
		Protocol:    parts[4],
	}, nil

}
