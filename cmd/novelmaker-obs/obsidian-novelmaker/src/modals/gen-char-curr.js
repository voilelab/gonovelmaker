const { Modal, Notice, Setting } = require('obsidian');

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

module.exports = GenCharCurrModal;
