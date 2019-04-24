package controllers

import (
    // "github.com/astaxie/beego"
    //"github.com/astaxie/beego/validation"
    // "github.com/lisijie/goblog/models"
    // "github.com/lisijie/goblog/util"
    // "strconv"
    "encoding/json"
    "myearAdmin/models"
    "strings"
    "time"
)

type ManagersController struct {
    baseController
}

type ManagerController struct {
    baseController
}

type GetManagerListRsp struct {
    Managers []models.Manager `json:"managers"`
}

type GetManagerRsp struct {
    Manager models.Manager `json:"manager"`
}

//CreateManagerReq 创建manager请求
type CreateManagerReq struct {
    Name  string `json:"name" valid:"Required"`
    Pass  string `json:"pass" valid:"Required"`
    Role  int    `json:"role" valid:"Required"`
    Phone string `json:"phone"`
    Email string `json:"email"`
}

//UpdateManagerReq 更新manager请求
type UpdateManagerReq struct {
    Pass  string `json:"pass" valid:"Required`
    Role  int    `json:"role" valid:"Required`
    Phone string `json:"phone" valid:"Required`
    Email string `json:"email" valid:"Required`
}

func (c *ManagersController) GetManagerList() {
    mList, err := models.GetAllMangerList()
    if err != nil {
        c.RspErr(1, "管理员列表不存在")
        return
    }

    c.RspData(GetManagerListRsp{Managers: mList})
}

func (c *ManagersController) CreateManager() {
    createReq := CreateManagerReq{}
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &createReq); err != nil {
        c.RspErr(ReqParaErr, err.Error())
        return
    }

    if c.ValidRequest(&createReq) == false {
        return
    }

    if err := models.AddManger(models.Manager{Name: createReq.Name,
        Pass:       createReq.Pass,
        Role:       createReq.Role,
        Email:      createReq.Email,
        Phone:      createReq.Phone,
        CreateTime: time.Now()}); err != nil {
        c.RspErr(2, err.Error())
        return
    }

    c.RspData(struct{}{})
}

func (c *ManagersController) DeleteManager() {
    ids := c.GetString("ids")
    if ids == "" {
        c.RspErr(ReqParaErr, "缺少参数")
        return
    }

    num, err := models.DeleteManger(strings.Split(ids, ","))
    if err != nil {
        c.RspErr(1, err.Error())
        return
    }
    c.RspData(map[string]int64{"num": num})
}

func (c *ManagerController) GetManagerByID() {
    id := c.Ctx.Input.Param(":id")
    m, err := models.GetMangerByID(id)
    if err != nil {
        c.RspErr(1, err.Error())
        return
    }

    c.RspData(GetManagerRsp{Manager: m})
}

func (c *ManagerController) UpdateManager() {
    updateReq := UpdateManagerReq{}
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &updateReq); err != nil {
        c.RspErr(ReqParaErr, err.Error())
        return
    }

    if c.ValidRequest(&updateReq) == false {
        return
    }

    id := c.Ctx.Input.Param(":id")

    err := models.UpdateManger(id, updateReq.Pass, updateReq.Role, updateReq.Phone, updateReq.Email)
    if err != nil {
        c.RspErr(1, err.Error())
        return
    }
    c.RspData(struct{}{})
}
