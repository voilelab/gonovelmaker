const { Plugin, Notice, Setting, PluginSettingTab } = require('obsidian');
const { exec } = require('child_process');
const { promisify } = require('util');

const execAsync = promisify(exec);

// Import modals
const LoadingModal = require('./src/modals/loading');
const ResultModal = require('./src/modals/result');
const GenNextModal = require('./src/modals/gen-next');
const GenNextEmptyModal = require('./src/modals/gen-next-empty');
const GenCurrModal = require('./src/modals/gen-curr');
const GenCharModal = require('./src/modals/gen-char');
const GenCharCurrModal = require('./src/modals/gen-char-curr');
const GenCharImgModal = require('./src/modals/gen-char-img');
const ExportModal = require('./src/modals/export');
const BackendModal = require('./src/modals/backend');

// Default settings
const DEFAULT_SETTINGS = {
	cliPath: 'novelmaker-obs',
	backend: '',
	openAfterGen: false,
	openAfterGenMs: 500,
};

class NovelMakerPlugin extends Plugin {
	async onload() {
		// Load settings
		await this.loadSettings();
		
		// Add settings tab
		this.addSettingTab(new NovelMakerSettingTab(this.app, this));
		// Register gen-next command
		this.addCommand({
			id: 'gen-next-chapter',
			name: '生成下一章',
			callback: () => {
				new GenNextModal(this.app, async (title, prompt, prevCount) => {
					const loadingModal = new LoadingModal(this.app, '正在生成下一章...請稍候');
					loadingModal.open();

					try {
						// Get vault path
						const vaultPath = this.app.vault.adapter.basePath;
						
						// Build the command using configured CLI path
						let cmd = `${this.settings.cliPath} gen-next --json --title "${title}"`;
						if (prompt && prompt.trim()) {
							cmd += ` --prompt "${prompt}"`;
						}
						if (prevCount !== null && prevCount !== undefined) {
							cmd += ` --prev-chapters ${prevCount}`;
						}
						if (this.settings.backend && this.settings.backend.trim()) {
							cmd += ` --backend "${this.settings.backend}"`;
						}
						
						// Call CLI with the input, setting cwd to vault path
						const { stdout, stderr } = await execAsync(cmd, { cwd: vaultPath });

						// Parse JSON output
						let output;
						try {
							output = JSON.parse(stdout);
						} catch (jsonError) {
							throw new Error('無法解析生成結果的 JSON 輸出');
						}

						if (this.settings.openAfterGen) {
							if (output && output.filepath) {
								setTimeout(() => {
									const file = this.app.vault.getAbstractFileByPath(output.filepath);
									if (file) {
										this.app.workspace.getLeaf().openFile(file);
									} else {
										new Notice(`⚠ 無法在 Vault 中找到生成的檔案: ${output.filepath}`);
									}
								}, this.settings.openAfterGenMs);
							} else {
								new Notice('⚠ 生成結果中缺少 filepath 資訊');
							}
						}
						
						// Show result modal with token usage
						new ResultModal(this.app, '章節生成完成', output).open();
						
						// Optionally show output in console
						if (stdout) console.log(stdout);
						if (stderr) console.error(stderr);
					} catch (error) {
						new Notice(`❌ 錯誤: ${error.message}`);
						console.error(error);
					} finally {
						loadingModal.forceCloseNow();
					}
				}).open();
			}
		});

		// Register gen-next-empty command
		this.addCommand({
			id: 'gen-next-empty-chapter',
			name: '生成空白下一章',
			callback: () => {
				new GenNextEmptyModal(this.app, async (title, prompt) => {
					const loadingModal = new LoadingModal(this.app, '正在建立空白章節...請稍候');
					loadingModal.open();

					try {
						// Get vault path
						const vaultPath = this.app.vault.adapter.basePath;
						
						// Build the command using configured CLI path
						let cmd = `${this.settings.cliPath} gen-next-empty --json --title "${title}"`;
						if (prompt && prompt.trim()) {
							cmd += ` --prompt "${prompt}"`;
						}
						
						// Call CLI with the input, setting cwd to vault path
						const { stdout, stderr } = await execAsync(cmd, { cwd: vaultPath });

						if (this.settings.openAfterGen) {
							// Parse JSON output to get filepath
							let output;
							try {
								output = JSON.parse(stdout);
							} catch (jsonError) {
								throw new Error('無法解析生成結果的 JSON 輸出');
							}

							if (output && output.filepath) {
								setTimeout(() => {
									const file = this.app.vault.getAbstractFileByPath(output.filepath);
									if (file) {
										this.app.workspace.getLeaf().openFile(file);
									} else {
										new Notice(`⚠ 無法在 Vault 中找到生成的檔案: ${output.filepath}`);
									}
								}, this.settings.openAfterGenMs);
							} else {
								new Notice('⚠ 生成結果中缺少 filepath 資訊');
							}
						}
						
						// Show success notification
						new Notice('✅ 空白章節建立成功！');
						
						// Optionally show output in console
						if (stdout) console.log(stdout);
						if (stderr) console.error(stderr);
					} catch (error) {
						new Notice(`❌ 錯誤: ${error.message}`);
						console.error(error);
					} finally {
						loadingModal.forceCloseNow();
					}
				}).open();
			}
		});

		// Register gen-char command
		this.addCommand({
			id: 'gen-character',
			name: '生成角色',
			callback: () => {
				new GenCharModal(this.app, async (name, prompt) => {
					const loadingModal = new LoadingModal(this.app, '正在生成角色...請稍候');
					loadingModal.open();

					try {
						// Get vault path
						const vaultPath = this.app.vault.adapter.basePath;
						
						// Build the command using configured CLI path
						let cmd = `${this.settings.cliPath} gen-char --json --name "${name}"`;
						if (prompt && prompt.trim()) {
							cmd += ` --prompt "${prompt}"`;
						}
						if (this.settings.backend && this.settings.backend.trim()) {
							cmd += ` --backend "${this.settings.backend}"`;
						}
						
						// Call CLI with the input, setting cwd to vault path
						const { stdout, stderr } = await execAsync(cmd, { cwd: vaultPath });

						// Parse JSON output
						let output;
						try {
							output = JSON.parse(stdout);
						} catch (jsonError) {
							throw new Error('無法解析生成結果的 JSON 輸出');
						}

						if (this.settings.openAfterGen) {
							if (output && output.filepath) {
								setTimeout(() => {
									const file = this.app.vault.getAbstractFileByPath(output.filepath);
									if (file) {
										this.app.workspace.getLeaf().openFile(file);
									} else {
										new Notice(`⚠ 無法在 Vault 中找到生成的檔案: ${output.filepath}`);
									}
								}, this.settings.openAfterGenMs);
							} else {
								new Notice('⚠ 生成結果中缺少 filepath 資訊');
							}
						}

						// Show result modal with token usage
						new ResultModal(this.app, '角色生成完成', output).open();
						
						// Optionally show output in console
						if (stdout) console.log(stdout);
						if (stderr) console.error(stderr);
					} catch (error) {
						new Notice(`❌ 錯誤: ${error.message}`);
						console.error(error);
					} finally {
						loadingModal.forceCloseNow();
					}
				}).open();
			}
		});

		// Register gen-curr command
		this.addCommand({
			id: 'gen-curr-chapter',
			name: '重新生成當前章節',
			checkCallback: (checking) => {
				const activeFile = this.app.workspace.getActiveFile();
				if (activeFile && activeFile.path.startsWith('Story/') && activeFile.extension === 'md') {
					if (!checking) {
						new GenCurrModal(this.app, activeFile, async (filepath, prevCount) => {
							const loadingModal = new LoadingModal(this.app, '正在重新生成章節...請稍候');
							loadingModal.open();

							try {
								// Get vault path
								const vaultPath = this.app.vault.adapter.basePath;
								
								// Build the command using configured CLI path
								let cmd = `${this.settings.cliPath} gen-curr --json --filepath "${filepath}"`;
								if (prevCount !== null && prevCount !== undefined) {
									cmd += ` --prev-chapters ${prevCount}`;
								}
								if (this.settings.backend && this.settings.backend.trim()) {
									cmd += ` --backend "${this.settings.backend}"`;
								}
								if (this.settings.timeout && !isNaN(this.settings.timeout)) {
									cmd += ` --timeout ${this.settings.timeout}`;
								}
								
								// Call CLI with the input, setting cwd to vault path
								const { stdout, stderr } = await execAsync(cmd, { cwd: vaultPath });
								
								// Parse JSON output and show result modal
								try {
									const output = JSON.parse(stdout);
									new ResultModal(this.app, '章節重新生成完成', output).open();
								} catch (jsonError) {
									// Fallback to simple notice if JSON parsing fails
									new Notice('✅ 章節重新生成成功！');
								}
								
								// Optionally show output in console
								if (stdout) console.log(stdout);
								if (stderr) console.error(stderr);
							} catch (error) {
								new Notice(`❌ 錯誤: ${error.message}`);
								console.error(error);
							} finally {
								loadingModal.forceCloseNow();
							}
						}).open();
					}
					return true;
				}
				return false;
			}
		});

		// Register gen-char-curr command
		this.addCommand({
			id: 'gen-char-curr',
			name: '重新生成當前角色',
			checkCallback: (checking) => {
				const activeFile = this.app.workspace.getActiveFile();
				if (activeFile && activeFile.path.startsWith('Character/') && activeFile.extension === 'md') {
					if (!checking) {
						new GenCharCurrModal(this.app, activeFile, async (filepath) => {
							const loadingModal = new LoadingModal(this.app, '正在重新生成角色...請稍候');
							loadingModal.open();

							try {
								// Get vault path
								const vaultPath = this.app.vault.adapter.basePath;
								
								// Build the command using configured CLI path
								let cmd = `${this.settings.cliPath} gen-char-curr --json --filepath "${filepath}"`;
								if (this.settings.backend && this.settings.backend.trim()) {
									cmd += ` --backend "${this.settings.backend}"`;
								}
								
								// Call CLI with the input, setting cwd to vault path
								const { stdout, stderr } = await execAsync(cmd, { cwd: vaultPath });
								
								// Parse JSON output and show result modal
								try {
									const output = JSON.parse(stdout);
									new ResultModal(this.app, '角色重新生成完成', output).open();
								} catch (jsonError) {
									// Fallback to simple notice if JSON parsing fails
									new Notice('✅ 角色重新生成成功！');
								}
								
								// Optionally show output in console
								if (stdout) console.log(stdout);
								if (stderr) console.error(stderr);
							} catch (error) {
								new Notice(`❌ 錯誤: ${error.message}`);
								console.error(error);
							} finally {
								loadingModal.forceCloseNow();
							}
						}).open();
					}
					return true;
				}
				return false;
			}
		});

		// Register gen-char-img command
		this.addCommand({
			id: 'gen-char-img',
			name: '生成角色圖片',
			checkCallback: (checking) => {
				const activeFile = this.app.workspace.getActiveFile();
				if (activeFile && activeFile.path.startsWith('Character/') && activeFile.extension === 'md') {
					if (!checking) {
						new GenCharImgModal(this.app, activeFile, async (name, prompt) => {
							const loadingModal = new LoadingModal(this.app, '正在生成角色圖片...請稍候');
							loadingModal.open();

							try {
								// Get vault path
								const vaultPath = this.app.vault.adapter.basePath;
								
								// Build the command using configured CLI path
								let cmd = `${this.settings.cliPath} gen-char-img --json --name "${name}"`;
								if (prompt && prompt.trim()) {
									cmd += ` --prompt "${prompt}"`;
								}
								if (this.settings.backend && this.settings.backend.trim()) {
									cmd += ` --backend "${this.settings.backend}"`;
								}
								
								// Call CLI with the input, setting cwd to vault path
								const { stdout, stderr } = await execAsync(cmd, { cwd: vaultPath });

								// Parse JSON output to get filepath
								let output;
								try {
									output = JSON.parse(stdout);
								} catch (jsonError) {
									throw new Error('無法解析生成結果的 JSON 輸出');
								}

								// Show success notification with image path
								if (output && output.filepath) {
									new Notice(`✅ 角色圖片生成成功！\n檔案：${output.filepath}`);
								} else {
									new Notice('✅ 角色圖片生成成功！');
								}
								
								// Optionally show output in console
								if (stdout) console.log(stdout);
								if (stderr) console.error(stderr);
							} catch (error) {
								new Notice(`❌ 錯誤: ${error.message}`);
								console.error(error);
							} finally {
								loadingModal.forceCloseNow();
							}
						}).open();
					}
					return true;
				}
				return false;
			}
		});

		// Register export command
		this.addCommand({
			id: 'export-novel',
			name: '匯出小說',
			callback: () => {
				new ExportModal(this.app, async (outputPath) => {
					const loadingModal = new LoadingModal(this.app, '正在匯出小說...請稍候');
					loadingModal.open();

					try {
						// Get vault path
						const vaultPath = this.app.vault.adapter.basePath;
						
						// Build the command using configured CLI path
						let cmd = `${this.settings.cliPath} export --output "${outputPath}" --type txt`;
						
						// Call CLI with the input, setting cwd to vault path
						const { stdout, stderr } = await execAsync(cmd, { cwd: vaultPath });

						// Show success notification
						new Notice('✅ 小說匯出成功！');
						
						// Optionally show output in console
						if (stdout) console.log(stdout);
						if (stderr) console.error(stderr);
					} catch (error) {
						new Notice(`❌ 錯誤: ${error.message}`);
						console.error(error);
					} finally {
						loadingModal.forceCloseNow();
					}
				}).open();
			}
		});
	}

