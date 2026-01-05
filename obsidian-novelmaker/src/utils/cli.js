const { execFile } = require('child_process');
const { promisify } = require('util');
const { Notice } = require('obsidian');

const execFileAsync = promisify(execFile);

/**
 * Execute CLI command and return JSON output
 */
async function executeCLI(plugin, cliPath, args, vaultPath) {
	const { stdout, stderr } = await execFileAsync(cliPath, args, { cwd: vaultPath });
	
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
 * Build CLI command arguments array with common options
 */
function buildCLICommand(baseCommand, options = {}) {
	// Split baseCommand by spaces to handle multi-word commands like "character gen"
	const args = baseCommand.split(/\s+/);
	
	if (options.json) {
		args.push('--json');
	}
	
	if (options.title) {
		args.push('--title', options.title);
	}
	
	if (options.name) {
		args.push('--name', options.name);
	}
	
	if (options.filepath) {
		args.push('--filepath', options.filepath);
	}
	
	if (options.prompt && options.prompt.trim()) {
		args.push('--prompt', options.prompt);
	}
	
	if (options.promptGoal && options.promptGoal.trim()) {
		args.push('--prompt', options.promptGoal);
	}
	
	if (options.prevCount !== null && options.prevCount !== undefined) {
		args.push('--prev-chapters', String(options.prevCount));
	}
	
	if (options.backend && options.backend.trim()) {
		args.push('--backend', options.backend);
	}
	
	if (options.timeout && !isNaN(options.timeout)) {
		args.push('--timeout', String(options.timeout));
	}
	
	if (options.output) {
		args.push('--output', options.output);
	}
	
	if (options.type) {
		args.push('--type', options.type);
	}
	
	if (options.lineStart !== null && options.lineStart !== undefined) {
		args.push('--line-start', String(options.lineStart));
	}
	
	if (options.lineEnd !== null && options.lineEnd !== undefined) {
		args.push('--line-end', String(options.lineEnd));
	}
	
	if (options.contextPrev !== null && options.contextPrev !== undefined) {
		args.push('--context-prev', String(options.contextPrev));
	}
	
	if (options.contextNext !== null && options.contextNext !== undefined) {
		args.push('--context-next', String(options.contextNext));
	}
	
	return args;
}

module.exports = {
	execFileAsync,
	executeCLI,
	parseJSONOutput,
	openGeneratedFile,
	buildCLICommand,
};
