// middleware/jwt.go
// 职责：JWT 校验中间件 —— 拦截未登录/伪造 token 的请求
// 对应前端：路由守卫里 userStore.token 的检查——只是这次检查发生在服务端
// 铁律：前端的鉴权只是"体验优化"（看不见页面），服务端的鉴权才是"真正的安全"（拿不到数据）
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret：与 handler/login.go 里签发用的密钥必须一致！
// 签发用这把钥匙盖章，校验也用这把钥匙验章——不一致则所有 token 都失效
var jwtSecret = []byte("my-admin-secret")

// JwtAuth：JWT 校验中间件
// 逻辑：① 取请求头里的 token → ② 验证签名和过期时间 → ③ 通过则放行，失败则拦截
// 关键变量：tokenStr —— 从 x-access-token 请求头读取
//
//	（这个请求头名是前端阶段 5 在 RequestHttp 拦截器里注入的，前后端约定）
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ── 第 1 步：取 token ──
		tokenStr := c.GetHeader("x-access-token")
		if tokenStr == "" {
			// 没有 token → 拦截
			// 注意：返回 code: 401 而不是 HTTP 401——
			// 对应前端响应拦截器 data.code === ResultEnum.OVERDUE 的判断（跳登录页）
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "未登录"})
			c.Abort() // ★ Abort：拦截请求，后面的 handler 不会执行
			return
		}

		// ── 第 2 步：验证签名 ──
		// jwt.Parse：解析并验证 token
		// 第二个参数是"密钥提供函数"：告诉库用哪把钥匙验签
		// 验证内容：签名是否被篡改 + 是否过期（exp 字段）
		_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			// 验签失败：token 是伪造的 / 已过期 / 格式错误
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "token 无效或已过期"})
			c.Abort()
			return
		}

		// ── 第 3 步：放行 ──
		c.Next()
	}
}
