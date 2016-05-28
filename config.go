// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	i "github.com/caixw/apidoc/input"
)

type config struct {
	Version string  `json:"version"` // 兼容的 apidoc 版本
	Input   *input  `json:"input"`
	Output  *output `json:"output"`
	Doc     *doc    `json:"doc"`
}

type input struct {
	Type      string   `json:"type"`      // 输入的目标语言
	Dir       string   `json:"dir"`       // 源代码目录
	Exts      []string `json:"exts"`      // 需要扫描的文件扩展名
	Recursive bool     `json:"recursive"` // 是否查找Dir的子目录
}

type output struct {
	Dir string `json:"dir"`
	//Type string   `json:"type"` // 输出的语言格式
	//Groups     []string `json:"groups"`     // 需要打印的分组内容。
	//Timezone   string   `json:"timezone"`   // 时区
}

type doc struct {
	Version string `json:"version"` // 文档版本号
	Title   string `json:"title"`   // 文档的标题，默认为apidoc
	BaseURL string `json:"baseURL"` // api文档中url的前缀，不指定，则为空
}

// 从配置文件中加载配置项。
func loadConfig() (*config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(wd + "/" + configFilename)
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if err = initInput(wd, cfg); err != nil {
		return nil, err
	}
	if err = initDoc(cfg); err != nil {
		return nil, err
	}

	if err = initOutput(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// 对config.Output中的变量做初始化
func initOutput(cfg *config) error {
	cfg.Output.Dir += string(os.PathSeparator)
	return nil
}

// 对config.Doc中的变量做初始化
func initDoc(cfg *config) error {
	if len(cfg.Doc.Title) == 0 {
		cfg.Doc.Title = "APIDOC"
	}

	return nil
}

// 对config.Input中的变量做初始化
func initInput(wd string, cfg *config) error {
	if len(cfg.Input.Dir) == 0 {
		cfg.Input.Dir = wd
	}
	cfg.Input.Dir += string(os.PathSeparator)

	if len(cfg.Input.Exts) > 0 {
		exts := make([]string, 0, len(cfg.Input.Exts))
		for _, ext := range cfg.Input.Exts {
			if len(ext) == 0 {
				continue
			}

			if ext[0] != '.' {
				ext = "." + ext
			}
			exts = append(exts, ext)
		}
		cfg.Input.Exts = exts
	}

	// 若没有指定Type，则根据exts和当前目录下的文件检测来确定其值
	if len(cfg.Input.Type) == 0 {
		return fmt.Errorf("必须指定参数 type")
	}
	cfg.Input.Type = strings.ToLower(cfg.Input.Type)

	if i.LangIsSupported(cfg.Input.Type) {
		return fmt.Errorf("暂不支持该类型[%v]的语言", cfg.Input.Type)
	}

	return nil
}

// 在当前目录下产生个默认的配置文件。
func genConfigFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := wd + string(os.PathSeparator) + configFilename
	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	cfg := &config{
		Input:  &input{Dir: "./", Recursive: true},
		Output: &output{},
		Doc:    &doc{},
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	_, err = fi.Write(data)
	return err
}
