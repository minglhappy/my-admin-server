// main.go
// 职责：后端入口——注册路由、挂中间件、启动 HTTP 服务
// 类比前端：main.ts 的 app.use(router) + app.mount()
package main

import (
	"my-admin-server/handler" // 自己项目的包（模块名/目录名）
	"my-admin-server/middleware"

	"github.com/gin-gonic/gin"
)

// main：Go 程序的唯一入口（和 C 语言一样，main 函数是起点）
func main() {
	// ── 第 1 步：创建 gin 引擎 ──
	// gin.Default()：带默认中间件（日志 + 异常恢复）的引擎
	// r 就是"路由表 + 中间件链"的容器——所有请求都经过它
	r := gin.Default()

	// ── 第 2 步：挂全局中间件 ──
	// r.Use：注册中间件。执行顺序 = 注册顺序：
	//   请求进来 → Cors（设跨域头）→ [匹配路由] → 路由组中间件 JwtAuth → 业务 handler
	r.Use(middleware.Cors())

	// ── 第 3 步：注册路由 ──
	// 登录接口：直接挂在引擎上，不需要鉴权（用户还没有 token）
	r.POST("/geeker/login", handler.Login)

	// 路由分组：需要鉴权的接口统一放进 auth 组
	// 分组的价值：JwtAuth 写一次，组内所有接口自动受保护
	auth := r.Group("/geeker")
	auth.Use(middleware.JwtAuth()) // ★ 组级中间件：组内每个请求都先过 JWT 校验
	{
		auth.GET("/menu/list", handler.GetMenuList)       // 菜单
		auth.GET("/auth/buttons", handler.GetAuthButtons) // 按钮权限
		auth.POST("/user/list", handler.GetUserList)      // 用户列表
		auth.GET("/dashboard/data", handler.GetDashboardData)
		auth.POST("/upload/chunk", handler.UploadChunk)
		auth.GET("/upload/check", handler.CheckChunks)
		auth.POST("/upload/merge", handler.MergeChunks)
	}

	// ── 第 4 步：启动服务 ──
	// Run(":3000")：监听 3000 端口（":3000" 表示监听所有网卡的 3000 端口）
	// 这个调用会阻塞——服务一直运行，直到 Ctrl+C
	r.Run(":3000")
}
