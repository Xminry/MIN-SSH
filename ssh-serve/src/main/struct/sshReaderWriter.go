package _struct

import (
	"bufio"
	"fmt"
	"net"
	"ssh-serve/src/main/utils"
)

type SshReaderWriter struct {
	con net.Conn
	reader *bufio.Reader
	key []byte
	useKey bool
}

func NewSshReaderWriter(conn net.Conn) *SshReaderWriter {
	return &SshReaderWriter{
		con: conn,
		reader: bufio.NewReader(conn),
		key:nil,
		useKey: false,
	}
}

func (s *SshReaderWriter) Read(p []byte) (int, error) {
	if !s.useKey{
		recvStr, err := utils.Decode(s.reader)
		if err!=nil {
			fmt.Println(err)
			return 0,err
		}
		copy(p, []byte(recvStr))
		fmt.Println("message:"+recvStr)
		return len(recvStr),nil
	}else{
		recvStr, err := utils.Decode(s.reader)
		encrypt := []byte(recvStr)
		if err!=nil {
			fmt.Println(err)
			return 0,err
		}
		// 先异或运算解密
		for i,item := range encrypt {
			encrypt[i] = item ^ 28
		}
		// 再AES解密
		decrypt,_ := utils.AesDecrypt(encrypt, s.key)
		fmt.Println("decrypt:", string(decrypt))
		copy(p, decrypt)
		return len(decrypt),nil
	}
}

func (s *SshReaderWriter) Write(buf []byte) (int,error) {
	message := string(buf)
	if !s.useKey{
		if buf==nil || len(buf)==0 || len(message)==0{
			buf = []byte("#")
		}
		encode, err := utils.Encode(message)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		s.con.Write(encode)
		return len(encode),nil
	}else {
		//先AES加密
		encrypt,_ := utils.AesEncrypt([]byte(message), s.key)
		// 再异或运算加密
		for i,item := range encrypt {
			encrypt[i] = item ^ 28
		}
		encode, err := utils.Encode(string(encrypt))
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		s.con.Write(encode)
		return len(encode),nil
	}

}

func (s *SshReaderWriter) Close() error{
	err := s.con.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *SshReaderWriter) SetKey(key []byte){
	s.key = key
	s.useKey = true
}
