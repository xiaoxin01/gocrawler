package model

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/gocolly/colly"
	"github.com/liuzl/gocc"
)

var (
	t2sService *gocc.OpenCC
)

// Field field to add to each item
type Field struct {
	Operator  string
	Parameter string
	Selector  string
	Regexp    *RegexOperation
	Sprintf   *string
	Action    *string
}

// RegexOperation regexp to change field value
type RegexOperation struct {
	Expression string
	Group      int
}

// InitField init field
func InitField(gccPath string) {
	var err error
	*gocc.Dir = `./module/gocc`
	if gccPath != "" {
		*gocc.Dir = fmt.Sprintf("%s%s", gccPath, *gocc.Dir)
	}
	t2sService, err = gocc.New("t2s")
	if err != nil {
		panic(err)
	}
}

// GetValue get field value
func (field Field) GetValue(e *colly.HTMLElement) (v interface{}, ok bool) {
	ok = true
	switch field.Operator {
	case "Attr":
		v = e.ChildAttr(field.Selector, field.Parameter)
	case "Attrs":
		if attr := e.ChildAttrs(field.Selector, field.Parameter); attr != nil {
			v = attr
		} else {
			ok = false
		}
	case "Text":
		v = e.ChildText(field.Selector)
	case "Const":
		v = field.Parameter
	case "Func":
		v, ok = field.getFuncValue(field.Parameter)
	default:
		ok = false
	}

	if ok && field.Regexp != nil {
		values := regexp.MustCompile(field.Regexp.Expression).FindStringSubmatch(v.(string))
		if len(values) > field.Regexp.Group {
			v = values[field.Regexp.Group]
		}
	}

	if ok && field.Sprintf != nil {
		v = fmt.Sprintf(*field.Sprintf, v)
	}

	if ok && field.Action != nil && *field.Action == "t2s" && v != "" {
		v, _ = t2sService.Convert(fmt.Sprintf("%v", v))
	}

	ok = ok && v != ""

	defer func() {
		err := recover()
		if err != nil {
			log.Print(err)
		}
	}()

	return
}

func (field Field) getFuncValue(funcName string) (v interface{}, ok bool) {
	ok = true

	switch funcName {
	case "time.Now().Unix()":
		v = time.Now().Unix()
	default:
		ok = false
	}

	return
}
