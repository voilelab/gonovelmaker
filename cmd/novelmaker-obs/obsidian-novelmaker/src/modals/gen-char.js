const { Modal, Notice, Setting } = require('obsidian');

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

module.exports = GenCharModal;
