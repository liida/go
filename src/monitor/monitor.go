package main
import (
	"fmt"
	"os"
	"strings"
	"os/exec"
	"time"
)


var bin	string
var run_script	string	
var log_file	string	


func init(){
	if len(os.Args) > 3 {
		bin = os.Args[1]
		run_script = os.Args[2]
		log_file = os.Args[3]
	}else{
		fmt.Println("进程文件范例\n",
			"./monitor php chat_queue.php /data/logs/chat_queue.log >>/data/logs/monitor.log\n")
		os.Exit(1)
	}	
}


func main() {
	for ; ; {
		cmd := exec.Command("sh","-c","ps -ef|grep -v grep|grep "+run_script+"|awk '{a[NR]=$0}END{for(i=1;i<NR;i++)print a[i]}'|awk '{print $2}'")
		r, err2 := cmd.CombinedOutput()
		if err2 != nil {
			//fmt.Println(err2)
			os.Exit(1)
		}else{
			i := 0
			if string(r) != "" {
				for _,pid := range strings.Split(string(r), "\n"){
					if pid !="" {
						if i >= 1 {
							//垃圾进程，杀死
							fmt.Println("发现垃圾进程"+pid+"，杀死")
							cmd := exec.Command("sh","-c","kill -9 "+pid)
							_, err4 := cmd.CombinedOutput()
							if err4 != nil {
								fmt.Println(err4)
							}
						}
						i ++
					}

				} 
			}
			if i < 1 {
				fmt.Println("启动进程"+bin+" "+run_script+" >>"+log_file)
				c := exec.Command("sh","-c",bin+" "+run_script+" >>"+log_file)
				go c.Run()
				
				
			}
		
		}
		//fmt.Println("sleep")
		time.Sleep(3*time.Second);
	}
}


