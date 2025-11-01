package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// {{ AURA-X: Add - 正则表达式验证器模块. Confirmed via 寸止 }}

// RegexPattern 预定义正则模式
type RegexPattern struct {
	Name        string `json:"name"`        // 模式名称
	Pattern     string `json:"pattern"`     // 正则表达式
	Description string `json:"description"` // 描述
	Example     string `json:"example"`     // 示例匹配
}

// RegexTestResult 正则测试结果
type RegexTestResult struct {
	Success     bool     `json:"success"`     // 是否匹配成功
	Matched     string   `json:"matched"`     // 匹配到的内容
	Groups      []string `json:"groups"`      // 捕获组
	Error       string   `json:"error"`       // 错误信息
	OriginalName string  `json:"originalName"` // 原始文件名
	ExtractedNumber string `json:"extractedNumber"` // 提取的番号
}

// RegexValidator 正则验证器
type RegexValidator struct {
	patterns []RegexPattern
}

// NewRegexValidator 创建新的正则验证器
func NewRegexValidator() *RegexValidator {
	return &RegexValidator{
		patterns: GetDefaultPatterns(),
	}
}

// GetDefaultPatterns 获取默认的正则模式列表
// {{ AURA-X: Modify - 增强说明和示例. Confirmed via 寸止 }}
func GetDefaultPatterns() []RegexPattern {
	return []RegexPattern{
		{
			Name:        "标准格式",
			Pattern:     `(?i)([a-z]+)[-_](\d+)`,
			Description: "📌 最常用的番号格式\n• 字母+破折号/下划线+数字\n• 大小写不敏感\n• 自动规范化为大写+破折号",
			Example:     "✅ ABC-123, ipx-456, SSIS_789\n✅ 4k2.com@ipzz-655 (会自动清理前缀)",
		},
		{
			Name:        "复杂格式",
			Pattern:     `(?i)([a-z]+[-_][a-z]+)[-_](\d+)`,
			Description: "📌 多字母组合的番号\n• 两组字母+破折号+数字\n• 常见于特定厂商系列",
			Example:     "✅ MKY-NS-001, T28-123, ABP-XYZ-456",
		},
		{
			Name:        "FC2格式",
			Pattern:     `(?i)FC2[-_]?(?:PPV[-_]?)?(\d+)`,
			Description: "📌 FC2专用格式\n• 支持 FC2/FC2-PPV 等变体\n• 自动提取纯数字部分",
			Example:     "✅ FC2-1234567, FC2PPV-1234567\n✅ fc2_ppv_1234567",
		},
		{
			Name:        "纯数字格式",
			Pattern:     `^(\d{6,})$`,
			Description: "📌 FANZA CID等纯数字番号\n• 至少6位数字\n• 常见于FANZA官方编号",
			Example:     "✅ 123456, 1234567890\n❌ abc123 (必须纯数字)",
		},
		{
			Name:        "一本道/加勒比格式",
			Pattern:     `(?i)(\d{6})[-_](\d{3})`,
			Description: "📌 无码片商专用格式\n• 6位数字+破折号+3位数字\n• 一本道、加勒比、Pacopacomama等",
			Example:     "✅ 123456-789, 010122_001\n✅ 1pondo 123456_789",
		},
		{
			Name:        "Tokyo Hot格式",
			Pattern:     `(?i)(cz|gedo|k|n|red-|se)(\d{2,4})`,
			Description: "📌 Tokyo Hot系列专用\n• 特定字母前缀+2-4位数字\n• n/k/cz/red等系列",
			Example:     "✅ n1234, k0123, red-123\n✅ cz012, se0456",
		},
		{
			Name:        "Heyzo格式",
			Pattern:     `(?i)heyzo[-_]?(\d{4})`,
			Description: "📌 Heyzo站点专用\n• heyzo+4位数字\n• 支持带/不带破折号",
			Example:     "✅ HEYZO-1234, heyzo1234\n✅ Heyzo_2345",
		},
		{
			Name:        "X-Art格式",
			Pattern:     `(?i)x-art\.(\d{2})\.(\d{2})\.(\d{2})`,
			Description: "📌 X-Art站点专用\n• 日期格式: YY.MM.DD\n• 欧美高端系列",
			Example:     "✅ x-art.20.01.15\n✅ X-Art.21.12.25",
		},
		{
			Name:        "Heydouga格式",
			Pattern:     `(?i)heydouga[-_]?(\d{4})[-_](\d{3,5})`,
			Description: "📌 Heydouga站点专用\n• heydouga+4位+3-5位数字\n• 素人投稿系列",
			Example:     "✅ heydouga-4030-1234\n✅ Heydouga_4017_12345",
		},
		{
			Name:        "通用提取（带捕获组）",
			Pattern:     `([A-Z]{2,}-\d{3,})`,
			Description: "📌 严格格式提取\n• 至少2个大写字母+破折号+至少3个数字\n• 适合已规范化的文件名",
			Example:     "✅ ABC-123, ABCD-1234\n❌ abc-123 (需大写)",
		},
		{
			Name:        "带网站前缀格式",
			Pattern:     `(?i)(?:\w+\.(?:com|net|cc|org|xyz)@)?([a-z]{3,}[-_]\d{3,})`,
			Description: "📌 自动清理网站前缀\n• 支持 xxx.com@, xxx.net@ 等\n• 提取真正的番号部分",
			Example:     "✅ 4k2.com@ipzz-655 → IPZZ-655\n✅ xxx.net@abc-123 → ABC-123",
		},
	}
}

