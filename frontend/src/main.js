// {{ AURA-X: Modify - 完整配置项读写逻辑,支持所有75个配置项 }}

// 全局状态
let currentPage = 'home';
let currentConfigTab = 'basic';
let isRunning = false;
let logs = [];
let currentLogFilter = 'ALL';
let regexPresets = []; // 存储预定义正则模式
let fileProcessingList = []; // 存储文件处理状态

// 页面初始化
document.addEventListener('DOMContentLoaded', () => {
    console.log('[GUI] 前端初始化...');
    
    // 初始化页面导航
    initNavigation();
    
    // 初始化配置标签页
    initConfigTabs();
    
    // 初始化事件监听
    initEventListeners();
    
    // 加载配置
    loadConfig();
    
    // 加载正则预设
    loadRegexPresets();
    
    // 监听后端事件
    listenToBackendEvents();
    
    console.log('[GUI] 前端初始化完成');
});

// 初始化导航
function initNavigation() {
    const navBtns = document.querySelectorAll('.nav-btn');
    
    navBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const pageName = btn.dataset.page;
            switchPage(pageName);
        });
    });
}

// 初始化配置标签页
function initConfigTabs() {
    const tabs = document.querySelectorAll('.config-tab');
    
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const tabName = tab.dataset.tab;
            switchConfigTab(tabName);
        });
    });
}

// 切换页面
function switchPage(pageName) {
    // 更新导航按钮状态
    document.querySelectorAll('.nav-btn').forEach(btn => {
        btn.classList.remove('active');
        if (btn.dataset.page === pageName) {
            btn.classList.add('active');
        }
    });
    
    // 更新页面显示
    document.querySelectorAll('.page').forEach(page => {
        page.classList.remove('active');
    });
    document.getElementById(`page-${pageName}`).classList.add('active');
    
    currentPage = pageName;
}

// 切换配置标签页
function switchConfigTab(tabName) {
    // 更新标签按钮状态
    document.querySelectorAll('.config-tab').forEach(tab => {
        tab.classList.remove('active');
        if (tab.dataset.tab === tabName) {
            tab.classList.add('active');
        }
    });
    
    // 更新标签页内容显示
    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.remove('active');
    });
    document.getElementById(`tab-${tabName}`).classList.add('active');
    
    currentConfigTab = tabName;
}

// 初始化事件监听器
function initEventListeners() {
    // 运行控制按钮
    document.getElementById('btn-start').addEventListener('click', startScraping);
    document.getElementById('btn-stop').addEventListener('click', stopScraping);
    
    // 配置页面按钮
    document.getElementById('btn-save-config').addEventListener('click', saveConfig);
    document.getElementById('btn-reload-config').addEventListener('click', loadConfig);
    document.getElementById('btn-reset-config').addEventListener('click', resetConfig);
    
    // 文件夹浏览按钮
    document.querySelectorAll('.btn-browse').forEach(btn => {
        btn.addEventListener('click', () => {
            const targetId = btn.dataset.target;
            selectFolder(targetId);
        });
    });
    
    // 导出导入配置
    const exportBtn = document.getElementById('btn-export-config');
    const importBtn = document.getElementById('btn-import-config');
    if (exportBtn) exportBtn.addEventListener('click', exportConfig);
    if (importBtn) importBtn.addEventListener('click', importConfig);
    
    // 日志页面按钮
    document.getElementById('log-filter').addEventListener('click', (e) => {
        currentLogFilter = e.target.value;
        renderLogs();
    });
    document.getElementById('btn-clear-logs').addEventListener('click', clearLogs);
    
    // 文件页面按钮
    document.getElementById('btn-refresh-files').addEventListener('click', loadFileList);
    
    // 正则测试页面按钮
    document.getElementById('btn-load-presets').addEventListener('click', loadRegexPresets);
    document.getElementById('btn-validate-regex').addEventListener('click', validateRegex);
    document.getElementById('btn-test-regex').addEventListener('click', testRegex);
    document.getElementById('btn-suggest-pattern').addEventListener('click', suggestPattern);
    document.getElementById('regex-preset-select').addEventListener('change', onPresetChange);
}

// 开始刮削
async function startScraping() {
    try {
        const sourcePath = getConfigValue('cfg-common-source_folder');
        
        if (!sourcePath) {
            showMessage('error', '请先选择源文件夹');
            return;
        }
        
        // {{ AURA-X: Modify - 启动时清空文件处理列表. Confirmed via 寸止 }}
        // 清空文件处理列表
        fileProcessingList = [];
        renderFileProcessingList();
        
        await window.go.gui.App.Start(sourcePath);
        
        document.getElementById('btn-start').disabled = true;
        document.getElementById('btn-stop').disabled = false;
        isRunning = true;
        
        updateStatus('运行中...', 'info');
        
    } catch (error) {
        showMessage('error', `启动失败: ${error}`);
        console.error('[GUI] 启动失败:', error);
    }
}

// 停止刮削
async function stopScraping() {
    try {
        await window.go.gui.App.Stop();
        
        document.getElementById('btn-start').disabled = false;
        document.getElementById('btn-stop').disabled = true;
        isRunning = false;
        
        updateStatus('已停止', 'warning');
        
    } catch (error) {
        showMessage('error', `停止失败: ${error}`);
        console.error('[GUI] 停止失败:', error);
    }
}

