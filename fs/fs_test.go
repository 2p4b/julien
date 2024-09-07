package fs

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsFileType tests the IsFileType function
func TestHasExt(t *testing.T) {
	// Test with a matching extension
	assert.True(t, HasExt("example.txt", "txt"))

	// Test with a non-matching extension
	assert.False(t, HasExt("example.txt", "jpg"))

	// Test with an empty extension (should match any file)
	assert.True(t, HasExt("example.txt", ""))
}

// TestGetNameParts tests the GetNameParts function
func TestGetNameParts(t *testing.T) {
	// Test with a regular file name
	name, ext := GetNameParts("example.txt")
	assert.Equal(t, "example", name)
	assert.Equal(t, "txt", ext)

	// Test with a file name without an extension
	name, ext = GetNameParts("example")
	assert.Equal(t, "example", name)
	assert.Equal(t, "", ext)

	// Test with a file name that starts with a dot
	name, ext = GetNameParts(".hiddenfile")
	assert.Equal(t, "", name)
	assert.Equal(t, "hiddenfile", ext)
}

// TestMount tests the Mount function
func TestMount(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Mount the directory
	disk := Mount(tmpDir, "index", "md")

	// Check the Disk fields
	assert.Equal(t, tmpDir, disk.Root())
	assert.Equal(t, "index", disk.Index())
	assert.Equal(t, "md", disk.Ext())
}

// TestList tests the List method of Disk
func TestList(t *testing.T) {
	// Create a temporary directory with some files
	tmpDir := t.TempDir()
	os.WriteFile(path.Join(tmpDir, "file1.md"), []byte("content"), 0644)
	os.WriteFile(path.Join(tmpDir, "file2.md"), []byte("content"), 0644)

	// Mount the directory
	disk := Mount(tmpDir, "index", "md")

	// List files with the .txt extension
	entries, err := disk.List("")
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "file1", entries[0].Name())
	assert.Equal(t, "md", entries[0].Ext())
}

// TestFind tests the Find method of Disk
func TestFind(t *testing.T) {
	// Create a temporary directory with a file
	tmpDir := t.TempDir()
	filePath := path.Join(tmpDir, "file1.md")
	os.WriteFile(filePath, []byte("content"), 0644)

	// Mount the directory
	disk := Mount(tmpDir, "index", "md")

	// Find the file
	entry, err := disk.Find("file1.md")
	assert.NoError(t, err)
	assert.Equal(t, "file1", entry.Name())
	assert.Equal(t, "md", entry.Ext())
	assert.True(t, entry.IsFile())
}
