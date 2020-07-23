package control

import (
	"errors"
	"fmt"
	"html/template"
	_ "log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/iam1912/XIE_1/model"
)

var (
	templates = template.Must(template.ParseFiles("view/index.html", "view/edit.html"))
	validPath = regexp.MustCompile("^/(edit|index)/([a-zA-Z0-9]+)$")
	Stu       = model.NewStuSlice()
)

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[1], nil
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	renderTemplate(w, title, nil)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	if r.Method == "GET" {
		renderTemplate(w, title, nil)
	} else {
		r.ParseForm()
		switch r.Form.Get("change") {
		case "学生信息列表":
			ShowHandle(w, title)
		case "查询":
			SearchHandle(w, title, r)
		case "添加":
			AddHandle(w, title, r)
		case "修改":
			ModifyHandle(w, title, r)
		case "删除":
			DeleteHandle(w, title, r)
		case "学生信息排序":
			SortHandle(w, title)
		}
	}
}

func ShowHandle(w http.ResponseWriter, tmpl string) {
	stu := Stu.List()
	renderTemplate(w, tmpl, stu)
}

func SearchHandle(w http.ResponseWriter, tmpl string, r *http.Request) {
	if len(r.Form.Get("ID")) == 5 {
		getint, err := mathvaild(r.Form.Get("ID"))
		if err != nil {
			redirect(w, r, err)
		} else {
			if err := Stu.FindIndex(getint); err == nil {
				stu := Stu.Search(getint)
				renderTemplate(w, tmpl, stu)
			} else {
				redirect(w, r, err)
				return
			}
		}
	} else {
		err := errors.New("该数据长度不满足所要求的")
		redirect(w, r, err)
	}
}

func AddHandle(w http.ResponseWriter, tmpl string, r *http.Request) {
	var zh []string
	var math []int
	for index, val := range r.Form["add"] {
		if index == 0 || index == 3 || index == 5 {
			getint, err := mathvaild(val)
			if err != nil {
				redirect(w, r, err)
				return
			} else {
				math = append(math, getint)
			}
		} else if index == 1 || index == 2 {
			getstr, err := zhvaild(val)
			if err != nil {
				redirect(w, r, err)
				return
			} else {
				zh = append(zh, getstr)
			}
		} else {
			zh = append(zh, val)
		}
	}
	stu := model.NewStu(math[0], zh[0], zh[1], math[1],
		zh[2], math[2], zh[3])
	err := Stu.Add(stu)
	if err != nil {
		redirect(w, r, err)
	} else {
		redirect(w, r, nil)
		fmt.Println("添加成功")
	}
}

func ModifyHandle(w http.ResponseWriter, tmpl string, r *http.Request) {
	var zh []string
	var math []int
	getint, err := mathvaild(r.Form["mod"][0])
	if err != nil {
		redirect(w, r, err)
		return
	} else {
		err = Stu.FindIndex(getint)
		if err != nil {
			redirect(w, r, err)
			return
		} else {
			for index, val := range r.Form["mod"] {
				if index == 3 || index == 5 {
					getint, err := mathvaild(val)
					if err != nil {
						redirect(w, r, err)
						return
					} else {
						math = append(math, getint)
					}
				} else if index == 1 || index == 2 {
					getstr, err := zhvaild(val)
					if err != nil {
						redirect(w, r, err)
						return
					} else {
						zh = append(zh, getstr)
					}
				} else {
					zh = append(zh, val)
				}
			}
		}
	}
	err = Stu.Modify(getint, zh[1], zh[2], math[0], zh[3], math[1], zh[4])
	if err != nil {
		redirect(w, r, err)
	} else {
		redirect(w, r, nil)
	}
}

func DeleteHandle(w http.ResponseWriter, tmpl string, r *http.Request) {
	getint, err := mathvaild(r.Form.Get("IDD"))
	if err != nil {
		redirect(w, r, err)
	} else {
		err := Stu.Delete(getint)
		if err != nil {
			redirect(w, r, err)
		} else {
			fmt.Println("删除成功")
			redirect(w, r, nil)
		}
	}
}

func SortHandle(w http.ResponseWriter, tmpl string) {
	stu := Stu.Sort()
	renderTemplate(w, tmpl, stu)
}

func mathvaild(val string) (int, error) {
	getint, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	} else {
		return getint, nil
	}
}

func zhvaild(val string) (string, error) {
	m, err := regexp.MatchString("^\\p{Han}+$", val)
	if !m {
		return "", err
	} else {
		return val, nil
	}
}

func redirect(w http.ResponseWriter, r *http.Request, err error) {
	http.Redirect(w, r, r.URL.Path, http.StatusFound)
	if err != nil {
		fmt.Println(err)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, stu []model.Student) {
	err := templates.ExecuteTemplate(w, tmpl+".html", stu)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
