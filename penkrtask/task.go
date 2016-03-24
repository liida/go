package main
import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/robfig/cron"
	"os/exec"
	"time"
)

var service string


type conf struct {
	key	string
	val	string
}
var confs []conf

type task struct {
	id	int
	title	string
	task_type int
	desc	string
	rule	string
	bin	string
	run_script	string
	status	int
}

func ReadConf(file string) {
	f, err := os.Open(file)
	if err == nil {
		defer f.Close()
		content,readerr := ioutil.ReadAll(f)
		if readerr == nil {
			var g conf
			for _,line := range strings.Split(string(content), "\n"){
				if line != "" {
					for key,val := range strings.Split(line,"="){
						if key == 0 {
							g.key = strings.Trim(val," ")
						}else if key == 1 {
							g.val = strings.Trim(val," ")
						}
					}
					if g.key != "" || g.val != "" {
						confs = append(confs,g)
					}
				}
			} 
		}
	}else{
		fmt.Println("配置文件不存在",file)
		os.Exit(1)
	}
}

func getConf(key string) string {
	for _,conf := range confs {
		if conf.key == key {
			return conf.val
		}
	}
	return ""
}


func (t *task) run() bool{
	cmd := exec.Command("sh","-c","ps -ef|grep -v grep|grep "+t.run_script+"|awk '{print $2}'")
	r, err2 := cmd.CombinedOutput()
	if err2 != nil {
		fmt.Println(err2)
		return false
	}else{
		if string(r) != "" {
			for _,pid := range strings.Split(string(r), "\n"){
				if pid !="" {
					//垃圾进程，杀死
					fmt.Println("发现垃圾进程"+pid+"，杀死")
					cmd := exec.Command("sh","-c","kill -9 "+pid)
					_, err4 := cmd.CombinedOutput()
					if err4 != nil {
						fmt.Println(err4)
						return false
					}
				}

			} 
		}
	}
	_,staterr := os.Stat(t.run_script)
	if staterr != nil{
		fmt.Println(staterr)
		return false
	}
	fmt.Println(time.Now().String(),t.bin, t.run_script)
	c := exec.Command(t.bin, t.run_script)
	_,err := c.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return false
	}else{
		return true
	}

	
}



func init(){
	if len(os.Args) > 1 {
		ReadConf(os.Args[1])
	}else{
		fmt.Println("配置文件不存在,请在后面指定配置文件，配置文件范例\n",
			"logpath = /data/logs/\n",
			"dbhost = 192.168.10.159\n",
			"dbport = 3306\n",
			"dbuser = liudan\n",
			"dbpass = liudan123\n",
			"dbname = weidaogou\n",
			"dbcharset = utf8\n")
		fmt.Println("MySQL数据库cron表结构\n",
			"CREATE TABLE `cron` (\n",
			"`id` int(11) unsigned NOT NULL AUTO_INCREMENT,\n",
			"`title` varchar(48) NOT NULL DEFAULT '' COMMENT '任务标题',\n",
			"`type` smallint(6) NOT NULL DEFAULT '1' COMMENT '1:crontab  2:守护进程',\n",
			"`desc` varchar(255) NOT NULL DEFAULT '' COMMENT '任务描述',\n",
			"`rule` varchar(128) NOT NULL DEFAULT '' COMMENT 'crontab格式，支持到秒',\n",
			"`bin` varchar(128) NOT NULL DEFAULT '' COMMENT '执行文件bin（php sh等）',\n",
			"`run_script` varchar(240) NOT NULL DEFAULT '' COMMENT '执行文件全路径',\n",
			"`status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态（0:未执行 1：已加入服务 2：已废弃）',\n",
			"PRIMARY KEY (`id`),\n",
			"KEY `status` (`status`)\n",
			") ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8\n")
		os.Exit(1)
	}
	service = os.Args[0]
	cmd := exec.Command("sh","-c","ps -ef|grep -v grep|grep "+service+"|awk '{a[NR]=$0}END{for(i=1;i<NR;i++)print a[i]}'|awk '{print $2}'")
	r, err2 := cmd.CombinedOutput()
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(1)
	}else{
		if string(r) != "" {
			for _,pid := range strings.Split(string(r), "\n"){
				if pid !="" {
					//垃圾进程，杀死
					fmt.Println("发现垃圾进程"+pid+"，杀死")
					cmd := exec.Command("sh","-c","kill -9 "+pid)
					_, err4 := cmd.CombinedOutput()
					if err4 != nil {
						fmt.Println(err4)
					}
				}

			} 
		}
	}
}


func main() {
	c := cron.New()
	db, err := sql.Open("mysql", getConf("dbuser")+":"+getConf("dbpass")+"@tcp("+getConf("dbhost")+":"+getConf("dbport")+")/"+getConf("dbname")+"?charset="+getConf("dbcharset"))
	if err == nil {
		defer db.Close()
		db.Query("update cron set status =0 where status=1")
		//查询数据
		rows, err := db.Query("SELECT id,title,`type` as task_type,`desc`,rule,`bin`,run_script,status FROM cron where status = 0")
		if err == nil {
			for rows.Next() {
				var task task
				rows.Scan(&task.id,&task.title,&task.task_type,&task.desc,&task.rule,&task.bin,&task.run_script,&task.status)
				if task.task_type == 1 { 
					c.AddFunc(task.rule, func() { 
						task.run()
					})
				}else if task.task_type == 2{
					go task.run()
				}
				db.Query("update cron set status =1 where id = ?",task.id)
			}
		}
		c.Start()
	}else{
		fmt.Println(err)
	}
}


