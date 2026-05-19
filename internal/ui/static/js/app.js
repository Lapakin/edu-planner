'use strict';

// ═══════════════════════════════════════════════════════════════
// TOAST
// ═══════════════════════════════════════════════════════════════
const Toast = {
  _show(msg, type) {
    const c = document.getElementById('toast-container');
    if (!c) return;
    const el = document.createElement('div');
    el.className = `toast toast-${type}`;
    el.textContent = msg;
    c.appendChild(el);
    setTimeout(() => {
      el.classList.add('toast-exit');
      el.addEventListener('animationend', () => el.remove(), { once: true });
    }, 3500);
  },
  success: (m) => Toast._show(m, 'success'),
  error:   (m) => Toast._show(m, 'error'),
  warning: (m) => Toast._show(m, 'warning'),
  info:    (m) => Toast._show(m, 'info'),
};

// ═══════════════════════════════════════════════════════════════
// MODAL
// ═══════════════════════════════════════════════════════════════
function openModal(title, bodyEl, footerEl) {
  document.getElementById('modal-title').textContent = title;
  const body = document.getElementById('modal-body');
  const foot = document.getElementById('modal-footer');
  body.innerHTML = '';
  foot.innerHTML = '';
  if (bodyEl instanceof HTMLElement) body.appendChild(bodyEl); else body.innerHTML = bodyEl || '';
  if (footerEl instanceof HTMLElement) foot.appendChild(footerEl); else foot.innerHTML = footerEl || '';
  document.getElementById('modal').classList.remove('modal-fullscreen');
  document.getElementById('modal-overlay').classList.remove('hidden');
}

function closeModal() {
  document.getElementById('modal-overlay').classList.add('hidden');
  document.getElementById('modal').classList.remove('modal-fullscreen');
}

function closeModalOnOverlay(e) {
  if (e.target === document.getElementById('modal-overlay')) closeModal();
}

function confirmDialog(msg) {
  return new Promise(resolve => {
    const overlay = document.getElementById('confirm-overlay');
    document.getElementById('confirm-message').textContent = msg;
    overlay.classList.remove('hidden');
    const okBtn     = document.getElementById('confirm-ok');
    const cancelBtn = document.getElementById('confirm-cancel');
    const cleanup = (val) => {
      overlay.classList.add('hidden');
      okBtn.replaceWith(okBtn.cloneNode(true));
      cancelBtn.replaceWith(cancelBtn.cloneNode(true));
      resolve(val);
    };
    document.getElementById('confirm-ok').addEventListener('click',     () => cleanup(true),  { once: true });
    document.getElementById('confirm-cancel').addEventListener('click', () => cleanup(false), { once: true });
  });
}

