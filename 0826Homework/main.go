package main

import (
	"chenyuanshan/0826Homework/CLI"
)

/*
	作业，
	1, 认真学习raft算法,改写mainyouhua.go为局域网环境下测试
	局域网环境下
	go run main.go node -local 192.168.0.100:5001 -addrs '["192.168.0.100:5000","192.168.0.100:5001","192.168.0.106:5000"]'
	只有本地节点时
	go run main.go node -local :5002 -addrs '[":5000",":5001",":5002"]'
 */


func main() {
	cli := CLI.CLIYS{}
	cli.RunYS()
}




