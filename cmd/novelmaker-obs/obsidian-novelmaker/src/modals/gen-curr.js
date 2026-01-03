const { Modal, Notice, Setting } = require('obsidian');

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

module.exports = GenCurrModal;
