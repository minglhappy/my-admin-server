// handler/upload.go
// 职责：分片上传三接口（chunk 接收/check 查询/merge 合并）
// 前端约定：index/filename 走 query 参数（不是 multipart body）
package handler

import (
	"net/http"
	"strconv" // 字符串转数字（query 参数都是字符串）
	"sync"    // 并发锁

	"github.com/gin-gonic/gin"
)

// uploadedChunks：内存版"服务端磁盘"——文件名 → 已收到的片序号集合
// 为什么用 map[int]bool 而不是 []int：查重 O(1)（分片序号判断"有没有"是高频操作）
var (
	uploadedChunks   = make(map[string]map[int]bool)
	uploadedChunksMu sync.Mutex // ★ 互斥锁：Go 的 map 不是并发安全的！
// 前端 3 个 worker 并发上传 → 3 个 goroutine 同时写 map → 不加锁会直接 panic
// Mutex 保证同一时刻只有一个 goroutine 能写
)

// UploadChunk：接收一个分片
// 参数：query 里的 filename（文件名）和 index（片序号）
// 逻辑：把片序号记录到 uploadedChunks（真实项目：片内容存磁盘/对象存储）
func UploadChunk(c *gin.Context) {
	filename := c.Query("filename")
	index, _ := strconv.Atoi(c.Query("index")) // "3" → 3

	uploadedChunksMu.Lock() // 上锁（写操作前）
	if uploadedChunks[filename] == nil {
		uploadedChunks[filename] = make(map[int]bool)
	}
	uploadedChunks[filename][index] = true
	uploadedChunksMu.Unlock() // 解锁（写操作后）

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "chunk ok"})
}

// CheckChunks：查询已传分片（断点续传）
// 逻辑：返回某个文件已收到的所有片序号数组
func CheckChunks(c *gin.Context) {
	filename := c.Query("filename")

	uploadedChunksMu.Lock()
	defer uploadedChunksMu.Unlock() // defer：函数返回前自动解锁（防止中途 return 忘解锁）

	chunks := []int{}
	for i := range uploadedChunks[filename] {
		chunks = append(chunks, i)
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": chunks, "msg": "ok"})
}

// MergeChunks：合并分片（教学版：清理记录，模拟"合并完成"）
func MergeChunks(c *gin.Context) {
	var req struct {
		Filename string `json:"filename"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	uploadedChunksMu.Lock()
	delete(uploadedChunks, req.Filename) // 合并完成 → 清理（和 mock 版行为一致）
	uploadedChunksMu.Unlock()

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"url": "/files/" + req.Filename}, "msg": "merge ok"})
}
