// package idgen

// import (
// 	"fmt"
// 	"sync"
// 	"testing"
// 	"time"

// 	"go-protos/pkg/idgen"
// )

// func TestTime(t *testing.T) {
// 	// 获取当前时间
// 	now := time.Now()

// 	// 测试本地时间
// 	t.Run("LocalTime", func(t *testing.T) {
// 		localTime := now.Local()
// 		t.Logf("Local time: %s", localTime)
// 		name, offset := localTime.Zone()
// 		t.Logf("Local zone: %s (UTC offset: %d)", name, offset)
// 	})

// 	// 测试UTC时间
// 	t.Run("UTCTime", func(t *testing.T) {
// 		utcTime := now.UTC()
// 		t.Logf("UTC time: %s", utcTime)
// 		name, offset := utcTime.Zone()
// 		t.Logf("UTC zone: %s (UTC offset: %d)", name, offset)
// 	})

// 	// 测试时间戳
// 	t.Run("Timestamp", func(t *testing.T) {
// 		unixSeconds := now.Unix()
// 		unixMilliseconds := now.UnixMilli()
// 		t.Logf("Unix seconds: %d", unixSeconds)
// 		t.Logf("Unix milliseconds: %d", unixMilliseconds)
// 	})

// 	// 测试时区转换
// 	t.Run("TimeZoneConversion", func(t *testing.T) {
// 		shanghaiLoc, err := time.LoadLocation("Asia/Shanghai")
// 		if err != nil {
// 			t.Fatalf("Failed to load Shanghai location: %v", err)
// 		}

// 		newYorkLoc, err := time.LoadLocation("America/New_York")
// 		if err != nil {
// 			t.Fatalf("Failed to load New York location: %v", err)
// 		}

// 		shanghaiTime := now.In(shanghaiLoc)
// 		newYorkTime := now.In(newYorkLoc)

// 		t.Logf("Shanghai time: %s", shanghaiTime)
// 		t.Logf("New York time: %s", newYorkTime)
// 		t.Logf("Time difference: %s", shanghaiTime.Sub(newYorkTime))
// 	})

// 	// 测试时间格式化
// 	t.Run("TimeFormatting", func(t *testing.T) {
// 		t.Logf("RFC3339 format: %s", now.Format(time.RFC3339))
// 		t.Logf("Custom format: %s", now.Format("2006-01-02 15:04:05 MST"))
// 	})
// }

// func TestIDGenerator(t *testing.T) {
// 	// 初始化生成器（每个服务节点需要不同的nodeID）
// 	generator := idgen.NewIDGenerator()

// 	// 生成用户ID
// 	userID, _ := generator.Generate(idgen.UserType)

// 	// 生成账户ID
// 	accountID, _ := generator.Generate(idgen.AccountType)

// 	// 生成账户ID
// 	orderID, _ := generator.Generate(idgen.OrderType)

// 	fmt.Println(userID, accountID, orderID)
// 	// 解析ID
// 	createTime, bizType, seq := idgen.ParseID(userID)
// 	fmt.Println(createTime, bizType, seq)
// }

// // | 32位时间戳（毫秒） | 8位节点ID | 4位业务类型 | 20位序列号 |
// // 50057639 8670 6137088
// // 50057639 8670 8234241

// func TestTimeZoneHandling(t *testing.T) {
// 	gen := idgen.NewIDGenerator()
// 	// 测试不同时区下生成ID的一致性
// 	t.Run("UTC", func(t *testing.T) {
// 		id, err := gen.Generate(idgen.UserType)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		tm, _, _ := idgen.ParseID(id)
// 		t.Logf("UTC time: %s", tm)
// 	})

// 	t.Run("Shanghai", func(t *testing.T) {
// 		loc, _ := time.LoadLocation("Asia/Shanghai")
// 		time.Local = loc // 模拟服务器在上海时区
// 		id, err := gen.Generate(idgen.UserType)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		tm, _, _ := idgen.ParseID(id)
// 		t.Logf("Parsed UTC time: %s", tm)
// 		t.Logf("Local time: %s", tm.In(loc))
// 	})
// }

// func TestIDGenerationAndParsing(t *testing.T) {
// 	gen := idgen.NewIDGenerator()

// 	id, err := gen.Generate(idgen.UserType)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	parsedTime, bizType, seq := idgen.ParseIDLocal(id)

// 	// 验证时间
// 	now := time.Now().UTC()
// 	if parsedTime.Sub(now).Abs() > time.Second {
// 		t.Errorf("解析时间与当前时间差异过大: %s vs %s", parsedTime, now)
// 	}

// 	if bizType != idgen.UserType {
// 		t.Errorf("业务类型不匹配: 期望 %d, 得到 %d", idgen.UserType, bizType)
// 	}
// 	if seq < 0 || seq >= (1<<24) {
// 		t.Errorf("序列号超出范围: %d", seq)
// 	}

