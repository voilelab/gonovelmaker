const { Modal, Notice, Setting } = require('obsidian');

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

module.exports = GenNextEmptyModal;
