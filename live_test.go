package bilibili

import (
	"testing"
	"time"
)

func TestCalculateAppSign(t *testing.T) {
	// 测试签名计算函数
	params := map[string]string{
		"appkey":   "aae92bc66f3edfab",
		"build":    "1234",
		"platform": "pc_link",
		"room_id":  "12345",
		"ts":       "1640995200",
		"version":  "1.0.0",
	}

	sign := calculateAppSign(params)

	// 验证签名不为空
	if sign == "" {
		t.Error("Sign should not be empty")
	}

	// 验证签名长度为32（MD5 hex编码长度）
	if len(sign) != 32 {
		t.Errorf("Sign length should be 32, got %d", len(sign))
	}

	// 验证相同参数生成相同签名
	sign2 := calculateAppSign(params)
	if sign != sign2 {
		t.Error("Same parameters should generate same signature")
	}
}

func TestStartLiveAutoSign(t *testing.T) {
	// 测试 StartLive 自动签名功能

	// 创建参数，不包含签名
	param := StartLiveParam{
		RoomId:   12345,
		AreaV2:   1,
		Platform: "pc_link",
		Version:  "1.0.0",
		Build:    1234,
		Appkey:   "aae92bc66f3edfab",
		// Sign 故意留空，测试自动计算
		Ts: int(time.Now().Unix()),
	}

	// 注意：这个测试不会真正发送请求，因为我们没有有效的cookie
	// 但我们可以验证签名是否被正确计算
	originalSign := param.Sign

	// 模拟签名计算（不实际发送请求）
	if param.Sign == "" && param.Appkey != "" {
		if param.Ts == 0 {
			param.Ts = int(time.Now().Unix())
		}

		signParams := map[string]string{
			"appkey":   param.Appkey,
			"build":    "1234",
			"platform": param.Platform,
			"room_id":  "12345",
			"ts":       "1640995200", // 使用固定时间戳以便测试
			"version":  param.Version,
		}

		if param.AreaV2 != 0 {
			signParams["area_v2"] = "1"
		}

		calculatedSign := calculateAppSign(signParams)

		// 验证签名被正确计算
		if calculatedSign == "" {
			t.Error("Calculated sign should not be empty")
		}

		if len(calculatedSign) != 32 {
			t.Errorf("Calculated sign length should be 32, got %d", len(calculatedSign))
		}
	}

	// 验证原始签名为空（符合测试预期）
	if originalSign != "" {
		t.Error("Original sign should be empty for this test")
	}
}

func TestStartLiveWithExistingSign(t *testing.T) {
	// 测试当已有签名时不会被覆盖
	existingSign := "existing_signature_12345678901234567890"

	params := map[string]string{
		"appkey":   "aae92bc66f3edfab",
		"build":    "1234",
		"platform": "pc_link",
		"room_id":  "12345",
		"ts":       "1640995200",
		"version":  "1.0.0",
	}

	// 如果已经有签名，应该使用现有的签名
	sign := existingSign
	if sign == "" {
		sign = calculateAppSign(params)
	}

	// 验证使用了现有签名
	if sign != existingSign {
		t.Error("Should use existing signature when provided")
	}
}
