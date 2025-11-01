package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"movie-data-capture/internal/config"
	"movie-data-capture/internal/scraper"
	"movie-data-capture/pkg/logger"
)

const (
	// Windows路径最大长度限制
	windowsMaxPathLength = 260
	// Windows长路径前缀（支持最多32,767字符）
	windowsLongPathPrefix = `\\?\`
	// 安全余量（为文件名预留空间）
	pathSafetyMargin = 50
)

// Storage 处理文件操作和文件夹创建
type Storage struct {
	config *config.Config
}

// New 创建一个新的存储实例
func New(cfg *config.Config) *Storage {
	return &Storage{
		config: cfg,
	}
}

// CreateFolder 根据位置规则创建输出文件夹
func (s *Storage) CreateFolder(data *scraper.MovieData) (string, error) {
	successFolder := s.config.Common.SuccessOutputFolder
	
	// 评估位置规则
	locationRule := s.config.NameRule.LocationRule
	folderPath := s.evaluateLocationRule(locationRule, data)
	
	// 调试：打印评估后的文件夹路径
	logger.Debug("Evaluated folder path: %s", folderPath)
	
	// 处理过长的演员名称
	if strings.Contains(locationRule, "actor") && len(data.Actor) > 100 {
		// 对于多演员电影，将演员替换为"多人作品"
		folderPath = strings.ReplaceAll(folderPath, data.Actor, "多人作品")
	}
	
	// 处理过长的标题
	maxTitleLen := s.config.NameRule.MaxTitleLen
	if maxTitleLen > 0 && strings.Contains(locationRule, "title") && len(data.Title) > maxTitleLen {
		shortTitle := data.Title[:maxTitleLen]
		folderPath = strings.ReplaceAll(folderPath, data.Title, shortTitle)
	}
	
	// 确保相对路径（添加 ./ 前缀）
	if !strings.HasPrefix(folderPath, ".") && !strings.HasPrefix(folderPath, "/") {
		folderPath = "./" + folderPath
	}
	
	fullPath := filepath.Join(successFolder, folderPath)
	fullPath = filepath.Clean(fullPath)
	
	// 转义有问题的字符
	fullPath = s.escapePath(fullPath)
	
	// Windows平台：检查并处理路径长度限制
	if runtime.GOOS == "windows" {
		fullPath = s.handleWindowsLongPath(fullPath, data)
	}
	
	// 创建目录
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		// 回退：仅使用编号创建
		fallbackPath := filepath.Join(successFolder, data.Number)
		fallbackPath = s.escapePath(fallbackPath)
		
		// 对回退路径也进行长路径处理
		if runtime.GOOS == "windows" {
			fallbackPath = s.handleWindowsLongPath(fallbackPath, data)
		}
		
		err = os.MkdirAll(fallbackPath, 0755)
		if err != nil {
			return "", fmt.Errorf("创建目录失败: %w", err)
		}
		logger.Warn("Used fallback path due to original path error: %s", fallbackPath)
		return fallbackPath, nil
	}
	
	return fullPath, nil
}

// evaluateLocationRule 评估位置规则模板
func (s *Storage) evaluateLocationRule(rule string, data *scraper.MovieData) string {
	result := rule
	
	// 定义字段映射
	fields := map[string]string{
		"number":   data.Number,
		"title":    data.Title,
		"actor":    data.Actor,
		"studio":   data.Studio,
		"director": data.Director,
		"release":  data.Release,
		"year":     data.Year,
		"series":   data.Series,
		"label":    data.Label,
	}
	
	// 处理Python风格的表达式，如 "actor + '/' + number"
	// 逐步解析表达式
	parts := strings.Split(result, " + ")
	var resultParts []string
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		
		// 处理单引号中的字面字符串
		if strings.HasPrefix(part, "'") && strings.HasSuffix(part, "'") {
			// 移除引号并添加字面字符串
			literal := part[1 : len(part)-1]
			resultParts = append(resultParts, literal)
		} else {
			// 替换字段占位符
			if value, exists := fields[part]; exists {
				resultParts = append(resultParts, value)
			} else {
				// 如果不是已知字段则保持原样
				resultParts = append(resultParts, part)
			}
		}
	}
	
	// 连接所有部分，将 '/' 视为路径分隔符
	var pathComponents []string
	currentComponent := ""
	
	for _, part := range resultParts {
		if part == "/" {
				// 遇到分隔符时，将当前组件添加到路径中
			if currentComponent != "" {
				pathComponents = append(pathComponents, currentComponent)
				currentComponent = ""
			}
		} else {
			// 累积非分隔符部分
			currentComponent += part
		}
	}
	
	// 添加最后一个组件
	if currentComponent != "" {
		pathComponents = append(pathComponents, currentComponent)
	}
	
	// 使用 filepath.Join 创建具有操作系统特定分隔符的正确路径
	if len(pathComponents) > 1 {
		result = filepath.Join(pathComponents...)
	} else if len(pathComponents) == 1 {
		result = pathComponents[0]
	} else {
		result = strings.Join(resultParts, "")
	}
	
	// 移除任何剩余的空格
	result = strings.TrimSpace(result)
	
	return result
}

// escapePath 转义文件路径中的有问题字符
func (s *Storage) escapePath(path string) string {
	literals := s.config.Escape.Literals
	
	result := path
	for _, char := range literals {
		// 不转义路径分隔符
		if char == '\\' || char == '/' {
			continue
		}
		result = strings.ReplaceAll(result, string(char), "")
	}
	
	return result
}

// sanitizeFileName 清理文件名中的非法字符（保持最大兼容性）
// Source: AURA-X Protocol - 确保文件名在所有文件系统上都有效
func (s *Storage) sanitizeFileName(fileName string) string {
	if fileName == "" {
		return fileName
	}
	
	// Windows文件系统禁止的字符: < > : " / \ | ? *
	// 还包括控制字符（0-31）和某些特殊字符
	// {{ AURA-X: Modify - 修复 map 值类型错误，将 rune 改为 string. Approval: 寸止(ID:20251101). }}
	illegalChars := map[rune]string{
		'<':  "＜", // 全角替换
		'>':  "＞",
		':':  "꞉", // 修饰符冒号
		'"':  "＂", // 全角引号
		'/':  "∕", // 除号斜杠
		'\\': "∖", // 集合减号
		'|':  "ǀ", // 齿音咔嗒
		'?':  "？", // 全角问号
		'*':  "∗", // 星号运算符
	}
	
	result := ""
	replaced := false
	
	for _, char := range fileName {
		// 检查是否是非法字符
		if replacement, isIllegal := illegalChars[char]; isIllegal {
			result += replacement
			replaced = true
		} else if char < 32 {
			// 跳过控制字符
			replaced = true
		} else {
			result += string(char)
		}
	}
	
	// 移除文件名末尾的点和空格（Windows限制）
	result = strings.TrimRight(result, ". ")
	
	// 如果文件名为空，使用默认名称
	if result == "" {
		result = "unnamed_file"
	}
	
	if replaced {
		logger.Debug("Sanitized filename: '%s' -> '%s'", fileName, result)
	}
	
	return result
}

// MoveFile 移动或链接文件到目标位置
func (s *Storage) MoveFile(sourcePath, destPath string) error {
	// Source: AURA-X Protocol - 清理目标文件名
	destDir := filepath.Dir(destPath)
	destFileName := filepath.Base(destPath)
	cleanDestFileName := s.sanitizeFileName(destFileName)
	cleanDestPath := filepath.Join(destDir, cleanDestFileName)
	
	// 检查目标文件是否已存在
	if _, err := os.Stat(cleanDestPath); err == nil {
		return fmt.Errorf("destination file already exists: %s", cleanDestPath)
	}
	
	// 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	
	linkMode := s.config.Common.LinkMode
	
	// 使用清理后的路径
	actualDestPath := cleanDestPath
	
	switch linkMode {
	case 0:
		// 移动文件
		return s.moveFile(sourcePath, actualDestPath)
	case 1:
		// 创建软链接
		return s.createSoftLink(sourcePath, actualDestPath)
	case 2:
		// 首先尝试硬链接，失败则回退到软链接
		err := s.createHardLink(sourcePath, actualDestPath)
		if err != nil {
			logger.Debug("Hard link failed, trying soft link: %v", err)
			return s.createSoftLink(sourcePath, actualDestPath)
		}
		return nil
	default:
		return s.moveFile(sourcePath, actualDestPath)
	}
}

// moveFile 将文件从源位置移动到目标位置
func (s *Storage) moveFile(sourcePath, destPath string) error {
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		// 如果重命名失败，尝试复制并删除
		return s.copyAndDelete(sourcePath, destPath)
	}
	
	logger.Info("Moved file: %s -> %s", sourcePath, destPath)
	return nil
}

// createSoftLink 创建符号链接
func (s *Storage) createSoftLink(sourcePath, destPath string) error {
	// 首先尝试相对路径
	destDir := filepath.Dir(destPath)
	relPath, err := filepath.Rel(destDir, sourcePath)
	if err == nil {
		err = os.Symlink(relPath, destPath)
		if err == nil {
			logger.Info("Created soft link: %s -> %s", destPath, relPath)
			return nil
		}
	}
	
	// 回退到绝对路径
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	
	err = os.Symlink(absPath, destPath)
	if err != nil {
		return fmt.Errorf("failed to create soft link: %w", err)
	}
	
	logger.Info("Created soft link: %s -> %s", destPath, absPath)
	return nil
}

// createHardLink 创建硬链接
func (s *Storage) createHardLink(sourcePath, destPath string) error {
	err := os.Link(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("failed to create hard link: %w", err)
	}
	
	logger.Info("Created hard link: %s -> %s", destPath, sourcePath)
	return nil
}

// copyAndDelete 复制文件并删除源文件
func (s *Storage) copyAndDelete(sourcePath, destPath string) error {
	// 打开源文件
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()
	
	// 创建目标文件
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()
	
	// 复制数据
	_, err = destFile.ReadFrom(srcFile)
	if err != nil {
		// 移除部分复制的文件
		os.Remove(destPath)
		return fmt.Errorf("failed to copy file: %w", err)
	}
	
	// 复制文件权限
	srcInfo, err := srcFile.Stat()
	if err == nil {
		os.Chmod(destPath, srcInfo.Mode())
	}
	
	// 删除源文件
	err = os.Remove(sourcePath)
	if err != nil {
		logger.Warn("Failed to delete source file %s: %v", sourcePath, err)
	}
	
	logger.Info("Copied and deleted: %s -> %s", sourcePath, destPath)
	return nil
}

// MoveToFailedFolder 将文件移动到失败文件夹
func (s *Storage) MoveToFailedFolder(filePath string) error {
	failedFolder := s.config.Common.FailedOutputFolder
	
	// 如果失败文件夹不存在则创建
	if err := os.MkdirAll(failedFolder, 0755); err != nil {
		return fmt.Errorf("failed to create failed folder: %w", err)
	}
	
	mainMode := s.config.Common.MainMode
	linkMode := s.config.Common.LinkMode
	
	// 模式3或链接模式：添加到失败列表而不是移动
	if mainMode == 3 || linkMode > 0 {
		return s.addToFailedList(filePath, failedFolder)
	}
	
	// 移动模式：如果配置了则实际移动文件
	if s.config.Common.FailedMove {
		return s.moveToFailedFolder(filePath, failedFolder)
	}
	
	return nil
}

// addToFailedList 将文件路径添加到失败列表
func (s *Storage) addToFailedList(filePath, failedFolder string) error {
	failedListPath := filepath.Join(failedFolder, "failed_list.txt")
	
	file, err := os.OpenFile(failedListPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open failed list: %w", err)
	}
	defer file.Close()
	
	_, err = file.WriteString(filePath + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to failed list: %w", err)
	}
	
	logger.Info("Added to failed list: %s", filePath)
	return nil
}

// moveToFailedFolder 将文件移动到失败文件夹
func (s *Storage) moveToFailedFolder(filePath, failedFolder string) error {
	fileName := filepath.Base(filePath)
	// Source: AURA-X Protocol - 清理文件名确保兼容性
	cleanFileName := s.sanitizeFileName(fileName)
	destPath := filepath.Join(failedFolder, cleanFileName)
	
	// Source: AURA-X Protocol - 增强错误处理，避免"找不到文件"错误
	
	// 首先检查源文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			logger.Warn("Source file no longer exists, skipping move: %s", filePath)
			return nil // 文件不存在，不算错误，直接返回
		}
		return fmt.Errorf("failed to check source file: %w", err)
	}
	
	// 检查目标是否存在
	if _, err := os.Stat(destPath); err == nil {
		logger.Warn("File already exists in failed folder: %s", fileName)
		// 源文件存在但目标已存在，删除源文件避免冲突
		if err := os.Remove(filePath); err != nil {
			logger.Warn("Failed to remove duplicate source file: %v", err)
		}
		return nil
	}
	
	// 确保失败文件夹存在（二次检查，以防并发问题）
	if err := os.MkdirAll(failedFolder, 0755); err != nil {
		return fmt.Errorf("failed to ensure failed folder exists: %w", err)
	}
	
	// 记录移动操作
	recordPath := filepath.Join(failedFolder, "where_was_i_before_being_moved.txt")
	file, err := os.OpenFile(recordPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		timestamp := time.Now().Format("2006-01-02 15:04")
		record := fmt.Sprintf("%s FROM[%s]TO[%s]\n", timestamp, filePath, destPath)
		file.WriteString(record)
		file.Close()
	} else {
		logger.Warn("Failed to write move record: %v", err)
	}
	
	// 移动文件
	err = os.Rename(filePath, destPath)
	if err != nil {
		// 如果是跨驱动器移动导致的错误，尝试复制后删除
		if strings.Contains(err.Error(), "cross-device") || 
		   strings.Contains(err.Error(), "different") {
			logger.Debug("Cross-device move detected, using copy method")
			if copyErr := s.copyAndRemove(filePath, destPath); copyErr != nil {
				return fmt.Errorf("failed to move file (copy method): %w", copyErr)
			}
		} else {
			return fmt.Errorf("failed to move file to failed folder: %w", err)
		}
	}
	
	logger.Info("Moved to failed folder: %s", fileName)
	return nil
}

// copyAndRemove 复制文件后删除源文件（用于跨驱动器移动）
func (s *Storage) copyAndRemove(src, dst string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()
	
	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()
	
	// 复制内容
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		// 复制失败，删除不完整的目标文件
		os.Remove(dst)
		return fmt.Errorf("failed to copy file content: %w", err)
	}
	
	// 确保内容写入磁盘
	err = destFile.Sync()
	if err != nil {
		logger.Warn("Failed to sync destination file: %v", err)
	}
	
	// 复制成功，删除源文件
	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("failed to remove source file after copy: %w", err)
	}
	
	logger.Debug("Successfully copied and removed: %s -> %s", src, dst)
	return nil
}

// RemoveEmptyFolders 移除空目录
func (s *Storage) RemoveEmptyFolders(rootPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 出错时继续
		}
		
		if !info.IsDir() {
			return nil
		}
		
		// 不要移除根路径本身
		if path == rootPath {
			return nil
		}
		
		// 检查目录是否为空
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil
		}
		
		if len(entries) == 0 {
			err = os.Remove(path)
			if err == nil {
				logger.Info("Removed empty folder: %s", path)
			}
		}
		
		return nil
	})
}

// FindSubtitleFiles 查找与视频文件匹配的字幕文件
// 返回找到的所有字幕文件路径列表
func (s *Storage) FindSubtitleFiles(videoFilePath string) []string {
	var subtitleFiles []string
	
	// 获取视频文件的目录和基础名称（无扩展名）
	videoDir := filepath.Dir(videoFilePath)
	videoBase := strings.TrimSuffix(filepath.Base(videoFilePath), filepath.Ext(videoFilePath))
	
	// 支持的字幕文件扩展名列表
	subtitleExts := strings.Split(s.config.Media.SubType, ",")
	
	// 查找目录中的所有文件
	entries, err := os.ReadDir(videoDir)
	if err != nil {
		logger.Warn("Failed to read directory for subtitle search: %v", err)
		return subtitleFiles
	}
	
	// 遍历文件，查找匹配的字幕文件
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		fileName := entry.Name()
		fileExt := filepath.Ext(fileName)
		fileBase := strings.TrimSuffix(fileName, fileExt)
		
		// 检查扩展名是否为字幕格式
		isSubtitle := false
		for _, subtitleExt := range subtitleExts {
			subtitleExt = strings.TrimSpace(subtitleExt)
			if strings.EqualFold(fileExt, subtitleExt) {
				isSubtitle = true
				break
			}
		}
		
		if !isSubtitle {
			continue
		}
		
		// 检查文件名是否匹配（支持多种匹配模式）
		// 模式1: 完全匹配 video.srt
		// 模式2: 语言后缀 video.zh.srt, video.chs.srt, video.eng.srt
		// 模式3: 强制/外挂标识 video.forced.srt, video.default.srt
		if strings.EqualFold(fileBase, videoBase) ||
			strings.HasPrefix(strings.ToLower(fileBase), strings.ToLower(videoBase)+".") ||
			strings.HasPrefix(strings.ToLower(fileBase), strings.ToLower(videoBase)+"_") {
			
			subtitlePath := filepath.Join(videoDir, fileName)
			subtitleFiles = append(subtitleFiles, subtitlePath)
			logger.Debug("Found subtitle file: %s", fileName)
		}
	}
	
	return subtitleFiles
}

// MoveSubtitleFiles 移动字幕文件到目标目录
// videoFileName: 目标视频文件名（用于重命名字幕文件）
// destDir: 目标目录
func (s *Storage) MoveSubtitleFiles(subtitleFiles []string, videoFileName, destDir string) error {
	if len(subtitleFiles) == 0 {
		return nil
	}
	
	videoBase := strings.TrimSuffix(videoFileName, filepath.Ext(videoFileName))
	
	for _, subtitlePath := range subtitleFiles {
		// 获取字幕文件信息
		subtitleName := filepath.Base(subtitlePath)
		subtitleExt := filepath.Ext(subtitleName)
		originalBase := strings.TrimSuffix(subtitleName, subtitleExt)
		
		// 提取语言/类型后缀（如 .zh, .chs, .eng, .forced 等）
		// 示例: movie.zh.srt -> .zh
		//      movie.chs.srt -> .chs
		//      movie.forced.srt -> .forced
		suffix := ""
		if idx := strings.Index(originalBase, "."); idx != -1 {
			suffix = originalBase[idx:] // 保留所有点号后的部分
		} else if idx := strings.Index(originalBase, "_"); idx != -1 {
			// 处理下划线分隔的情况
			suffix = "." + originalBase[idx+1:]
		}
		
		// 生成新的字幕文件名
		// 如果有语言后缀，保留它；否则直接使用视频基础名
		var newSubtitleName string
		if suffix != "" && suffix != "."+strings.ToLower(videoBase) {
			newSubtitleName = videoBase + suffix + subtitleExt
		} else {
			newSubtitleName = videoBase + subtitleExt
		}
		
		destPath := filepath.Join(destDir, newSubtitleName)
		
		// 检查目标文件是否已存在
		if _, err := os.Stat(destPath); err == nil {
			logger.Debug("Subtitle file already exists at destination: %s", newSubtitleName)
			continue
		}
		
		// 移动字幕文件（使用与视频文件相同的link_mode）
		err := s.MoveFile(subtitlePath, destPath)
		if err != nil {
			logger.Warn("Failed to move subtitle file %s: %v", subtitleName, err)
			// 继续处理其他字幕文件，不中断
			continue
		}
		
		logger.Info("Moved subtitle file: %s -> %s", subtitleName, newSubtitleName)
	}
	
	return nil
}

// handleWindowsLongPath 处理Windows平台的长路径问题
// 如果路径超过限制，采用智能策略缩短路径或添加长路径前缀
func (s *Storage) handleWindowsLongPath(fullPath string, data *scraper.MovieData) string {
	// 转换为绝对路径以准确计算长度
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		logger.Warn("Failed to get absolute path for %s: %v", fullPath, err)
		absPath = fullPath
	}
	
	// 检查路径长度（考虑安全余量）
	pathLength := len(absPath)
	if pathLength < windowsMaxPathLength-pathSafetyMargin {
		return fullPath // 路径长度安全，无需处理
	}
	
	logger.Warn("Path length (%d chars) exceeds safe limit, attempting to shorten: %s", pathLength, absPath)
	
	// 策略1: 智能缩短路径组件
	shortenedPath := s.shortenPathComponents(fullPath, data)
	
	// 再次检查缩短后的路径
	shortenedAbsPath, err := filepath.Abs(shortenedPath)
	if err == nil && len(shortenedAbsPath) < windowsMaxPathLength-pathSafetyMargin {
		logger.Info("Successfully shortened path to %d chars", len(shortenedAbsPath))
		return shortenedPath
	}
	
	// 策略2: 如果缩短失败，尝试使用长路径前缀（\\?\）
	// 注意：长路径前缀需要绝对路径
	if !strings.HasPrefix(absPath, windowsLongPathPrefix) {
		longPath := windowsLongPathPrefix + absPath
		logger.Info("Applied Windows long path prefix: %s", longPath)
		return longPath
	}
	
	// 如果已经有前缀，返回原路径
	return fullPath
}

// shortenPathComponents 智能缩短路径中的各个组件
func (s *Storage) shortenPathComponents(fullPath string, data *scraper.MovieData) string {
	// 分解路径为组件
	pathParts := strings.Split(filepath.ToSlash(fullPath), "/")
	
	// 保留基础路径（输出文件夹），只缩短后面的组件
	if len(pathParts) <= 2 {
		return fullPath // 路径太短，无法进一步缩短
	}
	
	// 找到需要缩短的部分（通常是演员名和标题）
	for i := len(pathParts) - 1; i >= 0; i-- {
		part := pathParts[i]
		
		// 跳过短组件
		if len(part) <= 20 {
			continue
		}
		
		// 缩短长组件
		shortened := s.shortenString(part, 30)
		if shortened != part {
			pathParts[i] = shortened
			logger.Debug("Shortened path component: %s -> %s", part, shortened)
		}
	}
	
	// 重新组合路径
	result := strings.Join(pathParts, string(filepath.Separator))
	return result
}

// shortenString 智能缩短字符串（保留重要信息）
// 尝试在单词边界处截断，保留开头和重要部分
func (s *Storage) shortenString(str string, maxLen int) string {
	// 如果字符串已经够短，直接返回
	if utf8.RuneCountInString(str) <= maxLen {
		return str
	}
	
	// 将字符串转换为rune数组以正确处理多字节字符（中文、日文等）
	runes := []rune(str)
	if len(runes) <= maxLen {
		return str
	}
	
	// 策略：保留前80%的maxLen长度
	keepLen := (maxLen * 8) / 10
	if keepLen < 10 {
		keepLen = maxLen - 3 // 至少保留一些字符
	}
	
	// 截断并添加省略号
	if keepLen > len(runes) {
		keepLen = len(runes)
	}
	
	shortened := string(runes[:keepLen])
	
	// 尝试在单词边界（空格、标点）处截断
	lastSpace := strings.LastIndexAny(shortened, " -_.,;")
	if lastSpace > maxLen/2 { // 确保不会截得太短
		shortened = shortened[:lastSpace]
	}
	
	return strings.TrimSpace(shortened)
}

// getPathLength 获取路径的实际长度（字节数）
func (s *Storage) getPathLength(path string) int {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return len(path)
	}
	return len(absPath)
}

// isPathTooLong 检查路径是否超过Windows限制
func (s *Storage) isPathTooLong(path string) bool {
	if runtime.GOOS != "windows" {
		return false // 非Windows平台不受此限制
	}
	
	length := s.getPathLength(path)
	// 考虑安全余量（为文件名留空间）
	return length >= windowsMaxPathLength-pathSafetyMargin
}