// ═══════════════════════════════════════════════════════════════
// APP
// ═══════════════════════════════════════════════════════════════
const App = {
  state: { token: null, user: null, yearId: null, semesterId: null },
  cache: {},
  currentMgr: null,
  _navGuardActive: false,

  // ── Init ───────────────────────────────────────────────────
  async init() {
    applyTranslations();

    // Restore auth
    const saved = localStorage.getItem('token');
    if (saved && Api.isTokenValid()) {
      this.setAuth(saved);
    } else {
      Api.setToken(null);
      this.showLogin();
      return;
    }

    this.showApp();
    await this.loadYears();
    await this.checkSetup();
    this.setupHashRouter();
    window.dispatchEvent(new HashChangeEvent('hashchange'));
  },

  // ── Auth ───────────────────────────────────────────────────
  setAuth(token) {
    this.state.token = token;
    Api.setToken(token);
    this.state.user = Api.parseJWT(token);
    localStorage.setItem('token', token);
  },

  clearAuth() {
    this.state.token = null;
    this.state.user  = null;
    this.cache = {};
    Api.setToken(null);
    localStorage.removeItem('token');
  },

  showLogin() {
    document.getElementById('login-page').classList.remove('hidden');
    document.getElementById('app-shell').classList.add('hidden');
  },

  showApp() {
    document.getElementById('login-page').classList.add('hidden');
    document.getElementById('app-shell').classList.remove('hidden');
    // Show admin-only nav items
    const role = this.state.user?.role;
    document.querySelectorAll('.admin-only').forEach(el => {
      el.classList.toggle('hidden', role !== 'admin');
    });
    // Set profile button initial letter
    const email   = this.state.user?.email || '';
    const initial = (email[0] || 'U').toUpperCase();
    const profileBtn = document.getElementById('profile-btn');
    if (profileBtn) profileBtn.textContent = initial;
    // Sync lang buttons
    applyTranslations();
  },

  // ── Router ─────────────────────────────────────────────────
  setupHashRouter() {
    window.addEventListener('hashchange', async (e) => {
      const newHash = location.hash || '#/home';
      if (this.currentMgr?.isDirty() && !this._navGuardActive) {
        this._navGuardActive = true;
        const ok = await confirmDialog(t('messages.unsavedChanges'));
        this._navGuardActive = false;
        if (!ok) {
          history.replaceState(null, '', e.oldURL.includes('#') ? '#' + e.oldURL.split('#')[1] : '#/home');
          return;
        }
        this.currentMgr.discard();
      }
      this.route(newHash);
    });
    window.addEventListener('auth:logout', () => { this.clearAuth(); this.showLogin(); });
  },

  route(hash) {
    const path  = hash.replace('#/', '').replace('#', '') || 'home';
    const parts = path.split('/');
    const page  = parts[0];
    const tab   = parts[1] || null;

    document.querySelectorAll('.nav-item').forEach(a => {
      a.classList.toggle('active', a.dataset.page === page);
    });

    const main = document.getElementById('main-content');

    if (page === 'home')                 this.renderHome(main);
    else if (page === 'reference')       this.renderReference(tab, main);
    else if (page === 'processing')      this.renderProcessing(tab, main);
    else if (page === 'schedule')        this.renderSchedule(tab, main);
    else if (page === 'academic-settings') this.renderAcademicSettings(tab, main);
    else if (page === 'profile')         this.renderProfile(main);
    else if (page === 'users')           this.renderUsers(main);
    else this.renderHome(main);
  },

  // ── Header ─────────────────────────────────────────────────
  async loadYears() {
    try {
      const years = await Api.get('/api/syllabus/academic-years');
      const sel = document.getElementById('year-select');
      sel.innerHTML = `<option value="">${t('header.selectYear')}</option>`;
      (years || []).forEach(y => {
        const o = document.createElement('option');
        o.value = y.id; o.textContent = y.name;
        sel.appendChild(o);
      });
      this.cache['academic-years'] = years || [];

      // Use saved year if still valid, otherwise auto-select active year (or first)
      const savedYear = localStorage.getItem('yearId');
      let selectedYear = (years || []).find(y => String(y.id) === savedYear);
      if (!selectedYear) {
        selectedYear = (years || []).find(y => y.is_active) || (years || [])[0];
      }
      if (selectedYear) {
        sel.value = selectedYear.id;
        this.state.yearId = String(selectedYear.id);
        localStorage.setItem('yearId', this.state.yearId);
        await this.loadSemesters(this.state.yearId);
      }

      sel.addEventListener('change', async () => {
        this.state.yearId = sel.value || null;
        localStorage.setItem('yearId', this.state.yearId || '');
        await this.loadSemesters(this.state.yearId);
        if (this.currentMgr) { await this._reloadCurrentMgr(); }
      });
    } catch { /* ignore */ }
  },

  async loadSemesters(yearId) {
    const sel = document.getElementById('semester-select');
    sel.innerHTML = `<option value="">${t('header.selectSemester')}</option>`;
    if (!yearId) { this.cache['semesters'] = []; return; }
    try {
      const all = await Api.get('/api/syllabus/semesters');
      const semesters = (all || []).filter(s => String(s.academic_year_id) === String(yearId));
      semesters.forEach(s => {
        const o = document.createElement('option');
        o.value = s.id;
        o.textContent = `${s.period_start?.slice(0,10)} – ${s.period_end?.slice(0,10)}`;
        sel.appendChild(o);
      });
      this.cache['semesters'] = all || [];
      const savedSem = localStorage.getItem('semesterId');
      if (savedSem && semesters.find(s => String(s.id) === savedSem)) {
        sel.value = savedSem;
        this.state.semesterId = savedSem;
      }
      sel.addEventListener('change', () => {
        this.state.semesterId = sel.value || null;
        localStorage.setItem('semesterId', this.state.semesterId || '');
      });
    } catch { /* ignore */ }
  },

  async _reloadCurrentMgr() {
    if (!this.currentMgr) return;
    try {
      const rows = await Api.get(this.currentMgr.config.apiPath);
      this.currentMgr.setRows(rows || []);
      this.currentMgr._renderTable();
    } catch { /* ignore */ }
  },

  // ── Cache ──────────────────────────────────────────────────
  async ensureCache(entityKey) {
    if (this.cache[entityKey]) return;
    const cfg = ENTITY_CONFIGS[entityKey];
    if (!cfg) return;
    try {
      const rows = await Api.get(cfg.apiPath);
      this.cache[entityKey] = rows || [];
    } catch { this.cache[entityKey] = []; }
  },

  async refreshCache(apiPath) {
    const key = Object.keys(ENTITY_CONFIGS).find(k => ENTITY_CONFIGS[k].apiPath === apiPath);
    if (key) {
      try {
        const rows = await Api.get(apiPath);
        this.cache[key] = rows || [];
      } catch { /* ignore */ }
    }
  },

  // ── First-login check ──────────────────────────────────────
  async checkSetup() {
    try {
      const years = await Api.get('/api/syllabus/academic-years');
      if (!years || years.length === 0) this.showSetupModal();
    } catch { /* ignore */ }
  },

  showSetupModal() {
    const body = document.createElement('div');
    body.innerHTML = `
      <p style="margin-bottom:16px;color:var(--text-muted)">${t('setupModal.description')}</p>
      <div class="form-field">
        <label>${t('setupModal.yearName')} <span class="req">*</span></label>
        <input type="text" id="setup-year-name" placeholder="2025-2026">
      </div>
      <div class="form-field">
        <label>${t('setupModal.yearStart')} <span class="req">*</span></label>
        <input type="date" id="setup-year-start">
      </div>
      <div class="form-field">
        <label>${t('setupModal.yearEnd')} <span class="req">*</span></label>
        <input type="date" id="setup-year-end">
      </div>
      <hr style="margin:16px 0">
      <p style="font-weight:600;margin-bottom:8px">${t('entities.semesters')} (${t('actions.optional') || 'optional'})</p>
      <div class="form-field">
        <label>${t('setupModal.semesterStart')}</label>
        <input type="date" id="setup-sem-start">
      </div>
      <div class="form-field">
        <label>${t('setupModal.semesterEnd')}</label>
        <input type="date" id="setup-sem-end">
      </div>`;

    const foot = document.createElement('div');
    foot.innerHTML = `<button class="btn btn-primary" id="setup-create-btn">${t('setupModal.createBtn')}</button>`;
    openModal(t('setupModal.title'), body, foot);

    foot.querySelector('#setup-create-btn').addEventListener('click', async () => {
      const name  = document.getElementById('setup-year-name').value.trim();
      const start = document.getElementById('setup-year-start').value;
      const end   = document.getElementById('setup-year-end').value;
      if (!name || !start || !end) { Toast.warning(t('messages.fillRequired') || 'Fill required fields'); return; }
      try {
        const payload = {
          name,
          start_date: start.length === 10 ? start + 'T00:00:00Z' : start,
          end_date: end.length === 10 ? end + 'T00:00:00Z' : end
        };
        const created = await Api.post('/api/syllabus/academic-years', [payload]);
        this.cache['academic-years'] = created || [];
        const yearId = created?.[0]?.id;
        // optional semester
        const ss = document.getElementById('setup-sem-start').value;
        const se = document.getElementById('setup-sem-end').value;
        if (yearId && ss && se) {
          const semPayload = {
            academic_year_id: yearId,
            period_start: ss.length === 10 ? ss + 'T00:00:00Z' : ss,
            period_end: se.length === 10 ? se + 'T00:00:00Z' : se
          };
          await Api.post('/api/syllabus/semesters', [semPayload]);
        }
        closeModal();
        Toast.success(t('messages.saved'));
        await this.loadYears();
      } catch (e) { Toast.error(e.message); }
    });
  },

  // ── Special row actions ────────────────────────────────────
  async toggleAcademicYear(row, mgr) {
    const activating = !row.is_active;
    const path = `/api/syllabus/academic-years/${row.id}/${row.is_active ? 'deactivate' : 'activate'}`;
    try {
      await Api.post(path, null);
      if (activating) {
        const allYears = await Api.get('/api/syllabus/academic-years');
        this.cache['academic-years'] = allYears || [];
        mgr.setRows(allYears || []);
        mgr._renderTable();
        // Update global year state so filters reflect the newly activated year
        this.state.yearId = String(row.id);
        localStorage.setItem('yearId', this.state.yearId);
        const sel = document.getElementById('year-select');
        if (sel) sel.value = row.id;
        await this.loadSemesters(this.state.yearId);
      } else {
        row.is_active = false;
        mgr.reloadRow(row);
      }
      Toast.success(t('messages.saved'));
    } catch (e) { Toast.error(e.message); }
  },

  async toggleScheduleTemplate(row, mgr) {
    const path = `/api/syllabus/schedule-templates/${row.id}/activate`;
    try {
      await Api.post(path, null);
      row.is_active = !(row.is_active);
      mgr.reloadRow(row);
      Toast.success(t('messages.saved'));
    } catch (e) { Toast.error(e.message); }
  },

  // ═══════════════════════════════════════════════════════════
  // HOME PAGE
  // ═══════════════════════════════════════════════════════════
  renderHome(main) {
    this.currentMgr = null;
    const user = this.state.user;
    main.innerHTML = `
      <div class="page-header">
        <h1>${t('pages.homeTitle')}</h1>
      </div>
      <div class="home-welcome">
        <div class="welcome-card">
          <div class="welcome-icon">
            <svg style="width:48px;height:48px;color:var(--accent)"><use href="#icon-book"/></svg>
          </div>
          <div>
            <h2>${t('pages.homeWelcome')}, ${user?.email || ''}!</h2>
            <p style="color:var(--text-muted);margin-top:8px">${t('roles.' + (user?.role || 'user'))}</p>
          </div>
        </div>
      </div>
      <div id="home-active-tpl" style="margin:16px 0"></div>
      <div class="home-cards" id="home-stats"></div>`;

    this._renderActiveTemplate(main.querySelector('#home-active-tpl'));

    const cards = [
      { i18nKey: 'entities.academicYears', cacheKey: 'academic-years', href: '#/academic-settings/academic-years', icon: 'icon-calendar' },
      { i18nKey: 'entities.departments',   cacheKey: 'departments',     href: '#/reference/departments',   icon: 'icon-book' },
      { i18nKey: 'entities.teachers',      cacheKey: 'teachers',        href: '#/reference/teachers',      icon: 'icon-users' },
      { i18nKey: 'entities.groups',        cacheKey: 'groups',          href: '#/reference/groups',        icon: 'icon-settings' },
      { i18nKey: 'entities.scheduleTemplates', cacheKey: 'schedule-templates', href: '#/schedule/schedule-templates', icon: 'icon-calendar' },
    ];

    const statsEl = document.getElementById('home-stats');
    cards.forEach(async card => {
      const div = document.createElement('a');
      div.href = card.href;
      div.className = 'stat-card';
      let count = '…';
      try {
        await this.ensureCache(card.cacheKey);
        count = (this.cache[card.cacheKey] || []).length;
      } catch { count = '?'; }
      div.innerHTML = `
        <svg class="icon stat-icon"><use href="#${card.icon}"/></svg>
        <div class="stat-count">${count}</div>
        <div class="stat-label">${t(card.i18nKey)}</div>`;
      statsEl.appendChild(div);
    });
  },

  async _renderActiveTemplate(container) {
    try {
      const templates = await Api.get('/api/syllabus/schedule-templates');
      const active = (templates || []).find(tmpl => tmpl.is_active);
      if (!active) {
        container.innerHTML = `<div class="empty-state-hint" style="margin:24px 0;padding:24px;text-align:center;background:var(--surface);border:1px solid var(--border);border-radius:var(--radius-lg);color:var(--text-muted)">
          <svg class="icon" style="width:32px;height:32px;margin-bottom:12px;opacity:.4"><use href="#icon-calendar"/></svg>
          <div>${t('home.noActiveTemplate')}</div>
        </div>`;
        return;
      }
      await this.ensureCache('semesters');
      const sem = (this.cache['semesters'] || []).find(s => s.id === active.semester_id);
      const semText = sem
        ? `${sem.period_start?.slice(0,10)} – ${sem.period_end?.slice(0,10)}`
        : `#${active.semester_id}`;
      container.innerHTML = `
        <div class="card active-tpl-card">
          <div class="card-header">
            <div class="card-title">
              <svg class="icon" style="color:var(--success)"><use href="#icon-check"/></svg>
              ${t('home.activeTemplate')}
            </div>
            <button class="btn btn-secondary btn-sm" id="active-tpl-view">${t('actions.view')}</button>
          </div>
          <div style="display:flex;align-items:center;gap:16px;padding-top:4px">
            <div>
              <div style="font-weight:600;font-size:.95rem">${active.name || ''}</div>
              <div style="color:var(--text-muted);font-size:.82rem;margin-top:2px">${semText}</div>
            </div>
          </div>
        </div>`;
      container.querySelector('#active-tpl-view').addEventListener('click', () => {
        this.viewScheduleTemplate(active);
      });
    } catch { /* ignore */ }
  },

  // ═══════════════════════════════════════════════════════════
  // REFERENCE DATA PAGE
  // ═══════════════════════════════════════════════════════════
  async renderReference(tab, main) {
    const tabs = [
      { key: 'departments',     label: t('tabs.departments') },
      { key: 'rooms',           label: t('tabs.rooms') },
      { key: 'specialties',     label: t('tabs.specialties') },
      { key: 'cycle-committees',label: t('tabs.cycleCommittees') },
      { key: 'teachers',        label: t('tabs.teachers') },
      { key: 'disciplines',     label: t('tabs.disciplines') },
      { key: 'groups',          label: t('tabs.groups') },
    ];
    await this._renderTabPage('reference', tabs, tab || 'departments', main);
  },

  // ═══════════════════════════════════════════════════════════
  // PROCESSING PAGE
  // ═══════════════════════════════════════════════════════════
  async renderProcessing(tab, main) {
    const tabs = [
      { key: 'study-plans',            label: t('tabs.studyPlans') },
      { key: 'workload-distributions', label: t('tabs.workloadDistributions') },
    ];
    await this._renderTabPage('processing', tabs, tab || 'study-plans', main);
  },

  // ═══════════════════════════════════════════════════════════
  // ACADEMIC SETTINGS PAGE (Academic Years, Semesters, Bell Schedules)
  // ═══════════════════════════════════════════════════════════
  async renderAcademicSettings(tab, main) {
    const tabs = [
      { key: 'academic-years',  label: t('tabs.academicYears') },
      { key: 'semesters',       label: t('tabs.semesters') },
      { key: 'bell-schedules',  label: t('tabs.bellSchedules') },
    ];
    await this._renderTabPage('academic-settings', tabs, tab || 'academic-years', main);
  },

  // ═══════════════════════════════════════════════════════════
  // PROFILE PAGE
  // ═══════════════════════════════════════════════════════════
  async renderProfile(main) {
    this.currentMgr = null;
    const jwtUser = this.state.user;
    const userId = jwtUser?.id;
    main.innerHTML = '<div class="loading-state"><div class="spinner"></div></div>';

    let user = jwtUser || {};
    try {
      if (userId) user = await Api.get(`/api/auth/users/${userId}`) || jwtUser;
    } catch { /* use jwt claims as fallback */ }

    const renderPage = (u) => {
      main.innerHTML = `
        <div class="page">
          <h1 class="page-title">${t('pages.profileTitle')}</h1>
          <div class="profile-card">
            <div class="profile-card-top">
              <div class="profile-page-avatar" style="width:64px;height:64px;font-size:1.4rem">${(u.email||'U')[0].toUpperCase()}</div>
              <div>
                <div class="profile-page-email">${u.email||''}</div>
                <span class="badge badge-${u.role==='admin'?'warning':'info'}" style="margin-top:6px">${t('roles.'+(u.role||'user'))}</span>
              </div>
            </div>
            <div class="profile-fields">
              <div class="form-field">
                <label>${t('fields.firstName')}</label>
                <input type="text" id="prof-first" value="${u.first_name||''}">
              </div>
              <div class="form-field">
                <label>${t('fields.lastName')}</label>
                <input type="text" id="prof-last" value="${u.last_name||''}">
              </div>
              <div class="form-field">
                <label>${t('fields.patronymic')} <span style="color:var(--text-muted);font-size:.78rem">(${t('actions.optional')})</span></label>
                <input type="text" id="prof-patron" value="${u.patronymic||''}">
              </div>
            </div>
            <div style="margin-top:16px">
              <button class="btn btn-primary" id="prof-save">
                <svg class="icon icon-sm"><use href="#icon-save"/></svg>
                ${t('actions.save')}
              </button>
            </div>
          </div>
        </div>`;

      main.querySelector('#prof-save').addEventListener('click', async () => {
        const btn = main.querySelector('#prof-save');
        btn.disabled = true;
        try {
          const updated = {
            ...u,
            first_name: main.querySelector('#prof-first').value.trim(),
            last_name: main.querySelector('#prof-last').value.trim(),
            patronymic: main.querySelector('#prof-patron').value.trim() || null,
          };
          await Api.put('/api/auth/users', [updated]);
          Toast.success(t('messages.saved'));
          renderPage(updated);
        } catch(e) { Toast.error(e.message); btn.disabled=false; }
      });
    };

    renderPage(user);
  },

  // ═══════════════════════════════════════════════════════════
  // SCHEDULE PAGE
  // ═══════════════════════════════════════════════════════════
  async renderSchedule(tab, main) {
    const tabs = [
      { key: 'generate',           label: t('tabs.generate') },
      { key: 'schedule-templates', label: t('tabs.templates') },
    ];
    await this._renderTabPage('schedule', tabs, tab || 'generate', main);
  },

  // ═══════════════════════════════════════════════════════════
  // GENERIC TAB PAGE RENDERER
  // ═══════════════════════════════════════════════════════════
  _kebabToCamel(s) {
    return s.replace(/-([a-z])/g, (_, c) => c.toUpperCase());
  },

  async _renderTabPage(pageKey, tabs, activeTabKey, main) {
    main.innerHTML = `
      <div class="page-header">
        <h1>${t('pages.' + this._kebabToCamel(pageKey) + 'Title')}</h1>
      </div>
      <div class="tabs-bar" id="tabs-bar"></div>
      <div id="tab-content" class="tab-content"></div>`;

    const tabsBar = document.getElementById('tabs-bar');
    tabs.forEach(tab => {
      const btn = document.createElement('button');
      btn.className = 'tab-btn' + (tab.key === activeTabKey ? ' active' : '');
      btn.dataset.tab = tab.key;
      btn.textContent = tab.label;
      btn.addEventListener('click', async () => {
        if (this.currentMgr?.isDirty()) {
          const ok = await confirmDialog(t('messages.unsavedChanges'));
          if (!ok) return;
          this.currentMgr.discard();
        }
        location.hash = `#/${pageKey}/${tab.key}`;
      });
      tabsBar.appendChild(btn);
    });

    const content = document.getElementById('tab-content');
    if (activeTabKey === 'generate') {
      await this.renderGenerateSchedule(content);
    } else if (activeTabKey === 'groups') {
      await this.renderGroupsMasterDetail(content);
    } else if (activeTabKey === 'workload-distributions') {
      await this.renderWorkloadMasterDetail(content);
    } else {
      await this.renderEntityTab(activeTabKey, content);
    }
  },

  async renderEntityTab(entityKey, container) {
    if (entityKey === 'groups') {
      await this.renderGroupsMasterDetail(container);
      return;
    }
    if (entityKey === 'workload-distributions') {
      await this.renderWorkloadMasterDetail(container);
      return;
    }

    const cfg = ENTITY_CONFIGS[entityKey];
    if (!cfg) { container.innerHTML = '<p>Unknown entity</p>'; return; }

    // Show loading
    container.innerHTML = '<div class="loading-state"><div class="spinner"></div></div>';

    // Load required caches for relations
    const refsNeeded = new Set();
    cfg.columns.forEach(c => { if (c.ref) refsNeeded.add(c.ref); });
    cfg.filters?.forEach(f => { if (f.ref) refsNeeded.add(f.ref); });

    if (refsNeeded.has('users')) await this.ensureUsersCache();
    if (refsNeeded.has('teacher-users')) await this.ensureTeacherUsersCache();
    for (const ref of refsNeeded) {
      if (ref !== 'users' && ref !== 'teacher-users') await this.ensureCache(ref);
    }

    // Load entity data, appending academic_year_id filter when configured
    let rows = [];
    try {
      let apiPath = cfg.apiPath;
      if (cfg.autoYearFilter && this.state.yearId) {
        apiPath += `?academic_year_id=${this.state.yearId}`;
      }
      rows = await Api.get(apiPath) || [];
    } catch (e) {
      container.innerHTML = `<div class="alert alert-error">${t('messages.loadError')}: ${e.message}</div>`;
      return;
    }

    // Build cache snapshot
    const cacheSnap = {};
    refsNeeded.forEach(ref => { cacheSnap[ref] = this.cache[ref] || []; });

    const mgr = new EntityManager(cfg, cacheSnap);
    mgr.setRows(rows);
    this.currentMgr = mgr;

    container.innerHTML = '';
    mgr.render(container);
  },

  async ensureUsersCache() {
    if (this.cache['users']) return;
    try {
      const users = await Api.get('/api/auth/users');
      this.cache['users'] = users || [];
    } catch { this.cache['users'] = []; }
  },

  async ensureTeacherUsersCache() {
    if (this.cache['teacher-users']) return;
    try {
      const users = await Api.get('/api/auth/users?role=teacher');
      this.cache['teacher-users'] = users || [];
    } catch { this.cache['teacher-users'] = []; }
  },

  // ═══════════════════════════════════════════════════════════
  // GENERATE SCHEDULE TAB
  // ═══════════════════════════════════════════════════════════
  async renderGenerateSchedule(container) {
    this.currentMgr = null;
    await this.ensureCache('semesters');
    await this.ensureCache('groups');

    const semesters = this.cache['semesters'] || [];
    const groups    = this.cache['groups']    || [];

    container.innerHTML = `
      <div class="gen-layout">
        <aside class="gen-form">
          <h3>${t('generate.title')}</h3>
          <div class="form-field">
            <label>${t('generate.semesterLabel')} <span class="req">*</span></label>
            <select id="gen-semester">
              <option value="">— ${t('header.selectSemester')} —</option>
              ${semesters.map(s =>
                `<option value="${s.id}">${s.period_start?.slice(0,10)} – ${s.period_end?.slice(0,10)}</option>`
              ).join('')}
            </select>
          </div>
          <div class="form-field">
            <label>${t('generate.groupsLabel')}</label>
            <div class="groups-checklist" id="gen-groups">
              ${groups.map(g =>
                `<label class="check-item">
                  <input type="checkbox" value="${g.id}" checked>
                  <span>${g.name}</span>
                </label>`
              ).join('')}
            </div>
          </div>
          <button class="btn btn-primary" style="width:100%;margin-bottom:12px" id="gen-btn">${t('generate.generateBtn')}</button>
          <div class="divider"></div>
          <button class="btn btn-secondary btn-sm" style="width:100%;margin-top:10px" id="gen-settings-btn">
            <svg class="icon"><use href="#icon-sliders"/></svg>
            ${t('entities.scheduleTemplateSettings')}
          </button>
          <button class="btn btn-secondary btn-sm" style="width:100%;margin-top:6px" id="gen-restrict-btn">
            <svg class="icon"><use href="#icon-filter"/></svg>
            ${t('entities.scheduleRestrictions')}
          </button>
        </aside>
        <section class="gen-main" id="gen-main">
          <div class="gen-empty">${t('generate.emptyHint')}</div>
        </section>
      </div>`;

    if (this.state.semesterId) {
      const sel = container.querySelector('#gen-semester');
      if (sel) sel.value = this.state.semesterId;
    }

    container.querySelector('#gen-settings-btn').addEventListener('click', () => {
      this._openSettingsModal('schedule-template-settings', t('entities.scheduleTemplateSettings'));
    });

    container.querySelector('#gen-restrict-btn').addEventListener('click', () => {
      this._openSettingsModal('schedule-restrictions', t('entities.scheduleRestrictions'));
    });

    container.querySelector('#gen-btn').addEventListener('click', async () => {
      const semId    = container.querySelector('#gen-semester').value;
      if (!semId) { Toast.warning(t('generate.noSemester')); return; }
      const groupIds = [...container.querySelectorAll('#gen-groups input:checked')].map(c => Number(c.value));
      await this._doGenerate(semId, groupIds, container.querySelector('#gen-main'));
    });
  },

  async _openSettingsModal(entityKey, title) {
    const body = document.createElement('div');
    body.innerHTML = '<div class="loading-state"><div class="spinner"></div></div>';
    const foot = document.createElement('div');
    foot.innerHTML = `<button class="btn btn-secondary" onclick="closeModal()">${t('actions.close')}</button>`;
    openModal(title, body, foot);
    document.getElementById('modal').classList.add('modal-fullscreen');
    await this.renderEntityTab(entityKey, body);
  },

  async _doGenerate(semId, groupIds, mainEl) {
    this.currentMgr = null;
    mainEl.innerHTML = `<div class="generate-progress"><div class="spinner"></div><span>${t('generate.generating')}</span></div>`;
    try {
      const data = await Api.post('/api/syllabus/schedule-templates/generate', {
        semester_id: Number(semId),
        group_ids: groupIds.length > 0 ? groupIds : undefined,
      });
      mainEl.innerHTML = '';

      const saveBar = document.createElement('div');
      saveBar.className = 'gen-save-bar';
      saveBar.innerHTML = `
        <span class="gen-save-label">${t('generate.success')}</span>
        <input type="text" class="gen-name-input" placeholder="${t('generate.saveNamePlaceholder')}">
        <button class="btn btn-success btn-sm" id="gen-save-btn">
          <svg class="icon"><use href="#icon-save"/></svg>
          ${t('generate.saveBtn')}
        </button>`;
      mainEl.appendChild(saveBar);

      const scheduleEl = document.createElement('div');
      mainEl.appendChild(scheduleEl);
      await this._renderScheduleData(data, scheduleEl);

      saveBar.querySelector('#gen-save-btn').addEventListener('click', async () => {
        await this._saveGeneratedSchedule(Number(semId), data, saveBar);
      });
    } catch (e) {
      mainEl.innerHTML = `<div class="alert alert-error">${e.message}</div>`;
    }
  },

  async _saveGeneratedSchedule(semId, scheduleData, saveBar) {
    const btn = saveBar.querySelector('#gen-save-btn');
    const nameInput = saveBar.querySelector('.gen-name-input');
    btn.disabled = true;
    try {
      const payload = { semester_id: semId, data: scheduleData };
      const name = nameInput?.value.trim();
      if (name) payload.name = name;
      await Api.post('/api/syllabus/schedule-templates/save', payload);
      Toast.success(t('generate.savedSuccess'));
      delete this.cache['schedule-templates'];
    } catch (e) {
      Toast.error(e.message);
    } finally {
      btn.disabled = false;
    }
  },

  async viewScheduleTemplate(row) {
    const body = document.createElement('div');
    body.innerHTML = '<div class="loading-state"><div class="spinner"></div></div>';
    const foot = document.createElement('div');
    foot.innerHTML = `<button class="btn btn-secondary" onclick="closeModal()">${t('actions.close')}</button>`;
    openModal(t('schedule.viewTitle'), body, foot);
    document.getElementById('modal').classList.add('modal-fullscreen');
    await this._renderScheduleData(row, body);
  },

  async _renderScheduleData(result, container) {
    await this.ensureCache('disciplines');
    await this.ensureCache('groups');
    await this.ensureCache('rooms');
    await this.ensureCache('teachers');
    await this.ensureCache('bell-schedules');
    await this.ensureUsersCache();

    const disciplines  = this.cache['disciplines']   || [];
    const groups       = this.cache['groups']        || [];
    const rooms        = this.cache['rooms']         || [];
    const teachers     = this.cache['teachers']      || [];
    const users        = this.cache['users']         || [];
    const bellSchedules = this.cache['bell-schedules'] || [];

    const subjectName = id => {
      const d = disciplines.find(x => x.id === id);
      return d?.name || d?.short_name || `#${id}`;
    };
    const groupName   = id => groups.find(g => g.id === id)?.name   || `#${id}`;
    const roomName    = id => rooms.find(r => r.id === id)?.name    || `#${id}`;
    const teacherLabel = id => {
      const teacher = teachers.find(tc => tc.id === id);
      if (!teacher) return `#${id}`;
      const user = users.find(u => u.id === teacher.user_id);
      return user?.email?.split('@')[0] || `#${id}`;
    };
    const bellTime = num => {
      const bs = bellSchedules.find(b => b.lesson_number === num);
      if (!bs) return String(num);
      const s = (bs.start_time || '').slice(0, 5);
      const e = (bs.end_time   || '').slice(0, 5);
      return s && e ? `${s}–${e}` : String(num);
    };

    let scheduleJson;
    try {
      if (result && result.data != null) {
        // Saved template: result.data is base64-encoded JSON bytes
        scheduleJson = JSON.parse(atob(result.data));
      } else {
        // Raw ScheduleData returned directly from the generate endpoint
        scheduleJson = result;
      }
    } catch {
      container.innerHTML = `<p class="text-muted">No schedule data.</p>`;
      return;
    }

    const DAYS = ['monday', 'tuesday', 'wednesday', 'thursday', 'friday'];
    const slots = bellSchedules.length
      ? [...new Set(bellSchedules.map(b => b.lesson_number))].sort((a, b) => a - b)
      : [1, 2, 3, 4, 5];

    // Collect only groups that appear in this schedule
    const scheduleGroupIds = new Set();
    for (const weekData of Object.values(scheduleJson)) {
      if (!weekData) continue;
      for (const dayData of Object.values(weekData)) {
        if (!dayData) continue;
        for (const lessons of Object.values(dayData)) {
          for (const l of lessons) {
            for (const sl of (l.sub_lessons || [])) scheduleGroupIds.add(sl.group_id);
          }
        }
      }
    }
    const scheduleGroups = groups.filter(g => scheduleGroupIds.has(g.id));

    const renderWeek = (weekData, groupId) => {
      let rows = '';
      for (const slot of slots) {
        let cells = `<td class="sched-slot"><span class="sched-slot-num">${slot}</span><span class="sched-slot-time">${bellTime(slot)}</span></td>`;
        for (const day of DAYS) {
          let lessons = weekData?.[day]?.[String(slot)] || [];
          if (groupId) {
            lessons = lessons
              .map(l => ({ ...l, sub_lessons: (l.sub_lessons || []).filter(sl => sl.group_id === groupId) }))
              .filter(l => l.sub_lessons.length > 0);
          }
          let cellContent = '';
          for (const lesson of lessons) {
            const subj = subjectName(lesson.subject_id);
            let details = '';
            for (const sl of (lesson.sub_lessons || [])) {
              details += `<div class="sched-lesson-detail">
                <span class="sched-badge sched-group">${groupName(sl.group_id)}</span>
                <span class="sched-badge sched-teacher">${teacherLabel(sl.teacher_id)}</span>
                <span class="sched-badge sched-room">${roomName(sl.room_id)}</span>
              </div>`;
            }
            cellContent += `<div class="sched-lesson-card sched-type-${lesson.type}">
              <div class="sched-lesson-subject">${subj}</div>
              ${details}
            </div>`;
          }
          cells += `<td class="sched-cell">${cellContent}</td>`;
        }
        rows += `<tr>${cells}</tr>`;
      }

      const dayHeaders = DAYS.map(d => `<th>${t('schedule.days.' + d)}</th>`).join('');
      return `<div class="sched-wrapper">
        <table class="sched-table">
          <thead><tr><th class="sched-slot-hdr">#</th>${dayHeaders}</tr></thead>
          <tbody>${rows}</tbody>
        </table>
      </div>`;
    };

    const groupOptions = scheduleGroups.map(g =>
      `<option value="${g.id}">${g.name}</option>`).join('');

    container.innerHTML = `
      <div class="sched-view">
        <div class="sched-name">${result.name || ''}</div>
        <div class="sched-toolbar">
          <div class="sched-week-tabs">
            <button class="sched-week-btn active" data-week="numerator">${t('schedule.numerator')}</button>
            <button class="sched-week-btn" data-week="denominator">${t('schedule.denominator')}</button>
          </div>
          <div class="sched-group-filter">
            <span>${t('schedule.groupFilter')}</span>
            <select id="sched-group-select">
              <option value="">— ${t('schedule.allGroups')} —</option>
              ${groupOptions}
            </select>
          </div>
        </div>
        <div id="sched-content"></div>
      </div>`;

    const content = container.querySelector('#sched-content');
    const activeGroupId = () => Number(container.querySelector('#sched-group-select').value) || 0;
    const activeWeek = () => container.querySelector('.sched-week-btn.active')?.dataset.week || 'numerator';

    const showWeek = week => {
      content.innerHTML = renderWeek(scheduleJson[week] || {}, activeGroupId());
      container.querySelectorAll('.sched-week-btn').forEach(btn =>
        btn.classList.toggle('active', btn.dataset.week === week));
    };

    container.querySelectorAll('.sched-week-btn').forEach(btn =>
      btn.addEventListener('click', () => showWeek(btn.dataset.week)));

    container.querySelector('#sched-group-select').addEventListener('change', () =>
      showWeek(activeWeek()));

    showWeek('numerator');
  },

  // ═══════════════════════════════════════════════════════════
  // GROUPS MASTER-DETAIL
  // ═══════════════════════════════════════════════════════════
  async renderGroupsMasterDetail(container) {
    await this.ensureCache('specialties');
    await this.ensureCache('semesters');

    let groups = [];
    let groupSemesters = [];
    try {
      const gpPath = ENTITY_CONFIGS['groups'].apiPath +
        (this.state.yearId ? `?academic_year_id=${this.state.yearId}` : '');
      groups = await Api.get(gpPath) || [];
      groupSemesters = await Api.get(ENTITY_CONFIGS['group-semesters'].apiPath) || [];
    } catch(e) {
      container.innerHTML = `<div class="alert alert-error">${e.message}</div>`;
      return;
    }
    this.cache['groups'] = groups;

    const specialties = this.cache['specialties'] || [];
    const semesters = this.cache['semesters'] || [];
    let selectedId = groups.length > 0 ? groups[0].id : null;
    let filterSpec = '';

    const redraw = () => {
      const filtered = groups.filter(g => {
        if (filterSpec && String(g.specialty_id) !== filterSpec) return false;
        return true;
      });
      if (!filtered.find(g => g.id === selectedId) && filtered.length) selectedId = filtered[0].id;
      const sel = groups.find(g => g.id === selectedId) || null;
      const selSems = groupSemesters.filter(gs => gs.group_id === selectedId).sort((a,b)=>a.start_date<b.start_date?-1:1);

      container.innerHTML = `
        <div class="gmd-wrap">
          <div class="gmd-toolbar">
            <select id="gmd-spec">
              <option value="">— ${t('fields.specialty')} —</option>
              ${specialties.map(s=>`<option value="${s.id}"${filterSpec===String(s.id)?' selected':''}>${s.name}</option>`).join('')}
            </select>
            <button class="btn btn-primary btn-sm" id="gmd-add">
              <svg class="icon icon-sm"><use href="#icon-plus"/></svg>${t('actions.add')}
            </button>
          </div>

          <div class="gmd-layout">
            <div class="gmd-list-panel">
              <div class="gmd-list-hdr">${t('entities.groups')} (${filtered.length})</div>
              ${filtered.length === 0 ? `<div class="gmd-empty">${t('messages.noData')}</div>` :
                filtered.map(g => {
                  const sp = specialties.find(s=>s.id===g.specialty_id);
                  const isSel = g.id === selectedId;
                  const gSemCount = groupSemesters.filter(gs=>gs.group_id===g.id).length;
                  return `<div class="gmd-item${isSel?' gmd-item-sel':''}" data-gid="${g.id}">
                    <div class="gmd-item-row">
                      <span class="gmd-item-name">${g.name}</span>
                      <span class="gmd-item-meta">${gSemCount} ${t('fields.semester').toLowerCase()}</span>
                    </div>
                    <div class="gmd-item-sub">${sp?.name||'—'} · ${(g.education_start||'').slice(0,4)}–${(g.education_end||'').slice(0,4)}</div>
                  </div>`;
                }).join('')}
            </div>

            <div class="gmd-detail-panel">
              ${sel ? `
                <div class="gmd-detail-header">
                  <div>
                    <h3 class="gmd-detail-name">${sel.name}</h3>
                    <div class="gmd-detail-sub">${specialties.find(s=>s.id===sel.specialty_id)?.name||'—'} · ${(sel.education_start||'').slice(0,10)} — ${(sel.education_end||'').slice(0,10)}</div>
                  </div>
                  <div style="display:flex;gap:6px">
                    <button class="btn btn-sm btn-secondary" id="gmd-edit-grp"><svg class="icon icon-sm"><use href="#icon-save"/></svg>${t('actions.edit')}</button>
                    <button class="btn btn-sm btn-danger" id="gmd-del-grp"><svg class="icon icon-sm"><use href="#icon-trash"/></svg></button>
                  </div>
                </div>
                <div class="gmd-props">
                  <div class="gmd-prop"><span class="gmd-prop-lbl">${t('fields.isContract')}</span><span>${sel.is_contract?'✓':'✗'}</span></div>
                  <div class="gmd-prop"><span class="gmd-prop-lbl">${t('fields.isSplitting')}</span><span>${sel.is_splitting?'✓':'✗'}</span></div>
                  <div class="gmd-prop"><span class="gmd-prop-lbl">${t('fields.shortName')}</span><span>${sel.short_name||'—'}</span></div>
                </div>

                <div class="gmd-sems-header">
                  <span class="gmd-sems-title">${t('entities.semesters')}</span>
                  <button class="btn btn-sm btn-secondary" id="gmd-add-sem">
                    <svg class="icon icon-sm"><use href="#icon-plus"/></svg>${t('actions.add')}
                  </button>
                </div>
                <div class="gmd-sems-table">
                  <div class="gmd-sems-head">
                    <span>${t('fields.startDate')}</span><span>${t('fields.endDate')}</span><span></span>
                  </div>
                  ${selSems.map(gs=>`
                    <div class="gmd-sems-row" data-gsid="${gs.id}">
                      <span>${(gs.start_date||'').slice(0,10)}</span>
                      <span>${(gs.end_date||'').slice(0,10)}</span>
                      <button class="btn btn-sm btn-danger gmd-del-sem" data-gsid="${gs.id}"><svg class="icon icon-sm"><use href="#icon-trash"/></svg></button>
                    </div>`).join('')}
                  ${selSems.length===0?`<div class="gmd-empty" style="padding:12px 0">${t('messages.noData')}</div>`:''}
                </div>
              ` : `<div class="gmd-empty">${t('messages.noData')}</div>`}
            </div>
          </div>
        </div>`;

      // wire events
      container.querySelector('#gmd-spec')?.addEventListener('change', e => { filterSpec=e.target.value; redraw(); });
      container.querySelectorAll('.gmd-item').forEach(el => {
        el.addEventListener('click', () => { selectedId=Number(el.dataset.gid); redraw(); });
      });

      // Add group
      container.querySelector('#gmd-add')?.addEventListener('click', () => {
        this._openGroupModal(null, specialties, async (saved) => {
          if (this.state.yearId) saved.academic_year_id = Number(this.state.yearId);
          const created = await Api.post(ENTITY_CONFIGS['groups'].apiPath, [saved]);
          groups = groups.concat(created||[]);
          this.cache['groups'] = groups;
          if(created?.[0]) selectedId = created[0].id;
          redraw();
        });
      });

      // Edit group
      container.querySelector('#gmd-edit-grp')?.addEventListener('click', () => {
        if(!sel) return;
        this._openGroupModal(sel, specialties, async (saved) => {
          await Api.put(ENTITY_CONFIGS['groups'].apiPath, [{ ...sel, ...saved }]);
          const idx = groups.findIndex(g=>g.id===sel.id);
          if(idx>=0) groups[idx] = { ...sel, ...saved };
          this.cache['groups'] = groups;
          redraw();
        });
      });

      // Delete group
      container.querySelector('#gmd-del-grp')?.addEventListener('click', async () => {
        if(!sel || !await confirmDialog(t('messages.confirmDelete'))) return;
        try {
          await Api.post(ENTITY_CONFIGS['groups'].apiPath+'/delete', [sel.id]);
          groups = groups.filter(g=>g.id!==sel.id);
          groupSemesters = groupSemesters.filter(gs=>gs.group_id!==sel.id);
          this.cache['groups'] = groups;
          selectedId = groups[0]?.id || null;
          redraw();
        } catch(e) { Toast.error(e.message); }
      });

      // Add semester
      container.querySelector('#gmd-add-sem')?.addEventListener('click', () => {
        if(!sel) return;
        this._openGroupSemesterModal(sel.id, null, async (saved) => {
          const created = await Api.post(ENTITY_CONFIGS['group-semesters'].apiPath, [saved]);
          groupSemesters = groupSemesters.concat(created||[]);
          redraw();
        });
      });

      // Delete semester
      container.querySelectorAll('.gmd-del-sem').forEach(btn => {
        btn.addEventListener('click', async (e) => {
          e.stopPropagation();
          const gsid = Number(btn.dataset.gsid);
          if(!await confirmDialog(t('messages.confirmDelete'))) return;
          try {
            await Api.post(ENTITY_CONFIGS['group-semesters'].apiPath+'/delete', [gsid]);
            groupSemesters = groupSemesters.filter(gs=>gs.id!==gsid);
            redraw();
          } catch(e) { Toast.error(e.message); }
        });
      });
    };

    redraw();
  },

  _openGroupModal(group, specialties, onSave) {
    const isEdit = !!group;
    const body = document.createElement('div');
    body.innerHTML = ENTITY_CONFIGS['groups'].columns.filter(c=>!c.readonly).map(col => {
      const val = group?.[col.key] ?? '';
      let inp = '';
      if(col.type==='select' && col.ref==='specialties') {
        inp = `<select data-field="${col.key}">
          <option value="">—</option>
          ${specialties.map(s=>`<option value="${s.id}"${String(s.id)===String(val)?' selected':''}>${s.name}</option>`).join('')}
        </select>`;
      } else if(col.type==='bool') {
        inp = `<select data-field="${col.key}">
          <option value="false"${!val?' selected':''}>✗</option>
          <option value="true"${val?' selected':''}>✓</option>
        </select>`;
      } else {
        inp = `<input type="${col.type==='date'?'date':'text'}" data-field="${col.key}" value="${col.type==='date'?String(val).slice(0,10):val}">`;
      }
      return `<div class="form-field"><label>${t(col.labelKey)}${col.required?' <span class="req">*</span>':''}</label>${inp}</div>`;
    }).join('');
    const foot = document.createElement('div');
    foot.innerHTML = `<button class="btn btn-secondary" onclick="closeModal()">${t('actions.cancel')}</button>
      <button class="btn btn-primary" id="grp-save">${t('actions.save')}</button>`;
    openModal(isEdit ? t('actions.edit') : t('actions.add'), body, foot);
    foot.querySelector('#grp-save').addEventListener('click', async () => {
      const saved = {};
      body.querySelectorAll('[data-field]').forEach(el => {
        const col = ENTITY_CONFIGS['groups'].columns.find(c=>c.key===el.dataset.field);
        let v = el.value;
        if(col?.type==='bool') v = v==='true';
        if(col?.type==='number') v = v===''?null:Number(v);
        if(col?.type==='select'&&col?.ref) v = v===''?null:Number(v);
        if(col?.type==='date'&&v&&v.length===10) v=v+'T00:00:00Z';
        saved[el.dataset.field] = v;
      });
      try {
        foot.querySelector('#grp-save').disabled = true;
        await onSave(saved);
        closeModal();
        Toast.success(t('messages.saved'));
      } catch(e) { Toast.error(e.message); foot.querySelector('#grp-save').disabled=false; }
    });
  },

  _openGroupSemesterModal(groupId, gs, onSave) {
    const body = document.createElement('div');
    body.innerHTML = `
      <div class="form-field"><label>${t('fields.startDate')} <span class="req">*</span></label>
        <input type="date" id="gs-start" value="${(gs?.start_date||'').slice(0,10)}"></div>
      <div class="form-field"><label>${t('fields.endDate')} <span class="req">*</span></label>
        <input type="date" id="gs-end" value="${(gs?.end_date||'').slice(0,10)}"></div>`;
    const foot = document.createElement('div');
    foot.innerHTML = `<button class="btn btn-secondary" onclick="closeModal()">${t('actions.cancel')}</button>
      <button class="btn btn-primary" id="gs-save">${t('actions.save')}</button>`;
    openModal(t('actions.add'), body, foot);
    foot.querySelector('#gs-save').addEventListener('click', async () => {
      const start = document.getElementById('gs-start').value;
      const end = document.getElementById('gs-end').value;
      if(!start||!end) { Toast.warning(t('messages.fillRequired')); return; }
      const payload = { group_id: groupId, start_date: start+'T00:00:00Z', end_date: end+'T00:00:00Z' };
      try {
        foot.querySelector('#gs-save').disabled = true;
        await onSave(payload);
        closeModal();
        Toast.success(t('messages.saved'));
      } catch(e) { Toast.error(e.message); foot.querySelector('#gs-save').disabled=false; }
    });
  },

  // ═══════════════════════════════════════════════════════════
  // WORKLOAD MASTER-DETAIL
  // ═══════════════════════════════════════════════════════════
  async renderWorkloadMasterDetail(container) {
    this.currentMgr = null;
    await this.ensureCache('study-plans');
    await this.ensureCache('groups');
    await this.ensureCache('specialties');
    await this.ensureCache('disciplines');
    await this.ensureCache('teachers');
    await this.ensureTeacherUsersCache();

    const studyPlans  = this.cache['study-plans']   || [];
    const groups      = this.cache['groups']        || [];
    const specialties = this.cache['specialties']   || [];
    const disciplines = this.cache['disciplines']   || [];
    const teachers    = this.cache['teachers']      || [];
    const tUsers      = this.cache['teacher-users'] || [];

    let distributions = [];
    let assignments   = [];
    try {
      distributions = await Api.get(ENTITY_CONFIGS['workload-distributions'].apiPath) || [];
      assignments   = await Api.get(ENTITY_CONFIGS['workload-assignments'].apiPath)   || [];
    } catch(e) {
      container.innerHTML = `<div class="alert alert-error">${e.message}</div>`;
      return;
    }

    let selectedId = distributions.length > 0 ? distributions[0].id : null;

    const distTitle = (d) => {
      const plan = studyPlans.find(p => p.id === d.study_plan_id);
      const disc = plan ? disciplines.find(x => x.id === plan.discipline_id) : null;
      return disc?.name || disc?.short_name || `#${d.study_plan_id}`;
    };
    const distSub = (d) => {
      const plan = studyPlans.find(p => p.id === d.study_plan_id);
      const spec = plan ? specialties.find(x => x.id === plan.specialty_id) : null;
      const grp  = groups.find(g => g.id === d.group_id);
      return `${spec?.name||'—'} · ${grp?.name||'—'}`;
    };
    const teacherName = (teacherId) => {
      const tc   = teachers.find(x => x.id === teacherId);
      const user = tc ? tUsers.find(u => u.id === tc.user_id) : null;
      if (!user) return `#${teacherId||'?'}`;
      return (`${user.first_name||''} ${user.last_name||''}`).trim() || user.email;
    };

    const redraw = () => {
      const sel     = distributions.find(d => d.id === selectedId) || null;
      const selAssigns = assignments.filter(a => a.workload_distribution_id === selectedId);

      container.innerHTML = `
        <div class="gmd-wrap">
          <div class="gmd-toolbar">
            <button class="btn btn-primary btn-sm" id="wmd-add">
              <svg class="icon icon-sm"><use href="#icon-plus"/></svg>${t('actions.add')}
            </button>
          </div>
          <div class="gmd-layout">
            <div class="gmd-list-panel">
              <div class="gmd-list-hdr">${t('entities.workloadDistributions')} (${distributions.length})</div>
              ${distributions.length === 0 ? `<div class="gmd-empty">${t('messages.noData')}</div>` :
                distributions.map(d => {
                  const isSel   = d.id === selectedId;
                  const aCount  = assignments.filter(a => a.workload_distribution_id === d.id).length;
                  return `<div class="gmd-item${isSel?' gmd-item-sel':''}" data-did="${d.id}">
                    <div class="gmd-item-row">
                      <span class="gmd-item-name">${distTitle(d)}</span>
                      <span class="gmd-item-meta">${aCount} ${t('entities.workloadAssignments').toLowerCase()}</span>
                    </div>
                    <div class="gmd-item-sub">${distSub(d)}</div>
                  </div>`;
                }).join('')}
            </div>
            <div class="gmd-detail-panel">
              ${sel ? `
                <div class="gmd-detail-header">
                  <div>
                    <h3 class="gmd-detail-name">${distTitle(sel)}</h3>
                    <div class="gmd-detail-sub">${distSub(sel)}</div>
                  </div>
                  <div style="display:flex;gap:6px">
                    <button class="btn btn-sm btn-secondary" id="wmd-edit-dist">${t('actions.edit')}</button>
                    <button class="btn btn-sm btn-danger" id="wmd-del-dist"><svg class="icon icon-sm"><use href="#icon-trash"/></svg></button>
                  </div>
                </div>
                <div class="gmd-props">
                  <div class="gmd-prop"><span class="gmd-prop-lbl">${t('fields.classroomWork')}</span><span>${sel.classroom_work??'—'}</span></div>
                  <div class="gmd-prop"><span class="gmd-prop-lbl">${t('fields.laboratoryHours')}</span><span>${sel.laboratory??'—'}</span></div>
                  <div class="gmd-prop"><span class="gmd-prop-lbl">${t('fields.practicalHours')}</span><span>${sel.practical??'—'}</span></div>
                  <div class="gmd-prop"><span class="gmd-prop-lbl">${t('fields.examHours')}</span><span>${sel.exam??'—'}</span></div>
                </div>
                <div class="gmd-sems-header">
                  <span class="gmd-sems-title">${t('entities.workloadAssignments')}</span>
                  <button class="btn btn-sm btn-secondary" id="wmd-add-assign">
                    <svg class="icon icon-sm"><use href="#icon-plus"/></svg>${t('actions.add')}
                  </button>
                </div>
                <div class="gmd-sems-table">
                  <div style="display:grid;grid-template-columns:1fr 1fr 80px 40px;padding:8px 12px;background:var(--surface-2);font-size:.72rem;font-weight:600;text-transform:uppercase;letter-spacing:.04em;color:var(--text-muted);border-bottom:1px solid var(--border)">
                    <span>${t('fields.teacher')}</span><span>${t('fields.roleType')}</span><span>${t('fields.assignedHours')}</span><span></span>
                  </div>
                  ${selAssigns.map(a => `
                    <div class="gmd-sems-row" style="grid-template-columns:1fr 1fr 80px 40px">
                      <span>${teacherName(a.teacher_id)}</span>
                      <span>${t('roleTypes.' + a.role_type) !== 'roleTypes.' + a.role_type ? t('roleTypes.' + a.role_type) : (a.role_type||'—')}</span>
                      <span>${a.assigned_hours??'—'}</span>
                      <button class="btn btn-sm btn-danger wmd-del-assign" data-aid="${a.id}"><svg class="icon icon-sm"><use href="#icon-trash"/></svg></button>
                    </div>`).join('')}
                  ${selAssigns.length===0?`<div class="gmd-empty" style="padding:12px 0">${t('messages.noData')}</div>`:''}
                </div>
              ` : `<div class="gmd-empty">${t('messages.noData')}</div>`}
            </div>
          </div>
        </div>`;

      container.querySelectorAll('.gmd-item').forEach(el => {
        el.addEventListener('click', () => { selectedId = Number(el.dataset.did); redraw(); });
      });

      container.querySelector('#wmd-add')?.addEventListener('click', () => {
        this._openWorkloadDistModal(null, studyPlans, groups, specialties, disciplines, async (saved) => {
          const created = await Api.post(ENTITY_CONFIGS['workload-distributions'].apiPath, [saved]);
          distributions = distributions.concat(created || []);
          if (created?.[0]) selectedId = created[0].id;
          redraw();
        });
      });

      container.querySelector('#wmd-edit-dist')?.addEventListener('click', () => {
        if (!sel) return;
        this._openWorkloadDistModal(sel, studyPlans, groups, specialties, disciplines, async (saved) => {
          await Api.put(ENTITY_CONFIGS['workload-distributions'].apiPath, [{ ...sel, ...saved }]);
          const idx = distributions.findIndex(d => d.id === sel.id);
          if (idx >= 0) distributions[idx] = { ...sel, ...saved };
          redraw();
        });
      });

      container.querySelector('#wmd-del-dist')?.addEventListener('click', async () => {
        if (!sel || !await confirmDialog(t('messages.confirmDelete'))) return;
        try {
          await Api.post(ENTITY_CONFIGS['workload-distributions'].apiPath + '/delete', [sel.id]);
          distributions = distributions.filter(d => d.id !== sel.id);
          assignments   = assignments.filter(a => a.workload_distribution_id !== sel.id);
          selectedId    = distributions[0]?.id || null;
          redraw();
        } catch(e) { Toast.error(e.message); }
      });

      container.querySelector('#wmd-add-assign')?.addEventListener('click', () => {
        if (!sel) return;
        this._openWorkloadAssignModal(sel.id, null, teachers, tUsers, async (saved) => {
          const created = await Api.post(ENTITY_CONFIGS['workload-assignments'].apiPath, [saved]);
          assignments = assignments.concat(created || []);
          redraw();
        });
      });

      container.querySelectorAll('.wmd-del-assign').forEach(btn => {
        btn.addEventListener('click', async (e) => {
          e.stopPropagation();
          const aid = Number(btn.dataset.aid);
          if (!await confirmDialog(t('messages.confirmDelete'))) return;
          try {
            await Api.post(ENTITY_CONFIGS['workload-assignments'].apiPath + '/delete', [aid]);
            assignments = assignments.filter(a => a.id !== aid);
            redraw();
          } catch(e) { Toast.error(e.message); }
        });
      });
    };

    redraw();
  },

  _openWorkloadDistModal(dist, studyPlans, groups, specialties, disciplines, onSave) {
    const isEdit = !!dist;
    const body = document.createElement('div');
    const studyPlanLabel = (p) => {
      const disc = disciplines.find(x => x.id === p.discipline_id);
      const spec = specialties.find(x => x.id === p.specialty_id);
      return `${disc?.name || '#' + p.id} (${spec?.short_name || spec?.name || '—'}, сем.${p.semester_number})`;
    };
    body.innerHTML = `
      <div class="form-field">
        <label>${t('fields.studyPlan')} <span class="req">*</span></label>
        <select data-field="study_plan_id">
          <option value="">—</option>
          ${studyPlans.map(p => `<option value="${p.id}"${dist && String(p.id)===String(dist.study_plan_id)?' selected':''}>${studyPlanLabel(p)}</option>`).join('')}
        </select>
      </div>
      <div class="form-field">
        <label>${t('fields.group')} <span class="req">*</span></label>
        <select data-field="group_id">
          <option value="">—</option>
          ${groups.map(g => `<option value="${g.id}"${dist && String(g.id)===String(dist.group_id)?' selected':''}>${g.name}</option>`).join('')}
        </select>
      </div>
      <div class="form-field"><label>${t('fields.classroomWork')}</label>
        <input type="number" data-field="classroom_work" value="${dist?.classroom_work ?? ''}"></div>
      <div class="form-field"><label>${t('fields.laboratoryHours')}</label>
        <input type="number" data-field="laboratory" value="${dist?.laboratory ?? ''}"></div>
      <div class="form-field"><label>${t('fields.practicalHours')}</label>
        <input type="number" data-field="practical" value="${dist?.practical ?? ''}"></div>
      <div class="form-field"><label>${t('fields.examHours')}</label>
        <input type="number" data-field="exam" value="${dist?.exam ?? ''}"></div>`;
    const foot = document.createElement('div');
    foot.innerHTML = `<button class="btn btn-secondary" onclick="closeModal()">${t('actions.cancel')}</button>
      <button class="btn btn-primary" id="wdist-save">${t('actions.save')}</button>`;
    openModal(isEdit ? t('actions.edit') : t('actions.add'), body, foot);
    foot.querySelector('#wdist-save').addEventListener('click', async () => {
      const get = (field) => body.querySelector(`[data-field="${field}"]`)?.value;
      const spId = get('study_plan_id');
      const grId = get('group_id');
      if (!spId || !grId) { Toast.warning(t('messages.fillRequired')); return; }
      const saved = {
        study_plan_id:  Number(spId),
        group_id:       Number(grId),
        classroom_work: get('classroom_work') !== '' ? Number(get('classroom_work')) : null,
        laboratory:     get('laboratory')     !== '' ? Number(get('laboratory'))     : null,
        practical:      get('practical')      !== '' ? Number(get('practical'))      : null,
        exam:           get('exam')           !== '' ? Number(get('exam'))           : null,
      };
      try {
        foot.querySelector('#wdist-save').disabled = true;
        await onSave(saved);
        closeModal();
        Toast.success(t('messages.saved'));
      } catch(e) { Toast.error(e.message); foot.querySelector('#wdist-save').disabled = false; }
    });
  },

  _openWorkloadAssignModal(distId, assign, teachers, tUsers, onSave) {
    const teacherLabel = (tc) => {
      const user = tUsers.find(u => u.id === tc.user_id);
      return user ? ((`${user.first_name||''} ${user.last_name||''}`).trim() || user.email) : `#${tc.id}`;
    };
    const body = document.createElement('div');
    body.innerHTML = `
      <div class="form-field">
        <label>${t('fields.teacher')} <span class="req">*</span></label>
        <select data-field="teacher_id">
          <option value="">—</option>
          ${teachers.map(tc => `<option value="${tc.id}"${assign && String(tc.id)===String(assign.teacher_id)?' selected':''}>${teacherLabel(tc)}</option>`).join('')}
        </select>
      </div>
      <div class="form-field">
        <label>${t('fields.roleType')}</label>
        <select data-field="role_type">
          <option value="">—</option>
          ${['lecture','practical','laboratory','exam','other'].map(v =>
            `<option value="${v}"${assign?.role_type===v?' selected':''}>${t('roleTypes.'+v)}</option>`
          ).join('')}
        </select>
      </div>
      <div class="form-field"><label>${t('fields.assignedHours')}</label>
        <input type="number" data-field="assigned_hours" value="${assign?.assigned_hours ?? ''}"></div>`;
    const foot = document.createElement('div');
    foot.innerHTML = `<button class="btn btn-secondary" onclick="closeModal()">${t('actions.cancel')}</button>
      <button class="btn btn-primary" id="wassign-save">${t('actions.save')}</button>`;
    openModal(t('actions.add'), body, foot);
    foot.querySelector('#wassign-save').addEventListener('click', async () => {
      const get = (f) => body.querySelector(`[data-field="${f}"]`)?.value;
      const tcId = get('teacher_id');
      if (!tcId) { Toast.warning(t('messages.fillRequired')); return; }
      const saved = {
        workload_distribution_id: distId,
        teacher_id:     Number(tcId),
        role_type:      get('role_type') || null,
        assigned_hours: get('assigned_hours') !== '' ? Number(get('assigned_hours')) : null,
      };
      try {
        foot.querySelector('#wassign-save').disabled = true;
        await onSave(saved);
        closeModal();
        Toast.success(t('messages.saved'));
      } catch(e) { Toast.error(e.message); foot.querySelector('#wassign-save').disabled = false; }
    });
  },

  async _openTransferToWorkloadModal(studyPlan) {
    await this.ensureCache('groups');
    const groups = this.cache['groups'] || [];
    const body = document.createElement('div');
    body.innerHTML = `
      <p style="color:var(--text-muted);margin-bottom:16px;font-size:.9rem">${t('workload.transfer')}: <strong>${
        (() => { const d = (this.cache['disciplines']||[]).find(x=>x.id===studyPlan.discipline_id); return d?.name||'#'+studyPlan.id; })()
      }</strong></p>
      <div class="form-field">
        <label>${t('workload.transferGroup')} <span class="req">*</span></label>
        <select id="transfer-group">
          <option value="">—</option>
          ${groups.map(g => `<option value="${g.id}">${g.name}</option>`).join('')}
        </select>
      </div>`;
    const foot = document.createElement('div');
    foot.innerHTML = `<button class="btn btn-secondary" onclick="closeModal()">${t('actions.cancel')}</button>
      <button class="btn btn-primary" id="transfer-save">${t('actions.save')}</button>`;
    openModal(t('workload.transfer'), body, foot);
    foot.querySelector('#transfer-save').addEventListener('click', async () => {
      const groupId = document.getElementById('transfer-group').value;
      if (!groupId) { Toast.warning(t('messages.fillRequired')); return; }
      const payload = [{
        study_plan_id:  studyPlan.id,
        group_id:       Number(groupId),
        classroom_work: studyPlan.lectures   || null,
        laboratory:     studyPlan.laboratory || null,
        practical:      studyPlan.practical  || null,
        exam:           studyPlan.exam       || null,
      }];
      try {
        foot.querySelector('#transfer-save').disabled = true;
        await Api.post(ENTITY_CONFIGS['workload-distributions'].apiPath, payload);
        closeModal();
        Toast.success(t('messages.saved'));
        location.hash = '#/processing/workload-distributions';
      } catch(e) { Toast.error(e.message); foot.querySelector('#transfer-save').disabled = false; }
    });
  },

  // ═══════════════════════════════════════════════════════════
  // USERS PAGE (admin only)
  // ═══════════════════════════════════════════════════════════
  async renderUsers(main) {
    if (this.state.user?.role !== 'admin') {
      main.innerHTML = `<div class="alert alert-error">${t('messages.forbidden') || 'Forbidden'}</div>`;
      return;
    }
    this.currentMgr = null;

    main.innerHTML = `
      <div class="page-header">
        <h1>${t('pages.usersTitle')}</h1>
        <button class="btn btn-primary" id="add-user-btn">
          <svg class="icon icon-sm"><use href="#icon-plus"/></svg>
          ${t('users.register')}
        </button>
      </div>
      <div class="loading-state"><div class="spinner"></div></div>`;

    main.querySelector('#add-user-btn').addEventListener('click', () => this._openRegisterModal());

    let users = [];
    try {
      users = await Api.get('/api/auth/users') || [];
    } catch (e) {
      main.querySelector('.loading-state').outerHTML = `<div class="alert alert-error">${e.message}</div>`;
      return;
    }

    const wrap = document.createElement('div');
    wrap.className = 'table-wrapper';
    main.querySelector('.loading-state')?.replaceWith(wrap);
    this._renderUsersTable(users, wrap);
  },

  _renderUsersTable(users, wrap) {
    if (!users.length) {
      wrap.innerHTML = `<div class="empty-state">—</div>`;
      return;
    }
    wrap.innerHTML = `
      <table class="entity-table">
        <thead>
          <tr>
            <th>${t('fields.email')}</th>
            <th>${t('fields.name')}</th>
            <th>${t('fields.role')}</th>
            <th>${t('fields.isActive')}</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          ${users.map(u => {
            const name = (`${u.first_name || ''} ${u.last_name || ''}`).trim() || `(${t('users.noName')})`;
            const roleBadge = `<span class="badge badge-${u.role==='admin'?'warning':'info'}">${t('roles.'+u.role)}</span>`;
            const statusBadge = u.is_active
              ? `<span class="badge badge-success">${t('status.active')}</span>`
              : `<span class="badge badge-neutral" style="opacity:.7">${t('users.pendingStatus')}</span>`;
            return `<tr data-uid="${u.id}">
              <td>${u.email}</td>
              <td>${name}</td>
              <td>${roleBadge}</td>
              <td>${statusBadge}</td>
              <td class="col-actions" style="width:auto;white-space:nowrap">
                ${!u.is_active ? `<button class="btn btn-sm btn-secondary" data-copy-invite="${u.id}" title="${t('users.inviteLink')}">
                  <svg class="icon icon-sm"><use href="#icon-refresh"/></svg>
                </button>` : ''}
                <button class="btn btn-sm btn-${u.is_active ? 'secondary' : 'primary'}" data-id="${u.id}" data-active="${u.is_active}">
                  ${u.is_active ? t('actions.deactivate') : t('actions.activate')}
                </button>
                <button class="btn btn-sm btn-danger" data-delete="${u.id}">
                  <svg class="icon icon-sm"><use href="#icon-trash"/></svg>
                </button>
              </td>
            </tr>`;
          }).join('')}
        </tbody>
      </table>`;

    wrap.querySelectorAll('[data-copy-invite]').forEach(btn => {
      btn.addEventListener('click', async () => {
        const id = btn.dataset.copyInvite;
        try {
          const resp = await Api.post(`/api/auth/users/${id}/reset-invite`, null);
          const link = window.location.origin + (resp?.invite_link || '');
          navigator.clipboard.writeText(link).then(() => Toast.info(t('messages.copied')));
        } catch(e) { Toast.error(e.message); }
      });
    });

    wrap.querySelectorAll('[data-id]').forEach(btn => {
      btn.addEventListener('click', async () => {
        const id     = btn.dataset.id;
        const active = btn.dataset.active === 'true';
        try {
          await Api.post(`/api/auth/users/${id}/${active ? 'deactivate' : 'activate'}`, null);
          Toast.success(t('messages.saved'));
          const updated = await Api.get('/api/auth/users');
          this._renderUsersTable(updated || [], wrap);
        } catch (e) { Toast.error(e.message); }
      });
    });

    wrap.querySelectorAll('[data-delete]').forEach(btn => {
      btn.addEventListener('click', async () => {
        const id = Number(btn.dataset.delete);
        if (!await confirmDialog(t('messages.confirmDelete'))) return;
        try {
          await Api.post('/api/auth/users/delete', [id]);
          Toast.success(t('messages.deleted'));
          const updated = await Api.get('/api/auth/users');
          this._renderUsersTable(updated || [], wrap);
        } catch (e) { Toast.error(e.message); }
      });
    });
  },

  _openRegisterModal() {
    const body = document.createElement('div');
    body.innerHTML = `
      <div class="form-field"><label>${t('users.email') || t('fields.email')} *</label>
        <input type="email" id="reg-email"></div>
      <div class="form-field"><label>${t('fields.firstName') || 'First name'} *</label>
        <input type="text" id="reg-first-name"></div>
      <div class="form-field"><label>${t('fields.lastName') || 'Last name'} *</label>
        <input type="text" id="reg-last-name"></div>
      <div class="form-field"><label>${t('fields.role')}</label>
        <select id="reg-role">
          <option value="teacher">${t('roles.teacher')}</option>
          <option value="user">${t('roles.user')}</option>
          <option value="dean">${t('roles.dean')}</option>
          <option value="admin">${t('roles.admin')}</option>
        </select></div>
      <div id="invite-link-section" style="display:none;margin-top:16px">
        <label style="font-weight:600;margin-bottom:6px;display:block">${t('users.inviteLink')}</label>
        <div style="display:flex;gap:8px;align-items:center">
          <input type="text" id="invite-link-input" readonly style="flex:1;background:var(--bg-secondary,#f8fafc)">
          <button class="btn btn-secondary btn-sm" id="copy-link-btn">${t('actions.copy')}</button>
        </div>
        <p style="color:var(--text-muted);font-size:.82rem;margin-top:6px">
          ${t('users.inviteLinkHint')}
        </p>
      </div>`;
    const foot = document.createElement('div');
    foot.innerHTML = `<button class="btn btn-secondary" id="reg-cancel">${t('actions.cancel')}</button>
      <button class="btn btn-primary" id="reg-save">${t('users.register')}</button>`;

    openModal(t('users.register'), body, foot);

    foot.querySelector('#reg-cancel').addEventListener('click', () => closeModal());

    foot.querySelector('#reg-save').addEventListener('click', async () => {
      const email     = document.getElementById('reg-email').value.trim();
      const firstName = document.getElementById('reg-first-name').value.trim();
      const lastName  = document.getElementById('reg-last-name').value.trim();
      const role      = document.getElementById('reg-role').value;
      if (!email || !firstName || !lastName) {
        Toast.warning(t('messages.fillRequired') || 'Fill required fields');
        return;
      }
      try {
        const btn = foot.querySelector('#reg-save');
        btn.disabled = true;
        const resp = await Api.post('/api/auth/auth/invite', { email, first_name: firstName, last_name: lastName, role });
        btn.disabled = false;

        const inviteLink = window.location.origin + (resp?.invite_link || '');
        const linkInput = document.getElementById('invite-link-input');
        linkInput.value = inviteLink;
        document.getElementById('invite-link-section').style.display = 'block';

        document.getElementById('copy-link-btn').addEventListener('click', () => {
          navigator.clipboard.writeText(inviteLink).then(() => Toast.info(t('messages.copied')));
        });

        foot.querySelector('#reg-save').style.display = 'none';

        Toast.success(t('messages.saved'));
        const updated = await Api.get('/api/auth/users');
        const wrap = document.querySelector('.table-wrapper');
        if (wrap) this._renderUsersTable(updated || [], wrap);
      } catch (e) {
        Toast.error(e.message);
      }
    });
  },
};

