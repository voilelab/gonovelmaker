const { Plugin, Modal, Notice, Setting, PluginSettingTab } = require('obsidian');
const { exec } = require('child_process');
const { promisify } = require('util');

const execAsync = promisify(exec);

// Default settings
const DEFAULT_SETTINGS = {
	cliPath: './novelmaker-obs',
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

class GenNextModal extends Modal {
	constructor(app, onSubmit) {
		super(app);
		this.title = '';
		this.prompt = '';
		this.prevCount = 3; // Default value
		this.onSubmit = onSubmit;

		this.prevSetting = null;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '生成下一章' });

		new Setting(contentEl)
			.setName('章節標題')
			.setDesc('下一章的標題（必填）')
			.addText((text) =>
				text
					.setPlaceholder('e.g., 第3章')
					.onChange((value) => {
						this.title = value;
					})
			);

		this.prevSetting = new Setting(contentEl)
			.setName('前幾章數量 (前3章)')
			.setDesc('要包含多少前面的章節作為上下文（預設：3，最大：10）')
			.addSlider((text) =>
				text
					.setValue(3)
					.setLimits(0, 10, 1)
					.onChange((num) => {
						if (isNaN(num)) {
							return;
						}

						this.prevCount = num;
						this.prevSetting.setName(`前幾章數量 (前${num}章)`);
					})
			);

		new Setting(contentEl)
			.setName('提示詞')
			.setDesc('章節生成的額外指示（選填）')
			.addTextArea((text) =>
				text
					.setPlaceholder('例如：著重於角色發展...')
					.onChange((value) => {
						this.prompt = value;
					})
			);

		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('取消')
					.onClick(() => {
						this.close();
					})
			)
			.addButton((btn) =>
				btn
					.setButtonText('生成')
					.setCta()
					.onClick(() => {
						if (!this.title || !this.title.trim()) {
							new Notice('❌ 請輸入章節標題');
							return;
						}
						this.close();
						this.onSubmit(this.title, this.prompt, this.prevCount);
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

class GenNextEmptyModal extends Modal {
	constructor(app, onSubmit) {
		super(app);
		this.title = '';
		this.prompt = '';
		this.onSubmit = onSubmit;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '生成空白下一章' });

		new Setting(contentEl)
			.setName('章節標題')
			.setDesc('下一章的標題（必填）')
			.addText((text) =>
				text
					.setPlaceholder('e.g., 第3章')
					.onChange((value) => {
						this.title = value;
					})
			);

		new Setting(contentEl)
			.setName('提示詞')
			.setDesc('章節的備註或提示（選填）')
			.addTextArea((text) =>
				text
					.setPlaceholder('例如：計劃寫一個戰鬥場景...')
					.onChange((value) => {
						this.prompt = value;
					})
			);

		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('取消')
					.onClick(() => {
						this.close();
					})
			)
			.addButton((btn) =>
				btn
					.setButtonText('建立')
					.setCta()
					.onClick(() => {
						if (!this.title || !this.title.trim()) {
							new Notice('❌ 請輸入章節標題');
							return;
						}
						this.close();
						this.onSubmit(this.title, this.prompt);
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

class GenCharModal extends Modal {
	constructor(app, onSubmit) {
		super(app);
		this.name = '';
		this.prompt = '';
		this.onSubmit = onSubmit;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '生成角色' });

		new Setting(contentEl)
			.setName('角色名稱')
			.setDesc('角色的名稱（必填）')
			.addText((text) =>
				text
					.setPlaceholder('e.g., 艾莉娜')
					.onChange((value) => {
						this.name = value;
					})
			);

		new Setting(contentEl)
			.setName('角色描述')
			.setDesc('角色的特徵或背景描述（選填）')
			.addTextArea((text) =>
				text
					.setPlaceholder('例如：一位年輕的魔法師，擅長元素魔法...')
					.onChange((value) => {
						this.prompt = value;
					})
			);

		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('取消')
					.onClick(() => {
						this.close();
					})
			)
			.addButton((btn) =>
				btn
					.setButtonText('生成')
					.setCta()
					.onClick(() => {
						if (!this.name || !this.name.trim()) {
							new Notice('❌ 請輸入角色名稱');
							return;
						}
						this.close();
						this.onSubmit(this.name, this.prompt);
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

class GenCharImgModal extends Modal {
	constructor(app, activeFile, onSubmit) {
		super(app);
		this.activeFile = activeFile;
		this.name = '';
		this.prompt = '';
		this.onSubmit = onSubmit;
	}

	async onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '生成角色圖片' });

		// Show current file info
		new Setting(contentEl)
			.setName('目標檔案')
			.setDesc('將為此角色生成圖片')
			.addText((text) => {
				text.setValue(this.activeFile.path);
				text.inputEl.disabled = true;
				text.inputEl.style.width = '100%';
			});

		// Try to extract character name from frontmatter
		try {
			const content = await this.app.vault.read(this.activeFile);
			const frontmatterMatch = content.match(/^---\n([\s\S]*?)\n---/);
			if (frontmatterMatch) {
				const frontmatter = frontmatterMatch[1];
				const nameMatch = frontmatter.match(/^name:\s*(.+)$/m);
				if (nameMatch) {
					this.name = nameMatch[1].trim();
				}
			}
		} catch (error) {
			console.error('Failed to read character name:', error);
		}

		new Setting(contentEl)
			.setName('角色名稱')
			.setDesc('角色的名稱（從檔案讀取）')
			.addText((text) => {
				text.setValue(this.name);
				text.inputEl.disabled = true;
				text.inputEl.style.width = '100%';
			});

		new Setting(contentEl)
			.setName('圖片描述')
			.setDesc('自訂圖片生成的提示詞（選填，留空則使用角色檔案中的描述）')
			.addTextArea((text) =>
				text
					.setPlaceholder('例如：一位年輕的魔法師肖像，背景是魔法學院...')
					.onChange((value) => {
						this.prompt = value;
					})
			);

		contentEl.createEl('p', { 
			text: '💡 提示：圖片將使用 DALL-E API 生成並儲存至 Character/ 目錄',
			cls: 'mod-info'
		});

		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('取消')
					.onClick(() => {
						this.close();
					})
			)
			.addButton((btn) =>
				btn
					.setButtonText('生成圖片')
					.setCta()
					.onClick(() => {
						if (!this.name || !this.name.trim()) {
							new Notice('❌ 無法取得角色名稱');
							return;
						}
						this.close();
						this.onSubmit(this.name, this.prompt);
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

class GenCurrModal extends Modal {
	constructor(app, activeFile, onSubmit) {
		super(app);
		this.activeFile = activeFile;
		this.prevCount = 3; // Default value
		this.onSubmit = onSubmit;
		this.prevSetting = null;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '重新生成章節' });

		// Show current file info
		new Setting(contentEl)
			.setName('目標檔案')
			.setDesc('將重新生成此檔案的內容')
			.addText((text) => {
				text.setValue(this.activeFile.path);
				text.inputEl.disabled = true;
				text.inputEl.style.width = '100%';
			});

		contentEl.createEl('p', { 
			text: '⚠️ 此操作將使用檔案中的 prompt 欄位重新生成章節內容，並覆蓋現有內容。',
			cls: 'mod-warning'
		});

		this.prevSetting = new Setting(contentEl)
			.setName('前幾章數量 (前3章)')
			.setDesc('要包含多少前面的章節作為上下文（預設：3，最大：10）')
			.addSlider((slider) =>
				slider
					.setValue(3)
					.setLimits(0, 10, 1)
					.onChange((num) => {
						if (isNaN(num)) {
							return;
						}
						this.prevCount = num;
						this.prevSetting.setName(`前幾章數量 (前${num}章)`);
					})
			);

		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('取消')
					.onClick(() => {
						this.close();
					})
			)
			.addButton((btn) =>
				btn
					.setButtonText('重新生成')
					.setCta()
					.onClick(() => {
						this.close();
						this.onSubmit(this.activeFile.path, this.prevCount);
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

class GenCharCurrModal extends Modal {
	constructor(app, activeFile, onSubmit) {
		super(app);
		this.activeFile = activeFile;
		this.onSubmit = onSubmit;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '重新生成角色' });

		// Show current file info
		new Setting(contentEl)
			.setName('目標檔案')
			.setDesc('將重新生成此檔案的內容')
			.addText((text) => {
				text.setValue(this.activeFile.path);
				text.inputEl.disabled = true;
				text.inputEl.style.width = '100%';
			});

		contentEl.createEl('p', { 
			text: '⚠️ 此操作將使用檔案中的 prompt 欄位重新生成角色資料，並覆蓋現有內容。',
			cls: 'mod-warning'
		});

		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('取消')
					.onClick(() => {
						this.close();
					})
			)
			.addButton((btn) =>
				btn
					.setButtonText('重新生成')
					.setCta()
					.onClick(() => {
						this.close();
						this.onSubmit(this.activeFile.path);
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

class ExportModal extends Modal {
	constructor(app, onSubmit) {
		super(app);
		this.outputPath = '';
		this.onSubmit = onSubmit;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '匯出小說' });

		const pathSetting = new Setting(contentEl)
			.setName('輸出檔案路徑')
			.setDesc('選擇小說匯出的檔案路徑')
			.addText((text) => {
				text
					.setPlaceholder('選擇檔案位置...')
					.setValue(this.outputPath)
					.onChange((value) => {
						this.outputPath = value;
					});
				text.inputEl.style.width = '100%';
			})
			.addButton((btn) =>
				btn
					.setButtonText('瀏覽...')
					.onClick(async () => {
						try {
							// Use Electron's dialog
							const { dialog } = require('electron').remote;
							const result = await dialog.showSaveDialog({
								title: '選擇匯出位置',
								defaultPath: 'novel.txt',
								filters: [
									{ name: '文字檔案', extensions: ['txt'] },
									{ name: '所有檔案', extensions: ['*'] }
								],
								properties: ['createDirectory', 'showOverwriteConfirmation']
							});

							if (!result.canceled && result.filePath) {
								this.outputPath = result.filePath;
								// Update the text input with selected path
								const textInput = pathSetting.controlEl.querySelector('input[type="text"]');
								if (textInput) {
									textInput.value = this.outputPath;
								}
							}
						} catch (error) {
							// Fallback if electron.remote is not available
							console.error('File dialog error:', error);
							new Notice('⚠ 無法開啟檔案對話框，請手動輸入路徑');
						}
					})
			);

		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('取消')
					.onClick(() => {
						this.close();
					})
			)
			.addButton((btn) =>
				btn
					.setButtonText('匯出')
					.setCta()
					.onClick(() => {
						if (!this.outputPath || !this.outputPath.trim()) {
							new Notice('❌ 請選擇或輸入輸出檔案路徑');
							return;
						}
						this.close();
						this.onSubmit(this.outputPath);
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

class LoadingModal extends Modal {
	constructor(app, message = "處理中…請稍候") {
		super(app);
		this.message = message;
	}

	onOpen() {
		// remove X
		const closeBtn = this.containerEl.querySelector(".modal-close-button");
		if (closeBtn) closeBtn.remove();

		// block ESC
		this.scope.register([], 'Escape', evt => {
			evt.preventDefault();
		});

		// block clicking background
		this.containerEl.onclick = (evt) => {
			evt.stopPropagation();
		};

		// UI
		this.contentEl.empty();
		this.setTitle(this.message);
		this.contentEl.createEl("div", { text: "⏳ 生成中…" });
	}

	// prevent close
	close() {
		if (this.forceClose) {
			super.close();
		}
	}

	// safe manual close
	forceCloseNow() {
		this.forceClose = true;
		this.close();
	}
}

class ResultModal extends Modal {
	constructor(app, title, output) {
		super(app);
		this.modalTitle = title;
		this.output = output;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: this.modalTitle });

		// Success message
		const successDiv = contentEl.createDiv({ cls: 'novelmaker-result-success' });
		successDiv.createEl('span', { text: '✅ ', cls: 'novelmaker-success-icon' });
		successDiv.createEl('span', { text: '生成成功！' });

		// File path if available
		if (this.output.filepath) {
			const fileDiv = contentEl.createDiv({ cls: 'novelmaker-result-info' });
			fileDiv.createEl('strong', { text: '檔案：' });
			fileDiv.createEl('span', { text: this.output.filepath });
		}

		// Token usage section
		if (this.output.input_tokens !== undefined || this.output.output_tokens !== undefined) {
			const usageDiv = contentEl.createDiv({ cls: 'novelmaker-token-usage' });
			usageDiv.createEl('h3', { text: '🔢 Token 使用量' });
			
			const usageGrid = usageDiv.createDiv({ cls: 'novelmaker-token-grid' });
			
			if (this.output.input_tokens !== undefined) {
				const inputRow = usageGrid.createDiv({ cls: 'novelmaker-token-row' });
				inputRow.createEl('span', { text: 'Input tokens:', cls: 'novelmaker-token-label' });
				inputRow.createEl('span', { text: this.output.input_tokens.toLocaleString(), cls: 'novelmaker-token-value' });
			}
			
			if (this.output.output_tokens !== undefined) {
				const outputRow = usageGrid.createDiv({ cls: 'novelmaker-token-row' });
				outputRow.createEl('span', { text: 'Output tokens:', cls: 'novelmaker-token-label' });
				outputRow.createEl('span', { text: this.output.output_tokens.toLocaleString(), cls: 'novelmaker-token-value' });
			}
			
			if (this.output.total_tokens !== undefined) {
				const totalRow = usageGrid.createDiv({ cls: 'novelmaker-token-row novelmaker-token-total' });
				totalRow.createEl('span', { text: 'Total tokens:', cls: 'novelmaker-token-label' });
				totalRow.createEl('span', { text: this.output.total_tokens.toLocaleString(), cls: 'novelmaker-token-value' });
			}
		}

		// Close button
		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('確定')
					.setCta()
					.onClick(() => {
						this.close();
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

class BackendModal extends Modal {
	constructor(app, plugin, backend = null, onSubmit) {
		super(app);
		this.plugin = plugin;
		this.backend = backend; // null for add, object for edit
		this.name = backend?.name || '';
		this.base_url = backend?.base_url || '';
		this.api_key = backend?.api_key || '';
		this.model = backend?.model || '';
		this.timeout = backend?.timeout || 60;
		this.onSubmit = onSubmit;
		this.apiKeyModified = false; // Track if user actually modified the API key
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: this.backend ? '編輯 Backend' : '新增 Backend' });

		new Setting(contentEl)
			.setName('Backend 名稱')
			.setDesc('Backend 的識別名稱（必填）')
			.addText((text) => {
				text
					.setPlaceholder('例如：openrouter, openai')
					.setValue(this.name)
					.onChange((value) => {
						this.name = value;
					});
				text.inputEl.style.width = '100%';
				// Disable name editing when editing existing backend
				if (this.backend) {
					text.inputEl.disabled = true;
				}
			});

		new Setting(contentEl)
			.setName('Base URL')
			.setDesc('API 的基礎 URL（必填）')
			.addText((text) => {
				text
					.setPlaceholder('例如：https://openrouter.ai/api/v1')
					.setValue(this.base_url)
					.onChange((value) => {
						this.base_url = value;
					});
				text.inputEl.style.width = '100%';
			});

		new Setting(contentEl)
			.setName('API Key')
			.setDesc(this.backend ? 'API 金鑰（留空表示不修改）' : 'API 金鑰（必填）')
			.addText((text) => {
				text
					.setPlaceholder(this.backend ? '留空表示保持原值...' : 'sk-...')
					.setValue(this.backend ? '' : this.api_key) // In edit mode, show empty by default
					.onChange((value) => {
						this.api_key = value;
						// Mark as modified if user enters any value
						if (this.backend && value.trim()) {
							this.apiKeyModified = true;
						}
					});
				text.inputEl.style.width = '100%';
				text.inputEl.type = 'password';
			});

		new Setting(contentEl)
			.setName('預設模型')
			.setDesc('此 backend 的預設模型名稱（選填）')
			.addText((text) => {
				text
					.setPlaceholder('例如：gpt-4, claude-3-opus-20240229')
					.setValue(this.model)
					.onChange((value) => {
						this.model = value;
					});
				text.inputEl.style.width = '100%';
			});

		new Setting(contentEl)
			.setName(`API 請求超時 (秒) (${this.timeout} 秒)`)
			.setDesc('此 backend 的超時時間（秒）（預設：60 秒）')
			.addSlider((slider) => {
				slider
					.setLimits(10, 300, 10)
					.setValue(this.timeout)
					.onChange((value) => {
						if (isNaN(value)) {
							return;
						}
						this.timeout = value;
						// Update the setting name to show current value
						const settingEl = contentEl.querySelector('.setting-item:last-of-type .setting-item-name');
						if (settingEl) {
							settingEl.textContent = `API 請求超時 (秒) (${value} 秒)`;
						}
					});
			});

		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('取消')
					.onClick(() => {
						this.close();
					})
			)
			.addButton((btn) =>
				btn
					.setButtonText(this.backend ? '更新' : '新增')
					.setCta()
					.onClick(() => {
						if (!this.name || !this.name.trim()) {
							new Notice('❌ 請輸入 Backend 名稱');
							return;
						}
						if (!this.base_url || !this.base_url.trim()) {
							new Notice('❌ 請輸入 Base URL');
							return;
						}
						// Only require API key for new backends
						if (!this.backend && (!this.api_key || !this.api_key.trim())) {
							new Notice('❌ 請輸入 API Key');
							return;
						}
						this.close();
						this.onSubmit({
							name: this.name,
							base_url: this.base_url,
							api_key: this.api_key,
							model: this.model,
							timeout: this.timeout,
							apiKeyModified: this.apiKeyModified // Pass the flag
						});
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
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
