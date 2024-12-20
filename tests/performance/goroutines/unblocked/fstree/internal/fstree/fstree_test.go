package fstree

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewTree(t *testing.T) {
	type args struct {
		rootPath string
	}
	tests := []struct {
		name string
		args args
		want *Tree
	}{
		{
			name: "Root path is set correctly",
			args: args{rootPath: "/test/root"},
			want: &Tree{Folders: Folder{Path: "/test/root"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTree(tt.args.rootPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_AddExclude(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	type args struct {
		exclude string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Tree
	}{
		{
			name: "Exclude is added correctly",
			fields: fields{
				Exclude: []string{"/test/exclude"},
			},
			args: args{exclude: "/test/exclude2"},
			want: &Tree{Exclude: []string{"/test/exclude", "/test/exclude2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			if got := tr.AddExclude(tt.args.exclude); !reflect.DeepEqual(got.Exclude, tt.want.Exclude) {
				t.Errorf("Tree.AddExclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_SetExcludes(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	type args struct {
		excludes []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Tree
	}{
		{
			name: "Excludes are set correctly",
			fields: fields{
				Exclude: []string{"/test/exclude"},
			},
			args: args{excludes: []string{"/test/exclude2", "/test/exclude3"}},
			want: &Tree{Exclude: []string{"/test/exclude", "/test/exclude2", "/test/exclude3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			if got := tr.SetExcludes(tt.args.excludes); !reflect.DeepEqual(got.Exclude, tt.want.Exclude) {
				t.Errorf("Tree.SetExcludes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_GetExcludes(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Returns the list of folders to exclude",
			fields: fields{
				Exclude: []string{"/test/exclude", "/test/exclude2"},
			},
			want: []string{"/test/exclude", "/test/exclude2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			if got := tr.GetExcludes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tree.GetExcludes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_Wait(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	tests := []struct {
		name   string
		fields fields
		want   *Tree
	}{
		{
			name: "Wait for all goroutines to finish",
			fields: fields{
				startTime: time.Now(),
			},
			want: &Tree{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			tr.Start()
			tr.wg.Add(1)
			go func() {
				defer tr.wg.Done()
				time.Sleep(2 * time.Second)
			}()
			if got := tr.Wait(); got.startTime.Sub(tt.want.startTime) < (2 * time.Second) {
				t.Errorf("Tree.Wait() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_Start(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	tests := []struct {
		name   string
		fields fields
		want   *Tree
	}{
		{
			name: "Defines start time the execution of the tree.",
			fields: fields{
				startTime: time.Now(),
			},
			want: &Tree{startTime: time.Now()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			// Compare time.Time values difference with a tolerance of 1 second
			if got := tr.Start(); got.startTime.Sub(tt.want.startTime) > time.Second {
				t.Errorf("Tree.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_Duration(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "Returns the duration of the execution.",
			fields: fields{
				startTime: time.Now(),
			},
			want: time.Since(time.Now()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			// Compare time.Duration values difference with a tolerance of 1 second
			if got := tr.Duration(); got-tt.want > time.Second {
				t.Errorf("Tree.Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_startExecution(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "This function is called when a goroutine starts.",
			fields: fields{
				Executing: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			tr.startExecution()
			if !tr.Executing || tr.inexec.Load() != 1 {
				t.Errorf("Tree.startExecution() = %v, want %v", tr.Executing, true)
			}
		})
	}
}

func TestTree_stopExecution(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "This function is called when a goroutine finishes.",
			fields: fields{
				Executing: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			tr.inexec.Add(1)
			tr.Start()
			tr.stopExecution()
			if tr.Executing {
				t.Errorf("Tree.stopExecution() = %v, want %v", tr.Executing, false)
			}
			if tr.inexec.Load() != 0 {
				t.Errorf("Tree.stopExecution() = %v, want %v", tr.inexec.Load(), 0)
			}
			if tr.duration == 0 {
				t.Errorf("Tree.stopExecution() = %v, want %v", tr.duration, "greater than 0")
			}
		})
	}
}

func TestTree_IsExecuting(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "IsExecuting returns true if the tree is being executed",
			fields: fields{
				Executing: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			if got := tr.IsExecuting(); got != tt.want {
				t.Errorf("Tree.IsExecuting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_InExclusion(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Path is in the exclusion list",
			fields: fields{
				Exclude: []string{"/test/exclude"},
			},
			args: args{path: "/test/exclude"},
			want: true,
		},
		{
			name: "Path is not in the exclusion list",
			fields: fields{
				Exclude: []string{"/test/exclude"},
			},
			args: args{path: "/test/exclude2"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			if got := tr.InExclusion(tt.args.path); got != tt.want {
				t.Errorf("Tree.InExclusion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_Fill(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	tests := []struct {
		name   string
		fields fields
		want   *Tree
	}{
		{
			name: "Start process of filling the tree with the files and folders",
			fields: fields{
				Folders: Folder{Path: "/test/folder"},
			},
			want: &Tree{Folders: Folder{Path: "/test/folder"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			got := tr.Fill().Wait()
			if !reflect.DeepEqual(got.Folders.Path, tt.want.Folders.Path) {
				t.Errorf("Tree.Fill() = %v, want %v", got.Folders.Path, tt.want.Folders.Path)
			}
			if len(got.Folders.Errors) == 0 {
				t.Errorf("Tree.Fill() = %v, want %v", got.Folders.Errors, tt.want.Folders.Errors)
			}
		})
	}
}

func TestTree_walkFolder(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	type args struct {
		folder *Folder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Walk through the folder",
			fields: fields{
				Folders: Folder{Path: "../"},
				Exclude: []string{},
			},
			args: args{folder: &Folder{Path: "../"}},
		},
		{
			name: "Walk through the folder in exclusion list",
			fields: fields{
				Folders: Folder{Path: "/test/folder"},
				Exclude: []string{"/test/folder"},
			},
			args: args{folder: &Folder{Path: "/test/folder"}},
		},
		{
			name: "Walk through the folder with errors",
			fields: fields{
				Folders: Folder{Path: "/test/folder"},
				Exclude: []string{},
			},
			args: args{folder: &Folder{Path: "/test/folder"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			tr.wg.Add(1)
			tr.Start()
			go tr.walkFolder(tt.args.folder /*, make(chan struct{}, maxGoroutines)*/)
			tr.Wait()
			if len(tt.args.folder.Errors) == 0 && len(tt.args.folder.Subfolders) == 0 && len(tt.fields.Exclude) == 0 {
				t.Errorf("Tree.walkFolder() = %v, want %v", tt.args.folder.Errors, "files count greater than 0")
			}
			if len(tt.args.folder.Errors) == 0 && len(tt.fields.Exclude) == 0 && len(tt.args.folder.Subfolders) == 0 {
				t.Errorf("Tree.walkFolder() = %v, want %v", tt.args.folder.Errors, "errors count greater than 0")
			}
			if len(tt.args.folder.Errors) > 0 && len(tt.fields.Exclude) != 0 {
				t.Errorf("Tree.walkFolder() = %v, want %v", tt.args.folder.Errors, "exclude folder")
			}
		})
	}
}

func TestTree_Print(t *testing.T) {
	type fields struct {
		startTime time.Time
		duration  time.Duration
		Executing bool
		Exclude   []string
		Folders   Folder
	}
	type args struct {
		outFileName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Tree
	}{
		{
			name: "Print the tree to a file",
			fields: fields{
				Folders: Folder{Path: "/test/folder"},
			},
			args: args{outFileName: "/tmp/test_output.txt"},
			want: &Tree{Folders: Folder{Path: "/test/folder"}},
		},
		{
			name: "Print the tree to a file with errors",
			fields: fields{
				Folders: Folder{Path: "/test/folder", Errors: []string{"Error 1", "Error 2"}},
			},
			args: args{outFileName: "/dev/null/test_output.txt"},
			want: &Tree{Folders: Folder{Path: "/test/folder", Errors: []string{"Error 1", "Error 2"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tree{
				startTime: tt.fields.startTime,
				duration:  tt.fields.duration,
				Executing: tt.fields.Executing,
				Exclude:   tt.fields.Exclude,
				Folders:   tt.fields.Folders,
			}
			os.Remove(tt.args.outFileName)
			got := tr.Print(tt.args.outFileName)
			if !reflect.DeepEqual(got.Folders, tt.want.Folders) {
				t.Errorf("Tree.Print() = %v, want %v", got, tt.want)
			}
			if _, err := os.Stat(tt.args.outFileName); !strings.Contains(tt.args.outFileName, "/dev/null") && os.IsNotExist(err) {
				t.Errorf("Tree.Print() = %v, want %v", tt.args.outFileName, "file exists")
			}
			os.Remove(tt.args.outFileName)
		})
	}
}

func Test_printFolder(t *testing.T) {
	type args struct {
		folder     *Folder
		level      int
		outputFile string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Print the folder to the output file",
			args: args{
				folder: &Folder{
					Path:   "/test/folder",
					Errors: []string{"Error 1", "Error 2"},
					Files:  []string{"/test/file1", "/test/file2"},
					Subfolders: []*Folder{
						{Path: "/test/folder1", Files: []string{"/test/file3"}},
						{Path: "/test/folder2", Files: []string{"/test/file4"}},
					}},
				level:      0,
				outputFile: "test_output.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outFile, err := os.CreateTemp("", tt.args.outputFile)
			if err != nil {
				t.Errorf("Failed to create output file: %v\n", err)
				return
			}
			defer outFile.Close()
			printFolder(tt.args.folder, tt.args.level, outFile)
			if _, err := os.Stat(outFile.Name()); os.IsNotExist(err) {
				t.Errorf("printFolder() = %v, want %v", tt.args.outputFile, "file exists")
			}
			os.Remove(tt.args.outputFile)
		})
	}
}
