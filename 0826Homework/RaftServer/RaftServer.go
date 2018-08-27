package RaftServer

import (
	"fmt"
	"net"
	"strings"
	"time"
	"math/rand"
	"net/http"
	"strconv"
)

/*
	节点角色
 */
const (
	LEADERYS    = iota
	CANDIDATEYS
	FOLLOWERYS
)

type AddrYS struct {
	AddrYS string //节点地址 host:port
}

/*
	Raft服务器结构体
 */
type ServerYS struct {
	VotesYS         int         //选票数量
	RoleYS          int         //角色 follower candidate leader
	NodesYS         []AddrYS    //保存已知节点s
	isElectingYS    bool        //判断当前节点是否处于选举中
	TimeoutYS       int         //选举间隔时间（也叫超时时间）
	AddressYS       string      //地址
	ElecChanYS      chan bool   //通道信号
	HeartBeatChanYS chan bool   //leader 的心跳信号
	CusMsgYS        chan string //网页接收到的参数 由主节点向子节点传参
}

/*
	创建一个节点,初始角色为follower
 */
func CreateNewRaftServerYS(nodeAddr string, addrs []AddrYS) *ServerYS {
	rs := ServerYS{}
	rs.isElectingYS = true
	rs.VotesYS = 0
	rs.changeRoleYS(FOLLOWERYS)
	rs.ElecChanYS = make(chan bool)
	rs.HeartBeatChanYS = make(chan bool)
	rs.CusMsgYS = make(chan string)
	rs.resetTimeoutYS()
	rs.NodesYS = addrs
	rs.AddressYS = nodeAddr

	return &rs
}

/*
	设置节点角色
 */
func (rs *ServerYS) changeRoleYS(role int) {
	switch role {
	case LEADERYS:
		fmt.Println("leader")
	case CANDIDATEYS:
		fmt.Println("candidate")
	case FOLLOWERYS:
		fmt.Println("follower")

	}
	rs.RoleYS = role
}

/*
	设置节点选举时间间隔
 */
func (rs *ServerYS) resetTimeoutYS() {
	//Raft系统一般为1500-3000毫秒选一次
	rs.TimeoutYS = 2000
}

/*
	运行服务器
 */
func (rs *ServerYS) RunYS() {
	//rs监听 是否有人 给我投票
	listen, _ := net.Listen("tcp", rs.AddressYS)

	defer listen.Close()

	//启动各个子线程
	go rs.electYS()

	go rs.electTimeDurationYS()

	go rs.sendHeartBeat()

	go rs.sendDataToOtherNodes()

	go rs.setHttpServer()


	//监听接受消息处理
	for {
		conn, _ := listen.Accept()
		go func() {

			for {
				message := messageOfConnYS(conn)

				if len(message) > 0 {
					fmt.Println("收到消息", message)

					if message == rs.AddressYS {
						rs.VotesYS++
						fmt.Println("当前票数：", rs.VotesYS)
						// leader 选举成功
						if voteSuccessYS(rs.VotesYS, 5) {
							fmt.Printf("我是 %s, 我被选举成leader", rs.AddressYS)

							//通知其他节点。停止选举
							rs.sendMessageYS("stopVote")
							//停止当前节点投票
							rs.isElectingYS = false
							//改变当前节点状态
							rs.changeRoleYS(LEADERYS)
							break
						}
					}

					//收到leader发来的消息
					if strings.HasPrefix(message, "stopVote") {
						//停止给别人投票
						rs.isElectingYS = false
						//回退自己的状态
						rs.changeRoleYS(FOLLOWERYS)
						break
					}
				}
			}

		}()
	}

}

/*
	判断节点票数是否足够成为leadder
 */
func voteSuccessYS(vote int, target int) bool {
	if vote >= target {
		return true
	}
	return false
}

/*
	发送数据
	投票/状态改变为leader时通知其他节点/leader心跳/leader发送消息
 */
func (rs *ServerYS) sendMessageYS(data string) {
	//这里遍历所有节点，如果某个节点没有响应，就会进入死循环直到全部非自己的节点连接上
	for _, k := range rs.NodesYS {
		if k.AddrYS != rs.AddressYS {
		label:
			conn, err := net.Dial("tcp", k.AddrYS)
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

/*
	投票给其他节点
 */
func (rs *ServerYS) electYS() {

	for {
		//通过通道确定现在可以给别人投票
		<-rs.ElecChanYS

		//获取一个随机的已知节点进行投票
		vote := rs.getVoteNumYS()

		//投票, 消息为要投票节点的地址
		rs.sendMessageYS(vote)
		// 设置选举状态
		if rs.RoleYS != LEADERYS {
			rs.changeRoleYS(CANDIDATEYS)
		} else {
			//是leader的情况
			return
		}

	}
}

/*
	随机生成投票给的节点地址
 */
func (rs *ServerYS) getVoteNumYS() string {

	//不能投票给自己
	for {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(rs.NodesYS))
		if rs.NodesYS[i].AddrYS != rs.AddressYS {
			return rs.NodesYS[i].AddrYS
		}
	}

}

/*
	投票时间间隔生效
 */
func (rs *ServerYS) electTimeDurationYS() {
	//
	fmt.Println("+++", rs.isElectingYS)
	for {
		if rs.isElectingYS {
			rs.ElecChanYS <- true
			time.Sleep(time.Duration(rs.TimeoutYS) * time.Millisecond)
		}

	}
}

/*
	打印当前对象的角色
 */
func (rs *ServerYS) printRoleYS() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(rs.AddressYS, "状态为", rs.RoleYS, rs.isElectingYS)
	}
}



/*
	主节点发送心跳信号给其他节点
 */
func (rs *ServerYS) sendHeartBeat() {
	// 每隔1s 发送一次心跳
	for {
		time.Sleep(1 * time.Second)
		if rs.RoleYS == LEADERYS {
			rs.sendMessageYS("heat beating")
		}
	}
}

/*
	leader节点发送消息到其他节点
 */
func (rs *ServerYS) sendDataToOtherNodes() {
	for {
		//Leader从http收到消息时,CusMsg通道收到后,发送到其他节点
		msg := <-rs.CusMsgYS
		if rs.RoleYS == LEADERYS {
			rs.sendMessageYS(msg)

		}
	}
}

/*
	开启http服务器
 */
func (rs *ServerYS) setHttpServer() {

	//http:localhost:5010/req?data=123456
	//调整port + 10
	http.HandleFunc("/req", rs.requestYS)
	//httpPort := rs.Port + 10
	array := strings.Split(rs.AddressYS, ":")
	host := array[0]
	port := array[1]
	//fmt.Println("host:", host, "port:", port)
	portInt, _ := strconv.Atoi(port)
	portInt = portInt + 10
	port = strconv.Itoa(portInt)
	finalAddr := host + ":" + port
	//fmt.Println("finalAddr:", finalAddr)
	if err := http.ListenAndServe(finalAddr, nil); err == nil {
		fmt.Println(err)
	}

}

/*
	leader向其他子节点发送数据
 */
func (rs *ServerYS) requestYS(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	if len(request.Form["data"][0]) > 0 {
		writer.Write([]byte("ok"))
		fmt.Println("---------request handle-----------")
		fmt.Println(request.Form["data"][0])
		fmt.Println("---------request handle-----------")
		rs.CusMsgYS <- request.Form["data"][0]
	}
}

/*
	将conn收到的内容转成string
 */
func messageOfConnYS(conn net.Conn) string {
	by := make([]byte, 1024)
	n, _ := conn.Read(by)
	value := string(by[:n])
	return value
}