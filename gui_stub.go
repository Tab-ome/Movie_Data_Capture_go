//go:build !gui
// +build !gui

package main

import "log"

// {{ AURA-X: Add - 标识这是CLI构建版本 }}
const isGUIBuild = false

// {{ AURA-X: Add - GUI模式存根(非GUI构建). Confirmed via 寸止 }}
// runGUI 在非GUI构建中的存根实现
func runGUI() {
	log.Fatal("此程序未包含GUI支持。请使用 'wails dev' 或 'wails build' 编译GUI版本。")
}

