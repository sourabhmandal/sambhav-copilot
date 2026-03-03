(function () {
  "use strict";

  /* ================================
     CONFIG
  =================================== */

  const CONFIG = {
    API_URL: "http://localhost:8080/translate",
    COMPANY_ID: 1,
    SOURCE_LANGUAGE: "en",
    TARGET_LANGUAGE: "bn",
    BATCH_SIZE: 50,
    DEBOUNCE_TIME: 150,
  };

  /* ================================
     STATE
  =================================== */

  let currentLanguage = CONFIG.TARGET_LANGUAGE;
  let isTranslating = false;
  let mutationObserver = null;
  let pendingNodes = new Set();
  let processedNodes = new WeakSet();
  let currentLocation = null;
  let routeVersion = 0;

  /* ================================
     UTILITIES
  =================================== */

  function normalize(text) {
    return text.trim().replace(/\s+/g, " ");
  }

  function isValidText(text) {
    if (!text) return false;
    const trimmed = text.trim();
    if (!trimmed) return false;
    if (trimmed.length < 2) return false;
    if (/^[\d\s.,$%€₹£¥-]+$/.test(trimmed)) return false;
    return true;
  }

  function shouldIgnore(node) {
    if (!node.parentElement) return true;

    const parent = node.parentElement;

    if (
      parent.closest("[data-i18n-ignore]") ||
      parent.closest("[data-i18n-owner='backend']") ||
      parent.closest("#__autonomous_i18n_ui")
    )
      return true;

    const tag = parent.tagName.toLowerCase();

    const blocked = [
      "script",
      "style",
      "noscript",
      "textarea",
      "input",
      "code",
      "pre",
      "svg",
    ];

    if (blocked.includes(tag)) return true;

    if (parent.isContentEditable) return true;

    return false;
  }

  /* ================================
     TEXT EXTRACTION
  =================================== */

  function collectTextNodes(root = document.body) {
    const walker = document.createTreeWalker(
      root,
      NodeFilter.SHOW_TEXT,
      null,
      false
    );

    const nodes = [];

    while (walker.nextNode()) {
      const node = walker.currentNode;
      if (processedNodes.has(node)) continue;
      if (shouldIgnore(node)) continue;

      const text = normalize(node.nodeValue);
      if (!isValidText(text)) continue;

      nodes.push(node);
    }

    return nodes;
  }

  /* ================================
     API CALL
  =================================== */

  async function translateBatch(textArray) {
    if (!textArray.length) return [];

    const payload = {
      company_id: CONFIG.COMPANY_ID,
      location: currentLocation,
      source_language: CONFIG.SOURCE_LANGUAGE,
      target_language: currentLanguage,
      text: textArray,
    };

    const response = await fetch(CONFIG.API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!response.ok) throw new Error("Translation failed");

    return (await response.json()).text || [];
  }

  /* ================================
     HYDRATION
  =================================== */

  function hydrateNodes(results, originalMap) {
    results.forEach((item) => {
      const original = normalize(item.original);
      const node = originalMap.get(original);
      if (!node) return;

      if (!item.translated) return;

      node.nodeValue = item.translated;
      processedNodes.add(node);
    });
  }

  /* ================================
     MAIN TRANSLATION FLOW
  =================================== */

  async function processNodes(nodes) {
    if (!nodes.length) return;

    isTranslating = true;

    const batches = [];
    for (let i = 0; i < nodes.length; i += CONFIG.BATCH_SIZE) {
      batches.push(nodes.slice(i, i + CONFIG.BATCH_SIZE));
    }

    for (const batch of batches) {
      const textArray = batch.map((n) => normalize(n.nodeValue));
      const originalMap = new Map();

      batch.forEach((node) => {
        originalMap.set(normalize(node.nodeValue), node);
      });

      try {
        const results = await translateBatch(textArray);
        hydrateNodes(results, originalMap);
      } catch (err) {
        console.error("Translation error:", err);
      }
    }

    isTranslating = false;
  }

  /* ================================
     INITIAL RUN
  =================================== */

  function initialScan() {
    const nodes = collectTextNodes();
    processNodes(nodes);
  }

  /* ================================
     MUTATION OBSERVER
  =================================== */

  function startObserver() {
    mutationObserver = new MutationObserver((mutations) => {
      if (isTranslating) return;

      mutations.forEach((mutation) => {
        mutation.addedNodes.forEach((node) => {
          if (node.nodeType === 3) {
            pendingNodes.add(node);
          } else if (node.nodeType === 1) {
            collectTextNodes(node).forEach((n) => pendingNodes.add(n));
          }
        });
      });

      debounceProcess();
    });

    mutationObserver.observe(document.body, {
      childList: true,
      subtree: true,
    });
  }

  let debounceTimer = null;

  function debounceProcess() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      const nodes = Array.from(pendingNodes);
      pendingNodes.clear();
      processNodes(nodes);
    }, CONFIG.DEBOUNCE_TIME);
  }

  /* ================================
     ROUTE CHANGE DETECTION (SPA)
  =================================== */

  function interceptHistory() {
    const pushState = history.pushState;
    history.pushState = function () {
      pushState.apply(history, arguments);
      routeChanged();
    };

    const replaceState = history.replaceState;
    history.replaceState = function () {
      replaceState.apply(history, arguments);
      routeChanged();
    };

    window.addEventListener("popstate", routeChanged);
  }

  function routeChanged() {
    routeVersion++;
    setTimeout(initialScan, 200);
  }

  /* ================================
     LOCATION HANDLING
  =================================== */

  function requestLocation() {
    if (!navigator.geolocation) {
      alert("Geolocation not supported");
      return;
    }

    navigator.geolocation.getCurrentPosition((position) => {
      currentLocation = {
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
        accuracy: position.coords.accuracy,
        altitude: position.coords.altitude,
        timestamp: new Date().toISOString(),
        metadata: {
          device: "Browser",
          provider: "GPS",
        },
      };

      retranslateAll();
    });
  }

  /* ================================
     LANGUAGE SWITCH
  =================================== */

  function retranslateAll() {
    processedNodes = new WeakSet();
    initialScan();
  }

  /* ================================
     UI INJECTION
  =================================== */

  function injectUI() {
    const container = document.createElement("div");
    container.id = "__autonomous_i18n_ui";
    container.style.position = "fixed";
    container.style.top = "0";
    container.style.left = "0";
    container.style.width = "100%";
    container.style.background = "#111";
    container.style.color = "#fff";
    container.style.padding = "8px 12px";
    container.style.display = "flex";
    container.style.justifyContent = "space-between";
    container.style.alignItems = "center";
    container.style.zIndex = "999999";
    container.style.fontFamily = "sans-serif";

    const select = document.createElement("select");
    ["en", "bn", "hi", "fr", "de"].forEach((lang) => {
      const opt = document.createElement("option");
      opt.value = lang;
      opt.textContent = lang.toUpperCase();
      if (lang === currentLanguage) opt.selected = true;
      select.appendChild(opt);
    });

    select.onchange = function () {
      currentLanguage = this.value;
      retranslateAll();
    };

    const locationBtn = document.createElement("button");
    locationBtn.textContent = "📍 Share Location";
    locationBtn.style.marginLeft = "10px";
    locationBtn.onclick = requestLocation;

    container.appendChild(select);
    container.appendChild(locationBtn);

    document.body.appendChild(container);

    document.body.style.marginTop = "50px";
  }

  /* ================================
     START ENGINE
  =================================== */

  function start() {
    if (document.readyState === "complete") {
      boot();
    } else {
      window.addEventListener("load", boot);
    }
  }

  function boot() {
    injectUI();
    initialScan();
    startObserver();
    interceptHistory();
  }

  start();
})();