package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type SoundConfig struct {
	Appkey     string `ini:"appkey"`
	Token      string `ini:"token"`
	Format     string `ini:"format"`
	Voice      string `ini:"voice"`
	SpeechRate string `ini:"speech_rate"`
	Volume     string `ini:"volume"`
}

func LoadIni(fileName string, data interface{}) (err error) {
	// 校验data参数，必须为结构体指针类型
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		err = errors.New("data不是一个指针类型")
		return
	}
	if t.Elem().Kind() != reflect.Struct {
		err = errors.New("data不是一个结构体指针类型")
		return
	}
	// 1. 读文件得到字节类型的数据
	fileStr, err := ioutil.ReadFile(fileName)
	if err != nil {
		err = errors.New("打开文件失败，请检查同级目录下是否有conf.in配置文件")
		return
	}
	// 2. 一行一行得到数据
	lineSlice := strings.Split(string(fileStr), "\n")
	for i, line := range lineSlice {
		// 去除每行首尾的空格
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		// 2.1 如果是注释就跳过
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		// 2.2 如果是[开头的就是节
		var section string
		if strings.HasPrefix(line, "[") {
			right := len(line) - 1
			if line[right] != ']' {
				err = fmt.Errorf("配置文件第%d行, [语法错误] 没有节名称", i+1)
				return
			}
			// 获得节名称
			section = strings.TrimSpace(line[1:right])
			if len(section) == 0 {
				err = fmt.Errorf("配置文件第%d行, [语法错误] 没有节名称", i+1)
				return
			}
		} else {
			// 2.3 如果不是[开头并且有=号分隔的键值对
			split := strings.Split(line, "=")
			if len(split) != 2 {
				err = fmt.Errorf("配置文件第%d行, [语法错误] 遗失等号或出现多个等号", i+1)
				return
			}
			// 函数入参是结构体指针，故这里得到的值是个结构体指针的值，需要通过.Elem()来获取对应的结构体值
			sValue := reflect.ValueOf(data).Elem() // 结构体的值
			sType := sValue.Type()                 // 结构体的类型

			if sValue.Kind() != reflect.Struct {
				err = fmt.Errorf("内部运行出错")
				return
			}
			// 根据key去结构体中找对应的字段名
			key := strings.TrimSpace(split[0]) // 只是配置文件中的关键字
			value := strings.TrimSpace(split[1])
			var fieldName string // 真正的字段名称
			var fieldType reflect.StructField
			// 遍历字段，判断tag是否等于key
			for i := 0; i < sValue.NumField(); i++ {
				field := sType.Field(i)
				fieldType = field
				if field.Tag.Get("ini") == key {
					fieldName = field.Name // 得到真正的字段名称
					break
				}
			}
			// 取出该字段
			fieldObj := sValue.FieldByName(fieldName)
			// 对其赋值
			switch fieldType.Type.Kind() {
			case reflect.String:
				fieldObj.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					err = fmt.Errorf("配置文件第%d行, [数值错误] 不为整数", i+1)
				}
				fieldObj.SetInt(intValue)
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					err = fmt.Errorf("配置文件第%d行, [数值错误] 不为正确的布尔值", i+1)
				}
				fieldObj.SetBool(boolValue)
			}
		}
	}
	return
}

// 返回一个字符串切片，每个字符串长度小于300
func LoadFile(fileName string) []string {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("%s文件不存在，请重试", fileName)
	}
	lines := strings.Split(string(file), "\n")
	slice := make([]string, 0, 10)
	// 按段拆分
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		runeLine := []rune(line) // 用于中英文混输统计长度，故统一转成[]rune
		runeSize := len(runeLine)
		if runeSize == 0 {
			continue
		}
		if runeSize < 300 {
			slice = append(slice, line)
		} else {
			// 一段超过300则按句拆分
			base := 0   // 每次正式切分时的位置
			offset := 0 // 上次的位置
			length := 0 // 当前已计算的字符长度
			for base < runeSize {
				count := getIndex(runeLine[offset:], '。', '！', '？', '”', '.', '!', '?', ',', '"') + 1
				length += count // 更新长度
				offset += count // 更新相对位置
				if length < 300 {
					// 继续遍历累加长度
					// 如果到末尾为止都不足300，则要算进去
					if runeSize-base <= 300 {
						slice = append(slice, string(runeLine[base:]))
						break // 算完了
					}
				} else {
					// length已经大于300了，应该把之前的加进去
					length -= count
					offset -= count
					slice = append(slice, string(runeLine[base:base+length]))
					base += offset
					length = 0
				}
			}
		}
	}
	return slice
}

func getIndex(slice []rune, token ...rune) int {
	for i := range token {
		for j := range slice {
			if slice[j] == token[i] {
				return j
			}
		}
	}
	return -1
}
