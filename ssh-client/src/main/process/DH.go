package process

import (
	"fmt"
	"math/big"
	"math/rand"
	_struct "ssh-client/src/main/struct"
	"strconv"
	"time"
)

// DH密钥协商算法,底数P，模板G,返回会话密钥K
func Dh(p int,g int,sshReaderWriter *_struct.SshReaderWriter) ([]byte,error){
	// p为16位质数
	//p := big.NewInt(7000487615661733)
	//g := big.NewInt(5925845745820835)
	P := big.NewInt(int64(p))
	G := big.NewInt(int64(g))

	// EXP(a, b, c) = (a ** b) % c
	// 服务器生成A发给客户端
	rand.Seed(time.Now().Unix())
	targer := rand.Intn(900) + 100
	a := big.NewInt(int64(targer))
	A := big.NewInt(0).Exp(G, a, P)
	//fmt.Println("生成客户端公钥：",  A)
	//发送公钥至服务端
	sshReaderWriter.Write([]byte(A.String()))
	buf := make([]byte, 1024)
	n, err := sshReaderWriter.Read(buf)
	if err != nil {
		return nil, nil
	}
	B, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		fmt.Println("DH算法内部错误，格式转换异常！")
		return nil, nil
	}
	//fmt.Println("收到服务端公钥：",B)

	// 服务器拿到客户端的B，生成密钥K
	K := big.NewInt(0).Exp(big.NewInt(int64(B)), a, P)
	//fmt.Println("生成会话密钥K", K)

	// 最终SA=SB，完成密钥协商
	key := []byte(K.String())
	//将16位key拼接成32位key
	key = append(key, key...)
	return key,nil

	////先AES加密
	//encrypt,_ := utils.AesEncrypt([]byte("123456"), key)
	//// 再异或运算加密
	//for i,item := range encrypt {
	//	encrypt[i] = item ^ 28
	//}
	//fmt.Println("DH-encrypt:", encrypt)
	//
	//// 先异或运算解密
	//for i,item := range encrypt {
	//	encrypt[i] = item ^ 28
	//}
	//// 再AES解密
	//decrypt,_ := utils.AesDecrypt(encrypt, key)
	//fmt.Println("decrypt:", string(decrypt))

}
