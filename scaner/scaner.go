// tscan
package main

import (
	"bufio"
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"github.com/vyegres/nscaner/server/storage"
	proxy "go.net/proxy"
	"net"
	"net/rpc/jsonrpc"
	"os"
)

const (
	BLOCK_SIZE = 256
)

func scan(ip string) {
	conn, err := net.Dial("tcp", ip+":80")
	if err != nil {
	} else {
		conn.Close()
		fmt.Println("Подключились", ip)
	}
}

func scanDiapazon(ch chan storage.IpItem, quit chan uint32, ipstart, ipend uint32) {
	var ip uint32
	var ipItem storage.IpItem
	proxy, _ := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	for ip = ipstart; ip < ipend; ip++ {
		ipItem.IP = ip
		conn, err := proxy.Dial("tcp", ipItem.String()+":23")
		if err != nil {
		} else {
			hello, err := bufio.NewReader(conn).ReadString(':')
			if err == nil {
				output, _ := iconv.ConvertString(hello, "latin1", "utf-8")
				ipItem.Hello = output
				conn.Close()
				ch <- ipItem
			}
		}
	}
	quit <- ipstart / BLOCK_SIZE
}

func main() {
	var reply uint32
	var ch chan storage.IpItem = make(chan storage.IpItem, 10)
	var quit chan uint32 = make(chan uint32, 10)
	var i uint32

	/*
		if c, err := proxy.Dial("tcp", "www.silisoftware.com:80"); err != nil {
			fmt.Println("SOCKS5.Dial failed: %v", err)
		} else {

			fmt.Fprintf(c, "GET /tools/ip.php HTTP/1.0\r\nHost:www.silisoftware.com\r\n\r\n")
			status, _ := bufio.NewReader(c).ReadBytes(128)
			//ReadString('\n')
			fmt.Println("ol", string(status))
			c.Close()
		}
		os.Exit(0)*/

	client, err := jsonrpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		os.Exit(1)
	}
	//i = 16777216

	var goCount int = 0

	for j := 0; j < 100; j++ {
		err = client.Call("BlockStorage.GetNew", nil, &reply)
		if err != nil {
			fmt.Println("Блок не получен:", err)
		}
		fmt.Println("GET BLOCK:", reply)
		i = reply * BLOCK_SIZE
		go scanDiapazon(ch, quit, i, i+BLOCK_SIZE)
		goCount++
	}

	for goCount > 0 {
		select {
		case ipItem := <-ch:
			fmt.Println("IP:", ipItem.String())
			err = client.Call("IpStorage.AddIP", ipItem, &reply)
			if err != nil {
				fmt.Println("Ошибка добавления IP:", ipItem, err)
			}
		case block := <-quit:
			fmt.Println("Block Done:", block)
			err = client.Call("BlockStorage.SetDone", block, &reply)
			if err != nil {
				fmt.Println("Ошибка вызова функции BlockStorage.SetDone:", block, err)
			}
			goCount--
		}
	}
	fmt.Println("End Progamm. Last IP is ", i)

	var input string
	fmt.Scanln(&input)
}
