const { Plugin } = require('obsidian');

// Import settings
const DEFAULT_SETTINGS = require('./src/settings/constants');
const NovelMakerSettingTab = require('./src/settings/settings-tab');

// Import commands
const { registerChapterCommands } = require('./src/commands/chapter');
const { registerCharacterCommands } = require('./src/commands/character');
const { registerExportCommand } = require('./src/commands/export');
const { registerLogsCommand } = require('./src/commands/logs');
const { registerRewriteCommand } = require('./src/commands/rewrite');

class NovelMakerPlugin extends Plugin {
	async onload() {
		// Load settings
		await this.loadSettings();
		
		// Add settings tab
		this.addSettingTab(new NovelMakerSettingTab(this.app, this));
		
		// Register commands
		registerChapterCommands(this);
		registerCharacterCommands(this);
		registerExportCommand(this);
		registerLogsCommand(this);
		registerRewriteCommand(this);
	}

	onunload() {
		// Cleanup if needed
	}
	
	async loadSettings() {
		this.settings = Object.assign({}, DEFAULT_SETTINGS, await this.loadData());
	}

	async saveSettings() {
		await this.saveData(this.settings);
	}
}

module.exports = NovelMakerPlugin;
