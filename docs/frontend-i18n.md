## CyberStrikeAI Frontend Internationalization

This document describes the internationalization design and development conventions for the CyberStrikeAI web frontend (`web/templates/index.html` + `web/static/js/*.js`). The goal is to provide extensible, low-rework multilingual support without introducing a bundler or changing backend routes.

Current goals:

- **Support Chinese and English switching (`zh-CN` / `en-US`).**
- Make it easy to add more languages later, such as `ja-JP` and `ko-KR`.

---

## 1. Overall Design Principles

- **Client-side internationalization led by the frontend:** all UI text is rendered dynamically in the browser according to the current language. The Go backend provides structure and data but does not distribute languages.
- **One HTML template:** continue using one `index.html` template instead of copying templates for different languages.
- **Separate text from logic:** all visible text is managed through key-value tables (multilingual JSON); HTML and JS use keys rather than hard-coded Chinese or English constants.
- **Progressive migration:** cover critical areas such as the header, login, sidebar, and system settings first, then migrate other pages module by module to avoid a one-time large rewrite.
- **Fallback language:** even when the target language is not fully translated, fall back to the default Chinese language instead of displaying raw keys.

---

## 2. Technology Choices and Directory Structure

### 2.1 Technology choices

- **i18n engine:** use the browser UMD build of [i18next](https://www.i18next.com/) loaded through a CDN; no bundler is required.
- **Resource format:** use one JSON file per language, with hierarchical keys organized by domain and meaning, for example:
  - `common.ok`
  - `nav.dashboard`
  - `header.apiDocs`
  - `settings.robot.wecom.token`

### 2.2 Directory structure

- `web/templates/index.html`
  - Page skeleton and all static text locations, progressively converted to `data-i18n` markers.
- `web/static/js/i18n.js`
  - Frontend i18n initialization and DOM application logic (introduced by this design).
- `web/static/i18n/` (new directory)
  - `zh-CN.json`: Chinese text (default language).
  - `en-US.json`: English text.
  - More languages can be added later, such as `ja-JP.json` and `ko-KR.json`.

---

## 3. Text Organization Conventions

### 3.1 Key naming conventions

- Use the **`module.meaning`** form, with at most two or three levels to keep keys readable:
  - Navigation: `nav.dashboard`, `nav.chat`, `nav.settings`.
  - Header: `header.title`, `header.apiDocs`, `header.logout`.
  - Login: `login.title`, `login.subtitle`, `login.passwordLabel`, `login.submit`.
  - Dashboard: `dashboard.title`, `dashboard.refresh`, `dashboard.runningTasks`.
  - System settings: `settings.title`, `settings.nav.basic`, `settings.nav.robot`, `settings.apply`.
  - Robot configuration: `settings.robot.wecom.enabled`, `settings.robot.wecom.token`, and so on.
- Divide domains by UI area rather than by filename whenever possible, so non-developers can understand them more easily.

### 3.2 JSON example

Example from `web/static/i18n/zh-CN.json`:

```json
{
  "common": {
    "ok": "OK",
    "cancel": "Cancel"
  },
  "nav": {
    "dashboard": "Dashboard",
    "chat": "Chat",
    "infoCollect": "Recon",
    "tasks": "Tasks",
    "vulnerabilities": "Vulnerabilities",
    "settings": "Settings"
  },
  "header": {
    "title": "CyberStrikeAI",
    "apiDocs": "API Docs",
    "logout": "Sign out",
    "language": "Interface language"
  },
  "login": {
    "title": "Sign in to CyberStrikeAI",
    "subtitle": "Enter the access password from config",
    "passwordLabel": "Password",
    "passwordPlaceholder": "Enter password",
    "submit": "Sign in"
  }
}
```

The English file, `en-US.json`, keeps the same keys with different values:

```json
{
  "common": {
    "ok": "OK",
    "cancel": "Cancel"
  },
  "nav": {
    "dashboard": "Dashboard",
    "chat": "Chat",
    "infoCollect": "Recon",
    "tasks": "Tasks",
    "vulnerabilities": "Vulnerabilities",
    "settings": "Settings"
  },
  "header": {
    "title": "CyberStrikeAI",
    "apiDocs": "API Docs",
    "logout": "Sign out",
    "language": "Interface language"
  },
  "login": {
    "title": "Sign in to CyberStrikeAI",
    "subtitle": "Enter the access password from config",
    "passwordLabel": "Password",
    "passwordPlaceholder": "Enter password",
    "submit": "Sign in"
  }
}
```

> Convention: **when adding a UI, define the i18n key first, then use that key in HTML or JS. Do not hard-code Chinese or English text.**

---

## 4. HTML Markup Conventions (`data-i18n`)

### 4.1 Basic rules

- Bind an element's text to a key with `data-i18n`:

```html
<span data-i18n="nav.dashboard">Dashboard</span>
```

- By default, the script replaces the element's `textContent`.
- To translate attributes as well, add `data-i18n-attr` with multiple attribute names separated by commas:

```html
<button
  class="openapi-doc-btn"
  onclick="window.open('/api-docs', '_blank')"
  data-i18n="header.apiDocs"
  data-i18n-attr="title"
  title="API Docs">
  <span data-i18n="header.apiDocs">API Docs</span>
</button>
```

### 4.2 Purpose of default text

- Chinese defaults inside HTML act as **placeholders when JS is unavailable or before initialization**:
  - The page does not become blank or display raw keys while JS is still loading.
  - After initialization, JS overwrites the text with the current language.

---

## 5. Text Conventions in JavaScript

### 5.1 Global translation function `t()`

`i18n.js` exposes these global functions:

- `window.t(key: string): string`
  - Returns the translation for the current language. If it is missing, it falls back to the default language, and finally returns the key itself.
- `window.changeLanguage(lang: string): Promise<void>`
  - Switches language and refreshes page text without reloading the entire page.

Example using `web/static/js/settings.js`:

```js
// Before
alert('Failed to load config: ' + error.message);

// After
alert(t('settings.loadConfigFailed') + ': ' + error.message);
```

> Convention: **all user-facing notices, button labels, and dialog titles in JS should be obtained through `t()`. Do not hard-code Chinese or English constants.**

### 5.2 Progressive migration recommendations

- Migrate these areas first:
  - Frequently displayed error and success notices.
  - Login-related and system-settings-related text.
- Lower priority:
  - Debug notices intended only for operators may temporarily retain English or Chinese constants.

---

## 6. i18n Initialization and Language Switching

### 6.1 Language selection strategy

- Default language: `zh-CN`.
- Priority from highest to lowest:
  1. The user's selection in `localStorage` (key: `csai_lang`).
  2. The browser's `navigator.language` (`zh` prefix maps to `zh-CN`; otherwise use `en-US`).
  3. The default `zh-CN`.

### 6.2 Initialization flow (`i18n.js`)

1. Read the initial language.
2. Initialize i18next:
   - Set `lng` to the current language.
   - Set `fallbackLng` to `zh-CN`.
   - Leave resources empty initially and load them on demand.
3. Fetch `/static/i18n/{lng}.json` and call `i18next.addResources`.
4. Update:
   - The `<html lang="...">` attribute.
   - Every element with `data-i18n` or `data-i18n-attr`.
5. Expose `window.t` and `window.changeLanguage`.

### 6.3 DOM application logic

Pseudocode:

```js
function applyTranslations(root = document) {
  const elements = root.querySelectorAll('[data-i18n]');
  elements.forEach(el => {
    const key = el.getAttribute('data-i18n');
    if (!key) return;
    const text = i18next.t(key);
    if (text) {
      el.textContent = text;
    }

    const attrList = el.getAttribute('data-i18n-attr');
    if (attrList) {
      attrList.split(',').map(s => s.trim()).forEach(attr => {
        if (!attr) return;
        const val = i18next.t(key);
        if (val) el.setAttribute(attr, val);
      });
    }
  });
}
```

> For elements inserted dynamically by JS, call `applyTranslations(newContainer)` again after insertion.

---

## 7. Language-Switcher UI Conventions

### 7.1 Position and form

- Position: near the `API Docs` button on the right side of the `index.html` header, close to the user avatar.
- Interaction:
  - Use a compact language switcher, such as a dropdown button with a `🌐` icon and the current language text (`Chinese` / `English`).
  - List all available languages in the dropdown.

### 7.2 Example structure

```html
<div class="lang-switcher">
  <button class="btn-secondary lang-switcher-btn" onclick="toggleLangDropdown()" data-i18n="header.language">
    <span class="lang-switcher-icon">🌐</span>
    <span id="current-lang-label">Chinese</span>
  </button>
  <div id="lang-dropdown" class="lang-dropdown" style="display: none;">
    <div class="lang-option" data-lang="zh-CN" onclick="onLanguageSelect('zh-CN')">Chinese</div>
    <div class="lang-option" data-lang="en-US" onclick="onLanguageSelect('en-US')">English</div>
  </div>
</div>
```

Corresponding JS in `i18n.js`:

```js
function onLanguageSelect(lang) {
  changeLanguage(lang).then(updateLangLabel).catch(console.error);
  closeLangDropdown();
}

function updateLangLabel() {
  const labelEl = document.getElementById('current-lang-label');
  if (!labelEl) return;
  const lang = i18next.language || 'zh-CN';
  labelEl.textContent = lang.startsWith('zh') ? 'Chinese' : 'English';
}
```

> Convention: **language switching updates text only; it does not reload the entire page or modify the URL hash.**

---

## 8. Recommended Development Workflow

### 8.1 Workflow for adding or changing a UI

1. List all text while designing the UI.
2. Add or update the keys and translations in the corresponding language JSON files.
3. Use `data-i18n` in HTML and `t('...')` in JS.
4. Switch between Chinese and English in the browser and confirm that both languages render correctly.

### 8.2 Recommended progressive migration order

1. **Phase 1 (planned):**
   - Introduce i18next and `i18n.js`.
   - Create `zh-CN.json` and `en-US.json`, initially covering the header, login, and left navigation.
   - Implement the language-switcher component in the header.
2. **Phase 2 (completed):**
   - Internationalize all text on the system-settings pages, including the robot configuration page.
   - Use `t()` for notices and errors in `settings.js`.
3. **Phase 3 (in progress):**
   - Progressively migrate the dashboard, task management, vulnerability management, MCP, Skills, Roles, and other pages module by module.
4. **Phase 4:**
   - Remove hard-coded text remaining in JS and HTML and route it consistently through i18n.

---

## 9. Adding More Languages Later

When a new language is needed:

1. Add `{lang}.json` under `web/static/i18n/`, copy the structure of the existing English or Chinese file, and add the corresponding translations.
2. Add the language to the language-switcher dropdown, for example:
   - `data-lang="ja-JP"` with the text `Japanese`.
3. No changes to `i18n.js` or the existing HTML/JS logic should be required.

---

## 10. Notes and Common Pitfalls

- **Do not copy multiple HTML templates** for multilingual support; that creates a high maintenance cost. Use frontend i18n consistently.
- **Do not use Chinese or English sentences directly as keys.** Use short `module.meaning` keys to make diffs and searches easier.
- Avoid hard-coded text in CSS, such as `content: "xxx"`. If needed, set it through JS and route it through i18n.
- For localizable error text returned by the backend (if supported in the future), prefer having the backend return the appropriate language according to `Accept-Language`; the frontend should only display it.
