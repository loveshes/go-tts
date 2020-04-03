package main

import (
	"fmt"
	"github.com/loveshes/go-tts/core"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const FORMAT = "http://nls-gateway.cn-shanghai.aliyuncs.com/stream/v1/tts?appkey=%s&token=%s&format=%s&voice=%s&volume=%s&speech_rate=%s&text="

func main() {
	basePath := ""
	var sc core.SoundConfig
	err := core.LoadIni("./conf.ini", &sc)
	if err != nil {
		fmt.Printf("初始化配置文件失败, 错误信息:%v\n", err)
	}
	baseurl := fmt.Sprintf(FORMAT, sc.Appkey, sc.Token, sc.Format, sc.Voice, sc.Volume, sc.SpeechRate)
	txtFiles := loadTxt()
	for _, txtName := range txtFiles {
		fmt.Printf("\r正在分析文件%s中     ", txtName)
		urls := core.LoadFile(basePath + txtName)
		fmt.Printf("\r分析文件%s完成       ", txtName)
		// 创建文件
		soundName := strings.Split(txtName, ".")[0] + "." + sc.Format
		file, err := os.OpenFile(soundName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}
		defer file.Close()
		total := float64(len(urls))
		for i, text := range urls {
			// 中文字符一定要转码
			text = url.QueryEscape(text)
			core.GetSoundAndSave(baseurl+text, file, float64(i), total)
		}
		fmt.Printf("\r合成文件%s成功.                      \n", soundName)
	}
}

// 得到当前路径下的txt文件
func loadTxt() []string {
	fileNames := make([]string, 0, 10)
	pwd, _ := os.Getwd()
	fileSlice, _ := filepath.Glob(filepath.Join(pwd, "*.txt"))
	for _, slice := range fileSlice {
		_, fileName := filepath.Split(slice)
		fileNames = append(fileNames, fileName)
	}
	return fileNames
}
