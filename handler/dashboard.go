// handler/dashboard.go
// 职责：处理 GET /geeker/dashboard/data —— 仪表盘数据（每次随机，模拟实时变化）
package handler

import (
	"math/rand"   // 随机数
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetDashboardData：仪表盘数据
// 逻辑：随机生成折线/柱状/饼图数据 + 当前时间（对应前端轮询的"数据变化"效果）
func GetDashboardData(c *gin.Context) {
	// random：闭包辅助函数——返回 [base, base+rng) 的随机数
	random := func(base, rng int) int { return base + rand.Intn(rng) }

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"line": gin.H{
				"days":   []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"},
				"values": []int{random(800, 600), random(800, 600), random(800, 600), random(800, 600), random(800, 600), random(800, 600), random(800, 600)},
			},
			"bar": gin.H{
				"categories": []string{"苹果", "香蕉", "橙子", "葡萄", "西瓜"},
				"values":     []int{random(60, 150), random(60, 150), random(60, 150), random(60, 150), random(60, 150)},
			},
			"pie": []gin.H{
				{"value": random(800, 400), "name": "搜索引擎"},
				{"value": random(500, 300), "name": "直接访问"},
				{"value": random(400, 200), "name": "邮件营销"},
				{"value": random(300, 200), "name": "联盟广告"},
				{"value": random(200, 150), "name": "视频广告"},
			},
			"updatedAt": time.Now().Format("15:04:05"), // 时:分:秒（在 data 内部！）
		},
		"msg": "success",
	})
}