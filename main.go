package main

import (
	"fmt"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main(){
	var argLen = len(os.Args)
	if argLen < 2 {
		fmt.Println("启动传入参数错误，第一个参数为扫描的路径如：/diskD/movie/,第二个参数为替换的后缀名如：.srt|.ass")
		return
	}

	_,err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Println("未找到路径：", os.Args[1])
		return
	}

	var extList = []string{".srt", ".ass"}
	if argLen > 2 {
		extList = strings.Split(os.Args[2], "|")
	}

	filepath.Walk(os.Args[1],
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				fmt.Println("dir:", path)
				return nil
			}

			for _,ext := range extList {
				if filepath.Ext(path) == ext {
					ConvertUft8(path)
				}
			}

			return nil
		})
}

func ConvertUft8(path string){
	file, err := os.Open(path) // For read access.
	if err != nil {
		fmt.Println("打开文件出错", err)
		return
	}
	defer file.Close()
	buffer := make([]byte, 32<<10)
	textDetector := chardet.NewTextDetector()
	size, err := io.ReadFull(file, buffer)
	if err == io.EOF {
		fmt.Println("读取文件内容为空", err)
		return
	}
	input := buffer[:size]
	var detector = textDetector
	result, err := detector.DetectBest(input)
	if err != nil {
		fmt.Println("解析编码出错", err)
		return
	}

	content,err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取文件所有内容出错", err)
		return
	}

	//fmt.Println("原始编码内容：", string(content), result.Charset)
	if strings.Contains(result.Charset, "ISO-") || result.Charset == "GB-18030" {
		contentUtf8 := ConvertByte2String(content, "GB-18030")
		//fmt.Println("utf-8编码：",contentUtf8)
		ioutil.WriteFile(path, []byte(contentUtf8), 0666)
		fmt.Println("成功转码：", path)
	}
}

func ConvertByte2String(byte []byte, charset string) string {
	var str string
	switch charset {
	case "GB-18030":
		var decodeBytes,_= simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str= string(decodeBytes)
	case "GBK":
		var decodeBytes,_= simplifiedchinese.GBK.NewDecoder().Bytes(byte)
		str= string(decodeBytes)
	case "HZ-GB2312":
		var decodeBytes,_= simplifiedchinese.HZGB2312.NewDecoder().Bytes(byte)
		str= string(decodeBytes)
	case "UTF-8":
		str = string(byte)
	default: str = string(byte)
	}
	return str
}