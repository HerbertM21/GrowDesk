(function(global2, factory) {
  typeof exports === "object" && typeof module !== "undefined" ? module.exports = factory(require("vue")) : typeof define === "function" && define.amd ? define(["vue"], factory) : (global2 = typeof globalThis !== "undefined" ? globalThis : global2 || self, global2.GrowDeskWidget = factory(global2.Vue));
})(this, function(vue) {
  "use strict";
  var _a;
  function getDevtoolsGlobalHook() {
    return getTarget().__VUE_DEVTOOLS_GLOBAL_HOOK__;
  }
  function getTarget() {
    return typeof navigator !== "undefined" && typeof window !== "undefined" ? window : typeof globalThis !== "undefined" ? globalThis : {};
  }
  const isProxyAvailable = typeof Proxy === "function";
  const HOOK_SETUP = "devtools-plugin:setup";
  const HOOK_PLUGIN_SETTINGS_SET = "plugin:settings:set";
  let supported;
  let perf;
  function isPerformanceSupported() {
    var _a2;
    if (supported !== void 0) {
      return supported;
    }
    if (typeof window !== "undefined" && window.performance) {
      supported = true;
      perf = window.performance;
    } else if (typeof globalThis !== "undefined" && ((_a2 = globalThis.perf_hooks) === null || _a2 === void 0 ? void 0 : _a2.performance)) {
      supported = true;
      perf = globalThis.perf_hooks.performance;
    } else {
      supported = false;
    }
    return supported;
  }
  function now() {
    return isPerformanceSupported() ? perf.now() : Date.now();
  }
  class ApiProxy {
    constructor(plugin, hook) {
      this.target = null;
      this.targetQueue = [];
      this.onQueue = [];
      this.plugin = plugin;
      this.hook = hook;
      const defaultSettings = {};
      if (plugin.settings) {
        for (const id in plugin.settings) {
          const item = plugin.settings[id];
          defaultSettings[id] = item.defaultValue;
        }
      }
      const localSettingsSaveId = `__vue-devtools-plugin-settings__${plugin.id}`;
      let currentSettings = Object.assign({}, defaultSettings);
      try {
        const raw = localStorage.getItem(localSettingsSaveId);
        const data = JSON.parse(raw);
        Object.assign(currentSettings, data);
      } catch (e) {
      }
      this.fallbacks = {
        getSettings() {
          return currentSettings;
        },
        setSettings(value) {
          try {
            localStorage.setItem(localSettingsSaveId, JSON.stringify(value));
          } catch (e) {
          }
          currentSettings = value;
        },
        now() {
          return now();
        }
      };
      if (hook) {
        hook.on(HOOK_PLUGIN_SETTINGS_SET, (pluginId, value) => {
          if (pluginId === this.plugin.id) {
            this.fallbacks.setSettings(value);
          }
        });
      }
      this.proxiedOn = new Proxy({}, {
        get: (_target, prop) => {
          if (this.target) {
            return this.target.on[prop];
          } else {
            return (...args) => {
              this.onQueue.push({
                method: prop,
                args
              });
            };
          }
        }
      });
      this.proxiedTarget = new Proxy({}, {
        get: (_target, prop) => {
          if (this.target) {
            return this.target[prop];
          } else if (prop === "on") {
            return this.proxiedOn;
          } else if (Object.keys(this.fallbacks).includes(prop)) {
            return (...args) => {
              this.targetQueue.push({
                method: prop,
                args,
                resolve: () => {
                }
              });
              return this.fallbacks[prop](...args);
            };
          } else {
            return (...args) => {
              return new Promise((resolve) => {
                this.targetQueue.push({
                  method: prop,
                  args,
                  resolve
                });
              });
            };
          }
        }
      });
    }
    async setRealTarget(target) {
      this.target = target;
      for (const item of this.onQueue) {
        this.target.on[item.method](...item.args);
      }
      for (const item of this.targetQueue) {
        item.resolve(await this.target[item.method](...item.args));
      }
    }
  }
  function setupDevtoolsPlugin(pluginDescriptor, setupFn) {
    const descriptor = pluginDescriptor;
    const target = getTarget();
    const hook = getDevtoolsGlobalHook();
    const enableProxy = isProxyAvailable && descriptor.enableEarlyProxy;
    if (hook && (target.__VUE_DEVTOOLS_PLUGIN_API_AVAILABLE__ || !enableProxy)) {
      hook.emit(HOOK_SETUP, pluginDescriptor, setupFn);
    } else {
      const proxy = enableProxy ? new ApiProxy(descriptor, hook) : null;
      const list = target.__VUE_DEVTOOLS_PLUGINS__ = target.__VUE_DEVTOOLS_PLUGINS__ || [];
      list.push({
        pluginDescriptor: descriptor,
        setupFn,
        proxy
      });
      if (proxy) {
        setupFn(proxy.proxiedTarget);
      }
    }
  }
  const piniaSymbol = Symbol("pinia");
  var MutationType;
  (function(MutationType2) {
    MutationType2["direct"] = "direct";
    MutationType2["patchObject"] = "patch object";
    MutationType2["patchFunction"] = "patch function";
  })(MutationType || (MutationType = {}));
  const IS_CLIENT = typeof window !== "undefined";
  const _global$1 = /* @__PURE__ */ (() => typeof window === "object" && window.window === window ? window : typeof self === "object" && self.self === self ? self : typeof global === "object" && global.global === global ? global : typeof globalThis === "object" ? globalThis : { HTMLElement: null })();
  function bom(blob, { autoBom = false } = {}) {
    if (autoBom && /^\s*(?:text\/\S*|application\/xml|\S*\/\S*\+xml)\s*;.*charset\s*=\s*utf-8/i.test(blob.type)) {
      return new Blob([String.fromCharCode(65279), blob], { type: blob.type });
    }
    return blob;
  }
  function download(url, name, opts) {
    const xhr = new XMLHttpRequest();
    xhr.open("GET", url);
    xhr.responseType = "blob";
    xhr.onload = function() {
      saveAs(xhr.response, name, opts);
    };
    xhr.onerror = function() {
      console.error("could not download file");
    };
    xhr.send();
  }
  function corsEnabled(url) {
    const xhr = new XMLHttpRequest();
    xhr.open("HEAD", url, false);
    try {
      xhr.send();
    } catch (e) {
    }
    return xhr.status >= 200 && xhr.status <= 299;
  }
  function click(node) {
    try {
      node.dispatchEvent(new MouseEvent("click"));
    } catch (e) {
      const evt = document.createEvent("MouseEvents");
      evt.initMouseEvent("click", true, true, window, 0, 0, 0, 80, 20, false, false, false, false, 0, null);
      node.dispatchEvent(evt);
    }
  }
  const _navigator$1 = typeof navigator === "object" ? navigator : { userAgent: "" };
  const isMacOSWebView = /* @__PURE__ */ (() => /Macintosh/.test(_navigator$1.userAgent) && /AppleWebKit/.test(_navigator$1.userAgent) && !/Safari/.test(_navigator$1.userAgent))();
  const saveAs = !IS_CLIENT ? () => {
  } : (
    // Use download attribute first if possible (#193 Lumia mobile) unless this is a macOS WebView or mini program
    typeof HTMLAnchorElement !== "undefined" && "download" in HTMLAnchorElement.prototype && !isMacOSWebView ? downloadSaveAs : (
      // Use msSaveOrOpenBlob as a second approach
      "msSaveOrOpenBlob" in _navigator$1 ? msSaveAs : (
        // Fallback to using FileReader and a popup
        fileSaverSaveAs
      )
    )
  );
  function downloadSaveAs(blob, name = "download", opts) {
    const a = document.createElement("a");
    a.download = name;
    a.rel = "noopener";
    if (typeof blob === "string") {
      a.href = blob;
      if (a.origin !== location.origin) {
        if (corsEnabled(a.href)) {
          download(blob, name, opts);
        } else {
          a.target = "_blank";
          click(a);
        }
      } else {
        click(a);
      }
    } else {
      a.href = URL.createObjectURL(blob);
      setTimeout(function() {
        URL.revokeObjectURL(a.href);
      }, 4e4);
      setTimeout(function() {
        click(a);
      }, 0);
    }
  }
  function msSaveAs(blob, name = "download", opts) {
    if (typeof blob === "string") {
      if (corsEnabled(blob)) {
        download(blob, name, opts);
      } else {
        const a = document.createElement("a");
        a.href = blob;
        a.target = "_blank";
        setTimeout(function() {
          click(a);
        });
      }
    } else {
      navigator.msSaveOrOpenBlob(bom(blob, opts), name);
    }
  }
  function fileSaverSaveAs(blob, name, opts, popup) {
    popup = popup || open("", "_blank");
    if (popup) {
      popup.document.title = popup.document.body.innerText = "downloading...";
    }
    if (typeof blob === "string")
      return download(blob, name, opts);
    const force = blob.type === "application/octet-stream";
    const isSafari = /constructor/i.test(String(_global$1.HTMLElement)) || "safari" in _global$1;
    const isChromeIOS = /CriOS\/[\d]+/.test(navigator.userAgent);
    if ((isChromeIOS || force && isSafari || isMacOSWebView) && typeof FileReader !== "undefined") {
      const reader = new FileReader();
      reader.onloadend = function() {
        let url = reader.result;
        if (typeof url !== "string") {
          popup = null;
          throw new Error("Wrong reader.result type");
        }
        url = isChromeIOS ? url : url.replace(/^data:[^;]*;/, "data:attachment/file;");
        if (popup) {
          popup.location.href = url;
        } else {
          location.assign(url);
        }
        popup = null;
      };
      reader.readAsDataURL(blob);
    } else {
      const url = URL.createObjectURL(blob);
      if (popup)
        popup.location.assign(url);
      else
        location.href = url;
      popup = null;
      setTimeout(function() {
        URL.revokeObjectURL(url);
      }, 4e4);
    }
  }
  function toastMessage(message, type) {
    const piniaMessage = "🍍 " + message;
    if (typeof __VUE_DEVTOOLS_TOAST__ === "function") {
      __VUE_DEVTOOLS_TOAST__(piniaMessage, type);
    } else if (type === "error") {
      console.error(piniaMessage);
    } else if (type === "warn") {
      console.warn(piniaMessage);
    } else {
      console.log(piniaMessage);
    }
  }
  function isPinia(o) {
    return "_a" in o && "install" in o;
  }
  function checkClipboardAccess() {
    if (!("clipboard" in navigator)) {
      toastMessage(`Your browser doesn't support the Clipboard API`, "error");
      return true;
    }
  }
  function checkNotFocusedError(error) {
    if (error instanceof Error && error.message.toLowerCase().includes("document is not focused")) {
      toastMessage('You need to activate the "Emulate a focused page" setting in the "Rendering" panel of devtools.', "warn");
      return true;
    }
    return false;
  }
  async function actionGlobalCopyState(pinia) {
    if (checkClipboardAccess())
      return;
    try {
      await navigator.clipboard.writeText(JSON.stringify(pinia.state.value));
      toastMessage("Global state copied to clipboard.");
    } catch (error) {
      if (checkNotFocusedError(error))
        return;
      toastMessage(`Failed to serialize the state. Check the console for more details.`, "error");
      console.error(error);
    }
  }
  async function actionGlobalPasteState(pinia) {
    if (checkClipboardAccess())
      return;
    try {
      loadStoresState(pinia, JSON.parse(await navigator.clipboard.readText()));
      toastMessage("Global state pasted from clipboard.");
    } catch (error) {
      if (checkNotFocusedError(error))
        return;
      toastMessage(`Failed to deserialize the state from clipboard. Check the console for more details.`, "error");
      console.error(error);
    }
  }
  async function actionGlobalSaveState(pinia) {
    try {
      saveAs(new Blob([JSON.stringify(pinia.state.value)], {
        type: "text/plain;charset=utf-8"
      }), "pinia-state.json");
    } catch (error) {
      toastMessage(`Failed to export the state as JSON. Check the console for more details.`, "error");
      console.error(error);
    }
  }
  let fileInput;
  function getFileOpener() {
    if (!fileInput) {
      fileInput = document.createElement("input");
      fileInput.type = "file";
      fileInput.accept = ".json";
    }
    function openFile() {
      return new Promise((resolve, reject) => {
        fileInput.onchange = async () => {
          const files = fileInput.files;
          if (!files)
            return resolve(null);
          const file = files.item(0);
          if (!file)
            return resolve(null);
          return resolve({ text: await file.text(), file });
        };
        fileInput.oncancel = () => resolve(null);
        fileInput.onerror = reject;
        fileInput.click();
      });
    }
    return openFile;
  }
  async function actionGlobalOpenStateFile(pinia) {
    try {
      const open2 = getFileOpener();
      const result = await open2();
      if (!result)
        return;
      const { text, file } = result;
      loadStoresState(pinia, JSON.parse(text));
      toastMessage(`Global state imported from "${file.name}".`);
    } catch (error) {
      toastMessage(`Failed to import the state from JSON. Check the console for more details.`, "error");
      console.error(error);
    }
  }
  function loadStoresState(pinia, state) {
    for (const key in state) {
      const storeState = pinia.state.value[key];
      if (storeState) {
        Object.assign(storeState, state[key]);
      } else {
        pinia.state.value[key] = state[key];
      }
    }
  }
  function formatDisplay(display) {
    return {
      _custom: {
        display
      }
    };
  }
  const PINIA_ROOT_LABEL = "🍍 Pinia (root)";
  const PINIA_ROOT_ID = "_root";
  function formatStoreForInspectorTree(store) {
    return isPinia(store) ? {
      id: PINIA_ROOT_ID,
      label: PINIA_ROOT_LABEL
    } : {
      id: store.$id,
      label: store.$id
    };
  }
  function formatStoreForInspectorState(store) {
    if (isPinia(store)) {
      const storeNames = Array.from(store._s.keys());
      const storeMap = store._s;
      const state2 = {
        state: storeNames.map((storeId) => ({
          editable: true,
          key: storeId,
          value: store.state.value[storeId]
        })),
        getters: storeNames.filter((id) => storeMap.get(id)._getters).map((id) => {
          const store2 = storeMap.get(id);
          return {
            editable: false,
            key: id,
            value: store2._getters.reduce((getters, key) => {
              getters[key] = store2[key];
              return getters;
            }, {})
          };
        })
      };
      return state2;
    }
    const state = {
      state: Object.keys(store.$state).map((key) => ({
        editable: true,
        key,
        value: store.$state[key]
      }))
    };
    if (store._getters && store._getters.length) {
      state.getters = store._getters.map((getterName) => ({
        editable: false,
        key: getterName,
        value: store[getterName]
      }));
    }
    if (store._customProperties.size) {
      state.customProperties = Array.from(store._customProperties).map((key) => ({
        editable: true,
        key,
        value: store[key]
      }));
    }
    return state;
  }
  function formatEventData(events) {
    if (!events)
      return {};
    if (Array.isArray(events)) {
      return events.reduce((data, event) => {
        data.keys.push(event.key);
        data.operations.push(event.type);
        data.oldValue[event.key] = event.oldValue;
        data.newValue[event.key] = event.newValue;
        return data;
      }, {
        oldValue: {},
        keys: [],
        operations: [],
        newValue: {}
      });
    } else {
      return {
        operation: formatDisplay(events.type),
        key: formatDisplay(events.key),
        oldValue: events.oldValue,
        newValue: events.newValue
      };
    }
  }
  function formatMutationType(type) {
    switch (type) {
      case MutationType.direct:
        return "mutation";
      case MutationType.patchFunction:
        return "$patch";
      case MutationType.patchObject:
        return "$patch";
      default:
        return "unknown";
    }
  }
  let isTimelineActive = true;
  const componentStateTypes = [];
  const MUTATIONS_LAYER_ID = "pinia:mutations";
  const INSPECTOR_ID = "pinia";
  const { assign: assign$1 } = Object;
  const getStoreType = (id) => "🍍 " + id;
  function registerPiniaDevtools(app, pinia) {
    setupDevtoolsPlugin({
      id: "dev.esm.pinia",
      label: "Pinia 🍍",
      logo: "https://pinia.vuejs.org/logo.svg",
      packageName: "pinia",
      homepage: "https://pinia.vuejs.org",
      componentStateTypes,
      app
    }, (api2) => {
      if (typeof api2.now !== "function") {
        toastMessage("You seem to be using an outdated version of Vue Devtools. Are you still using the Beta release instead of the stable one? You can find the links at https://devtools.vuejs.org/guide/installation.html.");
      }
      api2.addTimelineLayer({
        id: MUTATIONS_LAYER_ID,
        label: `Pinia 🍍`,
        color: 15064968
      });
      api2.addInspector({
        id: INSPECTOR_ID,
        label: "Pinia 🍍",
        icon: "storage",
        treeFilterPlaceholder: "Search stores",
        actions: [
          {
            icon: "content_copy",
            action: () => {
              actionGlobalCopyState(pinia);
            },
            tooltip: "Serialize and copy the state"
          },
          {
            icon: "content_paste",
            action: async () => {
              await actionGlobalPasteState(pinia);
              api2.sendInspectorTree(INSPECTOR_ID);
              api2.sendInspectorState(INSPECTOR_ID);
            },
            tooltip: "Replace the state with the content of your clipboard"
          },
          {
            icon: "save",
            action: () => {
              actionGlobalSaveState(pinia);
            },
            tooltip: "Save the state as a JSON file"
          },
          {
            icon: "folder_open",
            action: async () => {
              await actionGlobalOpenStateFile(pinia);
              api2.sendInspectorTree(INSPECTOR_ID);
              api2.sendInspectorState(INSPECTOR_ID);
            },
            tooltip: "Import the state from a JSON file"
          }
        ],
        nodeActions: [
          {
            icon: "restore",
            tooltip: 'Reset the state (with "$reset")',
            action: (nodeId) => {
              const store = pinia._s.get(nodeId);
              if (!store) {
                toastMessage(`Cannot reset "${nodeId}" store because it wasn't found.`, "warn");
              } else if (typeof store.$reset !== "function") {
                toastMessage(`Cannot reset "${nodeId}" store because it doesn't have a "$reset" method implemented.`, "warn");
              } else {
                store.$reset();
                toastMessage(`Store "${nodeId}" reset.`);
              }
            }
          }
        ]
      });
      api2.on.inspectComponent((payload, ctx) => {
        const proxy = payload.componentInstance && payload.componentInstance.proxy;
        if (proxy && proxy._pStores) {
          const piniaStores = payload.componentInstance.proxy._pStores;
          Object.values(piniaStores).forEach((store) => {
            payload.instanceData.state.push({
              type: getStoreType(store.$id),
              key: "state",
              editable: true,
              value: store._isOptionsAPI ? {
                _custom: {
                  value: vue.toRaw(store.$state),
                  actions: [
                    {
                      icon: "restore",
                      tooltip: "Reset the state of this store",
                      action: () => store.$reset()
                    }
                  ]
                }
              } : (
                // NOTE: workaround to unwrap transferred refs
                Object.keys(store.$state).reduce((state, key) => {
                  state[key] = store.$state[key];
                  return state;
                }, {})
              )
            });
            if (store._getters && store._getters.length) {
              payload.instanceData.state.push({
                type: getStoreType(store.$id),
                key: "getters",
                editable: false,
                value: store._getters.reduce((getters, key) => {
                  try {
                    getters[key] = store[key];
                  } catch (error) {
                    getters[key] = error;
                  }
                  return getters;
                }, {})
              });
            }
          });
        }
      });
      api2.on.getInspectorTree((payload) => {
        if (payload.app === app && payload.inspectorId === INSPECTOR_ID) {
          let stores = [pinia];
          stores = stores.concat(Array.from(pinia._s.values()));
          payload.rootNodes = (payload.filter ? stores.filter((store) => "$id" in store ? store.$id.toLowerCase().includes(payload.filter.toLowerCase()) : PINIA_ROOT_LABEL.toLowerCase().includes(payload.filter.toLowerCase())) : stores).map(formatStoreForInspectorTree);
        }
      });
      globalThis.$pinia = pinia;
      api2.on.getInspectorState((payload) => {
        if (payload.app === app && payload.inspectorId === INSPECTOR_ID) {
          const inspectedStore = payload.nodeId === PINIA_ROOT_ID ? pinia : pinia._s.get(payload.nodeId);
          if (!inspectedStore) {
            return;
          }
          if (inspectedStore) {
            if (payload.nodeId !== PINIA_ROOT_ID)
              globalThis.$store = vue.toRaw(inspectedStore);
            payload.state = formatStoreForInspectorState(inspectedStore);
          }
        }
      });
      api2.on.editInspectorState((payload, ctx) => {
        if (payload.app === app && payload.inspectorId === INSPECTOR_ID) {
          const inspectedStore = payload.nodeId === PINIA_ROOT_ID ? pinia : pinia._s.get(payload.nodeId);
          if (!inspectedStore) {
            return toastMessage(`store "${payload.nodeId}" not found`, "error");
          }
          const { path } = payload;
          if (!isPinia(inspectedStore)) {
            if (path.length !== 1 || !inspectedStore._customProperties.has(path[0]) || path[0] in inspectedStore.$state) {
              path.unshift("$state");
            }
          } else {
            path.unshift("state");
          }
          isTimelineActive = false;
          payload.set(inspectedStore, path, payload.state.value);
          isTimelineActive = true;
        }
      });
      api2.on.editComponentState((payload) => {
        if (payload.type.startsWith("🍍")) {
          const storeId = payload.type.replace(/^🍍\s*/, "");
          const store = pinia._s.get(storeId);
          if (!store) {
            return toastMessage(`store "${storeId}" not found`, "error");
          }
          const { path } = payload;
          if (path[0] !== "state") {
            return toastMessage(`Invalid path for store "${storeId}":
${path}
Only state can be modified.`);
          }
          path[0] = "$state";
          isTimelineActive = false;
          payload.set(store, path, payload.state.value);
          isTimelineActive = true;
        }
      });
    });
  }
  function addStoreToDevtools(app, store) {
    if (!componentStateTypes.includes(getStoreType(store.$id))) {
      componentStateTypes.push(getStoreType(store.$id));
    }
    setupDevtoolsPlugin({
      id: "dev.esm.pinia",
      label: "Pinia 🍍",
      logo: "https://pinia.vuejs.org/logo.svg",
      packageName: "pinia",
      homepage: "https://pinia.vuejs.org",
      componentStateTypes,
      app,
      settings: {
        logStoreChanges: {
          label: "Notify about new/deleted stores",
          type: "boolean",
          defaultValue: true
        }
        // useEmojis: {
        //   label: 'Use emojis in messages ⚡️',
        //   type: 'boolean',
        //   defaultValue: true,
        // },
      }
    }, (api2) => {
      const now2 = typeof api2.now === "function" ? api2.now.bind(api2) : Date.now;
      store.$onAction(({ after, onError, name, args }) => {
        const groupId = runningActionId++;
        api2.addTimelineEvent({
          layerId: MUTATIONS_LAYER_ID,
          event: {
            time: now2(),
            title: "🛫 " + name,
            subtitle: "start",
            data: {
              store: formatDisplay(store.$id),
              action: formatDisplay(name),
              args
            },
            groupId
          }
        });
        after((result) => {
          activeAction = void 0;
          api2.addTimelineEvent({
            layerId: MUTATIONS_LAYER_ID,
            event: {
              time: now2(),
              title: "🛬 " + name,
              subtitle: "end",
              data: {
                store: formatDisplay(store.$id),
                action: formatDisplay(name),
                args,
                result
              },
              groupId
            }
          });
        });
        onError((error) => {
          activeAction = void 0;
          api2.addTimelineEvent({
            layerId: MUTATIONS_LAYER_ID,
            event: {
              time: now2(),
              logType: "error",
              title: "💥 " + name,
              subtitle: "end",
              data: {
                store: formatDisplay(store.$id),
                action: formatDisplay(name),
                args,
                error
              },
              groupId
            }
          });
        });
      }, true);
      store._customProperties.forEach((name) => {
        vue.watch(() => vue.unref(store[name]), (newValue, oldValue) => {
          api2.notifyComponentUpdate();
          api2.sendInspectorState(INSPECTOR_ID);
          if (isTimelineActive) {
            api2.addTimelineEvent({
              layerId: MUTATIONS_LAYER_ID,
              event: {
                time: now2(),
                title: "Change",
                subtitle: name,
                data: {
                  newValue,
                  oldValue
                },
                groupId: activeAction
              }
            });
          }
        }, { deep: true });
      });
      store.$subscribe(({ events, type }, state) => {
        api2.notifyComponentUpdate();
        api2.sendInspectorState(INSPECTOR_ID);
        if (!isTimelineActive)
          return;
        const eventData = {
          time: now2(),
          title: formatMutationType(type),
          data: assign$1({ store: formatDisplay(store.$id) }, formatEventData(events)),
          groupId: activeAction
        };
        if (type === MutationType.patchFunction) {
          eventData.subtitle = "⤵️";
        } else if (type === MutationType.patchObject) {
          eventData.subtitle = "🧩";
        } else if (events && !Array.isArray(events)) {
          eventData.subtitle = events.type;
        }
        if (events) {
          eventData.data["rawEvent(s)"] = {
            _custom: {
              display: "DebuggerEvent",
              type: "object",
              tooltip: "raw DebuggerEvent[]",
              value: events
            }
          };
        }
        api2.addTimelineEvent({
          layerId: MUTATIONS_LAYER_ID,
          event: eventData
        });
      }, { detached: true, flush: "sync" });
      const hotUpdate = store._hotUpdate;
      store._hotUpdate = vue.markRaw((newStore) => {
        hotUpdate(newStore);
        api2.addTimelineEvent({
          layerId: MUTATIONS_LAYER_ID,
          event: {
            time: now2(),
            title: "🔥 " + store.$id,
            subtitle: "HMR update",
            data: {
              store: formatDisplay(store.$id),
              info: formatDisplay(`HMR update`)
            }
          }
        });
        api2.notifyComponentUpdate();
        api2.sendInspectorTree(INSPECTOR_ID);
        api2.sendInspectorState(INSPECTOR_ID);
      });
      const { $dispose } = store;
      store.$dispose = () => {
        $dispose();
        api2.notifyComponentUpdate();
        api2.sendInspectorTree(INSPECTOR_ID);
        api2.sendInspectorState(INSPECTOR_ID);
        api2.getSettings().logStoreChanges && toastMessage(`Disposed "${store.$id}" store 🗑`);
      };
      api2.notifyComponentUpdate();
      api2.sendInspectorTree(INSPECTOR_ID);
      api2.sendInspectorState(INSPECTOR_ID);
      api2.getSettings().logStoreChanges && toastMessage(`"${store.$id}" store installed 🆕`);
    });
  }
  let runningActionId = 0;
  let activeAction;
  function patchActionForGrouping(store, actionNames, wrapWithProxy) {
    const actions = actionNames.reduce((storeActions, actionName) => {
      storeActions[actionName] = vue.toRaw(store)[actionName];
      return storeActions;
    }, {});
    for (const actionName in actions) {
      store[actionName] = function() {
        const _actionId = runningActionId;
        const trackedStore = wrapWithProxy ? new Proxy(store, {
          get(...args) {
            activeAction = _actionId;
            return Reflect.get(...args);
          },
          set(...args) {
            activeAction = _actionId;
            return Reflect.set(...args);
          }
        }) : store;
        activeAction = _actionId;
        const retValue = actions[actionName].apply(trackedStore, arguments);
        activeAction = void 0;
        return retValue;
      };
    }
  }
  function devtoolsPlugin({ app, store, options }) {
    if (store.$id.startsWith("__hot:")) {
      return;
    }
    store._isOptionsAPI = !!options.state;
    if (!store._p._testing) {
      patchActionForGrouping(store, Object.keys(options.actions), store._isOptionsAPI);
      const originalHotUpdate = store._hotUpdate;
      vue.toRaw(store)._hotUpdate = function(newStore) {
        originalHotUpdate.apply(this, arguments);
        patchActionForGrouping(store, Object.keys(newStore._hmrPayload.actions), !!store._isOptionsAPI);
      };
    }
    addStoreToDevtools(
      app,
      // FIXME: is there a way to allow the assignment from Store<Id, S, G, A> to StoreGeneric?
      store
    );
  }
  function createPinia() {
    const scope = vue.effectScope(true);
    const state = scope.run(() => vue.ref({}));
    let _p = [];
    let toBeInstalled = [];
    const pinia = vue.markRaw({
      install(app) {
        {
          pinia._a = app;
          app.provide(piniaSymbol, pinia);
          app.config.globalProperties.$pinia = pinia;
          if (IS_CLIENT) {
            registerPiniaDevtools(app, pinia);
          }
          toBeInstalled.forEach((plugin) => _p.push(plugin));
          toBeInstalled = [];
        }
      },
      use(plugin) {
        if (!this._a && true) {
          toBeInstalled.push(plugin);
        } else {
          _p.push(plugin);
        }
        return this;
      },
      _p,
      // it's actually undefined here
      // @ts-expect-error
      _a: null,
      _e: scope,
      _s: /* @__PURE__ */ new Map(),
      state
    });
    if (IS_CLIENT && typeof Proxy !== "undefined") {
      pinia.use(devtoolsPlugin);
    }
    return pinia;
  }
  function bind(fn, thisArg) {
    return function wrap() {
      return fn.apply(thisArg, arguments);
    };
  }
  const { toString } = Object.prototype;
  const { getPrototypeOf } = Object;
  const kindOf = /* @__PURE__ */ ((cache) => (thing) => {
    const str = toString.call(thing);
    return cache[str] || (cache[str] = str.slice(8, -1).toLowerCase());
  })(/* @__PURE__ */ Object.create(null));
  const kindOfTest = (type) => {
    type = type.toLowerCase();
    return (thing) => kindOf(thing) === type;
  };
  const typeOfTest = (type) => (thing) => typeof thing === type;
  const { isArray } = Array;
  const isUndefined = typeOfTest("undefined");
  function isBuffer(val) {
    return val !== null && !isUndefined(val) && val.constructor !== null && !isUndefined(val.constructor) && isFunction(val.constructor.isBuffer) && val.constructor.isBuffer(val);
  }
  const isArrayBuffer = kindOfTest("ArrayBuffer");
  function isArrayBufferView(val) {
    let result;
    if (typeof ArrayBuffer !== "undefined" && ArrayBuffer.isView) {
      result = ArrayBuffer.isView(val);
    } else {
      result = val && val.buffer && isArrayBuffer(val.buffer);
    }
    return result;
  }
  const isString = typeOfTest("string");
  const isFunction = typeOfTest("function");
  const isNumber = typeOfTest("number");
  const isObject = (thing) => thing !== null && typeof thing === "object";
  const isBoolean = (thing) => thing === true || thing === false;
  const isPlainObject = (val) => {
    if (kindOf(val) !== "object") {
      return false;
    }
    const prototype2 = getPrototypeOf(val);
    return (prototype2 === null || prototype2 === Object.prototype || Object.getPrototypeOf(prototype2) === null) && !(Symbol.toStringTag in val) && !(Symbol.iterator in val);
  };
  const isDate = kindOfTest("Date");
  const isFile = kindOfTest("File");
  const isBlob = kindOfTest("Blob");
  const isFileList = kindOfTest("FileList");
  const isStream = (val) => isObject(val) && isFunction(val.pipe);
  const isFormData = (thing) => {
    let kind;
    return thing && (typeof FormData === "function" && thing instanceof FormData || isFunction(thing.append) && ((kind = kindOf(thing)) === "formdata" || // detect form-data instance
    kind === "object" && isFunction(thing.toString) && thing.toString() === "[object FormData]"));
  };
  const isURLSearchParams = kindOfTest("URLSearchParams");
  const [isReadableStream, isRequest, isResponse, isHeaders] = ["ReadableStream", "Request", "Response", "Headers"].map(kindOfTest);
  const trim = (str) => str.trim ? str.trim() : str.replace(/^[\s\uFEFF\xA0]+|[\s\uFEFF\xA0]+$/g, "");
  function forEach(obj, fn, { allOwnKeys = false } = {}) {
    if (obj === null || typeof obj === "undefined") {
      return;
    }
    let i;
    let l;
    if (typeof obj !== "object") {
      obj = [obj];
    }
    if (isArray(obj)) {
      for (i = 0, l = obj.length; i < l; i++) {
        fn.call(null, obj[i], i, obj);
      }
    } else {
      const keys = allOwnKeys ? Object.getOwnPropertyNames(obj) : Object.keys(obj);
      const len = keys.length;
      let key;
      for (i = 0; i < len; i++) {
        key = keys[i];
        fn.call(null, obj[key], key, obj);
      }
    }
  }
  function findKey(obj, key) {
    key = key.toLowerCase();
    const keys = Object.keys(obj);
    let i = keys.length;
    let _key;
    while (i-- > 0) {
      _key = keys[i];
      if (key === _key.toLowerCase()) {
        return _key;
      }
    }
    return null;
  }
  const _global = (() => {
    if (typeof globalThis !== "undefined") return globalThis;
    return typeof self !== "undefined" ? self : typeof window !== "undefined" ? window : global;
  })();
  const isContextDefined = (context) => !isUndefined(context) && context !== _global;
  function merge() {
    const { caseless } = isContextDefined(this) && this || {};
    const result = {};
    const assignValue = (val, key) => {
      const targetKey = caseless && findKey(result, key) || key;
      if (isPlainObject(result[targetKey]) && isPlainObject(val)) {
        result[targetKey] = merge(result[targetKey], val);
      } else if (isPlainObject(val)) {
        result[targetKey] = merge({}, val);
      } else if (isArray(val)) {
        result[targetKey] = val.slice();
      } else {
        result[targetKey] = val;
      }
    };
    for (let i = 0, l = arguments.length; i < l; i++) {
      arguments[i] && forEach(arguments[i], assignValue);
    }
    return result;
  }
  const extend = (a, b, thisArg, { allOwnKeys } = {}) => {
    forEach(b, (val, key) => {
      if (thisArg && isFunction(val)) {
        a[key] = bind(val, thisArg);
      } else {
        a[key] = val;
      }
    }, { allOwnKeys });
    return a;
  };
  const stripBOM = (content) => {
    if (content.charCodeAt(0) === 65279) {
      content = content.slice(1);
    }
    return content;
  };
  const inherits = (constructor, superConstructor, props, descriptors2) => {
    constructor.prototype = Object.create(superConstructor.prototype, descriptors2);
    constructor.prototype.constructor = constructor;
    Object.defineProperty(constructor, "super", {
      value: superConstructor.prototype
    });
    props && Object.assign(constructor.prototype, props);
  };
  const toFlatObject = (sourceObj, destObj, filter, propFilter) => {
    let props;
    let i;
    let prop;
    const merged = {};
    destObj = destObj || {};
    if (sourceObj == null) return destObj;
    do {
      props = Object.getOwnPropertyNames(sourceObj);
      i = props.length;
      while (i-- > 0) {
        prop = props[i];
        if ((!propFilter || propFilter(prop, sourceObj, destObj)) && !merged[prop]) {
          destObj[prop] = sourceObj[prop];
          merged[prop] = true;
        }
      }
      sourceObj = filter !== false && getPrototypeOf(sourceObj);
    } while (sourceObj && (!filter || filter(sourceObj, destObj)) && sourceObj !== Object.prototype);
    return destObj;
  };
  const endsWith = (str, searchString, position) => {
    str = String(str);
    if (position === void 0 || position > str.length) {
      position = str.length;
    }
    position -= searchString.length;
    const lastIndex = str.indexOf(searchString, position);
    return lastIndex !== -1 && lastIndex === position;
  };
  const toArray = (thing) => {
    if (!thing) return null;
    if (isArray(thing)) return thing;
    let i = thing.length;
    if (!isNumber(i)) return null;
    const arr = new Array(i);
    while (i-- > 0) {
      arr[i] = thing[i];
    }
    return arr;
  };
  const isTypedArray = /* @__PURE__ */ ((TypedArray) => {
    return (thing) => {
      return TypedArray && thing instanceof TypedArray;
    };
  })(typeof Uint8Array !== "undefined" && getPrototypeOf(Uint8Array));
  const forEachEntry = (obj, fn) => {
    const generator = obj && obj[Symbol.iterator];
    const iterator = generator.call(obj);
    let result;
    while ((result = iterator.next()) && !result.done) {
      const pair = result.value;
      fn.call(obj, pair[0], pair[1]);
    }
  };
  const matchAll = (regExp, str) => {
    let matches;
    const arr = [];
    while ((matches = regExp.exec(str)) !== null) {
      arr.push(matches);
    }
    return arr;
  };
  const isHTMLForm = kindOfTest("HTMLFormElement");
  const toCamelCase = (str) => {
    return str.toLowerCase().replace(
      /[-_\s]([a-z\d])(\w*)/g,
      function replacer(m, p1, p2) {
        return p1.toUpperCase() + p2;
      }
    );
  };
  const hasOwnProperty = (({ hasOwnProperty: hasOwnProperty2 }) => (obj, prop) => hasOwnProperty2.call(obj, prop))(Object.prototype);
  const isRegExp = kindOfTest("RegExp");
  const reduceDescriptors = (obj, reducer) => {
    const descriptors2 = Object.getOwnPropertyDescriptors(obj);
    const reducedDescriptors = {};
    forEach(descriptors2, (descriptor, name) => {
      let ret;
      if ((ret = reducer(descriptor, name, obj)) !== false) {
        reducedDescriptors[name] = ret || descriptor;
      }
    });
    Object.defineProperties(obj, reducedDescriptors);
  };
  const freezeMethods = (obj) => {
    reduceDescriptors(obj, (descriptor, name) => {
      if (isFunction(obj) && ["arguments", "caller", "callee"].indexOf(name) !== -1) {
        return false;
      }
      const value = obj[name];
      if (!isFunction(value)) return;
      descriptor.enumerable = false;
      if ("writable" in descriptor) {
        descriptor.writable = false;
        return;
      }
      if (!descriptor.set) {
        descriptor.set = () => {
          throw Error("Can not rewrite read-only method '" + name + "'");
        };
      }
    });
  };
  const toObjectSet = (arrayOrString, delimiter) => {
    const obj = {};
    const define2 = (arr) => {
      arr.forEach((value) => {
        obj[value] = true;
      });
    };
    isArray(arrayOrString) ? define2(arrayOrString) : define2(String(arrayOrString).split(delimiter));
    return obj;
  };
  const noop = () => {
  };
  const toFiniteNumber = (value, defaultValue) => {
    return value != null && Number.isFinite(value = +value) ? value : defaultValue;
  };
  function isSpecCompliantForm(thing) {
    return !!(thing && isFunction(thing.append) && thing[Symbol.toStringTag] === "FormData" && thing[Symbol.iterator]);
  }
  const toJSONObject = (obj) => {
    const stack = new Array(10);
    const visit = (source, i) => {
      if (isObject(source)) {
        if (stack.indexOf(source) >= 0) {
          return;
        }
        if (!("toJSON" in source)) {
          stack[i] = source;
          const target = isArray(source) ? [] : {};
          forEach(source, (value, key) => {
            const reducedValue = visit(value, i + 1);
            !isUndefined(reducedValue) && (target[key] = reducedValue);
          });
          stack[i] = void 0;
          return target;
        }
      }
      return source;
    };
    return visit(obj, 0);
  };
  const isAsyncFn = kindOfTest("AsyncFunction");
  const isThenable = (thing) => thing && (isObject(thing) || isFunction(thing)) && isFunction(thing.then) && isFunction(thing.catch);
  const _setImmediate = ((setImmediateSupported, postMessageSupported) => {
    if (setImmediateSupported) {
      return setImmediate;
    }
    return postMessageSupported ? ((token, callbacks) => {
      _global.addEventListener("message", ({ source, data }) => {
        if (source === _global && data === token) {
          callbacks.length && callbacks.shift()();
        }
      }, false);
      return (cb) => {
        callbacks.push(cb);
        _global.postMessage(token, "*");
      };
    })(`axios@${Math.random()}`, []) : (cb) => setTimeout(cb);
  })(
    typeof setImmediate === "function",
    isFunction(_global.postMessage)
  );
  const asap = typeof queueMicrotask !== "undefined" ? queueMicrotask.bind(_global) : typeof process !== "undefined" && process.nextTick || _setImmediate;
  const utils$1 = {
    isArray,
    isArrayBuffer,
    isBuffer,
    isFormData,
    isArrayBufferView,
    isString,
    isNumber,
    isBoolean,
    isObject,
    isPlainObject,
    isReadableStream,
    isRequest,
    isResponse,
    isHeaders,
    isUndefined,
    isDate,
    isFile,
    isBlob,
    isRegExp,
    isFunction,
    isStream,
    isURLSearchParams,
    isTypedArray,
    isFileList,
    forEach,
    merge,
    extend,
    trim,
    stripBOM,
    inherits,
    toFlatObject,
    kindOf,
    kindOfTest,
    endsWith,
    toArray,
    forEachEntry,
    matchAll,
    isHTMLForm,
    hasOwnProperty,
    hasOwnProp: hasOwnProperty,
    // an alias to avoid ESLint no-prototype-builtins detection
    reduceDescriptors,
    freezeMethods,
    toObjectSet,
    toCamelCase,
    noop,
    toFiniteNumber,
    findKey,
    global: _global,
    isContextDefined,
    isSpecCompliantForm,
    toJSONObject,
    isAsyncFn,
    isThenable,
    setImmediate: _setImmediate,
    asap
  };
  function AxiosError$1(message, code, config, request, response) {
    Error.call(this);
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor);
    } else {
      this.stack = new Error().stack;
    }
    this.message = message;
    this.name = "AxiosError";
    code && (this.code = code);
    config && (this.config = config);
    request && (this.request = request);
    if (response) {
      this.response = response;
      this.status = response.status ? response.status : null;
    }
  }
  utils$1.inherits(AxiosError$1, Error, {
    toJSON: function toJSON() {
      return {
        // Standard
        message: this.message,
        name: this.name,
        // Microsoft
        description: this.description,
        number: this.number,
        // Mozilla
        fileName: this.fileName,
        lineNumber: this.lineNumber,
        columnNumber: this.columnNumber,
        stack: this.stack,
        // Axios
        config: utils$1.toJSONObject(this.config),
        code: this.code,
        status: this.status
      };
    }
  });
  const prototype$1 = AxiosError$1.prototype;
  const descriptors = {};
  [
    "ERR_BAD_OPTION_VALUE",
    "ERR_BAD_OPTION",
    "ECONNABORTED",
    "ETIMEDOUT",
    "ERR_NETWORK",
    "ERR_FR_TOO_MANY_REDIRECTS",
    "ERR_DEPRECATED",
    "ERR_BAD_RESPONSE",
    "ERR_BAD_REQUEST",
    "ERR_CANCELED",
    "ERR_NOT_SUPPORT",
    "ERR_INVALID_URL"
    // eslint-disable-next-line func-names
  ].forEach((code) => {
    descriptors[code] = { value: code };
  });
  Object.defineProperties(AxiosError$1, descriptors);
  Object.defineProperty(prototype$1, "isAxiosError", { value: true });
  AxiosError$1.from = (error, code, config, request, response, customProps) => {
    const axiosError = Object.create(prototype$1);
    utils$1.toFlatObject(error, axiosError, function filter(obj) {
      return obj !== Error.prototype;
    }, (prop) => {
      return prop !== "isAxiosError";
    });
    AxiosError$1.call(axiosError, error.message, code, config, request, response);
    axiosError.cause = error;
    axiosError.name = error.name;
    customProps && Object.assign(axiosError, customProps);
    return axiosError;
  };
  const httpAdapter = null;
  function isVisitable(thing) {
    return utils$1.isPlainObject(thing) || utils$1.isArray(thing);
  }
  function removeBrackets(key) {
    return utils$1.endsWith(key, "[]") ? key.slice(0, -2) : key;
  }
  function renderKey(path, key, dots) {
    if (!path) return key;
    return path.concat(key).map(function each(token, i) {
      token = removeBrackets(token);
      return !dots && i ? "[" + token + "]" : token;
    }).join(dots ? "." : "");
  }
  function isFlatArray(arr) {
    return utils$1.isArray(arr) && !arr.some(isVisitable);
  }
  const predicates = utils$1.toFlatObject(utils$1, {}, null, function filter(prop) {
    return /^is[A-Z]/.test(prop);
  });
  function toFormData$1(obj, formData, options) {
    if (!utils$1.isObject(obj)) {
      throw new TypeError("target must be an object");
    }
    formData = formData || new FormData();
    options = utils$1.toFlatObject(options, {
      metaTokens: true,
      dots: false,
      indexes: false
    }, false, function defined(option, source) {
      return !utils$1.isUndefined(source[option]);
    });
    const metaTokens = options.metaTokens;
    const visitor = options.visitor || defaultVisitor;
    const dots = options.dots;
    const indexes = options.indexes;
    const _Blob = options.Blob || typeof Blob !== "undefined" && Blob;
    const useBlob = _Blob && utils$1.isSpecCompliantForm(formData);
    if (!utils$1.isFunction(visitor)) {
      throw new TypeError("visitor must be a function");
    }
    function convertValue(value) {
      if (value === null) return "";
      if (utils$1.isDate(value)) {
        return value.toISOString();
      }
      if (!useBlob && utils$1.isBlob(value)) {
        throw new AxiosError$1("Blob is not supported. Use a Buffer instead.");
      }
      if (utils$1.isArrayBuffer(value) || utils$1.isTypedArray(value)) {
        return useBlob && typeof Blob === "function" ? new Blob([value]) : Buffer.from(value);
      }
      return value;
    }
    function defaultVisitor(value, key, path) {
      let arr = value;
      if (value && !path && typeof value === "object") {
        if (utils$1.endsWith(key, "{}")) {
          key = metaTokens ? key : key.slice(0, -2);
          value = JSON.stringify(value);
        } else if (utils$1.isArray(value) && isFlatArray(value) || (utils$1.isFileList(value) || utils$1.endsWith(key, "[]")) && (arr = utils$1.toArray(value))) {
          key = removeBrackets(key);
          arr.forEach(function each(el, index) {
            !(utils$1.isUndefined(el) || el === null) && formData.append(
              // eslint-disable-next-line no-nested-ternary
              indexes === true ? renderKey([key], index, dots) : indexes === null ? key : key + "[]",
              convertValue(el)
            );
          });
          return false;
        }
      }
      if (isVisitable(value)) {
        return true;
      }
      formData.append(renderKey(path, key, dots), convertValue(value));
      return false;
    }
    const stack = [];
    const exposedHelpers = Object.assign(predicates, {
      defaultVisitor,
      convertValue,
      isVisitable
    });
    function build(value, path) {
      if (utils$1.isUndefined(value)) return;
      if (stack.indexOf(value) !== -1) {
        throw Error("Circular reference detected in " + path.join("."));
      }
      stack.push(value);
      utils$1.forEach(value, function each(el, key) {
        const result = !(utils$1.isUndefined(el) || el === null) && visitor.call(
          formData,
          el,
          utils$1.isString(key) ? key.trim() : key,
          path,
          exposedHelpers
        );
        if (result === true) {
          build(el, path ? path.concat(key) : [key]);
        }
      });
      stack.pop();
    }
    if (!utils$1.isObject(obj)) {
      throw new TypeError("data must be an object");
    }
    build(obj);
    return formData;
  }
  function encode$1(str) {
    const charMap = {
      "!": "%21",
      "'": "%27",
      "(": "%28",
      ")": "%29",
      "~": "%7E",
      "%20": "+",
      "%00": "\0"
    };
    return encodeURIComponent(str).replace(/[!'()~]|%20|%00/g, function replacer(match) {
      return charMap[match];
    });
  }
  function AxiosURLSearchParams(params, options) {
    this._pairs = [];
    params && toFormData$1(params, this, options);
  }
  const prototype = AxiosURLSearchParams.prototype;
  prototype.append = function append(name, value) {
    this._pairs.push([name, value]);
  };
  prototype.toString = function toString2(encoder) {
    const _encode = encoder ? function(value) {
      return encoder.call(this, value, encode$1);
    } : encode$1;
    return this._pairs.map(function each(pair) {
      return _encode(pair[0]) + "=" + _encode(pair[1]);
    }, "").join("&");
  };
  function encode(val) {
    return encodeURIComponent(val).replace(/%3A/gi, ":").replace(/%24/g, "$").replace(/%2C/gi, ",").replace(/%20/g, "+").replace(/%5B/gi, "[").replace(/%5D/gi, "]");
  }
  function buildURL(url, params, options) {
    if (!params) {
      return url;
    }
    const _encode = options && options.encode || encode;
    if (utils$1.isFunction(options)) {
      options = {
        serialize: options
      };
    }
    const serializeFn = options && options.serialize;
    let serializedParams;
    if (serializeFn) {
      serializedParams = serializeFn(params, options);
    } else {
      serializedParams = utils$1.isURLSearchParams(params) ? params.toString() : new AxiosURLSearchParams(params, options).toString(_encode);
    }
    if (serializedParams) {
      const hashmarkIndex = url.indexOf("#");
      if (hashmarkIndex !== -1) {
        url = url.slice(0, hashmarkIndex);
      }
      url += (url.indexOf("?") === -1 ? "?" : "&") + serializedParams;
    }
    return url;
  }
  class InterceptorManager {
    constructor() {
      this.handlers = [];
    }
    /**
     * Add a new interceptor to the stack
     *
     * @param {Function} fulfilled The function to handle `then` for a `Promise`
     * @param {Function} rejected The function to handle `reject` for a `Promise`
     *
     * @return {Number} An ID used to remove interceptor later
     */
    use(fulfilled, rejected, options) {
      this.handlers.push({
        fulfilled,
        rejected,
        synchronous: options ? options.synchronous : false,
        runWhen: options ? options.runWhen : null
      });
      return this.handlers.length - 1;
    }
    /**
     * Remove an interceptor from the stack
     *
     * @param {Number} id The ID that was returned by `use`
     *
     * @returns {Boolean} `true` if the interceptor was removed, `false` otherwise
     */
    eject(id) {
      if (this.handlers[id]) {
        this.handlers[id] = null;
      }
    }
    /**
     * Clear all interceptors from the stack
     *
     * @returns {void}
     */
    clear() {
      if (this.handlers) {
        this.handlers = [];
      }
    }
    /**
     * Iterate over all the registered interceptors
     *
     * This method is particularly useful for skipping over any
     * interceptors that may have become `null` calling `eject`.
     *
     * @param {Function} fn The function to call for each interceptor
     *
     * @returns {void}
     */
    forEach(fn) {
      utils$1.forEach(this.handlers, function forEachHandler(h) {
        if (h !== null) {
          fn(h);
        }
      });
    }
  }
  const transitionalDefaults = {
    silentJSONParsing: true,
    forcedJSONParsing: true,
    clarifyTimeoutError: false
  };
  const URLSearchParams$1 = typeof URLSearchParams !== "undefined" ? URLSearchParams : AxiosURLSearchParams;
  const FormData$1 = typeof FormData !== "undefined" ? FormData : null;
  const Blob$1 = typeof Blob !== "undefined" ? Blob : null;
  const platform$1 = {
    isBrowser: true,
    classes: {
      URLSearchParams: URLSearchParams$1,
      FormData: FormData$1,
      Blob: Blob$1
    },
    protocols: ["http", "https", "file", "blob", "url", "data"]
  };
  const hasBrowserEnv = typeof window !== "undefined" && typeof document !== "undefined";
  const _navigator = typeof navigator === "object" && navigator || void 0;
  const hasStandardBrowserEnv = hasBrowserEnv && (!_navigator || ["ReactNative", "NativeScript", "NS"].indexOf(_navigator.product) < 0);
  const hasStandardBrowserWebWorkerEnv = (() => {
    return typeof WorkerGlobalScope !== "undefined" && // eslint-disable-next-line no-undef
    self instanceof WorkerGlobalScope && typeof self.importScripts === "function";
  })();
  const origin = hasBrowserEnv && window.location.href || "http://localhost";
  const utils = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
    __proto__: null,
    hasBrowserEnv,
    hasStandardBrowserEnv,
    hasStandardBrowserWebWorkerEnv,
    navigator: _navigator,
    origin
  }, Symbol.toStringTag, { value: "Module" }));
  const platform = {
    ...utils,
    ...platform$1
  };
  function toURLEncodedForm(data, options) {
    return toFormData$1(data, new platform.classes.URLSearchParams(), Object.assign({
      visitor: function(value, key, path, helpers) {
        if (platform.isNode && utils$1.isBuffer(value)) {
          this.append(key, value.toString("base64"));
          return false;
        }
        return helpers.defaultVisitor.apply(this, arguments);
      }
    }, options));
  }
  function parsePropPath(name) {
    return utils$1.matchAll(/\w+|\[(\w*)]/g, name).map((match) => {
      return match[0] === "[]" ? "" : match[1] || match[0];
    });
  }
  function arrayToObject(arr) {
    const obj = {};
    const keys = Object.keys(arr);
    let i;
    const len = keys.length;
    let key;
    for (i = 0; i < len; i++) {
      key = keys[i];
      obj[key] = arr[key];
    }
    return obj;
  }
  function formDataToJSON(formData) {
    function buildPath(path, value, target, index) {
      let name = path[index++];
      if (name === "__proto__") return true;
      const isNumericKey = Number.isFinite(+name);
      const isLast = index >= path.length;
      name = !name && utils$1.isArray(target) ? target.length : name;
      if (isLast) {
        if (utils$1.hasOwnProp(target, name)) {
          target[name] = [target[name], value];
        } else {
          target[name] = value;
        }
        return !isNumericKey;
      }
      if (!target[name] || !utils$1.isObject(target[name])) {
        target[name] = [];
      }
      const result = buildPath(path, value, target[name], index);
      if (result && utils$1.isArray(target[name])) {
        target[name] = arrayToObject(target[name]);
      }
      return !isNumericKey;
    }
    if (utils$1.isFormData(formData) && utils$1.isFunction(formData.entries)) {
      const obj = {};
      utils$1.forEachEntry(formData, (name, value) => {
        buildPath(parsePropPath(name), value, obj, 0);
      });
      return obj;
    }
    return null;
  }
  function stringifySafely(rawValue, parser, encoder) {
    if (utils$1.isString(rawValue)) {
      try {
        (parser || JSON.parse)(rawValue);
        return utils$1.trim(rawValue);
      } catch (e) {
        if (e.name !== "SyntaxError") {
          throw e;
        }
      }
    }
    return (encoder || JSON.stringify)(rawValue);
  }
  const defaults = {
    transitional: transitionalDefaults,
    adapter: ["xhr", "http", "fetch"],
    transformRequest: [function transformRequest(data, headers) {
      const contentType = headers.getContentType() || "";
      const hasJSONContentType = contentType.indexOf("application/json") > -1;
      const isObjectPayload = utils$1.isObject(data);
      if (isObjectPayload && utils$1.isHTMLForm(data)) {
        data = new FormData(data);
      }
      const isFormData2 = utils$1.isFormData(data);
      if (isFormData2) {
        return hasJSONContentType ? JSON.stringify(formDataToJSON(data)) : data;
      }
      if (utils$1.isArrayBuffer(data) || utils$1.isBuffer(data) || utils$1.isStream(data) || utils$1.isFile(data) || utils$1.isBlob(data) || utils$1.isReadableStream(data)) {
        return data;
      }
      if (utils$1.isArrayBufferView(data)) {
        return data.buffer;
      }
      if (utils$1.isURLSearchParams(data)) {
        headers.setContentType("application/x-www-form-urlencoded;charset=utf-8", false);
        return data.toString();
      }
      let isFileList2;
      if (isObjectPayload) {
        if (contentType.indexOf("application/x-www-form-urlencoded") > -1) {
          return toURLEncodedForm(data, this.formSerializer).toString();
        }
        if ((isFileList2 = utils$1.isFileList(data)) || contentType.indexOf("multipart/form-data") > -1) {
          const _FormData = this.env && this.env.FormData;
          return toFormData$1(
            isFileList2 ? { "files[]": data } : data,
            _FormData && new _FormData(),
            this.formSerializer
          );
        }
      }
      if (isObjectPayload || hasJSONContentType) {
        headers.setContentType("application/json", false);
        return stringifySafely(data);
      }
      return data;
    }],
    transformResponse: [function transformResponse(data) {
      const transitional = this.transitional || defaults.transitional;
      const forcedJSONParsing = transitional && transitional.forcedJSONParsing;
      const JSONRequested = this.responseType === "json";
      if (utils$1.isResponse(data) || utils$1.isReadableStream(data)) {
        return data;
      }
      if (data && utils$1.isString(data) && (forcedJSONParsing && !this.responseType || JSONRequested)) {
        const silentJSONParsing = transitional && transitional.silentJSONParsing;
        const strictJSONParsing = !silentJSONParsing && JSONRequested;
        try {
          return JSON.parse(data);
        } catch (e) {
          if (strictJSONParsing) {
            if (e.name === "SyntaxError") {
              throw AxiosError$1.from(e, AxiosError$1.ERR_BAD_RESPONSE, this, null, this.response);
            }
            throw e;
          }
        }
      }
      return data;
    }],
    /**
     * A timeout in milliseconds to abort a request. If set to 0 (default) a
     * timeout is not created.
     */
    timeout: 0,
    xsrfCookieName: "XSRF-TOKEN",
    xsrfHeaderName: "X-XSRF-TOKEN",
    maxContentLength: -1,
    maxBodyLength: -1,
    env: {
      FormData: platform.classes.FormData,
      Blob: platform.classes.Blob
    },
    validateStatus: function validateStatus(status) {
      return status >= 200 && status < 300;
    },
    headers: {
      common: {
        "Accept": "application/json, text/plain, */*",
        "Content-Type": void 0
      }
    }
  };
  utils$1.forEach(["delete", "get", "head", "post", "put", "patch"], (method) => {
    defaults.headers[method] = {};
  });
  const ignoreDuplicateOf = utils$1.toObjectSet([
    "age",
    "authorization",
    "content-length",
    "content-type",
    "etag",
    "expires",
    "from",
    "host",
    "if-modified-since",
    "if-unmodified-since",
    "last-modified",
    "location",
    "max-forwards",
    "proxy-authorization",
    "referer",
    "retry-after",
    "user-agent"
  ]);
  const parseHeaders = (rawHeaders) => {
    const parsed = {};
    let key;
    let val;
    let i;
    rawHeaders && rawHeaders.split("\n").forEach(function parser(line) {
      i = line.indexOf(":");
      key = line.substring(0, i).trim().toLowerCase();
      val = line.substring(i + 1).trim();
      if (!key || parsed[key] && ignoreDuplicateOf[key]) {
        return;
      }
      if (key === "set-cookie") {
        if (parsed[key]) {
          parsed[key].push(val);
        } else {
          parsed[key] = [val];
        }
      } else {
        parsed[key] = parsed[key] ? parsed[key] + ", " + val : val;
      }
    });
    return parsed;
  };
  const $internals = Symbol("internals");
  function normalizeHeader(header) {
    return header && String(header).trim().toLowerCase();
  }
  function normalizeValue(value) {
    if (value === false || value == null) {
      return value;
    }
    return utils$1.isArray(value) ? value.map(normalizeValue) : String(value);
  }
  function parseTokens(str) {
    const tokens = /* @__PURE__ */ Object.create(null);
    const tokensRE = /([^\s,;=]+)\s*(?:=\s*([^,;]+))?/g;
    let match;
    while (match = tokensRE.exec(str)) {
      tokens[match[1]] = match[2];
    }
    return tokens;
  }
  const isValidHeaderName = (str) => /^[-_a-zA-Z0-9^`|~,!#$%&'*+.]+$/.test(str.trim());
  function matchHeaderValue(context, value, header, filter, isHeaderNameFilter) {
    if (utils$1.isFunction(filter)) {
      return filter.call(this, value, header);
    }
    if (isHeaderNameFilter) {
      value = header;
    }
    if (!utils$1.isString(value)) return;
    if (utils$1.isString(filter)) {
      return value.indexOf(filter) !== -1;
    }
    if (utils$1.isRegExp(filter)) {
      return filter.test(value);
    }
  }
  function formatHeader(header) {
    return header.trim().toLowerCase().replace(/([a-z\d])(\w*)/g, (w, char, str) => {
      return char.toUpperCase() + str;
    });
  }
  function buildAccessors(obj, header) {
    const accessorName = utils$1.toCamelCase(" " + header);
    ["get", "set", "has"].forEach((methodName) => {
      Object.defineProperty(obj, methodName + accessorName, {
        value: function(arg1, arg2, arg3) {
          return this[methodName].call(this, header, arg1, arg2, arg3);
        },
        configurable: true
      });
    });
  }
  let AxiosHeaders$1 = class AxiosHeaders {
    constructor(headers) {
      headers && this.set(headers);
    }
    set(header, valueOrRewrite, rewrite) {
      const self2 = this;
      function setHeader(_value, _header, _rewrite) {
        const lHeader = normalizeHeader(_header);
        if (!lHeader) {
          throw new Error("header name must be a non-empty string");
        }
        const key = utils$1.findKey(self2, lHeader);
        if (!key || self2[key] === void 0 || _rewrite === true || _rewrite === void 0 && self2[key] !== false) {
          self2[key || _header] = normalizeValue(_value);
        }
      }
      const setHeaders = (headers, _rewrite) => utils$1.forEach(headers, (_value, _header) => setHeader(_value, _header, _rewrite));
      if (utils$1.isPlainObject(header) || header instanceof this.constructor) {
        setHeaders(header, valueOrRewrite);
      } else if (utils$1.isString(header) && (header = header.trim()) && !isValidHeaderName(header)) {
        setHeaders(parseHeaders(header), valueOrRewrite);
      } else if (utils$1.isHeaders(header)) {
        for (const [key, value] of header.entries()) {
          setHeader(value, key, rewrite);
        }
      } else {
        header != null && setHeader(valueOrRewrite, header, rewrite);
      }
      return this;
    }
    get(header, parser) {
      header = normalizeHeader(header);
      if (header) {
        const key = utils$1.findKey(this, header);
        if (key) {
          const value = this[key];
          if (!parser) {
            return value;
          }
          if (parser === true) {
            return parseTokens(value);
          }
          if (utils$1.isFunction(parser)) {
            return parser.call(this, value, key);
          }
          if (utils$1.isRegExp(parser)) {
            return parser.exec(value);
          }
          throw new TypeError("parser must be boolean|regexp|function");
        }
      }
    }
    has(header, matcher) {
      header = normalizeHeader(header);
      if (header) {
        const key = utils$1.findKey(this, header);
        return !!(key && this[key] !== void 0 && (!matcher || matchHeaderValue(this, this[key], key, matcher)));
      }
      return false;
    }
    delete(header, matcher) {
      const self2 = this;
      let deleted = false;
      function deleteHeader(_header) {
        _header = normalizeHeader(_header);
        if (_header) {
          const key = utils$1.findKey(self2, _header);
          if (key && (!matcher || matchHeaderValue(self2, self2[key], key, matcher))) {
            delete self2[key];
            deleted = true;
          }
        }
      }
      if (utils$1.isArray(header)) {
        header.forEach(deleteHeader);
      } else {
        deleteHeader(header);
      }
      return deleted;
    }
    clear(matcher) {
      const keys = Object.keys(this);
      let i = keys.length;
      let deleted = false;
      while (i--) {
        const key = keys[i];
        if (!matcher || matchHeaderValue(this, this[key], key, matcher, true)) {
          delete this[key];
          deleted = true;
        }
      }
      return deleted;
    }
    normalize(format) {
      const self2 = this;
      const headers = {};
      utils$1.forEach(this, (value, header) => {
        const key = utils$1.findKey(headers, header);
        if (key) {
          self2[key] = normalizeValue(value);
          delete self2[header];
          return;
        }
        const normalized = format ? formatHeader(header) : String(header).trim();
        if (normalized !== header) {
          delete self2[header];
        }
        self2[normalized] = normalizeValue(value);
        headers[normalized] = true;
      });
      return this;
    }
    concat(...targets) {
      return this.constructor.concat(this, ...targets);
    }
    toJSON(asStrings) {
      const obj = /* @__PURE__ */ Object.create(null);
      utils$1.forEach(this, (value, header) => {
        value != null && value !== false && (obj[header] = asStrings && utils$1.isArray(value) ? value.join(", ") : value);
      });
      return obj;
    }
    [Symbol.iterator]() {
      return Object.entries(this.toJSON())[Symbol.iterator]();
    }
    toString() {
      return Object.entries(this.toJSON()).map(([header, value]) => header + ": " + value).join("\n");
    }
    get [Symbol.toStringTag]() {
      return "AxiosHeaders";
    }
    static from(thing) {
      return thing instanceof this ? thing : new this(thing);
    }
    static concat(first, ...targets) {
      const computed = new this(first);
      targets.forEach((target) => computed.set(target));
      return computed;
    }
    static accessor(header) {
      const internals = this[$internals] = this[$internals] = {
        accessors: {}
      };
      const accessors = internals.accessors;
      const prototype2 = this.prototype;
      function defineAccessor(_header) {
        const lHeader = normalizeHeader(_header);
        if (!accessors[lHeader]) {
          buildAccessors(prototype2, _header);
          accessors[lHeader] = true;
        }
      }
      utils$1.isArray(header) ? header.forEach(defineAccessor) : defineAccessor(header);
      return this;
    }
  };
  AxiosHeaders$1.accessor(["Content-Type", "Content-Length", "Accept", "Accept-Encoding", "User-Agent", "Authorization"]);
  utils$1.reduceDescriptors(AxiosHeaders$1.prototype, ({ value }, key) => {
    let mapped = key[0].toUpperCase() + key.slice(1);
    return {
      get: () => value,
      set(headerValue) {
        this[mapped] = headerValue;
      }
    };
  });
  utils$1.freezeMethods(AxiosHeaders$1);
  function transformData(fns, response) {
    const config = this || defaults;
    const context = response || config;
    const headers = AxiosHeaders$1.from(context.headers);
    let data = context.data;
    utils$1.forEach(fns, function transform(fn) {
      data = fn.call(config, data, headers.normalize(), response ? response.status : void 0);
    });
    headers.normalize();
    return data;
  }
  function isCancel$1(value) {
    return !!(value && value.__CANCEL__);
  }
  function CanceledError$1(message, config, request) {
    AxiosError$1.call(this, message == null ? "canceled" : message, AxiosError$1.ERR_CANCELED, config, request);
    this.name = "CanceledError";
  }
  utils$1.inherits(CanceledError$1, AxiosError$1, {
    __CANCEL__: true
  });
  function settle(resolve, reject, response) {
    const validateStatus = response.config.validateStatus;
    if (!response.status || !validateStatus || validateStatus(response.status)) {
      resolve(response);
    } else {
      reject(new AxiosError$1(
        "Request failed with status code " + response.status,
        [AxiosError$1.ERR_BAD_REQUEST, AxiosError$1.ERR_BAD_RESPONSE][Math.floor(response.status / 100) - 4],
        response.config,
        response.request,
        response
      ));
    }
  }
  function parseProtocol(url) {
    const match = /^([-+\w]{1,25})(:?\/\/|:)/.exec(url);
    return match && match[1] || "";
  }
  function speedometer(samplesCount, min) {
    samplesCount = samplesCount || 10;
    const bytes = new Array(samplesCount);
    const timestamps = new Array(samplesCount);
    let head = 0;
    let tail = 0;
    let firstSampleTS;
    min = min !== void 0 ? min : 1e3;
    return function push(chunkLength) {
      const now2 = Date.now();
      const startedAt = timestamps[tail];
      if (!firstSampleTS) {
        firstSampleTS = now2;
      }
      bytes[head] = chunkLength;
      timestamps[head] = now2;
      let i = tail;
      let bytesCount = 0;
      while (i !== head) {
        bytesCount += bytes[i++];
        i = i % samplesCount;
      }
      head = (head + 1) % samplesCount;
      if (head === tail) {
        tail = (tail + 1) % samplesCount;
      }
      if (now2 - firstSampleTS < min) {
        return;
      }
      const passed = startedAt && now2 - startedAt;
      return passed ? Math.round(bytesCount * 1e3 / passed) : void 0;
    };
  }
  function throttle(fn, freq) {
    let timestamp = 0;
    let threshold = 1e3 / freq;
    let lastArgs;
    let timer;
    const invoke = (args, now2 = Date.now()) => {
      timestamp = now2;
      lastArgs = null;
      if (timer) {
        clearTimeout(timer);
        timer = null;
      }
      fn.apply(null, args);
    };
    const throttled = (...args) => {
      const now2 = Date.now();
      const passed = now2 - timestamp;
      if (passed >= threshold) {
        invoke(args, now2);
      } else {
        lastArgs = args;
        if (!timer) {
          timer = setTimeout(() => {
            timer = null;
            invoke(lastArgs);
          }, threshold - passed);
        }
      }
    };
    const flush = () => lastArgs && invoke(lastArgs);
    return [throttled, flush];
  }
  const progressEventReducer = (listener, isDownloadStream, freq = 3) => {
    let bytesNotified = 0;
    const _speedometer = speedometer(50, 250);
    return throttle((e) => {
      const loaded = e.loaded;
      const total = e.lengthComputable ? e.total : void 0;
      const progressBytes = loaded - bytesNotified;
      const rate = _speedometer(progressBytes);
      const inRange = loaded <= total;
      bytesNotified = loaded;
      const data = {
        loaded,
        total,
        progress: total ? loaded / total : void 0,
        bytes: progressBytes,
        rate: rate ? rate : void 0,
        estimated: rate && total && inRange ? (total - loaded) / rate : void 0,
        event: e,
        lengthComputable: total != null,
        [isDownloadStream ? "download" : "upload"]: true
      };
      listener(data);
    }, freq);
  };
  const progressEventDecorator = (total, throttled) => {
    const lengthComputable = total != null;
    return [(loaded) => throttled[0]({
      lengthComputable,
      total,
      loaded
    }), throttled[1]];
  };
  const asyncDecorator = (fn) => (...args) => utils$1.asap(() => fn(...args));
  const isURLSameOrigin = platform.hasStandardBrowserEnv ? /* @__PURE__ */ ((origin2, isMSIE) => (url) => {
    url = new URL(url, platform.origin);
    return origin2.protocol === url.protocol && origin2.host === url.host && (isMSIE || origin2.port === url.port);
  })(
    new URL(platform.origin),
    platform.navigator && /(msie|trident)/i.test(platform.navigator.userAgent)
  ) : () => true;
  const cookies = platform.hasStandardBrowserEnv ? (
    // Standard browser envs support document.cookie
    {
      write(name, value, expires, path, domain, secure) {
        const cookie = [name + "=" + encodeURIComponent(value)];
        utils$1.isNumber(expires) && cookie.push("expires=" + new Date(expires).toGMTString());
        utils$1.isString(path) && cookie.push("path=" + path);
        utils$1.isString(domain) && cookie.push("domain=" + domain);
        secure === true && cookie.push("secure");
        document.cookie = cookie.join("; ");
      },
      read(name) {
        const match = document.cookie.match(new RegExp("(^|;\\s*)(" + name + ")=([^;]*)"));
        return match ? decodeURIComponent(match[3]) : null;
      },
      remove(name) {
        this.write(name, "", Date.now() - 864e5);
      }
    }
  ) : (
    // Non-standard browser env (web workers, react-native) lack needed support.
    {
      write() {
      },
      read() {
        return null;
      },
      remove() {
      }
    }
  );
  function isAbsoluteURL(url) {
    return /^([a-z][a-z\d+\-.]*:)?\/\//i.test(url);
  }
  function combineURLs(baseURL, relativeURL) {
    return relativeURL ? baseURL.replace(/\/?\/$/, "") + "/" + relativeURL.replace(/^\/+/, "") : baseURL;
  }
  function buildFullPath(baseURL, requestedURL, allowAbsoluteUrls) {
    let isRelativeUrl = !isAbsoluteURL(requestedURL);
    if (baseURL && (isRelativeUrl || allowAbsoluteUrls == false)) {
      return combineURLs(baseURL, requestedURL);
    }
    return requestedURL;
  }
  const headersToObject = (thing) => thing instanceof AxiosHeaders$1 ? { ...thing } : thing;
  function mergeConfig$1(config1, config2) {
    config2 = config2 || {};
    const config = {};
    function getMergedValue(target, source, prop, caseless) {
      if (utils$1.isPlainObject(target) && utils$1.isPlainObject(source)) {
        return utils$1.merge.call({ caseless }, target, source);
      } else if (utils$1.isPlainObject(source)) {
        return utils$1.merge({}, source);
      } else if (utils$1.isArray(source)) {
        return source.slice();
      }
      return source;
    }
    function mergeDeepProperties(a, b, prop, caseless) {
      if (!utils$1.isUndefined(b)) {
        return getMergedValue(a, b, prop, caseless);
      } else if (!utils$1.isUndefined(a)) {
        return getMergedValue(void 0, a, prop, caseless);
      }
    }
    function valueFromConfig2(a, b) {
      if (!utils$1.isUndefined(b)) {
        return getMergedValue(void 0, b);
      }
    }
    function defaultToConfig2(a, b) {
      if (!utils$1.isUndefined(b)) {
        return getMergedValue(void 0, b);
      } else if (!utils$1.isUndefined(a)) {
        return getMergedValue(void 0, a);
      }
    }
    function mergeDirectKeys(a, b, prop) {
      if (prop in config2) {
        return getMergedValue(a, b);
      } else if (prop in config1) {
        return getMergedValue(void 0, a);
      }
    }
    const mergeMap = {
      url: valueFromConfig2,
      method: valueFromConfig2,
      data: valueFromConfig2,
      baseURL: defaultToConfig2,
      transformRequest: defaultToConfig2,
      transformResponse: defaultToConfig2,
      paramsSerializer: defaultToConfig2,
      timeout: defaultToConfig2,
      timeoutMessage: defaultToConfig2,
      withCredentials: defaultToConfig2,
      withXSRFToken: defaultToConfig2,
      adapter: defaultToConfig2,
      responseType: defaultToConfig2,
      xsrfCookieName: defaultToConfig2,
      xsrfHeaderName: defaultToConfig2,
      onUploadProgress: defaultToConfig2,
      onDownloadProgress: defaultToConfig2,
      decompress: defaultToConfig2,
      maxContentLength: defaultToConfig2,
      maxBodyLength: defaultToConfig2,
      beforeRedirect: defaultToConfig2,
      transport: defaultToConfig2,
      httpAgent: defaultToConfig2,
      httpsAgent: defaultToConfig2,
      cancelToken: defaultToConfig2,
      socketPath: defaultToConfig2,
      responseEncoding: defaultToConfig2,
      validateStatus: mergeDirectKeys,
      headers: (a, b, prop) => mergeDeepProperties(headersToObject(a), headersToObject(b), prop, true)
    };
    utils$1.forEach(Object.keys(Object.assign({}, config1, config2)), function computeConfigValue(prop) {
      const merge2 = mergeMap[prop] || mergeDeepProperties;
      const configValue = merge2(config1[prop], config2[prop], prop);
      utils$1.isUndefined(configValue) && merge2 !== mergeDirectKeys || (config[prop] = configValue);
    });
    return config;
  }
  const resolveConfig = (config) => {
    const newConfig = mergeConfig$1({}, config);
    let { data, withXSRFToken, xsrfHeaderName, xsrfCookieName, headers, auth } = newConfig;
    newConfig.headers = headers = AxiosHeaders$1.from(headers);
    newConfig.url = buildURL(buildFullPath(newConfig.baseURL, newConfig.url, newConfig.allowAbsoluteUrls), config.params, config.paramsSerializer);
    if (auth) {
      headers.set(
        "Authorization",
        "Basic " + btoa((auth.username || "") + ":" + (auth.password ? unescape(encodeURIComponent(auth.password)) : ""))
      );
    }
    let contentType;
    if (utils$1.isFormData(data)) {
      if (platform.hasStandardBrowserEnv || platform.hasStandardBrowserWebWorkerEnv) {
        headers.setContentType(void 0);
      } else if ((contentType = headers.getContentType()) !== false) {
        const [type, ...tokens] = contentType ? contentType.split(";").map((token) => token.trim()).filter(Boolean) : [];
        headers.setContentType([type || "multipart/form-data", ...tokens].join("; "));
      }
    }
    if (platform.hasStandardBrowserEnv) {
      withXSRFToken && utils$1.isFunction(withXSRFToken) && (withXSRFToken = withXSRFToken(newConfig));
      if (withXSRFToken || withXSRFToken !== false && isURLSameOrigin(newConfig.url)) {
        const xsrfValue = xsrfHeaderName && xsrfCookieName && cookies.read(xsrfCookieName);
        if (xsrfValue) {
          headers.set(xsrfHeaderName, xsrfValue);
        }
      }
    }
    return newConfig;
  };
  const isXHRAdapterSupported = typeof XMLHttpRequest !== "undefined";
  const xhrAdapter = isXHRAdapterSupported && function(config) {
    return new Promise(function dispatchXhrRequest(resolve, reject) {
      const _config = resolveConfig(config);
      let requestData = _config.data;
      const requestHeaders = AxiosHeaders$1.from(_config.headers).normalize();
      let { responseType, onUploadProgress, onDownloadProgress } = _config;
      let onCanceled;
      let uploadThrottled, downloadThrottled;
      let flushUpload, flushDownload;
      function done() {
        flushUpload && flushUpload();
        flushDownload && flushDownload();
        _config.cancelToken && _config.cancelToken.unsubscribe(onCanceled);
        _config.signal && _config.signal.removeEventListener("abort", onCanceled);
      }
      let request = new XMLHttpRequest();
      request.open(_config.method.toUpperCase(), _config.url, true);
      request.timeout = _config.timeout;
      function onloadend() {
        if (!request) {
          return;
        }
        const responseHeaders = AxiosHeaders$1.from(
          "getAllResponseHeaders" in request && request.getAllResponseHeaders()
        );
        const responseData = !responseType || responseType === "text" || responseType === "json" ? request.responseText : request.response;
        const response = {
          data: responseData,
          status: request.status,
          statusText: request.statusText,
          headers: responseHeaders,
          config,
          request
        };
        settle(function _resolve(value) {
          resolve(value);
          done();
        }, function _reject(err) {
          reject(err);
          done();
        }, response);
        request = null;
      }
      if ("onloadend" in request) {
        request.onloadend = onloadend;
      } else {
        request.onreadystatechange = function handleLoad() {
          if (!request || request.readyState !== 4) {
            return;
          }
          if (request.status === 0 && !(request.responseURL && request.responseURL.indexOf("file:") === 0)) {
            return;
          }
          setTimeout(onloadend);
        };
      }
      request.onabort = function handleAbort() {
        if (!request) {
          return;
        }
        reject(new AxiosError$1("Request aborted", AxiosError$1.ECONNABORTED, config, request));
        request = null;
      };
      request.onerror = function handleError() {
        reject(new AxiosError$1("Network Error", AxiosError$1.ERR_NETWORK, config, request));
        request = null;
      };
      request.ontimeout = function handleTimeout() {
        let timeoutErrorMessage = _config.timeout ? "timeout of " + _config.timeout + "ms exceeded" : "timeout exceeded";
        const transitional = _config.transitional || transitionalDefaults;
        if (_config.timeoutErrorMessage) {
          timeoutErrorMessage = _config.timeoutErrorMessage;
        }
        reject(new AxiosError$1(
          timeoutErrorMessage,
          transitional.clarifyTimeoutError ? AxiosError$1.ETIMEDOUT : AxiosError$1.ECONNABORTED,
          config,
          request
        ));
        request = null;
      };
      requestData === void 0 && requestHeaders.setContentType(null);
      if ("setRequestHeader" in request) {
        utils$1.forEach(requestHeaders.toJSON(), function setRequestHeader(val, key) {
          request.setRequestHeader(key, val);
        });
      }
      if (!utils$1.isUndefined(_config.withCredentials)) {
        request.withCredentials = !!_config.withCredentials;
      }
      if (responseType && responseType !== "json") {
        request.responseType = _config.responseType;
      }
      if (onDownloadProgress) {
        [downloadThrottled, flushDownload] = progressEventReducer(onDownloadProgress, true);
        request.addEventListener("progress", downloadThrottled);
      }
      if (onUploadProgress && request.upload) {
        [uploadThrottled, flushUpload] = progressEventReducer(onUploadProgress);
        request.upload.addEventListener("progress", uploadThrottled);
        request.upload.addEventListener("loadend", flushUpload);
      }
      if (_config.cancelToken || _config.signal) {
        onCanceled = (cancel) => {
          if (!request) {
            return;
          }
          reject(!cancel || cancel.type ? new CanceledError$1(null, config, request) : cancel);
          request.abort();
          request = null;
        };
        _config.cancelToken && _config.cancelToken.subscribe(onCanceled);
        if (_config.signal) {
          _config.signal.aborted ? onCanceled() : _config.signal.addEventListener("abort", onCanceled);
        }
      }
      const protocol = parseProtocol(_config.url);
      if (protocol && platform.protocols.indexOf(protocol) === -1) {
        reject(new AxiosError$1("Unsupported protocol " + protocol + ":", AxiosError$1.ERR_BAD_REQUEST, config));
        return;
      }
      request.send(requestData || null);
    });
  };
  const composeSignals = (signals, timeout) => {
    const { length } = signals = signals ? signals.filter(Boolean) : [];
    if (timeout || length) {
      let controller = new AbortController();
      let aborted;
      const onabort = function(reason) {
        if (!aborted) {
          aborted = true;
          unsubscribe();
          const err = reason instanceof Error ? reason : this.reason;
          controller.abort(err instanceof AxiosError$1 ? err : new CanceledError$1(err instanceof Error ? err.message : err));
        }
      };
      let timer = timeout && setTimeout(() => {
        timer = null;
        onabort(new AxiosError$1(`timeout ${timeout} of ms exceeded`, AxiosError$1.ETIMEDOUT));
      }, timeout);
      const unsubscribe = () => {
        if (signals) {
          timer && clearTimeout(timer);
          timer = null;
          signals.forEach((signal2) => {
            signal2.unsubscribe ? signal2.unsubscribe(onabort) : signal2.removeEventListener("abort", onabort);
          });
          signals = null;
        }
      };
      signals.forEach((signal2) => signal2.addEventListener("abort", onabort));
      const { signal } = controller;
      signal.unsubscribe = () => utils$1.asap(unsubscribe);
      return signal;
    }
  };
  const streamChunk = function* (chunk, chunkSize) {
    let len = chunk.byteLength;
    if (len < chunkSize) {
      yield chunk;
      return;
    }
    let pos = 0;
    let end;
    while (pos < len) {
      end = pos + chunkSize;
      yield chunk.slice(pos, end);
      pos = end;
    }
  };
  const readBytes = async function* (iterable, chunkSize) {
    for await (const chunk of readStream(iterable)) {
      yield* streamChunk(chunk, chunkSize);
    }
  };
  const readStream = async function* (stream) {
    if (stream[Symbol.asyncIterator]) {
      yield* stream;
      return;
    }
    const reader = stream.getReader();
    try {
      for (; ; ) {
        const { done, value } = await reader.read();
        if (done) {
          break;
        }
        yield value;
      }
    } finally {
      await reader.cancel();
    }
  };
  const trackStream = (stream, chunkSize, onProgress, onFinish) => {
    const iterator = readBytes(stream, chunkSize);
    let bytes = 0;
    let done;
    let _onFinish = (e) => {
      if (!done) {
        done = true;
        onFinish && onFinish(e);
      }
    };
    return new ReadableStream({
      async pull(controller) {
        try {
          const { done: done2, value } = await iterator.next();
          if (done2) {
            _onFinish();
            controller.close();
            return;
          }
          let len = value.byteLength;
          if (onProgress) {
            let loadedBytes = bytes += len;
            onProgress(loadedBytes);
          }
          controller.enqueue(new Uint8Array(value));
        } catch (err) {
          _onFinish(err);
          throw err;
        }
      },
      cancel(reason) {
        _onFinish(reason);
        return iterator.return();
      }
    }, {
      highWaterMark: 2
    });
  };
  const isFetchSupported = typeof fetch === "function" && typeof Request === "function" && typeof Response === "function";
  const isReadableStreamSupported = isFetchSupported && typeof ReadableStream === "function";
  const encodeText = isFetchSupported && (typeof TextEncoder === "function" ? /* @__PURE__ */ ((encoder) => (str) => encoder.encode(str))(new TextEncoder()) : async (str) => new Uint8Array(await new Response(str).arrayBuffer()));
  const test = (fn, ...args) => {
    try {
      return !!fn(...args);
    } catch (e) {
      return false;
    }
  };
  const supportsRequestStream = isReadableStreamSupported && test(() => {
    let duplexAccessed = false;
    const hasContentType = new Request(platform.origin, {
      body: new ReadableStream(),
      method: "POST",
      get duplex() {
        duplexAccessed = true;
        return "half";
      }
    }).headers.has("Content-Type");
    return duplexAccessed && !hasContentType;
  });
  const DEFAULT_CHUNK_SIZE = 64 * 1024;
  const supportsResponseStream = isReadableStreamSupported && test(() => utils$1.isReadableStream(new Response("").body));
  const resolvers = {
    stream: supportsResponseStream && ((res) => res.body)
  };
  isFetchSupported && ((res) => {
    ["text", "arrayBuffer", "blob", "formData", "stream"].forEach((type) => {
      !resolvers[type] && (resolvers[type] = utils$1.isFunction(res[type]) ? (res2) => res2[type]() : (_, config) => {
        throw new AxiosError$1(`Response type '${type}' is not supported`, AxiosError$1.ERR_NOT_SUPPORT, config);
      });
    });
  })(new Response());
  const getBodyLength = async (body) => {
    if (body == null) {
      return 0;
    }
    if (utils$1.isBlob(body)) {
      return body.size;
    }
    if (utils$1.isSpecCompliantForm(body)) {
      const _request = new Request(platform.origin, {
        method: "POST",
        body
      });
      return (await _request.arrayBuffer()).byteLength;
    }
    if (utils$1.isArrayBufferView(body) || utils$1.isArrayBuffer(body)) {
      return body.byteLength;
    }
    if (utils$1.isURLSearchParams(body)) {
      body = body + "";
    }
    if (utils$1.isString(body)) {
      return (await encodeText(body)).byteLength;
    }
  };
  const resolveBodyLength = async (headers, body) => {
    const length = utils$1.toFiniteNumber(headers.getContentLength());
    return length == null ? getBodyLength(body) : length;
  };
  const fetchAdapter = isFetchSupported && (async (config) => {
    let {
      url,
      method,
      data,
      signal,
      cancelToken,
      timeout,
      onDownloadProgress,
      onUploadProgress,
      responseType,
      headers,
      withCredentials = "same-origin",
      fetchOptions
    } = resolveConfig(config);
    responseType = responseType ? (responseType + "").toLowerCase() : "text";
    let composedSignal = composeSignals([signal, cancelToken && cancelToken.toAbortSignal()], timeout);
    let request;
    const unsubscribe = composedSignal && composedSignal.unsubscribe && (() => {
      composedSignal.unsubscribe();
    });
    let requestContentLength;
    try {
      if (onUploadProgress && supportsRequestStream && method !== "get" && method !== "head" && (requestContentLength = await resolveBodyLength(headers, data)) !== 0) {
        let _request = new Request(url, {
          method: "POST",
          body: data,
          duplex: "half"
        });
        let contentTypeHeader;
        if (utils$1.isFormData(data) && (contentTypeHeader = _request.headers.get("content-type"))) {
          headers.setContentType(contentTypeHeader);
        }
        if (_request.body) {
          const [onProgress, flush] = progressEventDecorator(
            requestContentLength,
            progressEventReducer(asyncDecorator(onUploadProgress))
          );
          data = trackStream(_request.body, DEFAULT_CHUNK_SIZE, onProgress, flush);
        }
      }
      if (!utils$1.isString(withCredentials)) {
        withCredentials = withCredentials ? "include" : "omit";
      }
      const isCredentialsSupported = "credentials" in Request.prototype;
      request = new Request(url, {
        ...fetchOptions,
        signal: composedSignal,
        method: method.toUpperCase(),
        headers: headers.normalize().toJSON(),
        body: data,
        duplex: "half",
        credentials: isCredentialsSupported ? withCredentials : void 0
      });
      let response = await fetch(request);
      const isStreamResponse = supportsResponseStream && (responseType === "stream" || responseType === "response");
      if (supportsResponseStream && (onDownloadProgress || isStreamResponse && unsubscribe)) {
        const options = {};
        ["status", "statusText", "headers"].forEach((prop) => {
          options[prop] = response[prop];
        });
        const responseContentLength = utils$1.toFiniteNumber(response.headers.get("content-length"));
        const [onProgress, flush] = onDownloadProgress && progressEventDecorator(
          responseContentLength,
          progressEventReducer(asyncDecorator(onDownloadProgress), true)
        ) || [];
        response = new Response(
          trackStream(response.body, DEFAULT_CHUNK_SIZE, onProgress, () => {
            flush && flush();
            unsubscribe && unsubscribe();
          }),
          options
        );
      }
      responseType = responseType || "text";
      let responseData = await resolvers[utils$1.findKey(resolvers, responseType) || "text"](response, config);
      !isStreamResponse && unsubscribe && unsubscribe();
      return await new Promise((resolve, reject) => {
        settle(resolve, reject, {
          data: responseData,
          headers: AxiosHeaders$1.from(response.headers),
          status: response.status,
          statusText: response.statusText,
          config,
          request
        });
      });
    } catch (err) {
      unsubscribe && unsubscribe();
      if (err && err.name === "TypeError" && /fetch/i.test(err.message)) {
        throw Object.assign(
          new AxiosError$1("Network Error", AxiosError$1.ERR_NETWORK, config, request),
          {
            cause: err.cause || err
          }
        );
      }
      throw AxiosError$1.from(err, err && err.code, config, request);
    }
  });
  const knownAdapters = {
    http: httpAdapter,
    xhr: xhrAdapter,
    fetch: fetchAdapter
  };
  utils$1.forEach(knownAdapters, (fn, value) => {
    if (fn) {
      try {
        Object.defineProperty(fn, "name", { value });
      } catch (e) {
      }
      Object.defineProperty(fn, "adapterName", { value });
    }
  });
  const renderReason = (reason) => `- ${reason}`;
  const isResolvedHandle = (adapter) => utils$1.isFunction(adapter) || adapter === null || adapter === false;
  const adapters = {
    getAdapter: (adapters2) => {
      adapters2 = utils$1.isArray(adapters2) ? adapters2 : [adapters2];
      const { length } = adapters2;
      let nameOrAdapter;
      let adapter;
      const rejectedReasons = {};
      for (let i = 0; i < length; i++) {
        nameOrAdapter = adapters2[i];
        let id;
        adapter = nameOrAdapter;
        if (!isResolvedHandle(nameOrAdapter)) {
          adapter = knownAdapters[(id = String(nameOrAdapter)).toLowerCase()];
          if (adapter === void 0) {
            throw new AxiosError$1(`Unknown adapter '${id}'`);
          }
        }
        if (adapter) {
          break;
        }
        rejectedReasons[id || "#" + i] = adapter;
      }
      if (!adapter) {
        const reasons = Object.entries(rejectedReasons).map(
          ([id, state]) => `adapter ${id} ` + (state === false ? "is not supported by the environment" : "is not available in the build")
        );
        let s = length ? reasons.length > 1 ? "since :\n" + reasons.map(renderReason).join("\n") : " " + renderReason(reasons[0]) : "as no adapter specified";
        throw new AxiosError$1(
          `There is no suitable adapter to dispatch the request ` + s,
          "ERR_NOT_SUPPORT"
        );
      }
      return adapter;
    },
    adapters: knownAdapters
  };
  function throwIfCancellationRequested(config) {
    if (config.cancelToken) {
      config.cancelToken.throwIfRequested();
    }
    if (config.signal && config.signal.aborted) {
      throw new CanceledError$1(null, config);
    }
  }
  function dispatchRequest(config) {
    throwIfCancellationRequested(config);
    config.headers = AxiosHeaders$1.from(config.headers);
    config.data = transformData.call(
      config,
      config.transformRequest
    );
    if (["post", "put", "patch"].indexOf(config.method) !== -1) {
      config.headers.setContentType("application/x-www-form-urlencoded", false);
    }
    const adapter = adapters.getAdapter(config.adapter || defaults.adapter);
    return adapter(config).then(function onAdapterResolution(response) {
      throwIfCancellationRequested(config);
      response.data = transformData.call(
        config,
        config.transformResponse,
        response
      );
      response.headers = AxiosHeaders$1.from(response.headers);
      return response;
    }, function onAdapterRejection(reason) {
      if (!isCancel$1(reason)) {
        throwIfCancellationRequested(config);
        if (reason && reason.response) {
          reason.response.data = transformData.call(
            config,
            config.transformResponse,
            reason.response
          );
          reason.response.headers = AxiosHeaders$1.from(reason.response.headers);
        }
      }
      return Promise.reject(reason);
    });
  }
  const VERSION$1 = "1.8.4";
  const validators$1 = {};
  ["object", "boolean", "number", "function", "string", "symbol"].forEach((type, i) => {
    validators$1[type] = function validator2(thing) {
      return typeof thing === type || "a" + (i < 1 ? "n " : " ") + type;
    };
  });
  const deprecatedWarnings = {};
  validators$1.transitional = function transitional(validator2, version, message) {
    function formatMessage(opt, desc) {
      return "[Axios v" + VERSION$1 + "] Transitional option '" + opt + "'" + desc + (message ? ". " + message : "");
    }
    return (value, opt, opts) => {
      if (validator2 === false) {
        throw new AxiosError$1(
          formatMessage(opt, " has been removed" + (version ? " in " + version : "")),
          AxiosError$1.ERR_DEPRECATED
        );
      }
      if (version && !deprecatedWarnings[opt]) {
        deprecatedWarnings[opt] = true;
        console.warn(
          formatMessage(
            opt,
            " has been deprecated since v" + version + " and will be removed in the near future"
          )
        );
      }
      return validator2 ? validator2(value, opt, opts) : true;
    };
  };
  validators$1.spelling = function spelling(correctSpelling) {
    return (value, opt) => {
      console.warn(`${opt} is likely a misspelling of ${correctSpelling}`);
      return true;
    };
  };
  function assertOptions(options, schema, allowUnknown) {
    if (typeof options !== "object") {
      throw new AxiosError$1("options must be an object", AxiosError$1.ERR_BAD_OPTION_VALUE);
    }
    const keys = Object.keys(options);
    let i = keys.length;
    while (i-- > 0) {
      const opt = keys[i];
      const validator2 = schema[opt];
      if (validator2) {
        const value = options[opt];
        const result = value === void 0 || validator2(value, opt, options);
        if (result !== true) {
          throw new AxiosError$1("option " + opt + " must be " + result, AxiosError$1.ERR_BAD_OPTION_VALUE);
        }
        continue;
      }
      if (allowUnknown !== true) {
        throw new AxiosError$1("Unknown option " + opt, AxiosError$1.ERR_BAD_OPTION);
      }
    }
  }
  const validator = {
    assertOptions,
    validators: validators$1
  };
  const validators = validator.validators;
  let Axios$1 = class Axios {
    constructor(instanceConfig) {
      this.defaults = instanceConfig;
      this.interceptors = {
        request: new InterceptorManager(),
        response: new InterceptorManager()
      };
    }
    /**
     * Dispatch a request
     *
     * @param {String|Object} configOrUrl The config specific for this request (merged with this.defaults)
     * @param {?Object} config
     *
     * @returns {Promise} The Promise to be fulfilled
     */
    async request(configOrUrl, config) {
      try {
        return await this._request(configOrUrl, config);
      } catch (err) {
        if (err instanceof Error) {
          let dummy = {};
          Error.captureStackTrace ? Error.captureStackTrace(dummy) : dummy = new Error();
          const stack = dummy.stack ? dummy.stack.replace(/^.+\n/, "") : "";
          try {
            if (!err.stack) {
              err.stack = stack;
            } else if (stack && !String(err.stack).endsWith(stack.replace(/^.+\n.+\n/, ""))) {
              err.stack += "\n" + stack;
            }
          } catch (e) {
          }
        }
        throw err;
      }
    }
    _request(configOrUrl, config) {
      if (typeof configOrUrl === "string") {
        config = config || {};
        config.url = configOrUrl;
      } else {
        config = configOrUrl || {};
      }
      config = mergeConfig$1(this.defaults, config);
      const { transitional, paramsSerializer, headers } = config;
      if (transitional !== void 0) {
        validator.assertOptions(transitional, {
          silentJSONParsing: validators.transitional(validators.boolean),
          forcedJSONParsing: validators.transitional(validators.boolean),
          clarifyTimeoutError: validators.transitional(validators.boolean)
        }, false);
      }
      if (paramsSerializer != null) {
        if (utils$1.isFunction(paramsSerializer)) {
          config.paramsSerializer = {
            serialize: paramsSerializer
          };
        } else {
          validator.assertOptions(paramsSerializer, {
            encode: validators.function,
            serialize: validators.function
          }, true);
        }
      }
      if (config.allowAbsoluteUrls !== void 0) ;
      else if (this.defaults.allowAbsoluteUrls !== void 0) {
        config.allowAbsoluteUrls = this.defaults.allowAbsoluteUrls;
      } else {
        config.allowAbsoluteUrls = true;
      }
      validator.assertOptions(config, {
        baseUrl: validators.spelling("baseURL"),
        withXsrfToken: validators.spelling("withXSRFToken")
      }, true);
      config.method = (config.method || this.defaults.method || "get").toLowerCase();
      let contextHeaders = headers && utils$1.merge(
        headers.common,
        headers[config.method]
      );
      headers && utils$1.forEach(
        ["delete", "get", "head", "post", "put", "patch", "common"],
        (method) => {
          delete headers[method];
        }
      );
      config.headers = AxiosHeaders$1.concat(contextHeaders, headers);
      const requestInterceptorChain = [];
      let synchronousRequestInterceptors = true;
      this.interceptors.request.forEach(function unshiftRequestInterceptors(interceptor) {
        if (typeof interceptor.runWhen === "function" && interceptor.runWhen(config) === false) {
          return;
        }
        synchronousRequestInterceptors = synchronousRequestInterceptors && interceptor.synchronous;
        requestInterceptorChain.unshift(interceptor.fulfilled, interceptor.rejected);
      });
      const responseInterceptorChain = [];
      this.interceptors.response.forEach(function pushResponseInterceptors(interceptor) {
        responseInterceptorChain.push(interceptor.fulfilled, interceptor.rejected);
      });
      let promise;
      let i = 0;
      let len;
      if (!synchronousRequestInterceptors) {
        const chain = [dispatchRequest.bind(this), void 0];
        chain.unshift.apply(chain, requestInterceptorChain);
        chain.push.apply(chain, responseInterceptorChain);
        len = chain.length;
        promise = Promise.resolve(config);
        while (i < len) {
          promise = promise.then(chain[i++], chain[i++]);
        }
        return promise;
      }
      len = requestInterceptorChain.length;
      let newConfig = config;
      i = 0;
      while (i < len) {
        const onFulfilled = requestInterceptorChain[i++];
        const onRejected = requestInterceptorChain[i++];
        try {
          newConfig = onFulfilled(newConfig);
        } catch (error) {
          onRejected.call(this, error);
          break;
        }
      }
      try {
        promise = dispatchRequest.call(this, newConfig);
      } catch (error) {
        return Promise.reject(error);
      }
      i = 0;
      len = responseInterceptorChain.length;
      while (i < len) {
        promise = promise.then(responseInterceptorChain[i++], responseInterceptorChain[i++]);
      }
      return promise;
    }
    getUri(config) {
      config = mergeConfig$1(this.defaults, config);
      const fullPath = buildFullPath(config.baseURL, config.url, config.allowAbsoluteUrls);
      return buildURL(fullPath, config.params, config.paramsSerializer);
    }
  };
  utils$1.forEach(["delete", "get", "head", "options"], function forEachMethodNoData(method) {
    Axios$1.prototype[method] = function(url, config) {
      return this.request(mergeConfig$1(config || {}, {
        method,
        url,
        data: (config || {}).data
      }));
    };
  });
  utils$1.forEach(["post", "put", "patch"], function forEachMethodWithData(method) {
    function generateHTTPMethod(isForm) {
      return function httpMethod(url, data, config) {
        return this.request(mergeConfig$1(config || {}, {
          method,
          headers: isForm ? {
            "Content-Type": "multipart/form-data"
          } : {},
          url,
          data
        }));
      };
    }
    Axios$1.prototype[method] = generateHTTPMethod();
    Axios$1.prototype[method + "Form"] = generateHTTPMethod(true);
  });
  let CancelToken$1 = class CancelToken2 {
    constructor(executor) {
      if (typeof executor !== "function") {
        throw new TypeError("executor must be a function.");
      }
      let resolvePromise;
      this.promise = new Promise(function promiseExecutor(resolve) {
        resolvePromise = resolve;
      });
      const token = this;
      this.promise.then((cancel) => {
        if (!token._listeners) return;
        let i = token._listeners.length;
        while (i-- > 0) {
          token._listeners[i](cancel);
        }
        token._listeners = null;
      });
      this.promise.then = (onfulfilled) => {
        let _resolve;
        const promise = new Promise((resolve) => {
          token.subscribe(resolve);
          _resolve = resolve;
        }).then(onfulfilled);
        promise.cancel = function reject() {
          token.unsubscribe(_resolve);
        };
        return promise;
      };
      executor(function cancel(message, config, request) {
        if (token.reason) {
          return;
        }
        token.reason = new CanceledError$1(message, config, request);
        resolvePromise(token.reason);
      });
    }
    /**
     * Throws a `CanceledError` if cancellation has been requested.
     */
    throwIfRequested() {
      if (this.reason) {
        throw this.reason;
      }
    }
    /**
     * Subscribe to the cancel signal
     */
    subscribe(listener) {
      if (this.reason) {
        listener(this.reason);
        return;
      }
      if (this._listeners) {
        this._listeners.push(listener);
      } else {
        this._listeners = [listener];
      }
    }
    /**
     * Unsubscribe from the cancel signal
     */
    unsubscribe(listener) {
      if (!this._listeners) {
        return;
      }
      const index = this._listeners.indexOf(listener);
      if (index !== -1) {
        this._listeners.splice(index, 1);
      }
    }
    toAbortSignal() {
      const controller = new AbortController();
      const abort = (err) => {
        controller.abort(err);
      };
      this.subscribe(abort);
      controller.signal.unsubscribe = () => this.unsubscribe(abort);
      return controller.signal;
    }
    /**
     * Returns an object that contains a new `CancelToken` and a function that, when called,
     * cancels the `CancelToken`.
     */
    static source() {
      let cancel;
      const token = new CancelToken2(function executor(c) {
        cancel = c;
      });
      return {
        token,
        cancel
      };
    }
  };
  function spread$1(callback) {
    return function wrap(arr) {
      return callback.apply(null, arr);
    };
  }
  function isAxiosError$1(payload) {
    return utils$1.isObject(payload) && payload.isAxiosError === true;
  }
  const HttpStatusCode$1 = {
    Continue: 100,
    SwitchingProtocols: 101,
    Processing: 102,
    EarlyHints: 103,
    Ok: 200,
    Created: 201,
    Accepted: 202,
    NonAuthoritativeInformation: 203,
    NoContent: 204,
    ResetContent: 205,
    PartialContent: 206,
    MultiStatus: 207,
    AlreadyReported: 208,
    ImUsed: 226,
    MultipleChoices: 300,
    MovedPermanently: 301,
    Found: 302,
    SeeOther: 303,
    NotModified: 304,
    UseProxy: 305,
    Unused: 306,
    TemporaryRedirect: 307,
    PermanentRedirect: 308,
    BadRequest: 400,
    Unauthorized: 401,
    PaymentRequired: 402,
    Forbidden: 403,
    NotFound: 404,
    MethodNotAllowed: 405,
    NotAcceptable: 406,
    ProxyAuthenticationRequired: 407,
    RequestTimeout: 408,
    Conflict: 409,
    Gone: 410,
    LengthRequired: 411,
    PreconditionFailed: 412,
    PayloadTooLarge: 413,
    UriTooLong: 414,
    UnsupportedMediaType: 415,
    RangeNotSatisfiable: 416,
    ExpectationFailed: 417,
    ImATeapot: 418,
    MisdirectedRequest: 421,
    UnprocessableEntity: 422,
    Locked: 423,
    FailedDependency: 424,
    TooEarly: 425,
    UpgradeRequired: 426,
    PreconditionRequired: 428,
    TooManyRequests: 429,
    RequestHeaderFieldsTooLarge: 431,
    UnavailableForLegalReasons: 451,
    InternalServerError: 500,
    NotImplemented: 501,
    BadGateway: 502,
    ServiceUnavailable: 503,
    GatewayTimeout: 504,
    HttpVersionNotSupported: 505,
    VariantAlsoNegotiates: 506,
    InsufficientStorage: 507,
    LoopDetected: 508,
    NotExtended: 510,
    NetworkAuthenticationRequired: 511
  };
  Object.entries(HttpStatusCode$1).forEach(([key, value]) => {
    HttpStatusCode$1[value] = key;
  });
  function createInstance(defaultConfig) {
    const context = new Axios$1(defaultConfig);
    const instance = bind(Axios$1.prototype.request, context);
    utils$1.extend(instance, Axios$1.prototype, context, { allOwnKeys: true });
    utils$1.extend(instance, context, null, { allOwnKeys: true });
    instance.create = function create(instanceConfig) {
      return createInstance(mergeConfig$1(defaultConfig, instanceConfig));
    };
    return instance;
  }
  const axios = createInstance(defaults);
  axios.Axios = Axios$1;
  axios.CanceledError = CanceledError$1;
  axios.CancelToken = CancelToken$1;
  axios.isCancel = isCancel$1;
  axios.VERSION = VERSION$1;
  axios.toFormData = toFormData$1;
  axios.AxiosError = AxiosError$1;
  axios.Cancel = axios.CanceledError;
  axios.all = function all2(promises) {
    return Promise.all(promises);
  };
  axios.spread = spread$1;
  axios.isAxiosError = isAxiosError$1;
  axios.mergeConfig = mergeConfig$1;
  axios.AxiosHeaders = AxiosHeaders$1;
  axios.formToJSON = (thing) => formDataToJSON(utils$1.isHTMLForm(thing) ? new FormData(thing) : thing);
  axios.getAdapter = adapters.getAdapter;
  axios.HttpStatusCode = HttpStatusCode$1;
  axios.default = axios;
  const {
    Axios,
    AxiosError,
    CanceledError,
    isCancel,
    CancelToken,
    VERSION,
    all,
    Cancel,
    isAxiosError,
    spread,
    toFormData,
    AxiosHeaders,
    HttpStatusCode,
    formToJSON,
    getAdapter,
    mergeConfig
  } = axios;
  /*! js-cookie v3.0.5 | MIT */
  function assign(target) {
    for (var i = 1; i < arguments.length; i++) {
      var source = arguments[i];
      for (var key in source) {
        target[key] = source[key];
      }
    }
    return target;
  }
  var defaultConverter = {
    read: function(value) {
      if (value[0] === '"') {
        value = value.slice(1, -1);
      }
      return value.replace(/(%[\dA-F]{2})+/gi, decodeURIComponent);
    },
    write: function(value) {
      return encodeURIComponent(value).replace(
        /%(2[346BF]|3[AC-F]|40|5[BDE]|60|7[BCD])/g,
        decodeURIComponent
      );
    }
  };
  function init(converter, defaultAttributes) {
    function set(name, value, attributes) {
      if (typeof document === "undefined") {
        return;
      }
      attributes = assign({}, defaultAttributes, attributes);
      if (typeof attributes.expires === "number") {
        attributes.expires = new Date(Date.now() + attributes.expires * 864e5);
      }
      if (attributes.expires) {
        attributes.expires = attributes.expires.toUTCString();
      }
      name = encodeURIComponent(name).replace(/%(2[346B]|5E|60|7C)/g, decodeURIComponent).replace(/[()]/g, escape);
      var stringifiedAttributes = "";
      for (var attributeName in attributes) {
        if (!attributes[attributeName]) {
          continue;
        }
        stringifiedAttributes += "; " + attributeName;
        if (attributes[attributeName] === true) {
          continue;
        }
        stringifiedAttributes += "=" + attributes[attributeName].split(";")[0];
      }
      return document.cookie = name + "=" + converter.write(value, name) + stringifiedAttributes;
    }
    function get(name) {
      if (typeof document === "undefined" || arguments.length && !name) {
        return;
      }
      var cookies2 = document.cookie ? document.cookie.split("; ") : [];
      var jar = {};
      for (var i = 0; i < cookies2.length; i++) {
        var parts = cookies2[i].split("=");
        var value = parts.slice(1).join("=");
        try {
          var found = decodeURIComponent(parts[0]);
          jar[found] = converter.read(value, found);
          if (name === found) {
            break;
          }
        } catch (e) {
        }
      }
      return name ? jar[name] : jar;
    }
    return Object.create(
      {
        set,
        get,
        remove: function(name, attributes) {
          set(
            name,
            "",
            assign({}, attributes, {
              expires: -1
            })
          );
        },
        withAttributes: function(attributes) {
          return init(this.converter, assign({}, this.attributes, attributes));
        },
        withConverter: function(converter2) {
          return init(assign({}, this.converter, converter2), this.attributes);
        }
      },
      {
        attributes: { value: Object.freeze(defaultAttributes) },
        converter: { value: Object.freeze(converter) }
      }
    );
  }
  var api = init(defaultConverter, { path: "/" });
  const SESSION_COOKIE_NAME = "growdesk_session";
  const SESSION_EXPIRY_DAYS = 7;
  const LOCAL_STORAGE_KEY = "growdesk_session_data";
  console.log("[WIDGET] GrowDeskConfig al inicializar widgetApi:", window.GrowDeskConfig);
  let apiConfig = {
    apiUrl: ((_a = window.GrowDeskConfig) == null ? void 0 : _a.apiUrl) || "http://localhost:8082",
    widgetId: "demo-widget",
    widgetToken: "demo-token"
  };
  console.log("[WIDGET] Configuración inicial de API:", apiConfig);
  const initializeFromScript = () => {
    var _a2, _b, _c, _d, _e;
    const script = document.getElementById("growdesk-widget");
    if (script) {
      console.log("[WIDGET] Obteniendo configuración desde script:", script);
      const apiUrl = script.getAttribute("data-api-url");
      const widgetId = script.getAttribute("data-widget-id");
      const widgetToken = script.getAttribute("data-widget-token");
      const brandName = script.getAttribute("data-brand-name");
      const welcomeMessage = script.getAttribute("data-welcome-message");
      const primaryColor = script.getAttribute("data-primary-color");
      const position = script.getAttribute("data-position");
      if (apiUrl) {
        console.log("[WIDGET] Usando apiUrl desde script:", apiUrl);
        apiConfig.apiUrl = apiUrl;
      }
      if (widgetId) apiConfig.widgetId = widgetId;
      if (widgetToken) apiConfig.widgetToken = widgetToken;
      window.GrowDeskConfig = {
        ...window.GrowDeskConfig || {},
        apiUrl: apiUrl || ((_a2 = window.GrowDeskConfig) == null ? void 0 : _a2.apiUrl),
        brandName: brandName || ((_b = window.GrowDeskConfig) == null ? void 0 : _b.brandName),
        welcomeMessage: welcomeMessage || ((_c = window.GrowDeskConfig) == null ? void 0 : _c.welcomeMessage),
        primaryColor: primaryColor || ((_d = window.GrowDeskConfig) == null ? void 0 : _d.primaryColor),
        position: position || ((_e = window.GrowDeskConfig) == null ? void 0 : _e.position)
      };
      console.log("[WIDGET] Widget configurado desde script:", apiConfig, window.GrowDeskConfig);
    } else {
      console.warn("[WIDGET] No se encontró el script para configurar el widget.");
    }
  };
  const configureWidgetApi = (config) => {
    if (config.apiUrl) apiConfig.apiUrl = config.apiUrl;
    if (config.widgetId) apiConfig.widgetId = config.widgetId;
    if (config.widgetToken) apiConfig.widgetToken = config.widgetToken;
    console.log("API del widget configurada:", apiConfig);
  };
  const saveSession = (data) => {
    const now2 = Math.floor(Date.now() / 1e3);
    const expiration = now2 + SESSION_EXPIRY_DAYS * 24 * 60 * 60;
    const sessionData = {
      ...data,
      exp: expiration
    };
    try {
      api.set(SESSION_COOKIE_NAME, JSON.stringify(sessionData), {
        expires: SESSION_EXPIRY_DAYS,
        path: "/",
        secure: window.location.protocol === "https:"
      });
    } catch (error) {
      console.warn("No se pudo guardar sesión en cookies, usando localStorage como alternativa");
    }
    try {
      localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(sessionData));
    } catch (error) {
      console.warn("No se pudo guardar en localStorage:", error);
    }
    return sessionData;
  };
  const getSession = () => {
    try {
      const sessionCookie = api.get(SESSION_COOKIE_NAME);
      if (sessionCookie) {
        const sessionData = JSON.parse(sessionCookie);
        if (sessionData.exp && sessionData.exp < Math.floor(Date.now() / 1e3)) {
          clearSession();
          return null;
        }
        return sessionData;
      }
      const localData = localStorage.getItem(LOCAL_STORAGE_KEY);
      if (localData) {
        const sessionData = JSON.parse(localData);
        if (sessionData.exp && sessionData.exp < Math.floor(Date.now() / 1e3)) {
          clearSession();
          return null;
        }
        return sessionData;
      }
      return null;
    } catch (error) {
      console.error("Error al obtener sesión:", error);
      clearSession();
      return null;
    }
  };
  const clearSession = () => {
    try {
      api.remove(SESSION_COOKIE_NAME, { path: "/" });
    } catch (error) {
      console.warn("Error al eliminar cookie:", error);
    }
    try {
      localStorage.removeItem(LOCAL_STORAGE_KEY);
    } catch (error) {
      console.warn("Error al eliminar datos del localStorage:", error);
    }
  };
  const useWidgetApi = () => {
    var _a2, _b;
    console.log("[WIDGET DEBUG] Creando cliente Axios con apiConfig:", apiConfig);
    console.log("[WIDGET DEBUG] URL base que se va a usar:", `${((_a2 = window.GrowDeskConfig) == null ? void 0 : _a2.apiUrl) || apiConfig.apiUrl}`);
    const apiClient = axios.create({
      baseURL: ((_b = window.GrowDeskConfig) == null ? void 0 : _b.apiUrl) || apiConfig.apiUrl,
      timeout: 1e4,
      headers: {
        "Content-Type": "application/json",
        "X-Widget-ID": apiConfig.widgetId,
        "X-Widget-Token": apiConfig.widgetToken
      },
      // No usar withCredentials para evitar problemas de CORS
      withCredentials: false
    });
    console.log("[WIDGET DEBUG] Cliente Axios creado con baseURL:", apiClient.defaults.baseURL);
    apiClient.interceptors.request.use((config) => {
      var _a3;
      config.baseURL = ((_a3 = window.GrowDeskConfig) == null ? void 0 : _a3.apiUrl) || apiConfig.apiUrl;
      console.log("[WIDGET DEBUG] Interceptor - Request a:", config.baseURL + config.url);
      const session = getSession();
      if (session) {
        config.headers["X-User-Name"] = session.name;
        config.headers["X-User-Email"] = session.email;
        if (session.ticketId) {
          config.headers["X-Ticket-ID"] = session.ticketId;
        }
      }
      return config;
    });
    const hasActiveSession = () => {
      return getSession() !== null;
    };
    const createTicket = async (data) => {
      var _a3, _b2;
      try {
        console.log("[WIDGET] Intentando crear ticket con los siguientes datos:", data);
        console.log("[WIDGET] URL de API configurada:", apiConfig.apiUrl);
        let response;
        try {
          const baseUrl = apiConfig.apiUrl;
          let fullUrl = baseUrl;
          if (!fullUrl.endsWith("/")) {
            fullUrl += "/";
          }
          if (!fullUrl.includes("/widget/")) {
            fullUrl += "widget/";
          }
          const url = `${fullUrl}tickets`;
          console.log("[WIDGET] URL final:", url);
          console.log("[WIDGET] Headers:", {
            "Content-Type": "application/json",
            "X-Widget-ID": apiConfig.widgetId,
            "X-Widget-Token": apiConfig.widgetToken
          });
          response = await axios.post(url, data, {
            headers: {
              "Content-Type": "application/json",
              "X-Widget-ID": apiConfig.widgetId,
              "X-Widget-Token": apiConfig.widgetToken
            }
          });
          console.log("[WIDGET] Ticket creado con éxito:", response.data);
        } catch (axiosError) {
          console.warn("[WIDGET] Error al usar Axios, intentando con fetch:", axiosError);
          console.warn("[WIDGET] Detalles del error:", {
            status: (_a3 = axiosError == null ? void 0 : axiosError.response) == null ? void 0 : _a3.status,
            message: axiosError == null ? void 0 : axiosError.message,
            data: (_b2 = axiosError == null ? void 0 : axiosError.response) == null ? void 0 : _b2.data
          });
          const baseUrl = apiConfig.apiUrl;
          let fullUrl = baseUrl;
          if (!fullUrl.endsWith("/")) {
            fullUrl += "/";
          }
          if (!fullUrl.includes("/widget/")) {
            fullUrl += "widget/";
          }
          const fetchUrl = `${fullUrl}tickets`;
          console.log("[WIDGET] Intentando fetch a:", fetchUrl);
          const fetchResponse = await fetch(fetchUrl, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              "X-Widget-ID": apiConfig.widgetId,
              "X-Widget-Token": apiConfig.widgetToken
            },
            body: JSON.stringify(data)
          });
          if (!fetchResponse.ok) {
            console.error(`[WIDGET] Error en fetch: ${fetchResponse.status} ${fetchResponse.statusText || "Error desconocido"}`);
            throw new Error(`Error ${fetchResponse.status}: ${fetchResponse.statusText || "Error desconocido"}`);
          }
          const responseData = await fetchResponse.json();
          console.log("[WIDGET] Datos recibidos con fetch:", responseData);
          response = { data: responseData };
        }
        if (response.data && response.data.id) {
          saveSession({
            name: data.name,
            email: data.email,
            ticketId: response.data.id
          });
        }
        return response.data;
      } catch (error) {
        console.error("Error creating ticket:", error);
        if (window.location.hostname === "localhost" || window.location.hostname === "127.0.0.1") {
          console.warn("Creando ticket simulado para entorno de desarrollo");
          const mockTicketId = `MOCK-${Date.now()}`;
          saveSession({
            name: data.name,
            email: data.email,
            ticketId: mockTicketId
          });
          return {
            ticketId: mockTicketId,
            status: "open",
            message: "Ticket simulado creado con éxito"
          };
        }
        throw error;
      }
    };
    const sendMessage = async (data) => {
      try {
        console.log("[WIDGET] Intentando enviar mensaje a ticket:", data.ticketId);
        console.log("[WIDGET] URL de API configurada:", apiConfig.apiUrl);
        const baseUrl = apiConfig.apiUrl.endsWith("/") ? apiConfig.apiUrl : `${apiConfig.apiUrl}/`;
        const adjustedBaseUrl = baseUrl.includes("widget") ? baseUrl : `${baseUrl}widget/`;
        const messagesUrl = `${adjustedBaseUrl}messages`;
        console.log("[WIDGET] URL base ajustada para enviar mensaje:", messagesUrl);
        console.log("[WIDGET] Datos a enviar:", {
          ticketId: data.ticketId,
          message: data.message
        });
        let response;
        try {
          console.log("[WIDGET] Usando axios.post para enviar mensaje");
          response = await axios.post(messagesUrl, {
            ticketId: data.ticketId,
            message: data.message
          }, {
            headers: {
              "Content-Type": "application/json",
              "X-Widget-ID": apiConfig.widgetId,
              "X-Widget-Token": apiConfig.widgetToken,
              "X-User-Name": data.userName || "",
              "X-User-Email": data.userEmail || ""
            }
          });
          console.log("[WIDGET] Mensaje enviado con éxito:", response.data);
        } catch (axiosError) {
          console.warn("[WIDGET] Error al usar Axios para enviar mensaje, intentando con fetch:", axiosError);
          console.log("[WIDGET] Intentando fetch a:", messagesUrl);
          const fetchResponse = await fetch(messagesUrl, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              "X-Widget-ID": apiConfig.widgetId,
              "X-Widget-Token": apiConfig.widgetToken,
              "X-User-Name": data.userName || "",
              "X-User-Email": data.userEmail || ""
            },
            body: JSON.stringify({
              ticketId: data.ticketId,
              message: data.message
            })
          });
          if (!fetchResponse.ok) {
            console.error(`[WIDGET] Error en fetch: ${fetchResponse.status} ${fetchResponse.statusText || "Error desconocido"}`);
            throw new Error(`Error ${fetchResponse.status}: ${fetchResponse.statusText || "Error desconocido"}`);
          }
          const responseData = await fetchResponse.json();
          console.log("[WIDGET] Datos recibidos con fetch:", responseData);
          response = { data: responseData };
        }
        return response.data;
      } catch (error) {
        console.error("[WIDGET] Error sending message:", error);
        throw error;
      }
    };
    const getMessageHistory = async (ticketId) => {
      try {
        console.log("[WIDGET] Intentando obtener historial de mensajes para ticket:", ticketId);
        console.log("[WIDGET] URL de API configurada:", apiConfig.apiUrl);
        let response;
        try {
          console.log("[WIDGET] Obteniendo mensajes de:", `${apiConfig.apiUrl}/widget/tickets/${ticketId}/messages`);
          response = await axios.get(`${apiConfig.apiUrl}/widget/tickets/${ticketId}/messages`, {
            headers: {
              "Content-Type": "application/json",
              "X-Widget-ID": apiConfig.widgetId,
              "X-Widget-Token": apiConfig.widgetToken,
              "X-Ticket-ID": ticketId
            }
          });
          console.log("[WIDGET] Mensajes obtenidos con éxito:", response.data);
        } catch (axiosError) {
          console.warn("[WIDGET] Error al usar Axios para obtener mensajes, intentando con fetch:", axiosError);
          const fetchUrl = `${apiConfig.apiUrl}/widget/tickets/${ticketId}/messages`;
          console.log("[WIDGET] Intentando fetch a:", fetchUrl);
          const fetchResponse = await fetch(fetchUrl, {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
              "X-Widget-ID": apiConfig.widgetId,
              "X-Widget-Token": apiConfig.widgetToken,
              "X-Ticket-ID": ticketId
            }
          });
          if (!fetchResponse.ok) {
            throw new Error(`Error ${fetchResponse.status}: ${fetchResponse.statusText || "Error desconocido"}`);
          }
          response = { data: await fetchResponse.json() };
        }
        return response.data;
      } catch (error) {
        console.error("Error getting message history:", error);
        if (window.location.hostname === "localhost" || window.location.hostname === "127.0.0.1") {
          console.warn("Devolviendo historial de mensajes simulado para entorno de desarrollo");
          return {
            messages: [
              {
                id: "mock-msg-1",
                content: "Bienvenido al chat de soporte simulado.",
                isClient: false,
                timestamp: (/* @__PURE__ */ new Date()).toISOString()
              }
            ]
          };
        }
        throw error;
      }
    };
    const getFaqs = async () => {
      try {
        console.log("[WIDGET] Intentando obtener FAQs");
        console.log("[WIDGET] URL de API configurada:", apiConfig.apiUrl);
        let response;
        try {
          console.log("[WIDGET] Obteniendo FAQs de:", `${apiConfig.apiUrl}/widget/faqs`);
          response = await axios.get(`${apiConfig.apiUrl}/widget/faqs`, {
            headers: {
              "Content-Type": "application/json",
              "X-Widget-ID": apiConfig.widgetId,
              "X-Widget-Token": apiConfig.widgetToken
            }
          });
          console.log("[WIDGET] FAQs obtenidas con éxito:", response.data);
        } catch (axiosError) {
          console.warn("[WIDGET] Error al usar Axios para obtener FAQs, intentando con fetch:", axiosError);
          const fetchUrl = `${apiConfig.apiUrl}/widget/faqs`;
          console.log("[WIDGET] Intentando fetch a:", fetchUrl);
          const fetchResponse = await fetch(fetchUrl, {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
              "X-Widget-ID": apiConfig.widgetId,
              "X-Widget-Token": apiConfig.widgetToken
            }
          });
          if (!fetchResponse.ok) {
            const statusText = fetchResponse.statusText || "Error desconocido";
            console.error(`[WIDGET] Error en fetch: ${fetchResponse.status} ${statusText}`);
            throw new Error(`Error ${fetchResponse.status}: ${statusText}`);
          }
          const responseData = await fetchResponse.json();
          console.log("[WIDGET] Datos recibidos con fetch:", responseData);
          response = { data: responseData };
        }
        console.log("FAQs obtenidas correctamente:", response.data);
        return response.data;
      } catch (error) {
        console.error("Error obteniendo FAQs:", error);
        console.error("Configuración actual:", {
          apiUrl: apiConfig.apiUrl,
          widgetId: apiConfig.widgetId
        });
        if (window.location.hostname === "localhost" || window.location.hostname === "127.0.0.1") {
          console.warn("Devolviendo FAQs simuladas para entorno de desarrollo");
          return [
            {
              id: 1,
              question: "¿Cómo puedo usar este widget?",
              answer: "Este widget se puede integrar en cualquier sitio web añadiendo el script correspondiente.",
              category: "General",
              isPublished: true
            },
            {
              id: 2,
              question: "¿Cómo puedo contactar con soporte?",
              answer: "Puedes hacerlo directamente a través de este chat.",
              category: "Soporte",
              isPublished: true
            }
          ];
        }
        return [];
      }
    };
    const logout = () => {
      clearSession();
    };
    return {
      hasActiveSession,
      getSession,
      createTicket,
      sendMessage,
      getMessageHistory,
      logout,
      getFaqs
    };
  };
  const _sfc_main = /* @__PURE__ */ vue.defineComponent({
    __name: "ChatWidget",
    props: {
      primaryColor: {
        type: String,
        default: "#6200ea"
      },
      brandName: {
        type: String,
        default: "Chat Support"
      },
      position: {
        type: String,
        default: "bottom-right"
      },
      welcomeMessage: {
        type: String,
        default: "Para iniciar tu consulta, por favor completa el siguiente formulario:"
      }
    },
    setup(__props, { expose: __expose }) {
      __expose();
      const props = __props;
      const darkenColor = (color, amount) => {
        let usePound = false;
        if (color[0] === "#") {
          color = color.slice(1);
          usePound = true;
        }
        const num = parseInt(color, 16);
        let r = (num >> 16) - amount;
        if (r > 255) r = 255;
        else if (r < 0) r = 0;
        let g = (num >> 8 & 255) - amount;
        if (g > 255) g = 255;
        else if (g < 0) g = 0;
        let b = (num & 255) - amount;
        if (b > 255) b = 255;
        else if (b < 0) b = 0;
        return (usePound ? "#" : "") + (g | r << 8 | b << 16).toString(16).padStart(6, "0");
      };
      const formatMessageTime = (timestamp) => {
        try {
          const date = new Date(timestamp);
          return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
        } catch (error) {
          return "";
        }
      };
      const isOpen = vue.ref(false);
      const isRegistered = vue.ref(false);
      const loading = vue.ref(false);
      const currentTicketId = vue.ref("");
      const messages = vue.ref([]);
      const newMessage = vue.ref("");
      const webSocket = vue.ref(null);
      const showChatView = vue.ref(false);
      const faqs = vue.ref([]);
      const loadingFaqs = vue.ref(false);
      const expandedFaqs = vue.ref([]);
      const userData = vue.ref({
        name: "",
        email: "",
        initialMessage: ""
      });
      const api2 = useWidgetApi();
      const hasSession = vue.computed(() => {
        return api2.hasActiveSession();
      });
      const faqCategories = vue.computed(() => {
        const categories = /* @__PURE__ */ new Set();
        faqs.value.forEach((faq) => categories.add(faq.category));
        return Array.from(categories);
      });
      const getFaqsByCategory = (category) => {
        return faqs.value.filter((faq) => faq.category === category);
      };
      const toggleFaq = (id) => {
        const index = expandedFaqs.value.indexOf(id);
        if (index === -1) {
          expandedFaqs.value.push(id);
        } else {
          expandedFaqs.value.splice(index, 1);
        }
      };
      const loadFaqs = async () => {
        loadingFaqs.value = true;
        console.log("Iniciando carga de FAQs desde el servidor...");
        try {
          console.log("Llamando a api.getFaqs()...");
          const response = await api2.getFaqs();
          console.log("Respuesta de API para FAQs:", response);
          if (response && Array.isArray(response)) {
            const publishedFaqs = response.filter((faq) => faq.isPublished);
            console.log(`FAQs publicadas encontradas: ${publishedFaqs.length}`);
            faqs.value = publishedFaqs;
            if (publishedFaqs.length === 0) {
              console.warn("No se encontraron FAQs publicadas para mostrar");
            } else {
              const categories = /* @__PURE__ */ new Set();
              const firstFaqsIds = [];
              faqs.value.forEach((faq) => {
                if (!categories.has(faq.category)) {
                  categories.add(faq.category);
                  firstFaqsIds.push(faq.id);
                }
              });
              expandedFaqs.value = firstFaqsIds;
              console.log("Categorías de FAQs encontradas:", Array.from(categories));
            }
          } else {
            console.error("Formato inesperado en la respuesta de FAQs:", response);
            faqs.value = [];
          }
        } catch (error) {
          console.error("Error al cargar FAQs:", error);
          faqs.value = [];
          if (isRegistered.value) {
            showChatView.value = true;
            console.log("Cambiando a vista de chat debido a error en FAQs");
          }
        } finally {
          loadingFaqs.value = false;
          console.log("Carga de FAQs finalizada. Estado:", {
            faqsCount: faqs.value.length,
            expandedFaqsCount: expandedFaqs.value.length,
            showingChat: showChatView.value
          });
        }
      };
      const refreshFaqs = () => {
        loadFaqs();
      };
      const setInitialQuestion = (question) => {
        newMessage.value = question;
        showChatView.value = true;
        setTimeout(() => {
          sendMessage();
        }, 100);
      };
      const connectWebSocket = (ticketId) => {
        if (!ticketId) {
          console.error("No se puede conectar al WebSocket sin un ID de ticket");
          return;
        }
        if (webSocket.value) {
          webSocket.value.close();
          webSocket.value = null;
        }
        const effectiveTicketId = ticketId;
        const apiUrl = new URL(apiConfig.apiUrl);
        const wsProtocol = apiUrl.protocol === "https:" ? "wss:" : "ws:";
        const hostPart = apiUrl.hostname;
        const portPart = apiUrl.port || (apiUrl.protocol === "https:" ? "443" : "80");
        const wsUrl = `${wsProtocol}//${hostPart}:${portPart}/api/ws/chat/${effectiveTicketId}`;
        console.log("Información de conexión WebSocket:", {
          apiUrl: apiConfig.apiUrl,
          parsedHost: hostPart,
          parsedPort: portPart,
          wsProtocol,
          finalWsUrl: wsUrl,
          ticketId: effectiveTicketId
        });
        try {
          webSocket.value = new WebSocket(wsUrl);
          webSocket.value.onopen = () => {
            console.log("Conexión WebSocket establecida correctamente");
          };
          webSocket.value.onmessage = (event) => {
            var _a2, _b, _c, _d;
            console.log("Mensaje WebSocket recibido:", event.data);
            try {
              const data = JSON.parse(event.data);
              console.log("Datos WebSocket procesados:", {
                type: data.type,
                hasData: !!data.data,
                hasMessage: !!data.message,
                dataContent: (_a2 = data.data) == null ? void 0 : _a2.content,
                messageContent: (_b = data.message) == null ? void 0 : _b.content,
                dataIsClient: (_c = data.data) == null ? void 0 : _c.isClient,
                messageIsClient: (_d = data.message) == null ? void 0 : _d.isClient
              });
              let messageContent = "";
              let isClientMessage = false;
              let messageObj = null;
              if (data.type === "new_message") {
                if (data.message && data.message.content) {
                  messageContent = data.message.content;
                  isClientMessage = data.message.isClient === true;
                  messageObj = data.message;
                } else if (data.data && data.data.content) {
                  messageContent = data.data.content;
                  isClientMessage = data.data.isClient === true;
                  messageObj = data.data;
                }
              } else if (data.type === "message_received") {
                if (data.data && data.data.content) {
                  messageContent = data.data.content;
                  isClientMessage = data.data.isClient === true;
                  messageObj = data.data;
                }
              } else if (data.content) {
                messageContent = data.content;
                isClientMessage = data.isClient === true;
                messageObj = data;
              }
              if (messageContent) {
                console.log("Mensaje extraído para mostrar:", {
                  content: messageContent,
                  isClientMessage,
                  messageId: (messageObj == null ? void 0 : messageObj.id) || "no-id"
                });
                const isDuplicate = messages.value.some(
                  (msg) => msg.text === messageContent && msg.isUser === isClientMessage
                );
                if (!isDuplicate) {
                  messages.value.push({
                    text: messageContent,
                    isUser: isClientMessage
                    // isUser=true para mensajes del cliente, false para mensajes del agente
                  });
                  scrollToBottom();
                } else {
                  console.log("Mensaje duplicado detectado y omitido");
                }
              } else if (data.type === "error") {
                console.error("Error del servidor WebSocket:", data.message || data.data || "Error desconocido");
              } else if (data.type === "connection_established" || data.type === "identify_success") {
                console.log("Conexión WebSocket confirmada:", data.type);
              } else {
                console.warn("Formato de mensaje WebSocket no reconocido:", data);
              }
            } catch (error) {
              console.error("Error al procesar mensaje WebSocket:", error);
            }
          };
          webSocket.value.onerror = (error) => {
            console.error("Error en conexión WebSocket:", error);
          };
          webSocket.value.onclose = (event) => {
            console.log("Conexión WebSocket cerrada. Código:", event.code, "Razón:", event.reason);
            if (isOpen.value && isRegistered.value) {
              setTimeout(() => {
                connectWebSocket(currentTicketId.value);
              }, 5e3);
            }
          };
        } catch (error) {
          console.error("Error al crear conexión WebSocket:", error);
        }
      };
      vue.onMounted(() => {
        const container = document.getElementById("growdesk-widget-container");
        if (container) {
          container.addEventListener("open-widget", () => {
            isOpen.value = true;
          });
          container.addEventListener("close-widget", () => {
            isOpen.value = false;
          });
        }
        if (hasSession.value) {
          const session = api2.getSession();
          if (session) {
            userData.value.name = session.name;
            userData.value.email = session.email;
            currentTicketId.value = session.ticketId;
            isRegistered.value = true;
            showChatView.value = false;
            loadFaqs();
            loadPreviousMessages(session.ticketId);
            connectWebSocket(session.ticketId);
          }
        }
      });
      vue.onBeforeUnmount(() => {
        if (webSocket.value) {
          webSocket.value.close();
          webSocket.value = null;
        }
      });
      const toggleChat = () => {
        isOpen.value = !isOpen.value;
      };
      const submitRegistration = async () => {
        if (!userData.value.name || !userData.value.email || !userData.value.initialMessage) return;
        loading.value = true;
        try {
          const ticketData = {
            name: userData.value.name,
            email: userData.value.email,
            message: userData.value.initialMessage,
            // Usamos el mensaje inicial proporcionado por el usuario
            metadata: {
              url: window.location.href,
              referrer: document.referrer,
              userAgent: navigator.userAgent,
              screenSize: `${window.innerWidth}x${window.innerHeight}`
            }
          };
          const response = await api2.createTicket(ticketData);
          currentTicketId.value = response.ticketId;
          isRegistered.value = true;
          showChatView.value = false;
          await loadFaqs();
          console.log("FAQs cargadas después del registro:", faqs.value);
          messages.value.push({
            text: userData.value.initialMessage,
            isUser: true
          });
          messages.value.push({
            text: `Hola ${userData.value.name}, gracias por contactarnos. Puedes consultar nuestras preguntas frecuentes o iniciar un chat con nosotros.`,
            isUser: false
          });
          connectWebSocket(response.ticketId);
          scrollToBottom();
        } catch (error) {
          console.error("Error al registrar al usuario:", error);
          alert("Lo sentimos, hubo un problema al iniciar el chat. Por favor, inténtalo de nuevo más tarde.");
        } finally {
          loading.value = false;
        }
      };
      const loadPreviousMessages = async (ticketId) => {
        try {
          loading.value = true;
          const response = await api2.getMessageHistory(ticketId);
          if (response && response.messages && response.messages.length > 0) {
            console.log("Histórico de mensajes recibido:", response.messages);
            const formattedMessages = response.messages.map((msg) => {
              const isClientMessage = msg.isClient === true || typeof msg.isClient === "string" && msg.isClient.toLowerCase() === "true";
              console.log(`Mensaje procesado - contenido: "${msg.text || msg.content}", isClient: ${isClientMessage}`);
              return {
                text: msg.text || msg.content,
                isUser: isClientMessage
                // Los mensajes del cliente (isClient=true) son los del usuario
              };
            });
            messages.value = formattedMessages;
            scrollToBottom();
          } else {
            messages.value.push({
              text: `Hola ${userData.value.name}, bienvenido de nuevo. ¿En qué podemos ayudarte hoy?`,
              isUser: false
            });
          }
        } catch (error) {
          console.error("Error al cargar mensajes anteriores:", error);
          messages.value.push({
            text: `Hola ${userData.value.name}, parece que hubo un problema al cargar tus mensajes anteriores. ¿En qué podemos ayudarte hoy?`,
            isUser: false
          });
        } finally {
          loading.value = false;
        }
      };
      const scrollToBottom = () => {
        setTimeout(() => {
          const messagesContainer = document.querySelector(".chat-messages");
          if (messagesContainer) {
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
          }
        }, 100);
      };
      const logout = () => {
        api2.logout();
        isRegistered.value = false;
        currentTicketId.value = "";
        messages.value = [];
        faqs.value = [];
        expandedFaqs.value = [];
        userData.value = {
          name: "",
          email: "",
          initialMessage: ""
        };
      };
      const sendMessage = async () => {
        if (newMessage.value.trim() === "" || !currentTicketId.value) return;
        const userMessage = newMessage.value;
        newMessage.value = "";
        const messageId = `temp-${Date.now()}`;
        messages.value.push({
          id: messageId,
          text: userMessage,
          isUser: true,
          pending: true
        });
        scrollToBottom();
        loading.value = true;
        try {
          const response = await api2.sendMessage({
            ticketId: currentTicketId.value,
            message: userMessage,
            userName: userData.value.name,
            userEmail: userData.value.email
          });
          const index = messages.value.findIndex((msg) => msg.id === messageId);
          if (index >= 0) {
            messages.value[index].pending = false;
            if (response.messageId === "duplicate-ignored") {
              console.log("Mensaje detectado como duplicado por el servidor");
            }
          }
          if (response.messageId !== "duplicate-ignored") {
          }
          loading.value = false;
        } catch (error) {
          console.error("Error al enviar el mensaje:", error);
          const index = messages.value.findIndex((msg) => msg.id === messageId);
          if (index >= 0) {
            messages.value[index].error = true;
            messages.value[index].pending = false;
          }
          messages.value.push({
            text: "Lo sentimos, hubo un problema al enviar tu mensaje. Por favor, inténtalo de nuevo.",
            isUser: false
          });
          scrollToBottom();
          loading.value = false;
        }
      };
      const __returned__ = { props, darkenColor, formatMessageTime, isOpen, isRegistered, loading, currentTicketId, messages, newMessage, webSocket, showChatView, faqs, loadingFaqs, expandedFaqs, userData, api: api2, hasSession, faqCategories, getFaqsByCategory, toggleFaq, loadFaqs, refreshFaqs, setInitialQuestion, connectWebSocket, toggleChat, submitRegistration, loadPreviousMessages, scrollToBottom, logout, sendMessage };
      Object.defineProperty(__returned__, "__isScriptSetup", { enumerable: false, value: true });
      return __returned__;
    }
  });
  const _export_sfc = (sfc, props) => {
    const target = sfc.__vccOpts || sfc;
    for (const [key, val] of props) {
      target[key] = val;
    }
    return target;
  };
  const _hoisted_1 = {
    key: 0,
    class: "pi pi-comments text-white text-3xl"
  };
  const _hoisted_2 = {
    key: 1,
    class: "pi pi-times text-white text-3xl"
  };
  const _hoisted_3 = {
    key: 0,
    class: "fixed bottom-24 right-6 w-[420px] h-[600px] bg-white rounded-2xl shadow-2xl flex flex-col overflow-hidden border border-gray-200 z-40 transition-all duration-300 ease-in-out animate-slide-up"
  };
  const _hoisted_4 = { class: "flex items-center" };
  const _hoisted_5 = { class: "font-semibold text-white text-lg" };
  const _hoisted_6 = { class: "flex items-center space-x-3" };
  const _hoisted_7 = {
    key: 0,
    class: "flex-1 p-8 overflow-y-auto bg-white flex flex-col animate-fade-in"
  };
  const _hoisted_8 = { class: "text-center text-gray-700 mb-8" };
  const _hoisted_9 = { class: "flex flex-col" };
  const _hoisted_10 = { class: "flex flex-col" };
  const _hoisted_11 = { class: "flex flex-col" };
  const _hoisted_12 = ["disabled"];
  const _hoisted_13 = {
    key: 0,
    class: "flex items-center"
  };
  const _hoisted_14 = {
    key: 1,
    class: "flex items-center"
  };
  const _hoisted_15 = { class: "flex-1 overflow-y-auto bg-white" };
  const _hoisted_16 = {
    key: 0,
    class: "p-6 chat-messages"
  };
  const _hoisted_17 = {
    key: 0,
    class: "flex flex-col items-center justify-center h-full"
  };
  const _hoisted_18 = { class: "text-center text-gray-500" };
  const _hoisted_19 = { class: "flex items-start" };
  const _hoisted_20 = {
    key: 0,
    class: "w-9 h-9 rounded-full bg-white border border-gray-200 flex items-center justify-center mr-3 flex-shrink-0"
  };
  const _hoisted_21 = { class: "flex-1" };
  const _hoisted_22 = {
    key: 0,
    class: "font-semibold text-sm mb-1 text-gray-600"
  };
  const _hoisted_23 = { class: "text-base" };
  const _hoisted_24 = {
    key: 0,
    class: "text-right mt-1"
  };
  const _hoisted_25 = {
    key: 1,
    class: "text-right mt-1 text-xs text-red-500"
  };
  const _hoisted_26 = { class: "p-6" };
  const _hoisted_27 = { class: "text-center mb-6" };
  const _hoisted_28 = { class: "text-xl font-bold mb-5 text-gray-700 flex items-center" };
  const _hoisted_29 = {
    key: 0,
    class: "text-center py-6"
  };
  const _hoisted_30 = {
    key: 1,
    class: "text-center py-6 text-gray-500"
  };
  const _hoisted_31 = {
    key: 2,
    class: "space-y-6"
  };
  const _hoisted_32 = { class: "font-semibold text-gray-700 mb-4 pb-2 border-b border-gray-200" };
  const _hoisted_33 = { class: "space-y-4" };
  const _hoisted_34 = ["onClick"];
  const _hoisted_35 = { class: "font-medium text-gray-800 flex-1" };
  const _hoisted_36 = {
    key: 0,
    class: "p-4 bg-gray-50 animate-fade-in"
  };
  const _hoisted_37 = { class: "text-gray-700" };
  const _hoisted_38 = { class: "mt-4 flex justify-end" };
  const _hoisted_39 = ["onClick"];
  const _hoisted_40 = {
    key: 2,
    class: "p-5 border-t border-gray-200 bg-white"
  };
  const _hoisted_41 = { class: "flex items-center" };
  const _hoisted_42 = {
    key: 3,
    class: "p-5 border-t border-gray-200 bg-white"
  };
  const _hoisted_43 = { class: "flex justify-between items-center" };
  function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
    return vue.openBlock(), vue.createElementBlock("div", null, [
      vue.createCommentVNode(" Chat bubble con animación de pulso restaurada "),
      vue.createElementVNode(
        "div",
        {
          class: "fixed bottom-6 right-6 w-20 h-20 rounded-full flex items-center justify-center cursor-pointer shadow-xl transition-all z-50 hover:scale-110 animate-pulse-slow",
          style: vue.normalizeStyle({ backgroundColor: $props.primaryColor }),
          onClick: $setup.toggleChat
        },
        [
          !$setup.isOpen ? (vue.openBlock(), vue.createElementBlock("i", _hoisted_1)) : (vue.openBlock(), vue.createElementBlock("i", _hoisted_2))
        ],
        4
        /* STYLE */
      ),
      vue.createCommentVNode(" Chat interface con sombras restauradas "),
      $setup.isOpen ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_3, [
        vue.createCommentVNode(" Chat header - Diseño plano (sin gradiente) "),
        vue.createElementVNode(
          "div",
          {
            class: "p-5 flex justify-between items-center shadow-sm",
            style: vue.normalizeStyle({ backgroundColor: $props.primaryColor })
          },
          [
            vue.createElementVNode("div", _hoisted_4, [
              vue.createElementVNode(
                "h3",
                _hoisted_5,
                vue.toDisplayString($props.brandName),
                1
                /* TEXT */
              )
            ]),
            vue.createElementVNode("div", _hoisted_6, [
              $setup.isRegistered ? (vue.openBlock(), vue.createElementBlock("i", {
                key: 0,
                class: "pi pi-sign-out cursor-pointer text-white hover:scale-110 transition-transform",
                onClick: vue.withModifiers($setup.logout, ["stop"]),
                title: "Cerrar sesión"
              })) : vue.createCommentVNode("v-if", true),
              vue.createElementVNode("i", {
                class: "pi pi-times cursor-pointer text-white hover:scale-110 transition-transform",
                onClick: $setup.toggleChat
              })
            ])
          ],
          4
          /* STYLE */
        ),
        vue.createCommentVNode(" Registration form sin efectos de sombra "),
        !$setup.isRegistered ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_7, [
          vue.createElementVNode("div", _hoisted_8, [
            vue.createElementVNode(
              "div",
              {
                class: "text-2xl font-semibold mb-2",
                style: vue.normalizeStyle({ color: $props.primaryColor })
              },
              "Bienvenido",
              4
              /* STYLE */
            ),
            vue.createElementVNode(
              "p",
              null,
              vue.toDisplayString($props.welcomeMessage),
              1
              /* TEXT */
            )
          ]),
          vue.createElementVNode(
            "form",
            {
              onSubmit: vue.withModifiers($setup.submitRegistration, ["prevent"]),
              class: "flex flex-col gap-6"
            },
            [
              vue.createElementVNode("div", _hoisted_9, [
                _cache[6] || (_cache[6] = vue.createElementVNode(
                  "label",
                  {
                    for: "name",
                    class: "text-sm font-medium text-gray-700 mb-2"
                  },
                  "Nombre",
                  -1
                  /* HOISTED */
                )),
                vue.withDirectives(vue.createElementVNode(
                  "input",
                  {
                    "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => $setup.userData.name = $event),
                    type: "text",
                    id: "name",
                    class: "border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:border-2 transition-all",
                    style: vue.normalizeStyle({ "--tw-border-opacity": 1, borderColor: $props.primaryColor }),
                    required: ""
                  },
                  null,
                  4
                  /* STYLE */
                ), [
                  [vue.vModelText, $setup.userData.name]
                ])
              ]),
              vue.createElementVNode("div", _hoisted_10, [
                _cache[7] || (_cache[7] = vue.createElementVNode(
                  "label",
                  {
                    for: "email",
                    class: "text-sm font-medium text-gray-700 mb-2"
                  },
                  "Email",
                  -1
                  /* HOISTED */
                )),
                vue.withDirectives(vue.createElementVNode(
                  "input",
                  {
                    "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => $setup.userData.email = $event),
                    type: "email",
                    id: "email",
                    class: "border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:border-2 transition-all",
                    style: vue.normalizeStyle({ "--tw-border-opacity": 1, borderColor: $props.primaryColor }),
                    required: ""
                  },
                  null,
                  4
                  /* STYLE */
                ), [
                  [vue.vModelText, $setup.userData.email]
                ])
              ]),
              vue.createElementVNode("div", _hoisted_11, [
                _cache[8] || (_cache[8] = vue.createElementVNode(
                  "label",
                  {
                    for: "firstMessage",
                    class: "text-sm font-medium text-gray-700 mb-2"
                  },
                  "¿En qué podemos ayudarte?",
                  -1
                  /* HOISTED */
                )),
                vue.withDirectives(vue.createElementVNode(
                  "textarea",
                  {
                    "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => $setup.userData.initialMessage = $event),
                    id: "firstMessage",
                    rows: "5",
                    class: "border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:border-2 transition-all",
                    style: vue.normalizeStyle({ "--tw-border-opacity": 1, borderColor: $props.primaryColor }),
                    required: "",
                    placeholder: "Describe brevemente tu consulta..."
                  },
                  null,
                  4
                  /* STYLE */
                ), [
                  [vue.vModelText, $setup.userData.initialMessage]
                ])
              ]),
              vue.createElementVNode("button", {
                type: "submit",
                class: "mt-2 text-white rounded-lg py-4 font-medium focus:outline-none transition-all hover:shadow-lg flex items-center justify-center",
                style: vue.normalizeStyle({ backgroundColor: $props.primaryColor }),
                disabled: $setup.loading
              }, [
                !$setup.loading ? (vue.openBlock(), vue.createElementBlock("span", _hoisted_13, _cache[9] || (_cache[9] = [
                  vue.createElementVNode(
                    "i",
                    { class: "pi pi-send mr-2" },
                    null,
                    -1
                    /* HOISTED */
                  ),
                  vue.createTextVNode(" Iniciar chat")
                ]))) : (vue.openBlock(), vue.createElementBlock("span", _hoisted_14, _cache[10] || (_cache[10] = [
                  vue.createElementVNode(
                    "i",
                    { class: "pi pi-spin pi-spinner mr-2" },
                    null,
                    -1
                    /* HOISTED */
                  ),
                  vue.createTextVNode(" Procesando... ")
                ])))
              ], 12, _hoisted_12)
            ],
            32
            /* NEED_HYDRATION */
          )
        ])) : (vue.openBlock(), vue.createElementBlock(
          vue.Fragment,
          { key: 1 },
          [
            vue.createCommentVNode(" FAQ and Chat sections (displayed after registration) "),
            vue.createElementVNode("div", _hoisted_15, [
              vue.createCommentVNode(" Chat View con mensajes mejorados "),
              $setup.showChatView ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_16, [
                $setup.messages.length === 0 ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_17, [
                  vue.createElementVNode("div", _hoisted_18, [
                    vue.createElementVNode(
                      "i",
                      {
                        class: "pi pi-comments text-5xl mb-4",
                        style: vue.normalizeStyle({ color: $props.primaryColor })
                      },
                      null,
                      4
                      /* STYLE */
                    ),
                    _cache[11] || (_cache[11] = vue.createElementVNode(
                      "p",
                      null,
                      "Inicia una conversación escribiendo un mensaje.",
                      -1
                      /* HOISTED */
                    ))
                  ])
                ])) : vue.createCommentVNode("v-if", true),
                (vue.openBlock(true), vue.createElementBlock(
                  vue.Fragment,
                  null,
                  vue.renderList($setup.messages, (message, index) => {
                    return vue.openBlock(), vue.createElementBlock("div", {
                      key: index,
                      class: "mb-5 animate-fade-in"
                    }, [
                      vue.createElementVNode(
                        "div",
                        {
                          class: vue.normalizeClass([
                            "max-w-[90%] p-5 rounded-lg",
                            message.isUser ? "text-white ml-auto" : "bg-gray-100 text-gray-800"
                          ]),
                          style: vue.normalizeStyle(message.isUser ? { backgroundColor: $props.primaryColor } : {})
                        },
                        [
                          vue.createElementVNode("div", _hoisted_19, [
                            vue.createCommentVNode(" Avatar solo para mensajes que no son del usuario "),
                            !message.isUser ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_20, [..._cache[12] || (_cache[12] = [
                              vue.createElementVNode(
                                "i",
                                { class: "pi pi-user text-gray-500" },
                                null,
                                -1
                                /* HOISTED */
                              )
                            ])])) : vue.createCommentVNode("v-if", true),
                            vue.createElementVNode("div", _hoisted_21, [
                              !message.isUser ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_22, "Soporte")) : vue.createCommentVNode("v-if", true),
                              vue.createElementVNode(
                                "div",
                                _hoisted_23,
                                vue.toDisplayString(message.text),
                                1
                                /* TEXT */
                              ),
                              vue.createElementVNode(
                                "div",
                                {
                                  class: vue.normalizeClass(["text-right text-xs mt-2", message.isUser ? "text-white opacity-70" : "text-gray-500"])
                                },
                                vue.toDisplayString(message.timestamp ? $setup.formatMessageTime(message.timestamp) : ""),
                                3
                                /* TEXT, CLASS */
                              )
                            ])
                          ]),
                          vue.createCommentVNode(" Indicador de estado del mensaje "),
                          message.pending ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_24, [..._cache[13] || (_cache[13] = [
                            vue.createElementVNode(
                              "i",
                              { class: "pi pi-spin pi-spinner text-white text-xs" },
                              null,
                              -1
                              /* HOISTED */
                            )
                          ])])) : vue.createCommentVNode("v-if", true),
                          message.error ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_25, " Error al enviar ")) : vue.createCommentVNode("v-if", true)
                        ],
                        6
                        /* CLASS, STYLE */
                      )
                    ]);
                  }),
                  128
                  /* KEYED_FRAGMENT */
                ))
              ])) : (vue.openBlock(), vue.createElementBlock(
                vue.Fragment,
                { key: 1 },
                [
                  vue.createCommentVNode(" FAQ View mejorado "),
                  vue.createElementVNode("div", _hoisted_26, [
                    vue.createElementVNode("div", _hoisted_27, [
                      vue.createElementVNode(
                        "button",
                        {
                          onClick: _cache[3] || (_cache[3] = ($event) => $setup.showChatView = true),
                          class: "text-white rounded-lg py-3 px-5 font-medium focus:outline-none w-full transition-all hover:opacity-90 flex items-center justify-center",
                          style: vue.normalizeStyle({ backgroundColor: $props.primaryColor })
                        },
                        _cache[14] || (_cache[14] = [
                          vue.createElementVNode(
                            "i",
                            { class: "pi pi-comments mr-2" },
                            null,
                            -1
                            /* HOISTED */
                          ),
                          vue.createTextVNode("Iniciar chat ")
                        ]),
                        4
                        /* STYLE */
                      )
                    ]),
                    vue.createElementVNode("h3", _hoisted_28, [
                      vue.createElementVNode(
                        "i",
                        {
                          class: "pi pi-question-circle mr-2",
                          style: vue.normalizeStyle({ color: $props.primaryColor })
                        },
                        null,
                        4
                        /* STYLE */
                      ),
                      _cache[15] || (_cache[15] = vue.createTextVNode(" Preguntas Frecuentes "))
                    ]),
                    $setup.loadingFaqs ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_29, [
                      vue.createElementVNode(
                        "i",
                        {
                          class: "pi pi-spin pi-spinner text-3xl",
                          style: vue.normalizeStyle({ color: $props.primaryColor })
                        },
                        null,
                        4
                        /* STYLE */
                      ),
                      _cache[16] || (_cache[16] = vue.createElementVNode(
                        "p",
                        { class: "text-sm text-gray-500 mt-3" },
                        "Cargando preguntas frecuentes...",
                        -1
                        /* HOISTED */
                      ))
                    ])) : $setup.faqs.length === 0 ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_30, [
                      vue.createElementVNode(
                        "i",
                        {
                          class: "pi pi-inbox text-3xl mb-3",
                          style: vue.normalizeStyle({ color: $props.primaryColor })
                        },
                        null,
                        4
                        /* STYLE */
                      ),
                      _cache[17] || (_cache[17] = vue.createElementVNode(
                        "p",
                        null,
                        "No hay preguntas frecuentes disponibles.",
                        -1
                        /* HOISTED */
                      ))
                    ])) : (vue.openBlock(), vue.createElementBlock("div", _hoisted_31, [
                      (vue.openBlock(true), vue.createElementBlock(
                        vue.Fragment,
                        null,
                        vue.renderList($setup.faqCategories, (category, index) => {
                          return vue.openBlock(), vue.createElementBlock(
                            "div",
                            {
                              key: index,
                              class: "mb-6 animate-fade-in",
                              style: vue.normalizeStyle({ animationDelay: `${index * 0.1}s` })
                            },
                            [
                              vue.createElementVNode(
                                "h4",
                                _hoisted_32,
                                vue.toDisplayString(category),
                                1
                                /* TEXT */
                              ),
                              vue.createElementVNode("div", _hoisted_33, [
                                (vue.openBlock(true), vue.createElementBlock(
                                  vue.Fragment,
                                  null,
                                  vue.renderList($setup.getFaqsByCategory(category), (faq) => {
                                    return vue.openBlock(), vue.createElementBlock("div", {
                                      key: faq.id,
                                      class: "border border-gray-200 rounded-lg overflow-hidden hover:border-gray-300 transition-all"
                                    }, [
                                      vue.createElementVNode("div", {
                                        class: "p-4 bg-white cursor-pointer flex justify-between items-center",
                                        onClick: ($event) => $setup.toggleFaq(faq.id),
                                        style: vue.normalizeStyle($setup.expandedFaqs.includes(faq.id) ? { borderLeft: `4px solid ${$props.primaryColor}` } : {})
                                      }, [
                                        vue.createElementVNode(
                                          "span",
                                          _hoisted_35,
                                          vue.toDisplayString(faq.question),
                                          1
                                          /* TEXT */
                                        ),
                                        vue.createElementVNode(
                                          "i",
                                          {
                                            class: vue.normalizeClass(["pi transition-transform", $setup.expandedFaqs.includes(faq.id) ? "pi-chevron-up" : "pi-chevron-down"]),
                                            style: vue.normalizeStyle({ color: $props.primaryColor })
                                          },
                                          null,
                                          6
                                          /* CLASS, STYLE */
                                        )
                                      ], 12, _hoisted_34),
                                      $setup.expandedFaqs.includes(faq.id) ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_36, [
                                        vue.createElementVNode(
                                          "p",
                                          _hoisted_37,
                                          vue.toDisplayString(faq.answer),
                                          1
                                          /* TEXT */
                                        ),
                                        vue.createElementVNode("div", _hoisted_38, [
                                          vue.createElementVNode("button", {
                                            onClick: ($event) => $setup.setInitialQuestion(faq.question),
                                            class: "text-sm text-white rounded-lg py-2 px-4 focus:outline-none flex items-center transition-all hover:opacity-90",
                                            style: vue.normalizeStyle({ backgroundColor: $props.primaryColor })
                                          }, [..._cache[18] || (_cache[18] = [
                                            vue.createElementVNode(
                                              "i",
                                              { class: "pi pi-comments mr-2" },
                                              null,
                                              -1
                                              /* HOISTED */
                                            ),
                                            vue.createTextVNode(" Consultar más ")
                                          ])], 12, _hoisted_39)
                                        ])
                                      ])) : vue.createCommentVNode("v-if", true)
                                    ]);
                                  }),
                                  128
                                  /* KEYED_FRAGMENT */
                                ))
                              ])
                            ],
                            4
                            /* STYLE */
                          );
                        }),
                        128
                        /* KEYED_FRAGMENT */
                      ))
                    ]))
                  ])
                ],
                2112
                /* STABLE_FRAGMENT, DEV_ROOT_FRAGMENT */
              ))
            ])
          ],
          2112
          /* STABLE_FRAGMENT, DEV_ROOT_FRAGMENT */
        )),
        vue.createCommentVNode(" Chat input (displayed after registration and when in chat view) mejorado "),
        $setup.isRegistered && $setup.showChatView ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_40, [
          vue.createElementVNode("div", _hoisted_41, [
            vue.withDirectives(vue.createElementVNode(
              "input",
              {
                "onUpdate:modelValue": _cache[4] || (_cache[4] = ($event) => $setup.newMessage = $event),
                type: "text",
                placeholder: "Escribe un mensaje...",
                class: "flex-1 border border-gray-300 rounded-lg px-5 py-4 focus:outline-none focus:border-2 transition-all",
                style: vue.normalizeStyle({ "--tw-border-opacity": 1, borderColor: $props.primaryColor }),
                onKeyup: vue.withKeys($setup.sendMessage, ["enter"])
              },
              null,
              36
              /* STYLE, NEED_HYDRATION */
            ), [
              [vue.vModelText, $setup.newMessage]
            ]),
            vue.createElementVNode(
              "button",
              {
                onClick: $setup.sendMessage,
                class: "ml-3 text-white rounded-lg p-4 focus:outline-none transition-all hover:opacity-90",
                style: vue.normalizeStyle({ backgroundColor: $props.primaryColor, border: "2px solid #000" })
              },
              _cache[19] || (_cache[19] = [
                vue.createElementVNode(
                  "i",
                  { class: "pi pi-send" },
                  null,
                  -1
                  /* HOISTED */
                )
              ]),
              4
              /* STYLE */
            )
          ])
        ])) : vue.createCommentVNode("v-if", true),
        vue.createCommentVNode(" Bottom navigation when in FAQ view mejorado "),
        $setup.isRegistered && !$setup.showChatView ? (vue.openBlock(), vue.createElementBlock("div", _hoisted_42, [
          vue.createElementVNode("div", _hoisted_43, [
            vue.createElementVNode("button", {
              onClick: $setup.refreshFaqs,
              class: "text-sm text-gray-600 focus:outline-none flex items-center hover:text-gray-800 transition-colors"
            }, _cache[20] || (_cache[20] = [
              vue.createElementVNode(
                "i",
                { class: "pi pi-refresh mr-1" },
                null,
                -1
                /* HOISTED */
              ),
              vue.createTextVNode(" Actualizar FAQs ")
            ])),
            vue.createElementVNode(
              "button",
              {
                onClick: _cache[5] || (_cache[5] = ($event) => $setup.showChatView = true),
                class: "text-white rounded-lg py-2 px-4 text-sm focus:outline-none transition-all hover:opacity-90",
                style: vue.normalizeStyle({ backgroundColor: $props.primaryColor, border: "2px solid #000" })
              },
              _cache[21] || (_cache[21] = [
                vue.createElementVNode(
                  "i",
                  { class: "pi pi-comments mr-2" },
                  null,
                  -1
                  /* HOISTED */
                ),
                vue.createTextVNode(" Ir al chat ")
              ]),
              4
              /* STYLE */
            )
          ])
        ])) : vue.createCommentVNode("v-if", true)
      ])) : vue.createCommentVNode("v-if", true)
    ]);
  }
  const ChatWidget = /* @__PURE__ */ _export_sfc(_sfc_main, [["render", _sfc_render], ["__file", "/app/src/components/ChatWidget.vue"]]);
  console.log("GrowDesk Widget - Iniciando...");
  const createWidgetContainer = () => {
    console.log("Creando contenedor para el widget");
    const container = document.createElement("div");
    container.id = "growdesk-widget-container";
    document.body.appendChild(container);
    return container;
  };
  const initializeWidget = () => {
    var _a2, _b, _c, _d, _e, _f;
    console.log("Inicializando widget");
    initializeFromScript();
    if ((_a2 = window.GrowDeskConfig) == null ? void 0 : _a2.apiUrl) {
      console.log("Configurando API URL desde GrowDeskConfig:", window.GrowDeskConfig.apiUrl);
      let apiUrl = window.GrowDeskConfig.apiUrl;
      configureWidgetApi({
        apiUrl
      });
    }
    const pinia = createPinia();
    const container = createWidgetContainer();
    const app = vue.createApp(ChatWidget, {
      // Propiedades que se pueden personalizar
      brandName: ((_b = window.GrowDeskConfig) == null ? void 0 : _b.brandName) || "GrowDesk",
      welcomeMessage: ((_c = window.GrowDeskConfig) == null ? void 0 : _c.welcomeMessage) || "¿Necesitas ayuda? Estamos aquí para ayudarte.",
      primaryColor: ((_d = window.GrowDeskConfig) == null ? void 0 : _d.primaryColor) || "#1976d2",
      position: ((_e = window.GrowDeskConfig) == null ? void 0 : _e.position) || "bottom-right",
      logoUrl: ((_f = window.GrowDeskConfig) == null ? void 0 : _f.logoUrl) || ""
    });
    app.use(pinia);
    app.mount(container);
    console.log("Widget montado correctamente");
    const widgetApi = {
      open: () => {
        console.log("Abriendo widget vía API");
        const widgetComponent = document.querySelector("#growdesk-widget-container");
        if (widgetComponent) {
          const event = new CustomEvent("open-widget");
          widgetComponent.dispatchEvent(event);
        }
      },
      close: () => {
        console.log("Cerrando widget vía API");
        const widgetComponent = document.querySelector("#growdesk-widget-container");
        if (widgetComponent) {
          const event = new CustomEvent("close-widget");
          widgetComponent.dispatchEvent(event);
        }
      },
      configure: (config) => {
        console.log("Reconfigurando widget", config);
        if (config.apiUrl || config.widgetId || config.widgetToken) {
          configureWidgetApi({
            apiUrl: config.apiUrl,
            widgetId: config.widgetId,
            widgetToken: config.widgetToken
          });
        }
        window.GrowDeskConfig = {
          ...window.GrowDeskConfig,
          ...config
        };
      }
    };
    console.log("API del widget creada", widgetApi);
    return widgetApi;
  };
  let widgetInitialized = false;
  const initialize = () => {
    if (widgetInitialized) return;
    console.log("Inicializando GrowDeskWidget automáticamente");
    const api2 = initializeWidget();
    window.GrowDeskWidget = api2;
    widgetInitialized = true;
  };
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initialize);
  } else {
    initialize();
  }
  return initializeWidget;
});
//# sourceMappingURL=growdesk-widget.umd.js.map
