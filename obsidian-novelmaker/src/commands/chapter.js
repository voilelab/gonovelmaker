const { Notice } = require('obsidian');
const LoadingModal = require('../modals/loading');
const ResultModal = require('../modals/result');
const GenNextModal = require('../modals/gen-next');
const GenNextEmptyModal = require('../modals/gen-next-empty');
const GenCurrModal = require('../modals/gen-curr');
const { executeCLI, parseJSONOutput, openGeneratedFile, buildCLICommand } = require('../utils/cli');

/**
 * Register all chapter-related commands
 */
function registerChapterCommands(plugin) {
	registerGenNextCommand(plugin);
	registerGenNextEmptyCommand(plugin);
	registerGenCurrCommand(plugin);
}

/**
 * Register gen-next command
 */
function registerGenNextCommand(plugin) {
	plugin.addCommand({
		id: 'gen-next-chapter',
		name: '生成下一章',
		callback: () => {
			new GenNextModal(plugin.app, async (title, prompt, prevCount) => {
				const loadingModal = new LoadingModal(plugin.app, '正在生成下一章...請稍候');
				loadingModal.open();

				try {
					const vaultPath = plugin.app.vault.adapter.basePath;
					
					const args = buildCLICommand('gen-next', {
						json: true,
						title,
						prompt,
						prevCount,
						backend: plugin.settings.backend,
					});
					
					const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
					const output = parseJSONOutput(stdout);

					openGeneratedFile(plugin.app, plugin.settings, output.filepath);
					new ResultModal(plugin.app, '章節生成完成', output).open();
				} catch (error) {
					new Notice(`❌ 錯誤: ${error.message}`);
					console.error(error);
				} finally {
					loadingModal.forceCloseNow();
				}
			}).open();
		}
	});
}

/**
 * Register gen-next-empty command
 */
function registerGenNextEmptyCommand(plugin) {
	plugin.addCommand({
		id: 'gen-next-empty-chapter',
		name: '生成空白下一章',
		callback: () => {
			new GenNextEmptyModal(plugin.app, async (title, prompt) => {
				const loadingModal = new LoadingModal(plugin.app, '正在建立空白章節...請稍候');
				loadingModal.open();

				try {
					const vaultPath = plugin.app.vault.adapter.basePath;
					
					const args = buildCLICommand('gen-next-empty', {
						json: true,
						title,
						prompt,
					});
					
					const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
					const output = parseJSONOutput(stdout);

					openGeneratedFile(plugin.app, plugin.settings, output.filepath);
					new Notice('✅ 空白章節建立成功！');
				} catch (error) {
					new Notice(`❌ 錯誤: ${error.message}`);
					console.error(error);
				} finally {
					loadingModal.forceCloseNow();
				}
			}).open();
		}
	});
}

/**
 * Register gen-curr command
 */
function registerGenCurrCommand(plugin) {
	plugin.addCommand({
		id: 'gen-curr-chapter',
		name: '重新生成當前章節',
		checkCallback: (checking) => {
			const activeFile = plugin.app.workspace.getActiveFile();
			if (activeFile && activeFile.path.startsWith('Story/') && activeFile.extension === 'md') {
				if (!checking) {
					new GenCurrModal(plugin.app, activeFile, async (filepath, prevCount) => {
						const loadingModal = new LoadingModal(plugin.app, '正在重新生成章節...請稍候');
						loadingModal.open();

						try {
							const vaultPath = plugin.app.vault.adapter.basePath;
							
							const args = buildCLICommand('gen-curr', {
								json: true,
								filepath,
								prevCount,
								backend: plugin.settings.backend,
								timeout: plugin.settings.timeout,
							});
							
							const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
							
							try {
								const output = parseJSONOutput(stdout);
								new ResultModal(plugin.app, '章節重新生成完成', output).open();
							} catch (jsonError) {
								new Notice('✅ 章節重新生成成功！');
							}
						} catch (error) {
							new Notice(`❌ 錯誤: ${error.message}`);
							console.error(error);
						} finally {
							loadingModal.forceCloseNow();
						}
					}).open();
				}
				return true;
			}
			return false;
		}
	});
}

module.exports = {
	registerChapterCommands,
};
