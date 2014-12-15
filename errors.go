package golidera

import (
	"errors"
	"fmt"
	"strings"
)

func GetErrorMessage(d string, t ...string) string {
	if len(t) == 0 {
		return d
	}
	return t[0]
}

func RenderError(t string, data map[string]interface{}) error {
	for key, value := range data {
		t = strings.Replace(t, "{:"+key+"}", fmt.Sprint(value), -1)
	}
	return errors.New(t)
}

type Translator interface {
	Translate(string, string) string
	TranslateOk(string, string) (string, bool)
}

type TranslatorMap map[string]map[string]string

func (this TranslatorMap) Translate(locale string, msg string) string {
	m, _ := this.TranslateOk(locale, msg)
	return m
}

func (this TranslatorMap) TranslateOk(locale string, msg string) (string, bool) {
	m, ok := this[locale][msg]
	return m, ok
}

func T(fn CheckerFunc, s string, t Translator) CheckerFunc {
	return CheckerFunc(func(g GolideraFielder, fi Fielder, f Former, locale ...string) (bool, error) {
		if len(locale) == 0 {
			return fn.Check(g, fi, f)
		}
		m, ok := t.TranslateOk(locale[0], s)
		if ok {
			return fn.Check(g, fi, f, m)
		}
		return fn.Check(g, fi, f)
	})
}
