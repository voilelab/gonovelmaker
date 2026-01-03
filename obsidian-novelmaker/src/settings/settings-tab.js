const { Notice, Setting, PluginSettingTab } = require('obsidian');
const LoadingModal = require('../modals/loading');
const BackendModal = require('../modals/backend');
const { execAsync } = require('../utils/cli');

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

		this.renderCLIPathSetting(containerEl);
		await this.renderBackendManagement(containerEl);
		this.renderGenerationSettings(containerEl);
	}

	renderCLIPathSetting(containerEl) {
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
	}

	async renderBackendManagement(containerEl) {
		containerEl.createEl('h3', { text: 'Backend 管理' });

		// Load available backends
		const { backends, backendList, defaultBackend } = await this.loadBackends();

		this.renderAddBackendButton(containerEl);
		this.renderBackendList(containerEl, backendList);
		this.renderBackendSelection(containerEl, backends, defaultBackend);
		this.renderBackendInfo(containerEl);
	}

	async loadBackends() {
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

		return { backends, backendList, defaultBackend };
	}

	renderAddBackendButton(containerEl) {
		new Setting(containerEl)
			.setName('管理 Backends')
			.setDesc('新增、編輯或刪除 AI backend 配置')
			.addButton((btn) =>
				btn
					.setButtonText('新增 Backend')
					.setCta()
					.onClick(() => this.handleAddBackend())
			);
	}

	async handleAddBackend() {
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
				new Notice(`✅ Backend "${data.name}" 新增成功！`);
				this.display();
			} catch (error) {
				new Notice(`❌ 新增 Backend 失敗: ${error.message}`);
				console.error(error);
			} finally {
				loadingModal.forceCloseNow();
			}
		}).open();
	}

	renderBackendList(containerEl, backendList) {
		if (backendList.length > 0) {
			containerEl.createEl('h4', { text: '已設定的 Backends' });
			
			for (const backend of backendList) {
				this.renderBackendItem(containerEl, backend);
			}
		} else {
			containerEl.createEl('p', { 
				text: '尚未設定任何 backend，請點擊上方按鈕新增。',
				cls: 'mod-info'
			});
		}
	}

	renderBackendItem(containerEl, backend) {
		const backendDiv = containerEl.createDiv({ cls: 'novelmaker-backend-item' });
		
		// Backend info
		const backendInfo = backendDiv.createDiv({ cls: 'novelmaker-backend-info' });
		const nameEl = backendInfo.createEl('strong', { text: backend.name });
		if (backend.is_default) {
			nameEl.createEl('span', { text: ' (預設)', cls: 'novelmaker-backend-default' });
		}
		backendInfo.createEl('br');
		backendInfo.createEl('small', { text: backend.base_url });
		
		// Backend actions
		const backendActions = backendDiv.createDiv({ cls: 'novelmaker-backend-actions' });
		
		if (!backend.is_default) {
			this.addSetDefaultButton(backendActions, backend);
		}
		this.addCheckButton(backendActions, backend);
		this.addEditButton(backendActions, backend);
		this.addDeleteButton(backendActions, backend);
	}

	addSetDefaultButton(container, backend) {
		const useBtn = container.createEl('button', { text: '設為預設' });
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

	addCheckButton(container, backend) {
		const checkBtn = container.createEl('button', { text: '檢查' });
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
	}

	addEditButton(container, backend) {
		const editBtn = container.createEl('button', { text: '編輯' });
		editBtn.onclick = () => {
			new BackendModal(this.app, this.plugin, backend, async (data) => {
				const loadingModal = new LoadingModal(this.app, '正在更新 Backend...');
				loadingModal.open();

				try {
					const vaultPath = this.app.vault.adapter.basePath;
					let cmd = `${this.plugin.settings.cliPath} backend add "${data.name}" --base_url "${data.base_url}"`;
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
					new Notice(`✅ Backend "${data.name}" 更新成功！`);
					this.display();
				} catch (error) {
					new Notice(`❌ 更新 Backend 失敗: ${error.message}`);
					console.error(error);
				} finally {
					loadingModal.forceCloseNow();
				}
			}).open();
		};
	}

	addDeleteButton(container, backend) {
		const deleteBtn = container.createEl('button', { text: '刪除', cls: 'mod-warning' });
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

	renderBackendSelection(containerEl, backends, defaultBackend) {
		containerEl.createEl('h3', { text: '生成設定' });

		const backendSetting = new Setting(containerEl)
			.setName('Backend 名稱')
			.setDesc(`LLM Backend 的名稱。留空則使用設定檔中的預設 backend${defaultBackend ? ` (${defaultBackend})` : ''}。`);

		if (backends.length > 0) {
			backendSetting.addDropdown((dropdown) => {
				dropdown.addOption('', `使用預設${defaultBackend ? ` (${defaultBackend})` : ''}`);
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
	}

	renderBackendInfo(containerEl) {
		const infoEl = containerEl.createDiv({ cls: 'novelmaker-api-warning' });
		infoEl.createEl('span', { text: 'ℹ️ 提示： ', cls: 'novelmaker-warning-icon' });
		infoEl.createEl('span', { 
			text: 'Backend 配置（包括 API Key）儲存在 ~/.novelmaker/config.toml 中，不會同步到 vault。每個 backend 可以設定自己的預設模型和超時時間。',
			cls: 'novelmaker-warning-text'
		});
	}

	renderGenerationSettings(containerEl) {
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

module.exports = NovelMakerSettingTab;