// 加载配置
async function loadConfig() {
    try {
        const config = await window.go.gui.App.GetConfig();
        
        if (!config) {
            console.warn('[GUI] 配置为空');
            return;
        }
        
        console.log('[GUI] 加载的配置:', config);
        
        // 填充所有配置字段
        fillConfigFields(config);
        
        showMessage('info', '配置加载成功');
        console.log('[GUI] 配置加载成功');
        
    } catch (error) {
        showMessage('error', `加载配置失败: ${error}`);
        console.error('[GUI] 加载配置失败:', error);
    }
}

// 填充配置字段到表单
function fillConfigFields(config) {
    // Common配置
    if (config.common) {
        setConfigValue('cfg-common-main_mode', config.common.main_mode);
        setConfigValue('cfg-common-source_folder', config.common.source_folder);
        setConfigValue('cfg-common-failed_output_folder', config.common.failed_output_folder);
        setConfigValue('cfg-common-success_output_folder', config.common.success_output_folder);
        setConfigValue('cfg-common-link_mode', config.common.link_mode);
        setConfigValue('cfg-common-scan_hardlink', config.common.scan_hardlink);
        setConfigValue('cfg-common-failed_move', config.common.failed_move);
        setConfigValue('cfg-common-auto_exit', config.common.auto_exit);
        setConfigValue('cfg-common-translate_to_sc', config.common.translate_to_sc);
        setConfigValue('cfg-common-actor_gender', config.common.actor_gender);
        setConfigValue('cfg-common-del_empty_folder', config.common.del_empty_folder);
        setConfigValue('cfg-common-nfo_skip_days', config.common.nfo_skip_days);
        setConfigValue('cfg-common-ignore_failed_list', config.common.ignore_failed_list);
        setConfigValue('cfg-common-download_only_missing_images', config.common.download_only_missing_images);
        setConfigValue('cfg-common-mapping_table_validity', config.common.mapping_table_validity);
        setConfigValue('cfg-common-jellyfin', config.common.jellyfin);
        setConfigValue('cfg-common-actor_only_tag', config.common.actor_only_tag);
        setConfigValue('cfg-common-sleep', config.common.sleep);
        setConfigValue('cfg-common-anonymous_fill', config.common.anonymous_fill);
        setConfigValue('cfg-common-multi_threading', config.common.multi_threading);
        setConfigValue('cfg-common-stop_counter', config.common.stop_counter);
        setConfigValue('cfg-common-rerun_delay', config.common.rerun_delay);
    }
    
    // Proxy配置
    if (config.proxy) {
        setConfigValue('cfg-proxy-switch', config.proxy.switch);
        setConfigValue('cfg-proxy-proxy', config.proxy.proxy);
        setConfigValue('cfg-proxy-timeout', config.proxy.timeout);
        setConfigValue('cfg-proxy-retry', config.proxy.retry);
        setConfigValue('cfg-proxy-type', config.proxy.type);
        setConfigValue('cfg-proxy-cacert_file', config.proxy.cacert_file);
    }
    
    // NameRule配置
    if (config.name_rule) {
        setConfigValue('cfg-name_rule-location_rule', config.name_rule.location_rule);
        setConfigValue('cfg-name_rule-naming_rule', config.name_rule.naming_rule);
        setConfigValue('cfg-name_rule-max_title_len', config.name_rule.max_title_len);
        setConfigValue('cfg-name_rule-image_naming_with_number', config.name_rule.image_naming_with_number);
        setConfigValue('cfg-name_rule-number_uppercase', config.name_rule.number_uppercase);
        setConfigValue('cfg-name_rule-number_regexs', config.name_rule.number_regexs);
    }
    
    // Update配置
    if (config.update) {
        setConfigValue('cfg-update-update_check', config.update.update_check);
    }
    
    // Priority配置
    if (config.priority) {
        setConfigValue('cfg-priority-website', config.priority.website);
    }
    
    // Escape配置
    if (config.escape) {
        setConfigValue('cfg-escape-literals', config.escape.literals);
        setConfigValue('cfg-escape-folders', config.escape.folders);
    }
    
    // DebugMode配置
    if (config.debug_mode) {
        setConfigValue('cfg-debug_mode-switch', config.debug_mode.switch);
    }
    
    // Translate配置
    if (config.translate) {
        setConfigValue('cfg-translate-switch', config.translate.switch);
        setConfigValue('cfg-translate-engine', config.translate.engine);
        setConfigValue('cfg-translate-target_language', config.translate.target_language);
        setConfigValue('cfg-translate-key', config.translate.key);
        setConfigValue('cfg-translate-delay', config.translate.delay);
        setConfigValue('cfg-translate-values', config.translate.values);
        setConfigValue('cfg-translate-service_site', config.translate.service_site);
    }
    
    // Trailer配置
    if (config.trailer) {
        setConfigValue('cfg-trailer-switch', config.trailer.switch);
    }
    
    // Uncensored配置
    if (config.uncensored) {
        setConfigValue('cfg-uncensored-uncensored_prefix', config.uncensored.uncensored_prefix);
    }
    
    // Media配置
    if (config.media) {
        setConfigValue('cfg-media-media_type', config.media.media_type);
        setConfigValue('cfg-media-sub_type', config.media.sub_type);
    }
    
    // Watermark配置
    if (config.watermark) {
        setConfigValue('cfg-watermark-switch', config.watermark.switch);
        setConfigValue('cfg-watermark-water', config.watermark.water);
    }
    
    // Extrafanart配置
    if (config.extrafanart) {
        setConfigValue('cfg-extrafanart-switch', config.extrafanart.switch);
        setConfigValue('cfg-extrafanart-extrafanart_folder', config.extrafanart.extrafanart_folder);
        setConfigValue('cfg-extrafanart-parallel_download', config.extrafanart.parallel_download);
    }
    
    // Storyline配置
    if (config.storyline) {
        setConfigValue('cfg-storyline-switch', config.storyline.switch);
        setConfigValue('cfg-storyline-site', config.storyline.site);
        setConfigValue('cfg-storyline-censored_site', config.storyline.censored_site);
        setConfigValue('cfg-storyline-uncensored_site', config.storyline.uncensored_site);
        setConfigValue('cfg-storyline-show_result', config.storyline.show_result);
        setConfigValue('cfg-storyline-run_mode', config.storyline.run_mode);
    }
    
    // CCConvert配置
    if (config.cc_convert) {
        setConfigValue('cfg-cc_convert-mode', config.cc_convert.mode);
        setConfigValue('cfg-cc_convert-vars', config.cc_convert.vars);
    }
    
    // Javdb配置
    if (config.javdb) {
        setConfigValue('cfg-javdb-sites', config.javdb.sites);
    }
    
    // Face配置
    if (config.face) {
        setConfigValue('cfg-face-locations_model', config.face.locations_model);
        setConfigValue('cfg-face-uncensored_only', config.face.uncensored_only);
        setConfigValue('cfg-face-always_imagecut', config.face.always_imagecut);
        setConfigValue('cfg-face-aspect_ratio', config.face.aspect_ratio);
    }
    
    // Jellyfin配置
    if (config.jellyfin) {
        setConfigValue('cfg-jellyfin-multi_part_fanart', config.jellyfin.multi_part_fanart);
    }
    
    // ActorPhoto配置
    if (config.actor_photo) {
        setConfigValue('cfg-actor_photo-download_for_kodi', config.actor_photo.download_for_kodi);
    }
    
    // STRM配置
    if (config.strm) {
        setConfigValue('cfg-strm-enable', config.strm.enable);
        setConfigValue('cfg-strm-path_type', config.strm.path_type);
        setConfigValue('cfg-strm-content_mode', config.strm.content_mode);
        setConfigValue('cfg-strm-multipart_mode', config.strm.multipart_mode);
        setConfigValue('cfg-strm-network_base_path', config.strm.network_base_path);
        setConfigValue('cfg-strm-use_windows_path', config.strm.use_windows_path);
        setConfigValue('cfg-strm-validate_files', config.strm.validate_files);
        setConfigValue('cfg-strm-strict_validation', config.strm.strict_validation);
        setConfigValue('cfg-strm-output_suffix', config.strm.output_suffix);
    }
    
    // Scraper配置
    if (config.scraper) {
        setConfigValue('cfg-scraper-mode', config.scraper.mode);
        setConfigValue('cfg-scraper-metatube_url', config.scraper.metatube_url);
        setConfigValue('cfg-scraper-metatube_token', config.scraper.metatube_token);
        setConfigValue('cfg-scraper-fallback_to_legacy', config.scraper.fallback_to_legacy);
    }
}

