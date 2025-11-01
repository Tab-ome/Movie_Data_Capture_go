# 配置项映射表

用于GUI开发的配置项完整列表和说明

## 1. Common（基础配置）- 19项

| 字段 | 类型 | 说明 | 默认值 | 选项/范围 |
|------|------|------|--------|-----------|
| main_mode | int | 运行模式 | 1 | 1=刮削,2=整理,3=分析 |
| source_folder | string | 源文件夹路径 | "./" | 文件夹路径 |
| failed_output_folder | string | 失败文件输出文件夹 | "failed" | 文件夹路径 |
| success_output_folder | string | 成功文件输出文件夹 | "JAV_output" | 文件夹路径 |
| link_mode | int | 文件处理模式 | 0 | 0=移动,1=软链接,2=硬链接 |
| scan_hardlink | bool | 扫描硬链接文件 | false | true/false |
| failed_move | bool | 移动失败文件 | true | true/false |
| auto_exit | bool | 自动退出 | false | true/false |
| translate_to_sc | bool | 繁体转简体 | true | true/false |
| actor_gender | string | 演员性别过滤 | "female" | female/male/both/all |
| del_empty_folder | bool | 删除空文件夹 | true | true/false |
| nfo_skip_days | int | NFO跳过天数 | 30 | 0-365 |
| ignore_failed_list | bool | 忽略失败列表 | false | true/false |
| download_only_missing_images | bool | 仅下载缺失图片 | true | true/false |
| mapping_table_validity | int | 映射表有效期(天) | 7 | 1-30 |
| jellyfin | int | Jellyfin兼容模式 | 0 | 0=Kodi,1=Jellyfin |
| actor_only_tag | bool | 仅使用演员名作为标签 | false | true/false |
| sleep | int | 请求间隔(秒) | 3 | 1-60 |
| anonymous_fill | int | 匿名填充模式 | 0 | 0-2 |
| multi_threading | int | 多线程数量 | 0 | 0-20 (0=顺序) |
| stop_counter | int | 停止计数器 | 0 | >= 0 |
| rerun_delay | string | 重运行延迟 | "0" | "0" 或时间格式 |

## 2. Proxy（代理配置）- 6项

| 字段 | 类型 | 说明 | 默认值 | 选项/范围 |
|------|------|------|--------|-----------|
| switch | bool | 启用代理 | false | true/false |
| proxy | string | 代理地址 | "" | host:port |
| timeout | int | 超时时间(秒) | 5 | 5-60 |
| retry | int | 重试次数 | 3 | 1-10 |
| type | string | 代理类型 | "socks5" | socks5/socks5h/http/https |
| cacert_file | string | CA证书文件路径 | "" | 文件路径 |

## 3. NameRule（命名规则）- 6项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| location_rule | string | 文件夹命名规则 | "actor + '/' + number" |
| naming_rule | string | 文件命名规则 | "number + '-' + title" |
| max_title_len | int | 标题最大长度 | 50 |
| image_naming_with_number | bool | 图片名包含番号 | false |
| number_uppercase | bool | 番号大写 | false |
| number_regexs | string | 自定义番号正则 | "" |

## 4. Update（更新配置）- 1项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| update_check | bool | 检查更新 | true |

## 5. Priority（优先级配置）- 1项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| website | string | 数据源优先级 | "javbus,javdb,fanza,..." |

## 6. Escape（转义配置）- 2项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| literals | string | 转义字符 | "\\()/  " |
| folders | string | 转义文件夹 | "failed, JAV_output" |

## 7. DebugMode（调试模式）- 1项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| switch | bool | 启用调试模式 | false |

## 8. Translate（翻译配置）- 7项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| switch | bool | 启用翻译 | false |
| engine | string | 翻译引擎 | "google-free" |
| target_language | string | 目标语言 | "zh_cn" |
| key | string | API密钥 | "" |
| delay | int | 翻译延迟(秒) | 1 |
| values | string | 翻译字段 | "title,outline" |
| service_site | string | 服务站点 | "translate.google.cn" |

## 9. Trailer（预告片）- 1项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| switch | bool | 下载预告片 | false |

## 10. Uncensored（无码配置）- 1项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| uncensored_prefix | string | 无码前缀 | "S2M,BT,LAF,SMD" |

## 11. Media（媒体配置）- 2项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| media_type | string | 视频文件扩展名 | ".mp4,.avi,..." |
| sub_type | string | 字幕文件扩展名 | ".smi,.srt,..." |

## 12. Watermark（水印配置）- 2项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| switch | bool | 启用水印 | true |
| water | int | 水印位置 | 2 |

## 13. Extrafanart（额外剧照）- 3项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| switch | bool | 下载额外剧照 | true |
| extrafanart_folder | string | 剧照文件夹名 | "extrafanart" |
| parallel_download | int | 并行下载数 | 1 |

## 14. Storyline（剧情配置）- 5项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| switch | bool | 启用剧情刮削 | true |
| site | string | 默认站点 | "1:avno1" |
| censored_site | string | 有码站点 | "5:xcity,6:amazon" |
| uncensored_site | string | 无码站点 | "3:58avgo" |
| show_result | int | 显示结果 | 0 |
| run_mode | int | 运行模式 | 1 |

## 15. CCConvert（简繁转换）- 2项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| mode | int | 转换模式 | 1 |
| vars | string | 转换字段 | "actor,director,..." |

## 16. Javdb（Javdb配置）- 1项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| sites | string | 站点列表 | "38,39" |

## 17. Face（人脸识别）- 4项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| locations_model | string | 识别模型 | "hog" |
| uncensored_only | bool | 仅无码 | true |
| always_imagecut | bool | 总是裁剪 | false |
| aspect_ratio | float | 宽高比 | 2.12 |

## 18. Jellyfin（Jellyfin扩展）- 1项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| multi_part_fanart | bool | 多分片剧照 | false |

## 19. ActorPhoto（演员照片）- 1项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| download_for_kodi | bool | 为Kodi下载 | false |

## 20. STRM（STRM文件）- 9项

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| enable | bool | 启用STRM生成 | false |
| path_type | string | 路径类型 | "absolute" |
| content_mode | string | 内容模式 | "simple" |
| multipart_mode | string | 分片模式 | "separate" |
| network_base_path | string | 网络基础路径 | "" |
| use_windows_path | bool | 使用Windows路径 | false |
| validate_files | bool | 验证文件 | true |
| strict_validation | bool | 严格验证 | false |
| output_suffix | string | 输出后缀 | "" |

---

**总计：约75个配置项，分为20个分组**

