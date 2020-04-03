package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func GetSoundAndSave(urlPath string, file *os.File, i, total float64) {
	//通过url请求资源
	resp, err := http.Get(urlPath)
	if err != nil {
		fmt.Println("请求参数有误")
	}
	// 避免频繁请求
	time.Sleep(time.Millisecond * 300)
	defer resp.Body.Close()
	sound, err := ioutil.ReadAll(resp.Body)
	io.Copy(file, bytes.NewReader(sound))

	fmt.Printf("\r正在合成文件%s中...进度:%5.2f%%", file.Name(), i/total*100)
	return
}