// 从表单收集配置
function collectConfigFromForm() {
    return {
        common: {
            main_mode: parseInt(getConfigValue('cfg-common-main_mode')) || 1,
            source_folder: getConfigValue('cfg-common-source_folder') || './',
            failed_output_folder: getConfigValue('cfg-common-failed_output_folder') || 'failed',
            success_output_folder: getConfigValue('cfg-common-success_output_folder') || 'JAV_output',
            link_mode: parseInt(getConfigValue('cfg-common-link_mode')) || 0,
            scan_hardlink: getConfigValue('cfg-common-scan_hardlink') || false,
            failed_move: getConfigValue('cfg-common-failed_move') !== false,
            auto_exit: getConfigValue('cfg-common-auto_exit') || false,
            translate_to_sc: getConfigValue('cfg-common-translate_to_sc') !== false,
            actor_gender: getConfigValue('cfg-common-actor_gender') || 'female',
            del_empty_folder: getConfigValue('cfg-common-del_empty_folder') !== false,
            nfo_skip_days: parseInt(getConfigValue('cfg-common-nfo_skip_days')) || 30,
            ignore_failed_list: getConfigValue('cfg-common-ignore_failed_list') || false,
            download_only_missing_images: getConfigValue('cfg-common-download_only_missing_images') !== false,
            mapping_table_validity: parseInt(getConfigValue('cfg-common-mapping_table_validity')) || 7,
            jellyfin: parseInt(getConfigValue('cfg-common-jellyfin')) || 0,
            actor_only_tag: getConfigValue('cfg-common-actor_only_tag') || false,
            sleep: parseInt(getConfigValue('cfg-common-sleep')) || 3,
            anonymous_fill: parseInt(getConfigValue('cfg-common-anonymous_fill')) || 0,
            multi_threading: parseInt(getConfigValue('cfg-common-multi_threading')) || 0,
            stop_counter: parseInt(getConfigValue('cfg-common-stop_counter')) || 0,
            rerun_delay: getConfigValue('cfg-common-rerun_delay') || '0',
        },
        proxy: {
            switch: getConfigValue('cfg-proxy-switch') || false,
            proxy: getConfigValue('cfg-proxy-proxy') || '',
            timeout: parseInt(getConfigValue('cfg-proxy-timeout')) || 30,
            retry: parseInt(getConfigValue('cfg-proxy-retry')) || 5,
            type: getConfigValue('cfg-proxy-type') || 'socks5',
            cacert_file: getConfigValue('cfg-proxy-cacert_file') || '',
        },
        name_rule: {
            location_rule: getConfigValue('cfg-name_rule-location_rule') || "actor + '/' + number",
            naming_rule: getConfigValue('cfg-name_rule-naming_rule') || "number + '-' + title",
            max_title_len: parseInt(getConfigValue('cfg-name_rule-max_title_len')) || 50,
            image_naming_with_number: getConfigValue('cfg-name_rule-image_naming_with_number') || false,
            number_uppercase: getConfigValue('cfg-name_rule-number_uppercase') || false,
            number_regexs: getConfigValue('cfg-name_rule-number_regexs') || '',
        },
        update: {
            update_check: getConfigValue('cfg-update-update_check') !== false,
        },
        priority: {
            website: getConfigValue('cfg-priority-website') || 'javbus,fanza,fc2,fc2club,javdb',
        },
        escape: {
            literals: getConfigValue('cfg-escape-literals') || '\\()/ ',
            folders: getConfigValue('cfg-escape-folders') || 'failed, JAV_output',
        },
        debug_mode: {
            switch: getConfigValue('cfg-debug_mode-switch') || false,
        },
        translate: {
            switch: getConfigValue('cfg-translate-switch') || false,
            engine: getConfigValue('cfg-translate-engine') || 'google-free',
            target_language: getConfigValue('cfg-translate-target_language') || 'zh_cn',
            key: getConfigValue('cfg-translate-key') || '',
            delay: parseInt(getConfigValue('cfg-translate-delay')) || 1,
            values: getConfigValue('cfg-translate-values') || 'title,outline',
            service_site: getConfigValue('cfg-translate-service_site') || 'translate.google.cn',
        },
        trailer: {
            switch: getConfigValue('cfg-trailer-switch') || false,
        },
        uncensored: {
            uncensored_prefix: getConfigValue('cfg-uncensored-uncensored_prefix') || 'S2M,BT,LAF,SMD',
        },
        media: {
            media_type: getConfigValue('cfg-media-media_type') || '.mp4,.avi,.rmvb,.wmv,.mov,.mkv,.flv,.ts,.webm,.iso',
            sub_type: getConfigValue('cfg-media-sub_type') || '.smi,.srt,.idx,.sub,.sup,.psb,.ssa,.ass',
        },
        watermark: {
            switch: getConfigValue('cfg-watermark-switch') !== false,
            water: parseInt(getConfigValue('cfg-watermark-water')) || 2,
        },
        extrafanart: {
            switch: getConfigValue('cfg-extrafanart-switch') !== false,
            extrafanart_folder: getConfigValue('cfg-extrafanart-extrafanart_folder') || 'extrafanart',
            parallel_download: parseInt(getConfigValue('cfg-extrafanart-parallel_download')) || 1,
        },
        storyline: {
            switch: getConfigValue('cfg-storyline-switch') !== false,
            site: getConfigValue('cfg-storyline-site') || '1:avno1',
            censored_site: getConfigValue('cfg-storyline-censored_site') || '5:xcity,6:amazon',
            uncensored_site: getConfigValue('cfg-storyline-uncensored_site') || '3:58avgo',
            show_result: parseInt(getConfigValue('cfg-storyline-show_result')) || 0,
            run_mode: parseInt(getConfigValue('cfg-storyline-run_mode')) || 1,
        },
        cc_convert: {
            mode: parseInt(getConfigValue('cfg-cc_convert-mode')) || 1,
            vars: getConfigValue('cfg-cc_convert-vars') || 'actor,director,label,outline,series,studio,tag,title',
        },
        javdb: {
            sites: getConfigValue('cfg-javdb-sites') || '38,39',
        },
        face: {
            locations_model: getConfigValue('cfg-face-locations_model') || 'hog',
            uncensored_only: getConfigValue('cfg-face-uncensored_only') !== false,
            always_imagecut: getConfigValue('cfg-face-always_imagecut') || false,
            aspect_ratio: parseFloat(getConfigValue('cfg-face-aspect_ratio')) || 2.12,
        },
        jellyfin: {
            multi_part_fanart: getConfigValue('cfg-jellyfin-multi_part_fanart') || false,
        },
        actor_photo: {
            download_for_kodi: getConfigValue('cfg-actor_photo-download_for_kodi') || false,
        },
        strm: {
            enable: getConfigValue('cfg-strm-enable') || false,
            path_type: getConfigValue('cfg-strm-path_type') || 'absolute',
            content_mode: getConfigValue('cfg-strm-content_mode') || 'simple',
            multipart_mode: getConfigValue('cfg-strm-multipart_mode') || 'separate',
            network_base_path: getConfigValue('cfg-strm-network_base_path') || '',
            use_windows_path: getConfigValue('cfg-strm-use_windows_path') || false,
            validate_files: getConfigValue('cfg-strm-validate_files') !== false,
            strict_validation: getConfigValue('cfg-strm-strict_validation') || false,
            output_suffix: getConfigValue('cfg-strm-output_suffix') || '',
        },
        scraper: {
            mode: getConfigValue('cfg-scraper-mode') || 'legacy',
            metatube_url: getConfigValue('cfg-scraper-metatube_url') || 'http://localhost:8080',
            metatube_token: getConfigValue('cfg-scraper-metatube_token') || '',
            fallback_to_legacy: getConfigValue('cfg-scraper-fallback_to_legacy') !== false,
        },
    };
}

