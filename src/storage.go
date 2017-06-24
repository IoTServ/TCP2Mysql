package main

import ("runtime"
	"net"
	"fmt"
	"time"
	"strings"
	"mcunode/models"
	"github.com/astaxie/beego/orm"
)
var (
	mcu map[string]net.Conn
				)
func runstoragehandler(l net.Listener)  {
	for {	clientconnn, _ := l.Accept()
		var recv []byte = make([]byte, 10240)
		n, _ := clientconnn.Read(recv)
		id := string(recv[:n])
		go storagehandler(id, clientconnn)
	}
}

func storagehandler(id string,conn net.Conn) {
	fmt.Printf("storage-mcu-id:"+id)
	if v, ok := mcu[id]; ok {
		v.Close()
		delete(mcu,id)
	}
	mcu[id]=conn
	o := orm.NewOrm()
	for {
		var recv []byte = make([]byte, 10240)
		if conn.SetReadDeadline(time.Now().Add(time.Minute*60))!=nil{
			conn.Close()
			fmt.Printf("设置超时退出！")
			runtime.Goexit()
		}
		n,e:=conn.Read(recv)
		if e!=nil{
			conn.Close()
			//	delete(mcu,id)
			fmt.Printf("读取错误退出！")
			runtime.Goexit()
		}
		if conn.SetReadDeadline(time.Time{}) !=nil{
			conn.Close()
			fmt.Printf("取消超时退出！")
			runtime.Goexit()
		}
		msg:=recv[:n]
		if string(msg)!="" && string(msg)!="<h1></h1>" {
			str := strings.Replace(string(msg), "<h1></h1>", "", -1)

			user := models.DEVICE{MCUID: id , DATA:str , DATE:time.Now()}

			// insert
			id, _ := o.Insert(&user)
			fmt.Printf("Insert "+string(id)+" DATA: "+str)
		}
	}
}

func main()  {
	fmt.Printf("BY McuNode.com,TCP to MYSQL Storage!")
	fmt.Printf("USE 8002 port to connect ,First, send your id ONCE，then SEND your DATE\n")
	fmt.Printf("First TIME USE,creat MySQL DATABASE (named:mcunode,user:root,password:root) @ THIS SERVER!!!\n")
	runtime.GOMAXPROCS(runtime.NumCPU())
	mcu = make(map[string]net.Conn)
	l, _ := net.Listen("tcp", ":8002")
	for {	clientconnn, _ := l.Accept()
		var recv []byte = make([]byte, 10240)
		n, _ := clientconnn.Read(recv)
		id := string(recv[:n])
		go storagehandler(id, clientconnn)
	}
}
