package communicate

import (
	"net"
)

func Read(con net.Conn, buf []byte) (int, error) {
	n, err := con.Read(buf)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func Write(con net.Conn, buf []byte) (int, error) {
	n, err := con.Write(buf)
	if err != nil {
		return 0, err
	}
	return n, nil
}