	onunload() {
		// Cleanup if needed
	}
	
	async loadSettings() {
		this.settings = Object.assign({}, DEFAULT_SETTINGS, await this.loadData());
	}

	async saveSettings() {
		await this.saveData(this.settings);
	}
}

class NovelMakerSettingTab extends PluginSettingTab {
	constructor(app, plugin) {
		super(app, plugin);
		this.plugin = plugin;

		this.openAfterGenMsSetting = null;
	}

	async display() {
		const { containerEl } = this;

		containerEl.empty();

		containerEl.createEl('h2', { text: 'Novel Maker 設定' });

		new Setting(containerEl)
			.setName('CLI 指令路徑')
			.setDesc('設定 novelmaker-obs 指令的路徑。可使用相對路徑（如 ./novelmaker-obs）、絕對路徑（如 /usr/local/bin/novelmaker-obs）或指令名稱（如 novelmaker-obs，需在 PATH 中）')
			.addText((text) =>
				text
					.setPlaceholder('./novelmaker-obs')
					.setValue(this.plugin.settings.cliPath)
					.onChange(async (value) => {
						this.plugin.settings.cliPath = value || './novelmaker-obs';
						await this.plugin.saveSettings();
					})
			);

		// Backend Management Section
		containerEl.createEl('h3', { text: 'Backend 管理' });

		// Load available backends
		let backends = [];
		let backendList = [];
		let defaultBackend = '';
		try {
			const vaultPath = this.app.vault.adapter.basePath;
			const { stdout } = await execAsync(`${this.plugin.settings.cliPath} backend list --json`, { cwd: vaultPath });
			backendList = JSON.parse(stdout);
			if (Array.isArray(backendList)) {
				backends = backendList.map(b => b.name);
				const defaultBackendObj = backendList.find(b => b.is_default);
				if (defaultBackendObj) {
					defaultBackend = defaultBackendObj.name;
				}
			}
		} catch (error) {
			console.error('Failed to load backends:', error);
			new Notice('⚠ 無法載入 backend 列表，請確認 CLI 路徑是否正確');
		}

		// Add Backend button
		new Setting(containerEl)
			.setName('管理 Backends')
			.setDesc('新增、編輯或刪除 AI backend 配置')
			.addButton((btn) =>
				btn
					.setButtonText('新增 Backend')
					.setCta()
					.onClick(() => {
						new BackendModal(this.app, this.plugin, null, async (data) => {
							const loadingModal = new LoadingModal(this.app, '正在新增 Backend...');
							loadingModal.open();

							try {
								const vaultPath = this.app.vault.adapter.basePath;
								let cmd = `${this.plugin.settings.cliPath} backend add "${data.name}" --base_url "${data.base_url}" --api_key "${data.api_key}"`;
								if (data.model && data.model.trim()) {
									cmd += ` --model "${data.model}"`;
								}
								if (data.timeout) {
									cmd += ` --timeout ${data.timeout}`;
								}
								await execAsync(cmd, { cwd: vaultPath });
								new Notice(`\u2705 Backend "${data.name}" \u65b0\u589e\u6210\u529f\uff01`);
								// Refresh the settings page
								this.display();
							} catch (error) {
								new Notice(`\u274c \u65b0\u589e Backend \u5931\u6557: ${error.message}`);
								console.error(error);
							} finally {
								loadingModal.forceCloseNow();
							}
						}).open();
					})
			);

		// Display existing backends
		if (backendList.length > 0) {
			containerEl.createEl('h4', { text: '已設定的 Backends' });
			
			for (const backend of backendList) {
				const backendDiv = containerEl.createDiv({ cls: 'novelmaker-backend-item' });
				
				const backendInfo = backendDiv.createDiv({ cls: 'novelmaker-backend-info' });
				const nameEl = backendInfo.createEl('strong', { text: backend.name });
				if (backend.is_default) {
					nameEl.createEl('span', { text: ' (預設)', cls: 'novelmaker-backend-default' });
				}
				backendInfo.createEl('br');
				backendInfo.createEl('small', { text: backend.base_url });
				
				const backendActions = backendDiv.createDiv({ cls: 'novelmaker-backend-actions' });
				
				// Set as default button
				if (!backend.is_default) {
					const useBtn = backendActions.createEl('button', { text: '設為預設' });
					useBtn.onclick = async () => {
						try {
							const vaultPath = this.app.vault.adapter.basePath;
							await execAsync(`${this.plugin.settings.cliPath} backend use "${backend.name}"`, { cwd: vaultPath });
							new Notice(`✅ 已將 "${backend.name}" 設為預設 backend`);
							this.display();
						} catch (error) {
							new Notice(`❌ 設定預設 backend 失敗: ${error.message}`);
							console.error(error);
						}
					};
				}
				
				// Check button
				const checkBtn = backendActions.createEl('button', { text: '檢查' });
				checkBtn.onclick = async () => {
					const loadingModal = new LoadingModal(this.app, '正在檢查 Backend...');
					loadingModal.open();
					try {
						const vaultPath = this.app.vault.adapter.basePath;
						await execAsync(`${this.plugin.settings.cliPath} backend check "${backend.name}"`, { cwd: vaultPath });
						new Notice(`✅ Backend "${backend.name}" 連線正常！`);
					} catch (error) {
						new Notice(`❌ Backend "${backend.name}" 連線失敗: ${error.message}`);
						console.error(error);
					} finally {
						loadingModal.forceCloseNow();
					}
				};
				
				// Edit button
				const editBtn = backendActions.createEl('button', { text: '編輯' });
				editBtn.onclick = () => {
					new BackendModal(this.app, this.plugin, backend, async (data) => {
						const loadingModal = new LoadingModal(this.app, '正在更新 Backend...');
						loadingModal.open();

						try {
							const vaultPath = this.app.vault.adapter.basePath;
							// Update backend using add command (it handles both add and update)
							let cmd = `${this.plugin.settings.cliPath} backend add "${data.name}" --base_url "${data.base_url}"`;
							// Only include api_key if user actually modified it
							if (data.apiKeyModified && data.api_key && data.api_key.trim()) {
								cmd += ` --api_key "${data.api_key}"`;
							}
							if (data.model && data.model.trim()) {
								cmd += ` --model "${data.model}"`;
							}
							if (data.timeout) {
								cmd += ` --timeout ${data.timeout}`;
							}
							await execAsync(cmd, { cwd: vaultPath });
							new Notice(`\u2705 Backend "${data.name}" \u66f4\u65b0\u6210\u529f\uff01`);
							this.display();
						} catch (error) {
							new Notice(`\u274c \u66f4\u65b0 Backend \u5931\u6557: ${error.message}`);
							console.error(error);
						} finally {
							loadingModal.forceCloseNow();
						}
					}).open();
				};
				
				// Delete button
				const deleteBtn = backendActions.createEl('button', { text: '刪除', cls: 'mod-warning' });
				deleteBtn.onclick = async () => {
					if (!confirm(`確定要刪除 backend "${backend.name}" 嗎？`)) {
						return;
					}
					try {
						const vaultPath = this.app.vault.adapter.basePath;
						await execAsync(`${this.plugin.settings.cliPath} backend remove "${backend.name}"`, { cwd: vaultPath });
						new Notice(`✅ Backend "${backend.name}" 已刪除`);
						this.display();
					} catch (error) {
						new Notice(`❌ 刪除 Backend 失敗: ${error.message}`);
						console.error(error);
					}
				};
			}
		} else {
			containerEl.createEl('p', { 
				text: '尚未設定任何 backend，請點擊上方按鈕新增。',
				cls: 'mod-info'
			});
		}

		containerEl.createEl('h3', { text: '生成設定' });

		const backendSetting = new Setting(containerEl)
			.setName('Backend 名稱')
			.setDesc(`LLM Backend 的名稱。留空則使用設定檔中的預設 backend${defaultBackend ? ` (${defaultBackend})` : ''}。`);

		if (backends.length > 0) {
			backendSetting.addDropdown((dropdown) => {
				// Add empty option for using default
				dropdown.addOption('', `使用預設${defaultBackend ? ` (${defaultBackend})` : ''}`);
				
				// Add all available backends
				backends.forEach(name => {
					dropdown.addOption(name, name);
				});
				
				dropdown
					.setValue(this.plugin.settings.backend || '')
					.onChange(async (value) => {
						this.plugin.settings.backend = value;
						await this.plugin.saveSettings();
					});
			});
		} else {
			// Fallback to text input if we can't load backends
			backendSetting.addText((text) =>
				text
					.setPlaceholder('openrouter')
					.setValue(this.plugin.settings.backend)
					.onChange(async (value) => {
						this.plugin.settings.backend = value;
						await this.plugin.saveSettings();
					})
			);
		}

		// Info message for backend management
		const infoEl = containerEl.createDiv({ cls: 'novelmaker-api-warning' });
		infoEl.createEl('span', { text: 'ℹ️ 提示： ', cls: 'novelmaker-warning-icon' });
		infoEl.createEl('span', { 
			text: 'Backend 配置（包括 API Key）儲存在 ~/.novelmaker/config.toml 中，不會同步到 vault。每個 backend 可以設定自己的預設模型和超時時間。',
			cls: 'novelmaker-warning-text'
		});

		new Setting(containerEl)
			.setName('生成後自動打開檔案')
			.setDesc('生成章節或角色後，自動在 Obsidian 中打開生成的檔案')
			.addToggle((toggle) =>
				toggle
					.setValue(this.plugin.settings.openAfterGen)
					.onChange(async (value) => {
						this.plugin.settings.openAfterGen = value;
						await this.plugin.saveSettings();
					})
			);
		
		this.openAfterGenMsSetting = new Setting(containerEl)
			.setName(`生成後打開檔案延遲時間 (${this.plugin.settings.openAfterGenMs} 毫秒)`)
			.setDesc('生成章節或角色後，自動打開檔案前的延遲時間（毫秒）（預設：500 毫秒）')
			.addSlider((slider) =>
				slider
					.setLimits(0, 10000, 500)
					.setValue(this.plugin.settings.openAfterGenMs)
					.onChange(async (value) => {
						if (isNaN(value)) {
							return;
						}
						this.plugin.settings.openAfterGenMs = value;
						this.openAfterGenMsSetting.setName(`生成後打開檔案延遲時間 (${value} 毫秒)`);
						await this.plugin.saveSettings();
					})
			);
		
	}
}

module.exports = NovelMakerPlugin;
