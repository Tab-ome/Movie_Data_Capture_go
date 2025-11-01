package gui

import (
	"movie-data-capture/internal/config"
	"movie-data-capture/pkg/parser"
)

// {{ AURA-X: Add - 正则表达式测试API接口. Confirmed via 寸止 }}

// RegexTestRequest 正则测试请求
type RegexTestRequest struct {
	Pattern   string   `json:"pattern"`   // 正则表达式
	Filenames []string `json:"filenames"` // 要测试的文件名列表
}

// GetDefaultRegexPatterns 获取预定义的正则模式列表
func (a *App) GetDefaultRegexPatterns() []parser.RegexPattern {
	validator := parser.NewRegexValidator()
	return validator.GetAllPatterns()
}

// ValidateRegex 验证正则表达式语法
func (a *App) ValidateRegex(pattern string) map[string]interface{} {
	validator := parser.NewRegexValidator()
	valid, message := validator.ValidateRegex(pattern)
	
	return map[string]interface{}{
		"valid":   valid,
		"message": message,
	}
}

// TestRegexPattern 测试正则表达式对文件名的匹配效果
func (a *App) TestRegexPattern(request RegexTestRequest) []parser.RegexTestResult {
	validator := parser.NewRegexValidator()
	return validator.TestMultipleFiles(request.Pattern, request.Filenames)
}

// SuggestRegexPattern 根据文件名建议合适的正则模式
func (a *App) SuggestRegexPattern(filename string) []parser.RegexPattern {
	validator := parser.NewRegexValidator()
	return validator.SuggestPattern(filename)
}

// TestNumberExtraction 测试番号提取效果（使用完整的NumberParser）
func (a *App) TestNumberExtraction(filename string, customRegex string) map[string]interface{} {
	// 创建临时配置
	tempConfig := a.config
	if tempConfig == nil {
		// 如果没有配置，创建一个默认配置
		tempConfig = &config.Config{}
	}
	
	// 设置自定义正则
	if customRegex != "" {
		tempConfig.NameRule.NumberRegexs = customRegex
	}
	
	// 使用NumberParser提取番号
	numberParser := parser.NewNumberParser(tempConfig)
	extractedNumber := numberParser.GetNumber(filename)
	
	return map[string]interface{}{
		"filename":        filename,
		"extractedNumber": extractedNumber,
		"customRegex":     customRegex,
	}
}

