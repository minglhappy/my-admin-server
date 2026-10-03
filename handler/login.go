// handler/login.go
// 职责：处理 POST /geeker/login —— 校验账号密码，签发 JWT token
package handler

import (
	"net/http" // HTTP 状态码常量（http.StatusOK = 200）
	"time"     // 时间处理（token 过期时间）

	"github.com/gin-gonic/gin"     // gin 框架：处理请求/响应
	"github.com/golang-jwt/jwt/v5" // JWT 库：签发与校验 token
)

// jwtSecret：JWT 的签名密钥
// 作用：token 的"防伪印章"——签发和校验必须用同一个密钥，别人拿不到密钥就无法伪造 token
// 注意：真实项目里这个密钥要从配置/环境变量读，绝不能硬编码后提交到公开仓库
var jwtSecret = []byte("my-admin-secret")

// LoginReq：登录请求的结构体
// 功能：定义前端会传什么参数（gin 用它自动解析 JSON 请求体）
// json tag 的含义："username" 是前端 JSON 里的字段名，
// Go 的字段名是 Username——tag 就是两者之间的翻译器
type LoginReq struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// Login：登录处理函数
// 参数：c *gin.Context —— 一次 HTTP 请求的完整上下文（请求数据、响应写入、中间件控制都在里面）
// 逻辑：① 解析请求体 → ② 校验（教学版任意账号通过）→ ③ 签发 JWT → ④ 返回统一结构
func Login(c *gin.Context) {
	// ── 第 1 步：解析请求体 ──
	// ShouldBindJSON：把请求体的 JSON 自动映射进 LoginReq 结构体
	// 映射规则：JSON 字段名 ↔ json tag。前端发 {"username":"admin"} → req.Username = "admin"
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 解析失败（比如前端发了非 JSON）→ 返回错误
		// 注意：HTTP 状态码用 200，业务状态用 code: 400
		// —— 这是和前端约定的统一响应结构 { code, data, msg }
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// ── 第 2 步：校验账号密码 ──
	// 教学版：任意账号密码都放行。真实项目在这里查数据库校验
	// （留给你以后扩展：SELECT * FROM users WHERE username = ? AND password = ?）

	// ── 第 3 步：签发 JWT ──
	// JWT 由三部分组成（Base64 编码的）：头部.载荷.签名
	// 这里只设置"载荷"（payload），头部和签名由库自动生成
	// 载荷里放什么：username（标识是谁）+ exp（过期时间）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,                          // 用户标识：以后可以从 token 里解析出"这是谁的请求
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 过期时间：24 小时后 token 失效（Unix 秒）
	})
	// SignedString：用密钥对 token 签名，生成最终的字符串
	tokenStr, _ := token.SignedString(jwtSecret)

	// ── 第 4 步：返回统一结构 ──
	// 前端约定：{ code: 200, data: {...}, msg: "success" }
	// data.access_token 对应前端 loginApi 的 LoginResult 接口
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"access_token": tokenStr},
		"msg":  "success",
	})
}
