package main

import (
	"./mock"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func dialToServer(snKey string, groups int, batterys int) {
	conn, err := net.Dial("tcp", "127.0.0.1:60026")
	if err != nil {
		log.Println("连接服务器出错", err)
		return
	}
	defer conn.Close()
	fmt.Println("dialToServer", snKey)
	// aChan := make(chan int, 1)
	stationDataTicker := time.NewTimer(time.Hour * 2) // 两个小时发送一次站数据
	heatbeatTicker := time.NewTimer(time.Second * 1)  // 十秒发送一个心跳包
	errorTicker := time.NewTimer(time.Second * 1)     // 十秒随机测试是否需要发送错误数据包

	msg := []byte(mock.GetStation(snKey, 2, 3))
	conn.Write(msg)
	ch := make(chan string)
	go func() {
		for {
			select {
			case <-stationDataTicker.C:
				msg := []byte(mock.GetStation(snKey, groups, batterys))
				conn.Write(msg)
				fmt.Printf("send station data %v %v\n", snKey, time.Now())
				stationDataTicker.Reset(time.Hour * 2)
			case <-heatbeatTicker.C:
				msg := []byte("<{\"sn_key\":\"" + snKey + "\", \"sid\":123}>")
				conn.Write(msg)
				fmt.Printf("send heatbeat %v\n", time.Now())
				heatbeatTicker.Reset(time.Second * 1)
			case <-errorTicker.C:
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				randSeed := r.Float32()

				fmt.Printf("send error msg %v %v\n", snKey, randSeed)
				// 10%的概率发送错误数据
				if randSeed > 0.9 {
					msg := []byte(mock.GetError(snKey, groups, batterys, ""))
					conn.Write(msg)
					fmt.Printf("send error msg %v %v\n", snKey, time.Now())
				}

				errorTicker.Reset(time.Second * 1)

			}
		}
	}()
	ch <- "Hi"
	// ch := make(chan string)
	go func() {
		for {
			fmt.Println("ready to read from server", snKey)
			var b []byte = make([]byte, 1024)
			n, err := conn.Read(b)
			if err != nil {
				fmt.Println(err.Error())
				// log.Fatal(err)
				continue
			}
			fmt.Println(string(b[:n]))
			if string(b[:n]) == `<{"FuncSel":{"Operator":3}}>` {
				msg := []byte(mock.GetStation(snKey, 2, 3))
				conn.Write(msg)
				fmt.Printf("send station data\n")
			}
		}
	}()
	// ch <- "Hi"
}

func main() {
	ch := make(chan string)
	for i := 100000000; i <= 100001000; i++ {
		go dialToServer(strconv.Itoa(i), 2, 3)
	}

	ch <- "Hi"
}
