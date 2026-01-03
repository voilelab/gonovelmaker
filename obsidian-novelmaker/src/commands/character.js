const { Notice } = require('obsidian');
const LoadingModal = require('../modals/loading');
const ResultModal = require('../modals/result');
const GenCharModal = require('../modals/gen-char');
const GenCharCurrModal = require('../modals/gen-char-curr');
const GenCharImgModal = require('../modals/gen-char-img');
const { executeCLI, parseJSONOutput, openGeneratedFile, buildCLICommand } = require('../utils/cli');

/**
 * Register all character-related commands
 */
function registerCharacterCommands(plugin) {
	registerGenCharCommand(plugin);
	registerGenCharCurrCommand(plugin);
	registerGenCharImgCommand(plugin);
}

/**
 * Register gen-char command
 */
function registerGenCharCommand(plugin) {
	plugin.addCommand({
		id: 'gen-character',
		name: '生成角色',
		callback: () => {
			new GenCharModal(plugin.app, async (name, prompt) => {
				const loadingModal = new LoadingModal(plugin.app, '正在生成角色...請稍候');
				loadingModal.open();

				try {
					const vaultPath = plugin.app.vault.adapter.basePath;
					
					const cmd = buildCLICommand(plugin.settings.cliPath, 'gen-char', {
						json: true,
						name,
						prompt,
						backend: plugin.settings.backend,
					});
					
					const { stdout, stderr } = await executeCLI(plugin, cmd, vaultPath);
					const output = parseJSONOutput(stdout);

					openGeneratedFile(plugin.app, plugin.settings, output.filepath);
					new ResultModal(plugin.app, '角色生成完成', output).open();
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
 * Register gen-char-curr command
 */
function registerGenCharCurrCommand(plugin) {
	plugin.addCommand({
		id: 'gen-char-curr',
		name: '重新生成當前角色',
		checkCallback: (checking) => {
			const activeFile = plugin.app.workspace.getActiveFile();
			if (activeFile && activeFile.path.startsWith('Character/') && activeFile.extension === 'md') {
				if (!checking) {
					new GenCharCurrModal(plugin.app, activeFile, async (filepath) => {
						const loadingModal = new LoadingModal(plugin.app, '正在重新生成角色...請稍候');
						loadingModal.open();

						try {
							const vaultPath = plugin.app.vault.adapter.basePath;
							
							const cmd = buildCLICommand(plugin.settings.cliPath, 'gen-char-curr', {
								json: true,
								filepath,
								backend: plugin.settings.backend,
							});
							
							const { stdout, stderr } = await executeCLI(plugin, cmd, vaultPath);
							
							try {
								const output = parseJSONOutput(stdout);
								new ResultModal(plugin.app, '角色重新生成完成', output).open();
							} catch (jsonError) {
								new Notice('✅ 角色重新生成成功！');
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

/**
 * Register gen-char-img command
 */
function registerGenCharImgCommand(plugin) {
	plugin.addCommand({
		id: 'gen-char-img',
		name: '生成角色圖片',
		checkCallback: (checking) => {
			const activeFile = plugin.app.workspace.getActiveFile();
			if (activeFile && activeFile.path.startsWith('Character/') && activeFile.extension === 'md') {
				if (!checking) {
					new GenCharImgModal(plugin.app, activeFile, async (name, prompt) => {
						const loadingModal = new LoadingModal(plugin.app, '正在生成角色圖片...請稍候');
						loadingModal.open();

						try {
							const vaultPath = plugin.app.vault.adapter.basePath;
							
							const cmd = buildCLICommand(plugin.settings.cliPath, 'gen-char-img', {
								json: true,
								name,
								prompt,
								backend: plugin.settings.backend,
							});
							
							const { stdout, stderr } = await executeCLI(plugin, cmd, vaultPath);
							const output = parseJSONOutput(stdout);

							if (output && output.filepath) {
								new Notice(`✅ 角色圖片生成成功！\n檔案：${output.filepath}`);
							} else {
								new Notice('✅ 角色圖片生成成功！');
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
	registerCharacterCommands,
};
