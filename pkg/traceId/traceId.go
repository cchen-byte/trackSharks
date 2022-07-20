package traceId

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

type TraceId struct {
	IpData        string // 8位, 产生TraceId的机器的 IP, 这是一个十六进制的数字, 每两位代表 IP 中的一段, 我们把这个数字, 按每两位转成 10 进制即可得到常见的 IP 地址表示方式 xx.xx.xx.xx
	TimeStamp     string // 13位, 毫秒级时间戳
	AutoIncrement int32 // 4位, 自增的序列, 从 1000 涨到 9000, 到达 9000 后回到 1000 再开始往上涨. 用于避免多线程并发时 TraceId 碰撞(可以保证在 qps=9000000 以下不碰撞)
	ProcessId     string // 5位, 当前的进程 ID
}

// GetIpData 获取当前机器IP
func (id *TraceId) GetIpData() string {
	//var ip16Data []string
	//conn, _ := net.Dial("udp", "8.8.8.8:53")
	//localAddr := conn.LocalAddr().(*net.UDPAddr)
	//ipStr := strings.Split(localAddr.String(), ":")[0]
	//ipList := strings.Split(ipStr, ".")
	//for _, ipStr := range ipList{
	//	ipInt, _ := strconv.Atoi(ipStr)
	//	ip16Data = append(ip16Data, strconv.FormatInt(int64(ipInt), 16))
	//}
	//return strings.Join(ip16Data, "")
	return "c0a81414"
}

// GetAutoIncrement 获取自增的序列
func (id *TraceId) GetAutoIncrement() int32 {
	if id.AutoIncrement == 0 {
		atomic.AddInt32(&id.AutoIncrement, 1000)
	} else if id.AutoIncrement >= 9000 {
		atomic.AddInt32(&id.AutoIncrement, -8000)
	}else{
		atomic.AddInt32(&id.AutoIncrement, 1)
	}
	return id.AutoIncrement
}

// GetTimeStamp 返回当前时间毫秒级时间戳
func (id *TraceId) GetTimeStamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}


// GetProcessId 当前的进程 ID
func (id *TraceId) GetProcessId() string {
	return strconv.FormatInt(int64(os.Getppid()), 10)
}

func (id *TraceId) GetTraceId() string {
	return fmt.Sprintf("%s%s%v%s", id.GetIpData(), id.GetTimeStamp(), id.GetAutoIncrement(), id.GetProcessId())
}