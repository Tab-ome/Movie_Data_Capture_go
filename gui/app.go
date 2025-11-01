package gui

import (
	"context"
	"fmt"
	"time"

	"movie-data-capture/internal/config"
	"movie-data-capture/pkg/logger"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// {{ AURA-X: Add - 创建Wails应用核心结构. Confirmed via 寸止 }}
// App 结构体包含应用的核心状态和配置
type App struct {
	ctx        context.Context
	config     *config.Config
	configPath string
	isRunning  bool
	stats      *Stats
	runner     *Runner
}

// Stats 存储运行时统计信息
type Stats struct {
	Total     int       `json:"total"`
	Success   int       `json:"success"`
	Failed    int       `json:"failed"`
	Skipped   int       `json:"skipped"`
	StartTime time.Time `json:"startTime"`
	Duration  string    `json:"duration"`
}

// NewApp 创建新的应用实例
func NewApp() *App {
	return &App{
		configPath: "config.yaml",
		stats: &Stats{
			Total:   0,
			Success: 0,
			Failed:  0,
			Skipped: 0,
		},
	}
}

// Startup 在应用启动时被调用
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	logger.Info("[GUI] 应用启动中...")
	
	// 加载配置
	cfg, err := config.Load(a.configPath)
	if err != nil {
		logger.Error("[GUI] 加载配置失败: %v", err)
		runtime.EventsEmit(ctx, "log", map[string]interface{}{
			"level":   "ERROR",
			"message": fmt.Sprintf("加载配置失败: %v", err),
			"time":    time.Now().Format("15:04:05"),
		})
	} else {
		a.config = cfg
		logger.Info("[GUI] 配置加载成功")
		runtime.EventsEmit(ctx, "log", map[string]interface{}{
			"level":   "INFO",
			"message": "配置加载成功",
			"time":    time.Now().Format("15:04:05"),
		})
	}
	
	// 初始化Runner
	a.runner = NewRunner(a)
}

// DomReady 在前端DOM准备好后被调用
func (a *App) DomReady(ctx context.Context) {
	logger.Info("[GUI] 前端界面已就绪")
}

// Shutdown 在应用关闭前被调用
func (a *App) Shutdown(ctx context.Context) {
	logger.Info("[GUI] 应用关闭中...")
	if a.isRunning {
		a.runner.Stop()
	}
}

// SendLog 向前端发送日志消息
func (a *App) SendLog(level, message string) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "log", map[string]interface{}{
			"level":   level,
			"message": message,
			"time":    time.Now().Format("15:04:05"),
		})
	}
}

// SendProgress 向前端发送进度更新
func (a *App) SendProgress() {
	if a.ctx != nil {
		duration := ""
		if !a.stats.StartTime.IsZero() {
			duration = time.Since(a.stats.StartTime).Round(time.Second).String()
		}
		
		runtime.EventsEmit(a.ctx, "progress", map[string]interface{}{
			"total":    a.stats.Total,
			"success":  a.stats.Success,
			"failed":   a.stats.Failed,
			"skipped":  a.stats.Skipped,
			"running":  a.isRunning,
			"duration": duration,
		})
	}
}

// SendFileStatus 向前端发送单个文件的处理状态
// {{ AURA-X: Add - 实时发送文件处理状态. Confirmed via 寸止 }}
func (a *App) SendFileStatus(fileInfo *FileInfo) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "file_status", map[string]interface{}{
			"path":     fileInfo.Path,
			"name":     fileInfo.Name,
			"size":     fileInfo.Size,
			"number":   fileInfo.Number,
			"status":   fileInfo.Status,
			"error":    fileInfo.Error,
			"duration": fileInfo.Duration,
		})
	}
}

// GetStats 获取当前统计信息
func (a *App) GetStats() *Stats {
	if !a.stats.StartTime.IsZero() {
		a.stats.Duration = time.Since(a.stats.StartTime).Round(time.Second).String()
	}
	return a.stats
}

// IsRunning 返回是否正在运行
func (a *App) IsRunning() bool {
	return a.isRunning
}

