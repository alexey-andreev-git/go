package fstree

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// FsTree is the interface that wraps the basic methods for a filesystem tree.
type FsTree interface {
	AddExclude(exclude string) *Tree
	SetExcludes(excludes []string) *Tree
	GetExcludes() []string
	Fill() *Tree
	IsExecuting() bool
	Wait() *Tree
	Start() *Tree
	Duration() time.Duration
	Print(outFileName string) *Tree
}

type Folder struct {
	Path       string    // Path to the folder
	Files      []string  // List of files in the folder
	Subfolders []*Folder // List of subfolders
	Errors     []string  // List of errors
}

type Tree struct {
	wg        sync.WaitGroup // WaitGroup to wait for all goroutines to finish
	inexec    atomic.Int64   // Counter for the number of goroutines executing
	startTime time.Time      // Time when the execution started
	duration  time.Duration  // Duration of the execution
	Executing bool           // Flag to indicate if the tree is being executed
	Exclude   []string       // List of folders to exclude
	Folders   Folder         // Root folder
}

var maxGoroutines = 1024

// Creates a new Tree with the root path.
func NewTree(rootPath string) *Tree {
	return &Tree{Folders: Folder{Path: rootPath}}
}

// Adds a folder to the list of folders to exclude.
func (t *Tree) AddExclude(exclude string) *Tree {
	t.Exclude = append(t.Exclude, exclude)
	return t
}

// Sets the list of folders to exclude.
func (t *Tree) SetExcludes(excludes []string) *Tree {
	t.Exclude = append(t.Exclude, excludes...)
	return t
}

// Returns the list of folders to exclude.
func (t *Tree) GetExcludes() []string {
	return t.Exclude
}

// Waits for all goroutines to finish.
func (t *Tree) Wait() *Tree {
	t.wg.Wait()
	return t
}

// Defines start time the execution of the tree.
func (t *Tree) Start() *Tree {
	t.startTime = time.Now()
	return t
}

// Returns the duration of the execution.
func (t *Tree) Duration() time.Duration {
	return t.duration
}

// This function is called when a goroutine starts.
// Increment the counter of goroutines executing.
// Set the execution flag to true.
func (t *Tree) startExecution() {
	t.inexec.Add(1)
	t.Executing = true
}

// This function is called when a goroutine finishes.
// Decrement the counter of goroutines executing.
// If the counter reaches 0, set the execution flag to false and calculate the duration.
// Calculate the duration if the start time is not zero (was set by Start function).
func (t *Tree) stopExecution() {
	if t.inexec.Add(-1) == 0 {
		t.Executing = false
		if !t.startTime.IsZero() {
			t.duration = time.Since(t.startTime)
		}
	}
}

// IsExecuting returns true if the tree is being executed.
func (t *Tree) IsExecuting() bool {
	return t.Executing
}

// InExclusion checks if a path is in the list of folders to exclude.
func (t *Tree) InExclusion(path string) bool {
	if len(t.Exclude) != 0 {
		for _, ex := range t.Exclude {
			if path == ex {
				return true // Path is in the exclusion list
			}
		}
	}
	return false // Path is not in the exclusion list
}

// Starts process of filling the tree with the files and folders.
// It creates a new goroutine for the root folder.
// It returns the tree to allow chaining of methods.
// It calls the walkFolder function internally.
func (t *Tree) Fill() *Tree {
	t.wg.Add(1)
	// subSem := make(chan struct{}, maxGoroutines)
	go t.walkFolder(&t.Folders /*, subSem*/)
	return t
}

// Is called internally by the walkFolder function.
// walkFolder walks through the folder and its subfolders.
// It creates a new goroutine for each subfolder.
// If the folder is in the exclusion list, it returns without processing the folder.
// If there is an error reading the directory, it adds an error to the folder.
func (t *Tree) walkFolder(folder *Folder /*, sem chan struct{}*/) {
	defer t.wg.Done()
	// sem <- struct{}{}
	// defer func() {
	// 	<-sem
	// }()

	// subSem := make(chan struct{}, maxGoroutines)
	// var subSem chan struct{}

	t.startExecution()
	defer t.stopExecution()

	if t.InExclusion(folder.Path) {
		return
	}

	entries, err := os.ReadDir(folder.Path)
	if err != nil {
		folder.Errors = append(folder.Errors, fmt.Sprintf("Failed to read directory %s: %v", folder.Path, err))
		return
	}

	for _, entry := range entries {
		entryPath := filepath.Join(folder.Path, entry.Name())
		if entry.IsDir() {
			subfolder := &Folder{Path: entryPath}
			folder.Subfolders = append(folder.Subfolders, subfolder)
			t.wg.Add(1)
			go t.walkFolder(subfolder /*, subSem*/)
		} else {
			folder.Files = append(folder.Files, entryPath)
		}
	}
}

// Prints the tree to a file defined in outFileName.
// It calls the printFolder function internally.
func (t *Tree) Print(outFileName string) *Tree {
	outputFile, err := os.Create(outFileName)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		return t
	}
	defer outputFile.Close()

	printFolder(&t.Folders, 0, outputFile)
	return t
}

// This function is called by the Print function internally.
func printFolder(folder *Folder, level int, outputFile *os.File) {
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}

	fmt.Fprintf(outputFile, "%s[%s]\n", indent, folder.Path)
	if len(folder.Errors) > 0 {
		for _, err := range folder.Errors {
			fmt.Fprintf(outputFile, "%s  ! Error: %s\n", indent, err)
		}
	}

	for _, file := range folder.Files {
		fmt.Fprintf(outputFile, "%s  - %s\n", indent, file)
	}

	for _, subfolder := range folder.Subfolders {
		printFolder(subfolder, level+1, outputFile)
	}
}

func (t *Tree) ShowDataAddresses() {
	fmt.Printf("tree: %p\n", t)
	fmt.Printf("tree.Folders: %p\n", &t.Folders)
	fmt.Printf("tree.Folders.Path: %p\n", &t.Folders.Path)
	fmt.Printf("tree.Folders.Files: %p\n", &t.Folders.Files)
	fmt.Printf("tree.Folders.Subfolders: %p\n", &t.Folders.Subfolders)
	fmt.Printf("tree.Folders.Errors: %p\n", &t.Folders.Errors)
	fmt.Printf("tree.wg: %p\n", &t.wg)
	fmt.Printf("tree.inexec: %p\n", &t.inexec)
	fmt.Printf("tree.startTime: %p\n", &t.startTime)
	fmt.Printf("tree.duration: %p\n", &t.duration)
	fmt.Printf("tree.Exclude: %p\n", &t.Exclude)
	fmt.Printf("tree.Executing: %p\n", &t.Executing)
}
