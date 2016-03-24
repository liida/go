package penkr
import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)

type Guide struct{
	Guide_identity	int64
	Guide_name	string
	Password	string
	Salt  string
}

type GuideIfc interface {
    InitGuideByUserName(username string)
    Display()
    GetField(key string) interface{}
    GetGuides() []Guide
}

func (g *Guide) Display() {
    fmt.Printf("%d:%s:%s:%s\n",g.Guide_identity,g.Guide_name,g.Password,g.Salt)
}

func (g *Guide) GetId() int64{
    return g.Guide_identity
}

func (g *Guide) GetField(key string) interface{}{
    if key == "Guide_identity" {
	return g.Guide_identity
    }else if key == "Guide_name" {
	return g.Guide_name
    }else if key == "Password" {
	return g.Password
    }else if key == "Salt" {
	return g.Salt
    }else{
	return ""
    }
    
}


func (g *Guide) InitGuideByUserName(username string){
	db, err := sql.Open("mysql", "liudan:liudan123@tcp(192.168.10.159:3306)/weidaogou?charset=utf8")
	if err == nil {
		defer db.Close()
		//查询数据
		rows, err := db.Query("SELECT guide_identity,guide_name,password,salt FROM guide where cellphone=? limit 1",username)
		if err == nil {
			for rows.Next() {
				rows.Scan(&g.Guide_identity,&g.Guide_name,&g.Password,&g.Salt)
				
			}
		}
	}else{
		fmt.Println(err)
	}	

}


func (g *Guide) GetGuides() []Guide{
	var GuideList []Guide
	db, err := sql.Open("mysql", "liudan:liudan123@tcp(192.168.10.159:3306)/weidaogou?charset=utf8")
	if err == nil {
		defer db.Close()
		//查询数据
		rows, err := db.Query("SELECT guide_identity,guide_name,password,salt FROM guide")
		if err == nil {
			for rows.Next() {
				rows.Scan(&g.Guide_identity,&g.Guide_name,&g.Password,&g.Salt)
				GuideList = append(GuideList,*g)
			}
		}
	}
	return GuideList

}

