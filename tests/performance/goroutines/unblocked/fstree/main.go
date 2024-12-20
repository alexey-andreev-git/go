package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strings"

	"./internal/fstree"
)

var cmdLineFlags map[string]string = make(map[string]string)

var tree fstree.FsTree

func main() {
	addFlags()

	// debug.SetGCPercent(10000)
	// tmparray := make([]int, 0, 100*1024*1024)
	// for i := 0; i < 400; i++ {
	// 	tmparray = append(tmparray, i)
	// }
	// fmt.Println("tmparray len:", len(tmparray))

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// gcValue := m.NextGC
	// debug.SetGCPercent(-1)

	// Print memory stats
	printMemStat()

	perfExit := false
	go func() {
		fmt.Println("Press 'Enter' to exit...")
		fmt.Scanln()
		perfExit = true
	}()

	for !perfExit {
		tree = fstree.
			NewTree(cmdLineFlags["rootPath"]).
			SetExcludes(strings.Split(cmdLineFlags["excludes"], ",")).
			Start().Fill().Wait() //.Print("output.txt")
	}
	// prepareTree()

	// debug.SetGCPercent(int(gcValue))

	// tree.(*fstree.Tree).ShowDataAddresses()

	// Print memory stats
	printMemStat()

	fmt.Printf("Execution time: %v\n", tree.Duration())
	// time.Sleep(120 * time.Second)
	// fmt.Println("Press 'Enter' to exit...")
	// fmt.Scanln()
}

func prepareTree() {
	tree := fstree.
		NewTree(cmdLineFlags["rootPath"]).
		SetExcludes(strings.Split(cmdLineFlags["excludes"], ",")).
		Start().Fill().Wait()
	tree.ShowDataAddresses()
}

func addFlags() {
	// Make short and long flags for the exclude option
	rootPathShort := flag.String("r", "/", "Short flag for the root path")
	rootPathLong := flag.String("root", "/", "Root path to start the tree")
	excludesShort := flag.String("e", "/proc,/mnt", "Short flag for folders to exclude")
	excludesLong := flag.String("exclude", "/proc,/mnt", "Long flag for folders to exclude")
	flag.Parse()
	cmdLineFlags["excludes"] = *excludesShort
	if *excludesLong != "/proc,/mnt" {
		cmdLineFlags["excludes"] = *excludesLong
	}
	cmdLineFlags["rootPath"] = *rootPathShort
	if *rootPathLong != "/" {
		cmdLineFlags["rootPath"] = *rootPathLong
	}
}

func printMemStat() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
