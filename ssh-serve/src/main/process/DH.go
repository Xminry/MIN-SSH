package process

import (
	"fmt"
	"math/big"
	"math/rand"
	_struct "ssh-serve/src/main/struct"
	"strconv"
)

// DH密钥协商算法,底数P，模板G
func Dh(p int,g int,sshReaderWriter *_struct.SshReaderWriter) ([]byte,error){
	// p为16位质数
	//p := big.NewInt(7000487615661733)
	//g := big.NewInt(5925845745820835)
	P := big.NewInt(int64(p))
	G := big.NewInt(int64(g))

	// EXP(a, b, c) = (a ** b) % c
	//rand.Seed(time.Now().Unix())    一边设置种子，一边不设置，使得随机数不可预测
	targer := rand.Intn(900) + 100
	//私钥a
	a := big.NewInt(int64(targer))
	A := big.NewInt(0).Exp(G, a, P)
	fmt.Println("生成服务端公钥：",  A)
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
	fmt.Println("收到客户端公钥：",B)
	//发送公钥至客户端
	sshReaderWriter.Write([]byte(A.String()))


	// 服务器拿到客户端的公钥，生成密钥K
	K := big.NewInt(0).Exp(big.NewInt(int64(B)), a, P)
	key := []byte(K.String())
	//将16位key拼接成32位会话密钥
	key = append(key, key...)
	fmt.Println("生成会话密钥：", key)
	return key,nil

	// 先AES加密
	//encrypt,_ := utils.AesEncrypt([]byte("123456"), key)
	//// 再异或运算加密
	//for i,item := range encrypt {
	//	encrypt[i] = item ^ 28
	//}
	//fmt.Println("encrypt:", encrypt)
	//
	//// 先异或运算解密
	//for i,item := range encrypt {
	//	encrypt[i] = item ^ 28
	//}
	//// 再AES解密
	//decrypt,_ := utils.AesDecrypt(encrypt, key)
	//fmt.Println("decrypt:", string(decrypt))

}
