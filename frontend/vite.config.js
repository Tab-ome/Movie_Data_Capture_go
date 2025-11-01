// {{ AURA-X: Add - Vite 配置文件，适配 Wails 开发环境 }}

import { defineConfig } from 'vite';

export default defineConfig({
  // 服务器配置
  server: {
    port: 34115, // Wails 默认端口
    strictPort: true, // 端口被占用时不自动尝试下一个
    host: '127.0.0.1',
    cors: true,
  },
  
  // 构建配置
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    rollupOptions: {
      input: './src/index.html',
    },
  },
  
  // 开发模式下的根目录
  root: './src',
  
  // 基础路径
  base: './',
});

