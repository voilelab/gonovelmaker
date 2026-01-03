var __getOwnPropNames = Object.getOwnPropertyNames;
var __commonJS = (cb, mod) => function __require() {
  return mod || (0, cb[__getOwnPropNames(cb)[0]])((mod = { exports: {} }).exports, mod), mod.exports;
};

// src/settings/constants.js
var require_constants = __commonJS({
  "src/settings/constants.js"(exports2, module2) {
    var DEFAULT_SETTINGS2 = {
      cliPath: "novelmaker-obs",
      backend: "",
      openAfterGen: false,
      openAfterGenMs: 500
    };
    module2.exports = DEFAULT_SETTINGS2;
  }
});

// src/modals/loading.js
var require_loading = __commonJS({
  "src/modals/loading.js"(exports2, module2) {
    var { Modal } = require("obsidian");
    var LoadingModal = class extends Modal {
      constructor(app, message = "\u8655\u7406\u4E2D\u2026\u8ACB\u7A0D\u5019") {
        super(app);
        this.message = message;
      }
      onOpen() {
        const closeBtn = this.containerEl.querySelector(".modal-close-button");
        if (closeBtn)
          closeBtn.remove();
        this.scope.register([], "Escape", (evt) => {
          evt.preventDefault();
        });
        this.containerEl.onclick = (evt) => {
          evt.stopPropagation();
        };
        this.contentEl.empty();
        this.setTitle(this.message);
        this.contentEl.createEl("div", { text: "\u23F3 \u751F\u6210\u4E2D\u2026" });
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
    };
    module2.exports = LoadingModal;
  }
});

// src/utils/cli.js
var require_cli = __commonJS({
  "src/utils/cli.js"(exports2, module2) {
    var { execFile } = require("child_process");
    var { promisify } = require("util");
    var { Notice } = require("obsidian");
    var execFileAsync = promisify(execFile);
    async function executeCLI(plugin, cliPath, args, vaultPath) {
      const { stdout, stderr } = await execFileAsync(cliPath, args, { cwd: vaultPath });
      if (stdout)
        console.log(stdout);
      if (stderr)
        console.error(stderr);
      return { stdout, stderr };
    }
    function parseJSONOutput(stdout) {
      try {
        return JSON.parse(stdout);
      } catch (jsonError) {
        throw new Error("\u7121\u6CD5\u89E3\u6790\u751F\u6210\u7D50\u679C\u7684 JSON \u8F38\u51FA");
      }
    }
    function openGeneratedFile(app, settings, filepath, delay = null) {
      if (!settings.openAfterGen) {
        return;
      }
      if (!filepath) {
        new Notice("\u26A0 \u751F\u6210\u7D50\u679C\u4E2D\u7F3A\u5C11 filepath \u8CC7\u8A0A");
        return;
      }
      const delayMs = delay !== null ? delay : settings.openAfterGenMs;
      setTimeout(() => {
        const file = app.vault.getAbstractFileByPath(filepath);
        if (file) {
          app.workspace.getLeaf().openFile(file);
        } else {
          new Notice(`\u26A0 \u7121\u6CD5\u5728 Vault \u4E2D\u627E\u5230\u751F\u6210\u7684\u6A94\u6848: ${filepath}`);
        }
      }, delayMs);
    }
    function buildCLICommand(baseCommand, options = {}) {
      const args = baseCommand.split(/\s+/);
      if (options.json) {
        args.push("--json");
      }
      if (options.title) {
        args.push("--title", options.title);
      }
      if (options.name) {
        args.push("--name", options.name);
      }
      if (options.filepath) {
        args.push("--filepath", options.filepath);
      }
      if (options.prompt && options.prompt.trim()) {
        args.push("--prompt", options.prompt);
      }
      if (options.prevCount !== null && options.prevCount !== void 0) {
        args.push("--prev-chapters", String(options.prevCount));
      }
      if (options.backend && options.backend.trim()) {
        args.push("--backend", options.backend);
      }
      if (options.timeout && !isNaN(options.timeout)) {
        args.push("--timeout", String(options.timeout));
      }
      if (options.output) {
        args.push("--output", options.output);
      }
      if (options.type) {
        args.push("--type", options.type);
      }
      return args;
    }
    module2.exports = {
      execFileAsync,
      executeCLI,
      parseJSONOutput,
      openGeneratedFile,
      buildCLICommand
    };
  }
});

// src/modals/backend.js
var require_backend = __commonJS({
  "src/modals/backend.js"(exports2, module2) {
    var { Modal, Notice, Setting } = require("obsidian");
    var { execFileAsync } = require_cli();
    var BackendModal = class extends Modal {
      constructor(app, plugin, backend = null, onSubmit) {
        super(app);
        this.plugin = plugin;
        this.backend = backend;
        this.name = (backend == null ? void 0 : backend.name) || "";
        this.base_url = (backend == null ? void 0 : backend.base_url) || "";
        this.api_key = (backend == null ? void 0 : backend.api_key) || "";
        this.model = (backend == null ? void 0 : backend.model) || "";
        this.timeout = (backend == null ? void 0 : backend.timeout) || 60;
        this.onSubmit = onSubmit;
        this.apiKeyModified = false;
        this.availableModels = [];
        this.modelsLoading = false;
      }
      async fetchAvailableModels() {
        if (!this.backend || !this.backend.name) {
          return [];
        }
        try {
          this.modelsLoading = true;
          const vaultPath = this.app.vault.adapter.basePath;
          const { stdout } = await execFileAsync(
            this.plugin.settings.cliPath,
            ["backend", "list-available-models", this.backend.name, "--json"],
            { cwd: vaultPath }
          );
          const result = JSON.parse(stdout);
          if (result.success && Array.isArray(result.models)) {
            return result.models;
          }
          return [];
        } catch (error) {
          console.error("Failed to fetch available models:", error);
          return [];
        } finally {
          this.modelsLoading = false;
        }
      }
      async onOpen() {
        const { contentEl } = this;
        contentEl.createEl("h2", { text: this.backend ? "\u7DE8\u8F2F Backend" : "\u65B0\u589E Backend" });
        if (this.backend) {
          this.availableModels = await this.fetchAvailableModels();
        }
        new Setting(contentEl).setName("Backend \u540D\u7A31").setDesc("Backend \u7684\u8B58\u5225\u540D\u7A31\uFF08\u5FC5\u586B\uFF09").addText((text) => {
          text.setPlaceholder("\u4F8B\u5982\uFF1Aopenrouter, openai").setValue(this.name).onChange((value) => {
            this.name = value;
          });
          text.inputEl.style.width = "100%";
          if (this.backend) {
            text.inputEl.disabled = true;
          }
        });
        new Setting(contentEl).setName("Base URL").setDesc("API \u7684\u57FA\u790E URL\uFF08\u5FC5\u586B\uFF09").addText((text) => {
          text.setPlaceholder("\u4F8B\u5982\uFF1Ahttps://openrouter.ai/api/v1").setValue(this.base_url).onChange((value) => {
            this.base_url = value;
          });
          text.inputEl.style.width = "100%";
        });
        new Setting(contentEl).setName("API Key").setDesc(this.backend ? "API \u91D1\u9470\uFF08\u7559\u7A7A\u8868\u793A\u4E0D\u4FEE\u6539\uFF09" : "API \u91D1\u9470\uFF08\u5FC5\u586B\uFF09").addText((text) => {
          text.setPlaceholder(this.backend ? "\u7559\u7A7A\u8868\u793A\u4FDD\u6301\u539F\u503C..." : "sk-...").setValue(this.backend ? "" : this.api_key).onChange((value) => {
            this.api_key = value;
            if (this.backend && value.trim()) {
              this.apiKeyModified = true;
            }
          });
          text.inputEl.style.width = "100%";
          text.inputEl.type = "password";
        });
        const modelSetting = new Setting(contentEl).setName("\u9810\u8A2D\u6A21\u578B").setDesc("\u6B64 backend \u7684\u9810\u8A2D\u6A21\u578B\u540D\u7A31\uFF08\u9078\u586B\uFF09");
        if (this.availableModels.length > 0) {
          modelSetting.addDropdown((dropdown) => {
            dropdown.addOption("", "\uFF08\u9078\u64C7\u6A21\u578B...\uFF09");
            this.availableModels.forEach((model) => {
              dropdown.addOption(model, model);
            });
            dropdown.setValue(this.model).onChange((value) => {
              this.model = value;
            });
          });
        } else {
          modelSetting.addText((text) => {
            text.setPlaceholder("\u4F8B\u5982\uFF1Agpt-4, claude-3-opus-20240229").setValue(this.model).onChange((value) => {
              this.model = value;
            });
            text.inputEl.style.width = "100%";
          });
        }
        new Setting(contentEl).setName(`API \u8ACB\u6C42\u8D85\u6642 (\u79D2) (${this.timeout} \u79D2)`).setDesc("\u6B64 backend \u7684\u8D85\u6642\u6642\u9593\uFF08\u79D2\uFF09\uFF08\u9810\u8A2D\uFF1A60 \u79D2\uFF09").addSlider((slider) => {
          slider.setLimits(10, 300, 10).setValue(this.timeout).onChange((value) => {
            if (isNaN(value)) {
              return;
            }
            this.timeout = value;
            const settingEl = contentEl.querySelector(".setting-item:last-of-type .setting-item-name");
            if (settingEl) {
              settingEl.textContent = `API \u8ACB\u6C42\u8D85\u6642 (\u79D2) (${value} \u79D2)`;
            }
          });
        });
        new Setting(contentEl).addButton(
          (btn) => btn.setButtonText("\u53D6\u6D88").onClick(() => {
            this.close();
          })
        ).addButton(
          (btn) => btn.setButtonText(this.backend ? "\u66F4\u65B0" : "\u65B0\u589E").setCta().onClick(() => {
            if (!this.name || !this.name.trim()) {
              new Notice("\u274C \u8ACB\u8F38\u5165 Backend \u540D\u7A31");
              return;
            }
            if (!this.base_url || !this.base_url.trim()) {
              new Notice("\u274C \u8ACB\u8F38\u5165 Base URL");
              return;
            }
            if (!this.backend && (!this.api_key || !this.api_key.trim())) {
              new Notice("\u274C \u8ACB\u8F38\u5165 API Key");
              return;
            }
            this.close();
            this.onSubmit({
              name: this.name,
              base_url: this.base_url,
              api_key: this.api_key,
              model: this.model,
              timeout: this.timeout,
              apiKeyModified: this.apiKeyModified
              // Pass the flag
            });
          })
        );
      }
      onClose() {
        const { contentEl } = this;
        contentEl.empty();
      }
    };
    module2.exports = BackendModal;
  }
});

// src/settings/settings-tab.js
var require_settings_tab = __commonJS({
  "src/settings/settings-tab.js"(exports2, module2) {
    var { Notice, Setting, PluginSettingTab } = require("obsidian");
    var LoadingModal = require_loading();
    var BackendModal = require_backend();
    var { execFileAsync } = require_cli();
    var NovelMakerSettingTab2 = class extends PluginSettingTab {
      constructor(app, plugin) {
        super(app, plugin);
        this.plugin = plugin;
        this.openAfterGenMsSetting = null;
      }
      async display() {
        const { containerEl } = this;
        containerEl.empty();
        containerEl.createEl("h2", { text: "Novel Maker \u8A2D\u5B9A" });
        this.renderCLIPathSetting(containerEl);
        await this.renderBackendManagement(containerEl);
        this.renderGenerationSettings(containerEl);
      }
      renderCLIPathSetting(containerEl) {
        new Setting(containerEl).setName("CLI \u6307\u4EE4\u8DEF\u5F91").setDesc("\u8A2D\u5B9A novelmaker-obs \u6307\u4EE4\u7684\u8DEF\u5F91\u3002\u53EF\u4F7F\u7528\u76F8\u5C0D\u8DEF\u5F91\uFF08\u5982 ./novelmaker-obs\uFF09\u3001\u7D55\u5C0D\u8DEF\u5F91\uFF08\u5982 /usr/local/bin/novelmaker-obs\uFF09\u6216\u6307\u4EE4\u540D\u7A31\uFF08\u5982 novelmaker-obs\uFF0C\u9700\u5728 PATH \u4E2D\uFF09").addText(
          (text) => text.setPlaceholder("./novelmaker-obs").setValue(this.plugin.settings.cliPath).onChange(async (value) => {
            this.plugin.settings.cliPath = value || "./novelmaker-obs";
            await this.plugin.saveSettings();
          })
        );
      }
      async renderBackendManagement(containerEl) {
        containerEl.createEl("h3", { text: "Backend \u7BA1\u7406" });
        const { backends, backendList, defaultBackend } = await this.loadBackends();
        this.renderAddBackendButton(containerEl);
        this.renderBackendList(containerEl, backendList);
        this.renderBackendSelection(containerEl, backends, defaultBackend);
        this.renderBackendInfo(containerEl);
      }
      async loadBackends() {
        let backends = [];
        let backendList = [];
        let defaultBackend = "";
        try {
          const vaultPath = this.app.vault.adapter.basePath;
          const { stdout } = await execFileAsync(this.plugin.settings.cliPath, ["backend", "list", "--json"], { cwd: vaultPath });
          backendList = JSON.parse(stdout);
          if (Array.isArray(backendList)) {
            backends = backendList.map((b) => b.name);
            const defaultBackendObj = backendList.find((b) => b.is_default);
            if (defaultBackendObj) {
              defaultBackend = defaultBackendObj.name;
            }
          }
        } catch (error) {
          console.error("Failed to load backends:", error);
          new Notice("\u26A0 \u7121\u6CD5\u8F09\u5165 backend \u5217\u8868\uFF0C\u8ACB\u78BA\u8A8D CLI \u8DEF\u5F91\u662F\u5426\u6B63\u78BA");
        }
        return { backends, backendList, defaultBackend };
      }
      renderAddBackendButton(containerEl) {
        new Setting(containerEl).setName("\u7BA1\u7406 Backends").setDesc("\u65B0\u589E\u3001\u7DE8\u8F2F\u6216\u522A\u9664 AI backend \u914D\u7F6E").addButton(
          (btn) => btn.setButtonText("\u65B0\u589E Backend").setCta().onClick(() => this.handleAddBackend())
        );
      }
      async handleAddBackend() {
        new BackendModal(this.app, this.plugin, null, async (data) => {
          const loadingModal = new LoadingModal(this.app, "\u6B63\u5728\u65B0\u589E Backend...");
          loadingModal.open();
          try {
            const vaultPath = this.app.vault.adapter.basePath;
            const args = ["backend", "add", data.name, "--base_url", data.base_url, "--api_key", data.api_key];
            if (data.model && data.model.trim()) {
              args.push("--model", data.model);
            }
            if (data.timeout) {
              args.push("--timeout", String(data.timeout));
            }
            await execFileAsync(this.plugin.settings.cliPath, args, { cwd: vaultPath });
            new Notice(`\u2705 Backend "${data.name}" \u65B0\u589E\u6210\u529F\uFF01`);
            this.display();
          } catch (error) {
            new Notice(`\u274C \u65B0\u589E Backend \u5931\u6557: ${error.message}`);
            console.error(error);
          } finally {
            loadingModal.forceCloseNow();
          }
        }).open();
      }
      renderBackendList(containerEl, backendList) {
        if (backendList.length > 0) {
          containerEl.createEl("h4", { text: "\u5DF2\u8A2D\u5B9A\u7684 Backends" });
          for (const backend of backendList) {
            this.renderBackendItem(containerEl, backend);
          }
        } else {
          containerEl.createEl("p", {
            text: "\u5C1A\u672A\u8A2D\u5B9A\u4EFB\u4F55 backend\uFF0C\u8ACB\u9EDE\u64CA\u4E0A\u65B9\u6309\u9215\u65B0\u589E\u3002",
            cls: "mod-info"
          });
        }
      }
      renderBackendItem(containerEl, backend) {
        const backendDiv = containerEl.createDiv({ cls: "novelmaker-backend-item" });
        const backendInfo = backendDiv.createDiv({ cls: "novelmaker-backend-info" });
        const nameEl = backendInfo.createEl("strong", { text: backend.name });
        if (backend.is_default) {
          nameEl.createEl("span", { text: " (\u9810\u8A2D)", cls: "novelmaker-backend-default" });
        }
        backendInfo.createEl("br");
        backendInfo.createEl("small", { text: backend.base_url });
        const backendActions = backendDiv.createDiv({ cls: "novelmaker-backend-actions" });
        if (!backend.is_default) {
          this.addSetDefaultButton(backendActions, backend);
        }
        this.addCheckButton(backendActions, backend);
        this.addEditButton(backendActions, backend);
        this.addDeleteButton(backendActions, backend);
      }
      addSetDefaultButton(container, backend) {
        const useBtn = container.createEl("button", { text: "\u8A2D\u70BA\u9810\u8A2D" });
        useBtn.onclick = async () => {
          try {
            const vaultPath = this.app.vault.adapter.basePath;
            await execFileAsync(this.plugin.settings.cliPath, ["backend", "use", backend.name], { cwd: vaultPath });
            new Notice(`\u2705 \u5DF2\u5C07 "${backend.name}" \u8A2D\u70BA\u9810\u8A2D backend`);
            this.display();
          } catch (error) {
            new Notice(`\u274C \u8A2D\u5B9A\u9810\u8A2D backend \u5931\u6557: ${error.message}`);
            console.error(error);
          }
        };
      }
      addCheckButton(container, backend) {
        const checkBtn = container.createEl("button", { text: "\u6AA2\u67E5" });
        checkBtn.onclick = async () => {
          const loadingModal = new LoadingModal(this.app, "\u6B63\u5728\u6AA2\u67E5 Backend...");
          loadingModal.open();
          try {
            const vaultPath = this.app.vault.adapter.basePath;
            await execFileAsync(this.plugin.settings.cliPath, ["backend", "check", backend.name], { cwd: vaultPath });
            new Notice(`\u2705 Backend "${backend.name}" \u9023\u7DDA\u6B63\u5E38\uFF01`);
          } catch (error) {
            new Notice(`\u274C Backend "${backend.name}" \u9023\u7DDA\u5931\u6557: ${error.message}`);
            console.error(error);
          } finally {
            loadingModal.forceCloseNow();
          }
        };
      }
      addEditButton(container, backend) {
        const editBtn = container.createEl("button", { text: "\u7DE8\u8F2F" });
        editBtn.onclick = () => {
          new BackendModal(this.app, this.plugin, backend, async (data) => {
            const loadingModal = new LoadingModal(this.app, "\u6B63\u5728\u66F4\u65B0 Backend...");
            loadingModal.open();
            try {
              const vaultPath = this.app.vault.adapter.basePath;
              const args = ["backend", "add", data.name, "--base_url", data.base_url];
              if (data.apiKeyModified && data.api_key && data.api_key.trim()) {
                args.push("--api_key", data.api_key);
              }
              if (data.model && data.model.trim()) {
                args.push("--model", data.model);
              }
              if (data.timeout) {
                args.push("--timeout", String(data.timeout));
              }
              await execFileAsync(this.plugin.settings.cliPath, args, { cwd: vaultPath });
              new Notice(`\u2705 Backend "${data.name}" \u66F4\u65B0\u6210\u529F\uFF01`);
              this.display();
            } catch (error) {
              new Notice(`\u274C \u66F4\u65B0 Backend \u5931\u6557: ${error.message}`);
              console.error(error);
            } finally {
              loadingModal.forceCloseNow();
            }
          }).open();
        };
      }
      addDeleteButton(container, backend) {
        const deleteBtn = container.createEl("button", { text: "\u522A\u9664", cls: "mod-warning" });
        deleteBtn.onclick = async () => {
          if (!confirm(`\u78BA\u5B9A\u8981\u522A\u9664 backend "${backend.name}" \u55CE\uFF1F`)) {
            return;
          }
          try {
            const vaultPath = this.app.vault.adapter.basePath;
            await execFileAsync(this.plugin.settings.cliPath, ["backend", "remove", backend.name], { cwd: vaultPath });
            new Notice(`\u2705 Backend "${backend.name}" \u5DF2\u522A\u9664`);
            this.display();
          } catch (error) {
            new Notice(`\u274C \u522A\u9664 Backend \u5931\u6557: ${error.message}`);
            console.error(error);
          }
        };
      }
      renderBackendSelection(containerEl, backends, defaultBackend) {
        containerEl.createEl("h3", { text: "\u751F\u6210\u8A2D\u5B9A" });
        const backendSetting = new Setting(containerEl).setName("Backend \u540D\u7A31").setDesc(`LLM Backend \u7684\u540D\u7A31\u3002\u7559\u7A7A\u5247\u4F7F\u7528\u8A2D\u5B9A\u6A94\u4E2D\u7684\u9810\u8A2D backend${defaultBackend ? ` (${defaultBackend})` : ""}\u3002`);
        if (backends.length > 0) {
          backendSetting.addDropdown((dropdown) => {
            dropdown.addOption("", `\u4F7F\u7528\u9810\u8A2D${defaultBackend ? ` (${defaultBackend})` : ""}`);
            backends.forEach((name) => {
              dropdown.addOption(name, name);
            });
            dropdown.setValue(this.plugin.settings.backend || "").onChange(async (value) => {
              this.plugin.settings.backend = value;
              await this.plugin.saveSettings();
            });
          });
        } else {
          backendSetting.addText(
            (text) => text.setPlaceholder("openrouter").setValue(this.plugin.settings.backend).onChange(async (value) => {
              this.plugin.settings.backend = value;
              await this.plugin.saveSettings();
            })
          );
        }
      }
      renderBackendInfo(containerEl) {
        const infoEl = containerEl.createDiv({ cls: "novelmaker-api-warning" });
        infoEl.createEl("span", { text: "\u2139\uFE0F \u63D0\u793A\uFF1A ", cls: "novelmaker-warning-icon" });
        infoEl.createEl("span", {
          text: "Backend \u914D\u7F6E\uFF08\u5305\u62EC API Key\uFF09\u5132\u5B58\u5728 ~/.novelmaker/config.toml \u4E2D\uFF0C\u4E0D\u6703\u540C\u6B65\u5230 vault\u3002\u6BCF\u500B backend \u53EF\u4EE5\u8A2D\u5B9A\u81EA\u5DF1\u7684\u9810\u8A2D\u6A21\u578B\u548C\u8D85\u6642\u6642\u9593\u3002",
          cls: "novelmaker-warning-text"
        });
      }
      renderGenerationSettings(containerEl) {
        new Setting(containerEl).setName("\u751F\u6210\u5F8C\u81EA\u52D5\u6253\u958B\u6A94\u6848").setDesc("\u751F\u6210\u7AE0\u7BC0\u6216\u89D2\u8272\u5F8C\uFF0C\u81EA\u52D5\u5728 Obsidian \u4E2D\u6253\u958B\u751F\u6210\u7684\u6A94\u6848").addToggle(
          (toggle) => toggle.setValue(this.plugin.settings.openAfterGen).onChange(async (value) => {
            this.plugin.settings.openAfterGen = value;
            await this.plugin.saveSettings();
          })
        );
        this.openAfterGenMsSetting = new Setting(containerEl).setName(`\u751F\u6210\u5F8C\u6253\u958B\u6A94\u6848\u5EF6\u9072\u6642\u9593 (${this.plugin.settings.openAfterGenMs} \u6BEB\u79D2)`).setDesc("\u751F\u6210\u7AE0\u7BC0\u6216\u89D2\u8272\u5F8C\uFF0C\u81EA\u52D5\u6253\u958B\u6A94\u6848\u524D\u7684\u5EF6\u9072\u6642\u9593\uFF08\u6BEB\u79D2\uFF09\uFF08\u9810\u8A2D\uFF1A500 \u6BEB\u79D2\uFF09").addSlider(
          (slider) => slider.setLimits(0, 1e4, 500).setValue(this.plugin.settings.openAfterGenMs).onChange(async (value) => {
            if (isNaN(value)) {
              return;
            }
            this.plugin.settings.openAfterGenMs = value;
            this.openAfterGenMsSetting.setName(`\u751F\u6210\u5F8C\u6253\u958B\u6A94\u6848\u5EF6\u9072\u6642\u9593 (${value} \u6BEB\u79D2)`);
            await this.plugin.saveSettings();
          })
        );
      }
    };
    module2.exports = NovelMakerSettingTab2;
  }
});

// src/modals/result.js
var require_result = __commonJS({
  "src/modals/result.js"(exports2, module2) {
    var { Modal, Setting } = require("obsidian");
    var ResultModal = class extends Modal {
      constructor(app, title, output) {
        super(app);
        this.modalTitle = title;
        this.output = output;
      }
      onOpen() {
        const { contentEl } = this;
        contentEl.createEl("h2", { text: this.modalTitle });
        const successDiv = contentEl.createDiv({ cls: "novelmaker-result-success" });
        successDiv.createEl("span", { text: "\u2705 ", cls: "novelmaker-success-icon" });
        successDiv.createEl("span", { text: "\u751F\u6210\u6210\u529F\uFF01" });
        if (this.output.filepath) {
          const fileDiv = contentEl.createDiv({ cls: "novelmaker-result-info" });
          fileDiv.createEl("strong", { text: "\u6A94\u6848\uFF1A" });
          fileDiv.createEl("span", { text: this.output.filepath });
        }
        if (this.output.input_tokens !== void 0 || this.output.output_tokens !== void 0) {
          const usageDiv = contentEl.createDiv({ cls: "novelmaker-token-usage" });
          usageDiv.createEl("h3", { text: "\u{1F522} Token \u4F7F\u7528\u91CF" });
          const usageGrid = usageDiv.createDiv({ cls: "novelmaker-token-grid" });
          if (this.output.input_tokens !== void 0) {
            const inputRow = usageGrid.createDiv({ cls: "novelmaker-token-row" });
            inputRow.createEl("span", { text: "Input tokens:", cls: "novelmaker-token-label" });
            inputRow.createEl("span", { text: this.output.input_tokens.toLocaleString(), cls: "novelmaker-token-value" });
          }
          if (this.output.output_tokens !== void 0) {
            const outputRow = usageGrid.createDiv({ cls: "novelmaker-token-row" });
            outputRow.createEl("span", { text: "Output tokens:", cls: "novelmaker-token-label" });
            outputRow.createEl("span", { text: this.output.output_tokens.toLocaleString(), cls: "novelmaker-token-value" });
          }
          if (this.output.total_tokens !== void 0) {
            const totalRow = usageGrid.createDiv({ cls: "novelmaker-token-row novelmaker-token-total" });
            totalRow.createEl("span", { text: "Total tokens:", cls: "novelmaker-token-label" });
            totalRow.createEl("span", { text: this.output.total_tokens.toLocaleString(), cls: "novelmaker-token-value" });
          }
        }
        new Setting(contentEl).addButton(
          (btn) => btn.setButtonText("\u78BA\u5B9A").setCta().onClick(() => {
            this.close();
          })
        );
      }
      onClose() {
        const { contentEl } = this;
        contentEl.empty();
      }
    };
    module2.exports = ResultModal;
  }
});

// src/modals/gen-next.js
var require_gen_next = __commonJS({
  "src/modals/gen-next.js"(exports2, module2) {
    var { Modal, Notice, Setting } = require("obsidian");
    var GenNextModal = class extends Modal {
      constructor(app, onSubmit) {
        super(app);
        this.title = "";
        this.prompt = "";
        this.prevCount = 3;
        this.onSubmit = onSubmit;
        this.prevSetting = null;
      }
      onOpen() {
        const { contentEl } = this;
        contentEl.createEl("h2", { text: "\u751F\u6210\u4E0B\u4E00\u7AE0" });
        new Setting(contentEl).setName("\u7AE0\u7BC0\u6A19\u984C").setDesc("\u4E0B\u4E00\u7AE0\u7684\u6A19\u984C\uFF08\u5FC5\u586B\uFF09").addText(
          (text) => text.setPlaceholder("e.g., \u7B2C3\u7AE0").onChange((value) => {
            this.title = value;
          })
        );
        this.prevSetting = new Setting(contentEl).setName("\u524D\u5E7E\u7AE0\u6578\u91CF (\u524D3\u7AE0)").setDesc("\u8981\u5305\u542B\u591A\u5C11\u524D\u9762\u7684\u7AE0\u7BC0\u4F5C\u70BA\u4E0A\u4E0B\u6587\uFF08\u9810\u8A2D\uFF1A3\uFF0C\u6700\u5927\uFF1A10\uFF09").addSlider(
          (text) => text.setValue(3).setLimits(0, 10, 1).onChange((num) => {
            if (isNaN(num)) {
              return;
            }
            this.prevCount = num;
            this.prevSetting.setName(`\u524D\u5E7E\u7AE0\u6578\u91CF (\u524D${num}\u7AE0)`);
          })
        );
        new Setting(contentEl).setName("\u63D0\u793A\u8A5E").setDesc("\u7AE0\u7BC0\u751F\u6210\u7684\u984D\u5916\u6307\u793A\uFF08\u9078\u586B\uFF09").addTextArea(
          (text) => text.setPlaceholder("\u4F8B\u5982\uFF1A\u8457\u91CD\u65BC\u89D2\u8272\u767C\u5C55...").onChange((value) => {
            this.prompt = value;
          })
        );
        new Setting(contentEl).addButton(
          (btn) => btn.setButtonText("\u53D6\u6D88").onClick(() => {
            this.close();
          })
        ).addButton(
          (btn) => btn.setButtonText("\u751F\u6210").setCta().onClick(() => {
            if (!this.title || !this.title.trim()) {
              new Notice("\u274C \u8ACB\u8F38\u5165\u7AE0\u7BC0\u6A19\u984C");
              return;
            }
            this.close();
            this.onSubmit(this.title, this.prompt, this.prevCount);
          })
        );
      }
      onClose() {
        const { contentEl } = this;
        contentEl.empty();
      }
    };
    module2.exports = GenNextModal;
  }
});

// src/modals/gen-next-empty.js
var require_gen_next_empty = __commonJS({
  "src/modals/gen-next-empty.js"(exports2, module2) {
    var { Modal, Notice, Setting } = require("obsidian");
    var GenNextEmptyModal = class extends Modal {
      constructor(app, onSubmit) {
        super(app);
        this.title = "";
        this.prompt = "";
        this.onSubmit = onSubmit;
      }
      onOpen() {
        const { contentEl } = this;
        contentEl.createEl("h2", { text: "\u751F\u6210\u7A7A\u767D\u4E0B\u4E00\u7AE0" });
        new Setting(contentEl).setName("\u7AE0\u7BC0\u6A19\u984C").setDesc("\u4E0B\u4E00\u7AE0\u7684\u6A19\u984C\uFF08\u5FC5\u586B\uFF09").addText(
          (text) => text.setPlaceholder("e.g., \u7B2C3\u7AE0").onChange((value) => {
            this.title = value;
          })
        );
        new Setting(contentEl).setName("\u63D0\u793A\u8A5E").setDesc("\u7AE0\u7BC0\u7684\u5099\u8A3B\u6216\u63D0\u793A\uFF08\u9078\u586B\uFF09").addTextArea(
          (text) => text.setPlaceholder("\u4F8B\u5982\uFF1A\u8A08\u5283\u5BEB\u4E00\u500B\u6230\u9B25\u5834\u666F...").onChange((value) => {
            this.prompt = value;
          })
        );
        new Setting(contentEl).addButton(
          (btn) => btn.setButtonText("\u53D6\u6D88").onClick(() => {
            this.close();
          })
        ).addButton(
          (btn) => btn.setButtonText("\u5EFA\u7ACB").setCta().onClick(() => {
            if (!this.title || !this.title.trim()) {
              new Notice("\u274C \u8ACB\u8F38\u5165\u7AE0\u7BC0\u6A19\u984C");
              return;
            }
            this.close();
            this.onSubmit(this.title, this.prompt);
          })
        );
      }
      onClose() {
        const { contentEl } = this;
        contentEl.empty();
      }
    };
    module2.exports = GenNextEmptyModal;
  }
});

// src/modals/gen-curr.js
var require_gen_curr = __commonJS({
  "src/modals/gen-curr.js"(exports2, module2) {
    var { Modal, Setting } = require("obsidian");
    var GenCurrModal = class extends Modal {
      constructor(app, activeFile, onSubmit) {
        super(app);
        this.activeFile = activeFile;
        this.prevCount = 3;
        this.onSubmit = onSubmit;
        this.prevSetting = null;
      }
      onOpen() {
        const { contentEl } = this;
        contentEl.createEl("h2", { text: "\u91CD\u65B0\u751F\u6210\u7AE0\u7BC0" });
        new Setting(contentEl).setName("\u76EE\u6A19\u6A94\u6848").setDesc("\u5C07\u91CD\u65B0\u751F\u6210\u6B64\u6A94\u6848\u7684\u5167\u5BB9").addText((text) => {
          text.setValue(this.activeFile.path);
          text.inputEl.disabled = true;
          text.inputEl.style.width = "100%";
        });
        contentEl.createEl("p", {
          text: "\u26A0\uFE0F \u6B64\u64CD\u4F5C\u5C07\u4F7F\u7528\u6A94\u6848\u4E2D\u7684 prompt \u6B04\u4F4D\u91CD\u65B0\u751F\u6210\u7AE0\u7BC0\u5167\u5BB9\uFF0C\u4E26\u8986\u84CB\u73FE\u6709\u5167\u5BB9\u3002",
          cls: "mod-warning"
        });
        this.prevSetting = new Setting(contentEl).setName("\u524D\u5E7E\u7AE0\u6578\u91CF (\u524D3\u7AE0)").setDesc("\u8981\u5305\u542B\u591A\u5C11\u524D\u9762\u7684\u7AE0\u7BC0\u4F5C\u70BA\u4E0A\u4E0B\u6587\uFF08\u9810\u8A2D\uFF1A3\uFF0C\u6700\u5927\uFF1A10\uFF09").addSlider(
          (slider) => slider.setValue(3).setLimits(0, 10, 1).onChange((num) => {
            if (isNaN(num)) {
              return;
            }
            this.prevCount = num;
            this.prevSetting.setName(`\u524D\u5E7E\u7AE0\u6578\u91CF (\u524D${num}\u7AE0)`);
          })
        );
        new Setting(contentEl).addButton(
          (btn) => btn.setButtonText("\u53D6\u6D88").onClick(() => {
            this.close();
          })
        ).addButton(
          (btn) => btn.setButtonText("\u91CD\u65B0\u751F\u6210").setCta().onClick(() => {
            this.close();
            this.onSubmit(this.activeFile.path, this.prevCount);
          })
        );
      }
      onClose() {
        const { contentEl } = this;
        contentEl.empty();
      }
    };
    module2.exports = GenCurrModal;
  }
});

// src/commands/chapter.js
var require_chapter = __commonJS({
  "src/commands/chapter.js"(exports2, module2) {
    var { Notice } = require("obsidian");
    var LoadingModal = require_loading();
    var ResultModal = require_result();
    var GenNextModal = require_gen_next();
    var GenNextEmptyModal = require_gen_next_empty();
    var GenCurrModal = require_gen_curr();
    var { executeCLI, parseJSONOutput, openGeneratedFile, buildCLICommand } = require_cli();
    function registerChapterCommands2(plugin) {
      registerGenNextCommand(plugin);
      registerGenNextEmptyCommand(plugin);
      registerGenCurrCommand(plugin);
    }
    function registerGenNextCommand(plugin) {
      plugin.addCommand({
        id: "gen-next-chapter",
        name: "\u751F\u6210\u4E0B\u4E00\u7AE0",
        callback: () => {
          new GenNextModal(plugin.app, async (title, prompt, prevCount) => {
            const loadingModal = new LoadingModal(plugin.app, "\u6B63\u5728\u751F\u6210\u4E0B\u4E00\u7AE0...\u8ACB\u7A0D\u5019");
            loadingModal.open();
            try {
              const vaultPath = plugin.app.vault.adapter.basePath;
              const args = buildCLICommand("chapter gen-next", {
                json: true,
                title,
                prompt,
                prevCount,
                backend: plugin.settings.backend
              });
              const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
              const output = parseJSONOutput(stdout);
              openGeneratedFile(plugin.app, plugin.settings, output.filepath);
              new ResultModal(plugin.app, "\u7AE0\u7BC0\u751F\u6210\u5B8C\u6210", output).open();
            } catch (error) {
              new Notice(`\u274C \u932F\u8AA4: ${error.message}`);
              console.error(error);
            } finally {
              loadingModal.forceCloseNow();
            }
          }).open();
        }
      });
    }
    function registerGenNextEmptyCommand(plugin) {
      plugin.addCommand({
        id: "gen-next-empty-chapter",
        name: "\u751F\u6210\u7A7A\u767D\u4E0B\u4E00\u7AE0",
        callback: () => {
          new GenNextEmptyModal(plugin.app, async (title, prompt) => {
            const loadingModal = new LoadingModal(plugin.app, "\u6B63\u5728\u5EFA\u7ACB\u7A7A\u767D\u7AE0\u7BC0...\u8ACB\u7A0D\u5019");
            loadingModal.open();
            try {
              const vaultPath = plugin.app.vault.adapter.basePath;
              const args = buildCLICommand("chapter gen-empty", {
                json: true,
                title,
                prompt
              });
              const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
              const output = parseJSONOutput(stdout);
              openGeneratedFile(plugin.app, plugin.settings, output.filepath);
              new Notice("\u2705 \u7A7A\u767D\u7AE0\u7BC0\u5EFA\u7ACB\u6210\u529F\uFF01");
            } catch (error) {
              new Notice(`\u274C \u932F\u8AA4: ${error.message}`);
              console.error(error);
            } finally {
              loadingModal.forceCloseNow();
            }
          }).open();
        }
      });
    }
    function registerGenCurrCommand(plugin) {
      plugin.addCommand({
        id: "gen-curr-chapter",
        name: "\u91CD\u65B0\u751F\u6210\u7576\u524D\u7AE0\u7BC0",
        checkCallback: (checking) => {
          const activeFile = plugin.app.workspace.getActiveFile();
          if (activeFile && activeFile.path.startsWith("Story/") && activeFile.extension === "md") {
            if (!checking) {
              new GenCurrModal(plugin.app, activeFile, async (filepath, prevCount) => {
                const loadingModal = new LoadingModal(plugin.app, "\u6B63\u5728\u91CD\u65B0\u751F\u6210\u7AE0\u7BC0...\u8ACB\u7A0D\u5019");
                loadingModal.open();
                try {
                  const vaultPath = plugin.app.vault.adapter.basePath;
                  const args = buildCLICommand("chapter regen", {
                    json: true,
                    filepath,
                    prevCount,
                    backend: plugin.settings.backend,
                    timeout: plugin.settings.timeout
                  });
                  const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
                  try {
                    const output = parseJSONOutput(stdout);
                    new ResultModal(plugin.app, "\u7AE0\u7BC0\u91CD\u65B0\u751F\u6210\u5B8C\u6210", output).open();
                  } catch (jsonError) {
                    new Notice("\u2705 \u7AE0\u7BC0\u91CD\u65B0\u751F\u6210\u6210\u529F\uFF01");
                  }
                } catch (error) {
                  new Notice(`\u274C \u932F\u8AA4: ${error.message}`);
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
    module2.exports = {
      registerChapterCommands: registerChapterCommands2
    };
  }
});

// src/modals/gen-char.js
var require_gen_char = __commonJS({
  "src/modals/gen-char.js"(exports2, module2) {
    var { Modal, Notice, Setting } = require("obsidian");
    var GenCharModal = class extends Modal {
      constructor(app, onSubmit) {
        super(app);
        this.name = "";
        this.prompt = "";
        this.onSubmit = onSubmit;
      }
      onOpen() {
        const { contentEl } = this;
        contentEl.createEl("h2", { text: "\u751F\u6210\u89D2\u8272" });
        new Setting(contentEl).setName("\u89D2\u8272\u540D\u7A31").setDesc("\u89D2\u8272\u7684\u540D\u7A31\uFF08\u5FC5\u586B\uFF09").addText(
          (text) => text.setPlaceholder("e.g., \u827E\u8389\u5A1C").onChange((value) => {
            this.name = value;
          })
        );
        new Setting(contentEl).setName("\u89D2\u8272\u63CF\u8FF0").setDesc("\u89D2\u8272\u7684\u7279\u5FB5\u6216\u80CC\u666F\u63CF\u8FF0\uFF08\u9078\u586B\uFF09").addTextArea(
          (text) => text.setPlaceholder("\u4F8B\u5982\uFF1A\u4E00\u4F4D\u5E74\u8F15\u7684\u9B54\u6CD5\u5E2B\uFF0C\u64C5\u9577\u5143\u7D20\u9B54\u6CD5...").onChange((value) => {
            this.prompt = value;
          })
        );
        new Setting(contentEl).addButton(
          (btn) => btn.setButtonText("\u53D6\u6D88").onClick(() => {
            this.close();
          })
        ).addButton(
          (btn) => btn.setButtonText("\u751F\u6210").setCta().onClick(() => {
            if (!this.name || !this.name.trim()) {
              new Notice("\u274C \u8ACB\u8F38\u5165\u89D2\u8272\u540D\u7A31");
              return;
            }
            this.close();
            this.onSubmit(this.name, this.prompt);
          })
        );
      }
      onClose() {
        const { contentEl } = this;
        contentEl.empty();
      }
    };
    module2.exports = GenCharModal;
  }
});

// src/modals/gen-char-curr.js
var require_gen_char_curr = __commonJS({
  "src/modals/gen-char-curr.js"(exports2, module2) {
    var { Modal, Setting } = require("obsidian");
    var GenCharCurrModal = class extends Modal {
      constructor(app, activeFile, onSubmit) {
        super(app);
        this.activeFile = activeFile;
        this.onSubmit = onSubmit;
      }
      onOpen() {
        const { contentEl } = this;
        contentEl.createEl("h2", { text: "\u91CD\u65B0\u751F\u6210\u89D2\u8272" });
        new Setting(contentEl).setName("\u76EE\u6A19\u6A94\u6848").setDesc("\u5C07\u91CD\u65B0\u751F\u6210\u6B64\u6A94\u6848\u7684\u5167\u5BB9").addText((text) => {
          text.setValue(this.activeFile.path);
          text.inputEl.disabled = true;
          text.inputEl.style.width = "100%";
        });
        contentEl.createEl("p", {
          text: "\u26A0\uFE0F \u6B64\u64CD\u4F5C\u5C07\u4F7F\u7528\u6A94\u6848\u4E2D\u7684 prompt \u6B04\u4F4D\u91CD\u65B0\u751F\u6210\u89D2\u8272\u8CC7\u6599\uFF0C\u4E26\u8986\u84CB\u73FE\u6709\u5167\u5BB9\u3002",
          cls: "mod-warning"
        });
        new Setting(contentEl).addButton(
          (btn) => btn.setButtonText("\u53D6\u6D88").onClick(() => {
            this.close();
          })
        ).addButton(
          (btn) => btn.setButtonText("\u91CD\u65B0\u751F\u6210").setCta().onClick(() => {
            this.close();
            this.onSubmit(this.activeFile.path);
          })
        );
      }
      onClose() {
        const { contentEl } = this;
        contentEl.empty();
      }
    };
    module2.exports = GenCharCurrModal;
  }
});

// src/modals/gen-char-img.js
var require_gen_char_img = __commonJS({
  "src/modals/gen-char-img.js"(exports2, module2) {
    var { Modal, Notice, Setting } = require("obsidian");
    var GenCharImgModal = class extends Modal {
      constructor(app, activeFile, onSubmit) {
        super(app);
        this.activeFile = activeFile;
        this.name = "";
        this.prompt = "";
        this.onSubmit = onSubmit;
      }
      async onOpen() {
        const { contentEl } = this;
        contentEl.createEl("h2", { text: "\u751F\u6210\u89D2\u8272\u5716\u7247" });
        new Setting(contentEl).setName("\u76EE\u6A19\u6A94\u6848").setDesc("\u5C07\u70BA\u6B64\u89D2\u8272\u751F\u6210\u5716\u7247").addText((text) => {
          text.setValue(this.activeFile.path);
          text.inputEl.disabled = true;
          text.inputEl.style.width = "100%";
        });
        try {
          const content = await this.app.vault.read(this.activeFile);
          const frontmatterMatch = content.match(/^---\n([\s\S]*?)\n---/);
          if (frontmatterMatch) {
            const frontmatter = frontmatterMatch[1];
            const nameMatch = frontmatter.match(/^name:\s*(.+)$/m);
            if (nameMatch) {
              this.name = nameMatch[1].trim();
            }
          }
        } catch (error) {
          console.error("Failed to read character name:", error);
        }
        new Setting(contentEl).setName("\u89D2\u8272\u540D\u7A31").setDesc("\u89D2\u8272\u7684\u540D\u7A31\uFF08\u5F9E\u6A94\u6848\u8B80\u53D6\uFF09").addText((text) => {
          text.setValue(this.name);
          text.inputEl.disabled = true;
          text.inputEl.style.width = "100%";
        });
        new Setting(contentEl).setName("\u5716\u7247\u63CF\u8FF0").setDesc("\u81EA\u8A02\u5716\u7247\u751F\u6210\u7684\u63D0\u793A\u8A5E\uFF08\u9078\u586B\uFF0C\u7559\u7A7A\u5247\u4F7F\u7528\u89D2\u8272\u6A94\u6848\u4E2D\u7684\u63CF\u8FF0\uFF09").addTextArea(
          (text) => text.setPlaceholder("\u4F8B\u5982\uFF1A\u4E00\u4F4D\u5E74\u8F15\u7684\u9B54\u6CD5\u5E2B\u8096\u50CF\uFF0C\u80CC\u666F\u662F\u9B54\u6CD5\u5B78\u9662...").onChange((value) => {
            this.prompt = value;
          })
        );
        contentEl.createEl("p", {
          text: "\u{1F4A1} \u63D0\u793A\uFF1A\u5716\u7247\u5C07\u4F7F\u7528 DALL-E API \u751F\u6210\u4E26\u5132\u5B58\u81F3 Character/ \u76EE\u9304",
          cls: "mod-info"
        });
        new Setting(contentEl).addButton(
          (btn) => btn.setButtonText("\u53D6\u6D88").onClick(() => {
            this.close();
          })
        ).addButton(
          (btn) => btn.setButtonText("\u751F\u6210\u5716\u7247").setCta().onClick(() => {
            if (!this.name || !this.name.trim()) {
              new Notice("\u274C \u7121\u6CD5\u53D6\u5F97\u89D2\u8272\u540D\u7A31");
              return;
            }
            this.close();
            this.onSubmit(this.name, this.prompt);
          })
        );
      }
      onClose() {
        const { contentEl } = this;
        contentEl.empty();
      }
    };
    module2.exports = GenCharImgModal;
  }
});

// src/commands/character.js
var require_character = __commonJS({
  "src/commands/character.js"(exports2, module2) {
    var { Notice } = require("obsidian");
    var LoadingModal = require_loading();
    var ResultModal = require_result();
    var GenCharModal = require_gen_char();
    var GenCharCurrModal = require_gen_char_curr();
    var GenCharImgModal = require_gen_char_img();
    var { executeCLI, parseJSONOutput, openGeneratedFile, buildCLICommand } = require_cli();
    function registerCharacterCommands2(plugin) {
      registerGenCharCommand(plugin);
      registerGenCharCurrCommand(plugin);
      registerGenCharImgCommand(plugin);
    }
    function registerGenCharCommand(plugin) {
      plugin.addCommand({
        id: "gen-character",
        name: "\u751F\u6210\u89D2\u8272",
        callback: () => {
          new GenCharModal(plugin.app, async (name, prompt) => {
            const loadingModal = new LoadingModal(plugin.app, "\u6B63\u5728\u751F\u6210\u89D2\u8272...\u8ACB\u7A0D\u5019");
            loadingModal.open();
            try {
              const vaultPath = plugin.app.vault.adapter.basePath;
              const args = buildCLICommand("character gen", {
                json: true,
                name,
                prompt,
                backend: plugin.settings.backend
              });
              const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
              const output = parseJSONOutput(stdout);
              openGeneratedFile(plugin.app, plugin.settings, output.filepath);
              new ResultModal(plugin.app, "\u89D2\u8272\u751F\u6210\u5B8C\u6210", output).open();
            } catch (error) {
              new Notice(`\u274C \u932F\u8AA4: ${error.message}`);
              console.error(error);
            } finally {
              loadingModal.forceCloseNow();
            }
          }).open();
        }
      });
    }
    function registerGenCharCurrCommand(plugin) {
      plugin.addCommand({
        id: "gen-char-curr",
        name: "\u91CD\u65B0\u751F\u6210\u7576\u524D\u89D2\u8272",
        checkCallback: (checking) => {
          const activeFile = plugin.app.workspace.getActiveFile();
          if (activeFile && activeFile.path.startsWith("Character/") && activeFile.extension === "md") {
            if (!checking) {
              new GenCharCurrModal(plugin.app, activeFile, async (filepath) => {
                const loadingModal = new LoadingModal(plugin.app, "\u6B63\u5728\u91CD\u65B0\u751F\u6210\u89D2\u8272...\u8ACB\u7A0D\u5019");
                loadingModal.open();
                try {
                  const vaultPath = plugin.app.vault.adapter.basePath;
                  const args = buildCLICommand("character regen", {
                    json: true,
                    filepath,
                    backend: plugin.settings.backend
                  });
                  const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
                  try {
                    const output = parseJSONOutput(stdout);
                    new ResultModal(plugin.app, "\u89D2\u8272\u91CD\u65B0\u751F\u6210\u5B8C\u6210", output).open();
                  } catch (jsonError) {
                    new Notice("\u2705 \u89D2\u8272\u91CD\u65B0\u751F\u6210\u6210\u529F\uFF01");
                  }
                } catch (error) {
                  new Notice(`\u274C \u932F\u8AA4: ${error.message}`);
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
    function registerGenCharImgCommand(plugin) {
      plugin.addCommand({
        id: "gen-char-img",
        name: "\u751F\u6210\u89D2\u8272\u5716\u7247",
        checkCallback: (checking) => {
          const activeFile = plugin.app.workspace.getActiveFile();
          if (activeFile && activeFile.path.startsWith("Character/") && activeFile.extension === "md") {
            if (!checking) {
              new GenCharImgModal(plugin.app, activeFile, async (name, prompt) => {
                const loadingModal = new LoadingModal(plugin.app, "\u6B63\u5728\u751F\u6210\u89D2\u8272\u5716\u7247...\u8ACB\u7A0D\u5019");
                loadingModal.open();
                try {
                  const vaultPath = plugin.app.vault.adapter.basePath;
                  const args = buildCLICommand("character gen-img", {
                    json: true,
                    name,
                    prompt,
                    backend: plugin.settings.backend
                  });
                  const { stdout } = await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
                  const output = parseJSONOutput(stdout);
                  if (output && output.filepath) {
                    new Notice(`\u2705 \u89D2\u8272\u5716\u7247\u751F\u6210\u6210\u529F\uFF01
\u6A94\u6848\uFF1A${output.filepath}`);
                  } else {
                    new Notice("\u2705 \u89D2\u8272\u5716\u7247\u751F\u6210\u6210\u529F\uFF01");
                  }
                } catch (error) {
                  new Notice(`\u274C \u932F\u8AA4: ${error.message}`);
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
    module2.exports = {
      registerCharacterCommands: registerCharacterCommands2
    };
  }
});

// src/modals/export.js
var require_export = __commonJS({
  "src/modals/export.js"(exports2, module2) {
    var { Modal, Notice, Setting } = require("obsidian");
    var ExportModal = class extends Modal {
      constructor(app, onSubmit) {
        super(app);
        this.outputPath = "";
        this.onSubmit = onSubmit;
      }
      onOpen() {
        const { contentEl } = this;
        contentEl.createEl("h2", { text: "\u532F\u51FA\u5C0F\u8AAA" });
        const pathSetting = new Setting(contentEl).setName("\u8F38\u51FA\u6A94\u6848\u8DEF\u5F91").setDesc("\u9078\u64C7\u5C0F\u8AAA\u532F\u51FA\u7684\u6A94\u6848\u8DEF\u5F91").addText((text) => {
          text.setPlaceholder("\u9078\u64C7\u6A94\u6848\u4F4D\u7F6E...").setValue(this.outputPath).onChange((value) => {
            this.outputPath = value;
          });
          text.inputEl.style.width = "100%";
          this.pathTextComponent = text;
        }).addButton(
          (btn) => btn.setButtonText("\u700F\u89BD...").onClick(async () => {
            try {
              const remote = require("@electron/remote");
              const dialog = remote.dialog;
              const result = await dialog.showSaveDialog({
                title: "\u9078\u64C7\u532F\u51FA\u4F4D\u7F6E",
                defaultPath: "novel.txt",
                filters: [
                  { name: "\u6587\u5B57\u6A94\u6848", extensions: ["txt"] },
                  { name: "\u6240\u6709\u6A94\u6848", extensions: ["*"] }
                ],
                properties: ["createDirectory", "showOverwriteConfirmation"]
              });
              if (!result.canceled && result.filePath) {
                this.outputPath = result.filePath;
                if (this.pathTextComponent) {
                  this.pathTextComponent.setValue(result.filePath);
                }
              }
            } catch (error) {
              console.error("File dialog error:", error);
              new Notice("\u26A0 \u7121\u6CD5\u958B\u555F\u6A94\u6848\u5C0D\u8A71\u6846\uFF0C\u8ACB\u624B\u52D5\u8F38\u5165\u8DEF\u5F91");
            }
          })
        );
        new Setting(contentEl).addButton(
          (btn) => btn.setButtonText("\u53D6\u6D88").onClick(() => {
            this.close();
          })
        ).addButton(
          (btn) => btn.setButtonText("\u532F\u51FA").setCta().onClick(() => {
            if (!this.outputPath || !this.outputPath.trim()) {
              new Notice("\u274C \u8ACB\u9078\u64C7\u6216\u8F38\u5165\u8F38\u51FA\u6A94\u6848\u8DEF\u5F91");
              return;
            }
            this.close();
            this.onSubmit(this.outputPath);
          })
        );
      }
      onClose() {
        const { contentEl } = this;
        contentEl.empty();
      }
    };
    module2.exports = ExportModal;
  }
});

// src/commands/export.js
var require_export2 = __commonJS({
  "src/commands/export.js"(exports2, module2) {
    var { Notice } = require("obsidian");
    var LoadingModal = require_loading();
    var ExportModal = require_export();
    var { executeCLI, buildCLICommand } = require_cli();
    function registerExportCommand2(plugin) {
      plugin.addCommand({
        id: "export-novel",
        name: "\u532F\u51FA\u5C0F\u8AAA",
        callback: () => {
          new ExportModal(plugin.app, async (outputPath) => {
            const loadingModal = new LoadingModal(plugin.app, "\u6B63\u5728\u532F\u51FA\u5C0F\u8AAA...\u8ACB\u7A0D\u5019");
            loadingModal.open();
            try {
              const vaultPath = plugin.app.vault.adapter.basePath;
              const args = buildCLICommand("export", {
                output: outputPath,
                type: "txt"
              });
              await executeCLI(plugin, plugin.settings.cliPath, args, vaultPath);
              new Notice("\u2705 \u5C0F\u8AAA\u532F\u51FA\u6210\u529F\uFF01");
            } catch (error) {
              new Notice(`\u274C \u932F\u8AA4: ${error.message}`);
              console.error(error);
            } finally {
              loadingModal.forceCloseNow();
            }
          }).open();
        }
      });
    }
    module2.exports = {
      registerExportCommand: registerExportCommand2
    };
  }
});

// main.js
var { Plugin } = require("obsidian");
var DEFAULT_SETTINGS = require_constants();
var NovelMakerSettingTab = require_settings_tab();
var { registerChapterCommands } = require_chapter();
var { registerCharacterCommands } = require_character();
var { registerExportCommand } = require_export2();
var NovelMakerPlugin = class extends Plugin {
  async onload() {
    await this.loadSettings();
    this.addSettingTab(new NovelMakerSettingTab(this.app, this));
    registerChapterCommands(this);
    registerCharacterCommands(this);
    registerExportCommand(this);
  }
  onunload() {
  }
  async loadSettings() {
    this.settings = Object.assign({}, DEFAULT_SETTINGS, await this.loadData());
  }
  async saveSettings() {
    await this.saveData(this.settings);
  }
};
module.exports = NovelMakerPlugin;
