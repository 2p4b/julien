package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

// Disk represents a file system mounted at a specific root with a given extension and index file.
type Disk struct {
	fs    fs.FS       // The file system interface.
	root  string      // The root directory of the file system.
	ext   string      // Default file extension to use.
	index string      // Default index file name.
	fmode fs.FileMode // File mode permissions for files.
	dmode fs.FileMode // File mode permissions for directories.
}

// HasExt checks if the given filename has the specified extension.
//
// Parameters:
// - filename: The name of the file to check.
// - ext: The extension to compare against (e.g., "txt").
//
// Returns:
// - true if the file has the specified extension, false otherwise.
func HasExt(filename string, ext string) bool {
	if ext == "" {
		return true
	}
	_, name_ext := GetNameParts(filename)
	return name_ext == ext
}

// GetNameParts splits the filename into the base name and extension.
//
// Parameters:
// - filename: The full name of the file.
//
// Returns:
// - base: The base name of the file (without extension).
// - ext: The extension of the file (without the leading dot).
func GetNameParts(filename string) (string, string) {
	ext := path.Ext(filename)
	base := path.Base(filename)

	ext_len := len(ext)
	base_len := len(base)

	if ext_len == 0 {
		return base, ""
	}

	if base_len == 0 {
		return "", ext[1:]
	}

	return filename[0 : base_len-ext_len], ext[1:]
}

// Mount creates a new Disk mounted at the given root directory.
//
// Parameters:
// - root: The root directory to mount.
// - index: The default index file name (e.g., "index.html").
// - ext: The default file extension to use.
//
// Returns:
// - A pointer to the newly created Disk instance.
func Mount(root string, index string, ext string) *Disk {
	return &Disk{
		root:  root,
		fs:    os.DirFS(root),
		index: index,
		fmode: 0644,
		dmode: 0755,
		ext:   ext,
	}
}

// create_entry creates a new Entry from the given FileInfo and path.
//
// Parameters:
// - entry: The FileInfo object containing details about the file or directory.
// - fpath: The full path to the entry.
//
// Returns:
// - A new Entry instance representing the file or directory.
func (d *Disk) create_entry(entry fs.FileInfo, fpath string) Entry {
	filename := path.Base(fpath)
	fname, ext := GetNameParts(filename)
	info := entry
	return Entry{
		ext:     ext,
		info:    info,
		path:    fpath,
		name:    fname,
		is_file: !entry.IsDir(),
		disk:    d,
	}
}

// Root returns the root directory of the Disk.
//
// Returns:
// - A string representing the root directory path.
func (disk *Disk) Root() string {
	return disk.root
}

// Ext returns the default file extension of the Disk.
//
// Returns:
// - A string representing the default file extension.
func (disk *Disk) Ext() string {
	return disk.ext
}

// Index returns the default index name of the Disk.
//
// Returns:
// - A string representing the default index name.
func (disk *Disk) Index() string {
	return disk.index
}

// Index returns the default index filename of the Disk.
//
// Returns:
// - A string representing the default index filename.
func (disk *Disk) IName() string {
	return disk.index + "." + disk.ext
}

// Fmode returns the file mode permissions for files.
//
// Returns:
// - The file mode permissions as a fs.FileMode value.
func (disk *Disk) Fmode() fs.FileMode {
	return disk.fmode
}

// Dmode returns the file mode permissions for directories.
//
// Returns:
// - The directory mode permissions as a fs.FileMode value.
func (disk *Disk) Dmode() fs.FileMode {
	return disk.dmode
}

// List lists all entries in the specified directory.
//
// Parameters:
// - rpath: The relative path within the disk's root to list.
//
// Returns:
// - A slice of pointers to Entry objects representing the files and directories in the specified path.
// - An error if the directory cannot be read.
func (disk *Disk) List(rpath string) ([]*Entry, error) {
	dirpath := path.Join(disk.root, rpath)
	dentries := make([]*Entry, 0)

	entries, err := os.ReadDir(dirpath)
	if err != nil {
		return dentries, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			fullpath := path.Join(dirpath, entry.Name())
			indexpath := path.Join(fullpath, disk.IName())
			_, err := os.Stat(path.Join(indexpath))
			if err != nil {
				continue
			}

		} else if !entry.Type().IsRegular() || !HasExt(entry.Name(), disk.ext) {
			continue
		}
		info, _ := entry.Info()
		ppath := path.Join(rpath, entry.Name())
		dentry := disk.create_entry(info, ppath)
		dentries = append(dentries, &dentry)
	}

	return dentries, nil
}

