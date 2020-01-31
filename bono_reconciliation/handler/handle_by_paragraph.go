package handler

import (
    "github.com/tealeg/xlsx"
    "log"
    "strconv"
    "strings"
    "fmt"
    // "os"
)

var _= fmt.Println

const (
    bank_sheet_index int = 1
    sap_sheet_index int = 0
    bank_head_index int = 4
    sap_head_index int = 0
    sheet_max int = 2
    /*
    ** 银行借方
    */
    bank_debit_index int = 8
    /*
    ** 银行贷方
    */
    bank_lender_index int = 9
    /*
    ** 银行序号
    */
    bank_no int = 0
    bank_col_max int = 10
    /*
    ** sap借方
    */
    sap_debit_index int = 11
    /*
    ** sap贷方
    */
    sap_lender_index int = 12
    /*
    ** sap打印序号
    */
    sap_printno_index int = 6
    /*
    ** sap制单人
    */
    sap_prepared_by_index int = 8
    /*
    ** sap凭证
    */
    sap_cert_id_index int = 5
    sap_col_max int = 13
)

type CBankData struct {
    debitKey float64
    debitValue float64
    lenderKey float64
    lenderValue float64
    no string
}

type CSapData struct {
    debitKey float64
    debitValue float64
    lenderKey float64
    lenderValue float64
    printNo int64
    preparedby string
    certId string
}

/*
** 通过 借方 和 贷方进行比对
*/
type CByParagraph struct {
    bankPath string
    sapPath string
}

func (self *CByParagraph) writeXlsx(file *xlsx.File, sheetName string, header *[]string, rows *[]*[]string) error {
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
    // var file *xlsx.File
    // file = xlsx.NewFile()
    // err = file.Save(fileName)
    // if err != nil {
    //     return err
    // }
    return nil
}

