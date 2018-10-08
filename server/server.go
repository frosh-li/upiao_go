package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ConnMap struct {
	rwmutex *sync.RWMutex
	Data    map[string]string
}

// Map 的PUT方法
func (cmap *ConnMap) Put(key string, value string) {
	cmap.rwmutex.Lock()
	defer cmap.rwmutex.Unlock()
	cmap.Data[key] = value
}

// Map 的GET方法
func (cmap *ConnMap) Get(key string) string {
	cmap.rwmutex.RLock()
	defer cmap.rwmutex.RUnlock()
	oldValue, err := cmap.Data[key]
	if err {
		return oldValue

	}
	return ""
}

// Map 的Len方法
func (cmap *ConnMap) Len() int {
	cmap.rwmutex.Lock()
	defer cmap.rwmutex.Unlock()
	clen := len(cmap.Data)
	return clen
}

// Map 的Remove方法
func (cmap *ConnMap) Remove(key string) {
	cmap.rwmutex.Lock()
	defer cmap.rwmutex.Unlock()
	delete(cmap.Data, key)
}

func NewSafeMap() *ConnMap {
	sm := new(ConnMap)
	sm.rwmutex = new(sync.RWMutex)
	sm.Data = make(map[string]string)
	return sm
}

var lk *sync.RWMutex

func main() {

	fmt.Println("start tcp server ====>")
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":60026")
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	ConnMapData := NewSafeMap()

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
			// ConnMapData.Remove(tcpConn.RemoteAddr().String())
			continue
		} else {

			tcpConn.SetKeepAlive(true)
			tcpConn.SetNoDelay(true)
			go HandleConn(tcpConn, *ConnMapData)

		}

	}

}

func HandleConn(tcpConn net.Conn, ConnMapData ConnMap) {

	defer tcpConn.Close()
	// fmt.Println("连接的客户端信息：", remoteAddress)
	for {
		var buf = make([]byte, 128)
		n, err := tcpConn.Read(buf[0:])
		remoteAddress := tcpConn.RemoteAddr().String()
		if err != nil {
			log.Println("read error", remoteAddress, err)
			ConnMapData.Remove(remoteAddress)
			return
		} else {

			// oldData := ConnMapData.Data[remoteAddress]
			oldData := ConnMapData.Get(remoteAddress)
			socketDatas := oldData + string(buf[:n])
			reg := regexp.MustCompile(`<[^>]*>`)
			matchs := reg.FindString(socketDatas)
			if matchs != "" { // 匹配到了，进行替换
				fmt.Println("完整匹配到数据", matchs)
				newDatas := strings.Replace(socketDatas, matchs, "", 0)
				ConnMapData.Put(remoteAddress, newDatas)
			} else {
				if len(socketDatas) > 10000 {
				} else {
					// fmt.Println("put", remoteAddress, "len < 10000")
					//lk.Lock()
					/// fmt.Println("ERRORS", remoteAddress, socketDatas)
					// ConnMapData.Data[remoteAddress] = socketDatas

					ConnMapData.Put(remoteAddress, socketDatas)
				}

			}
		}

	}
}
