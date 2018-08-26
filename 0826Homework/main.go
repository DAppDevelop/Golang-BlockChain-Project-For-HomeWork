package main

import (
	"fmt"
	"time"
	"math/rand"
	"flag"
	"net"
	"strconv"
	"strings"
	"net/http"
)

/*
	作业，
	1, 认真学习raft算法,改写mainyouhua.go为局域网环境下测试

	2，自己测试etcd的局域网分布式部署（不用交
 */


const (
	LEADER    = iota
	CANDIDATE
	FOLLOWER
)

//声明地址信息
type Addr struct {
	Host string //ip
	Port int
	Addr string
}

type RaftServer struct {
	Votes         int //选票
	Role          int // 角色 follower candidate leader
	Nodes         []Addr
	isElecting    bool      //判断当前节点是否处于选举中
	Timeout       int       //选举间隔时间（也叫超时时间）
	ElecChan      chan bool //通道信号
	HeartBeatChan chan bool //leader 的心跳信号
	Port          int       //端口号

	//网页接收到的参数 由主节点向子节点传参
	CusMsg chan string
}

func (rs *RaftServer) changeRole(role int) {
	switch role {
	case LEADER:
		fmt.Println("leader")
	case CANDIDATE:
		fmt.Println("candidate")
	case FOLLOWER:
		fmt.Println("follower")

	}
	rs.Role = role
}

func (rs *RaftServer) resetTimeout() {
	//Raft系统一般为1500-3000毫秒选一次
	rs.Timeout = 2000
}

//运行服务器
func (rs *RaftServer) Run() {
	//rs监听 是否有人 给我投票
	listen, _ := net.Listen("tcp", ":"+strconv.Itoa(rs.Port))

	defer listen.Close()

	go rs.elect()

	//控制投票时间
	go rs.electTimeDuration()

	//go rs.printRole()

	// 主节点发送心跳
	go rs.sendHeartBeat()
	//
	go rs.sendDataToOtherNodes()

	//监听http协议
	go rs.setHttpServer()

	for {
		conn, _ := listen.Accept()
		go func() {

			for {
				by := make([]byte, 1024)
				n, _ := conn.Read(by)
				value := string(by[:n])
				if len(value) > 0 {
					fmt.Println("收到消息", value)
				}

				v, _ := strconv.Atoi(value)
				if v == rs.Port {
					rs.Votes++
					fmt.Println("当前票数：", rs.Votes)
					// leader 选举成功
					if VoteSuccess(rs.Votes, 5) == true {
						fmt.Printf("我是 %v, 我被选举成leader", rs.Port)

						//通知其他节点。停止选举
						//重置其他节点状态和票数
						rs.VoteToOther("stopVote")
						rs.isElecting = false
						//改变当前节点状态

						rs.changeRole(LEADER)
						break
					}
				}

				//收到leader发来的消息
				if strings.HasPrefix(string(by[:n]), "stopVote") {
					//停止给别人投票
					rs.isElecting = false
					//回退自己的状态
					rs.changeRole(FOLLOWER)
					break
				}

			}

		}()
	}

}

func VoteSuccess(vote int, target int) bool {
	if vote >= target {
		return true
	}
	return false
}

//发送数据)
func (rs *RaftServer) VoteToOther(data string) {
	//这里遍历所有节点，如果某个节点没有响应，就会进入死循环直到全部非自己的节点连接上
	for _, k := range rs.Nodes {
		if k.Port != rs.Port {
			if data == "1234" {
				fmt.Println("-------------", k.Port)
			}

		label:
			conn, err := net.Dial("tcp", ":"+strconv.Itoa(k.Port))
			for {
				if err != nil {
					time.Sleep(1 * time.Second)
					goto label
				}
				break
			}
			conn.Write([]byte(data))

		}
	}
}

//给别人投票
func (rs *RaftServer) elect() {

	for {
		//通过通道确定现在可以给别人投票

		<-rs.ElecChan

		//给其他节点投票，不能投给自己
		vote := getVoteNum()

		rs.VoteToOther(strconv.Itoa(vote))
		// 设置选举状态
		if rs.Role != LEADER {
			rs.changeRole(CANDIDATE)
		} else {
			//是leader的情况
			return
		}

	}
}

//随机生成投票给的端口号
func getVoteNum() int {

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(3) + 5000
}

func (rs *RaftServer) electTimeDuration() {
	//
	fmt.Println("+++", rs.isElecting)
	for {
		if rs.isElecting {

			rs.ElecChan <- true
			time.Sleep(time.Duration(rs.Timeout) * time.Millisecond)

		}

	}
}

//打印当前对象的角色
func (rs *RaftServer) printRole() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(rs.Port, "状态为", rs.Role, rs.isElecting)
	}
}

func main() {

	//获取参数
	//运行  go run main.go -p 5000  (p 后面就是要启动的端口)
	port := flag.Int("go", 1234, "port")
	flag.Parse()
	fmt.Println(*port)

	rs := RaftServer{}
	rs.isElecting = true
	rs.Votes = 0
	rs.Role = FOLLOWER
	//控制是否开始投票
	rs.ElecChan = make(chan bool)
	rs.HeartBeatChan = make(chan bool)
	rs.CusMsg = make(chan string)
	rs.resetTimeout()
	rs.Nodes = []Addr{
		{"127.0.0.1", 5000, "5000"},
		{"127.0.0.1", 5001, "5001"},
		{"127.0.0.1", 5002, "5002"},
		//{"127.0.0.1", 5003, "5003"},
	}
	rs.Port = *port

	rs.Run()


}

//主节点发送心跳信号给其他节点
func (rs *RaftServer) sendHeartBeat() {
	// 每隔1s 发送一次心跳
	for {
		time.Sleep(1 * time.Second)
		if rs.Role == LEADER {
			//发送消息
			rs.VoteToOther("heat beating")
		}
	}
}

//通过leader 给其他所有子节点发送数据
func (rs *RaftServer) sendDataToOtherNodes() {
	for {
		msg := <-rs.CusMsg
		if rs.Role == LEADER {
			//发送消息
			rs.VoteToOther(msg)

		}
	}
}

//开启http服务器
func (rs *RaftServer) setHttpServer() {

	//http:localhost:5010/req?data=123456
	http.HandleFunc("/req", rs.request)
	httpPort := rs.Port + 10
	if err := http.ListenAndServe(":"+strconv.Itoa(httpPort), nil); err == nil {
		fmt.Println(err)
	}

}

//leader向其他子节点发送数据
func (rs *RaftServer) request(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	if len(request.Form["data"][0]) > 0 {
		writer.Write([]byte("ok"))
		fmt.Println("request handle")
		fmt.Println(request.Form["data"][0])
		rs.CusMsg <- request.Form["data"][0]
	}

}
