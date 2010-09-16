package thread

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ThreadError string
func (this ThreadError) String() string {
	return string(this)
}

type Item struct {
	Id			string
	Match		string
}

type Match struct {
	Reg			regexp.Regexp
	Items		[]Item
}

// レス構造体
type Res struct {
	Number		int
	Name		string
	From		string
	Id			string
	Body		string
	Point		int
	Next		[]*Res
	Back		[]*Res
}

type Thread struct {
	Name		string
	Path		string
	Saba		string
	Ita			string
	Sure		string
	Point		int
	Reses		[]Res
}

func NewThread(base, saba, ita, sure string) *Thread {
	this := new(Thread)
	path := base + "/" + saba + "/" + ita + "/" + sure[0:4] + "/" + sure
	this.Saba = saba
	this.Ita = ita
	this.Sure = sure
	this.Path = path
	return this
}

func (this *Thread) GetData() (bool, os.Error) {
	reg_id, id_err := regexp.Compile(" ID:(........)")
	if id_err != nil { return false, ThreadError("id") }
	reg_from, from_err := regexp.Compile(" </b>(.+)¥((.+)¥)<b>")
	if from_err != nil { return false, ThreadError("from") }
	reg_res, res_err := regexp.Compile("&gt;&gt;([0-9]+)")
	if res_err != nil { return false, ThreadError("res") }
	data, err := fileGetContents(this.Path)
	if err != nil { return false, ThreadError("thread") }
	list := strings.Split(string(data), "\n", -1)
	list_length := len(list)
	this.Reses = make([]Res, list_length)
	line := strings.Split(list[0], "<>", -1)
	this.Name = line[4]
	for key := range list {
		it := &(this.Reses[key])
		it.Number = key + 1
		if line = strings.Split(list[key], "<>", -1); len(line) > 3 {
			it.Name = line[0]
			it.Body = line[3]
			if id := reg_id.FindStringSubmatch(line[2]); len(id) > 1 {
				it.Id = id[1]
			}
			if from := reg_from.FindStringSubmatch(line[0]); len(from) > 2 {
				it.From = from[2]
			}
			if res := reg_res.FindAllStringSubmatch(line[3], -1); len(res) > 0 {
				next := make([]*Res, len(res))
				i := 0
				for _, item := range res {
					resban, _ := strconv.Atoi(item[1])
					if resban < list_length {
						next[i] = &(this.Reses[resban])
						i++
					}
				}
				it.Next = next[0:i]
			}
		}
	}
	return true, nil
}

func (this *Thread) GetPoint() (bool, os.Error){
	point := 0
	key := len(this.Reses)
	for key--; key >= 0; key-- {
		it := &(this.Reses[key])
		if it.Point == 0 && len(it.Next) > 0 {
			point_r(it, 10)
		}
		point += it.Point
	}
	this.Point = point
	return true, nil
}

func point_r(res *Res, p int){
	for _, it := range res.Next {
		if res.Id != it.Id && res.Number > it.Number {
			it.Point += p
			point_r(it, p + 3)
		}
	}
}

func fileGetContents(filename string) ([]byte, os.Error){
	fp, open_err := os.Open(filename, os.O_RDONLY, 0777)
	if open_err != nil {
		return nil, ThreadError("open")
	}
	defer fp.Close()
	fileinfo, stat_err := fp.Stat()
	if stat_err != nil {
		return nil, ThreadError("stat")
	}
	data := make([]byte, fileinfo.Size)
  	if _, read_err := fp.Read(data); read_err != nil {
		return nil, ThreadError("read")
	}
	return data, nil
}

