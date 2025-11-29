const { Plugin, Modal, Notice, Setting, PluginSettingTab } = require('obsidian');
const { exec } = require('child_process');
const { promisify } = require('util');

const execAsync = promisify(exec);

// Default settings
const DEFAULT_SETTINGS = {
	cliPath: './novelmaker-obs',
	baseUrl: '',
	apiKey: '',
	model: '',
	timeout: 60, // default timeout in seconds
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
					try {
						new Notice('正在生成下一章...');
						
						// Get vault path
						const vaultPath = this.app.vault.adapter.basePath;
						
						// Build the command using configured CLI path
						let cmd = `${this.settings.cliPath} gen-next --title "${title}"`;
						if (prompt && prompt.trim()) {
							cmd += ` --prompt "${prompt}"`;
						}
						if (prevCount !== null && prevCount !== undefined) {
							cmd += ` --prev-chapters ${prevCount}`;
						}
						if (this.settings.baseUrl && this.settings.baseUrl.trim()) {
							cmd += ` --base-url "${this.settings.baseUrl}"`;
						}
						if (this.settings.apiKey && this.settings.apiKey.trim()) {
							cmd += ` --api-key "${this.settings.apiKey}"`;
						}
						if (this.settings.model && this.settings.model.trim()) {
							cmd += ` --model "${this.settings.model}"`;
						}
						if (this.settings.timeout && !isNaN(this.settings.timeout)) {
							cmd += ` --timeout ${this.settings.timeout}`;
						}
						
						// Call CLI with the input, setting cwd to vault path
						const { stdout, stderr } = await execAsync(cmd, { cwd: vaultPath });
						
						// Show success notification
						new Notice('✅ 章節生成成功！');
						
						// Optionally show output in console
						if (stdout) console.log(stdout);
						if (stderr) console.error(stderr);
					} catch (error) {
						new Notice(`❌ 錯誤: ${error.message}`);
						console.error(error);
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
					try {
						new Notice('正在生成角色...');
						
						// Get vault path
						const vaultPath = this.app.vault.adapter.basePath;
						
						// Build the command using configured CLI path
						let cmd = `${this.settings.cliPath} gen-char --name "${name}"`;
						if (prompt && prompt.trim()) {
							cmd += ` --prompt "${prompt}"`;
						}
						if (this.settings.baseUrl && this.settings.baseUrl.trim()) {
							cmd += ` --base-url "${this.settings.baseUrl}"`;
						}
						if (this.settings.apiKey && this.settings.apiKey.trim()) {
							cmd += ` --api-key "${this.settings.apiKey}"`;
						}
						if (this.settings.model && this.settings.model.trim()) {
							cmd += ` --model "${this.settings.model}"`;
						}
						if (this.settings.timeout && !isNaN(this.settings.timeout)) {
							cmd += ` --timeout ${this.settings.timeout}`;
						}
						
						// Call CLI with the input, setting cwd to vault path
						const { stdout, stderr } = await execAsync(cmd, { cwd: vaultPath });
						
						// Show success notification
						new Notice('✅ 角色生成成功！');
						
						// Optionally show output in console
						if (stdout) console.log(stdout);
						if (stderr) console.error(stderr);
					} catch (error) {
						new Notice(`❌ 錯誤: ${error.message}`);
						console.error(error);
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

class NovelMakerSettingTab extends PluginSettingTab {
	constructor(app, plugin) {
		super(app, plugin);
		this.plugin = plugin;

		this.timeoutSetting = null;
	}

	display() {
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

		new Setting(containerEl)
			.setName('API Base URL')
			.setDesc('LLM API 的基礎 URL（選填，例如：https://api.openai.com/v1）')
			.addText((text) =>
				text
					.setPlaceholder('https://api.openai.com/v1')
					.setValue(this.plugin.settings.baseUrl)
					.onChange(async (value) => {
						this.plugin.settings.baseUrl = value;
						await this.plugin.saveSettings();
					})
			);

		new Setting(containerEl)
			.setName('API Key')
			.setDesc('LLM API 的金鑰（選填）')
			.addText((text) =>
				text
					.setPlaceholder('sk-...')
					.setValue(this.plugin.settings.apiKey)
					.onChange(async (value) => {
						this.plugin.settings.apiKey = value;
						await this.plugin.saveSettings();
					})
			);

		// Warning message for API Key storage
		const warningEl = containerEl.createDiv({ cls: 'novelmaker-api-warning' });
		warningEl.createEl('span', { text: '⚠ 提示： ', cls: 'novelmaker-warning-icon' });
		warningEl.createEl('span', { 
			text: '此設定會被寫入 .obsidian/plugins/obsidian-novelmaker/data.json。' +
				'若你的 vault 使用 git / Sync，請避免把這個檔案推上遠端。' +
				'建議使用 ~/.novelmaker/config.toml 來設定 API Key，避免敏感資訊被同步。',
			cls: 'novelmaker-warning-text'
		});

		new Setting(containerEl)
			.setName('Model')
			.setDesc('使用的 LLM 模型名稱（選填，例如：gpt-4、claude-3-opus-20240229）')
			.addText((text) =>
				text
					.setPlaceholder('gpt-4')
					.setValue(this.plugin.settings.model)
					.onChange(async (value) => {
						this.plugin.settings.model = value;
						await this.plugin.saveSettings();
					})
			);
		
		this.timeoutSetting = new Setting(containerEl)
			.setName(`API 請求超時 (秒) (${this.plugin.settings.timeout} 秒)`)
			.setDesc('設定與 LLM API 通訊的超時時間（秒）（預設：60 秒）')
			.addSlider((slider) =>
				slider
					.setLimits(10, 300, 10)
					.setValue(this.plugin.settings.timeout)
					.onChange(async (value) => {
						if (isNaN(value)) {
							return;
						}
						this.plugin.settings.timeout = value;
						this.timeoutSetting.setName(`API 請求超時 (秒) (${value} 秒)`);
						await this.plugin.saveSettings();
					})
			);
	}
}

module.exports = NovelMakerPlugin;
