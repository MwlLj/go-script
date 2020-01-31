package main

import (
    "main/export"
    "flag"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    host := flag.String("host", "127.0.0.1", "host")
    port := flag.Int("port", 3306, "port")
    userName := flag.String("user-name", "root", "user name")
    userPwd := flag.String("user-pwd", "123", "user pwd")
    dbName := flag.String("db-name", "", "db name")
    prjName := flag.String("prj-name", "", "project name")
    flag.Parse()
    obj := export.NewStaffInsfrastructure()
    obj.SetDbInfo(&export.CDbInfo{
        Host: *host,
        Port: *port,
        UserName: *userName,
        UserPwd: *userPwd,
    })
    obj.SetDbName(dbName)
    obj.SetProjectName(prjName)
    obj.Init()
    if err := obj.Get(); err != nil {
        log.Println("get error, err:", err)
        return
    }
    log.Println("success")
}
