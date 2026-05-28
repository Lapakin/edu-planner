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
  body.removeAttribute('style');
  if (bodyEl instanceof HTMLElement) body.appendChild(bodyEl); else body.innerHTML = bodyEl || '';
  if (footerEl instanceof HTMLElement) foot.appendChild(footerEl); else foot.innerHTML = footerEl || '';
  const modal = document.getElementById('modal');
  modal.classList.remove('modal-fullscreen', 'modal-lg', 'modal-sm');
  document.getElementById('modal-overlay').classList.remove('hidden');
}

function closeModal() {
  document.getElementById('modal-overlay').classList.add('hidden');
  document.getElementById('modal').classList.remove('modal-fullscreen', 'modal-lg', 'modal-sm');
  document.getElementById('modal-body').removeAttribute('style');
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
    const role = this.state.user?.role;
    // admin-only nav items (users page)
    document.querySelectorAll('.admin-only').forEach(el => {
      el.classList.toggle('hidden', role !== 'admin' && role !== 'dean');
    });
    // teacher-only: show only home and schedule
    const teacherAllowed = ['home', 'profile'];
    document.querySelectorAll('.nav-item[data-page]').forEach(el => {
      if (role === 'teacher') {
        el.classList.toggle('hidden', !teacherAllowed.includes(el.dataset.page));
      } else {
        el.classList.remove('hidden');
      }
    });
    // Set profile button initial letter
    const email   = this.state.user?.email || '';
    const initial = (email[0] || 'U').toUpperCase();
    const profileBtn = document.getElementById('profile-btn');
    if (profileBtn) profileBtn.textContent = initial;
    // Set role label
    const headerRole = document.getElementById('header-role');
    if (headerRole) {
      headerRole.textContent = t('roles.' + (role || 'user'));
      headerRole.className = 'badge ' + (role === 'admin' ? 'badge-warning' : 'badge-info');
    }
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

    const role = this.state.user?.role;
    const teacherAllowed = ['home', 'profile'];
    if (role === 'teacher' && !teacherAllowed.includes(page)) {
      location.hash = '#/home';
      return;
    }
    if (page === 'users' && role !== 'admin' && role !== 'dean') {
      location.hash = '#/home';
      return;
    }

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
        const rows = await Api.get(apiPath) || [];
        if (key === 'users') {
          rows.forEach(u => {
            u.full_name = [u.first_name, u.last_name].filter(Boolean).join(' ') || u.email;
          });
        }
        this.cache[key] = rows;
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
    if (!activating) return;
    const path = `/api/syllabus/academic-years/${row.id}/activate`;
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

  async activateEntity(entityKey, row, mgr) {
    const cfg = ENTITY_CONFIGS[entityKey];
    if (!cfg) return;
    try {
      const payload = [{ ...row, is_active: true }];
      await Api.put(cfg.apiPath, payload);
      row.is_active = true;
      mgr.reloadRow(row);
      Toast.success(t('messages.saved'));
    } catch (e) { Toast.error(e.message); }
  },

  async resetUserInvite(row) {
    try {
      const resp = await Api.post(`/api/auth/users/${row.id}/reset-invite`, null);
      const inviteLink = window.location.origin + (resp?.invite_link || '');
      const body = document.createElement('div');
      body.innerHTML = `
        <p style="font-size:.88rem;color:var(--text-muted);margin-bottom:12px">${t('users.inviteLinkHint')}</p>
        <div style="display:flex;gap:8px;align-items:center">
          <input type="text" id="reset-link-input" readonly value="${inviteLink}" style="flex:1;background:var(--bg-secondary,#f8fafc)">
          <button class="btn btn-secondary btn-sm" id="reset-copy-btn">${t('actions.copy')}</button>
        </div>`;
      const foot = document.createElement('div');
      foot.innerHTML = `<button class="btn btn-primary" onclick="closeModal()">${t('actions.close')}</button>`;
      openModal(t('users.inviteLink').replace(':', ''), body, foot);
      body.querySelector('#reset-copy-btn').addEventListener('click', () => {
        navigator.clipboard.writeText(inviteLink).then(() => Toast.info(t('messages.copied')));
      });
    } catch (e) { Toast.error(e.message); }
  },

  async resetUserPassword(row) {
    try {
      const resp = await Api.post(`/api/auth/users/${row.id}/generate-reset-password-link`, null);
      const inviteLink = window.location.origin + (resp?.invite_link || '');
      const body = document.createElement('div');
      body.innerHTML = `
        <p style="font-size:.88rem;color:var(--text-muted);margin-bottom:12px">${t('users.inviteLinkHint')}</p>
        <div style="display:flex;gap:8px;align-items:center">
          <input type="text" id="reset-link-input" readonly value="${inviteLink}" style="flex:1;background:var(--bg-secondary,#f8fafc)">
          <button class="btn btn-secondary btn-sm" id="reset-copy-btn">${t('actions.copy')}</button>
        </div>`;
      const foot = document.createElement('div');
      foot.innerHTML = `<button class="btn btn-primary" onclick="closeModal()">${t('actions.close')}</button>`;
      openModal(t('users.inviteLink').replace(':', ''), body, foot);
      body.querySelector('#reset-copy-btn').addEventListener('click', () => {
        navigator.clipboard.writeText(inviteLink).then(() => Toast.info(t('messages.copied')));
      });
    } catch (e) { Toast.error(e.message); }
  },

  async toggleUser(row, mgr) {
    const activating = !row.is_active;
    if (!activating) {
      const name = [row.first_name, row.last_name].filter(Boolean).join(' ') || row.email;
      const ok = await confirmDialog(t('users.deactivateConfirm').replace('{name}', name));
      if (!ok) return;
    }
    const path = `/api/auth/users/${row.id}/${activating ? 'activate' : 'deactivate'}`;
    try {
      await Api.post(path, null);
      row.is_active = activating;
      mgr.reloadRow(row);
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
      <div class="home-welcome">
        <div class="welcome-card">
          <div class="welcome-icon">
            <svg style="width:22px;height:22px"><use href="#icon-book"/></svg>
          </div>
          <div>
            <h2>${t('pages.homeWelcome')}, ${user?.email || ''}!</h2>
            <p>${t('roles.' + (user?.role || 'user'))}</p>
          </div>
        </div>
      </div>
      <div id="home-active-year"></div>
      <div id="home-active-tpl"></div>
      <div class="home-cards" id="home-stats"></div>`;

    this._renderActiveYearInfo(main.querySelector('#home-active-year'));
    this._renderActiveTemplate(main.querySelector('#home-active-tpl'));

    if (user?.role === 'teacher') return;

    const cards = [
      { i18nKey: 'entities.academicYears',    cacheKey: 'academic-years',    href: '#/academic-settings/academic-years', icon: 'icon-calendar', color: 'blue'    },
      { i18nKey: 'entities.departments',      cacheKey: 'departments',        href: '#/reference/departments',            icon: 'icon-book',     color: 'indigo'  },
      { i18nKey: 'entities.groups',           cacheKey: 'groups',             href: '#/reference/groups',                icon: 'icon-settings', color: 'sky'     },
      { i18nKey: 'entities.disciplines',      cacheKey: 'disciplines',        href: '#/reference/disciplines',           icon: 'icon-book',     color: 'emerald' },
      { i18nKey: 'entities.scheduleTemplates',cacheKey: 'schedule-templates', href: '#/schedule/schedule-templates',     icon: 'icon-calendar', color: 'slate'   },
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
        <div class="stat-icon-wrap ${card.color}">
          <svg class="stat-icon"><use href="#${card.icon}"/></svg>
        </div>
        <div class="stat-count">${count}</div>
        <div class="stat-label">${t(card.i18nKey)}</div>`;
      statsEl.appendChild(div);
    });
  },

  _renderActiveYearInfo(container) {
    const years = this.cache['academic-years'] || [];
    const activeYear = years.find(y => y.is_active);
    if (activeYear) {
      const startDate = activeYear.start_date?.slice(0, 10) || '';
      const endDate   = activeYear.end_date?.slice(0, 10)   || '';
      container.innerHTML = `
        <div class="card" style="margin:16px 0 12px;padding:12px 16px">
          <div style="display:flex;align-items:center;gap:10px">
            <svg class="icon" style="color:var(--success);flex-shrink:0"><use href="#icon-calendar"/></svg>
            <span style="font-weight:600">${t('home.activeYear')}:</span>
            <span>${activeYear.name}</span>
            <span style="color:var(--text-muted);font-size:.85rem">${startDate} – ${endDate}</span>
          </div>
        </div>`;
    } else {
      container.innerHTML = `
        <div class="card" style="margin:16px 0 12px;padding:12px 16px;border-color:var(--warning)">
          <div style="display:flex;align-items:center;gap:10px">
            <svg class="icon" style="color:var(--warning);flex-shrink:0"><use href="#icon-calendar"/></svg>
            <span style="color:var(--warning)">${t('home.noActiveYear')}</span>
            <a href="#/academic-settings/academic-years" style="color:var(--accent)">${t('home.createYear')}</a>
          </div>
        </div>`;
    }
  },

  async _renderActiveTemplate(container) {
    try {
      const templates = await Api.get('/api/syllabus/schedule-templates');
      const active = (templates || []).find(tmpl => tmpl.is_active);
      if (!active) {
        container.innerHTML = `<div class="home-no-template">
          <div class="home-no-template-icon">
            <svg style="width:20px;height:20px;color:var(--accent)"><use href="#icon-calendar"/></svg>
          </div>
          <p>${t('home.noActiveTemplate')}</p>
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
    const userId = jwtUser?.user_id;
    main.innerHTML = '<div class="loading-state"><div class="spinner"></div></div>';

    let user = {};
    try {
      if (userId) user = await Api.get(`/api/auth/users/${userId}`) || {};
    } catch { /* ignore */ }
    if (!user.id) {
      user = { id: userId, email: jwtUser?.email || '', role: jwtUser?.role || 'user', ...user };
    }

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
            id: u.id,
            email: u.email,
            role: u.role,
            is_active: u.is_active ?? true,
            is_deleted: u.is_deleted ?? false,
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
      (users || []).forEach(u => {
        u.full_name = [u.first_name, u.last_name].filter(Boolean).join(' ') || u.email;
      });
      this.cache['users'] = users || [];
    } catch { this.cache['users'] = []; }
  },

  async ensureTeacherUsersCache() {
    if (this.cache['teacher-users']) return;
    try {
      const users = await Api.get('/api/auth/users?role=teacher');
      (users || []).forEach(u => {
        u.full_name = [u.first_name, u.last_name].filter(Boolean).join(' ') || u.email;
      });
      this.cache['teacher-users'] = users || [];
    } catch { this.cache['teacher-users'] = []; }
  },

  // ═══════════════════════════════════════════════════════════
  // GENERATE SCHEDULE TAB
  // ═══════════════════════════════════════════════════════════
  async renderGenerateSchedule(container) {
    this.currentMgr = null;
    await this.ensureCache('semesters');

    const semesters = this.cache['semesters'] || [];

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
      this._openTemplateSettingsModal();
    });

    container.querySelector('#gen-restrict-btn').addEventListener('click', () => {
      this._openScheduleRestrictionsModal();
    });

    container.querySelector('#gen-btn').addEventListener('click', async () => {
      const semId = container.querySelector('#gen-semester').value;
      if (!semId) { Toast.warning(t('generate.noSemester')); return; }
      await this._doGenerate(semId, container.querySelector('#gen-main'));
    });
  },

  // ── Template Settings Modal ─────────────────────────────────
  async _openTemplateSettingsModal() {
    const yearId = this.state.yearId;
    if (!yearId) { Toast.warning(t('messages.noAcademicYear')); return; }

    const bodyEl = document.createElement('div');
    bodyEl.style.cssText = 'display:flex;flex-direction:column;flex:1;min-height:0;overflow:hidden;';
    bodyEl.innerHTML = '<div class="loading-state" style="padding:20px"><div class="spinner"></div></div>';
    const footEl = document.createElement('div');
    footEl.style.cssText = 'display:flex;align-items:center;justify-content:space-between;width:100%;';

    openModal(t('entities.scheduleTemplateSettings'), bodyEl, footEl);
    document.getElementById('modal').classList.add('modal-lg');
    const mb = document.getElementById('modal-body');
    mb.style.padding = '0';
    mb.style.display = 'flex';
    mb.style.flexDirection = 'column';
    mb.style.overflow = 'hidden';

    let setting;
    try {
      const list = await Api.get(`/api/syllabus/schedule-template-settings?academic_year_id=${yearId}`);
      setting = (list && list.length > 0) ? { ...list[list.length - 1] } : {
        id: 0, academic_year_id: Number(yearId),
        lessons_per_class: 2, study_days_mask: 31,
        max_identical_lessons_per_day: 2, max_study_hours_per_day: 8,
        max_teacher_hours_per_week: 20, max_group_lesson_hours_per_week: 36,
      };
    } catch (e) {
      bodyEl.innerHTML = `<div class="alert alert-error" style="margin:20px">${e.message}</div>`;
      return;
    }

    const DAYS = [
      { label: t('weekdaysShort.monday'),    bit: 1 },
      { label: t('weekdaysShort.tuesday'),   bit: 2 },
      { label: t('weekdaysShort.wednesday'), bit: 4 },
      { label: t('weekdaysShort.thursday'),  bit: 8 },
      { label: t('weekdaysShort.friday'),    bit: 16 },
      { label: t('weekdaysShort.saturday'),  bit: 32 },
      { label: t('weekdaysShort.sunday'),    bit: 64, disabled: true },
    ];
    const dayButtonsHTML = DAYS.map(d => {
      const isActive = (setting.study_days_mask & d.bit) !== 0;
      const cls = ['sm-day', isActive ? 'active' : '', d.disabled ? 'sm-day-disabled' : ''].filter(Boolean).join(' ');
      return `<div class="${cls}" data-bit="${d.bit}">${d.label}</div>`;
    }).join('');

    bodyEl.innerHTML = `
      <div class="sm-tabs">
        <button class="sm-tab active" data-panel="format">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg>
          ${t('modal.tabFormat')} <span class="sm-tab-badge">2</span>
        </button>
        <button class="sm-tab" data-panel="calendar">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>
          ${t('modal.tabCalendar')} <span class="sm-tab-badge">1</span>
        </button>
        <button class="sm-tab" data-panel="daily">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M3 3v18h18"/><path d="M7 14l4-4 4 4 5-5"/></svg>
          ${t('modal.tabDaily')} <span class="sm-tab-badge">2</span>
        </button>
        <button class="sm-tab" data-panel="weekly">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="8.5" cy="7" r="4"/><path d="M20 8v6M23 11h-6"/></svg>
          ${t('modal.tabWeekly')} <span class="sm-tab-badge">2</span>
        </button>
      </div>
      <div class="sm-panels">
        <div class="sm-panel active" data-panel="format">
          <div class="sm-constraint">
            <div class="sm-constraint-info">
              <div class="sm-constraint-label">
                <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg></div>
                ${t('modal.lessonDuration')}
              </div>
              <div class="sm-constraint-desc">${t('modal.lessonDurationDesc')}</div>
            </div>
            <div class="sm-control">
              <select class="sm-select" id="sms-lessons-per-class">
                <option value="1" ${setting.lessons_per_class === 1 ? 'selected' : ''}>${t('lessonsPerClass.one')}</option>
                <option value="2" ${setting.lessons_per_class === 2 ? 'selected' : ''}>${t('lessonsPerClass.two')}</option>
                <option value="3" ${setting.lessons_per_class === 3 ? 'selected' : ''}>${t('lessonsPerClass.three')}</option>
              </select>
            </div>
          </div>
        </div>

        <div class="sm-panel" data-panel="calendar">
          <div class="sm-constraint">
            <div class="sm-constraint-info">
              <div class="sm-constraint-label">
                <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg></div>
                ${t('modal.studyDays')}
              </div>
              <div class="sm-constraint-desc">${t('modal.studyDaysDesc')}</div>
            </div>
            <div class="sm-control">
              <div class="sm-day-strip" id="sms-day-strip">${dayButtonsHTML}</div>
            </div>
          </div>
        </div>

        <div class="sm-panel" data-panel="daily">
          <div class="sm-constraint">
            <div class="sm-constraint-info">
              <div class="sm-constraint-label">
                <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg></div>
                ${t('modal.maxIdentical')}
              </div>
              <div class="sm-constraint-desc">${t('modal.maxIdenticalDesc')}</div>
            </div>
            <div class="sm-control">
              <input type="number" class="sm-num-input" id="sms-max-identical" value="${setting.max_identical_lessons_per_day}" min="1" max="6">
              <span class="sm-unit-label">${t('units.lessons')}</span>
            </div>
          </div>
          <div class="sm-constraint">
            <div class="sm-constraint-info">
              <div class="sm-constraint-label">
                <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M3 3v18h18"/><path d="M7 14l4-4 4 4 5-5"/></svg></div>
                ${t('modal.maxHoursDay')}
              </div>
              <div class="sm-constraint-desc">${t('modal.maxHoursDayDesc')}</div>
            </div>
            <div class="sm-control">
              <input type="number" class="sm-num-input" id="sms-max-hours-day" value="${setting.max_study_hours_per_day}" min="2" max="12">
              <span class="sm-unit-label">${t('units.hours')}</span>
            </div>
          </div>
        </div>

        <div class="sm-panel" data-panel="weekly">
          <div class="sm-constraint">
            <div class="sm-constraint-info">
              <div class="sm-constraint-label">
                <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M20 21v-2a4 4 0 00-3-3.87M4 21v-2a4 4 0 013-3.87"/><circle cx="12" cy="7" r="4"/></svg></div>
                ${t('modal.maxTeacherWeek')}
              </div>
              <div class="sm-constraint-desc">${t('modal.maxTeacherWeekDesc')}</div>
            </div>
            <div class="sm-control">
              <input type="number" class="sm-num-input" id="sms-max-teacher-week" value="${setting.max_teacher_hours_per_week}" min="1" max="40">
              <span class="sm-unit-label">${t('units.hours')}</span>
            </div>
          </div>
          <div class="sm-constraint">
            <div class="sm-constraint-info">
              <div class="sm-constraint-label">
                <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M16 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="8.5" cy="7" r="4"/><path d="M20 8v6M23 11h-6"/></svg></div>
                ${t('modal.maxGroupWeek')}
              </div>
              <div class="sm-constraint-desc">${t('modal.maxGroupWeekDesc')}</div>
            </div>
            <div class="sm-control">
              <input type="number" class="sm-num-input" id="sms-max-group-week" value="${setting.max_group_lesson_hours_per_week ?? 36}" min="1" max="60">
              <span class="sm-unit-label">${t('units.hours')}</span>
            </div>
          </div>
        </div>
      </div>`;

    // Tab switching
    bodyEl.querySelectorAll('.sm-tab').forEach(tab => {
      tab.addEventListener('click', () => {
        bodyEl.querySelectorAll('.sm-tab').forEach(t => t.classList.remove('active'));
        bodyEl.querySelectorAll('.sm-panel').forEach(p => p.classList.remove('active'));
        tab.classList.add('active');
        bodyEl.querySelector(`.sm-panel[data-panel="${tab.dataset.panel}"]`).classList.add('active');
      });
    });

    // Day picker
    bodyEl.querySelectorAll('.sm-day:not(.sm-day-disabled)').forEach(day => {
      day.addEventListener('click', () => day.classList.toggle('active'));
    });

    // Footer
    const isNew = !setting.id || setting.id === 0;
    const statusDot = document.createElement('span');
    statusDot.className = isNew ? 'sm-status-dot' : 'sm-status-dot saved';
    const statusText = document.createElement('span');
    statusText.textContent = isNew ? t('modal.unsaved') : t('modal.saved');
    const statusDiv = document.createElement('div');
    statusDiv.className = 'sm-modal-status';
    statusDiv.append(statusDot, statusText);

    const cancelBtn = document.createElement('button');
    cancelBtn.className = 'btn btn-secondary';
    cancelBtn.textContent = t('actions.cancel');
    cancelBtn.addEventListener('click', closeModal);

    const saveBtn = document.createElement('button');
    saveBtn.className = 'btn btn-primary';
    saveBtn.innerHTML = `<svg class="icon"><use href="#icon-save"/></svg> ${t('actions.save')}`;

    const actions = document.createElement('div');
    actions.style.cssText = 'display:flex;gap:8px;';
    actions.append(cancelBtn, saveBtn);
    footEl.append(statusDiv, actions);

    saveBtn.addEventListener('click', async () => {
      saveBtn.disabled = true;
      try {
        let mask = 0;
        bodyEl.querySelectorAll('.sm-day:not(.sm-day-disabled)').forEach(d => {
          if (d.classList.contains('active')) mask |= Number(d.dataset.bit);
        });
        const payload = {
          ...setting,
          academic_year_id: Number(yearId),
          lessons_per_class: Number(bodyEl.querySelector('#sms-lessons-per-class').value),
          study_days_mask: mask,
          max_identical_lessons_per_day: Number(bodyEl.querySelector('#sms-max-identical').value),
          max_study_hours_per_day: Number(bodyEl.querySelector('#sms-max-hours-day').value),
          max_teacher_hours_per_week: Number(bodyEl.querySelector('#sms-max-teacher-week').value),
          max_group_lesson_hours_per_week: Number(bodyEl.querySelector('#sms-max-group-week').value),
        };
        if (setting.id && setting.id > 0) {
          await Api.put('/api/syllabus/schedule-template-settings', [payload]);
        } else {
          const created = await Api.post('/api/syllabus/schedule-template-settings', [payload]);
          if (created && created[0]) setting.id = created[0].id;
        }
        delete this.cache['schedule-template-settings'];
        Toast.success(t('messages.saved'));
        statusDot.classList.add('saved');
        statusText.textContent = t('modal.saved');
      } catch (e) {
        Toast.error(e.message);
      } finally {
        saveBtn.disabled = false;
      }
    });
  },

  // ── Schedule Restrictions Modal ──────────────────────────────
  // ── Schedule Restrictions Modal ──────────────────────────────
  async _openScheduleRestrictionsModal() {
    const yearId = this.state.yearId;
    if (!yearId) { Toast.warning(t('messages.noAcademicYear')); return; }

    const bodyEl = document.createElement('div');
    bodyEl.style.cssText = 'display:flex;flex-direction:column;flex:1;min-height:0;overflow:hidden;';
    bodyEl.innerHTML = '<div class="loading-state" style="padding:20px"><div class="spinner"></div></div>';
    const footEl = document.createElement('div');
    footEl.style.cssText = 'display:flex;align-items:center;justify-content:space-between;width:100%;';

    openModal(t('entities.scheduleRestrictions'), bodyEl, footEl);
    document.getElementById('modal').classList.add('modal-lg');
    const mb = document.getElementById('modal-body');
    mb.style.padding = '0';
    mb.style.display = 'flex';
    mb.style.flexDirection = 'column';
    mb.style.overflow = 'hidden';

    // Load all data in parallel
    let restriction, allPreferences, allLabRooms, cycleCommittees, rooms, teachers;
    try {
      const getCached = async (key, url) => {
        if (this.cache[key] && this.cache[key].length > 0) return this.cache[key];
        const data = await Api.get(url);
        this.cache[key] = data || [];
        return this.cache[key];
      };
      const [restrictList, prefs, labRooms, cc, rm, tch] = await Promise.all([
        Api.get(`/api/syllabus/schedule-restrictions?academic_year_id=${yearId}`),
        Api.get(`/api/syllabus/teacher-slot-preferences?academic_year_id=${yearId}`),
        Api.get(`/api/syllabus/cycle-committee-lab-rooms?academic_year_id=${yearId}`),
        getCached('cycle-committees', '/api/syllabus/cycle-committees'),
        getCached('rooms', '/api/syllabus/rooms'),
        Api.get(`/api/syllabus/teachers?academic_year_id=${yearId}`),
        this.ensureTeacherUsersCache(),
      ]);
      restriction = (restrictList && restrictList.length > 0) ? { ...restrictList[restrictList.length - 1] } : {
        id: 0, academic_year_id: Number(yearId),
        min_group_lessons_per_day: 2, max_group_lessons_per_day: 4,
        max_teacher_lessons_per_day: 5, max_consecutive_teacher_lessons: 2,
        time_priority: 'none', allow_flow_lessons: true, no_gaps_in_group_schedule: true,
      };
      allPreferences = prefs || [];
      allLabRooms = labRooms || [];
      cycleCommittees = cc || [];
      rooms = rm || [];
      teachers = tch || [];
    } catch (e) {
      bodyEl.innerHTML = `<div class="alert alert-error" style="margin:20px">${e.message}</div>`;
      return;
    }

    // Lab room state: Map "ccId_roomId" -> {id, cycle_committee_id, room_id, academic_year_id}
    const labRoomMap = new Map(
      allLabRooms.map(lr => [`${lr.cycle_committee_id}_${lr.room_id}`, { ...lr }])
    );
    const originalLabKeys = new Set(labRoomMap.keys());

    const weekdayShort = {
      monday:    t('weekdaysShort.monday'),
      tuesday:   t('weekdaysShort.tuesday'),
      wednesday: t('weekdaysShort.wednesday'),
      thursday:  t('weekdaysShort.thursday'),
      friday:    t('weekdaysShort.friday'),
      saturday:  t('weekdaysShort.saturday'),
    };

    // Two sub-views: main tabs and teacher preferences sub-panel
    const mainView = document.createElement('div');
    mainView.style.cssText = 'display:flex;flex-direction:column;flex:1;min-height:0;overflow:hidden;';
    const subPanel = document.createElement('div');
    subPanel.style.cssText = 'display:none;flex-direction:column;flex:1;min-height:0;overflow:hidden;';
    bodyEl.innerHTML = '';
    bodyEl.append(mainView, subPanel);

    const renderMainView = () => {
      const forbCount = allPreferences.filter(p => p.slot_type === 'forbidden').length;
      const prefCount = allPreferences.filter(p => p.slot_type === 'preferred').length;
      const ccOptions = cycleCommittees.map(cc =>
        `<option value="${cc.id}">${cc.name}</option>`
      ).join('');

      mainView.innerHTML = `
        <div class="sm-tabs">
          <button class="sm-tab active" data-panel="load">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M3 3v18h18"/><path d="M7 14l4-4 4 4 5-5"/></svg>
            ${t('modal.tabLoad')} <span class="sm-tab-badge">4</span>
          </button>
          <button class="sm-tab" data-panel="time">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg>
            ${t('modal.tabTime')} <span class="sm-tab-badge">4</span>
          </button>
          <button class="sm-tab" data-panel="space">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 21h18M5 21V7l7-4 7 4v14"/><path d="M9 9h.01M15 9h.01M9 13h.01M15 13h.01M9 17h.01M15 17h.01"/></svg>
            ${t('modal.tabSpace')} <span class="sm-tab-badge">1</span>
          </button>
          <button class="sm-tab" data-panel="org">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="9" cy="7" r="4"/><path d="M3 21v-2a4 4 0 014-4h4a4 4 0 014 4v2"/><circle cx="17" cy="7" r="3"/><path d="M21 21v-2a4 4 0 00-3-3.87"/></svg>
            ${t('modal.tabOrg')} <span class="sm-tab-badge">1</span>
          </button>
        </div>
        <div class="sm-panels">
          <div class="sm-panel active" data-panel="load">
            <div class="sm-constraint">
              <div class="sm-constraint-info">
                <div class="sm-constraint-label">
                  <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M16 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="8.5" cy="7" r="4"/><path d="M20 8v6M23 11h-6"/></svg></div>
                  ${t('modal.maxGroupDay')}
                </div>
                <div class="sm-constraint-desc">${t('modal.maxGroupDayDesc')}</div>
              </div>
              <div class="sm-control">
                <input type="number" class="sm-num-input" id="smr-max-group-day" value="${restriction.max_group_lessons_per_day}" min="1" max="8">
                <span class="sm-unit-label">${t('units.lessons')}</span>
              </div>
            </div>
            <div class="sm-constraint">
              <div class="sm-constraint-info">
                <div class="sm-constraint-label">
                  <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M20 21v-2a4 4 0 00-3-3.87M4 21v-2a4 4 0 013-3.87"/><circle cx="12" cy="7" r="4"/></svg></div>
                  ${t('modal.maxTeacherDay')}
                </div>
                <div class="sm-constraint-desc">${t('modal.maxTeacherDayDesc')}</div>
              </div>
              <div class="sm-control">
                <input type="number" class="sm-num-input" id="smr-max-teacher-day" value="${restriction.max_teacher_lessons_per_day}" min="1" max="8">
                <span class="sm-unit-label">${t('units.lessons')}</span>
              </div>
            </div>
            <div class="sm-constraint">
              <div class="sm-constraint-info">
                <div class="sm-constraint-label">
                  <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15A9 9 0 1 1 19 7.5"/></svg></div>
                  ${t('modal.consecutive')}
                </div>
                <div class="sm-constraint-desc">${t('modal.consecutiveDesc')}</div>
              </div>
              <div class="sm-control">
                <input type="number" class="sm-num-input" id="smr-consecutive" value="${restriction.max_consecutive_teacher_lessons}" min="0" max="8">
                <span class="sm-unit-label">${t('units.consecutive')}</span>
              </div>
            </div>
            <div class="sm-constraint">
              <div class="sm-constraint-info">
                <div class="sm-constraint-label">
                  <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M22 12H2M12 2v10"/><circle cx="12" cy="12" r="2" fill="currentColor"/></svg></div>
                  ${t('modal.minGroupDay')}
                </div>
                <div class="sm-constraint-desc">${t('modal.minGroupDayDesc')}</div>
              </div>
              <div class="sm-control">
                <input type="number" class="sm-num-input" id="smr-min-group-day" value="${restriction.min_group_lessons_per_day}" min="0" max="4">
                <span class="sm-unit-label">${t('units.lessons')}</span>
              </div>
            </div>
          </div>

          <div class="sm-panel" data-panel="time">
            <div class="sm-constraint">
              <div class="sm-constraint-info">
                <div class="sm-constraint-label">
                  <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><rect x="3" y="3" width="18" height="18" rx="2"/></svg></div>
                  ${t('modal.gapPrevention')}
                </div>
                <div class="sm-constraint-desc">${t('modal.gapPreventionDesc')}</div>
              </div>
              <div class="sm-control">
                <select class="sm-select" id="smr-no-gaps">
                  <option value="false" ${!restriction.no_gaps_in_group_schedule ? 'selected' : ''}>${t('modal.gapOff')}</option>
                  <option value="true" ${restriction.no_gaps_in_group_schedule ? 'selected' : ''}>${t('modal.gapOn')}</option>
                </select>
              </div>
            </div>
            <div class="sm-constraint">
              <div class="sm-constraint-info">
                <div class="sm-constraint-label">
                  <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><circle cx="12" cy="12" r="9"/><path d="M4.93 4.93l14.14 14.14"/></svg></div>
                  ${t('modal.forbiddenHours')}
                </div>
                <div class="sm-constraint-desc">${t('modal.forbiddenHoursDesc')}</div>
              </div>
              <div class="sm-control">
                <button class="btn btn-secondary btn-sm" id="smr-forbidden-btn">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M12 5v14M5 12h14"/></svg>
                  ${t('modal.configure')}${forbCount > 0 ? ` (${forbCount})` : ''}
                </button>
              </div>
            </div>
            <div class="sm-constraint">
              <div class="sm-constraint-info">
                <div class="sm-constraint-label">
                  <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg></div>
                  ${t('modal.teacherPrefs')}
                </div>
                <div class="sm-constraint-desc">${t('modal.teacherPrefsDesc')}</div>
              </div>
              <div class="sm-control">
                <button class="btn btn-secondary btn-sm" id="smr-preferred-btn">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M12 5v14M5 12h14"/></svg>
                  ${t('modal.configure')}${prefCount > 0 ? ` (${prefCount})` : ''}
                </button>
              </div>
            </div>
            <div class="sm-constraint">
              <div class="sm-constraint-info">
                <div class="sm-constraint-label">
                  <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="1"/><circle cx="19" cy="12" r="1"/><circle cx="5" cy="12" r="1"/></svg></div>
                  ${t('modal.timePriority')}
                </div>
                <div class="sm-constraint-desc">${t('modal.timePriorityDesc')}</div>
              </div>
              <div class="sm-control">
                <select class="sm-select" id="smr-time-priority">
                  <option value="none" ${restriction.time_priority === 'none' ? 'selected' : ''}>${t('timePriority.none')}</option>
                  <option value="morning" ${restriction.time_priority === 'morning' ? 'selected' : ''}>${t('timePriority.morning')}</option>
                  <option value="afternoon" ${restriction.time_priority === 'afternoon' ? 'selected' : ''}>${t('timePriority.afternoon')}</option>
                </select>
              </div>
            </div>
          </div>

          <div class="sm-panel" data-panel="space">
            <div class="sm-constraint" style="display:block;">
              <div class="sm-constraint-label" style="margin-bottom:4px;">
                <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8M12 17v4"/></svg></div>
                ${t('modal.labRooms')}
              </div>
              <div class="sm-constraint-desc" style="margin-left:34px;margin-bottom:12px;">${t('modal.labRoomsDesc')}</div>
              <div style="display:flex;gap:12px;flex-wrap:wrap;align-items:flex-start;">
                <select class="sm-select" id="smr-cycle-committee" style="min-width:200px;flex-shrink:0;">
                  <option value="">${t('modal.chooseCycleCommittee')}</option>
                  ${ccOptions}
                </select>
                <div class="sm-chips" id="smr-room-chips" style="flex:1;"></div>
              </div>
            </div>
          </div>

          <div class="sm-panel" data-panel="org">
            <div class="sm-constraint">
              <div class="sm-constraint-info">
                <div class="sm-constraint-label">
                  <div class="sm-constraint-icon"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="9" cy="7" r="4"/><path d="M3 21v-2a4 4 0 014-4h4a4 4 0 014 4v2"/><circle cx="17" cy="7" r="3"/><path d="M21 21v-2a4 4 0 00-3-3.87"/></svg></div>
                  ${t('modal.flowLessons')}
                </div>
                <div class="sm-constraint-desc">${t('modal.flowLessonsDesc')}</div>
              </div>
              <div class="sm-control">
                <span class="sm-unit-label" id="smr-flow-label">${restriction.allow_flow_lessons ? t('modal.flowOn') : t('modal.flowOff')}</span>
                <button class="sm-toggle ${restriction.allow_flow_lessons ? 'on' : ''}" id="smr-flow-toggle" role="switch" aria-checked="${restriction.allow_flow_lessons}"></button>
              </div>
            </div>
          </div>
        </div>`;

      // Tab switching
      mainView.querySelectorAll('.sm-tab').forEach(tab => {
        tab.addEventListener('click', () => {
          mainView.querySelectorAll('.sm-tab').forEach(t => t.classList.remove('active'));
          mainView.querySelectorAll('.sm-panel').forEach(p => p.classList.remove('active'));
          tab.classList.add('active');
          mainView.querySelector(`.sm-panel[data-panel="${tab.dataset.panel}"]`).classList.add('active');
        });
      });

      // Flow toggle
      const flowToggle = mainView.querySelector('#smr-flow-toggle');
      const flowLabel  = mainView.querySelector('#smr-flow-label');
      flowToggle.addEventListener('click', () => {
        flowToggle.classList.toggle('on');
        const isOn = flowToggle.classList.contains('on');
        flowToggle.setAttribute('aria-checked', String(isOn));
        flowLabel.textContent = isOn ? t('modal.flowOn') : t('modal.flowOff');
      });

      // Forbidden / preferred buttons
      mainView.querySelector('#smr-forbidden-btn').addEventListener('click', () => openSubPanel('forbidden'));
      mainView.querySelector('#smr-preferred-btn').addEventListener('click', () => openSubPanel('preferred'));

      // Cycle committee → room chips
      const ccSelect  = mainView.querySelector('#smr-cycle-committee');
      const chipsArea = mainView.querySelector('#smr-room-chips');
      const renderChips = (ccId) => {
        chipsArea.innerHTML = '';
        if (!ccId) return;
        rooms.forEach(room => {
          const key    = `${ccId}_${room.id}`;
          const active = labRoomMap.has(key);
          const chip   = document.createElement('button');
          chip.className = `sm-chip${active ? ' active' : ''}`;
          chip.textContent = room.name + (active ? ' ✓' : '');
          chip.dataset.roomId = room.id;
          chip.dataset.ccId   = ccId;
          chip.addEventListener('click', () => {
            const k = `${chip.dataset.ccId}_${chip.dataset.roomId}`;
            if (labRoomMap.has(k)) {
              labRoomMap.delete(k);
              chip.classList.remove('active');
              chip.textContent = room.name;
            } else {
              labRoomMap.set(k, {
                id: 0,
                cycle_committee_id: Number(chip.dataset.ccId),
                room_id: Number(chip.dataset.roomId),
                academic_year_id: Number(yearId),
              });
              chip.classList.add('active');
              chip.textContent = room.name + ' ✓';
            }
          });
          chipsArea.appendChild(chip);
        });
      };
      ccSelect.addEventListener('change', () => renderChips(ccSelect.value));
    };

    // Teacher preferences sub-panel
    const openSubPanel = (slotType) => {
      mainView.style.display = 'none';
      subPanel.style.display = 'flex';
      subPanel.innerHTML = '';

      const title = slotType === 'forbidden' ? t('modal.forbiddenTitle') : t('modal.preferredTitle');

      const header = document.createElement('div');
      header.className = 'sm-subpanel-header';
      const backBtn = document.createElement('button');
      backBtn.className = 'btn btn-secondary btn-sm';
      backBtn.innerHTML = `<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M19 12H5M12 5l-7 7 7 7"/></svg> ${t('modal.back')}`;
      backBtn.addEventListener('click', () => {
        subPanel.style.display = 'none';
        mainView.style.display = 'flex';
        renderMainView();
      });
      const titleEl = document.createElement('span');
      titleEl.style.cssText = 'font-size:13px;font-weight:600;color:var(--text);';
      titleEl.textContent = title;
      header.append(backBtn, titleEl);
      subPanel.appendChild(header);

      const body = document.createElement('div');
      body.className = 'sm-subpanel-body';
      subPanel.appendChild(body);

      const renderPrefList = () => {
        body.innerHTML = '';
        const filtered = allPreferences.filter(p => p.slot_type === slotType);

        if (filtered.length === 0) {
          const empty = document.createElement('div');
          empty.style.cssText = 'color:var(--text-muted);font-size:13px;text-align:center;padding:16px 0;';
          empty.textContent = t('modal.noSettings');
          body.appendChild(empty);
        } else {
          filtered.forEach(pref => {
            const row = document.createElement('div');
            row.className = 'sm-constraint';
            row.innerHTML = `
              <div class="sm-constraint-info">
                <div class="sm-constraint-label" style="font-size:12.5px;">
                  ${(() => { const tc = teachers.find(t => t.id === pref.teacher_id); const tUsers = App.cache['teacher-users'] || []; const u = tc ? tUsers.find(u => u.id === tc.user_id) : null; return u?.full_name || `#${pref.teacher_id}`; })()} &mdash; ${weekdayShort[pref.weekday] || pref.weekday}, ${t('modal.lessonNum')}${pref.lesson_number}
                </div>
              </div>
              <div class="sm-control">
                <button class="btn btn-sm" style="color:var(--danger);border:1px solid var(--danger);padding:4px 8px;line-height:1;" title="Видалити">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6"/></svg>
                </button>
              </div>`;
            row.querySelector('button').addEventListener('click', async () => {
              try {
                await Api.post('/api/syllabus/teacher-slot-preferences/delete', [pref.id]);
                allPreferences = allPreferences.filter(p => p.id !== pref.id);
                renderPrefList();
              } catch (err) { Toast.error(err.message); }
            });
            body.appendChild(row);
          });
        }

        // Add form
        const form = document.createElement('div');
        form.className = 'sm-add-form';
        const tUsers = App.cache['teacher-users'] || [];
        const teacherOptions = teachers.length > 0
          ? teachers.map(tc => { const u = tUsers.find(u => u.id === tc.user_id); return `<option value="${tc.id}">${u?.full_name || '#' + tc.id}</option>`; }).join('')
          : `<option value="" disabled>${t('modal.selectTeacher')}</option>`;
        form.innerHTML = `
          <div class="sm-add-form-field">
            <span class="sm-add-form-label">${t('fields.teacher')}</span>
            <select class="sm-select" id="sp-teacher" style="min-width:150px;">${teacherOptions}</select>
          </div>
          <div class="sm-add-form-field">
            <span class="sm-add-form-label">${t('fields.weekday')}</span>
            <select class="sm-select" id="sp-weekday" style="min-width:100px;">
              <option value="monday">${t('weekdays.monday')}</option>
              <option value="tuesday">${t('weekdays.tuesday')}</option>
              <option value="wednesday">${t('weekdays.wednesday')}</option>
              <option value="thursday">${t('weekdays.thursday')}</option>
              <option value="friday">${t('weekdays.friday')}</option>
              <option value="saturday">${t('weekdays.saturday')}</option>
            </select>
          </div>
          <div class="sm-add-form-field">
            <span class="sm-add-form-label">${t('modal.lessonNum')}</span>
            <input type="number" class="sm-num-input" id="sp-lesson" value="1" min="1" max="8" style="width:64px;">
          </div>
          <button class="btn btn-primary btn-sm" id="sp-add-btn" style="align-self:flex-end;">${t('actions.add')}</button>`;
        body.appendChild(form);

        form.querySelector('#sp-add-btn').addEventListener('click', async () => {
          const teacherId = Number(form.querySelector('#sp-teacher').value);
          const weekday   = form.querySelector('#sp-weekday').value;
          const lessonNum = Number(form.querySelector('#sp-lesson').value);
          if (!teacherId) { Toast.warning(t('modal.selectTeacher')); return; }
          const addBtn = form.querySelector('#sp-add-btn');
          addBtn.disabled = true;
          try {
            const created = await Api.post('/api/syllabus/teacher-slot-preferences', [{
              academic_year_id: Number(yearId),
              teacher_id: teacherId,
              weekday,
              lesson_number: lessonNum,
              slot_type: slotType,
            }]);
            if (created && created[0]) allPreferences.push(created[0]);
            renderPrefList();
          } catch (err) {
            Toast.error(err.message);
            addBtn.disabled = false;
          }
        });
      };

      renderPrefList();
    };

    renderMainView();

    // Footer
    const isNewRestrict = !restriction.id || restriction.id === 0;
    const statusDot  = document.createElement('span');
    statusDot.className = isNewRestrict ? 'sm-status-dot' : 'sm-status-dot saved';
    const statusText = document.createElement('span');
    statusText.textContent = isNewRestrict ? t('modal.unsaved') : t('modal.saved');
    const statusDiv  = document.createElement('div');
    statusDiv.className = 'sm-modal-status';
    statusDiv.append(statusDot, statusText);

    const cancelBtn = document.createElement('button');
    cancelBtn.className = 'btn btn-secondary';
    cancelBtn.textContent = t('actions.cancel');
    cancelBtn.addEventListener('click', closeModal);

    const saveBtn = document.createElement('button');
    saveBtn.className = 'btn btn-primary';
    saveBtn.innerHTML = `<svg class="icon"><use href="#icon-save"/></svg> ${t('actions.save')}`;

    const actions = document.createElement('div');
    actions.style.cssText = 'display:flex;gap:8px;';
    actions.append(cancelBtn, saveBtn);
    footEl.append(statusDiv, actions);

    saveBtn.addEventListener('click', async () => {
      saveBtn.disabled = true;
      try {
        // Save lab room changes
        const toDelete = [...originalLabKeys]
          .filter(k => !labRoomMap.has(k))
          .map(k => allLabRooms.find(lr => `${lr.cycle_committee_id}_${lr.room_id}` === k)?.id)
          .filter(Boolean);
        const toCreate = [...labRoomMap.values()].filter(lr => lr.id === 0);
        if (toDelete.length > 0) await Api.post('/api/syllabus/cycle-committee-lab-rooms/delete', toDelete);
        if (toCreate.length > 0) {
          const created = await Api.post('/api/syllabus/cycle-committee-lab-rooms', toCreate);
          if (created) {
            created.forEach(lr => {
              const k = `${lr.cycle_committee_id}_${lr.room_id}`;
              labRoomMap.set(k, lr);
              originalLabKeys.add(k);
            });
          }
        }

        // Save restriction
        const payload = {
          ...restriction,
          academic_year_id: Number(yearId),
          min_group_lessons_per_day: Number(mainView.querySelector('#smr-min-group-day').value),
          max_group_lessons_per_day: Number(mainView.querySelector('#smr-max-group-day').value),
          max_teacher_lessons_per_day: Number(mainView.querySelector('#smr-max-teacher-day').value),
          max_consecutive_teacher_lessons: Number(mainView.querySelector('#smr-consecutive').value),
          no_gaps_in_group_schedule: mainView.querySelector('#smr-no-gaps').value === 'true',
          time_priority: mainView.querySelector('#smr-time-priority').value,
          allow_flow_lessons: mainView.querySelector('#smr-flow-toggle').classList.contains('on'),
        };
        if (restriction.id && restriction.id > 0) {
          await Api.put('/api/syllabus/schedule-restrictions', [payload]);
        } else {
          const created = await Api.post('/api/syllabus/schedule-restrictions', [payload]);
          if (created && created[0]) restriction.id = created[0].id;
        }
        delete this.cache['schedule-restrictions'];
        Toast.success(t('messages.saved'));
        statusDot.classList.add('saved');
        statusText.textContent = t('modal.saved');
      } catch (e) {
        Toast.error(e.message);
      } finally {
        saveBtn.disabled = false;
      }
    });
  },

  async _doGenerate(semId, mainEl) {
    this.currentMgr = null;
    mainEl.innerHTML = `<div class="generate-progress"><div class="spinner"></div><span>${t('generate.generating')}</span></div>`;
    try {
      const data = await Api.post('/api/syllabus/schedule-templates/generate', {
        semester_id: Number(semId),
      });
      mainEl.innerHTML = '';

      const saveBar = document.createElement('div');
      saveBar.className = 'gen-save-bar';
      saveBar.innerHTML = `
        <span class="gen-save-label">${t('generate.success')}</span>
        <input type="text" class="gen-name-input" placeholder="${t('generate.saveNamePlaceholder')}">
        <button class="btn btn-secondary btn-sm" id="gen-view-btn">${t('actions.view')}</button>
        <button class="btn btn-success btn-sm" id="gen-save-btn">
          <svg class="icon"><use href="#icon-save"/></svg>
          ${t('generate.saveBtn')}
        </button>`;
      mainEl.appendChild(saveBar);

      const scheduleEl = document.createElement('div');
      mainEl.appendChild(scheduleEl);
      await this._renderScheduleData(data, scheduleEl);

      saveBar.querySelector('#gen-view-btn').addEventListener('click', () => {
        this.viewScheduleTemplate(data);
      });

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
      const savedTemplate = await Api.post('/api/syllabus/schedule-templates/save', payload);
      Toast.success(t('generate.savedSuccess'));
      delete this.cache['schedule-templates'];
      if (savedTemplate?.id) {
        this._openActivateScheduleModal(savedTemplate.id, semId, scheduleData);
      }
    } catch (e) {
      Toast.error(e.message);
    } finally {
      btn.disabled = false;
    }
  },

  async _openActivateScheduleModal(templateId, semId, scheduleData) {
    const groupIds = new Set();
    let totalLessons = 0;
    for (const weekData of Object.values(scheduleData)) {
      if (!weekData || typeof weekData !== 'object') continue;
      for (const dayData of Object.values(weekData)) {
        if (!dayData || typeof dayData !== 'object') continue;
        for (const lessons of Object.values(dayData)) {
          if (!Array.isArray(lessons)) continue;
          totalLessons += lessons.length;
          for (const l of lessons) {
            for (const sl of (l.sub_lessons || [])) {
              groupIds.add(sl.group_id);
            }
          }
        }
      }
    }

    await this.ensureCache('semesters');
    const sem = (this.cache['semesters'] || []).find(s => s.id === Number(semId));
    const semText = sem
      ? `${sem.period_start?.slice(0,10)} – ${sem.period_end?.slice(0,10)}`
      : `#${semId}`;

    const body = document.createElement('div');
    body.innerHTML = `
      <p style="margin-bottom:16px;color:var(--text-muted)">${t('generate.activateConfirmText')}</p>
      <div style="display:flex;flex-direction:column;gap:10px">
        <div style="display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid var(--border)">
          <span style="font-weight:500">${t('generate.activateSemester')}</span>
          <span>${semText}</span>
        </div>
        <div style="display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid var(--border)">
          <span style="font-weight:500">${t('generate.activateGroups')}</span>
          <span>${groupIds.size}</span>
        </div>
        <div style="display:flex;justify-content:space-between;padding:8px 0">
          <span style="font-weight:500">${t('generate.activateLessons')}</span>
          <span>${totalLessons}</span>
        </div>
      </div>`;

    const foot = document.createElement('div');
    foot.style.cssText = 'display:flex;gap:8px;justify-content:flex-end;width:100%';
    foot.innerHTML = `
      <button class="btn btn-secondary" id="act-cancel-btn">${t('actions.cancel')}</button>
      <button class="btn btn-primary" id="act-confirm-btn">${t('generate.activateBtn')}</button>`;

    openModal(t('generate.activateTitle'), body, foot);

    foot.querySelector('#act-cancel-btn').addEventListener('click', () => closeModal());
    foot.querySelector('#act-confirm-btn').addEventListener('click', async () => {
      const confirmBtn = foot.querySelector('#act-confirm-btn');
      confirmBtn.disabled = true;
      try {
        await Api.post(`/api/syllabus/schedule-templates/${templateId}/activate`, null);
        Toast.success(t('generate.activateSuccess'));
        delete this.cache['schedule-templates'];
        closeModal();
      } catch (e) {
        Toast.error(e.message);
        confirmBtn.disabled = false;
      }
    });
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
    const teacherName = id => {
      const teacher = teachers.find(tc => tc.id === id);
      if (!teacher) return `#${id}`;
      const user = users.find(u => u.id === teacher.user_id);
      return [user?.last_name, user?.first_name].filter(Boolean).join(' ') || user?.email?.split('@')[0] || `#${id}`;
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
    const scheduleGroupIds   = new Set();
    const scheduleTeacherIds = new Set();
    for (const weekData of Object.values(scheduleJson)) {
      if (!weekData) continue;
      for (const dayData of Object.values(weekData)) {
        if (!dayData) continue;
        for (const lessons of Object.values(dayData)) {
          for (const l of lessons) {
            for (const sl of (l.sub_lessons || [])) {
              scheduleGroupIds.add(sl.group_id);
              scheduleTeacherIds.add(sl.teacher_id);
            }
          }
        }
      }
    }
    const scheduleGroups   = groups.filter(g  => scheduleGroupIds.has(g.id));
    const scheduleTeachers = teachers.filter(tc => scheduleTeacherIds.has(tc.id));

    const typeLabel = (lesson) => {
      if (lesson.type === 'flow') return t('schedule.typeFlow');
      return lesson.is_lab ? t('schedule.typeLab') : t('schedule.typeLecture');
    };

    const icnTeacher = `<svg width="13" height="13" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 7a4 4 0 11-8 0 4 4 0 018 0z"/><path d="M12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>`;
    const icnGroup   = `<svg width="13" height="13" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 00-3-3.87"/><path d="M16 3.13a4 4 0 010 7.75"/></svg>`;
    const icnRoom    = `<svg width="13" height="13" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13S3 17 3 10a9 9 0 0118 0z"/><circle cx="12" cy="10" r="3"/></svg>`;

    const renderWeek = (weekData, groupFilterId, teacherFilterId) => {
      let rows = '';
      for (const slot of slots) {
        let cells = `<td class="sched-slot"><span class="sched-slot-num">${slot}</span><span class="sched-slot-time">${bellTime(slot)}</span></td>`;
        for (const day of DAYS) {
          let lessons = weekData?.[day]?.[String(slot)] || [];
          lessons = lessons
            .map(l => ({
              ...l,
              sub_lessons: (l.sub_lessons || []).filter(sl => {
                if (groupFilterId   && sl.group_id   !== groupFilterId)   return false;
                if (teacherFilterId && sl.teacher_id !== teacherFilterId) return false;
                return true;
              })
            }))
            .filter(l => l.sub_lessons.length > 0);

          let cellContent = '';
          for (const lesson of lessons) {
            const subj = subjectName(lesson.subject_id);
            const tLbl = typeLabel(lesson);
            const classKey = lesson.type === 'flow' ? 'flow' : (lesson.is_lab ? 'lab' : 'lecture');
            let detailRows = '';
            for (const sl of lesson.sub_lessons) {
              detailRows += `<div class="sched-lesson-detail">
                <div class="sched-detail-row">${icnTeacher}<span>${teacherName(sl.teacher_id)}</span></div>
                <div class="sched-detail-row">${icnGroup}<span>${groupName(sl.group_id)}</span></div>
                <div class="sched-detail-row">${icnRoom}<span>Ауд. ${roomName(sl.room_id)}</span></div>
              </div>`;
            }
            cellContent += `<div class="sched-lesson-card sched-type-${classKey}">
              <div class="sched-lesson-header">
                <div class="sched-lesson-subject">${subj}</div>
                <div class="sched-type-badge sched-type-badge-${classKey}">${tLbl}</div>
              </div>
              ${detailRows}
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

    const groupOptions   = scheduleGroups.map(g => `<option value="${g.id}">${g.name}</option>`).join('');
    const teacherOptions = scheduleTeachers.map(tc => `<option value="${tc.id}">${teacherName(tc.id)}</option>`).join('');

    container.innerHTML = `
      <div class="sched-view">
        <div class="sched-name">${result.name || ''}</div>
        <div class="sched-toolbar">
          <div class="sched-week-tabs">
            <button class="sched-week-btn active" data-week="numerator">${t('schedule.numerator')}</button>
            <button class="sched-week-btn" data-week="denominator">${t('schedule.denominator')}</button>
          </div>
        </div>
        <div class="filter-bar sched-filter-bar">
          <div class="filter-field">
            <label>${t('schedule.groupFilter')}</label>
            <select id="sched-group-select">
              <option value="">— ${t('schedule.allGroups')} —</option>
              ${groupOptions}
            </select>
          </div>
          <div class="filter-field">
            <label>${t('schedule.teacherFilter')}</label>
            <select id="sched-teacher-select">
              <option value="">— ${t('schedule.allTeachers')} —</option>
              ${teacherOptions}
            </select>
          </div>
        </div>
        <div id="sched-content"></div>
        <div class="sched-legend">
          <span class="sched-legend-item"><span class="sched-legend-dot sched-legend-united"></span>${t('schedule.legendLecture')}</span>
          <span class="sched-legend-item"><span class="sched-legend-dot sched-legend-split"></span>${t('schedule.legendLab')}</span>
          <span class="sched-legend-item"><span class="sched-legend-dot sched-legend-flow"></span>${t('schedule.legendFlow')}</span>
        </div>
      </div>`;

    const content        = container.querySelector('#sched-content');
    const activeGroupId   = () => Number(container.querySelector('#sched-group-select').value)   || 0;
    const activeTeacherId = () => Number(container.querySelector('#sched-teacher-select').value) || 0;
    const activeWeek      = () => container.querySelector('.sched-week-btn.active')?.dataset.week || 'numerator';

    const showWeek = week => {
      content.innerHTML = renderWeek(scheduleJson[week] || {}, activeGroupId(), activeTeacherId());
      container.querySelectorAll('.sched-week-btn').forEach(btn =>
        btn.classList.toggle('active', btn.dataset.week === week));
    };

    container.querySelectorAll('.sched-week-btn').forEach(btn =>
      btn.addEventListener('click', () => showWeek(btn.dataset.week)));

    container.querySelector('#sched-group-select').addEventListener('change', () => showWeek(activeWeek()));
    container.querySelector('#sched-teacher-select').addEventListener('change', () => showWeek(activeWeek()));

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
          <div class="filter-bar">
            <div class="filter-field">
              <label>${t('fields.specialty')}</label>
              <select id="gmd-spec">
                <option value="">—</option>
                ${specialties.map(s=>`<option value="${s.id}"${filterSpec===String(s.id)?' selected':''}>${s.name}</option>`).join('')}
              </select>
            </div>
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
          <div class="gmd-footer">
            <button class="btn btn-primary btn-sm" id="gmd-add">
              <svg class="icon icon-sm"><use href="#icon-plus"/></svg>${t('actions.add')}
            </button>
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

    // Fetch study-plans filtered by the selected academic year
    let studyPlans = [];
    try {
      const spPath = ENTITY_CONFIGS['study-plans'].apiPath +
        (this.state.yearId ? `?academic_year_id=${this.state.yearId}` : '');
      studyPlans = await Api.get(spPath) || [];
      this.cache['study-plans'] = studyPlans;
    } catch { studyPlans = []; }

    await this.ensureCache('groups');
    await this.ensureCache('specialties');
    await this.ensureCache('disciplines');
    await this.ensureCache('teachers');
    await this.ensureTeacherUsersCache();

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

      // Narrow to the current academic year via study plan IDs
      if (this.state.yearId) {
        const spIds   = new Set(studyPlans.map(sp => sp.id));
        distributions = distributions.filter(d => spIds.has(d.study_plan_id));
        const distIds = new Set(distributions.map(d => d.id));
        assignments   = assignments.filter(a => distIds.has(a.workload_distribution_id));
      }
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
          <div class="gmd-footer">
            <button class="btn btn-primary btn-sm" id="wmd-add">
              <svg class="icon icon-sm"><use href="#icon-plus"/></svg>${t('actions.add')}
            </button>
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

    // Check for pending transfer data from study plans
    if (this._pendingWorkload) {
      const pending = this._pendingWorkload;
      this._pendingWorkload = null;
      this._openWorkloadDistModal(pending, studyPlans, groups, specialties, disciplines, async (saved) => {
        const created = await Api.post(ENTITY_CONFIGS['workload-distributions'].apiPath, [saved]);
        distributions = distributions.concat(created || []);
        if (created?.[0]) selectedId = created[0].id;
        redraw();
      });
    }
  },

  _openWorkloadDistModal(dist, studyPlans, groups, specialties, disciplines, onSave) {
    const isEdit = !!(dist && dist.id);
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
          ${['assistant','lector'].map(v =>
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
      <button class="btn btn-primary" id="transfer-go">${t('workload.transfer')}</button>`;
    openModal(t('workload.transfer'), body, foot);
    foot.querySelector('#transfer-go').addEventListener('click', () => {
      const groupId = document.getElementById('transfer-group').value;
      if (!groupId) { Toast.warning(t('messages.fillRequired')); return; }
      this._pendingWorkload = {
        study_plan_id:  studyPlan.id,
        group_id:       Number(groupId),
        classroom_work: studyPlan.lectures   || null,
        laboratory:     studyPlan.laboratory || null,
        practical:      studyPlan.practical  || null,
        exam:           studyPlan.exam       || null,
      };
      closeModal();
      location.hash = '#/processing/workload-distributions';
    });
  },

  // ═══════════════════════════════════════════════════════════
  // USERS PAGE (admin only)
  // ═══════════════════════════════════════════════════════════
  async renderUsers(main) {
    const role = this.state.user?.role;
    if (role !== 'admin' && role !== 'dean') {
      main.innerHTML = `<div class="alert alert-error">${t('messages.forbidden')}</div>`;
      return;
    }
    this.currentMgr = null;

    main.innerHTML = `
      <div class="page-header"><h1>${t('pages.usersTitle')}</h1></div>
      <div id="users-content"></div>`;

    // Invalidate users cache so EntityManager loads fresh data
    delete this.cache['users'];
    await this.renderEntityTab('users', main.querySelector('#users-content'));

    const toolbarLeft = main.querySelector('#users-content .toolbar-left');
    if (toolbarLeft) {
      if (role === 'admin') {
        // Register button
        const btn = document.createElement('button');
        btn.className = 'btn btn-primary btn-sm';
        btn.id = 'add-user-btn';
        btn.innerHTML = `<svg class="icon icon-sm"><use href="#icon-plus"/></svg> ${t('users.register')}`;
        btn.addEventListener('click', () => this._openRegisterModal());
        toolbarLeft.appendChild(btn);
      } else {
        // Dean: hide delete button
        const delBtn = toolbarLeft.querySelector('#em-del-btn');
        if (delBtn) delBtn.style.display = 'none';
      }
    }
  },

  _openRegisterModal() {
    const body = document.createElement('div');
    body.innerHTML = `
      <div class="form-field"><label>${t('fields.firstName')} *</label>
        <input type="text" id="reg-first-name"></div>
      <div class="form-field"><label>${t('fields.lastName')} *</label>
        <input type="text" id="reg-last-name"></div>
      <div class="form-field"><label>${t('fields.patronymic')} <span style="color:var(--text-muted);font-size:.78rem">(${t('actions.optional')})</span></label>
        <input type="text" id="reg-patronymic"></div>
      <div class="form-field"><label>${t('fields.email')} *</label>
        <input type="email" id="reg-email"></div>
      <div class="form-field"><label>${t('fields.role')}</label>
        <select id="reg-role">
          <option value="teacher">${t('roles.teacher')}</option>
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
      const email      = document.getElementById('reg-email').value.trim();
      const firstName  = document.getElementById('reg-first-name').value.trim();
      const lastName   = document.getElementById('reg-last-name').value.trim();
      const patronymic = document.getElementById('reg-patronymic').value.trim() || null;
      const role       = document.getElementById('reg-role').value;
      if (!email || !firstName || !lastName) {
        Toast.warning(t('messages.fillRequired'));
        return;
      }
      try {
        const btn = foot.querySelector('#reg-save');
        btn.disabled = true;
        const payload = { email, first_name: firstName, last_name: lastName, role };
        if (patronymic) payload.patronymic = patronymic;
        const resp = await Api.post('/api/auth/auth/invite', payload);
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
        // Reload the users entity tab after registration
        const usersContent = document.querySelector('#users-content');
        if (usersContent) {
          delete this.cache['users'];
          await this.renderEntityTab('users', usersContent);
          const tl = usersContent.querySelector('.toolbar-left');
          if (tl && !tl.querySelector('#add-user-btn') && this.state.user?.role === 'admin') {
            const rb = document.createElement('button');
            rb.className = 'btn btn-primary btn-sm';
            rb.id = 'add-user-btn';
            rb.innerHTML = `<svg class="icon icon-sm"><use href="#icon-plus"/></svg> ${t('users.register')}`;
            rb.addEventListener('click', () => this._openRegisterModal());
            tl.appendChild(rb);
          }
        }
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
      let errMsg;
      if (e.message === 'unauthorized' || e.message.includes('401')) {
        errMsg = t('auth.loginError');
      } else if (e.message === 'account deactivated') {
        errMsg = t('auth.accountDeactivated');
      } else {
        errMsg = e.message;
      }
      errEl.textContent = errMsg;
      errEl.classList.remove('hidden');
    }
  });

  document.getElementById('login-password')?.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') document.getElementById('login-btn')?.click();
  });

  // Password visibility toggle on login page
  document.getElementById('login-pw-toggle')?.addEventListener('click', () => {
    const pw = document.getElementById('login-password');
    const eye = document.getElementById('login-pw-eye');
    if (pw.type === 'password') {
      pw.type = 'text';
      eye.innerHTML = '<path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94"/><path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19"/><line x1="1" y1="1" x2="23" y2="23"/>';
    } else {
      pw.type = 'password';
      eye.innerHTML = '<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>';
    }
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
