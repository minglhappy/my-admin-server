// middleware/cors.go
// 职责：跨域中间件 —— 允许前端的浏览器跨域调用本后端
// 背景：前端页面跑在 localhost:8848，后端在 localhost:3000
//
//	浏览器认为这是"跨域"（域名相同但端口不同也算跨域），默认拦截响应
//
// 原理：跨域是浏览器管的，但"钥匙"在服务端——
//
//	服务端在响应头里声明"我允许 8848 来访问"，浏览器才放行
package middleware

import (
	"github.com/gin-gonic/gin"
)

// Cors：跨域中间件（返回一个 gin.HandlerFunc）
// 中间件本质：一个"包裹"正常处理函数的函数——
//
//	请求进来 → 中间件先执行（设置响应头）→ c.Next() 放行到业务 handler
//
// 关键变量：三个 Access-Control-* 响应头
//
//	Allow-Origin：允许哪些来源访问（* = 全部，真实项目写具体域名）
//	Allow-Methods：允许哪些 HTTP 方法
//	Allow-Headers：允许浏览器带哪些自定义请求头（x-access-token 必须在这里！）
// func Cors() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, x-access-token")

// 		// 预检请求（OPTIONS）：浏览器在"真正的跨域请求"发出前，先发一个 OPTIONS 探路
// 		// 问服务端："我能不能发这种请求？"——这里直接放行（204 = 无内容）
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204) // Abort：终止后续处理（不走业务 handler）
// 			return
// 		}

// 		c.Next() // 放行：继续执行后面的中间件和业务 handler
// 	}
// }


func Cors() gin.HandlerFunc {
    // 白名单：允许哪些"来访者"跨域（Origin 是前端页面的地址，不是后端自己的地址）
    allowOrigins := map[string]bool{
        "http://localhost:8848": true, // 前端 dev server（Vite）
    }

    return func(c *gin.Context) {
        // 拿到请求是从哪个网站发出来的
        origin := c.Request.Header.Get("Origin")

        // 查验白名单，只有在白名单里的才允许跨域
        if allowOrigins[origin] {
            c.Header("Access-Control-Allow-Origin", origin) // 动态设置为那个具体的域名
        }

        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, x-access-token, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true") // 允许前端带 Cookie

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}