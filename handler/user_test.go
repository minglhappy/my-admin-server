// handler/user_test.go
// 职责：演示 for range 的两种写法差异——为什么必须用 _, u := range
// 运行方式：go test ./handler -v -run TestRange
package handler

import "testing"

// TestRangeBasic：展示 range 返回 (下标, 元素) 两个值
// 参数：t 测试上下文，t.Log 输出的内容在 go test -v 时可见
// 实现逻辑：两个变量接收 range → i 是下标，u 是元素
func TestRangeBasic(t *testing.T) {
	for i, u := range users {
		// i 是 int 下标，u 是 UserItem 元素
		t.Logf("i=%d → u.Username=%s, u.Nickname=%s", i, u.Username, u.Nickname)
		if i == 1 && u.Username != "zhangsan" {
			t.Errorf("下标 1 应该是 zhangsan，实际 %s", u.Username)
		}
	}
}

// TestRangeSingleVar：展示只写一个变量时的陷阱——拿到的是下标而不是元素
// 参数：t 测试上下文
// 实现逻辑：for u := range users → u 的类型是 int（下标）
// 注意：循环里不能写 u.Username（int 没有字段），这里用 %T 打印类型来证明
func TestRangeSingleVar(t *testing.T) {
	for u := range users {
		// 关键证据：%T 打印 u 的实际类型 → int（下标），不是 UserItem
		if u == 0 {
			t.Logf("u 的类型是 %T，值是 %v（这是下标！不是用户）", u, u)
		}
		if u != 0 && u != 1 && u != 2 && u != 3 && u != 4 {
			t.Errorf("u 不是 0~4 的下标，而是 %v", u)
		}
	}
}

// TestRangeBlankIdent：展示 _ 的作用——吞掉不用的下标，只留元素
// 参数：t 测试上下文
// 实现逻辑：_, u := range → 下标扔进空白标识符，u 是真正的元素
// 注意：_ 不是变量，不能使用；声明的变量必须被使用，_ 是唯一例外
func TestRangeBlankIdent(t *testing.T) {
	count := 0
	for _, u := range users {
		// 这里 u 是 UserItem，可以正常访问字段
		count++
		if u.ID <= 0 {
			t.Errorf("用户 ID 应大于 0，实际 %d", u.ID)
		}
	}
	if count != len(users) {
		t.Errorf("遍历次数 %d != 用户数 %d", count, len(users))
	}
}

// TestRangeCopyTrap：展示 range 的复制陷阱——u 是元素副本，改它不影响原切片
// 参数：t 测试上下文
// 实现逻辑：循环里改 u.Nickname 后，users 原数据不变（因为 u 是值拷贝）
// 注意：想改原数据必须用下标：users[i].Nickname = "xxx"
func TestRangeCopyTrap(t *testing.T) {
	old := users[0].Nickname          // 先记住原来的昵称
	for _, u := range users {
		u.Nickname = "改了也没用" // 改的是副本，原切片不受影响
	}
	if users[0].Nickname != old {
		t.Errorf("u 是副本不该影响原切片，原值 %s 变成了 %s", old, users[0].Nickname)
	}
	t.Logf("改 u 后原数据仍是 %q → 证明 u 是副本", users[0].Nickname)

	// 对比：用下标访问才能改原数据
	for i := range users {
		users[i].Nickname = "改成功" // users[i] 才是切片里真实的那一份
	}
	if users[0].Nickname != "改成功" {
		t.Errorf("下标修改应生效，实际 %s", users[0].Nickname)
	}
	t.Logf("下标修改后原数据变成 %q → 证明 users[i] 才是真身", users[0].Nickname)
	users[0].Nickname = old // 测试完恢复原数据，避免影响其他测试
}
