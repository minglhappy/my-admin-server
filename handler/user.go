// handler/user.go
// 职责：处理 POST /geeker/user/list —— 用户列表（分页 + 条件过滤）
// 对应前端：ProTable 的 requestApi —— 前端传 {pageNum, pageSize, username, status}
//
//	后端返回 { list, total }
//
package handler

import (
	"net/http"
	"strings" // 字符串工具（Contains 判断包含）

	"github.com/gin-gonic/gin"
)

// UserItem：用户数据结构
// json tag 与前端 UserItem 接口逐字段对齐——tag 写错，前端收到的字段就是空的
type UserItem struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Status     int    `json:"status"`     // 1=启用 0=禁用
	CreateTime string `json:"createTime"` // 注意大写 T！前端字段就是 createTime
}

// users：内存中的用户数据（教学版"数据库"）
// 包级变量：整个程序共享一份，重启后恢复初始值
// 真实项目：这里换成数据库查询（后续阶段会做）
var users = []UserItem{
	{1, "admin", "管理员", 1, "2026-08-01 10:00:00"},
	{2, "zhangsan", "张三", 1, "2026-08-02 11:30:00"},
	{3, "lisi", "李四", 0, "2026-08-03 14:20:00"},
	{4, "wangwu", "王五", 1, "2026-08-04 09:15:00"},
	{5, "zhaoliu", "赵六", 0, "2026-08-05 16:40:00"},
}

// UserListReq：列表查询参数
// 关键知识点：Status 用指针 *int 而不是 int！
// 原因：int 的零值是 0——前端"没传 status"和"传了 status=0（禁用）"都会解析成 0，无法区分
// 指针的零值是 nil——没传 = nil，传了 = 指向具体值的指针，完美区分
type UserListReq struct {
	PageNum  int    `json:"pageNum"`  // 页码（从 1 开始）
	PageSize int    `json:"pageSize"` // 每页条数
	Username string `json:"username"` // 用户名模糊搜索
	Status   *int   `json:"status"`   // 状态过滤（可选字段 → 指针）
}

// GetUserList：用户列表处理函数
// 逻辑四步：① 解析参数 + 默认值 → ② 过滤 → ③ 分页切片 → ④ 返回
func GetUserList(c *gin.Context) {
	// ── 第 1 步：解析参数 + 兜底默认值 ──
	var req UserListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	// 前端没传分页参数时给默认值（防御性编程：后端不能假设前端一定守规矩）
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// ── 第 2 步：过滤 ──
	// filtered：过滤后的结果（先收集到新切片，不修改原始数据）
	var filtered []UserItem
	for _, u := range users {
		// 用户名模糊匹配：前端传了 username 且不包含 → 跳过这条
		if req.Username != "" && !strings.Contains(u.Username, req.Username) {
			continue // continue：跳过本次循环，处理下一条
		}
		// 状态精确匹配：Status 不是 nil（前端传了）且值不相等 → 跳过
		// *req.Status：解引用指针，取出它指向的值
		if req.Status != nil && u.Status != *req.Status {
			continue
		}
		filtered = append(filtered, u) // append：向切片追加元素
	}

	// ── 第 3 步：分页切片 ──
	total := len(filtered)                    // 总条数（过滤后）
	start := (req.PageNum - 1) * req.PageSize // 起始下标：第 2 页每页 10 条 → 从第 10 条开始
	if start > total {
		start = total // 页码超出范围 → 返回空列表（防切片越界）
	}
	end := start + req.PageSize // 结束下标
	if end > total {
		end = total
	}

	// ── 第 4 步：返回 ──
	// filtered[start:end]：Go 切片语法，取 [start, end) 区间——类似 JS 的 slice
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"list": filtered[start:end], "total": total},
		"msg":  "success",
	})
}

