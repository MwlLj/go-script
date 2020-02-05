package db

import (
	"bufio"
	"bytes"
	"database/sql"
	"io"
	"os"
	"regexp"
	"strconv"
	"fmt"
	"errors"
)

type CDbHandler struct  {
	m_db *sql.DB
}

func (this *CDbHandler) Connect(host string, port uint, username string, userpwd string, dbname string, dbtype string) (err error) {
	b := bytes.Buffer{}
	b.WriteString(username)
	b.WriteString(":")
	b.WriteString(userpwd)
	b.WriteString("@tcp(")
	b.WriteString(host)
	b.WriteString(":")
	b.WriteString(strconv.FormatUint(uint64(port), 10))
	b.WriteString(")/")
	b.WriteString(dbname)
	var name string
	if dbtype == "mysql" {
		name = b.String()
	} else if dbtype == "sqlite3" {
		name = dbname
	} else {
		return errors.New("dbtype not support")
	}
	this.m_db, err = sql.Open(dbtype, name)
	if err != nil {
		return err
	}
	this.m_db.SetMaxOpenConns(2000)
	this.m_db.SetMaxIdleConns(1000)
	this.m_db.Ping()
	return nil
}

func (this *CDbHandler) ConnectByRule(rule string, dbtype string) (err error) {
	this.m_db, err = sql.Open(dbtype, rule)
	if err != nil {
		return err
	}
	this.m_db.SetMaxOpenConns(2000)
	this.m_db.SetMaxIdleConns(1000)
	this.m_db.Ping()
	return nil
}

func (this *CDbHandler) ConnectByCfg(path string) error {
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	var host string = "localhost"
	var port uint = 3306
	var username string = "root"
	var userpwd string = "123456"
	var dbname string = "test"
	var dbtype string = "mysql"
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		content := string(a)
		r, _ := regexp.Compile("(.*)?=(.*)?")
		ret := r.FindStringSubmatch(content)
		if len(ret) != 3 {
			continue
		}
		k := ret[1]
		v := ret[2]
		switch k {
		case "host":
			host = v
		case "port":
			port_tmp, _ := strconv.ParseUint(v, 10, 32)
			port = uint(port_tmp)
		case "username":
			username = v
		case "userpwd":
			userpwd = v
		case "dbname":
			dbname = v
		case "dbtype":
			dbtype = v
		}
	}
	return this.Connect(host, port, username, userpwd, dbname, dbtype)
}

func (this *CDbHandler) Disconnect() {
	this.m_db.Close()
}

func (this *CDbHandler) Create() (error) {
	var err error = nil
	var _ error = err
	return nil
}

func (this *CDbHandler) GetAllStaffInfo(output0 *[]CGetAllStaffInfoOutput) (error, uint64) {
	var rowCount uint64 = 0
	if this.m_db == nil {
		return errors.New("db is nil"), 0
	}
	tx, err := this.m_db.Begin()
	if err != nil {
		return err, 0
	}
	var result sql.Result
	var _ = result
	var _ error = err
	rows0, err := tx.Query(fmt.Sprintf(`select tmp.infrastructureUuid
, tmp.staffName, tmp.address, tmp.nation, tmp.nationality, tmp.nativePlace, tmp.gender
, sc.contactType, sc.contactInfo
from
(
select si.staffUuid as suid, sir.infrastructureUuid
, si.staffName, si.address, si.nation, si.nationality, si.nativePlace, si.gender
from t_vss_staff_infrastructure_rl as sir
inner join t_vss_staff_info as si
on sir.staffUuid = si.staffUuid
) as tmp
left join t_vss_staff_contact as sc
on tmp.suid = sc.staffUuid
where sc.contactType = 'cellphone';`))
	if err != nil {
		tx.Rollback()
		return err, rowCount
	}
	defer rows0.Close()
	for rows0.Next() {
		rowCount += 1
		tmp := CGetAllStaffInfoOutput{}
		var infrastructureUuid sql.NullString
		var staffName sql.NullString
		var address sql.NullString
		var nation sql.NullString
		var nationality sql.NullString
		var nativePlace sql.NullString
		var gender sql.NullString
		var contactType sql.NullString
		var contactInfo sql.NullString
		scanErr := rows0.Scan(&infrastructureUuid, &staffName, &address, &nation, &nationality, &nativePlace, &gender, &contactType, &contactInfo)
		if scanErr != nil {
			continue
		}
		tmp.InfrastructureUuid = infrastructureUuid.String
		tmp.InfrastructureUuidIsValid = infrastructureUuid.Valid
		tmp.StaffName = staffName.String
		tmp.StaffNameIsValid = staffName.Valid
		tmp.Address = address.String
		tmp.AddressIsValid = address.Valid
		tmp.Nation = nation.String
		tmp.NationIsValid = nation.Valid
		tmp.Nationality = nationality.String
		tmp.NationalityIsValid = nationality.Valid
		tmp.NativePlace = nativePlace.String
		tmp.NativePlaceIsValid = nativePlace.Valid
		tmp.Gender = gender.String
		tmp.GenderIsValid = gender.Valid
		tmp.ContactType = contactType.String
		tmp.ContactTypeIsValid = contactType.Valid
		tmp.ContactInfo = contactInfo.String
		tmp.ContactInfoIsValid = contactInfo.Valid
		*output0 = append(*output0, tmp)
	}
	tx.Commit()
	return nil, rowCount
}

func (this *CDbHandler) GetInfrastructureInfoByUuid(input0 *CGetInfrastructureInfoByUuidInput, output0 *CGetInfrastructureInfoByUuidOutput) (error, uint64) {
	var rowCount uint64 = 0
	if this.m_db == nil {
		return errors.New("db is nil"), 0
	}
	tx, err := this.m_db.Begin()
	if err != nil {
		return err, 0
	}
	var result sql.Result
	var _ = result
	var _ error = err
	rows0, err := tx.Query(fmt.Sprintf(`select infrastructureName, parentUuid from t_vss_infrastructure_info where infrastructureUuid = ?;`), input0.InfrastructureUuid)
	if err != nil {
		tx.Rollback()
		return err, rowCount
	}
	defer rows0.Close()
	for rows0.Next() {
		rowCount += 1
		var infrastructureName sql.NullString
		var parentUuid sql.NullString
		scanErr := rows0.Scan(&infrastructureName, &parentUuid)
		if scanErr != nil {
			continue
		}
		output0.InfrastructureName = infrastructureName.String
		output0.InfrastructureNameIsValid = infrastructureName.Valid
		output0.ParentUuid = parentUuid.String
		output0.ParentUuidIsValid = parentUuid.Valid
	}
	tx.Commit()
	return nil, rowCount
}

