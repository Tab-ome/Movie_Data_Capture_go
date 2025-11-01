package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"movie-data-capture/internal/config"
	"movie-data-capture/pkg/logger"
)

// MetaTubeAdapter 适配器，将MetaTube API响应转换为MovieData
type MetaTubeAdapter struct {
	config     *config.Config
	httpClient *http.Client
	baseURL    string
	token      string
}

// MetaTubeMovieResponse MetaTube电影API响应结构
type MetaTubeMovieResponse struct {
	Data  *MetaTubeMovie `json:"data"`
	Error interface{}    `json:"error"`
}

// MetaTubeMovie MetaTube电影数据结构
type MetaTubeMovie struct {
	Provider    string                 `json:"provider"`
	ID          string                 `json:"id"`
	Number      string                 `json:"number"`
	Title       string                 `json:"title"`
	Actors      []MetaTubeActor        `json:"actors"`
	ReleaseDate string                 `json:"releaseDate"`
	Runtime     int                    `json:"runtime"`
	Director    string                 `json:"director"`
	Studio      string                 `json:"studio"`
	Label       string                 `json:"label"`
	Series      string                 `json:"series"`
	Tags        []string               `json:"tags"`
	Summary     string                 `json:"summary"`
	Cover       string                 `json:"cover"`
	Images      []string               `json:"images"`
	Trailer     string                 `json:"trailer"`
	Score       float64                `json:"score"`
	Homepage    string                 `json:"homepage"`
	Uncensored  bool                   `json:"uncensored"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// MetaTubeActor MetaTube演员数据结构
type MetaTubeActor struct {
	Name   string   `json:"name"`
	Images []string `json:"images"`
}

// MetaTubeSearchResult MetaTube搜索结果
type MetaTubeSearchResult struct {
	Data  []MetaTubeSearchItem `json:"data"`
	Error interface{}          `json:"error"`
}

// MetaTubeSearchItem MetaTube搜索项
type MetaTubeSearchItem struct {
	Provider string  `json:"provider"`
	ID       string  `json:"id"`
	Number   string  `json:"number"`
	Title    string  `json:"title"`
	Score    float64 `json:"score"`
}

// NewMetaTubeAdapter 创建新的MetaTube适配器
func NewMetaTubeAdapter(cfg *config.Config) *MetaTubeAdapter {
	return &MetaTubeAdapter{
		config:  cfg,
		baseURL: cfg.Scraper.MetaTubeURL,
		token:   cfg.Scraper.MetaTubeToken,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Proxy.Timeout) * time.Second,
		},
	}
}

// ScrapeByNumber 通过番号从MetaTube API抓取数据
func (m *MetaTubeAdapter) ScrapeByNumber(ctx context.Context, number string) (*MovieData, error) {
	// 首先搜索电影
	searchResults, err := m.searchMovie(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	if len(searchResults) == 0 {
		return nil, fmt.Errorf("no results found for: %s", number)
	}

	// 使用第一个结果（评分最高的）
	result := searchResults[0]
	logger.Debug("Found movie on MetaTube: provider=%s, id=%s, title=%s", 
		result.Provider, result.ID, result.Title)

	// 获取详细信息
	movie, err := m.getMovieInfo(ctx, result.Provider, result.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie info: %w", err)
	}

	// 转换为MovieData
	return m.convertToMovieData(movie), nil
}

// searchMovie 搜索电影
func (m *MetaTubeAdapter) searchMovie(ctx context.Context, query string) ([]MetaTubeSearchItem, error) {
	// 构建搜索URL
	apiURL := fmt.Sprintf("%s/v1/movies/search?q=%s", 
		strings.TrimRight(m.baseURL, "/"), 
		url.QueryEscape(query))

	logger.Debug("Searching MetaTube API: %s", apiURL)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// 添加认证令牌（如果有）
	if m.token != "" {
		req.Header.Set("Authorization", "Bearer "+m.token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var searchResult MetaTubeSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if searchResult.Error != nil {
		return nil, fmt.Errorf("API error: %v", searchResult.Error)
	}

	return searchResult.Data, nil
}

// getMovieInfo 获取电影详细信息
func (m *MetaTubeAdapter) getMovieInfo(ctx context.Context, provider, id string) (*MetaTubeMovie, error) {
	apiURL := fmt.Sprintf("%s/v1/movies/%s/%s?lazy=false", 
		strings.TrimRight(m.baseURL, "/"), 
		url.PathEscape(provider),
		url.PathEscape(id))

	logger.Debug("Getting movie info from MetaTube: %s", apiURL)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// 添加认证令牌（如果有）
	if m.token != "" {
		req.Header.Set("Authorization", "Bearer "+m.token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var movieResp MetaTubeMovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&movieResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if movieResp.Error != nil {
		return nil, fmt.Errorf("API error: %v", movieResp.Error)
	}

	if movieResp.Data == nil {
		return nil, fmt.Errorf("no movie data returned")
	}

	return movieResp.Data, nil
}

// convertToMovieData 将MetaTube数据转换为MovieData
func (m *MetaTubeAdapter) convertToMovieData(movie *MetaTubeMovie) *MovieData {
	data := &MovieData{
		Number:        movie.Number,
		Title:         movie.Title,
		OriginalTitle: movie.Title,
		Release:       m.formatDate(movie.ReleaseDate),
		Runtime:       m.formatRuntime(movie.Runtime),
		Director:      movie.Director,
		Studio:        movie.Studio,
		Label:         movie.Label,
		Series:        movie.Series,
		Tag:           movie.Tags,
		Outline:       movie.Summary,
		Cover:         movie.Cover,
		Trailer:       movie.Trailer,
		Website:       movie.Homepage,
		Source:        fmt.Sprintf("MetaTube(%s)", movie.Provider),
		Uncensored:    movie.Uncensored,
		UserRating:    movie.Score / 10.0, // MetaTube使用0-100，转换为0-10
		ActorList:     []string{},
		ActorPhoto:    make(map[string]string),
	}

	// 处理演员信息
	for _, actor := range movie.Actors {
		data.ActorList = append(data.ActorList, actor.Name)
		// 使用第一张图片作为演员照片
		if len(actor.Images) > 0 {
			data.ActorPhoto[actor.Name] = actor.Images[0]
		}
	}

	// 合并演员名称
	data.Actor = strings.Join(data.ActorList, ", ")

	// 处理额外fanart图片
	if len(movie.Images) > 0 {
		// 第一张通常是封面，其余的作为extrafanart
		if len(movie.Images) > 1 {
			data.Extrafanart = movie.Images[1:]
		}
		// 如果没有封面，使用第一张图片
		if data.Cover == "" && len(movie.Images) > 0 {
			data.Cover = movie.Images[0]
		}
	}

	// 解析年份
	if data.Release != "" && len(data.Release) >= 4 {
		data.Year = data.Release[:4]
	}

	return data
}

// formatDate 格式化日期
func (m *MetaTubeAdapter) formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// 尝试解析多种日期格式
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}

	// 如果解析失败，返回原始字符串
	return dateStr
}

// formatRuntime 格式化运行时间
func (m *MetaTubeAdapter) formatRuntime(minutes int) string {
	if minutes <= 0 {
		return ""
	}
	return fmt.Sprintf("%d", minutes)
}

// HealthCheck 检查MetaTube API服务器是否可访问
func (m *MetaTubeAdapter) HealthCheck(ctx context.Context) error {
	apiURL := fmt.Sprintf("%s/", strings.TrimRight(m.baseURL, "/"))
	
	logger.Debug("Checking MetaTube API health: %s", apiURL)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if m.token != "" {
		req.Header.Set("Authorization", "Bearer "+m.token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("MetaTube API unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("MetaTube API server error: status %d", resp.StatusCode)
	}

	logger.Info("MetaTube API health check passed")
	return nil
}

// Close 关闭适配器
func (m *MetaTubeAdapter) Close() error {
	// HTTP client不需要显式关闭
	return nil
}