// 保存配置
async function saveConfig() {
    try {
        // 收集表单数据
        const config = collectConfigFromForm();
        
        console.log('[GUI] 保存配置:', config);
        
        // 调用后端保存
        await window.go.gui.App.SaveConfig(config);
        
        showMessage('success', '配置保存成功');
        console.log('[GUI] 配置保存成功');
        
    } catch (error) {
        showMessage('error', `保存配置失败: ${error}`);
        console.error('[GUI] 保存配置失败:', error);
    }
}

// 重置配置
async function resetConfig() {
    if (!confirm('确定要重置配置为默认值吗？当前配置将被备份。')) {
        return;
    }
    
    try {
        await window.go.gui.App.ResetConfig();
        await loadConfig();
        showMessage('success', '配置已重置为默认值');
    } catch (error) {
        showMessage('error', `重置配置失败: ${error}`);
        console.error('[GUI] 重置配置失败:', error);
    }
}

// 选择文件夹
async function selectFolder(targetElementId) {
    try {
        const folder = await window.go.gui.App.SelectFolder('选择文件夹');
        
        if (folder) {
            setConfigValue(targetElementId, folder);
        }
        
    } catch (error) {
        showMessage('error', `选择文件夹失败: ${error}`);
        console.error('[GUI] 选择文件夹失败:', error);
    }
}

