const { Modal, Notice, Setting } = require('obsidian');

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

module.exports = GenCharImgModal;
