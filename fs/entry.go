package fs

import (
	"fmt"
	"os"
	"path"
	"time"
)

// Entry represents a single file or directory in the file system.
type Entry struct {
	ext     string      // File extension.
	path    string      // Full path to the entry.
	name    string      // Base name of the entry.
	is_file bool        // Whether the entry is a file.
	disk    *Disk       // Reference to the Disk this entry belongs to.
	info    os.FileInfo // File information (size, permissions, etc.).
}

func (entry *Entry) Dir() (*Entry, error) {
	return entry.disk.Find(path.Dir(entry.Path()))
}

// IsIndexed checks if the entry is indexable.
//
// Returns:
// - true if the entry is indexable, false otherwise.
func (entry *Entry) IsIndexed() bool {
	if entry.IsFile() {
		return true
	}
	file, err := entry.disk.fs.Open(entry.IndexPath())
	if err != nil {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (entry *Entry) HasExt(ext string) bool {
	return entry.Ext() == ext
}

// IsDir checks if the entry is a directory.
//
// Returns:
// - true if the entry is a directory, false otherwise.
func (entry *Entry) IsDir() bool {
	return entry.info.IsDir()
}

// IsFile checks if the entry is a file.
//
// Returns:
// - true if the entry is a file, false otherwise.
func (entry *Entry) IsFile() bool {
	return entry.is_file
}

// Size returns the size of the entry in bytes.
//
// Returns:
// - The size of the file in bytes.
func (entry *Entry) Size() int64 {
	return entry.info.Size()
}

// Name returns the base name of the entry.
//
// Returns:
// - The base name of the entry as a string.
func (entry *Entry) Name() string {
	return entry.name
}

// IsIndex checks if the entry is an index file.
//
// Returns:
// - true if the entry is an index file, false otherwise.
func (entry *Entry) IsIndex() bool {
	if entry.IsDir() {
		return false
	}
	disk := entry.disk
	ifname := disk.IName()
	filename := entry.Filename()
	if filename == disk.Index() || filename == (ifname) {
		return true
	}
	return false
}

// Ext returns the extension of the entry.
//
// Returns:
// - The extension of the entry as a string.
func (entry *Entry) Ext() string {
	return entry.ext
}

// Path returns the full path of the entry within the Disk.
//
// Returns:
// - The full path of the entry as a string.
func (entry *Entry) Path() string {
	return entry.path
}

// Filename returns the name of the file in the entry.
//
// Returns:
// - The name of the file as a string.
func (entry *Entry) Filename() string {
	return entry.info.Name()
}

// ModTime returns the modification time of the entry.
//
// Returns:
// - The modification time as a string.
func (entry *Entry) ModTime() string {
	return entry.info.ModTime().String()
}

// IndexPath returns the full path to the index file if the entry is a directory, or the file path if it's a file.
//
// Returns:
// - The full path to the index file or file path as a string.
func (entry *Entry) IndexPath() string {
	if entry.IsFile() {
		return entry.Path()
	}
	return path.Clean(path.Join(entry.Path(), entry.disk.index+"."+entry.disk.ext))
}

// Fullpath returns the absolute path of the entry in the Disk.
//
// Parameters:
// - ppath: The relative path within the disk.
//
// Returns:
// - The absolute path as a string.
func (entry *Entry) Datapath() string {
	return path.Clean(entry.disk.root + "/" + entry.IndexPath())
}

// Read reads the content of the file represented by the entry.
//
// Returns:
// - The content of the file as a byte slice.
// - An error if the file cannot be read.
func (entry *Entry) Read() ([]byte, error) {
	if !entry.IsIndexed() {
		return nil, fmt.Errorf("index not found at: %s", entry.IndexPath())
	}
	content, err := os.ReadFile(entry.Datapath())
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Write writes content to the file represented by the entry.
//
// Parameters:
// - content: The byte slice to write to the file.
//
// Returns:
// - An error if the file cannot be written.
func (entry *Entry) Write(content []byte) error {
	if !entry.IsIndexed() {
		return fmt.Errorf("index not found at: %s", entry.IndexPath())
	}
	return entry.disk.Dump(entry.IndexPath(), []byte(content))
}

// Append appends content to the file represented by the entry.
//
// Parameters:
// - content: The byte slice to append to the file.
//
// Returns:
// - An error if the file cannot be appended to.
func (entry *Entry) Append(content []byte) error {
	if !entry.IsIndexed() {
		return fmt.Errorf("index not found at: %s", entry.IndexPath())
	}
	return entry.disk.Append(entry.IndexPath(), []byte(content))
}

// Timestamp returns the modification time of the entry.
//
// Returns:
// - The modification time as a time.Time.
func (entry *Entry) Timestamp() time.Time {
	return entry.info.ModTime()
}