// 导出配置
function exportConfig() {
    const config = collectConfigFromForm();
    const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(config, null, 2));
    const downloadAnchorNode = document.createElement('a');
    downloadAnchorNode.setAttribute("href", dataStr);
    downloadAnchorNode.setAttribute("download", "config_export.json");
    document.body.appendChild(downloadAnchorNode);
    downloadAnchorNode.click();
    downloadAnchorNode.remove();
    showMessage('info', '配置已导出');
}

// 导入配置
function importConfig() {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json';
    input.onchange = e => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onload = event => {
            try {
                const config = JSON.parse(event.target.result);
                fillConfigFields(config);
                showMessage('success', '配置已导入，请点击保存按钮');
            } catch (error) {
                showMessage('error', '配置文件格式错误');
            }
        };
        reader.readAsText(file);
    };
    input.click();
}

// 获取配置值（通用方法）
function getConfigValue(elementId) {
    const element = document.getElementById(elementId);
    if (!element) {
        console.warn(`[GUI] 元素不存在: ${elementId}`);
        return null;
    }
    
    if (element.type === 'checkbox') {
        return element.checked;
    } else if (element.type === 'number') {
        const value = element.value;
        return value === '' ? null : (element.step && element.step.includes('.') ? parseFloat(value) : parseInt(value));
    } else {
        return element.value;
    }
}

// 设置配置值（通用方法）
function setConfigValue(elementId, value) {
    const element = document.getElementById(elementId);
    if (!element) {
        console.warn(`[GUI] 元素不存在: ${elementId}`);
        return;
    }
    
    if (element.type === 'checkbox') {
        element.checked = Boolean(value);
    } else {
        element.value = value !== null && value !== undefined ? value : '';
    }
}

// 监听后端事件
function listenToBackendEvents() {
    // 监听日志事件
    window.runtime.EventsOn('log', (data) => {
        addLog(data);
    });
    
    // 监听进度事件
    window.runtime.EventsOn('progress', (data) => {
        updateProgress(data);
    });
    
    // {{ AURA-X: Add - 监听文件处理状态事件. Confirmed via 寸止 }}
    // 监听文件处理状态事件
    window.runtime.EventsOn('file_status', (data) => {
        updateFileStatus(data);
    });
}

// 添加日志
function addLog(logData) {
    logs.push(logData);
    
    // 限制日志数量
    if (logs.length > 1000) {
        logs = logs.slice(-1000);
    }
    
    renderLogs();
}

// 渲染日志
function renderLogs() {
    const container = document.getElementById('logs-container');
    
    // 过滤日志
    const filteredLogs = logs.filter(log => {
        if (currentLogFilter === 'ALL') return true;
        return log.level === currentLogFilter;
    });
    
    // 渲染
    container.innerHTML = filteredLogs.map(log => {
        return `
            <div class="log-entry">
                <span class="log-time">${log.time}</span>
                <span class="log-level log-level-${log.level}">${log.level}</span>
                <span class="log-message">${escapeHtml(log.message)}</span>
            </div>
        `;
    }).join('');
    
    // 自动滚动到底部
    container.scrollTop = container.scrollHeight;
}

// 清空日志
function clearLogs() {
    logs = [];
    renderLogs();
}

// 更新进度
function updateProgress(data) {
    document.getElementById('stat-total').textContent = data.total || 0;
    document.getElementById('stat-success').textContent = data.success || 0;
    document.getElementById('stat-failed').textContent = data.failed || 0;
    document.getElementById('stat-skipped').textContent = data.skipped || 0;
    
    const progress = data.total > 0 ? ((data.success + data.failed + data.skipped) / data.total * 100) : 0;
    document.getElementById('progress-fill').style.width = `${progress}%`;
    
    if (data.running) {
        document.getElementById('progress-text').textContent = `处理中... ${Math.round(progress)}%`;
        document.getElementById('progress-time').textContent = data.duration || '';
        updateStatus('运行中...', 'info');
    } else {
        document.getElementById('progress-text').textContent = '就绪';
        document.getElementById('progress-time').textContent = '';
        updateStatus('就绪', 'success');
        
        // 任务完成后重置按钮
        document.getElementById('btn-start').disabled = false;
        document.getElementById('btn-stop').disabled = true;
    }
}

