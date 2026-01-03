const { Modal, Notice, Setting } = require('obsidian');

class ExportModal extends Modal {
	constructor(app, onSubmit) {
		super(app);
		this.outputPath = '';
		this.onSubmit = onSubmit;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '匯出小說' });

		new Setting(contentEl)
			.setName('輸出檔案路徑')
			.setDesc('請輸入小說匯出的完整檔案路徑（例如：C:\\Users\\novel.txt 或 /home/user/novel.txt）')
			.addText((text) => {
				text
					.setPlaceholder('請輸入完整檔案路徑...')
					.setValue(this.outputPath)
					.onChange((value) => {
						this.outputPath = value;
					});
				text.inputEl.style.width = '100%';
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

module.exports = ExportModal;
