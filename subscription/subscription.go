package subscription

import (
    "fmt"
    "database/sql"
	_ "github.com/mattn/go-sqlite3"
    "strconv"
    "net/http"
    "io/ioutil"
    "strings"
    "bufio"
    "../base64d"
)


//读取订阅链接(数据库)
func Get_subscription_link(sql_db_path string)[]string{
    db,err := sql.Open("sqlite3",sql_db_path)
    if err!=nil{
        fmt.Println(err)
    }
    defer db.Close()
    rows,err := db.Query("SELECT link FROM subscription_link;")
    if err != nil{
        fmt.Println(err)
    }

    var subscription_link[] string
    for rows.Next(){
        var link string
        rows.Scan(&link)
        subscription_link = append(subscription_link,link)
    }
    return subscription_link
}

//初始化订阅连接数据库
func Subscription_link_init(sql_db_path string){
    db,err := sql.Open("sqlite3",sql_db_path)
    if err!=nil{
        fmt.Println(err)
        return
    }
    //关闭数据库
    defer db.Close()
    //创建表
    db.Exec("CREATE TABLE IF NOT EXISTS subscription_link(link TEXT);")
}

//添加订阅链接
func Subscription_link_add(subscription_link,sql_db_path string){
    db,err := sql.Open("sqlite3",sql_db_path)
    if err!=nil{
        fmt.Println(err)
        return
    }
    defer db.Close()
    db.Exec("INSERT INTO subscription_link(link)values(?)",subscription_link)
}

//删除订阅链接(数据库)
func Subscription_link_delete(sql_db_path string){
    var subscription_link[] string
    db,err := sql.Open("sqlite3",sql_db_path)
    if err!=nil{
        fmt.Println(err)
        return
    }
    defer db.Close()

    rows,err := db.Query("SELECT link FROM subscription_link")
    var link string
    for rows.Next(){
        err = rows.Scan(&link)
        subscription_link = append(subscription_link,link)
    }
    //fmt.Println(subscription_link)
    for num,link_temp := range subscription_link{
        fmt.Println(strconv.Itoa(num+1)+"."+link_temp)
    }
    var select_delete int
    fmt.Scanln(&select_delete)
    if select_delete>=1&&select_delete<=len(subscription_link){
        db.Exec("DELETE FROM subscription_link WHERE link = ?",subscription_link[select_delete-1])
    }else{
        fmt.Println("enter error,please retry.")
        Subscription_link_delete(sql_db_path)
        return
    }
}

//更新订阅
func Http_get_subscription(url string)string{
    res,_ := http.Get(url)
    body,err := ioutil.ReadAll(res.Body)
    if err!=nil{
        fmt.Println(err)
        fmt.Println("可能出错原因,请检查能否成功访问订阅连接.")
    }
    return string(body)
    //ioutil.WriteFile(read_config().config_path,[]byte(body),0644)
}


//方便进行分割对字符串进行替换
func str_replace(str string)[]string{
    var config[] string
    scanner := bufio.NewScanner(strings.NewReader(strings.Replace(base64d.Base64d(str),"ssr://","",-1)))
    for scanner.Scan() {
    str_temp := strings.Replace(base64d.Base64d(scanner.Text()),"/?obfsparam=",":",-1)
    str_temp = strings.Replace(str_temp,"&protoparam=",":",-1)
    str_temp = strings.Replace(str_temp,"&remarks=",":",-1)
    str_temp = strings.Replace(str_temp,"&group=",":",-1)
    config = append(config,str_temp)
    }
    return config
}


//更新订阅(sqlite数据库)
func Update_config_db(str_2 []string,sql_db_path string){


    //访问数据库
    db,err := sql.Open("sqlite3",sql_db_path)
    if err!=nil{
        fmt.Println(err)
        return
    }

    defer db.Close()

    //删除表
    db.Exec("DROP TABLE IF EXISTS SSR_info;")

    //创建表
     sql_table := `
    CREATE TABLE IF NOT EXISTS SSR_info(
        id TEXT,
        remarks TEXT,
        server TEXT,
        server_port TEXT,
        protocol TEXT,
        method TEXT,
        obfs TEXT,
        password TEXT,
        obfsparam TEXT,
        protoparam TEXT
    );
    `
    db.Exec(sql_table)
    
    //config_middle_temp := str_replace(string(read_ssr_config()))
    //list_list(config_middle_temp)
    for num,config_temp := range str_2{
        config_split := strings.Split(config_temp,":")
        var server string
        if len(config_split) == 17 {
            server = config_split[0]+":"+config_split[1]+":"+config_split[2]+":"+config_split[3]+":"+config_split[4]+":"+config_split[5]+":"+config_split[6]+":"+config_split[7]
        } else if len(config_split) == 10 {
            server = config_split[0]
        }
        server_port := config_split[len(config_split)-9]
        protocol := config_split[len(config_split)-8]
        method := config_split[len(config_split)-7]
        obfs := config_split[len(config_split)-6]
        password := base64d.Base64d(config_split[len(config_split)-5])
        obfsparam := base64d.Base64d(config_split[len(config_split)-4])
        protoparam := base64d.Base64d(config_split[len(config_split)-3])
        remarks := base64d.Base64d(config_split[len(config_split)-2])
        //fmt.Println(num,remarks,server,server_port,protocol,method,obfs,password,obfsparam,protoparam)



        //向表中插入数据
        db.Exec("INSERT INTO SSR_info(id,remarks,server,server_port,protocol,method,obfs,password,obfsparam,protoparam)values(?,?,?,?,?,?,?,?,?,?)",num+1,remarks,server,server_port,protocol,method,obfs,password,obfsparam,protoparam)
        //stmt,_ := db.Prepare("INSERT INTO SSR_info(id,remarks,server,server_port,protocol,method,obfs,password,obfsparam,protoparam)values(?,?,?,?,?,?,?,?,?,?)")
        //res,_ := stmt.Exec(num+1,remarks,server,server_port,protocol,method,obfs,password,obfsparam,protoparam)
        //id,_ := res.LastInsertId()
        //fmt.Println(id)
    }
}