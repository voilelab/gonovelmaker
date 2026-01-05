const { Notice } = require('obsidian');
const LoadingModal = require('../modals/loading');
const RewritePromptModal = require('../modals/rewrite-prompt');
const RewriteConfirmModal = require('../modals/rewrite-confirm');
const { executeCLI, parseJSONOutput, buildCLICommand } = require('../utils/cli');

/**
 * Register rewrite command
 */
function registerRewriteCommand(plugin) {
	plugin.addCommand({
		id: 'rewrite-lines',
		name: '改寫選取的行',
		editorCheckCallback: (checking, editor, view) => {
			const file = view.file;
			
			// Only available for markdown files
			if (!file || file.extension !== 'md') {
				return false;
			}
			
			if (!checking) {
				// Get selection or current line
				const selection = editor.getSelection();
				let lineStart, lineEnd, selectedText;
				
				if (selection) {
					// User has selected text
					const from = editor.getCursor('from');
					const to = editor.getCursor('to');
					lineStart = from.line + 1; // Convert to 1-indexed
					lineEnd = to.line + 1;
					selectedText = selection;
				} else {
					// No selection, use current line
					const cursor = editor.getCursor();
					lineStart = cursor.line + 1;
					lineEnd = cursor.line + 1;
					selectedText = editor.getLine(cursor.line);
				}
				
				// Show prompt goal input modal
				new RewritePromptModal(plugin.app, async (promptGoal) => {
					await executeRewrite(plugin, editor, file, lineStart, lineEnd, selectedText, promptGoal);
				}).open();
			}
			
			return true;
		}
	});
}

async function executeRewrite(plugin, editor, file, lineStart, lineEnd, selectedText, promptGoal) {
	try {
		// Default context lines
		const contextPrev = 3;
		const contextNext = 3;
		
		const loadingModal = new LoadingModal(plugin.app, '正在改寫文字...請稍候');
		loadingModal.open();
		
		try {
			const vaultPath = plugin.app.vault.adapter.basePath;
			
			const args = buildCLICommand('rewrite', {
				json: true,
				filepath: file.path,
				lineStart,
				lineEnd,
				contextPrev,
				contextNext,
				promptGoal,
				backend: plugin.settings.backend,
				timeout: plugin.settings.timeout,
			});
			
			const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
			const output = parseJSONOutput(stdout);
			
			// Show confirmation modal with before/after comparison
			new RewriteConfirmModal(
				plugin.app,
				editor,
				file,
				selectedText,
				output.rewritten_content,
				lineStart,
				lineEnd,
				output
			).open();
			
		} catch (error) {
			new Notice(`❌ 錯誤: ${error.message}`);
			console.error(error);
		} finally {
			loadingModal.forceCloseNow();
		}
		
	} catch (error) {
		new Notice(`❌ 錯誤: ${error.message}`);
		console.error(error);
	}
}

module.exports = {
	registerRewriteCommand,
};
