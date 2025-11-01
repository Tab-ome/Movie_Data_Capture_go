# Build Assets

此目录用于存放构建应用所需的图标和资源文件。

## 当前状态

⚠️ **图标文件已被移除**：为避免构建错误，占位符图标已被删除。应用将使用Wails默认图标。

## 添加自定义图标（可选）

如需自定义应用图标，请按以下步骤操作：

### 步骤1：准备图标文件

准备一个 **512x512 像素或更大**的高质量PNG图标

### 步骤2：生成所需格式

**方法A - 使用在线工具**：
1. 访问 https://convertio.co/zh/png-ico/
2. 上传你的PNG图标
3. 转换并下载ICO文件

**方法B - 使用ImageMagick**（如已安装）：
```bash
# 生成多尺寸ICO文件
convert appicon.png -define icon:auto-resize=256,128,64,48,32,16 icon.ico
```

### 步骤3：放置文件

```
build/
  ├── appicon.png          # 放置PNG图标
  └── windows/
      └── icon.ico         # 放置Windows ICO图标
```

### 步骤4：重新构建

运行 `wails build -tags gui` 或 `wails dev -tags gui`

---

## 注意事项

- 图标文件必须是**有效格式**，否则会导致构建失败
- ICO文件至少需要几KB大小（包含多种分辨率）
- 如果不需要自定义图标，可以跳过此步骤