// 	t.Logf("生成ID: %d", id)
// 	t.Logf("解析结果: 时间=%s, 业务类型=%d, 序列号=%d",
// 		parsedTime, bizType, seq)
// }

// func TestMaxIDParsing(t *testing.T) {
// 	parsedTime, bizType, seq := idgen.ParseIDLocal(18446744073709551615)

// 	if seq < 0 || seq >= (1<<24) {
// 		t.Errorf("序列号超出范围: %d", seq)
// 	}

// 	// t.Logf("生成ID: %d", id)
// 	t.Logf("解析结果: 时间=%s, 业务类型=%d, 序列号=%d",
// 		parsedTime, bizType, seq)
// }

// func TestConcurrentSafety(t *testing.T) {
// 	gen := idgen.NewIDGenerator()
// 	const goroutines = 5
// 	const perRoutine = 10
// 	ids := make(chan uint64, goroutines*perRoutine)

// 	var wg sync.WaitGroup
// 	wg.Add(goroutines)

// 	// 并发生成
// 	for range goroutines {
// 		go func() {
// 			defer wg.Done()
// 			for range perRoutine {
// 				id, err := gen.Generate(idgen.UserType)
// 				if err != nil {
// 					t.Error(err)
// 					return
// 				}
// 				ids <- id
// 			}
// 		}()
// 	}
// 	wg.Wait()
// 	close(ids)

// 	// 验证唯一性
// 	idSet := make(map[uint64]bool)
// 	for id := range ids {
// 		// fmt.Println(id)
// 		if idSet[id] {
// 			t.Fatal("发现重复ID")
// 		}
// 		idSet[id] = true
// 	}

// 	for id := range idSet {
// 		parsedTime, bizType, seq := idgen.ParseIDLocal(id)

// 		t.Logf("生成ID: %d, 解析结果: 时间=%s, 业务类型=%d, 序列号=%d",
// 			id, parsedTime, bizType, seq)
// 	}
// }

// func TestPerturbationRecovery(t *testing.T) {
// 	// 测试扰动因子正确还原
// 	testCases := []struct {
// 		inputID   uint64
// 		expectSeq int
// 	}{
// 		{0x12345678<<16 | 0x3FF<<6 | 0x00, 0x00 - 0x3FF}, // 负值测试
// 		{0xABCDEF12<<16 | 0x200<<6 | 0x7F, 0x7F - 0x200},
// 		{0xDEADBEEF<<16 | 0x100<<6 | 0xFF, 0xFF - 0x100},
// 	}

// 	for _, tc := range testCases {
// 		_, _, actualSeq := idgen.ParseID(tc.inputID)
// 		expected := (tc.expectSeq + (1 << 16)) % (1 << 16)
// 		if actualSeq != expected {
// 			t.Errorf("扰动还原失败\n输入: %016X\n期望: %d\n实际: %d",
// 				tc.inputID, expected, actualSeq)
// 		}
// 	}
// }

// func TestPerturbationEffect(t *testing.T) {
// 	gen := idgen.NewIDGenerator()

// 	// 生成两个连续ID
// 	id1, _ := gen.Generate(idgen.UserType)
// 	id2, _ := gen.Generate(idgen.UserType)

// 	// 提取序列号部分
// 	seq1 := id1 & 0xFFFF
// 	seq2 := id2 & 0xFFFF

// 	// 计算实际步长
// 	perturb := (id1 >> 6) & 0x3FF
// 	expectedStep := (1 + perturb) % (1 << 16)
// 	actualStep := seq2 - seq1

// 	if actualStep != expectedStep {
// 		t.Errorf("扰动验证失败\n期望步长: %d\n实际步长: %d",
// 			expectedStep, actualStep)
// 	}
// }

// func TestMaxID(t *testing.T) {
// 	// 生成最大合法ID
// 	maxTimePart := uint64(1<<idgen.GetTimeBits() - 1)
// 	maxBizType := uint64(1<<idgen.GetBusinessBits() - 1)
// 	maxSequence := uint64(1<<idgen.GetSequenceBits() - 1)

// 	maxID := (maxTimePart << (idgen.GetBusinessBits() + idgen.GetSequenceBits())) |
// 		(maxBizType << idgen.GetSequenceBits()) |
// 		maxSequence

// 	// 验证无符号转换
// 	if uint64(maxID) < 0 {
// 		t.Fatal("生成的ID超过有符号BIGINT范围")
// 	}

// 	t.Logf("最大ID值: %d (0x%X)", maxID, maxID)
// 	t.Logf("BIGINT UNSIGNED最大值: 18446744073709551615")
// }
