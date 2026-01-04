const { Modal, Notice, Setting } = require('obsidian');
const fs = require('fs');
const path = require('path');
const os = require('os');

class LogsModal extends Modal {
	constructor(app, vaultPath) {
		super(app);
		this.vaultPath = vaultPath;
		this.selectedLogFile = '';
		this.logFiles = [];
		this.filterLevel = 'all';
	}

	getLevelColor(level) {
		const levelUpper = (level || 'INFO').toUpperCase();
		switch (levelUpper) {
			case 'ERROR': return '#f14c4c';
			case 'WARN': return '#cca700';
			case 'INFO': return '#4ec9b0';
			case 'DEBUG': return '#9cdcfe';
			default: return '#cccccc';
		}
	}

	getLevelBadge(level) {
		const levelUpper = (level || 'INFO').toUpperCase();
		const color = this.getLevelColor(level);
		return `<span style="display: inline-block; padding: 2px 6px; border-radius: 3px; background: ${color}; color: #000; font-weight: bold; font-size: 10px; margin-right: 8px;">${levelUpper}</span>`;
	}

	parseJSONLog(line) {
		try {
			const log = JSON.parse(line);
			return {
				valid: true,
				time: log.time || '',
				level: log.level || 'INFO',
				msg: log.msg || '',
				data: log
			};
		} catch (e) {
			return { valid: false, raw: line };
		}
	}

	formatLogEntry(logEntry) {
		if (!logEntry.valid) {
			return `<div style="padding: 8px; border-bottom: 1px solid var(--background-modifier-border); opacity: 0.7;">${this.escapeHtml(logEntry.raw)}</div>`;
		}

		const { time, level, msg, data } = logEntry;
		const timestamp = time ? new Date(time).toLocaleString('zh-TW', { 
			year: 'numeric',
			month: '2-digit', 
			day: '2-digit',
			hour: '2-digit', 
			minute: '2-digit', 
			second: '2-digit'
		}) : '';

		// Extract additional fields
		const additionalFields = Object.keys(data)
			.filter(key => !['time', 'level', 'msg'].includes(key))
			.map(key => `<span style="color: #569cd6;">${this.escapeHtml(key)}</span>=<span style="color: #ce9178;">${this.escapeHtml(String(data[key]))}</span>`)
			.join(' ');

		return `
			<div style="padding: 8px; border-bottom: 1px solid var(--background-modifier-border); font-family: monospace; font-size: 12px;">
				<div style="margin-bottom: 4px;">
					${this.getLevelBadge(level)}
					<span style="color: var(--text-muted); font-size: 11px;">${timestamp}</span>
				</div>
				<div style="margin-left: 0px; color: var(--text-normal);">
					${this.escapeHtml(msg)}
					${additionalFields ? `<div style="margin-top: 4px; color: var(--text-muted); font-size: 11px;">${additionalFields}</div>` : ''}
				</div>
			</div>
		`;
	}

	escapeHtml(text) {
		const div = document.createElement('div');
		div.textContent = text;
		return div.innerHTML;
	}

	async findLogFiles() {
		const logFiles = [];
		
		// Search for .log files in home directory .novelmaker/log
		const homeDir = os.homedir();
		const searchPath = path.join(homeDir, '.novelmaker', 'log');

        try {
            if (!fs.existsSync(searchPath)) {
                return logFiles;
            }
            
            const files = fs.readdirSync(searchPath);
            for (const file of files) {
                if (file.endsWith('.log')) {
                    const fullPath = path.join(searchPath, file);
                    const relativePath = path.relative(this.vaultPath, fullPath);
                    const stats = fs.statSync(fullPath);
                    logFiles.push({
                        name: file,
                        path: fullPath,
                        relativePath: relativePath,
                        mtime: stats.mtime
                    });
                }
            }
        } catch (error) {
            console.error(`Error reading directory ${searchPath}:`, error);
        }

		// Sort by modification time (newest first)
		logFiles.sort((a, b) => b.mtime - a.mtime);

		return logFiles;
	}

