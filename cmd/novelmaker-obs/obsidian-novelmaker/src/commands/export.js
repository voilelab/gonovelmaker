const { Notice } = require('obsidian');
const LoadingModal = require('../modals/loading');
const ExportModal = require('../modals/export');
const { executeCLI, buildCLICommand } = require('../utils/cli');

/**
 * Register export command
 */
function registerExportCommand(plugin) {
	plugin.addCommand({
		id: 'export-novel',
		name: '匯出小說',
		callback: () => {
			new ExportModal(plugin.app, async (outputPath) => {
				const loadingModal = new LoadingModal(plugin.app, '正在匯出小說...請稍候');
				loadingModal.open();

				try {
					const vaultPath = plugin.app.vault.adapter.basePath;
					
					const cmd = buildCLICommand(plugin.settings.cliPath, 'export', {
						output: outputPath,
						type: 'txt',
					});
					
					await executeCLI(plugin, cmd, vaultPath);
					new Notice('✅ 小說匯出成功！');
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

module.exports = {
	registerExportCommand,
};
