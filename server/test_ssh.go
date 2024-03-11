package main

//import (
//	"golang.org/x/crypto/ssh"
//	"log"
//)
//
//func main() {
//	// SSH 连接配置
//	config := &ssh.ClientConfig{
//		User: "root",
//		Auth: []ssh.AuthMethod{
//			ssh.Password("123456"),
//		},
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//	}
//
//	// 循环 SSH 连接服务器
//	for {
//		// 建立 SSH 连接
//		_, err := ssh.Dial("tcp", "141.11.86.94:22", config)
//		//client, err := ssh.Dial("tcp", "165.154.149.134:22", config)
//		if err != nil {
//			log.Println("Failed to dial: ", err)
//		}
//		//defer client.Close()
//
//		// 执行命令
//		//session, err := client.NewSession()
//		//if err != nil {
//		//	log.Fatal("Failed to create session: ", err)
//		//}
//		//defer session.Close()
//		//
//		//// 这里可以替换为你要执行的具体命令
//		//cmd := "echo Hello"
//		//
//		//// 执行命令并获取输出
//		//output, err := session.CombinedOutput(cmd)
//		//if err != nil {
//		//	log.Fatal("Failed to execute command: ", err)
//		//}
//		//
//		//fmt.Println(string(output))
//	}
//}
