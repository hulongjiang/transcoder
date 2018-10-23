# 文件转码工具

本小工具将各种中文编码文件转成utf-8，注意请勿使用大文件（超过500M）

## 用法

命令行说明：

    ./transcoder xxxx(扫码目录如：/var/html/) xxx(要转码的文件后缀多个后缀用|分割如：.srt|.ass)
    
列子：

    ./transcoder /diskD/movie .srt|.ass