func (self *CByParagraph) Calc() error {
    bankData, err := self.readBankData()
    if err != nil {
        return err
    }
    sapData, err := self.readSapData()
    if err != nil {
        return err
    }
    /*
    ** key: 银行借方
    */
    bankLenderMap := map[float64]*[]CBankData{}
    for _, bank := range *bankData {
        if bank.lenderValue == 0 || (bank.debitValue != 0 && bank.lenderValue != 0) {
            continue
        }
        if v, ok := bankLenderMap[bank.lenderKey]; ok {
            *v = append(*v, bank)
        } else {
            vec := []CBankData{
                bank,
            }
            bankLenderMap[bank.lenderKey] = &vec
        }
    }
    sapDebitMap := map[float64]*[]CSapData{}
    for _, sap := range *sapData {
        if sap.debitValue == 0 || (sap.debitValue != 0 && sap.lenderValue != 0) {
            continue
        }
        if v, ok := sapDebitMap[sap.debitKey]; ok {
            *v = append(*v, sap)
        } else {
            vec := []CSapData{
                sap,
            }
            sapDebitMap[sap.debitKey] = &vec
        }
    }
    // fmt.Println(sapDebitMap)
    /*
    ** 遍历银行借方, 查找sap贷方是否存在
    */
    bankExistSapNotexist := []*[]string{}
    bankExistSapNotexistHeader := []string{
        "借方", "贷方", "序号",
    }
    for _, bank := range *bankData {
        if bank.debitValue == 0 {
            continue
        }
        if v, ok := sapDebitMap[bank.debitKey]; ok {
            /*
            ** 银行借方 和 sap贷方都存在
            */
            /*
            ** 从sapLenderMap中删除都有的, 则剩下的就是sap有的, 但是银行没有的
            ** 1. 当vec为空时, 从map中移除
            */
            if len(*v) == 0 {
                delete(sapDebitMap, bank.debitKey)
                continue
            }
            *v = (*v)[1:]
            /*
            for i, va := range *v {
                if v.certId == va.certId {
                    v2 := (*v)[i+1:]
                    *v = (*v)[0:i]
                    *v = append(*v, v2)
                    break
                }
            }
            */
        } else {
            /*
            ** 银行借方存在, sap贷方不存在
            */
            // fmt.Println("银行借方存在, 但是sap贷方不存在:", bank)
            vs := []string{
                strconv.FormatFloat(bank.debitValue, 'f', 2, 64),
                strconv.FormatFloat(bank.lenderValue, 'f', 2, 64),
                bank.no,
            }
            bankExistSapNotexist = append(bankExistSapNotexist, &vs)
        }
    }
    // fmt.Println("-----------sap贷方存在, 但是银行借方不存在-----------")
    // f, err := os.OpenFile("sqp贷方存在_银行借方不存在.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil {
        return err
    }
    // defer f.Close()
    sapExistBankNotexistHeader := []string{
        "借方", "贷方", "打印序号", "制单人编号",
    }
    sapExistBankNotexist := []*[]string{}
    for _, value := range sapDebitMap {
        /*
        ** sap贷方存在, 银行借方不存在
        */
        for _, v := range *value {
            // s := fmt.Sprintf("%v, %v, %v, %s", v.debitValue, v.lenderValue, v.printNo, v.preparedby)
            // fmt.Println(key, s)
            // f.Write([]byte(s))
            vs := []string{
                strconv.FormatFloat(v.lenderValue, 'f', 2, 64),
                strconv.FormatFloat(v.debitValue, 'f', 2, 64),
                strconv.FormatInt(int64(v.printNo), 10),
                v.preparedby,
            }
            sapExistBankNotexist = append(sapExistBankNotexist, &vs)
        }
    }
    // self.writeXlsx("sap存在_银行不存在.xlsx", "sheet1", &sapLenderExistBankDebitNotexistHeader, &sapLenderExistBankDebitNotexist)
    // self.writeXlsx("银行存在_sap不存在.xlsx", "sheet1", &bankExistSapNotexistHeader, &bankExistSapNotexist)
    /*
    ** 遍历sap, 查找银行是否存在
    */
    // bankLenderExistSapDebitNotexist := []*[]string{}
    // bankLenderExistSapDebitNotexistHeader := []string{
    //     "借方", "贷方", "打印序号", "制单人编号",
    // }
    for _, sap := range *sapData {
        if sap.lenderValue == 0 {
            continue
        }
        if v, ok := bankLenderMap[sap.lenderKey]; ok {
            /*
            ** 银行贷方 和 sap借方都存在
            */
            /*
            ** 从bankDebitMap中删除都有的, 则剩下的就是sap有的, 但是银行没有的
            */
            if len(*v) == 0 {
                delete(bankLenderMap, sap.lenderKey)
                continue
            }
            *v = (*v)[1:]
        } else {
            /*
            ** 银行贷方存在, sap借方不存在
            */
            // fmt.Println("银行借方存在, 但是sap贷方不存在:", bank)
            vs := []string{
                strconv.FormatFloat(sap.lenderValue, 'f', 2, 64),
                strconv.FormatFloat(sap.debitValue, 'f', 2, 64),
                strconv.FormatInt(int64(sap.printNo), 10),
                sap.preparedby,
            }
            sapExistBankNotexist = append(sapExistBankNotexist, &vs)
        }
    }
    // fmt.Println("-----------sap贷方存在, 但是银行借方不存在-----------")
    // f, err := os.OpenFile("sqp贷方存在_银行借方不存在.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil {
        return err
    }
    // defer f.Close()
    // sapDebitExistBankLenderNotexistHeader := []string{
    //     "借方", "贷方", "序号",
    // }
    // sapDebitExistBankLenderNotexist := []*[]string{}
    for _, value := range bankLenderMap {
        /*
        ** sap贷方存在, 银行借方不存在
        */
        for _, v := range *value {
            // s := fmt.Sprintf("%v, %v, %v, %s", v.debitValue, v.lenderValue, v.printNo, v.preparedby)
            // fmt.Println(key, s)
            // f.Write([]byte(s))
            vs := []string{
                strconv.FormatFloat(v.debitValue, 'f', 2, 64),
                strconv.FormatFloat(v.lenderValue, 'f', 2, 64),
                v.no,
            }
            bankExistSapNotexist = append(bankExistSapNotexist, &vs)
        }
    }
    // self.writeXlsx("银行存在_sap不存在.xlsx", "sheet1", &bankExistSapNotexistHeader, &bankExistSapNotexist)
    // self.writeXlsx("sap存在_银行不存在.xlsx", "sheet1", &sapExistBankNotexistHeader, &sapExistBankNotexist)
    var file *xlsx.File
    file = xlsx.NewFile()
    self.writeXlsx(file, "银行存在, sap不存在", &bankExistSapNotexistHeader, &bankExistSapNotexist)
    self.writeXlsx(file, "sap存在, 银行不存在", &sapExistBankNotexistHeader, &sapExistBankNotexist)
    err = file.Save("resource/output.xlsx")
    if err != nil {
        return err
    }
    return nil
}

