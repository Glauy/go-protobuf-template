package idgen

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

/*
bigint 64位 最大值92233720 3685 47 75807
| 32位时间戳（毫秒） | 8位节点ID | 4位业务类型 | 20位序列号 |
50057639 8670 61 37088
50057639 8670 82 34241
92233720 3685 47 75807
*/

// 业务类型定义示例（保留注释与说明在上方）
/*
type BizType int

const (
	UserType       BizType = 0x1 // 用户业务
	ProfileType    BizType = 0x2 // 资料业务
	AccountType    BizType = 0x3 // 账户业务
	MembershipType BizType = 0x4 // 会员业务
	OrderType      BizType = 0x5 // 订单业务
)
*/

const (
	// 新位分配结构： (已调整业务类型到高位)
	// | 4位业务类型 | 44位时间戳（毫秒） | 16位序列号 |
	businessBits = 4  // 保持16种业务类型
	timeBits     = 44 // 可表示约557年（2^44/1000/31536000 ≈ 557年）
	sequenceBits = 16 // 保持6.5万/ms并发能力

	// 调整 epoch 为 2023-01-01 00:00:00 UTC 的毫秒时间戳
	epoch = 1672531200000 // 2023-01-01 UTC 的毫秒时间戳
)

// BizType 表示业务类型枚举（类型安全）
type BizType int

const (
	UserType       BizType = 0x1 // 用户业务
	ProfileType    BizType = 0x2 // 资料业务
	AccountType    BizType = 0x3 // 账户业务
	MembershipType BizType = 0x4 // 会员业务
	OrderType      BizType = 0x5 // 订单业务
)

// masks & shifts (uint64 safe)
var (
	sequenceMask uint64 = (uint64(1) << sequenceBits) - 1
	timeMask     uint64 = (uint64(1) << timeBits) - 1
	bizMask      uint64 = (uint64(1) << businessBits) - 1

	sequenceShift = 0
	timeShift     = sequenceBits
	bizShift      = sequenceBits + timeBits
)

type IDGenerator struct {
	mu       sync.Mutex
	lastTime int64
	sequence uint64
}

// NewIDGenerator 创建一个 ID 生成器实例（简单实现，不包含 nodeID）
// 若需要多节点安全（节点区分），可扩展此构造函数接受 nodeID 并把 nodeID 放入高位或专用位。
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// Generate 生成安全 ID（将 BizType 放在高位）
// 返回 uint64：| biz(4) | time(44) | seq(16) |
func (g *IDGenerator) Generate(biz BizType) (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UTC().UnixMilli()
	timePart := now - epoch

	if timePart < 0 {
		return 0, errors.New("系统时钟异常：当前时间早于 epoch 时间")
	}
	if uint64(timePart) > timeMask {
		return 0, errors.New("时间戳溢出，ID 生成器已达最大时间限制")
	}
	if biz < 0 || uint64(biz) > bizMask {
		return 0, fmt.Errorf("业务类型必须在0-%d之间", bizMask)
	}

	// 序列处理：同一毫秒内递增
	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & sequenceMask
		if g.sequence == 0 {
			// 序列溢出，等待到下一个毫秒
			for now <= g.lastTime {
				now = time.Now().UTC().UnixMilli()
			}
			timePart = now - epoch
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	id := (uint64(biz) << bizShift) |
		(uint64(timePart) << timeShift) |
		(g.sequence << sequenceShift)

	return id, nil
}

// ParseID 解析 ID，返回 (UTC time, BizType, sequence)
func ParseID(id uint64) (time.Time, BizType, int) {
	biz := BizType((id >> bizShift) & bizMask)
	timePart := (id >> timeShift) & timeMask
	seq := int((id >> sequenceShift) & sequenceMask)

	actualTime := int64(timePart) + epoch
	return time.UnixMilli(actualTime).UTC(), biz, seq
}

// ParseIDLocal 返回本地时区时间
func ParseIDLocal(id uint64) (time.Time, BizType, int) {
	utcTime, biz, seq := ParseID(id)
	return utcTime.In(time.Local), biz, seq
}

// ParseIDInLocation 支持自定义时区
func ParseIDInLocation(id uint64, loc *time.Location) (time.Time, BizType, int) {
	utcTime, biz, seq := ParseID(id)
	return utcTime.In(loc), biz, seq
}

// 辅助查询函数
func GetTimeBits() int {
	return timeBits
}

func GetBusinessBits() int {
	return businessBits
}

func GetSequenceBits() int {
	return sequenceBits
}
