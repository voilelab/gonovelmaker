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

		const pathSetting = new Setting(contentEl)
			.setName('輸出檔案路徑')
			.setDesc('選擇小說匯出的檔案路徑')
			.addText((text) => {
				text
					.setPlaceholder('選擇檔案位置...')
					.setValue(this.outputPath)
					.onChange((value) => {
						this.outputPath = value;
					});
				text.inputEl.style.width = '100%';
				// Store reference to update later
				this.pathTextComponent = text;
			})
			.addButton((btn) =>
				btn
					.setButtonText('瀏覽...')
					.onClick(async () => {
						try {
							// Try to get electron dialog using modern approach
							let dialog;
							
							// First, try @electron/remote (modern replacement package)
							try {
								const remote = require('@electron/remote');
								dialog = remote.dialog;
							} catch (e) {
								// Fallback: try getting dialog from electron directly
								const electron = require('electron');
								if (electron.dialog) {
									dialog = electron.dialog;
								} else if (electron.remote?.dialog) {
									// Last resort: use deprecated remote (with warning)
									console.warn('Using deprecated electron.remote as fallback');
									dialog = electron.remote.dialog;
								}
							}
							
							if (!dialog) {
								throw new Error('Dialog API not available');
							}
							
							const result = await dialog.showSaveDialog({
								title: '選擇匯出位置',
								defaultPath: 'novel.txt',
								filters: [
									{ name: '文字檔案', extensions: ['txt'] },
									{ name: '所有檔案', extensions: ['*'] }
								],
								properties: ['createDirectory', 'showOverwriteConfirmation']
							});

							if (!result.canceled && result.filePath) {
								this.outputPath = result.filePath;
								// Update the text input using the component reference
								if (this.pathTextComponent) {
									this.pathTextComponent.setValue(result.filePath);
								}
							}
						} catch (error) {
							// Fallback if dialog is not available
							console.error('File dialog error:', error);
							new Notice('⚠ 無法開啟檔案對話框，請手動輸入路徑');
						}
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