// Find locates a file or directory within the Disk.
//
// Parameters:
// - filepath: The relative path to the file or directory.
//
// Returns:
// - A pointer to the Entry if found.
// - An error if the file or directory cannot be found.
func (disk *Disk) Find(filepath string) (*Entry, error) {
	filepath = path.Clean(filepath)
	filepath = strings.TrimLeft(filepath, "/")
	filepath = strings.TrimLeft(filepath, `\`)
	if filepath == "/" || filepath == `\` {
		filepath = "."
	}
	index_path := filepath

	_, ext := GetNameParts(filepath)
	fileinfo, err := fs.Stat(disk.fs, filepath)

	// Check if file exist
	if err != nil {
		// File with ext not found
		if ext != "" {
			return nil, err
		}

		// Append disk extension  and try again
		index_path = filepath + "." + disk.ext

		fileinfo, err = fs.Stat(disk.fs, index_path)
		if err != nil {
			return nil, err
		}
		// Reject directory with extension
		// not something i feel like will be usefull
		if fileinfo.IsDir() {
			return nil, fmt.Errorf("index file not found: %s", filepath)
		}

		// Found file by appending ext
		filepath = index_path
	}

	if fileinfo.IsDir() {
		// Check dir for index file
		if filepath == "." {
			index_path = disk.index + "." + disk.ext
		} else {
			index_path = path.Join(filepath, disk.index+"."+disk.ext)
		}
		_, err := fs.Stat(disk.fs, index_path)

		if err != nil {
			return nil, err
		}
	}

	file, err := disk.fs.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	dentry := disk.create_entry(info, filepath)

	return &dentry, nil
}

// Read reads the content of a file in the Disk.
//
// Parameters:
// - filepath: The relative path to the file to be read.
//
// Returns:
// - The content of the file as a byte slice.
// - An error if the file cannot be read.
func (disk *Disk) Read(filepath string) ([]byte, error) {
	fullpath := path.Join(disk.root, filepath)
	content, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Write writes content to a file in the Disk.
//
// Parameters:
// - filepath: The relative path to the file to be written.
// - content: The byte slice to write to the file.
//
// Returns:
// - An error if the file cannot be written.
func (disk *Disk) Write(filepath string, content []byte) error {
	fullpath := path.Join(disk.root, filepath)
	err := os.WriteFile(fullpath, content, disk.fmode)
	if err != nil {
		return err
	}
	return nil
}

// Remove deletes a file or directory from the Disk.
//
// Parameters:
// - filepath: The relative path to the file or directory to be deleted.
//
// Returns:
// - An error if the file or directory cannot be deleted.
func (disk *Disk) Remove(filepath string) error {
	fullpath := path.Join(disk.root, filepath)
	err := os.RemoveAll(fullpath)
	if err != nil {
		return err
	}
	return nil
}

// Dump dumps file content to the Disk replacing original content
// if file already exist and is not empty.
//
// Parameters:
// - filepath: The relative path to the file to be appended to.
// - content: The byte slice to append to the file.
//
// Returns:
// - An error if the file contents cannot be dumped.
func (disk *Disk) Dump(ppath string, content []byte) error {
	cpath := path.Clean(disk.root + "/" + ppath)

	ext := path.Ext(cpath)
	if ext == "" {
		cpath += "." + disk.ext
	}

	dirpath := path.Dir(cpath)

	_, err := os.Stat(dirpath)

	if err != nil && dirpath != disk.root && dirpath != "" && dirpath != "." && dirpath != "/" {
		if err := os.MkdirAll(dirpath, disk.dmode); err != nil {
			log.Error(err)
			return err
		}
	}

	f, err := os.OpenFile(cpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, disk.fmode)
	if err != nil {
		log.Error(err)
		return err
	}
	defer f.Close()

	f.Truncate(0)
	f.Seek(0, 0)
	if _, err = f.Write(content); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Append appends content to a file in the Disk.
//
// Parameters:
// - filepath: The relative path to the file to be appended to.
// - content: The byte slice to append to the file.
//
// Returns:
// - An error if the file cannot be appended to.
func (disk *Disk) Append(ppath string, content []byte) error {
	cpath := path.Clean(disk.root + "/" + ppath)
	ext := path.Ext(cpath)
	if ext == "" {
		cpath += "." + disk.ext
	}
	f, err := os.OpenFile(cpath, os.O_APPEND|os.O_WRONLY, disk.fmode)
	if err != nil {
		return err
	}

	defer f.Close()
	if _, err = f.Write(content); err != nil {
		return err
	}
	return nil
}

// Create creates and return disk entry.
//
// Parameters:
// - filepath: The relative path to the file to be appended to.
// - content: The byte slice to append to the file.
//
// Returns:
// - The disk entry of the created file.
// - An error if the file cannot be created.
func (disk *Disk) Create(ppath string, content []byte) (*Entry, error) {
	cpath := path.Clean(disk.root + "/" + ppath)
	_, err := disk.Find(ppath)

	if err == nil {
		return nil, fmt.Errorf("file already exists: %s", cpath)
	}

	err = disk.Dump(ppath, content)
	if err != nil {
		return nil, err
	}
	return disk.Find(ppath)
}
