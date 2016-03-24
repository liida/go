package main
import (
	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"net/http"
	"os"
	"time"
	"penkr"
)

var globalSessions *session.Manager

type guide struct {
	penkr.Guide
}

type guideifc interface {
	penkr.GuideIfc
	makeEncryptPassword(password string) string
}

//密码
func (g *guide) makeEncryptPassword(password string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password))
	s := hex.EncodeToString(md5Ctx.Sum(nil))
	bs := []byte(s)
	var rs []byte
	for i := len(bs);i>0;i--{
		rs = append(rs,bs[i-1])
	}
	md5Ctx2 := md5.New()
	md5Ctx2.Write([]byte(string(rs)+g.Guide.Salt))
	return hex.EncodeToString(md5Ctx2.Sum(nil))
}

//欢迎页
func hello(w http.ResponseWriter, r *http.Request) {
	i := checkLogin(w,r)
	if i == 0 {
		http.Redirect(w,r,"/login",302)
		return
	}
	http.Redirect(w,r,"/guides",302)
}
//登录程序
func login(w http.ResponseWriter, r *http.Request) {
	writelog("test.log","login")
	defer func(){
		if err:=recover();err!=nil{
		    //fmt.Println(err)
		}
	}()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("/home/wwwroot/default/go/bin/login.gtpl")
		t.Execute(w, nil)
	}else{
		r.ParseForm()  
		if r.Form.Get("username") == "" {
			fmt.Fprintf(w, "登录失败,手机号不能为空!")
			panic("登录失败,手机号不能为空!")
		}
		if r.Form.Get("password") == "" {
			fmt.Fprintf(w, "登录失败,密码不能为空!")
			panic("登录失败,密码不能为空!")
		}
		var guide guideifc = new(guide)
		guide.InitGuideByUserName(r.Form.Get("username"))
		if guide.GetField("Guide_identity").(int64) <= 0 {
			fmt.Fprintf(w, "%s,登录失败,账号不存在!",r.Form.Get("username"))
		}else{
			if guide.GetField("Password").(string) != guide.makeEncryptPassword(r.Form.Get("password")){
				fmt.Fprintf(w, "%s,登录失败,密码错误!",guide.GetField("Guide_name"))
			}else{
				sess := globalSessions.SessionStart(w, r)
				sess.Set("Guide_identity", guide.GetField("Guide_identity").(int64))
				http.Redirect(w,r,"/guides",302)
			}
			
		}
	}
	
    
}

//获取所有导购
func guides(w http.ResponseWriter, r *http.Request) {
	defer func(){
		if err:=recover();err!=nil{
		    //fmt.Println(err)
		}
	}()
	i := checkLogin(w,r)
	if i == 0 {
		http.Redirect(w,r,"/login",302)
		return
	}
	var guide guideifc = new(guide)
	for _,v := range guide.GetGuides(){
		fmt.Fprintf(w, "%d:%s:%s:%s\n",v.Guide_identity,v.Guide_name,v.Password,v.Salt)
		
	}
	fmt.Fprintf(w, "%d:%s:%s:%s\n",guide.GetField("Guide_identity"),guide.GetField("Guide_name"),guide.GetField("Password"),guide.GetField("Salt"))
	
	
   
}

func writelog(file string,content string) bool{
	s := time.Unix(time.Now().Unix(),0).Format("2006/01/02")
	path := "/data/logs/go/"+s
	_,staterr := os.Stat(path)
	//不存在,创建
	if staterr != nil {
		//fmt.Println(path+"不存在，创建")
		mkerr := os.MkdirAll(path, 0777)
		if mkerr != nil {
			return false
		} 
	} 
	fi,openerr := os.OpenFile(path+"/"+file,os.O_WRONLY|os.O_APPEND,0666)
	if openerr != nil && os.IsNotExist(openerr) {
		//fmt.Println(path+"/"+file+"不存在，创建")
		fout,cerr := os.Create(path+"/"+file)
		defer fout.Close()
		if cerr == nil {
			fmt.Println(path+"/"+file+"创建成功，写内容")
			_,werr := fout.WriteString(content+"\r\n")
			if werr == nil {
				return false
			}else{
				fmt.Println(werr)
			}
			return true
		}
	}else{
		defer fi.Close()
		//fmt.Println(path+"/"+file+"存在，写内容")
		_,werr := fi.WriteString(content+"\r\n")
		if werr == nil {
			return false
		}else{
			fmt.Println(werr)
		}
		return true
	}
	return false
}

func readlog(){
	
	fi,openerr := os.Open("/data/logs/penkr.debug.Log")
	if openerr != nil && os.IsNotExist(openerr) {
		writelog("test.log","文件不存在")
	}else{
		defer fi.Close()
		
	}
}

func checkLogin(w http.ResponseWriter, r *http.Request) int64{
	sess := globalSessions.SessionStart(w, r)
	sessionval := sess.Get("Guide_identity")
	if sessionval == nil {
		return 0
	}
	return sessionval.(int64)
}

func logout(w http.ResponseWriter, r *http.Request){
	globalSessions.SessionDestroy(w, r)
	http.Redirect(w,r,"/",302)
}

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
	writelog("test.log","session start")
}

func main() {
	http.HandleFunc("/", hello)       //默认页
	http.HandleFunc("/login", login)         //登录
	http.HandleFunc("/guides", guides)         //获取所有用户
	http.HandleFunc("/logout", logout)         //退出登录
	err := http.ListenAndServe(":9090", nil) 
	if err != nil {
		//log.Fatal("ListenAndServe: ", err)
	}
}