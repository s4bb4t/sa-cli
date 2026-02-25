package scaffold

func swaggerHTMLTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Swagger</title>
    <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'><rect width='32' height='32' rx='8' fill='%236366f1'/><text x='16' y='22' font-size='16' font-weight='800' fill='white' text-anchor='middle' font-family='sans-serif'>AP</text></svg>">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui.css">
    <style>
        :root {
            --bg-primary: #0a0e17;
            --bg-secondary: #111827;
            --bg-card: #1a2233;
            --bg-input: #0f1623;
            --border: #1e2d45;
            --text-primary: #e2e8f0;
            --text-secondary: #94a3b8;
            --accent: #6366f1;
            --accent-hover: #818cf8;
            --accent-glow: rgba(99, 102, 241, 0.15);
            --green: #22c55e;
            --blue: #3b82f6;
            --orange: #f59e0b;
            --red: #ef4444;
        }

        * { box-sizing: border-box; margin: 0; padding: 0; }

        body {
            background: var(--bg-primary);
            color: var(--text-primary);
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
        }

        /* ── Header ── */
        .topbar-wrapper, .swagger-ui .topbar { display: none !important; }

        .web3-header {
            background: var(--bg-secondary);
            border-bottom: 1px solid var(--border);
            padding: 20px 40px;
            display: flex;
            align-items: center;
            gap: 16px;
            position: sticky;
            top: 0;
            z-index: 100;
            backdrop-filter: blur(12px);
        }

        .web3-header .logo {
            width: 36px;
            height: 36px;
            background: linear-gradient(135deg, var(--accent), #a855f7);
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 800;
            font-size: 16px;
            color: #fff;
            flex-shrink: 0;
        }

        .web3-header h1 {
            font-size: 20px;
            font-weight: 700;
            letter-spacing: -0.02em;
            color: var(--text-primary);
        }

        .web3-header h1 span {
            background: linear-gradient(135deg, var(--accent), #a855f7);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }

        /* Version selector chip */
        .version-wrapper {
            position: relative;
            display: flex;
            align-items: center;
        }

        .version-badge {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            font-size: 13px;
            font-weight: 600;
            color: var(--text-primary);
            background: var(--bg-input);
            border: 1px solid rgba(99, 102, 241, 0.25);
            padding: 6px 14px;
            border-radius: 20px;
            cursor: pointer;
            outline: none;
            appearance: none;
            -webkit-appearance: none;
            -moz-appearance: none;
            transition: all 0.15s ease;
        }

        .version-badge::before {
            content: '';
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: var(--green);
            flex-shrink: 0;
        }

        .version-badge option {
            background: var(--bg-secondary);
            color: var(--text-primary);
        }

        .version-badge:hover,
        .version-badge:focus {
            border-color: var(--accent);
            box-shadow: 0 0 0 3px var(--accent-glow);
        }

        /* Header actions */
        .header-actions {
            margin-left: auto;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .btn-copy-spec {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            font-size: 12px;
            font-weight: 600;
            color: var(--text-secondary);
            background: var(--bg-input);
            border: 1px solid var(--border);
            padding: 6px 12px;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.15s ease;
        }

        .btn-copy-spec:hover {
            color: var(--text-primary);
            border-color: var(--accent);
        }

        .btn-copy-spec.copied {
            color: var(--green);
            border-color: var(--green);
        }

        .kbd-hint {
            font-size: 11px;
            color: var(--text-secondary);
            opacity: 0.6;
        }

        .kbd-hint kbd {
            background: var(--bg-input);
            border: 1px solid var(--border);
            border-radius: 4px;
            padding: 1px 5px;
            font-family: inherit;
            font-size: 11px;
        }

        /* ── Loading spinner ── */
        .loading-overlay {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 80px 20px;
            gap: 16px;
        }

        .loading-overlay .spinner {
            width: 36px;
            height: 36px;
            border: 3px solid var(--border);
            border-top-color: var(--accent);
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        .loading-overlay .loading-text {
            color: var(--text-secondary);
            font-size: 14px;
        }

        /* ── Error state ── */
        .error-state {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 80px 20px;
            gap: 16px;
            text-align: center;
        }

        .error-state .error-icon {
            width: 48px;
            height: 48px;
            border-radius: 50%;
            background: rgba(239, 68, 68, 0.1);
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 24px;
        }

        .error-state .error-title {
            font-size: 16px;
            font-weight: 600;
            color: var(--text-primary);
        }

        .error-state .error-detail {
            font-size: 13px;
            color: var(--text-secondary);
            max-width: 400px;
        }

        .error-state .btn-retry {
            margin-top: 8px;
            font-size: 13px;
            font-weight: 600;
            color: #fff;
            background: var(--accent);
            border: none;
            padding: 8px 20px;
            border-radius: 8px;
            cursor: pointer;
            transition: background 0.15s ease;
        }

        .error-state .btn-retry:hover {
            background: var(--accent-hover);
        }

        /* ── Footer ── */
        .web3-footer {
            border-top: 1px solid var(--border);
            padding: 16px 40px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            font-size: 12px;
            color: var(--text-secondary);
        }

        .web3-footer a {
            color: var(--accent);
            text-decoration: none;
        }

        .web3-footer a:hover {
            text-decoration: underline;
        }

        /* ── Global overrides ── */
        .swagger-ui {
            font-family: inherit;
            color: var(--text-primary);
        }

        .swagger-ui .wrapper {
            max-width: 1200px;
            padding: 24px 40px;
        }

        /* ── Force light text everywhere ── */
        .swagger-ui,
        .swagger-ui p,
        .swagger-ui span,
        .swagger-ui label,
        .swagger-ui small,
        .swagger-ui td,
        .swagger-ui th,
        .swagger-ui li,
        .swagger-ui h1, .swagger-ui h2, .swagger-ui h3,
        .swagger-ui h4, .swagger-ui h5, .swagger-ui h6 {
            color: var(--text-primary) !important;
        }

        /* ── Info block ── */
        .swagger-ui .info {
            margin: 16px 0 32px;
        }

        .swagger-ui .info .title {
            font-size: 24px !important;
            font-weight: 700 !important;
        }

        .swagger-ui .info p, .swagger-ui .info li,
        .swagger-ui .info table td, .swagger-ui .info table th {
            color: var(--text-secondary) !important;
            font-size: 14px;
        }

        .swagger-ui .info a {
            color: var(--accent) !important;
        }

        .swagger-ui .info .base-url {
            color: var(--text-secondary) !important;
            font-size: 13px;
        }

        .swagger-ui .scheme-container {
            background: var(--bg-secondary) !important;
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 16px !important;
            box-shadow: none !important;
        }

        .swagger-ui .scheme-container label,
        .swagger-ui .scheme-container select {
            color: var(--text-primary) !important;
        }

        /* ── Tag sections ── */
        .swagger-ui .opblock-tag-section { margin-bottom: 8px; }

        .swagger-ui .opblock-tag {
            border-bottom: 1px solid var(--border) !important;
            font-size: 16px;
            font-weight: 600;
            padding: 14px 0;
            transition: opacity 0.15s ease;
        }

        .swagger-ui .opblock-tag:hover { background: transparent !important; }

        .swagger-ui .opblock-tag svg,
        .swagger-ui .opblock-tag small { fill: var(--text-secondary) !important; color: var(--text-secondary) !important; }

        /* ── Operation blocks ── */
        .swagger-ui .opblock {
            background: var(--bg-card) !important;
            border: 1px solid var(--border) !important;
            border-radius: 10px !important;
            margin-bottom: 8px !important;
            box-shadow: none !important;
            transition: box-shadow 0.2s ease, border-color 0.2s ease;
        }

        .swagger-ui .opblock:hover {
            box-shadow: 0 2px 12px rgba(0,0,0,0.15) !important;
        }

        .swagger-ui .opblock .opblock-summary {
            border: none !important;
            padding: 12px 16px !important;
        }

        .swagger-ui .opblock .opblock-summary-method {
            border-radius: 6px !important;
            font-size: 12px !important;
            font-weight: 700 !important;
            min-width: 70px !important;
            padding: 6px 12px !important;
            text-align: center;
            color: #fff !important;
        }

        .swagger-ui .opblock.opblock-get .opblock-summary-method { background: var(--blue) !important; }
        .swagger-ui .opblock.opblock-post .opblock-summary-method { background: var(--green) !important; }
        .swagger-ui .opblock.opblock-put .opblock-summary-method { background: var(--orange) !important; }
        .swagger-ui .opblock.opblock-delete .opblock-summary-method { background: var(--red) !important; }
        .swagger-ui .opblock.opblock-patch .opblock-summary-method { background: #a855f7 !important; }

        .swagger-ui .opblock.opblock-get { border-left: 3px solid var(--blue) !important; background: rgba(59, 130, 246, 0.04) !important; }
        .swagger-ui .opblock.opblock-post { border-left: 3px solid var(--green) !important; background: rgba(34, 197, 94, 0.04) !important; }
        .swagger-ui .opblock.opblock-put { border-left: 3px solid var(--orange) !important; background: rgba(245, 158, 11, 0.04) !important; }
        .swagger-ui .opblock.opblock-delete { border-left: 3px solid var(--red) !important; background: rgba(239, 68, 68, 0.04) !important; }
        .swagger-ui .opblock.opblock-patch { border-left: 3px solid #a855f7 !important; background: rgba(168, 85, 247, 0.04) !important; }

        .swagger-ui .opblock .opblock-summary-path,
        .swagger-ui .opblock .opblock-summary-path a {
            font-weight: 500;
            font-family: 'JetBrains Mono', 'Fira Code', monospace;
            font-size: 13px;
        }

        .swagger-ui .opblock .opblock-summary-description {
            color: var(--text-secondary) !important;
            font-size: 13px;
        }

        /* ── Expanded operation ── */
        .swagger-ui .opblock-body {
            background: var(--bg-secondary) !important;
            border-top: 1px solid var(--border) !important;
        }

        .swagger-ui .opblock-body pre,
        .swagger-ui .opblock-body pre.microlight {
            background: var(--bg-input) !important;
            color: var(--text-primary) !important;
            border: 1px solid var(--border) !important;
            border-radius: 8px !important;
            font-family: 'JetBrains Mono', 'Fira Code', monospace;
            font-size: 12px !important;
            padding: 16px !important;
        }

        .swagger-ui .opblock-section-header {
            background: transparent !important;
            border-bottom: 1px solid var(--border) !important;
            box-shadow: none !important;
        }

        .swagger-ui .tab li { color: var(--text-secondary) !important; }
        .swagger-ui .tab li.active { color: var(--text-primary) !important; }

        /* ── Live server response ── */
        .swagger-ui .live-responses-table {
            background: rgba(34, 197, 94, 0.06) !important;
            border: 1px solid rgba(34, 197, 94, 0.25) !important;
            border-radius: 10px !important;
            padding: 12px !important;
            margin-top: 8px !important;
            position: relative;
        }

        .swagger-ui .live-responses-table::before {
            content: 'LIVE';
            position: absolute;
            top: -9px;
            left: 14px;
            background: var(--green);
            color: #000;
            font-size: 10px;
            font-weight: 800;
            letter-spacing: 0.08em;
            padding: 1px 8px;
            border-radius: 4px;
        }

        .swagger-ui .responses-inner .responses-table.live-responses-table .response-col_status {
            color: var(--green) !important;
            font-weight: 700 !important;
        }

        .swagger-ui .live-responses-table pre {
            border-color: rgba(34, 197, 94, 0.2) !important;
            background: var(--bg-primary) !important;
        }

        .swagger-ui .request-url pre {
            background: var(--bg-primary) !important;
            border: 1px solid rgba(34, 197, 94, 0.2) !important;
            border-radius: 8px !important;
            color: var(--green) !important;
            font-family: 'JetBrains Mono', 'Fira Code', monospace;
            font-size: 12px !important;
        }

        .swagger-ui .request-url h4 {
            color: var(--green) !important;
            font-size: 12px !important;
        }

        /* ── Request duration badge ── */
        .swagger-ui .response-col_duration {
            font-family: 'JetBrains Mono', 'Fira Code', monospace;
            font-size: 12px;
            color: var(--text-secondary) !important;
        }

        /* ── Filter bar ── */
        .swagger-ui .filter-container {
            background: var(--bg-secondary) !important;
            border-bottom: 1px solid var(--border);
            padding: 12px 0 !important;
        }

        .swagger-ui .filter-container .operation-filter-input {
            background: var(--bg-input) !important;
            border: 1px solid var(--border) !important;
            border-radius: 8px !important;
            color: var(--text-primary) !important;
            padding: 8px 14px !important;
        }

        .swagger-ui .filter-container .operation-filter-input:focus {
            border-color: var(--accent) !important;
            box-shadow: 0 0 0 3px var(--accent-glow) !important;
        }

        /* ── Parameters & Tables ── */
        .swagger-ui table thead tr td, .swagger-ui table thead tr th {
            color: var(--text-secondary) !important;
            border-bottom: 1px solid var(--border) !important;
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .swagger-ui table tbody tr td {
            border-bottom: 1px solid var(--border) !important;
        }

        .swagger-ui .parameter__name.required::after {
            color: var(--red) !important;
        }

        .swagger-ui .parameter__deprecated { color: var(--red) !important; }

        /* ── Response codes ── */
        .swagger-ui .responses-table .response-col_status {
            font-weight: 600;
        }

        .swagger-ui .responses-table .response-col_description p {
            color: var(--text-secondary) !important;
        }

        .swagger-ui .response-col_links { color: var(--text-secondary) !important; }

        /* ── Models ── */
        .swagger-ui section.models {
            border: 1px solid var(--border) !important;
            border-radius: 12px !important;
            background: var(--bg-secondary) !important;
        }

        .swagger-ui section.models .model-container {
            background: var(--bg-card) !important;
            border-radius: 8px;
            margin: 4px 0 !important;
        }

        .swagger-ui .model .property { color: var(--text-secondary) !important; }
        .swagger-ui .model .property.primitive { color: var(--accent) !important; }
        .swagger-ui .model-toggle::after { background: none !important; }

        /* ── All selects, inputs, textareas ── */
        .swagger-ui input[type=text],
        .swagger-ui input[type=password],
        .swagger-ui input[type=search],
        .swagger-ui input[type=email],
        .swagger-ui input[type=file],
        .swagger-ui textarea,
        .swagger-ui select {
            background: var(--bg-input) !important;
            border: 1px solid var(--border) !important;
            border-radius: 8px !important;
            color: var(--text-primary) !important;
            padding: 8px 12px !important;
        }

        .swagger-ui select {
            background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='8'%3E%3Cpath d='M1 1l5 5 5-5' stroke='%23e2e8f0' stroke-width='1.5' fill='none'/%3E%3C/svg%3E") !important;
            background-repeat: no-repeat !important;
            background-position: right 10px center !important;
            padding-right: 28px !important;
            appearance: none !important;
            -webkit-appearance: none !important;
            -moz-appearance: none !important;
        }

        .swagger-ui select option {
            background: var(--bg-secondary) !important;
            color: var(--text-primary) !important;
        }

        .swagger-ui input:focus, .swagger-ui textarea:focus, .swagger-ui select:focus {
            border-color: var(--accent) !important;
            box-shadow: 0 0 0 3px var(--accent-glow) !important;
            outline: none !important;
        }

        /* ── Buttons ── */
        .swagger-ui .btn {
            border-radius: 8px !important;
            font-weight: 600 !important;
            font-size: 13px !important;
            transition: all 0.15s ease !important;
            color: var(--text-primary) !important;
        }

        .swagger-ui .btn.execute {
            background: var(--accent) !important;
            border-color: var(--accent) !important;
            color: #fff !important;
        }

        .swagger-ui .btn.execute:hover {
            background: var(--accent-hover) !important;
            box-shadow: 0 0 20px var(--accent-glow) !important;
        }

        .swagger-ui .btn.cancel {
            border-color: var(--border) !important;
            color: var(--text-secondary) !important;
        }

        .swagger-ui .btn.authorize {
            color: var(--accent) !important;
            border-color: var(--accent) !important;
        }

        .swagger-ui .btn.authorize svg { fill: var(--accent) !important; }

        /* ── Auth modal ── */
        .swagger-ui .dialog-ux .modal-ux {
            background: var(--bg-secondary) !important;
            border: 1px solid var(--border) !important;
            border-radius: 16px !important;
        }

        .swagger-ui .dialog-ux .modal-ux-header {
            border-bottom: 1px solid var(--border) !important;
        }

        .swagger-ui .dialog-ux .modal-ux-content p,
        .swagger-ui .dialog-ux .modal-ux-content label {
            color: var(--text-secondary) !important;
        }

        .swagger-ui .dialog-ux .modal-ux-content .btn-done {
            color: #fff !important;
        }

        .swagger-ui .dialog-ux .backdrop-ux {
            background: rgba(0, 0, 0, 0.7) !important;
        }

        /* ── Copy button ── */
        .swagger-ui .copy-to-clipboard { background: var(--bg-card) !important; }
        .swagger-ui .copy-to-clipboard button { color: var(--text-secondary) !important; }

        /* ── Scrollbar ── */
        ::-webkit-scrollbar { width: 6px; height: 6px; }
        ::-webkit-scrollbar-track { background: var(--bg-primary); }
        ::-webkit-scrollbar-thumb {
            background: var(--border);
            border-radius: 3px;
        }
        ::-webkit-scrollbar-thumb:hover { background: #2d3f5c; }

        /* ── SVG icons ── */
        .swagger-ui svg:not(.model-toggle) { fill: var(--text-secondary) !important; }
        .swagger-ui .expand-operation svg { fill: var(--text-secondary) !important; }
        .swagger-ui .arrow { fill: var(--text-primary) !important; }
        .swagger-ui .locked svg { fill: var(--green) !important; }
        .swagger-ui .unlocked svg { fill: var(--red) !important; }

        /* ── Loading ── */
        .swagger-ui .loading-container .loading::after { color: var(--text-secondary); }

        /* ── Misc text catch-all ── */
        .swagger-ui .opblock-description-wrapper p,
        .swagger-ui .opblock-external-docs-wrapper p {
            color: var(--text-secondary) !important;
        }

        .swagger-ui .markdown p, .swagger-ui .markdown li {
            color: var(--text-secondary) !important;
        }

        .swagger-ui .response-control-media-type__accept-message {
            color: var(--green) !important;
        }

        .swagger-ui .response-content-type.controls-accept-header select {
            border-color: var(--green) !important;
        }

        /* ── Highlighted JSON values ── */
        .swagger-ui .microlight .headerline,
        .swagger-ui .microlight .attrname,
        .swagger-ui .microlight .keyword { color: var(--accent) !important; }
        .swagger-ui .microlight .string { color: var(--green) !important; }
        .swagger-ui .microlight .number { color: var(--orange) !important; }

        /* ── Smooth transitions ── */
        .swagger-ui .opblock-body,
        .swagger-ui .model-box,
        .swagger-ui section.models {
            transition: max-height 0.25s ease, opacity 0.2s ease;
        }

        /* ── Accessibility focus rings ── */
        .swagger-ui .opblock-summary:focus-visible,
        .swagger-ui .opblock-tag:focus-visible,
        .swagger-ui .btn:focus-visible {
            outline: 2px solid var(--accent) !important;
            outline-offset: 2px !important;
        }

        /* ── Reduced motion ── */
        @media (prefers-reduced-motion: reduce) {
            *, *::before, *::after {
                animation-duration: 0.01ms !important;
                animation-iteration-count: 1 !important;
                transition-duration: 0.01ms !important;
            }
        }

        /* ── Mobile responsiveness ── */
        @media (max-width: 768px) {
            .web3-header {
                padding: 14px 16px;
                flex-wrap: wrap;
                gap: 10px;
            }
            .web3-header h1 { font-size: 16px; }
            .swagger-ui .wrapper { padding: 16px; }
            .web3-footer { padding: 12px 16px; flex-direction: column; gap: 8px; }
            .kbd-hint { display: none; }
        }
    </style>
</head>
<body>
    <header class="web3-header">
        <div class="logo">AP</div>
        <h1><span>API</span> Docs</h1>
        <div class="version-wrapper">
            <select class="version-badge" id="version-select"></select>
        </div>
        <div class="header-actions">
            <button class="btn-copy-spec" id="btn-copy-spec" title="Copy spec URL to clipboard">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
                Copy spec URL
            </button>
            <span class="kbd-hint"><kbd>Ctrl</kbd>+<kbd>K</kbd> search</span>
        </div>
    </header>
    <div id="swagger-ui">
        <div class="loading-overlay" id="loading-state">
            <div class="spinner"></div>
            <div class="loading-text">Loading API specification&hellip;</div>
        </div>
    </div>
    <footer class="web3-footer" id="web3-footer" style="display:none;">
        <div>API &mdash; <span id="footer-version"></span></div>
        <div>Built with <a href="https://swagger.io/tools/swagger-ui/" target="_blank" rel="noopener">Swagger UI</a></div>
    </footer>
    <script src="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui-bundle.js"></script>
    <script>
        var ui = null;
        var currentVersion = '';

        function showLoading() {
            var el = document.getElementById('swagger-ui');
            el.innerHTML = '<div class="loading-overlay"><div class="spinner"></div><div class="loading-text">Loading API specification&hellip;</div></div>';
        }

        function showError(message) {
            var el = document.getElementById('swagger-ui');
            el.innerHTML =
                '<div class="error-state">' +
                    '<div class="error-icon">&#x26A0;</div>' +
                    '<div class="error-title">Failed to load API spec</div>' +
                    '<div class="error-detail">' + escapeHtml(message) + '</div>' +
                    '<button class="btn-retry" onclick="initSwagger(currentVersion)">Retry</button>' +
                '</div>';
        }

        function escapeHtml(str) {
            var d = document.createElement('div');
            d.textContent = str;
            return d.innerHTML;
        }

        function initSwagger(version) {
            currentVersion = version;

            if (ui && typeof ui.destroy === 'function') {
                ui.destroy();
                ui = null;
            }

            showLoading();
            window.scrollTo({ top: 0, behavior: 'smooth' });

            var specUrl = '/docs/' + version + '/openapi.json';

            try {
                ui = SwaggerUIBundle({
                    url: specUrl,
                    dom_id: '#swagger-ui',
                    deepLinking: true,
                    presets: [SwaggerUIBundle.presets.apis],
                    defaultModelsExpandDepth: 1,
                    docExpansion: 'list',
                    syntaxHighlight: { theme: 'monokai' },
                    persistAuthorization: true,
                    filter: true,
                    tryItOutEnabled: true,
                    requestSnippetsEnabled: false,
                    displayRequestDuration: true,
                    responseInterceptor: function(response) {
                        if (response.status >= 400 && response.url && response.url.indexOf('openapi.json') !== -1) {
                            showError('Spec returned HTTP ' + response.status);
                        }
                        return response;
                    }
                });

                document.getElementById('footer-version').textContent = version;
                document.getElementById('web3-footer').style.display = '';
            } catch (e) {
                showError(e.message || 'Unknown error initializing Swagger UI');
            }
        }

        function fetchWithRetry(url, retries, delay) {
            return fetch(url).then(function(r) {
                if (!r.ok) {
                    if (retries > 0) {
                        return new Promise(function(resolve) {
                            setTimeout(resolve, delay);
                        }).then(function() {
                            return fetchWithRetry(url, retries - 1, delay * 2);
                        });
                    }
                    throw new Error('HTTP ' + r.status + ' fetching ' + url);
                }
                return r.json();
            }).catch(function(err) {
                if (retries > 0) {
                    return new Promise(function(resolve) {
                        setTimeout(resolve, delay);
                    }).then(function() {
                        return fetchWithRetry(url, retries - 1, delay * 2);
                    });
                }
                throw err;
            });
        }

        fetchWithRetry('/docs/versions.json', 2, 1000)
            .then(function(versions) {
                var sel = document.getElementById('version-select');
                versions.forEach(function(v) {
                    var opt = document.createElement('option');
                    opt.value = v;
                    opt.textContent = v;
                    sel.appendChild(opt);
                });

                sel.addEventListener('change', function() { initSwagger(sel.value); });
                initSwagger(versions[0]);
            })
            .catch(function(err) {
                showError('Failed to load API versions: ' + (err.message || err));
            });

        // Copy spec URL button
        document.getElementById('btn-copy-spec').addEventListener('click', function() {
            var btn = this;
            var url = window.location.origin + '/docs/' + currentVersion + '/openapi.json';
            navigator.clipboard.writeText(url).then(function() {
                btn.classList.add('copied');
                btn.textContent = 'Copied!';
                setTimeout(function() {
                    btn.classList.remove('copied');
                    btn.innerHTML = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg> Copy spec URL';
                }, 2000);
            });
        });

        // Keyboard shortcut: Ctrl+K focuses filter input
        document.addEventListener('keydown', function(e) {
            if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
                e.preventDefault();
                var input = document.querySelector('.swagger-ui .operation-filter-input');
                if (input) input.focus();
            }
        });
    </script>
</body>
</html>`
}
