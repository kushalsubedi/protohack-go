package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)
const (
	MAXBUF  = 1000
	VERSION = "Ken's Key-Value Store 1.0"
)
type KVStore struct {
	data map[string]string 
	mu sync.RWMutex
}

func NewKVStore() *KVStore{
	return &KVStore{
		data: make(map[string]string),
	}
}
func (kv *KVStore) Set(key,value string){
	if key == "version"{
		return
	}
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key]=value
}

func (kv *KVStore) Get(key string)string{
	if key == "version" {
		return VERSION
	}
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	return kv.data[key]
}


func handleConn(conn *net.UDPConn, kv *KVStore){
	defer conn.Close()
	buffer := make([]byte, MAXBUF)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err !=nil{
			fmt.Println("Read error:", err)
			continue
		}
		message :=string(buffer[:n])

		if strings.Contains(message, "="){
			parts := strings.SplitN(message, "=",2)
			key:= parts[0]
			value := parts[1]
			kv.Set(key,value)
		} else {
			value := kv.Get(message)
			resp := fmt.Sprintf("%s=%s", message,value)
			_, err:= conn.WriteToUDP([]byte(resp), clientAddr)
			if err != nil {
				fmt.Println("Write error:", err)
			}
		}
	}
}

func main (){
	if len(os.Args) < 2 {
		fmt.Printf("Usage %s:<Port>",os.Args[0])
	}
	addr,err := net.ResolveUDPAddr("udp",":"+os.Args[1])
	
	if err !=nil{
		log.Fatal("Failed to resolve UDP",err)
	}
	conn,err := net.ListenUDP("udp",addr)												
	if err != nil {
		log.Fatal("Error establishing Connection !!")
		return
	}
		defer conn.Close()
	fmt.Printf("UDP key-value server running on port %s\n", os.Args[1])
	kv := NewKVStore()
	handleConn(conn,kv)
}

