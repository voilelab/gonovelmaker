const LogsModal = require('../modals/logs');

/**
 * Register logs command
 */
function registerLogsCommand(plugin) {
	plugin.addCommand({
		id: 'show-logs',
		name: '查看日誌',
		callback: () => {
			const vaultPath = plugin.app.vault.adapter.basePath;
			new LogsModal(plugin.app, vaultPath).open();
		}
	});
}

module.exports = {
	registerLogsCommand,
};
