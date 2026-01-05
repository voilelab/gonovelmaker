const { Modal, Setting } = require('obsidian');

class RewritePromptModal extends Modal {
	constructor(app, onSubmit) {
		super(app);
		this.promptGoal = '去除 AI 味的慣用詞與過度評述語氣';
		this.onSubmit = onSubmit;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '改寫目標' });
		
		contentEl.createEl('p', { 
			text: '請描述您想要如何改寫選取的文字：',
			cls: 'setting-item-description'
		});

		new Setting(contentEl)
			.setName('改寫目標')
			.setDesc('例如：去除 AI 味、改善語氣、簡化句子、增加描述等')
			.addTextArea((text) => {
				text
					.setValue(this.promptGoal)
					.onChange((value) => {
						this.promptGoal = value;
					});
				text.inputEl.rows = 4;
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
					.setButtonText('開始改寫')
					.setCta()
					.onClick(() => {
						this.close();
						this.onSubmit(this.promptGoal);
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

module.exports = RewritePromptModal;
