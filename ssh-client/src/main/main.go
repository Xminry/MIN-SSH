package main

//import (
//	"bufio"
//	"fmt"
//	"net"
//	"os"
//	"ssh-client/src/main/communicate"
//	"strings"
//)
//
//func main()  {
//	con, err := net.Dial("tcp","www.involute.cn:8765")
//	if err != nil {
//		fmt.Println("客户端连接失败，错误：",err)
//		return
//	}
//	fmt.Println("连接成功")
//	reader := bufio.NewReader(os.Stdin)
//	for{
//		buf := make([]byte,65535)
//		line, err := reader.ReadString('\n')
//		if err != nil {
//			fmt.Println("终端输入失败，错误:",err)
//			return
//		}
//		line = strings.Trim(line,"\r\n")
//		if line == "exit" {
//			fmt.Println("客户端退出……")
//			break
//		}
//		_, err = communicate.Write(con, []byte(line+"\n"))
//		if err != nil {
//			fmt.Println("发送失败，错误：",err)
//		}
//		read, err := communicate.Read(con,buf)
//		if err != nil {
//			fmt.Println("服务器错误：",err)
//			return
//		}
//		fmt.Println(string(buf[:read]))
//	}
//}

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"ssh-client/src/main/process"
	_struct "ssh-client/src/main/struct"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp","www.involute.cn:8765")
	if err != nil {
		fmt.Println("客户端连接失败，错误：",err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("服务器连接成功")
	sshReaderWriter := _struct.NewSshReaderWriter(conn)

	//密钥交换阶段
	p := 7000487615661733
	g := 5925845745820835
	key, err := process.Dh(p,g,sshReaderWriter)
	if err != nil {
		fmt.Println("客户端会话密钥阶段异常，请重试",err)
		return
	}
	//输出key，临时
	//fmt.Println("生成会话密钥：",key)

	//密钥交换阶段结束

	fmt.Println("请输入用户名：")
	userName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	fmt.Println("请输入密码：")
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	sshReaderWriter.Write([]byte(userName+"\n"))
	sshReaderWriter.Write([]byte(password+"\n"))
	//这里开始，登录完成，会话阶段使用key加密
	sshReaderWriter.SetKey(key)

	go handleConnectionReader(conn)


	fmt.Printf("> ")
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		if strings.Compare("exit", strings.Replace(input, "\n", "", -1)) == 0 {
			break
		}

		if strings.HasPrefix(input,"cd") || strings.HasPrefix(input,"chmod"){
			sshReaderWriter.Write([]byte(input))
			sshReaderWriter.Write([]byte("pwd\n"))
		}else{
			sshReaderWriter.Write([]byte(input))
		}
		sshReaderWriter.Write([]byte("whoami\n"))
		sshReaderWriter.Write([]byte("pwd|xargs basename\n"))
	}
	conn.Close()
}

func handleConnectionReader(c net.Conn) {
	//sshReaderWriter := _struct.NewSshReaderWriter(c)
	count :=0
	result := ""
	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			os.Exit(1)
		}
		if string(buf[:n])=="登录成功！"{
			fmt.Printf("%s", buf[:n])
			fmt.Print("\n>")
			count=0
		}else{
			if count==0{
				fmt.Printf("%s", buf[:n])
				count++
			}else if count==1{
				result = result + "[" + string(buf[:n-1])
				count++
			}else {
				result = result + " " + string(buf[:n-1])
				count = 0
				fmt.Print(result + "]>")
				result = ""
			}
		}
	}
}







