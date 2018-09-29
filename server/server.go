package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type ConnMap struct {
	rwmutex sync.RWMutex
	Data    map[string]string
}

// Map 的PUT方法
func (cmap *ConnMap) Put(key string, value string) (string, bool) {

	cmap.rwmutex.Lock()
	defer cmap.rwmutex.Unlock()

	oldValue := cmap.Data[key]
	cmap.Data[key] = value
	return oldValue, true
}

// Map 的GET方法
func (cmap *ConnMap) Get(key string) string {

	cmap.rwmutex.RLock()
	defer cmap.rwmutex.RUnlock()

	oldValue := cmap.Data[key]
	return oldValue
}

// Map 的GET方法
func (cmap *ConnMap) Len() int {

	cmap.rwmutex.RLock()
	defer cmap.rwmutex.RUnlock()

	clen := len(cmap.Data)
	return clen
}

// Map 的GET方法
func (cmap *ConnMap) Remove(key string) string {

	cmap.rwmutex.Lock()
	defer cmap.rwmutex.Unlock()

	oldValue := cmap.Get(key)
	delete(cmap.Data, key)
	return oldValue
}

func main() {
	fmt.Println("start tcp server ====>")
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":60026")
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	ConnMapData := ConnMap{}
	ConnMapData.Data = make(map[string]string)

	ConnectTicker := time.NewTimer(time.Second * 10) // 十秒打印一次ConnMap长度

	go func() {
		for {
			select {
			case <-ConnectTicker.C:
				fmt.Printf("当前连接数 %v \n", ConnMapData.Len())
				ConnectTicker.Reset(time.Second * 10)
			}
		}
	}()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Println("tcpConn error", err)
			ConnMapData.Remove(tcpConn.RemoteAddr().String())
			tcpConn.Close()
			continue
		} else {
			// ConnMapData.Data[tcpConn.RemoteAddr().String()] = ""
			tcpConn.SetKeepAlive(true)
			tcpConn.SetNoDelay(true)
			go HandleConn(tcpConn, ConnMapData)
		}

	}
}

func HandleConn(tcpConn net.Conn, ConnMapData ConnMap) {
	defer tcpConn.Close()
	remoteAddress := tcpConn.RemoteAddr().String()
	fmt.Println("连接的客户端信息：", remoteAddress)
	for {
		var buf = make([]byte, 1024)
		n, err := tcpConn.Read(buf)

		if err != nil {
			log.Println("read error", remoteAddress, err)
			ConnMapData.Remove(remoteAddress)
			return
		} else {

			oldData := ConnMapData.Get(remoteAddress)
			ConnMapData.Put(remoteAddress, oldData+string(buf[:n]))
		}

		socketDatas := ConnMapData.Get(remoteAddress)
		fmt.Println(remoteAddress, socketDatas)

		if len(socketDatas) > 10000 {
			// 如果长时间无法消费内容，直接清理掉
			ConnMapData.Put(remoteAddress, "")
		}
	}
}
