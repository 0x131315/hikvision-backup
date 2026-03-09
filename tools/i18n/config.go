package main

import (
	"path/filepath"
	"regexp"
)

type language struct {
	Code       string
	TargetLang string
	GoogleLang string
}

type languageNavItem struct {
	Code string
	Name string
}

type tmFile struct {
	SourceHash string            `json:"source_hash"`
	UpdatedAt  string            `json:"updated_at"`
	Blocks     map[string]string `json:"blocks"`
}

var (
	sourceFile = "README.md"
	i18nDir    = "i18n"
	tmDir      = filepath.Join(i18nDir, "tm")
	langs      = []language{
		{Code: "ru", TargetLang: "RU", GoogleLang: "ru"},
		{Code: "zh", TargetLang: "ZH", GoogleLang: "zh-CN"},
	}
	languageNav = []languageNavItem{
		{Code: "en", Name: "English"},
		{Code: "ru", Name: "Русский"},
		{Code: "zh", Name: "中文"},
	}
	sepRe = regexp.MustCompile(`\n{2,}`)
)
