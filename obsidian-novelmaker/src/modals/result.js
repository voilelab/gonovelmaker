const { Modal, Setting } = require('obsidian');

class ResultModal extends Modal {
	constructor(app, title, output) {
		super(app);
		this.modalTitle = title;
		this.output = output;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: this.modalTitle });

		// Success message
		const successDiv = contentEl.createDiv({ cls: 'novelmaker-result-success' });
		successDiv.createEl('span', { text: '✅ ', cls: 'novelmaker-success-icon' });
		successDiv.createEl('span', { text: '生成成功！' });

		// File path if available
		if (this.output.filepath) {
			const fileDiv = contentEl.createDiv({ cls: 'novelmaker-result-info' });
			fileDiv.createEl('strong', { text: '檔案：' });
			fileDiv.createEl('span', { text: this.output.filepath });
		}

		// Token usage section
		if (this.output.input_tokens !== undefined || this.output.output_tokens !== undefined) {
			const usageDiv = contentEl.createDiv({ cls: 'novelmaker-token-usage' });
			usageDiv.createEl('h3', { text: '🔢 Token 使用量' });
			
			const usageGrid = usageDiv.createDiv({ cls: 'novelmaker-token-grid' });
			
			if (this.output.input_tokens !== undefined) {
				const inputRow = usageGrid.createDiv({ cls: 'novelmaker-token-row' });
				inputRow.createEl('span', { text: 'Input tokens:', cls: 'novelmaker-token-label' });
				inputRow.createEl('span', { text: this.output.input_tokens.toLocaleString(), cls: 'novelmaker-token-value' });
			}
			
			if (this.output.output_tokens !== undefined) {
				const outputRow = usageGrid.createDiv({ cls: 'novelmaker-token-row' });
				outputRow.createEl('span', { text: 'Output tokens:', cls: 'novelmaker-token-label' });
				outputRow.createEl('span', { text: this.output.output_tokens.toLocaleString(), cls: 'novelmaker-token-value' });
			}
			
			if (this.output.total_tokens !== undefined) {
				const totalRow = usageGrid.createDiv({ cls: 'novelmaker-token-row novelmaker-token-total' });
				totalRow.createEl('span', { text: 'Total tokens:', cls: 'novelmaker-token-label' });
				totalRow.createEl('span', { text: this.output.total_tokens.toLocaleString(), cls: 'novelmaker-token-value' });
			}
		}

		// Close button
		new Setting(contentEl)
			.addButton((btn) =>
				btn
					.setButtonText('確定')
					.setCta()
					.onClick(() => {
						this.close();
					})
			);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

module.exports = ResultModal;
