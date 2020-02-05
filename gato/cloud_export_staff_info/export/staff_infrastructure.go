package export

import (
    "github.com/tealeg/xlsx"
    "fmt"
    // "os"
    "main/db"
    "errors"
)

type CDbInfo struct {
    Host string
    Port int
    UserName string
    UserPwd string
}

type CStaffInsfrastructure struct {
    dbName *string
    projectName *string
    dbInfo CDbInfo
    dbHandler *db.CDbHandler
}

func (self *CStaffInsfrastructure) SetDbInfo(dbInfo *CDbInfo) {
    self.dbInfo = *dbInfo
}

func (self *CStaffInsfrastructure) SetDbName(dbName *string) {
    self.dbName = dbName
}

func (self *CStaffInsfrastructure) SetProjectName(projectName *string) {
    self.projectName = projectName
}

func (self *CStaffInsfrastructure) Init() {
    dbHandler := &db.CDbHandler{}
    err := dbHandler.Connect(self.dbInfo.Host, uint(self.dbInfo.Port), self.dbInfo.UserName, self.dbInfo.UserPwd, *self.dbName, "mysql")
    if err != nil {
        return
    }
    if self.dbHandler != nil {
        self.dbHandler.Disconnect()
    }
    self.dbHandler = dbHandler
}

func (self *CStaffInsfrastructure) writeXlsx(file *xlsx.File, sheetName string, header *[]string, rows *[]*[]string) error {
    var sheet *xlsx.Sheet
    var row *xlsx.Row
    var cell *xlsx.Cell
    var err error
    sheet, err = file.AddSheet(sheetName)
    if err != nil {
        return err
    }
    /*
    ** 添加头
    */
    row = sheet.AddRow()
    for _, head := range *header {
        cell = row.AddCell()
        cell.Value = head
    }
    for _, r := range *rows {
        row = sheet.AddRow()
        for _, c := range *r {
            cell = row.AddCell()
            cell.Value = c
        }
    }
    return nil
}

func (self *CStaffInsfrastructure) Get() error {
    if self.dbHandler == nil {
        return errors.New(fmt.Sprintf("db init failed"))
    }
    /*
    ** 获取所有人员信息
    */
    output := []db.CGetAllStaffInfoOutput{}
    err, _ := self.dbHandler.GetAllStaffInfo(&output)
    if err != nil {
        return err
    }
    var file *xlsx.File
    file = xlsx.NewFile()
    rows := []*[]string{}
    for _, item := range output {
        /*
        ** 根据基建uuid获取基建的完整信息
        */
        var full string
        if err := self.getFullInfrastructure(item.InfrastructureUuid, &full, true); err != nil {
            return err
        }
        rows = append(rows, &[]string{
            full, item.StaffName, item.Address, item.Gender, item.ContactInfo,
        })
    }
    headers := []string{
        "房间号", "姓名", "住址", "性别", "联系信息",
    }
    err = self.writeXlsx(file, "sheet", &headers, &rows)
    if err != nil {
        return err
    }
    file.Save(*self.projectName+".xlsx")
    return nil
}

func (self *CStaffInsfrastructure) getFullInfrastructure(uuid string, full *string, first bool) error {
    if self.dbHandler == nil {
        return errors.New(fmt.Sprintf("db init failed"))
    }
    input := db.CGetInfrastructureInfoByUuidInput{
        InfrastructureUuid: uuid,
    }
    output := db.CGetInfrastructureInfoByUuidOutput{}
    err, _ := self.dbHandler.GetInfrastructureInfoByUuid(&input, &output)
    if err != nil {
        return err
    }
    if !output.ParentUuidIsValid || output.ParentUuid == "" {
        /*
        ** 父节点为空 => 为根节点 (递归结束条件)
        */
        return nil
    }
    if first {
        *full = fmt.Sprintf("%s", output.InfrastructureName)
    } else {
        *full = fmt.Sprintf("%s/%s", output.InfrastructureName, *full)
    }
    return self.getFullInfrastructure(output.ParentUuid, full, false)
}

func NewStaffInsfrastructure() *CStaffInsfrastructure {
    obj := CStaffInsfrastructure{
    }
    return &obj
}
