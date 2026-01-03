const { Modal, Notice, Setting } = require('obsidian');
const { exec } = require('child_process');
const { promisify } = require('util');

const execAsync = promisify(exec);

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
		this.availableModels = []; // Store available models
		this.modelsLoading = false; // Track loading state
	}

	async fetchAvailableModels() {
		if (!this.backend || !this.backend.name) {
			return [];
		}

		try {
			this.modelsLoading = true;
			const vaultPath = this.app.vault.adapter.basePath;
			const { stdout } = await execAsync(
				`${this.plugin.settings.cliPath} backend list-available-models "${this.backend.name}" --json`,
				{ cwd: vaultPath }
			);
			const result = JSON.parse(stdout);
			if (result.success && Array.isArray(result.models)) {
				return result.models;
			}
			return [];
		} catch (error) {
			console.error('Failed to fetch available models:', error);
			return [];
		} finally {
			this.modelsLoading = false;
		}
	}

	async onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: this.backend ? '編輯 Backend' : '新增 Backend' });

		// Fetch available models if editing existing backend
		if (this.backend) {
			this.availableModels = await this.fetchAvailableModels();
		}

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

		// Model selection with dropdown if available models are loaded
		const modelSetting = new Setting(contentEl)
			.setName('預設模型')
			.setDesc('此 backend 的預設模型名稱（選填）');

		if (this.availableModels.length > 0) {
			// Use dropdown with available models
			modelSetting.addDropdown((dropdown) => {
				dropdown.addOption('', '（選擇模型...）');
				this.availableModels.forEach(model => {
					dropdown.addOption(model, model);
				});
				dropdown
					.setValue(this.model)
					.onChange((value) => {
						this.model = value;
					});
			});
		} else {
			// Fallback to text input
			modelSetting.addText((text) => {
				text
					.setPlaceholder('例如：gpt-4, claude-3-opus-20240229')
					.setValue(this.model)
					.onChange((value) => {
						this.model = value;
					});
				text.inputEl.style.width = '100%';
			});
		}

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

module.exports = BackendModal;
