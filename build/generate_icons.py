#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
{{ AURA-X: Create - 图标格式转换脚本 }}
用途：将SVG图标转换为PNG和ICO格式
依赖：pip install cairosvg pillow
"""

import os
from pathlib import Path

try:
    import cairosvg
    from PIL import Image
    HAS_DEPS = True
except ImportError:
    HAS_DEPS = False
    print("⚠️  缺少依赖库，请先安装：")
    print("   pip install cairosvg pillow")
    print()

def svg_to_png(svg_path, png_path, size):
    """将SVG转换为指定尺寸的PNG"""
    cairosvg.svg2png(
        url=svg_path,
        write_to=png_path,
        output_width=size,
        output_height=size
    )
    print(f"✅ 已生成：{png_path} ({size}x{size})")

def create_ico(png_path, ico_path):
    """创建Windows ICO文件（包含多尺寸）"""
    img = Image.open(png_path)
    
    # ICO需要的尺寸
    sizes = [(256, 256), (128, 128), (64, 64), (48, 48), (32, 32), (16, 16)]
    
    # 生成不同尺寸的图像
    icons = []
    for size in sizes:
        resized = img.resize(size, Image.Resampling.LANCZOS)
        icons.append(resized)
    
    # 保存为ICO
    icons[0].save(
        ico_path,
        format='ICO',
        sizes=sizes,
        append_images=icons[1:]
    )
    print(f"✅ 已生成：{ico_path} (包含 {len(sizes)} 个尺寸)")

def main():
    if not HAS_DEPS:
        return
    
    # 获取脚本所在目录
    script_dir = Path(__file__).parent
    svg_file = script_dir / "appicon.svg"
    
    if not svg_file.exists():
        print(f"❌ 错误：找不到 {svg_file}")
        return
    
    print("🎨 开始生成图标文件...\n")
    
    # 1. 生成主PNG文件 (512x512)
    png_512 = script_dir / "appicon.png"
    svg_to_png(str(svg_file), str(png_512), 512)
    
    # 2. 生成其他尺寸的PNG
    sizes = [256, 128, 64]
    for size in sizes:
        png_file = script_dir / f"icon_{size}.png"
        svg_to_png(str(svg_file), str(png_file), size)
    
    # 3. 生成Windows ICO
    ico_dir = script_dir / "windows"
    ico_dir.mkdir(exist_ok=True)
    ico_file = ico_dir / "icon.ico"
    
    print("\n🔧 生成 Windows ICO 文件...")
    create_ico(str(png_512), str(ico_file))
    
    # 4. 备份原图标（如果存在且未备份）
    original_png = script_dir / "appicon.png.backup"
    if png_512.exists() and not original_png.exists():
        import shutil
        shutil.copy2(str(png_512), str(original_png))
        print(f"💾 已备份原图标：{original_png}")
    
    print("\n✅ 图标生成完成！")
    print("\n📁 生成的文件：")
    print(f"   - {png_512}")
    print(f"   - {script_dir / 'icon_256.png'}")
    print(f"   - {script_dir / 'icon_128.png'}")
    print(f"   - {script_dir / 'icon_64.png'}")
    print(f"   - {ico_file}")
    print("\n🚀 现在可以重新编译应用查看新图标！")

if __name__ == "__main__":
    main()

