// Lightweight i18n for pgweb.
//
// English is the source language: dictionary keys are the original English
// strings (or short symbolic keys for HTML blocks). Ukrainian translations
// live in the `uk` table below. Selecting a language persists the choice in
// localStorage and reloads the page so both static markup and dynamically
// generated content pick up the new language.
var I18n = (function() {
  var translations = {
    uk: {
      // --- Navigation / tabs ---
      "Rows":         "Рядки",
      "Structure":    "Структура",
      "Indexes":      "Індекси",
      "Constraints":  "Обмеження",
      "Query":        "Запит",
      "History":      "Історія",
      "Activity":     "Активність",
      "Connection":   "Підключення",
      "Connect":      "Підключитися",
      "Disconnect":   "Відключитися",

      // --- Sidebar ---
      "Search database":         "Пошук бази даних",
      "Refresh tables list":     "Оновити список таблиць",
      "Filter database objects": "Фільтр обʼєктів бази даних",
      "Table Information":       "Інформація про таблицю",
      "Size:":                   "Розмір:",
      "Data size:":              "Розмір даних:",
      "Index size:":             "Розмір індексу:",
      "Estimated rows:":         "Орієнтовно рядків:",

      // --- Schema object groups ---
      "Tables":             "Таблиці",
      "Views":              "Подання",
      "Materialized Views": "Матеріалізовані подання",
      "Functions":          "Функції",
      "Sequences":          "Послідовності",

      // --- Query area ---
      "Run Query":                          "Виконати запит",
      "Explain Query":                      "Пояснити запит",
      "Analyze Query":                      "Аналізувати запит",
      "Template":                           "Шаблон",
      "Please wait, query is executing...": "Зачекайте, запит виконується...",

      // --- Pagination / filters ---
      "Search":                      "Пошук",
      "Select filter":               "Оберіть фільтр",
      "Filter value":                "Значення фільтра",
      "Apply":                       "Застосувати",
      "Click to change row limit":   "Натисніть, щоб змінити ліміт рядків",
      "rows":                        "рядків",
      "Select column":               "Оберіть стовпець",

      // --- Authorization notice ---
      "auth_required_notice":
        "Потрібна авторизація. Відкрийте цю сторінку за дійсним посиланням доступу.",
      "You have been disconnected. You can now close this tab.":
        "Ви відключилися. Тепер можна закрити цю вкладку.",

      // --- Cell content modal ---
      "Cell Content": "Вміст комірки",

      // --- Row editing ---
      "Add row":          "Додати рядок",
      "Edit row":         "Редагувати рядок",
      "Duplicate row":    "Дублювати рядок",
      "Delete row":       "Видалити рядок",
      "Save":             "Зберегти",
      "Are you sure you want to delete this row?":
        "Ви впевнені, що хочете видалити цей рядок?",
      "This table has no primary key, rows cannot be edited.":
        "Ця таблиця не має первинного ключа, рядки не можна редагувати.",

      // --- Create table ---
      "Create table":   "Створити таблицю",
      "Add column":     "Додати стовпець",
      "Create":         "Створити",
      "Column":         "Стовпець",
      "Type":           "Тип",
      "Default":        "За замовчуванням",
      "schema":         "схема",
      "table name":     "назва таблиці",
      "Table name is required": "Потрібна назва таблиці",
      "Add at least one named column":
        "Додайте принаймні один іменований стовпець",

      // --- Connection window ---
      "Scheme":                  "Схема",
      "Standard":                "Стандартне",
      "Enter server URL scheme": "Введіть URL-схему сервера",
      "Bookmark":                "Закладка",
      "Host":                    "Хост",
      "Username":                "Імʼя користувача",
      "Password":                "Пароль",
      "Database":                "База даних",
      "SSL Mode":                "Режим SSL",
      "SSH Connection":          "Підключення SSH",
      "Credentials":             "Облікові дані",
      "Auth Key":                "Ключ автентифікації",
      "Key path":                "Шлях до ключа",
      "Key password":            "Пароль ключа",
      "Cancel":                  "Скасувати",
      "connection_url_help":
        'Формат URL: postgres://user:password@host:port/db?sslmode=mode<br/>' +
        'Докладніше про <a href="https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING" target="_blank">формат рядка підключення</a> PostgreSQL.',

      // --- Context menus ---
      "Copy Table Name":            "Копіювати імʼя таблиці",
      "Analyze Table":              "Аналізувати таблицю",
      "Export to JSON":             "Експорт у JSON",
      "Export to CSV":              "Експорт у CSV",
      "Export to XML":              "Експорт у XML",
      "Export to SQL":              "Експорт у SQL",
      "Truncate Table":             "Очистити таблицю",
      "Delete Table":               "Видалити таблицю",
      "View Definition":            "Переглянути визначення",
      "Copy View Name":             "Копіювати імʼя подання",
      "Copy View Definition":       "Копіювати визначення подання",
      "Delete View":                "Видалити подання",
      "Show Database Stats":        "Показати статистику бази даних",
      "Download Database Stats":    "Завантажити статистику бази даних",
      "Show Server Settings":       "Показати налаштування сервера",
      "Export SQL dump":            "Експортувати SQL-дамп",
      "Unique Values":              "Унікальні значення",
      "Unique Values + Counts":     "Унікальні значення + кількість",
      "Numeric stats (min/max/avg)": "Числова статистика (min/max/avg)",
      "Copy Column Name":           "Копіювати імʼя стовпця",
      "Display Value":              "Показати значення",
      "Copy Value":                 "Копіювати значення",
      "Filter Rows By Value":       "Фільтрувати рядки за значенням",

      // --- Dynamic messages (app.js) ---
      "Sorry, something went wrong with your request. Refresh the page and try again!":
        "Вибачте, щось пішло не так із вашим запитом. Оновіть сторінку та спробуйте ще раз!",
      "Query timeout after {n}s":        "Час очікування запиту вичерпано через {n} с",
      "Failed to parse the JSON response.": "Не вдалося розібрати JSON-відповідь.",
      "Error while fetching schemas: {error}":
        "Помилка під час отримання схем: {error}",
      "Error while fetching database objects: {error}":
        "Помилка під час отримання обʼєктів бази даних: {error}",
      "truncate": "очистити",
      "delete":   "видалити",
      "Are you sure you want to {action} table {table} ?":
        "Ви впевнені, що хочете {action} таблицю {table}?",
      "Are you sure you want to {action} view {view} ?":
        "Ви впевнені, що хочете {action} подання {view}?",
      "Are you sure you want to stop the query?":
        "Ви впевнені, що хочете зупинити запит?",
      "Are you sure you want to disconnect?":
        "Ви впевнені, що хочете відключитися?",
      "ERROR:":           "ПОМИЛКА:",
      "No records found": "Записів не знайдено",
      "ms":               "мс",
      "{rows} rows in {ms} ms": "{rows} рядків за {ms} мс",
      "{rows} rows":            "{rows} рядків",
      "stop":             "зупинити",
      "Please select a table!": "Будь ласка, оберіть таблицю!",
      "Unknown":          "Невідомо",
      "{page} of {pages}": "{page} з {pages}",
      "Cant view rows for a function": "Неможливо переглянути рядки для функції",
      "View definition for:":     "Визначення подання:",
      "Function definition for:": "Визначення функції:",
      "Update available. Check out {tag} on {github}":
        "Доступне оновлення. Перегляньте {tag} на {github}",
      "Select a bookmarked database to connect to":
        "Оберіть базу даних із закладок для підключення",
      "Running in <b>bookmarks-only</b> mode but <b>NO</b> bookmarks configured.":
        "Працює в режимі <b>лише закладки</b>, але <b>ЖОДНОЇ</b> закладки не налаштовано.",
      "Shortcut: {keys}": "Комбінація: {keys}",
      "Please specify filter query":      "Будь ласка, вкажіть запит фільтра",
      "Please specify a new rows limit":  "Будь ласка, вкажіть новий ліміт рядків",
      "Please wait...":                   "Зачекайте...",
      "Unable to fetch app info: {error}. Please reload the browser page.":
        "Не вдалося отримати інформацію про застосунок: {error}. Будь ласка, перезавантажте сторінку браузера."
    }
  };

  var available = [["en", "EN"], ["uk", "UK"]];

  var lang = localStorage.getItem("pgweb_lang");
  if (!lang) {
    // Auto-detect from the browser, default to English.
    var nav = (navigator.language || navigator.userLanguage || "en").toLowerCase();
    lang = nav.indexOf("uk") === 0 ? "uk" : "en";
  }

  // Translate a string. `vars` is an optional map of {placeholder: value}
  // replacing `{placeholder}` tokens. Falls back to the original string when
  // no translation exists for the active language.
  function t(str, vars) {
    var dict = translations[lang];
    var out  = (dict && dict[str] != null) ? dict[str] : str;

    if (vars) {
      for (var k in vars) {
        out = out.replace(new RegExp("\\{" + k + "\\}", "g"), vars[k]);
      }
    }
    return out;
  }

  function setLang(next) {
    if (next === lang) return;
    localStorage.setItem("pgweb_lang", next);
    window.location.reload();
  }

  // Apply translations to static markup. Elements opt in via attributes:
  //   data-i18n            -> textContent
  //   data-i18n-html       -> innerHTML
  //   data-i18n-placeholder/title/value -> matching attribute
  // The original (English) markup is left untouched when no translation
  // exists, so the source language never needs its own table.
  function translatePage() {
    document.documentElement.lang = lang;

    document.querySelectorAll("[data-i18n]").forEach(function(el) {
      var v = t(el.getAttribute("data-i18n"));
      if (v !== el.getAttribute("data-i18n")) el.textContent = v;
    });

    document.querySelectorAll("[data-i18n-html]").forEach(function(el) {
      var key = el.getAttribute("data-i18n-html");
      var v = t(key);
      if (v !== key) el.innerHTML = v;
    });

    ["placeholder", "title", "value"].forEach(function(attr) {
      document.querySelectorAll("[data-i18n-" + attr + "]").forEach(function(el) {
        var key = el.getAttribute("data-i18n-" + attr);
        var v = t(key);
        if (v !== key) el.setAttribute(attr, v);
      });
    });
  }

  function renderSwitcher() {
    var box = document.createElement("div");
    box.id = "lang_switcher";

    available.forEach(function(pair) {
      var a = document.createElement("a");
      a.href = "#";
      a.className = "lang-option" + (pair[0] === lang ? " active" : "");
      a.textContent = pair[1];
      a.addEventListener("click", function(e) {
        e.preventDefault();
        setLang(pair[0]);
      });
      box.appendChild(a);
    });

    document.body.appendChild(box);
  }

  function injectStyles() {
    var css =
      "#lang_switcher{position:fixed;bottom:8px;right:10px;z-index:10000;" +
      "font-size:11px;background:rgba(255,255,255,0.9);border:1px solid #ccc;" +
      "border-radius:3px;overflow:hidden;box-shadow:0 1px 2px rgba(0,0,0,0.15);}" +
      "#lang_switcher .lang-option{display:inline-block;padding:3px 8px;" +
      "color:#555;text-decoration:none;cursor:pointer;}" +
      "#lang_switcher .lang-option:hover{background:#eee;}" +
      "#lang_switcher .lang-option.active{background:#4682b4;color:#fff;}";
    var style = document.createElement("style");
    style.appendChild(document.createTextNode(css));
    document.head.appendChild(style);
  }

  document.addEventListener("DOMContentLoaded", function() {
    injectStyles();
    translatePage();
    renderSwitcher();
  });

  return {
    t: t,
    setLang: setLang,
    translatePage: translatePage,
    lang: function() { return lang; }
  };
})();

// Global shorthand used throughout app.js.
function t(str, vars) { return I18n.t(str, vars); }
