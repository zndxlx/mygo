package routers

import (
    "github.com/astaxie/beego"
    "myearAdmin/controllers"
)

func init() {
    // 登陆 .
    beego.Router("/admin/login", &controllers.AuthController{}, "post:Login")
    beego.Router("/admin/logout", &controllers.AuthController{}, "post:Logout")

    // 管理员管理
    beego.Router("/admin/managers", &controllers.ManagersController{}, "get:GetManagerList")
    beego.Router("/admin/managers", &controllers.ManagersController{}, "post:CreateManager")
    beego.Router("/admin/managers", &controllers.ManagersController{}, "delete:DeleteManager")
    beego.Router("/admin/managers/:id", &controllers.ManagerController{}, "get:GetManagerByID")
    beego.Router("/admin/managers/:id", &controllers.ManagerController{}, "post:UpdateManager")

    // 用户管理.

    // beego.InsertFilter("/*", beego.FinishRouter, controllers.FinishRouter, false)
    // beego.InsertFilter("/*", beego.AfterExec, controllers.AfterExec)
    // beego.InsertFilter("/*", beego.BeforeRouter, controllers.BeforeRouter)
}
