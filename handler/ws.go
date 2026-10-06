// handler/ws.go
// 职责：WebSocket 聊天室——连接管理 + 消息广播（coder/websocket 版）
// 核心概念与 gorilla 版相同，但 API 全部围绕 context 设计
package handler

import (
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

// Client：一个在线用户
// 为什么要包一层 struct 而不是直接用 conn？
//
//	因为 coder/websocket 规定：同一连接同一时刻只允许一个 goroutine 写。
//	广播时多个 goroutine 可能同时给同一个连接写消息 → 必须给每个连接配一把"写锁"
type Client struct {
	conn    *websocket.Conn // WebSocket 连接
	writeMu sync.Mutex      // 写锁：保证同一时刻只有一个 goroutine 写这个连接
}

// 在线连接表 + 表锁
// 表锁保护的是 map 本身（增删条目）；写锁保护的是每个连接（并发写）——两层锁各管各的
var (
	clientsMu sync.Mutex
	clients   = make(map[*Client]bool)
)

// ChatWS：聊天室 WebSocket 端点
// 逻辑四步：① 升级连接 → ② 注册进在线表 → ③ 读循环（收到就广播）→ ④ 注销
func ChatWS(c *gin.Context) {
	// ── ① 升级：把 HTTP 请求升级为 WebSocket 连接 ──
	// Accept 对应 gorilla 的 upgrader.Upgrade
	// AcceptOptions.OriginPatterns：允许哪些来源建连（"*" = 全放行，教学版）
	//   真实项目写具体域名如 ["localhost:8848"]，防恶意网站盗用你的 WebSocket
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		return // 升级失败（不是 WebSocket 请求）
	}
	// Close：优雅关闭。StatusCode 是 WebSocket 协议的关闭码（1000 = 正常关闭）
	defer conn.Close(websocket.StatusNormalClosure, "bye")

	client := &Client{conn: conn}

	// ── ② 注册：加入在线表 ──
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()
	broadcast([]byte("有人加入了聊天室"))

	// ── ③ 读循环：这条连接的一生 ──
	// ctx：取请求的 context——客户端断开/服务端关闭时 ctx 自动取消，Read 立即返回错误
	//     这就是 context 的价值：不需要手动管理"停止读"的信号
	ctx := c.Request.Context()
	for {
		// Read：阻塞等待消息。返回 (消息类型, 消息内容, 错误)
		// 消息类型：MessageText(文本) / MessageBinary(二进制)——聊天用文本
		// IoT 场景常用二进制传传感器数据（更紧凑）
		_, msg, err := conn.Read(ctx)
		if err != nil {
			break // 客户端断开或连接出错 → 退出循环
		}
		broadcast(msg) // 收到一条 → 广播给所有人
	}

	// ── ④ 注销：清理在线表 ──
	clientsMu.Lock()
	delete(clients, client)
	clientsMu.Unlock()
	broadcast([]byte("有人离开了聊天室"))
}

// broadcast：广播——把消息发给所有在线连接
// 每个连接的写入流程：拿写锁 → 带超时写入 → 释放锁
// 超时的意义：某个客户端网络卡死（Write 阻塞），不能拖垮整个广播
func broadcast(msg []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		client.writeMu.Lock()

		// WithTimeout：给写入操作加 5 秒超时——超时则放弃这个连接
		// context 在这里的作用：让"可能卡死的网络写入"变得可控
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		client.conn.Write(ctx, websocket.MessageText, msg)
		cancel()

		client.writeMu.Unlock()
	}
}