func (self *CByParagraph) readBankData() (*[]CBankData, error) {
    f, err := xlsx.OpenFile(self.bankPath)
    if err != nil {
        log.Printf("readExcel error, path: %s, err: %v\n", self.bankPath, err)
        return nil, err
    }
    // defer f.Close()
    if len(f.Sheets) < sheet_max {
        log.Println("sheet less than sheet_max")
        return nil, err
    }
    sheet := f.Sheets[bank_sheet_index]
    datas := []CBankData{}
    for index, row := range sheet.Rows {
        cells := row.Cells
        if index <= bank_head_index {
            /*
            ** 去除前5行
            */
            continue
        }
        if len(cells) < bank_col_max {
            continue
        }
        debitStr := strings.Trim(cells[bank_debit_index].String(), " ")
        debitStr = strings.Replace(debitStr, "(", "-", -1)
        debitStr = strings.Replace(debitStr, ")", "", -1)
        debit, err := strconv.ParseFloat(debitStr, 64)
        if err != nil {
            log.Println(err)
            continue
        }
        debitKey := debit
        if debitKey < 0 {
            debitKey = -debitKey
        }
        lenderStr := strings.Trim(cells[bank_lender_index].String(), " ")
        lenderStr = strings.Replace(lenderStr, "(", "-", -1)
        lenderStr = strings.Replace(lenderStr, ")", "", -1)
        lender, err := strconv.ParseFloat(lenderStr, 64)
        if err != nil {
            log.Println(err)
            continue
        }
        lenderKey := lender
        if lenderKey < 0 {
            lenderKey = -lenderKey
        }
        bankNo := cells[bank_no].String()
        datas = append(datas, CBankData{
            debitKey: debitKey,
            debitValue: debit,
            lenderKey: lenderKey,
            lenderValue: lender,
            no: bankNo,
        })
    }
    return &datas, nil
}

func (self *CByParagraph) readSapData() (*[]CSapData, error) {
    f, err := xlsx.OpenFile(self.sapPath)
    if err != nil {
        log.Printf("readExcel error, path: %s, err: %v\n", self.sapPath, err)
        return nil, err
    }
    // defer f.Close()
    if len(f.Sheets) < sheet_max {
        log.Println("sheet less than sheet_max")
        return nil, err
    }
    sheet := f.Sheets[sap_sheet_index]
    datas := []CSapData{}
    for index, row := range sheet.Rows {
        cells := row.Cells
        if index <= sap_head_index {
            /*
            ** 去除前5行
            */
            continue
        }
        if len(cells) < sap_col_max {
            continue
        }
        debitStr := strings.Trim(cells[sap_debit_index].String(), " ")
        debitStr = strings.Replace(debitStr, "(", "-", -1)
        debitStr = strings.Replace(debitStr, ")", "", -1)
        debit, err := strconv.ParseFloat(debitStr, 64)
        if err != nil {
            log.Println(err)
            continue
        }
        debitKey := debit
        if debitKey < 0 {
            debitKey = -debitKey
        }
        lenderStr := strings.Trim(cells[sap_lender_index].String(), " ")
        lenderStr = strings.Replace(lenderStr, "(", "-", -1)
        lenderStr = strings.Replace(lenderStr, ")", "", -1)
        lender, err := strconv.ParseFloat(lenderStr, 64)
        if err != nil {
            log.Println(err)
            continue
        }
        lenderKey := lender
        if lenderKey < 0 {
            lenderKey = -lenderKey
        }
        printNoStr := strings.Trim(cells[sap_printno_index].String(), " ")
        printNo, err := strconv.ParseInt(printNoStr, 10, 64)
        if err != nil {
            log.Println(err)
            continue
        }
        preparedby := cells[sap_prepared_by_index].String()
        certId := cells[sap_cert_id_index].String()
        datas = append(datas, CSapData{
            debitKey: lenderKey,
            debitValue: lender,
            lenderKey: debitKey,
            lenderValue: debit,
            printNo: printNo,
            preparedby: preparedby,
            certId: certId,
        })
    }
    return &datas, nil
}

func NewByParagraph(bankPath string, sapPath string) *CByParagraph {
    obj := CByParagraph{
        bankPath: bankPath,
        sapPath: sapPath,
    }
    return &obj
}

