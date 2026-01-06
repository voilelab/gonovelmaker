const { Modal, Setting } = require('obsidian');

class RewriteConfirmModal extends Modal {
	constructor(app, editor, file, originalText, rewrittenText, lineStart, lineEnd, usageInfo) {
		super(app);
		this.editor = editor;
		this.file = file;
		this.originalText = originalText;
		this.rewrittenText = rewrittenText;
		this.lineStart = lineStart;
		this.lineEnd = lineEnd;
		this.usageInfo = usageInfo;
	}

	onOpen() {
		const { contentEl } = this;
		
		contentEl.createEl('h2', { text: '改寫確認' });
		
		// File info
		new Setting(contentEl)
			.setName('檔案')
			.setDesc(this.file.path);
		
		new Setting(contentEl)
			.setName('行數範圍')
			.setDesc(`第 ${this.lineStart} 行到第 ${this.lineEnd} 行`);
		
		// Before section
		contentEl.createEl('h3', { text: '原文' });
		const beforeDiv = contentEl.createDiv({ cls: 'rewrite-before' });
		beforeDiv.createEl('pre', { 
			text: this.originalText,
			cls: 'rewrite-text-box'
		});
		
		// After section
		contentEl.createEl('h3', { text: '改寫後' });
		const afterDiv = contentEl.createDiv({ cls: 'rewrite-after' });
		afterDiv.createEl('pre', { 
			text: this.rewrittenText,
			cls: 'rewrite-text-box'
		});
		
		// Usage info
		if (this.usageInfo) {
			const usageDiv = contentEl.createDiv({ cls: 'rewrite-usage-info' });
			usageDiv.createEl('p', { 
				text: `Token 使用量: 輸入 ${this.usageInfo.input_tokens}, 輸出 ${this.usageInfo.output_tokens}, 總計 ${this.usageInfo.total_tokens}`,
				cls: 'mod-muted'
			});
		}
		
		// Buttons
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
					.setButtonText('接受改寫')
					.setCta()
					.onClick(() => {
						this.applyRewrite();
						this.close();
					})
			);
	}

	applyRewrite() {
		// Get current document content
		const doc = this.editor.getDoc();
		
		// Convert from 1-indexed to 0-indexed
		const startLine = this.lineStart - 1;
		const endLine = this.lineEnd - 1;
		
		// Get the start of the first line and end of the last line
		const from = { line: startLine, ch: 0 };
		const to = { line: endLine, ch: this.editor.getLine(endLine).length };
		
		// Replace the content
		this.editor.replaceRange(this.rewrittenText, from, to);
		
		// Show success notice
		const { Notice } = require('obsidian');
		new Notice('✅ 已套用改寫內容！');
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

module.exports = RewriteConfirmModal;