	async onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '查看日誌' });

		// Find log files
		this.logFiles = await this.findLogFiles();

		if (this.logFiles.length === 0) {
			contentEl.createEl('p', { text: '未找到任何日誌檔案' });
			
			new Setting(contentEl)
				.addButton((btn) =>
					btn
						.setButtonText('關閉')
						.onClick(() => {
							this.close();
						})
				);
			return;
		}

		// Dropdown to select log file
		new Setting(contentEl)
			.setName('選擇日誌檔案')
			.setDesc(`找到 ${this.logFiles.length} 個日誌檔案`)
			.addDropdown((dropdown) => {
				// Add empty option
				dropdown.addOption('', '-- 選擇檔案 --');
				
				// Add all log files
				this.logFiles.forEach((logFile) => {
					dropdown.addOption(logFile.path, logFile.name);
				});
				
				// Auto-select the latest (first) log file
				if (this.logFiles.length > 0) {
					this.selectedLogFile = this.logFiles[0].path;
					dropdown.setValue(this.selectedLogFile);
				}
				
				dropdown.onChange((value) => {
					this.selectedLogFile = value;
					this.updateLogContent();
				});
				
				this.dropdownComponent = dropdown;
			});

		// Filter by level
		new Setting(contentEl)
			.setName('篩選層級')
			.addDropdown((dropdown) => {
				dropdown.addOption('all', '全部');
				dropdown.addOption('ERROR', 'ERROR');
				dropdown.addOption('WARN', 'WARN');
				dropdown.addOption('INFO', 'INFO');
				dropdown.addOption('DEBUG', 'DEBUG');
				
				dropdown.setValue(this.filterLevel);
				dropdown.onChange((value) => {
					this.filterLevel = value;
					this.updateLogContent();
				});
			});

		// Create container for log content
		const contentContainer = contentEl.createDiv({ cls: 'log-content-container' });
		contentContainer.style.maxHeight = '500px';
		contentContainer.style.overflow = 'auto';
		contentContainer.style.border = '1px solid var(--background-modifier-border)';
		contentContainer.style.marginTop = '10px';
		contentContainer.style.marginBottom = '10px';
		contentContainer.style.backgroundColor = 'var(--background-primary)';
		contentContainer.style.borderRadius = '4px';
		
		this.contentContainer = contentContainer;
		this.contentContainer.innerHTML = '<div style="padding: 20px; text-align: center; color: var(--text-muted);">請選擇一個日誌檔案</div>';

		// Load the latest log file content automatically
		if (this.logFiles.length > 0) {
			this.updateLogContent();
		}

		// Buttons
		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('重新整理')
					.onClick(async () => {
						// Refresh log files list
						this.logFiles = await this.findLogFiles();
						
						// Update dropdown
						if (this.dropdownComponent) {
							// Clear existing options
							this.dropdownComponent.selectEl.empty();
							
							// Re-add options
							const emptyOption = document.createElement('option');
							emptyOption.value = '';
							emptyOption.text = '-- 選擇檔案 --';
							this.dropdownComponent.selectEl.add(emptyOption);
							
							this.logFiles.forEach((logFile) => {
								const option = document.createElement('option');
								option.value = logFile.path;
								option.text = logFile.name;
								this.dropdownComponent.selectEl.add(option);
							});
						}
						
						// Clear content
						this.selectedLogFile = '';
						this.contentContainer.innerHTML = '<div style="padding: 20px; text-align: center; color: var(--text-muted);">請選擇一個日誌檔案</div>';
						
						new Notice(`找到 ${this.logFiles.length} 個日誌檔案`);
					})
			)
			.addButton((btn) =>
				btn
					.setButtonText('關閉')
					.onClick(() => {
						this.close();
					})
			);
	}

	updateLogContent() {
		if (!this.selectedLogFile) {
			this.contentContainer.innerHTML = '<div style="padding: 20px; text-align: center; color: var(--text-muted);">請選擇一個日誌檔案</div>';
			return;
		}

		try {
			const content = fs.readFileSync(this.selectedLogFile, 'utf8');
			
			if (content.trim() === '') {
				this.contentContainer.innerHTML = '<div style="padding: 20px; text-align: center; color: var(--text-muted);">(檔案為空)</div>';
				return;
			}

			// Parse JSON log lines
			const lines = content.trim().split('\n');
			const logEntries = lines.map(line => this.parseJSONLog(line));

			// Filter by level
			const filteredEntries = this.filterLevel === 'all' 
				? logEntries 
				: logEntries.filter(entry => entry.valid && entry.level.toUpperCase() === this.filterLevel);

			if (filteredEntries.length === 0) {
				this.contentContainer.innerHTML = '<div style="padding: 20px; text-align: center; color: var(--text-muted);">沒有符合篩選條件的日誌</div>';
				return;
			}

			// Format and display
			const html = filteredEntries.map(entry => this.formatLogEntry(entry)).join('');
			this.contentContainer.innerHTML = html;
			
			// Auto-scroll to bottom to show latest logs
			this.contentContainer.scrollTop = this.contentContainer.scrollHeight;
		} catch (error) {
			this.contentContainer.innerHTML = `<div style="padding: 20px; color: var(--text-error);">❌ 無法讀取檔案: ${this.escapeHtml(error.message)}</div>`;
			new Notice(`無法讀取日誌檔案: ${error.message}`);
		}
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

module.exports = LogsModal;
