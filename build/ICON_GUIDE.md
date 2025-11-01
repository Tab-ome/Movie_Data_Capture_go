# 🎨 Movie Data Capture 图标更换指南

## 📋 概述

本指南说明如何为 Movie Data Capture 应用生成和应用新的自定义图标。

## 🎯 设计说明

**图标主题**：电影胶片 + 数据融合
- **主体元素**：电影胶片卷轴造型
- **数据元素**：胶片画面区域内有金色数据点和扫描线
- **色彩方案**：
  - 主色：科技蓝 (#5B7FFF) → 紫色 (#9D6FFF) 渐变
  - 点缀：金色 (#FFB84D) 数据点
  - 背景：深色渐变 (#1a1a2e → #16213e)
- **底部标识**：MDC 字母标识

## 📁 文件说明

### 源文件
- `appicon.svg` - 矢量图标源文件（512x512，可编辑）

### 生成脚本
- `generate_icons.py` - Python 图标生成脚本
- `generate_icons.bat` - Windows 一键生成脚本

### 输出文件（脚本生成）
- `appicon.png` - 主图标 (512x512)
- `icon_256.png` - 中等尺寸 (256x256)
- `icon_128.png` - 小尺寸 (128x128)
- `icon_64.png` - 超小尺寸 (64x64)
- `windows/icon.ico` - Windows 图标（包含多尺寸）

## 🚀 快速开始

### 方法1：使用批处理脚本（推荐 - Windows用户）

1. 双击运行 `generate_icons.bat`
2. 脚本会自动：
   - 检查 Python 环境
   - 安装所需依赖（cairosvg, pillow）
   - 从 SVG 生成所有所需格式的图标
   - 备份原图标文件

### 方法2：手动执行Python脚本

```bash
# 1. 安装依赖
pip install cairosvg pillow

# 2. 运行脚本
cd build
python generate_icons.py
```

### 方法3：在线转换（无需编程环境）

如果不想安装Python，可使用在线工具：

1. **SVG → PNG 转换**
   - 访问：https://cloudconvert.com/svg-to-png
   - 上传 `appicon.svg`
   - 设置输出尺寸为 512x512
   - 下载并重命名为 `appicon.png`

2. **PNG → ICO 转换**
   - 访问：https://www.icoconverter.com/
   - 上传生成的 512x512 PNG
   - 选择生成多尺寸 ICO
   - 下载并放到 `windows/icon.ico`

## 🔧 Wails 图标配置

Wails 框架会自动使用以下图标文件：

### Windows 平台
- `build/appicon.png` - 应用程序图标
- `build/windows/icon.ico` - Windows 可执行文件图标

### macOS 平台（如需支持）
- `build/appicon.png` - 自动转换为 .icns

### Linux 平台
- `build/appicon.png` - 应用图标

**注意**：Wails 在构建时会自动查找这些文件，无需在 `wails.json` 中额外配置。

## 📝 自定义图标

如果需要修改图标设计：

1. **编辑 SVG 文件**
   - 使用 Inkscape、Adobe Illustrator 或任何 SVG 编辑器
   - 打开 `build/appicon.svg`
   - 修改颜色、形状或元素
   - 保存后重新运行生成脚本

2. **替换为自己的设计**
   - 准备一个 512x512 的 PNG 图标
   - 替换 `build/appicon.png`
   - 使用在线工具或 Python 脚本生成 ICO 文件

## 🏗️ 重新编译应用

图标生成完成后，需要重新编译应用：

```bash
# Windows
build-gui.bat

# Linux/macOS
./build-gui.sh
```

编译完成后，新图标会应用到可执行文件。

## ⚠️ 常见问题

### Q1：运行脚本时提示"找不到Python"
**解决方案**：
- 从 https://www.python.org/downloads/ 下载安装 Python 3.7+
- 安装时勾选 "Add Python to PATH"

### Q2：依赖安装失败
**解决方案**：
```bash
# 使用国内镜像加速
pip install -i https://pypi.tuna.tsinghua.edu.cn/simple cairosvg pillow
```

### Q3：生成的 ICO 文件无法使用
**解决方案**：
- 确保源 PNG 图标是正方形（512x512）
- 使用在线工具 https://www.icoconverter.com/ 重新生成

### Q4：编译后图标没有更新
**解决方案**：
- 清理构建缓存：删除 `build/bin` 目录
- 重新运行编译脚本
- Windows 可能需要重启文件资源管理器

## 📐 技术规格

### PNG 图标要求
- 尺寸：512x512 像素（推荐）
- 格式：PNG-24 带透明通道
- 色彩：RGB 模式
- 背景：建议使用圆角矩形背景（不透明）

### ICO 图标包含尺寸
- 256x256
- 128x128
- 64x64
- 48x48
- 32x32
- 16x16

## 🎨 设计建议

1. **可识别性**：图标在 16x16 小尺寸下仍能识别
2. **对比度**：保证在深色和浅色背景下都清晰可见
3. **圆角背景**：现代应用通常使用圆角矩形背景
4. **避免细节过多**：小尺寸时细节会模糊
5. **品牌一致性**：与应用主题和配色保持一致

## 📞 技术支持

如有问题，请：
1. 检查本指南的"常见问题"部分
2. 查看项目 README.md
3. 在 GitHub Issues 提交问题

---

**{{ AURA-X: Create - 图标更换完整指南文档 }}**

*最后更新：2025-11-01*

