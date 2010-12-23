package unlib

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"container/vector"
)

type Error string
func (this Error) String() string {
	return string(this)
}

func FileGetContents(filename string) ([]byte, os.Error){
	if filename == "" { return nil, Error("ファイル名がありません") }
	fp, open_err := os.Open(filename, os.O_RDONLY, 0777)
	if open_err != nil {
		return nil, open_err
	}
	defer fp.Close()
	fileinfo, stat_err := fp.Stat()
	if stat_err != nil {
		return nil, stat_err
	}
	data := make([]byte, fileinfo.Size)
	if _, read_err := fp.Read(data); read_err != nil {
		return nil, read_err
	}
	return data, nil
}

func FileSize(filename string) (int64, os.Error) {
	if filename == "" { return -1, Error("ファイル名がありません") }
	fi, err := os.Stat(filename)
	if err != nil { return -1, err }
	return fi.Size, nil
}

func FileMTime(filename string) (int64, os.Error) {
	if filename == "" { return 0, Error("ファイル名がありません") }
	fi, err := os.Stat(filename)
	if err != nil { return 0, err }
	return fi.Mtime_ns, nil
}

// ファイルに書き込む
func FilePutContents(filename string, data []byte, flag bool) os.Error {
	if filename == "" { return Error("ファイル名がありません") }
	var bitflag int
	if flag {
		bitflag = os.O_WRONLY | os.O_CREAT
	} else {
		bitflag = os.O_WRONLY | os.O_APPEND
	}
	fp, err := os.Open(filename, bitflag, 0777)
	if err != nil { return err }
	defer fp.Close()
	fp.Write(data)
	return nil
}

// 文字列配列に検索文字列が格納されているか確認
func InArray(str string, list []string) bool {
	for _, data := range list {
		if str == data {
			return true
		}
	}
	return false
}

// ファイルの存在確認
func FileExists(filename string) bool {
	if filename == "" { return false }
	_, err := os.Stat(filename)
	if err != nil { return false }
	return true
}

// フォルダ生成
func MakeFolder(p string) os.Error {
	if p == "" { return nil }
	dir, _ := path.Split(p)
	mkdir_err := os.MkdirAll(dir, 0777)
	if mkdir_err != nil { return mkdir_err }
	return nil
}

func Qsort(list []interface{}, cmp func(a, b interface{}) int) (ret []interface{}, err os.Error) {
	if len(list) <= 0 {
		err = Error("len")
		return
	}
	ret = make([]interface{}, len(list))
	ret = list
	stack := new(vector.IntVector)
	stack.Push(0)
	stack.Push(len(list) - 1)
	for len(*stack) != 0 {
		tail := stack.Pop()
		head := stack.Pop()
		pivot := ret[head + ((tail - head) >> 1)]
		i := head - 1
		j := tail + 1
		for {
			for i++; cmp(ret[i], pivot) < 0; i++ {}
			for j--; cmp(ret[j], pivot) > 0; j-- {}
			if i >= j { break }
			tmp := ret[i]
			ret[i] = ret[j]
			ret[j] = tmp
		}
		if head < (i - 1) {
			stack.Push(head)
			stack.Push(i - 1)
		}
		if (j + 1) < tail {
			stack.Push(j + 1)
			stack.Push(tail)
		}
	}
	return
}

func MemStatsPrint(){
	mem := runtime.MemStats
	fmt.Printf("Alloc\t\t:\t%d\n", mem.Alloc)
	fmt.Printf("TotalAlloc\t:\t%d\n", mem.TotalAlloc)
	fmt.Printf("Sys\t\t:\t%d\n", mem.Sys)
	fmt.Printf("Lookups\t\t:\t%d\n", mem.Lookups)
	fmt.Printf("Mallocs\t\t:\t%d\n", mem.Mallocs)
	fmt.Printf("HeapAlloc\t:\t%d\n", mem.HeapAlloc)
	fmt.Printf("HeapSys\t\t:\t%d\n", mem.HeapSys)
	fmt.Printf("HeapIdle\t:\t%d\n", mem.HeapIdle)
	fmt.Printf("HeapObjects\t:\t%d\n", mem.HeapObjects)
	fmt.Printf("StackInuse\t:\t%d\n", mem.StackInuse)
	fmt.Printf("StackSys\t:\t%d\n", mem.StackSys)
	fmt.Printf("MSpanInuse\t:\t%d\n", mem.MSpanInuse)
	fmt.Printf("MSpanSys\t:\t%d\n", mem.MSpanSys)
	fmt.Printf("MCacheInuse\t:\t%d\n", mem.MCacheInuse)
	fmt.Printf("MCacheSys\t:\t%d\n", mem.MCacheSys)
	fmt.Printf("MHeapMapSys\t:\t%d\n", mem.MHeapMapSys)
	fmt.Printf("BuckHashSys\t:\t%d\n", mem.BuckHashSys)
	fmt.Printf("NextGC\t\t:\t%d\n", mem.NextGC)
	fmt.Printf("PauseNs\t\t:\t%d\n", mem.PauseNs)
	fmt.Printf("NumGC\t\t:\t%d\n", mem.NumGC)
	println(mem.EnableGC)
	println(mem.DebugGC)
}

