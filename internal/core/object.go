package core

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
  )

  // CreateObject creates a Git-like object (blob/tree/commit).
func CreateObject(objType string, content []byte) (string, error) {
	header := fmt.Sprintf("%s %d\x00", objType, len(content))
	data := append([]byte(header), content...)
  
	hash := sha1.Sum(data)
	hashStr := hex.EncodeToString(hash[:])
  
	// Save to .gvc/objects/ab/cdef1234...
	objectPath := filepath.Join(".gvc", "objects", hashStr[:2], hashStr[2:])
	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
	  return "", err
	}
  
	return hashStr, os.WriteFile(objectPath, data, 0644)
  }