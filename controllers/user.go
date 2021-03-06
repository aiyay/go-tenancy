package controllers

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/snowlyg/go-tenancy/common"
	"github.com/snowlyg/go-tenancy/models"
	"github.com/snowlyg/go-tenancy/sysinit"
	"github.com/snowlyg/go-tenancy/transformer"
	"github.com/snowlyg/go-tenancy/validatas"
	"github.com/snowlyg/gotransformer"
)

type UserController struct {
	Ctx iris.Context
}

// GetUsers handles GET: http://localhost:8080/user/table.
func (c *UserController) GetTable() interface{} {

	var pagination common.Pagination
	if err := c.Ctx.ReadQuery(&pagination); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("分页参数获取错误：%v", err)}
	}

	args := map[string]interface{}{}
	count, users := sysinit.UserService.GetAll(args, &pagination, false)

	return common.Table{Code: 0, Msg: "", Count: count, Data: c.transformerTableUsers(users)}
}

// Get handles GET: http://localhost:8080/user.
func (c *UserController) Get() mvc.Result {
	return mvc.View{
		Name: "user/index.html",
	}
}

// Get handles GET: http://localhost:8080/user/create.
func (c *UserController) GetCreate() mvc.Result {
	return mvc.View{
		Name: "user/add.html",
	}
}

// GetRoletBy handles GET: http://localhost:8080/user/role/:id.
func (c *UserController) GetRoleBy(id uint) interface{} {

	user, _ := sysinit.UserService.GetByID(id)

	args := map[string]interface{}{}
	_, roles := sysinit.RoleService.GetAll(args, nil, false)

	return common.ActionResponse{Status: true, Msg: "", Data: c.transformerSelectRoles(roles, user.ID)}
}

// Get handles GET: http://localhost:8080/user/id.
func (c *UserController) GetBy(id uint) mvc.Result {
	user, _ := sysinit.UserService.GetByID(id)

	return mvc.View{
		Name: "user/edit.html",
		Data: iris.Map{
			"User": user,
		},
	}
}

// Get handles Post: http://localhost:8080/user.
// 使用 ReadJSON 获取数据前端数据需要格式化成json, JSON.stringify(data.field),
func (c *UserController) Post() interface{} {

	var user models.User

	if err := c.Ctx.ReadJSON(&user); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("数据获取错误：%v", err)}
	}

	if user.Password == "" {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("密码不能为空")}
	}

	if err := validatas.Vaild(user); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("数据验证错误：%v", err)}
	}

	if err := sysinit.UserService.Create(user.Password, &user); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("用户创建错误：%v", err)}
	}

	return common.ActionResponse{Status: true, Msg: "操作成功"}
}

// Get handles Post: http://localhost:8080/user/id.
func (c *UserController) PostBy(id uint) interface{} {

	var user transformer.UserUpdate

	if err := c.Ctx.ReadJSON(&user); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("数据获取错误：%v", err)}
	}

	if err := validatas.Vaild(user); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("数据验证错误：%v", err)}
	}

	if err := sysinit.UserService.UpdateUser(id, &user); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("用户更新错误：%v", err)}
	}

	return common.ActionResponse{Status: true, Msg: "操作成功"}
}

// Get handles Post: http://localhost:8080/user/id.
func (c *UserController) DeleteBy(id uint) interface{} {
	if err := sysinit.UserService.DeleteByID(id); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("用户删除错误：%v", err)}
	}

	return common.ActionResponse{Status: true, Msg: "操作成功"}
}

// Get handles Post: http://localhost:8080/user/deletes.
func (c *UserController) PostDeletes() interface{} {
	var userIds []common.Id

	if err := c.Ctx.ReadJSON(&userIds); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("数据获取错误：%v", err)}
	}

	if err := sysinit.UserService.DeleteMnutil(userIds); err != nil {
		return common.ActionResponse{Status: false, Msg: fmt.Sprintf("用户删除错误：%v", err)}
	}

	return common.ActionResponse{Status: true, Msg: "操作成功"}
}

// transformerTableUsers 菜单表格接口数据转换
func (c *UserController) transformerTableUsers(users []*models.User) []*transformer.UserTable {
	var tableusers []*transformer.UserTable
	for _, user := range users {
		tableuser := &transformer.UserTable{}
		g := gotransformer.NewTransform(tableuser, user, "")
		_ = g.Transformer()

		roles, err := sysinit.UserService.GetRolesByID(user.ID)
		if err == nil {
			for _, role := range roles {
				tableuser.RoleNames += role.DisplayName + " ; "
			}
		}

		tableusers = append(tableusers, tableuser)
	}

	return tableusers
}

// transformerTableUsers 菜单表格接口数据转换
func (c *UserController) transformerSelectRoles(roles []*models.Role, userId uint) []*transformer.RoleSelect {
	var tableroles []*transformer.RoleSelect
	userRoles, err := sysinit.UserService.GetRolesByID(userId)

	if err == nil {
		for _, role := range roles {
			tableuser := &transformer.RoleSelect{}
			tableuser.Name = role.DisplayName
			tableuser.Value = role.ID

			for _, userRole := range userRoles {
				if userRole.ID == role.ID {
					tableuser.Selected = true
				}
			}

			tableroles = append(tableroles, tableuser)
		}
	}

	return tableroles
}