// ValidateRegex 验证正则表达式语法
func (rv *RegexValidator) ValidateRegex(pattern string) (bool, string) {
	if pattern == "" {
		return false, "正则表达式不能为空"
	}
	
	_, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Sprintf("正则表达式语法错误: %v", err)
	}
	
	return true, "正则表达式语法正确"
}

// TestRegex 测试正则表达式是否能匹配给定的文件名
// {{ AURA-X: Modify - 支持取最后一个匹配，并使用normalizeNumber规范化. Confirmed via 寸止 }}
func (rv *RegexValidator) TestRegex(pattern string, filename string) *RegexTestResult {
	result := &RegexTestResult{
		OriginalName: filename,
	}
	
	// 验证正则语法
	re, err := regexp.Compile(pattern)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("正则表达式语法错误: %v", err)
		return result
	}
	
	// 使用FindAllStringSubmatch找到所有匹配
	allMatches := re.FindAllStringSubmatch(filename, -1)
	if len(allMatches) == 0 {
		result.Success = false
		result.Error = "未匹配到任何内容"
		return result
	}
	
	// 取最后一个匹配（类似Python版本的match[-1]）
	matches := allMatches[len(allMatches)-1]
	
	result.Success = true
	result.Matched = matches[0] // 完整匹配
	
	// 捕获组
	if len(matches) > 1 {
		result.Groups = matches[1:]
		// 使用第一个捕获组作为提取的番号，并规范化
		rawNumber := matches[1]
		result.ExtractedNumber = normalizeNumberForTest(rawNumber)
	} else {
		// 如果没有捕获组则使用整个匹配，并规范化
		result.ExtractedNumber = normalizeNumberForTest(matches[0])
	}
	
	return result
}

// normalizeNumberForTest 测试用的规范化函数（与NumberParser.normalizeNumber逻辑一致）
func normalizeNumberForTest(number string) string {
	// 1. 下划线统一转为破折号
	number = strings.ReplaceAll(number, "_", "-")
	
	// 2. 移除常见前缀
	prefixesToRemove := []string{
		"ppv-", "PPV-",
		"fc-", "FC-",
	}
	for _, prefix := range prefixesToRemove {
		if strings.HasPrefix(strings.ToLower(number), strings.ToLower(prefix)) {
			number = number[len(prefix):]
		}
	}
	
	// 3. 移除末尾的破折号
	number = strings.TrimSuffix(number, "-")
	
	// 4. 处理无破折号格式: abc234 -> ABC-234
	noDashRegex := regexp.MustCompile(`^([a-zA-Z]{3,})(\d{3,})$`)
	if matches := noDashRegex.FindStringSubmatch(number); len(matches) == 3 {
		number = strings.ToUpper(matches[1]) + "-" + matches[2]
	} else {
		// 否则统一转大写
		number = strings.ToUpper(number)
	}
	
	return number
}

// TestMultipleFiles 测试正则表达式对多个文件的匹配效果
func (rv *RegexValidator) TestMultipleFiles(pattern string, filenames []string) []RegexTestResult {
	results := make([]RegexTestResult, 0, len(filenames))
	
	for _, filename := range filenames {
		result := rv.TestRegex(pattern, filename)
		results = append(results, *result)
	}
	
	return results
}

// GetDefaultPattern 根据名称获取默认正则模式
func (rv *RegexValidator) GetDefaultPattern(name string) (RegexPattern, bool) {
	for _, pattern := range rv.patterns {
		if pattern.Name == name {
			return pattern, true
		}
	}
	return RegexPattern{}, false
}

// GetAllPatterns 获取所有预定义模式
func (rv *RegexValidator) GetAllPatterns() []RegexPattern {
	return rv.patterns
}

// ExtractNumberWithPattern 使用指定正则模式提取番号
func ExtractNumberWithPattern(pattern string, filename string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("正则表达式编译失败: %v", err)
	}
	
	matches := re.FindStringSubmatch(filename)
	if len(matches) == 0 {
		return "", fmt.Errorf("未匹配到任何内容")
	}
	
	// 优先返回第一个捕获组，如果没有捕获组则返回完整匹配
	if len(matches) > 1 {
		return strings.ToUpper(matches[1]), nil
	}
	
	return strings.ToUpper(matches[0]), nil
}

// SuggestPattern 根据文件名建议合适的正则模式
func (rv *RegexValidator) SuggestPattern(filename string) []RegexPattern {
	suggestions := make([]RegexPattern, 0)
	
	for _, pattern := range rv.patterns {
		result := rv.TestRegex(pattern.Pattern, filename)
		if result.Success {
			suggestions = append(suggestions, pattern)
		}
	}
	
	return suggestions
}

