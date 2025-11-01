package gui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"movie-data-capture/internal/core"
	"movie-data-capture/pkg/parser"
)

// {{ AURA-X: Add - 任务运行控制接口. Confirmed via 寸止 }}

// Runner 任务运行器
type Runner struct {
	app       *App
	ctx       context.Context
	cancel    context.CancelFunc
	processor *core.Processor
	mu        sync.Mutex
}

// FileInfo 文件信息结构
type FileInfo struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Number   string `json:"number"`
	Status   string `json:"status"` // pending, processing, success, failed, skipped
	Error    string `json:"error"`
	Duration string `json:"duration"`
}

// NewRunner 创建新的运行器
func NewRunner(app *App) *Runner {
	return &Runner{
		app: app,
	}
}

// Start 启动刮削任务
func (a *App) Start(sourcePath string) error {
	a.runner.mu.Lock()
	defer a.runner.mu.Unlock()
	
	if a.isRunning {
		return fmt.Errorf("任务正在运行中")
	}
	
	if a.config == nil {
		return fmt.Errorf("配置未加载")
	}
	
	// 如果提供了源路径，使用它；否则使用配置中的路径
	if sourcePath != "" {
		a.config.Common.SourceFolder = sourcePath
	}
	
	// 验证源文件夹
	if _, err := os.Stat(a.config.Common.SourceFolder); os.IsNotExist(err) {
		return fmt.Errorf("源文件夹不存在: %s", a.config.Common.SourceFolder)
	}
	
	// 重置统计信息
	a.stats = &Stats{
		Total:     0,
		Success:   0,
		Failed:    0,
		Skipped:   0,
		StartTime: time.Now(),
	}
	
	// 创建取消上下文
	a.runner.ctx, a.runner.cancel = context.WithCancel(context.Background())
	
	// 标记为运行中
	a.isRunning = true
	
	a.SendLog("INFO", fmt.Sprintf("开始处理文件夹: %s", a.config.Common.SourceFolder))
	a.SendProgress()
	
	// 在goroutine中运行任务
	go a.runner.run()
	
	return nil
}

// Stop 停止当前任务
func (a *App) Stop() error {
	return a.runner.Stop()
}

// Stop 停止运行器
func (r *Runner) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if !r.app.isRunning {
		return fmt.Errorf("没有正在运行的任务")
	}
	
	r.app.SendLog("WARN", "正在停止任务...")
	
	if r.cancel != nil {
		r.cancel()
	}
	
	r.app.isRunning = false
	r.app.SendLog("INFO", "任务已停止")
	r.app.SendProgress()
	
	return nil
}

// run 实际执行任务的函数
func (r *Runner) run() {
	defer func() {
		r.app.isRunning = false
		r.app.SendProgress()
		
		if rec := recover(); rec != nil {
			r.app.SendLog("ERROR", fmt.Sprintf("任务异常终止: %v", rec))
		}
	}()
	
	// 获取文件列表
	files, err := r.getVideoFiles(r.app.config.Common.SourceFolder)
	if err != nil {
		r.app.SendLog("ERROR", fmt.Sprintf("获取文件列表失败: %v", err))
		return
	}
	
	r.app.stats.Total = len(files)
	r.app.SendLog("INFO", fmt.Sprintf("找到 %d 个视频文件", len(files)))
	r.app.SendProgress()
	
	// 创建处理器
	r.processor = core.NewProcessor(r.app.config)
	
	// 处理每个文件
	// {{ AURA-X: Modify - 添加详细的文件处理状态推送. Confirmed via 寸止 }}
	for i, file := range files {
		// 检查是否需要停止
		select {
		case <-r.ctx.Done():
			r.app.SendLog("WARN", "任务被用户停止")
			return
		default:
		}
		
		// 创建文件信息
		fileInfo := &FileInfo{
			Path:   file,
			Name:   filepath.Base(file),
			Size:   r.getFileSize(file),
			Status: "processing",
		}
		
		startTime := time.Now()
		
		r.app.SendLog("INFO", fmt.Sprintf("[%d/%d] 正在处理: %s", i+1, len(files), filepath.Base(file)))
		r.app.SendFileStatus(fileInfo)
		
		// 解析番号
		numberParser := parser.NewNumberParser(r.app.config)
		number := numberParser.GetNumber(filepath.Base(file))
		fileInfo.Number = number
		
		if number == "" {
			r.app.SendLog("WARN", fmt.Sprintf("无法解析番号: %s", filepath.Base(file)))
			fileInfo.Status = "skipped"
			fileInfo.Error = "无法解析番号"
			fileInfo.Duration = time.Since(startTime).Round(time.Millisecond).String()
			r.app.SendFileStatus(fileInfo)
			r.app.stats.Skipped++
			r.app.SendProgress()
			continue
		}
		
		r.app.SendLog("INFO", fmt.Sprintf("识别番号: %s (原文件名: %s)", number, filepath.Base(file)))
		fileInfo.Number = number
		r.app.SendFileStatus(fileInfo)
		
		// 处理文件
		err := r.processor.ProcessSingleFile(file, number, "", "")
		fileInfo.Duration = time.Since(startTime).Round(time.Millisecond).String()
		
		if err != nil {
			r.app.SendLog("ERROR", fmt.Sprintf("处理失败: %s -> %s - %v", filepath.Base(file), number, err))
			fileInfo.Status = "failed"
			fileInfo.Error = err.Error()
			r.app.stats.Failed++
		} else {
			r.app.SendLog("INFO", fmt.Sprintf("处理成功: %s -> %s (耗时: %s)", filepath.Base(file), number, fileInfo.Duration))
			fileInfo.Status = "success"
			r.app.stats.Success++
		}
		
		r.app.SendFileStatus(fileInfo)
		r.app.SendProgress()
		
		// 延迟（防止请求过快）
		if r.app.config.Common.Sleep > 0 {
			time.Sleep(time.Duration(r.app.config.Common.Sleep) * time.Second)
		}
	}
	
	// 任务完成
	duration := time.Since(r.app.stats.StartTime).Round(time.Second)
	r.app.SendLog("INFO", fmt.Sprintf("========== 任务完成 =========="))
	r.app.SendLog("INFO", fmt.Sprintf("总计: %d | 成功: %d | 失败: %d | 跳过: %d",
		r.app.stats.Total, r.app.stats.Success, r.app.stats.Failed, r.app.stats.Skipped))
	r.app.SendLog("INFO", fmt.Sprintf("耗时: %s", duration))
}

