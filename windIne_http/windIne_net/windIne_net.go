package windIne_net

/*
Package windIne_net 网络工具
*/
import (
	"github.com/mxkcw/windIne/windIne_string"
	"net"
	"os/exec"
	"strings"
)

// WindIneGetLocalIPV4WithCurrentActive 获取当前活动网卡的IPV4地址
func WindIneGetLocalIPV4WithCurrentActive() string {
	// 尝试连接到公共地址但不发送数据
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// WindIneGetPublicIPV4 获取公网IP
func WindIneGetPublicIPV4() string {
	curl := exec.Command("curl", "https://ipinfo.io/ip")
	out, err := curl.Output()
	if err != nil {
		return ""
	}
	aStr := windIne_string.WindIneBytes2String(out)
	return aStr
}

// WindIneGetRandomTag 获取基于公网IP的随机字符
func WindIneGetRandomTag() string {
	if sArr := strings.Split(WindIneGetPublicIPV4(), "."); len(sArr) == 4 {
		ak := sArr[len(sArr)-1] + "x" + sArr[0]
		return ak
	}
	return ""
}
