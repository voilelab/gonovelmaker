const { exec } = require('child_process');
const { promisify } = require('util');
const { Notice } = require('obsidian');

const execAsync = promisify(exec);

/**
 * Execute CLI command and return JSON output
 */
async function executeCLI(plugin, command, vaultPath) {
	const { stdout, stderr } = await execAsync(command, { cwd: vaultPath });
	
	if (stdout) console.log(stdout);
	if (stderr) console.error(stderr);
	
	return { stdout, stderr };
}

/**
 * Parse JSON output from CLI
 */
function parseJSONOutput(stdout) {
	try {
		return JSON.parse(stdout);
	} catch (jsonError) {
		throw new Error('無法解析生成結果的 JSON 輸出');
	}
}

/**
 * Open generated file in Obsidian after a delay
 */
function openGeneratedFile(app, settings, filepath, delay = null) {
	if (!settings.openAfterGen) {
		return;
	}
	
	if (!filepath) {
		new Notice('⚠ 生成結果中缺少 filepath 資訊');
		return;
	}
	
	const delayMs = delay !== null ? delay : settings.openAfterGenMs;
	
	setTimeout(() => {
		const file = app.vault.getAbstractFileByPath(filepath);
		if (file) {
			app.workspace.getLeaf().openFile(file);
		} else {
			new Notice(`⚠ 無法在 Vault 中找到生成的檔案: ${filepath}`);
		}
	}, delayMs);
}

/**
 * Build CLI command with common options
 */
function buildCLICommand(cliPath, baseCommand, options = {}) {
	let cmd = `${cliPath} ${baseCommand}`;
	
	if (options.json) {
		cmd += ' --json';
	}
	
	if (options.title) {
		cmd += ` --title "${options.title}"`;
	}
	
	if (options.name) {
		cmd += ` --name "${options.name}"`;
	}
	
	if (options.filepath) {
		cmd += ` --filepath "${options.filepath}"`;
	}
	
	if (options.prompt && options.prompt.trim()) {
		cmd += ` --prompt "${options.prompt}"`;
	}
	
	if (options.prevCount !== null && options.prevCount !== undefined) {
		cmd += ` --prev-chapters ${options.prevCount}`;
	}
	
	if (options.backend && options.backend.trim()) {
		cmd += ` --backend "${options.backend}"`;
	}
	
	if (options.timeout && !isNaN(options.timeout)) {
		cmd += ` --timeout ${options.timeout}`;
	}
	
	if (options.output) {
		cmd += ` --output "${options.output}"`;
	}
	
	if (options.type) {
		cmd += ` --type ${options.type}`;
	}
	
	return cmd;
}

module.exports = {
	execAsync,
	executeCLI,
	parseJSONOutput,
	openGeneratedFile,
	buildCLICommand,
};
