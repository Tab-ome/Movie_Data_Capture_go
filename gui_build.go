//go:build gui
// +build gui

package main

import (
	"embed"
	"log"

	"movie-data-capture/gui"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// {{ AURA-X: Add - 标识这是GUI构建版本 }}
const isGUIBuild = true

// {{ AURA-X: Add - GUI模式实现(使用构建标签). Source: context7-mcp on 'Wails v2'. Confirmed via 寸止 }}
func runGUI() {
	// 创建应用实例
	app := gui.NewApp()

	// 创建Wails应用配置
	err := wails.Run(&options.App{
		Title:  "Movie Data Capture",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnDomReady:       app.DomReady,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		log.Fatal("启动GUI失败:", err)
	}
}