// ═══════════════════════════════════════════════════════════════
// SIDEBAR TOGGLE
// ═══════════════════════════════════════════════════════════════
document.addEventListener('DOMContentLoaded', () => {
  document.getElementById('sidebar-toggle')?.addEventListener('click', () => {
    document.getElementById('sidebar').classList.toggle('sidebar-collapsed');
    document.getElementById('app-shell').classList.toggle('sidebar-collapsed');
  });

  // Login form
  document.getElementById('login-btn')?.addEventListener('click', async () => {
    const email = document.getElementById('login-email').value.trim();
    const pass  = document.getElementById('login-password').value;
    const errEl = document.getElementById('login-error');
    errEl.classList.add('hidden');
    try {
      const data = await Api.login(email, pass);
      if (!data?.token) throw new Error('No token');
      App.setAuth(data.token);
      App.showApp();
      await App.loadYears();
      await App.checkSetup();
      App.setupHashRouter();
      location.hash = '#/home';
    } catch (e) {
      errEl.textContent = e.message === 'unauthorized' || e.message.includes('401')
        ? t('auth.loginError') : e.message;
      errEl.classList.remove('hidden');
    }
  });

  document.getElementById('login-password')?.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') document.getElementById('login-btn')?.click();
  });

  // Sidebar language buttons
  document.querySelectorAll('.lang-btn[data-lang]').forEach(btn => {
    btn.addEventListener('click', () => setLang(btn.getAttribute('data-lang')));
  });

  // Sidebar logout
  document.getElementById('sidebar-logout-btn')?.addEventListener('click', () => {
    App.clearAuth();
    App.showLogin();
    location.hash = '';
  });

  window.addEventListener('langchange', () => {
    // Re-render current page on lang switch
    App.route(location.hash || '#/home');
  });

  App.init();
});
