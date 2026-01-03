const { Modal, Notice, Setting } = require('obsidian');

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

module.exports = GenNextModal;
