const esbuild = require('esbuild');
const fs = require('fs');

// Build the plugin
esbuild.build({
	entryPoints: ['main.js'],
	bundle: true,
	external: ['obsidian', 'electron'],
	format: 'cjs',
	target: 'es2018',
	outfile: 'dist/main.js',
	platform: 'node',
	sourcemap: false,
}).then(() => {
	console.log('✅ Build complete!');
	// Copy manifest and styles
	fs.copyFileSync('manifest.json', 'dist/manifest.json');
	fs.copyFileSync('styles.css', 'dist/styles.css');
	console.log('✅ Copied manifest.json and styles.css');
}).catch((err) => {
	console.error('❌ Build failed:', err);
	process.exit(1);
});
