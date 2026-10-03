// handler/buttons.go
// 职责：处理 GET /geeker/auth/buttons —— 返回当前用户的按钮权限码列表
// 权限码约定：模块:操作（user:add、user:delete）——前端 v-auth 指令用它判断按钮显隐
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAuthButtons：按钮权限码列表
// 参数：c *gin.Context（GET 无请求体）
// 逻辑：教学版返回固定列表。真实项目：根据 token 解析出用户 → 查角色 → 查权限表
// 返回的权限码决定前端哪些按钮可见：
//
//	user:add ✅ → "新增用户"按钮显示
//	user:edit ❌（不在列表）→ "编辑"按钮被 v-auth 移除
func GetAuthButtons(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": []string{"user:add", "user:delete", "user:export"}, // 权限码数组
		"msg":  "success",
	})
}