// getFileSize 获取文件大小
func (r *Runner) getFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return info.Size()
}

// getVideoFiles 获取视频文件列表
func (r *Runner) getVideoFiles(folder string) ([]string, error) {
	var files []string
	
	// 获取支持的文件扩展名
	mediaTypes := r.getMediaTypes()
	
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			// 跳过排除的文件夹
			if r.shouldSkipFolder(path, info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		
		// 检查是否是视频文件
		ext := strings.ToLower(filepath.Ext(path))
		for _, mediaType := range mediaTypes {
			if ext == mediaType {
				files = append(files, path)
				break
			}
		}
		
		return nil
	})
	
	return files, err
}

// getMediaTypes 获取支持的媒体文件类型
func (r *Runner) getMediaTypes() []string {
	if r.app.config == nil || r.app.config.Media.MediaType == "" {
		// 默认支持的视频格式
		return []string{".mp4", ".avi", ".mkv", ".rmvb", ".wmv", ".mov", ".flv", ".ts", ".webm", ".iso"}
	}
	
	// 从配置中解析媒体类型
	types := strings.Split(r.app.config.Media.MediaType, ",")
	result := make([]string, 0, len(types))
	for _, t := range types {
		t = strings.TrimSpace(t)
		if t != "" {
			result = append(result, strings.ToLower(t))
		}
	}
	return result
}

// shouldSkipFolder 判断是否应该跳过文件夹
func (r *Runner) shouldSkipFolder(path, name string) bool {
	// 跳过隐藏文件夹
	if len(name) > 0 && name[0] == '.' {
		return true
	}
	
	// 跳过输出文件夹和失败文件夹
	skipFolders := []string{
		r.app.config.Common.SuccessOutputFolder,
		r.app.config.Common.FailedOutputFolder,
		"failed",
		"JAV_output",
	}
	
	for _, skip := range skipFolders {
		if name == skip || filepath.Base(path) == skip {
			return true
		}
	}
	
	return false
}

// GetFileList 获取待处理文件列表
func (a *App) GetFileList() ([]FileInfo, error) {
	if a.config == nil {
		return nil, fmt.Errorf("配置未加载")
	}
	
	files, err := a.runner.getVideoFiles(a.config.Common.SourceFolder)
	if err != nil {
		return nil, err
	}
	
	fileInfos := make([]FileInfo, 0, len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		
		numberParser := parser.NewNumberParser(a.config)
		number := numberParser.GetNumber(filepath.Base(file))
		
		fileInfos = append(fileInfos, FileInfo{
			Path:   file,
			Name:   filepath.Base(file),
			Size:   info.Size(),
			Number: number,
			Status: "pending",
		})
	}
	
	return fileInfos, nil
}