// 加载文件列表
async function loadFileList() {
    try {
        const files = await window.go.gui.App.GetFileList();
        const container = document.getElementById('files-container');
        
        if (!files || files.length === 0) {
            container.innerHTML = '<p class="empty-message">没有找到视频文件</p>';
            return;
        }
        
        container.innerHTML = files.map(file => {
            const sizeText = formatFileSize(file.size);
            return `
                <div class="file-item">
                    <div class="file-name">📹 ${escapeHtml(file.name)}</div>
                    <div class="file-info">
                        番号: ${escapeHtml(file.number || '未识别')} | 
                        大小: ${sizeText}
                    </div>
                </div>
            `;
        }).join('');
        
        showMessage('info', `找到 ${files.length} 个视频文件`);
        
    } catch (error) {
        showMessage('error', `加载文件列表失败: ${error}`);
        console.error('[GUI] 加载文件列表失败:', error);
    }
}

// 更新状态
function updateStatus(text, type) {
    const statusText = document.getElementById('status-text');
    statusText.textContent = text;
    statusText.className = `status-${type}`;
}

// 显示消息
function showMessage(type, message) {
    console.log(`[${type.toUpperCase()}] ${message}`);
    // 可以在这里添加Toast通知
    updateStatus(message, type);
}

// 辅助函数：格式化文件大小
function formatFileSize(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// 辅助函数：HTML转义
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ==================== 正则测试功能 ====================

// 加载预定义正则模式
async function loadRegexPresets() {
    try {
        const presets = await window.go.gui.App.GetDefaultRegexPatterns();
        regexPresets = presets;
        
        const select = document.getElementById('regex-preset-select');
        select.innerHTML = '<option value="">-- 选择预定义模式 --</option>';
        
        presets.forEach((preset, index) => {
            const option = document.createElement('option');
            option.value = index;
            option.textContent = preset.name;
            select.appendChild(option);
        });
        
        showMessage('info', `加载了 ${presets.length} 个预定义模式`);
        
    } catch (error) {
        showMessage('error', `加载预设失败: ${error}`);
        console.error('[GUI] 加载预设失败:', error);
    }
}

// 预设选择改变
function onPresetChange() {
    const select = document.getElementById('regex-preset-select');
    const index = parseInt(select.value);
    
    if (isNaN(index) || index < 0 || index >= regexPresets.length) {
        document.getElementById('preset-description').style.display = 'none';
        return;
    }
    
    const preset = regexPresets[index];
    
    // {{ AURA-X: Modify - 支持多行显示详细说明. Confirmed via 寸止 }}
    // 显示描述（支持换行符）
    const descElement = document.getElementById('preset-desc-text');
    descElement.innerHTML = escapeHtml(preset.description).replace(/\n/g, '<br>');
    
    const exampleElement = document.getElementById('preset-example-text');
    exampleElement.innerHTML = escapeHtml(preset.example).replace(/\n/g, '<br>');
    
    document.getElementById('preset-description').style.display = 'block';
    
    // 填充正则表达式
    document.getElementById('regex-pattern-input').value = preset.pattern;
    
    // 自动验证
    validateRegex();
}

// 验证正则表达式
async function validateRegex() {
    const pattern = document.getElementById('regex-pattern-input').value.trim();
    const validationDiv = document.getElementById('regex-validation');
    
    if (!pattern) {
        validationDiv.innerHTML = '';
        return;
    }
    
    try {
        const result = await window.go.gui.App.ValidateRegex(pattern);
        
        if (result.valid) {
            validationDiv.innerHTML = `<span style="color: green;">✓ ${escapeHtml(result.message)}</span>`;
        } else {
            validationDiv.innerHTML = `<span style="color: red;">✗ ${escapeHtml(result.message)}</span>`;
        }
        
    } catch (error) {
        validationDiv.innerHTML = `<span style="color: red;">验证失败: ${escapeHtml(error.toString())}</span>`;
        console.error('[GUI] 验证正则失败:', error);
    }
}

// 测试正则表达式
async function testRegex() {
    const pattern = document.getElementById('regex-pattern-input').value.trim();
    const filenamesText = document.getElementById('regex-filenames-input').value.trim();
    
    if (!pattern) {
        showMessage('error', '请输入正则表达式');
        return;
    }
    
    if (!filenamesText) {
        showMessage('error', '请输入要测试的文件名');
        return;
    }
    
    // 分割文件名
    const filenames = filenamesText.split('\n')
        .map(line => line.trim())
        .filter(line => line.length > 0);
    
    if (filenames.length === 0) {
        showMessage('error', '没有有效的文件名');
        return;
    }
    
    try {
        const results = await window.go.gui.App.TestRegexPattern({
            pattern: pattern,
            filenames: filenames
        });
        
        renderTestResults(results);
        showMessage('success', `测试完成，共测试 ${results.length} 个文件名`);
        
    } catch (error) {
        showMessage('error', `测试失败: ${error}`);
        console.error('[GUI] 测试正则失败:', error);
    }
}

// 渲染测试结果
function renderTestResults(results) {
    const container = document.getElementById('regex-test-results');
    
    if (!results || results.length === 0) {
        container.innerHTML = '<p class="empty-message">没有测试结果</p>';
        return;
    }
    
    let successCount = 0;
    let failCount = 0;
    
    let html = '<div class="test-results-summary">';
    
    results.forEach(result => {
        if (result.success) {
            successCount++;
            html += `
                <div class="test-result-item success">
                    <div class="result-header">
                        <span class="result-status">✓ 匹配成功</span>
                        <span class="result-filename">${escapeHtml(result.originalName)}</span>
                    </div>
                    <div class="result-details">
                        <div><strong>完整匹配:</strong> <code>${escapeHtml(result.matched)}</code></div>
                        <div><strong>提取番号:</strong> <code class="extracted-number">${escapeHtml(result.extractedNumber)}</code></div>
                        ${result.groups && result.groups.length > 0 ? `<div><strong>捕获组:</strong> ${result.groups.map(g => `<code>${escapeHtml(g)}</code>`).join(', ')}</div>` : ''}
                    </div>
                </div>
            `;
        } else {
            failCount++;
            html += `
                <div class="test-result-item failed">
                    <div class="result-header">
                        <span class="result-status">✗ 匹配失败</span>
                        <span class="result-filename">${escapeHtml(result.originalName)}</span>
                    </div>
                    <div class="result-details">
                        <div class="error-message">${escapeHtml(result.error)}</div>
                    </div>
                </div>
            `;
        }
    });
    
    html += '</div>';
    
    // 添加统计信息
    const summary = `
        <div class="test-summary">
            <h4>测试统计</h4>
            <p>总数: ${results.length} | 成功: <span style="color: green;">${successCount}</span> | 失败: <span style="color: red;">${failCount}</span> | 成功率: ${(successCount / results.length * 100).toFixed(1)}%</p>
        </div>
    `;
    
    container.innerHTML = summary + html;
}

// 智能推荐正则模式
async function suggestPattern() {
    const filenamesText = document.getElementById('regex-filenames-input').value.trim();
    
    if (!filenamesText) {
        showMessage('error', '请先输入文件名');
        return;
    }
    
    // 取第一个文件名进行推荐
    const firstFilename = filenamesText.split('\n')[0].trim();
    
    if (!firstFilename) {
        showMessage('error', '没有有效的文件名');
        return;
    }
    
    try {
        const suggestions = await window.go.gui.App.SuggestRegexPattern(firstFilename);
        
        if (!suggestions || suggestions.length === 0) {
            showMessage('warning', '没有找到合适的预定义模式，请尝试自定义正则');
            return;
        }
        
        // 使用第一个推荐
        const recommended = suggestions[0];
        
        document.getElementById('regex-pattern-input').value = recommended.pattern;
        
        // 支持多行显示
        const descElement = document.getElementById('preset-desc-text');
        descElement.innerHTML = escapeHtml(recommended.description).replace(/\n/g, '<br>');
        
        const exampleElement = document.getElementById('preset-example-text');
        exampleElement.innerHTML = escapeHtml(recommended.example).replace(/\n/g, '<br>');
        
        document.getElementById('preset-description').style.display = 'block';
        
        // 自动验证并测试
        await validateRegex();
        
        showMessage('success', `推荐使用: ${recommended.name}，找到 ${suggestions.length} 个匹配模式`);
        
    } catch (error) {
        showMessage('error', `推荐失败: ${error}`);
        console.error('[GUI] 推荐模式失败:', error);
    }
}

// ==================== 文件处理状态列表 ====================
// {{ AURA-X: Modify - 按状态分类展示文件处理列表. Approval: 寸止 }}

// 分类折叠状态
let categoryCollapseState = {
    processing: false,
    success: false,
    failed: false,
    skipped: false
};

// 更新文件处理状态
function updateFileStatus(data) {
    // 查找是否已存在
    const index = fileProcessingList.findIndex(f => f.path === data.path);
    
    if (index >= 0) {
        // 更新现有项
        fileProcessingList[index] = data;
    } else {
        // 添加新项
        fileProcessingList.push(data);
    }
    
    // 渲染列表
    renderFileProcessingList();
}

// 渲染文件处理列表（分类展示）
function renderFileProcessingList() {
    const container = document.getElementById('file-processing-list');
    
    if (!container) return;
    
    // 如果没有任何文件，显示空消息
    if (fileProcessingList.length === 0) {
        // 隐藏所有分类区域
        const categories = ['processing', 'success', 'failed', 'skipped'];
        categories.forEach(cat => {
            const section = document.getElementById(`category-${cat}`);
            if (section) section.style.display = 'none';
        });
        
        // 显示空消息（如果有的话）
        const emptyMsg = container.querySelector('.empty-message');
        if (emptyMsg) emptyMsg.style.display = 'block';
        return;
    }
    
    // 隐藏空消息
    const emptyMsg = container.querySelector('.empty-message');
    if (emptyMsg) emptyMsg.style.display = 'none';
    
    // 按状态分组
    const grouped = {
        processing: [],
        success: [],
        failed: [],
        skipped: []
    };
    
    fileProcessingList.forEach(file => {
        const status = file.status || 'skipped';
        if (grouped[status]) {
            grouped[status].push(file);
        }
    });
    
    // 渲染每个分类
    renderCategory('processing', grouped.processing);
    renderCategory('success', grouped.success);
    renderCategory('failed', grouped.failed);
    renderCategory('skipped', grouped.skipped);
}

// 渲染单个分类
function renderCategory(status, files) {
    const section = document.getElementById(`category-${status}`);
    const content = document.getElementById(`content-${status}`);
    const countEl = document.getElementById(`count-${status}`);
    
    if (!section || !content || !countEl) return;
    
    // 更新计数
    countEl.textContent = `(${files.length})`;
    
    // 如果该分类没有文件，隐藏整个区域
    if (files.length === 0) {
        section.style.display = 'none';
        return;
    }
    
    // 显示区域
    section.style.display = 'block';
    
    // 生成文件列表HTML（最新的在上面）
    let html = '';
    for (let i = files.length - 1; i >= 0; i--) {
        const file = files[i];
        const statusClass = getStatusClass(file.status);
        const sizeText = formatFileSize(file.size);
        
        html += `
            <div class="file-process-item ${statusClass}">
                <div class="file-process-header">
                    <span class="file-process-name">${escapeHtml(file.name)}</span>
                    ${file.duration ? `<span class="file-duration">${escapeHtml(file.duration)}</span>` : ''}
                </div>
                <div class="file-process-details">
                    ${file.number ? `<div class="file-detail-item"><strong>番号:</strong> <code>${escapeHtml(file.number)}</code></div>` : ''}
                    <div class="file-detail-item"><strong>大小:</strong> ${sizeText}</div>
                    ${file.error ? `<div class="file-detail-error">${escapeHtml(file.error)}</div>` : ''}
                </div>
            </div>
        `;
    }
    
    content.innerHTML = html;
    
    // 根据折叠状态显示/隐藏内容
    if (categoryCollapseState[status]) {
        content.style.display = 'none';
    } else {
        content.style.display = 'block';
    }
}

// 切换分类折叠状态
function toggleCategory(status) {
    categoryCollapseState[status] = !categoryCollapseState[status];
    
    const content = document.getElementById(`content-${status}`);
    const toggle = document.getElementById(`toggle-${status}`);
    
    if (!content || !toggle) return;
    
    if (categoryCollapseState[status]) {
        // 折叠
        content.style.display = 'none';
        toggle.textContent = '▶';
    } else {
        // 展开
        content.style.display = 'block';
        toggle.textContent = '▼';
    }
}

// 获取状态样式类
function getStatusClass(status) {
    switch (status) {
        case 'processing': return 'status-processing';
        case 'success': return 'status-success';
        case 'failed': return 'status-failed';
        case 'skipped': return 'status-skipped';
        default: return '';
    }
}

// 获取状态图标
function getStatusIcon(status) {
    switch (status) {
        case 'processing': return '⏳';
        case 'success': return '✅';
        case 'failed': return '❌';
        case 'skipped': return '⏭️';
        default: return '📄';
    }
}

// ==================== 命名规则模板功能 ====================
// {{ AURA-X: Add - 命名规则预设模板功能. Source: AURA-X协议 }}

// 命名规则预设模板
const namingTemplates = {
    jellyfin: {
        name: 'Jellyfin/Emby推荐',
        location_rule: "actor + '/' + number",
        naming_rule: "number + '-' + title",
        preview: '波多野结衣/SSIS-123/SSIS-123-美丽的诱惑.mp4'
    },
    simple: {
        name: '简洁格式（仅番号）',
        location_rule: "number",
        naming_rule: "number",
        preview: 'SSIS-123/SSIS-123.mp4'
    },
    detailed: {
        name: '详细格式（含制作商和演员）',
        location_rule: "studio + '/' + number",
        naming_rule: "number + ' ' + actor + ' ' + title",
        preview: 'S1/SSIS-123/SSIS-123 波多野结衣 美丽的诱惑.mp4'
    },
    by_year: {
        name: '按年份分类',
        location_rule: "year + '/' + actor + '/' + number",
        naming_rule: "number + '-' + title",
        preview: '2024/波多野结衣/SSIS-123/SSIS-123-美丽的诱惑.mp4'
    },
    by_studio: {
        name: '按制作商分类',
        location_rule: "studio + '/' + actor + '/' + number",
        naming_rule: "number + ' ' + title",
        preview: 'S1/波多野结衣/SSIS-123/SSIS-123 美丽的诱惑.mp4'
    }
};

// 应用命名规则模板
function applyNamingTemplate() {
    const select = document.getElementById('naming-template-select');
    const templateKey = select.value;
    
    if (!templateKey || templateKey === 'custom') {
        // 隐藏预览
        document.getElementById('template-preview').style.display = 'none';
        return;
    }
    
    const template = namingTemplates[templateKey];
    if (!template) {
        console.warn('[GUI] 未找到模板:', templateKey);
        return;
    }
    
    // 填充规则
    setConfigValue('cfg-name_rule-location_rule', template.location_rule);
    setConfigValue('cfg-name_rule-naming_rule', template.naming_rule);
    
    // 显示预览
    const previewDiv = document.getElementById('template-preview');
    const previewText = document.getElementById('template-preview-text');
    
    previewText.innerHTML = `
        <div style="line-height: 1.6;">
            <strong>${escapeHtml(template.name)}</strong><br>
            文件夹规则: <code>${escapeHtml(template.location_rule)}</code><br>
            文件命名规则: <code>${escapeHtml(template.naming_rule)}</code><br>
            <span style="color: #2e7d32;">→ ${escapeHtml(template.preview)}</span>
        </div>
    `;
    
    previewDiv.style.display = 'block';
    
    showMessage('success', `已应用模板：${template.name}`);
    console.log('[GUI] 应用命名规则模板:', templateKey, template);
}

// 将函数绑定到全局，以便HTML可以调用
if (typeof window !== 'undefined') {
    window.applyNamingTemplate = applyNamingTemplate;
}
