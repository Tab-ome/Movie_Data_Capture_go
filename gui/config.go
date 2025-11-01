package gui

import (
	"fmt"
	"os"
	"path/filepath"

	"movie-data-capture/internal/config"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

// {{ AURA-X: Add - 配置管理API接口. Confirmed via 寸止 }}

// LoadConfig 加载配置文件
func (a *App) LoadConfig() (map[string]interface{}, error) {
	cfg, err := config.Load(a.configPath)
	if err != nil {
		a.SendLog("ERROR", fmt.Sprintf("加载配置失败: %v", err))
		return nil, err
	}
	
	a.config = cfg
	a.SendLog("INFO", "配置文件加载成功")
	
	// 将配置转换为map返回给前端
	configMap := make(map[string]interface{})
	data, _ := yaml.Marshal(cfg)
	yaml.Unmarshal(data, &configMap)
	
	return configMap, nil
}

// SaveConfig 保存配置到文件
func (a *App) SaveConfig(configData map[string]interface{}) error {
	// 将map转换回配置结构
	data, err := yaml.Marshal(configData)
	if err != nil {
		a.SendLog("ERROR", fmt.Sprintf("配置序列化失败: %v", err))
		return err
	}
	
	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		a.SendLog("ERROR", fmt.Sprintf("配置解析失败: %v", err))
		return err
	}
	
	// 验证配置
	if err := a.ValidateConfig(&cfg); err != nil {
		a.SendLog("ERROR", fmt.Sprintf("配置验证失败: %v", err))
		return err
	}
	
	// 保存到文件
	file, err := os.Create(a.configPath)
	if err != nil {
		a.SendLog("ERROR", fmt.Sprintf("创建配置文件失败: %v", err))
		return err
	}
	defer file.Close()
	
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	if err := encoder.Encode(&cfg); err != nil {
		a.SendLog("ERROR", fmt.Sprintf("写入配置文件失败: %v", err))
		return err
	}
	
	a.config = &cfg
	a.SendLog("INFO", "配置保存成功")
	
	return nil
}

// GetConfig 获取当前配置
func (a *App) GetConfig() (map[string]interface{}, error) {
	if a.config == nil {
		return a.LoadConfig()
	}
	
	configMap := make(map[string]interface{})
	data, _ := yaml.Marshal(a.config)
	yaml.Unmarshal(data, &configMap)
	
	return configMap, nil
}

// ValidateConfig 验证配置的有效性
func (a *App) ValidateConfig(cfg *config.Config) error {
	// 检查必填字段
	if cfg.Common.SourceFolder == "" {
		return fmt.Errorf("源文件夹不能为空")
	}
	
	if cfg.Common.SuccessOutputFolder == "" {
		return fmt.Errorf("输出文件夹不能为空")
	}
	
	// 检查文件夹是否存在
	if _, err := os.Stat(cfg.Common.SourceFolder); os.IsNotExist(err) {
		return fmt.Errorf("源文件夹不存在: %s", cfg.Common.SourceFolder)
	}
	
	// 检查数值范围
	if cfg.Common.MainMode < 0 || cfg.Common.MainMode > 3 {
		return fmt.Errorf("运行模式值无效 (0-3): %d", cfg.Common.MainMode)
	}
	
	if cfg.Common.LinkMode < 0 || cfg.Common.LinkMode > 2 {
		return fmt.Errorf("链接模式值无效 (0-2): %d", cfg.Common.LinkMode)
	}
	
	return nil
}

// SelectFolder 打开文件夹选择对话框
func (a *App) SelectFolder(title string) (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	
	if err != nil {
		a.SendLog("ERROR", fmt.Sprintf("选择文件夹失败: %v", err))
		return "", err
	}
	
	if folder != "" {
		a.SendLog("INFO", fmt.Sprintf("已选择文件夹: %s", folder))
	}
	
	return folder, nil
}

// GetConfigValue 获取特定配置值
func (a *App) GetConfigValue(key string) interface{} {
	if a.config == nil {
		a.LoadConfig()
		if a.config == nil {
			return nil
		}
	}
	
	// 这里可以根据key返回特定的配置值
	// 示例实现，实际可以更复杂
	configMap := make(map[string]interface{})
	data, _ := yaml.Marshal(a.config)
	yaml.Unmarshal(data, &configMap)
	
	return configMap[key]
}

// SetConfigValue 设置特定配置值
func (a *App) SetConfigValue(key string, value interface{}) error {
	if a.config == nil {
		if _, err := a.LoadConfig(); err != nil {
			return err
		}
	}
	
	// 这里需要根据key设置特定的配置值
	// 实际实现会更复杂，需要反射或类型断言
	
	return nil
}

// ResetConfig 重置配置为默认值
func (a *App) ResetConfig() error {
	// 备份当前配置
	backupPath := a.configPath + ".backup"
	if err := copyFile(a.configPath, backupPath); err != nil {
		a.SendLog("WARN", fmt.Sprintf("备份配置文件失败: %v", err))
	} else {
		a.SendLog("INFO", fmt.Sprintf("配置已备份到: %s", backupPath))
	}
	
	// 复制模板配置
	templatePath := "config_template.yaml"
	if _, err := os.Stat(templatePath); err == nil {
		if err := copyFile(templatePath, a.configPath); err != nil {
			a.SendLog("ERROR", fmt.Sprintf("重置配置失败: %v", err))
			return err
		}
		a.SendLog("INFO", "配置已重置为默认值")
		_, err := a.LoadConfig()
		return err
	}
	
	return fmt.Errorf("配置模板文件不存在")
}

// copyFile 复制文件的辅助函数
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(dst, data, 0644)
}

