package main

//import (
//	"fmt"
//	"net"
//	"os/exec"
//	"ssh-serve/src/main/communicate"
//)
//
//func process(con net.Conn){
//	defer con.Close()
//
//	for{
//		buf := make([]byte,65535)
//		fmt.Printf("服务器等待客户端 %s 发送信息\n",con.RemoteAddr().String())
//		n, err := communicate.Read(con,buf)
//		if err != nil {
//			fmt.Println("客户端已退出")
//			return
//		}else{
//			fmt.Println("收到客户端 %s 数据：%s",con.RemoteAddr().String(),string(buf[:n]))
//			cmd := exec.Command("/bin/bash","/k",string(buf[:n]))
//			output, err := cmd.Output()
//			if err != nil {
//				fmt.Println("执行命令异常：",err)
//				communicate.Write(con, []byte("执行命令错误"))
//			}else{
//				result := string(output)
//				fmt.Println(result)
//				if result=="" || result=="\n"{
//					communicate.Write(con,[]byte(" \n"))
//				}else{
//					communicate.Write(con, output)
//				}
//			}
//		}
//	}
//}
//
//func main(){
//	fmt.Println("服务端开始监听……")
//	listen, err := net.Listen("tcp","0.0.0.0:8765")
//	if err != nil {
//		fmt.Println("监听失败，错误：",err)
//		return
//	}
//
//	defer listen.Close()
//
//	for{
//		fmt.Println("等待客户端连接……")
//		con, err := listen.Accept()
//		if err != nil {
//			fmt.Println("Accept失败,错误：",err)
//			return
//		}else{
//			fmt.Println("Accept()成功，客户端IP= %s ",con.RemoteAddr().String())
//		}
//		go process(con)
//	}
//}

import (
	"MINTCP/socket"
	"bytes"
	"fmt"
	"minlib/common"
	"minlib/minsecurity"
	"minlib/security"
	"os"
	"os/exec"
	"ssh-serve/src/main/communicate"
	"ssh-serve/src/main/process"
	_struct "ssh-serve/src/main/struct"
	"strings"
)

func main() {
	/*listen, err := net.Listen("tcp", "0.0.0.0:8765")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listen.Close()*/

	common.InitLogger(&common.LoggerParameters{
		ReportCaller: true,
		LogLevel:     "ALL",
		LogFormat:    "text",
	})
	keyChain := new(security.KeyChain)
	if err := keyChain.Init(); err != nil {
		common.LogFatal(err)
	}

	//
	identity := keyChain.GetIdentityByName("/localhost/operator")
	//identity := keyChain.GetIdentityByName("/test/server")
	// identity, err := keyChain.CreateIdentityByName("/test/server1", "123456")
	_, err := identity.UnLock("123456", minsecurity.SM4ECB)
	// TODO: require identity is unlock, if it's locked, please do unlock first
	listen, err := socket.Listen("min-push-tcp", 80, identity,
		"unix", "/tmp/mir-tcp-message-channel-stack.sock")
	if err != nil {
		common.LogFatal(err)
	} else {
		common.LogInfo("Listen Success")
	}

	con, err := listen.Accept()
	if err != nil {
		common.LogError(err)
		os.Exit(1)
	}
	sshReaderWriter := _struct.NewSshReaderWriter(con)

	//密钥交换阶段
	p := 7000487615661733
	g := 5925845745820835
	key, err := process.Dh(p, g, sshReaderWriter)
	if err != nil {
		fmt.Println("服务端会话密钥阶段异常，请重试", err)
		return
	}
	//密钥交换阶段结束

	buf := make([]byte, 65535)
	n, err := sshReaderWriter.Read(buf)
	if err != nil {
		fmt.Println("客户端已退出")
		return
	}
	username := string(buf[:n])
	fmt.Println("当前尝试登录用户名：" + username)

	//查询用户的加密密码
	shadowCmd := exec.Command("bash", "-c", "cat /etc/shadow | grep "+username)
	output, err := shadowCmd.Output()
	if err != nil {
		fmt.Println("执行命令异常：", err)
		communicate.Write(con, []byte("执行命令错误"))
	} else {
		result := string(output)
		//以盐值开头的子字符串,这里暂时跳过了加密算法，默认MD5
		saltStartStr := result[strings.Index(result, "$")+3:]
		fmt.Println(saltStartStr)
		salt := saltStartStr[:strings.Index(saltStartStr, "$")]
		fmt.Println("盐值：" + salt)
		//以加密密码开头的字符串
		encryPwdStr := saltStartStr[strings.Index(saltStartStr, "$")+1:]
		encryPwd := encryPwdStr[:strings.Index(encryPwdStr, ":")]
		fmt.Println("加密密码：" + encryPwd)

		//接下来检测密码
		n, err := sshReaderWriter.Read(buf)
		if err != nil {
			fmt.Println("客户端已退出")
			return
		}
		password := string(buf[:n])
		fmt.Println("用户输入的密码：" + password)
		pwdCmd := exec.Command("bash", "-c", "openssl passwd -6 -salt "+salt+" "+password)
		userEncryPwdByte, err := pwdCmd.Output()
		if err != nil {
			fmt.Println("执行命令异常：", err)
			communicate.Write(con, []byte("执行命令错误"))
		}
		//手动加密后的字符串做截取，截取完包含盐值，仍需再次截取
		includeSaltResult := string(userEncryPwdByte)[strings.Index(string(userEncryPwdByte), "$")+3:]
		userEncryPwd := includeSaltResult[strings.Index(includeSaltResult, "$")+1:]
		fmt.Println("用户输入的加密密码：" + userEncryPwd)
		if strings.Compare(userEncryPwd[:len(userEncryPwd)-1], encryPwd) == 0 {
			con.Write([]byte("登录成功！"))
			fmt.Println(username + "登录成功！")
		} else {
			con.Write([]byte("登录失败！"))
			fmt.Println(username + "登录失败！")
			return
		}
	}

	//登录成功了，设置会话交互阶段的key
	sshReaderWriter.SetKey(key)

	//cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = sshReaderWriter
	cmd.Stdout = con
	cmd.Stderr = sshReaderWriter
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("连接结束，断开")
}

func SetCommandStd(cmd *exec.Cmd) (stdout, stderr *bytes.Buffer) {
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return
}
