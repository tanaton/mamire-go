package main

import (
	"fmt"
	"os"
	"bufio"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"./unlib"
	"./thread"
)

type Board struct {
	Name		string
	Ita			string
}

type MiniThread struct {
	Name		string
	Ita			string
	Sure		string
	Point		int
}
func NewMiniThread(t *thread.Thread) (this MiniThread) {
	this.Name = t.Name
	this.Ita = t.Ita
	this.Sure = t.Sure
	this.Point = t.Point
	return
}

const g_base_path string = "/2ch/dat"
const g_output_path string = "/2ch/dat/2chpoint.tsv"
const g_board_list_path string = "/2ch/getboard.data"
const g_ita_data_path string = "/2ch/dat/ita.data"
const g_thread_list string = "subject.txt"

var LF_BYTE []byte = []byte{'\n'}
var HTML_DELI_BYTE []byte = []byte{'<', '>'}

func main() {
	cpu := 0
	if 1 < len(os.Args) {
		if i, err := strconv.Atoi(os.Args[1]); (err == nil) && (i > 0) {
			cpu = i
			runtime.GOMAXPROCS(cpu)
		}
	} else {
		cpu = 1
	}
	sl := serverList()
	bl := boardList(sl)
	tl := threadList(bl, cpu)
	pl, qsort_err := unlib.Qsort(tl, cmp)
	if qsort_err != nil {
		panic("qsort")
	}
	fp, open_err := os.Open(g_output_path, os.O_WRONLY | os.O_CREAT, 0777)
	if open_err != nil { panic("g_output_path") }
	defer fp.Close()
	bfp := bufio.NewWriter(fp)
	for _, p := range pl[0:5] {
		it := p.(MiniThread)
		dot := strings.Index(it.Sure, ".")
		if dot > 0 {
			bfp.WriteString(fmt.Sprintf("%d\t%s\t%s\t%s\n", it.Point, it.Name, it.Ita, it.Sure[0:dot]))
		}
	}
	bfp.Flush()
}

/*
func boardList(sl map[string]Board) ([]Board) {
	data, open_err := unlib.FileGetContents(g_board_list_path)
	if open_err != nil { panic("g_board_list_path") }
	bl := strings.Split(string(data), "\n", -1)
	list := make([]Board, 0, len(bl))
	for _, it := range bl {
		if board, ok := sl[it]; ok {
			list = append(list, board)
		}
	}
	return list
}
*/
func boardList(sl map[string]Board) ([]Board) {
	list := make([]Board, 0, len(sl))
	for _, it := range sl {
		list = append(list, it)
	}
	return list
}

func serverList() (map[string]Board) {
	var line Board
	list := make(map[string]Board, 1000)
	data, open_err := unlib.FileGetContents(g_ita_data_path)
	if open_err != nil { panic("g_ita_data_path") }
	reg, reg_err := regexp.Compile("(.+)/(.+)<>(.+)")
	if reg_err != nil { panic("reg err") }
	sl := strings.Split(string(data), "\n", -1)
	for _, it := range sl {
		if match := reg.FindStringSubmatch(it); len(match) > 2 {
			line.Name = match[3]
			line.Ita = match[2]
			list[line.Ita] = line
		}
	}
	return list
}

func threadList(bl []Board, cpu int) ([]interface{}) {
	tlist := make([]interface{}, 0, 400000)
	ch := make(chan MiniThread, cpu * 16)
	sync := make(chan bool, cpu)
	go func(){
		for {
			if data := <- ch; data.Sure != "" {
				tlist = append(tlist, data)
			} else {
				break
			}
		}
	}()
	for _, it := range bl {
		sync <- true
		runtime.GC()
		go func(){
			threadThread(it, ch)
			<- sync
		}()
	}
	for cpu > 0 {
		sync <- true
		runtime.GC()
		cpu--
	}
	close(ch)
	close(sync)
	return tlist
}

func threadThread(it Board, ch chan MiniThread) {
	base_path := g_base_path + "/" + it.Ita
	b_path := base_path + "/" + g_thread_list
	data, open_err := unlib.FileGetContents(b_path)
	if open_err != nil { return }
	list := strings.Split(string(data), "\n", -1)
	for _, line := range list {
		array := strings.Split(line, "<>", -1)
		if len(array) > 1 {
			t := thread.NewThread(g_base_path, it.Ita, array[0])
			if ok, _ := t.GetData(); ok {
				ch <- NewMiniThread(t)
			}
			t.Remove()
		}
	}
}

func cmp(a, b interface{}) int {
	aa := a.(MiniThread)
	bb := b.(MiniThread)
	return bb.Point - aa.Point
}

