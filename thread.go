package thread

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"./unlib"
)

var reg_id *regexp.Regexp = regexp.MustCompile(" ID:(........)")
var reg_from *regexp.Regexp = regexp.MustCompile(" </b>(.+)¥((.+)¥)<b>")
var reg_res *regexp.Regexp = regexp.MustCompile("&gt;&gt;([0-9]+)")

// レス構造体
type Res struct {
	Number		int
	Name		string
	From		string
	Mail		string
	Id			string
	Body		string
	Point		int
	Next		map[int]*Res
	Back		map[int]*Res
}

// スレッド構造体
type Thread struct {
	Name		string
	Path		string
	Ita			string
	Sure		string
	Length		int
	Point		int
	Reses		[]Res
	Ids			map[string][]int
}

func NewThread(base, ita, sure string) *Thread {
	this := new(Thread)
	this.Ita = ita
	this.Sure = sure
	this.Path = base + "/" + ita + "/" + sure[0:4] + "/" + sure
	return this
}

func (this *Thread) GetData() (bool, os.Error) {
	data, err := unlib.FileGetContents(this.Path)
	if err != nil { return false, unlib.Error("thread") }
	this.Ids = make(map[string][]int)
	list := strings.Split(string(data), "\n", -1)
	data = nil
	this.Length = len(list)
	this.Reses = make([]Res, this.Length)
	line := strings.Split(list[0], "<>", -1)
	this.Name = line[4]
	for key := range list {
		it := &(this.Reses[key])
		it.Number = key + 1
		if line = strings.Split(list[key], "<>", -1); len(line) > 3 {
			it.Name = line[0]
			it.Mail	= line[1]
			it.Body = line[3]
			this.fromSplit(it, line[0])
			this.idSplit(it, line[2])
			this.ankerSplit(it, line[3])
		}
	}
	for key := this.Length - 1; key >= 0; key-- {
		it := &(this.Reses[key])
		if it.Point == 0 && len(it.Next) > 0 {
			if p, ok := this.Ids[it.Id]; ok {
				num := len(p)
				if num > 3 {
					point_r(it, 10, 5)
				} else if num > 1 {
					point_r(it, 12, 3)
				} else {
					point_r(it, 15, 1)
				}
			} else {
				point_r(it, 15, 1)
			}
		}
		this.Point += it.Point
	}
	return true, nil
}

func (this *Thread) idSplit(res *Res, line string) (ret bool){
	ret = false
	if array := reg_id.FindStringSubmatch(line); len(array) > 1 {
		id := array[1]
		if m, ok := this.Ids[id]; ok {
			this.Ids[id] = append(m, res.Number)
		} else {
			this.Ids[id] = []int{res.Number}
		}
		res.Id = id
		ret = true
	}
	return
}

func (this *Thread) fromSplit(res *Res, line string) (ret bool){
	ret = false
	if from := reg_from.FindStringSubmatch(line); len(from) > 2 {
		res.From = from[2]
		ret = true
	}
	return
}

func (this *Thread) ankerSplit(res *Res, line string) (ret bool){
	ret = false
	if anker := reg_res.FindAllStringSubmatch(line, -1); len(anker) > 0 {
		res.Next = make(map[int]*Res)
		for _, item := range anker {
			resban, _ := strconv.Atoi(item[1])
			if resban < this.Length {
				anker := &(this.Reses[resban])
				res.Next[resban] = anker
				if anker.Back == nil {
					anker.Back = make(map[int]*Res)
				}
				anker.Back[res.Number] = res
			}
		}
		ret = true
	}
	return
}

func (this *Thread) Remove() {
	for key := range this.Reses {
		data := &(this.Reses[key])
		res_remove(&(data.Next))
		res_remove(&(data.Back))
		data.Next = nil
		data.Back = nil
		this.Reses[key] = Res{}
	}
	for key := range this.Ids {
		this.Ids[key] = nil, false
	}
	this.Reses = nil
	this.Ids = nil
}

func res_remove(m *(map[int]*Res)) {
	for k := range *m {
		it := (*m)[k]
		res_remove(&(it.Next))
		res_remove(&(it.Back))
		(*m)[k] = nil, false
	}
}

func point_r(res *Res, p, plus int){
	for _, it := range res.Next {
		if res.Id != it.Id && res.Number > it.Number {
			it.Point += p
			point_r(it, p + plus, plus + 1)
		}
	}
}

