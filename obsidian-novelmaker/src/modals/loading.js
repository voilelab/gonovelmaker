const { Modal } = require('obsidian');

class LoadingModal extends Modal {
	constructor(app, message = "處理中…請稍候") {
		super(app);
		this.message = message;
	}

	onOpen() {
		// remove X
		const closeBtn = this.containerEl.querySelector(".modal-close-button");
		if (closeBtn) closeBtn.remove();

		// block ESC
		this.scope.register([], 'Escape', evt => {
			evt.preventDefault();
		});

		// block clicking background
		this.containerEl.onclick = (evt) => {
			evt.stopPropagation();
		};

		// UI
		this.contentEl.empty();
		this.setTitle(this.message);
		this.contentEl.createEl("div", { text: "⏳ 生成中…" });
	}

	// prevent close
	close() {
		if (this.forceClose) {
			super.close();
		}
	}

	// safe manual close
	forceCloseNow() {
		this.forceClose = true;
		this.close();
	}
}

module.exports = LoadingModal;
