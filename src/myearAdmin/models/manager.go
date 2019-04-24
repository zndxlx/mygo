package models

import (
    //"crypto/rand"
    "fmt"
    "strings"
    "time"

    "database/sql"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
    //"golang.org/x/crypto/scrypt"
)

//Manager 管理员信息
type Manager struct {
    Id         int64     `json:"id"`
    Name       string    `json:"name"`
    Pass       string    `json:"-"`
    Role       int       `json:"role"`
    Email      string    `json:"email"`
    Phone      string    `json:"phone"`
    CreateTime time.Time `json:"createTime"`
}

//GetAllMangerList 返回所有的manager
func GetAllMangerList() (managers []Manager, err error) {
    o := orm.NewOrm()
    _, err = o.Raw("select id, name, pass, phone, email, role, create_time from manager").QueryRows(&managers)
    return
}

//GetMangerByName 根据用户名查找manager
func GetMangerByName(name string) (manager Manager, err error) {
    o := orm.NewOrm()
    err = o.Raw("select id, name, pass, phone, email, role, create_time from manager where name = ?", name).QueryRow(&manager)
    return
}

//GetMangerByID 根据用户名查找manager
func GetMangerByID(id string) (manager Manager, err error) {
    o := orm.NewOrm()
    err = o.Raw("select id, name, pass, phone, email, role, create_time from manager where id = ?", id).QueryRow(&manager)
    return
}

func AddManger(u Manager) (err error) {
    o := orm.NewOrm()
    _, err = o.Raw("insert into manager (name, pass, phone, email, role, create_time) values (?,?,?,?,?,?)",
        u.Name, u.Pass, u.Phone, u.Email, u.Role, u.CreateTime).Exec()
    return
}

func DeleteManger(ids []string) (num int64, err error) {
    o := orm.NewOrm()
    sqlQTmp := strings.Repeat("?,", len(ids))
    sqlQ := sqlQTmp[0 : len(sqlQTmp)-1]
    sqlFormat := fmt.Sprintf("delete from manager where id in (%s)", sqlQ)

    var res sql.Result
    res, err = o.Raw(sqlFormat, ids).Exec()
    if err == nil {
        num, _ = res.RowsAffected()
    }
    return
}

func UpdateManger(id string, pass string, role int, phone string, email string) (err error) {
    o := orm.NewOrm()

    _, err = o.Raw("update manager set pass = ?,role = ?,phone = ?, email = ? where id = ?",
        pass, role, phone, email, id).Exec()

    return
}
