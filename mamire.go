package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"container/vector"
	"./thread"
)

type Error string
func (this Error) String() string {
	return string(this)
}

type Board struct {
	Name		string
	Ita			string
	Saba		string
}

type MiniThread struct {
	Name		string
	Saba		string
	Ita			string
	Sure		string
	Point		int
}
func NewMiniThread(t *thread.Thread) *MiniThread {
	this := new(MiniThread)
	this.Name = t.Name
	this.Saba = t.Saba
	this.Ita = t.Ita
	this.Sure = t.Sure
	this.Point = t.Point
	return this
}

const g_base_path string = "/2ch/dat"
const g_output_path string = "/2ch/dat/2chpoint.tsv"
const g_board_list_path string = "/2ch/getboard.data"
const g_ita_data_path string = "/2ch/dat/ita.data"
const g_thread_list string = "subject.txt"

func main(){
	sl := serverList()
	bl := boardList(sl)
	tl := threadList(bl)
	tl = qsort(tl)
	fp, open_err := os.Open(g_output_path, os.O_WRONLY | os.O_CREAT, 0777)
	if open_err != nil { panic("g_output_path") }
	defer fp.Close()
	for _, it := range tl {
		dot := strings.Index(it.Sure, ".")
		if dot > 0 {
			fmt.Fprintf(fp, "%d\t%s\t%s\t%s\t%s\n", it.Point, it.Name, it.Saba, it.Ita, it.Sure[0:dot])
		}
	}
}

func boardList(sl map[string]Board) ([]Board){
	data, open_err := fileGetContents(g_board_list_path)
	if open_err != nil { panic("g_board_list_path") }
	bl := strings.Split(string(data), "\n", -1)
	list := make([]Board, len(bl))
	i := 0
	for _, it := range bl {
		if board, ok := sl[it]; ok {
			list[i] = board
			i++
		}
	}
	return list[0:i]
}

func serverList() (map[string]Board) {
	var line Board
	list := make(map[string]Board, 1000)
	data, open_err := fileGetContents(g_ita_data_path)
	if open_err != nil { panic("g_ita_data_path") }
	reg, reg_err := regexp.Compile("(.+)/(.+)<>(.+)")
	if reg_err != nil { panic("reg err") }
	sl := strings.Split(string(data), "\n", -1)
	for _, it := range sl {
		if match := reg.FindStringSubmatch(it); len(match) > 2 {
			line.Name = match[3]
			line.Ita = match[2]
			line.Saba = match[1]
			list[line.Ita] = line
		}
	}
	return list
}

func threadList(bl []Board) ([]*MiniThread) {
	tlist := new(vector.Vector)
	for _, it := range bl {
		base_path := g_base_path + "/" + it.Saba + "/" + it.Ita
		b_path := base_path + "/" + g_thread_list
		data, open_err := fileGetContents(b_path)
		if open_err != nil { continue }
		list := strings.Split(string(data), "\n", -1)
		for _, line := range list {
			array := strings.Split(line, "<>", -1)
			if len(array) > 1 {
				t := thread.NewThread(g_base_path, it.Saba, it.Ita, array[0])
				if ok, _ := t.GetData(); ok {
					t.GetPoint()
					tlist.Push(NewMiniThread(t))
				}
				t = nil
			}
		}
	}
	tl := make([]*MiniThread, len(*tlist))
	for i, it := range *tlist {
		tl[i] = it.(*MiniThread)
	}
	return tl
}

func fileGetContents(filename string) ([]byte, os.Error){
	fp, open_err := os.Open(filename, os.O_RDONLY, 0777)
	if open_err != nil {
		return nil, Error("open")
	}
	defer fp.Close()
	fileinfo, stat_err := fp.Stat()
	if stat_err != nil {
		return nil, Error("stat")
	}
	data := make([]byte, fileinfo.Size)
	if _, read_err := fp.Read(data); read_err != nil {
		return nil, Error("read")
	}
	return data, nil
}

func qsort(list []*MiniThread) []*MiniThread {
	cmp := func(a, b *MiniThread) int {
		return b.Point - a.Point
	}
	ret := make([]*MiniThread, len(list))
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
			tmp := ret[i];
			ret[i] = ret[j];
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
	length := len(ret)
	if length > 100 {
		length = 100
	}
	return ret[0:length]
}
