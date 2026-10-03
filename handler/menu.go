// handler/menu.go
// 职责：处理 GET /geeker/menu/list —— 返回当前用户的菜单树
// 注意：这个文件不查数据库，直接返回硬编码的菜单结构（对应前端 mock/menu/list.js 的内容）
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMenuList：菜单列表处理函数
// 参数：c *gin.Context（此接口是 GET，没有请求体，所以不需要解析参数）
// 逻辑：组装菜单树（含嵌套 children）→ 返回
// 菜单结构说明：
//   - 叶子菜单：有 component（前端用它映射到实际页面组件）
//   - 分组菜单：无 component、有 redirect + children（前端渲染成可展开的分组）
//   - meta.isKeepAlive：前端页面缓存的开关（true 的页面切走再切回保留状态）
func GetMenuList(c *gin.Context) {
	// menus：菜单树。用 []gin.H 表示 JSON 数组，gin.H 表示 JSON 对象
	// 嵌套 children 也是 []gin.H —— Go 里的 JSON 就是这么一层层套出来的
	menus := []gin.H{
		// 首页：叶子菜单
		{"path": "/home/index", "name": "home", "component": "home/index",
			"meta": gin.H{"title": "首页", "icon": "HomeFilled"}},

		// ProTable 分组：无 component，有 redirect 和 children
		{
			"path": "/proTable", "name": "proTable", "redirect": "/proTable/virtual",
			"meta": gin.H{"title": "ProTable演示", "icon": "Grid"},
			"children": []gin.H{
				// 子菜单 1：虚拟滚动（isKeepAlive: true = 页面缓存开关）
				{"path": "/proTable/virtual", "name": "virtualTable", "component": "proTable/virtualTable",
					"meta": gin.H{"title": "虚拟滚动演示", "icon": "DataLine", "isKeepAlive": true}},
				// 子菜单 2：表单联动
				{"path": "/proTable/formLinkage", "name": "formLinkage", "component": "proTable/formLinkage",
					"meta": gin.H{"title": "表单联动演示", "icon": "Connection"}},
			},
		},

		{"path": "/directives", "name": "directives", "component": "directives/index",
			"meta": gin.H{"title": "指令演示", "icon": "MagicStick"}},
		{"path": "/system/user", "name": "systemUser", "component": "system/user/index",
			"meta": gin.H{"title": "用户管理", "icon": "User"}},
		{"path": "/upload", "name": "upload", "component": "upload/index",
			"meta": gin.H{"title": "分片上传演示", "icon": "Upload"}},
		{"path": "/dashboard", "name": "dashboard", "component": "dashboard/index",
			"meta": gin.H{"title": "数据可视化", "icon": "TrendCharts"}},
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": menus, "msg": "success"})
}
