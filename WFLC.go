package main

//導入函式庫
import (
	"bufio"
	"database/sql" // 匯入資料庫元件
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	_ "github.com/go-sql-driver/mysql" // 匯入從Git 下載的MySQL function
)

//函數宣告
var (
	db          = &sql.DB{}
	valueReader = ""

	mesCnt         = "" //Message content
	mesTriggerTime = "" //Message triggered time
	hostIP         = "" //Computer IP
	computerUUID   = "" //Computer universally unique identifier

)

func init() {
	//db,_ = sql.Open("mysql", "[Username]:[Password]@tcp([IP:Port])/logcollector")//alanjian-big-pc 外部IP
}

//Main Fuction
func main() {
	GetDataFromPowerbuilder()
	GetMesTriggerTime()
	GetHostIP()
	GetComputerUUID()
	PushValueToServer()
}

func PushValueToServer() {
	tx, _ := db.Begin()
	tx.Exec("INSERT INTO logcollector(MesCnt, MesTriggerTime, HostIP, ComputerUUID) values(?,?,?,?)", mesCnt, mesTriggerTime, hostIP, computerUUID)
	//最后释放tx内部的连接
	tx.Commit()
}

//Get current system time
func GetMesTriggerTime() {
	mesTriggerTime = time.Now().Format("2006-01-02 15:04:05.00000")
}

//get computer IP
func GetHostIP() {
	_, hostIP = GetIPv4()
}

//get computer UUID
func GetComputerUUID() {
	id, err := machineid.ProtectedID("myAppName")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}


func GetIPv4() (string, string) {
	var getIPv4 string
	var getifaceType string
	list, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range list {
		//fmt.Printf("%d name=%s %v\n", i, iface.Name, iface)
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			//fmt.Printf(" %d %v\n", j, addr)
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && iface.Name == "乙太網路" {
				if ipnet.IP.To4() != nil {
					getifaceType = iface.Name
					getIPv4 = ipnet.IP.String()
					//fmt.Println("乙太網路, getIPv4:" + getIPv4)
					return getifaceType, getIPv4
				}
			} else if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && iface.Name == "區域連線" && getIPv4 == "" {
				if ipnet.IP.To4() != nil {
					getifaceType = iface.Name
					getIPv4 = ipnet.IP.String()
					//fmt.Println("區域連線, getIPv4:" + getIPv4)
					return getifaceType, getIPv4
				}
			} else if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && strings.ToLower(iface.Name) == "wi-fi" && getIPv4 == "" {
				if ipnet.IP.To4() != nil {
					getifaceType = iface.Name
					getIPv4 = ipnet.IP.String()
					//fmt.Println("wi-fi, getIPv4:" + getIPv4)
					return getifaceType, getIPv4
				}
			}
		}
	}

	if getIPv4 == "" {
		secondAddrs, err := net.InterfaceAddrs()
		if err != nil {
			os.Stderr.WriteString("Oops: " + err.Error() + "\n")
			os.Exit(1)
		}
		for _, a := range secondAddrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					getIPv4 = ipnet.IP.String()
				}
			}
		}
	}
	return getifaceType, getIPv4
}

//get the alert message content from CPOE
func GetDataFromPowerbuilder() {

	//參數初始化
	var mesCntInput string
	flag.StringVar(&mesCntInput, "mesCntInput", "null", "訊息內容")

	//解析cmd傳入參數
	flag.Parse()

	//將傳入參數指定給原本的全域變數
	mesCnt = mesCntInput

}


//go build -ldflags "-H windowsgui" WFLC.go
