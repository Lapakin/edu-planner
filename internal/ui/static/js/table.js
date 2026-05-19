'use strict';

class EntityManager {
  constructor(config, cache) {
    this.config   = config;
    this.cache    = cache;       // {entityKey: [rows]}
    this.rows     = [];          // loaded rows (originals)
    this.originals = new Map();  // id -> {...original}
    this.changes   = new Map();  // id -> {field: newValue}
    this.newRows   = [];         // rows staged for creation (have _tempId < 0)
    this.deletedIds = new Set(); // ids staged for deletion
    this._tempId   = -1;
    this.filterValues = {};
    this._container = null;
    this._page = 0;
    this._pageSize = 20;
  }

  // ── Data loading ──────────────────────────────────────────────────
  setRows(rows) {
    this.rows = rows;
    this.originals.clear();
    this.changes.clear();
    this.newRows = [];
    this.deletedIds.clear();
    rows.forEach(r => this.originals.set(r.id, JSON.parse(JSON.stringify(r))));
  }

  isDirty() {
    return this.changes.size > 0 || this.newRows.length > 0 || this.deletedIds.size > 0;
  }

  // ── Filter ────────────────────────────────────────────────────────
  _filteredRows() {
    return this.rows.filter(row => {
      for (const [key, val] of Object.entries(this.filterValues)) {
        if (!val) continue;
        const rv = String(row[key] ?? '');
        if (!rv.toLowerCase().includes(String(val).toLowerCase())) return false;
      }
      return true;
    });
  }

  // ── Dirty tracking ────────────────────────────────────────────────
  _markCellChanged(id, field, newVal, tdEl) {
    const orig = this.originals.get(id);
    if (!orig) return;
    const origStr = String(orig[field] ?? '');
    const newStr  = String(newVal ?? '');
    if (!this.changes.has(id)) this.changes.set(id, {});
    if (newStr !== origStr) {
      this.changes.get(id)[field] = newVal;
      tdEl && tdEl.classList.add('cell-changed');
    } else {
      delete this.changes.get(id)[field];
      if (Object.keys(this.changes.get(id)).length === 0) this.changes.delete(id);
      tdEl && tdEl.classList.remove('cell-changed');
    }
    // update row highlight
    const tr = tdEl && tdEl.closest('tr');
    if (tr) tr.classList.toggle('row-changed', this.changes.has(id));
    this._updateToolbar();
  }

  // ── Inline row creation ───────────────────────────────────────────
  _addInlineRow() {
    const rowData = { _tempId: this._tempId-- };
    this.config.columns.forEach(col => {
      if (col.readonly) return;
      rowData[col.key] = col.type === 'bool' ? false : null;
    });
    this.newRows.unshift(rowData);
    this._page = 0; // go to first page
    this._renderTable();
    this._updateToolbar();
    setTimeout(() => {
      const tr = this._container?.querySelector('tbody tr.row-new');
      const firstEditable = tr?.querySelector('td.editable');
      if (firstEditable) firstEditable.click();
    }, 0);
  }

  // ── Render ────────────────────────────────────────────────────────
  render(container) {
    this._container = container;
    container.innerHTML = '';

    const { config } = this;
    const noCreate = config.noCreate;
    const noEdit   = config.noEdit;

    // Toolbar (built first, appended last)
    const toolbar = document.createElement('div');
    toolbar.className = 'entity-toolbar';
    toolbar.innerHTML = `
      <div class="toolbar-left">
        ${!noCreate ? `<button class="btn btn-primary btn-sm" id="em-add-btn">
          <svg class="icon icon-sm"><use href="#icon-plus"/></svg>
          <span>${t('actions.add')}</span>
        </button>` : ''}
        <button class="btn btn-danger btn-sm" id="em-del-btn" disabled>
          <svg class="icon icon-sm"><use href="#icon-trash"/></svg>
          <span>${t('actions.deleteSelected')}</span>
        </button>
      </div>
      <div class="toolbar-right">
        <button class="btn btn-primary btn-sm hidden" id="em-save-btn">
          <svg class="icon icon-sm"><use href="#icon-save"/></svg>
          <span>${t('actions.save')}</span>
        </button>
        <button class="btn btn-secondary btn-sm hidden" id="em-discard-btn">
          <svg class="icon icon-sm"><use href="#icon-x"/></svg>
          <span>${t('actions.discard')}</span>
        </button>
      </div>`;

    this._saveBtn    = toolbar.querySelector('#em-save-btn');
    this._discardBtn = toolbar.querySelector('#em-discard-btn');
    this._delBtn     = toolbar.querySelector('#em-del-btn');

    if (!noCreate) toolbar.querySelector('#em-add-btn').addEventListener('click', () => this._addInlineRow());
    this._delBtn.addEventListener('click', () => this._deleteSelected());
    this._saveBtn.addEventListener('click', () => this.save());
    this._discardBtn.addEventListener('click', () => this.discard());

    // Unsaved warning bar (appended with toolbar at bottom)
    const warnBar = document.createElement('div');
    warnBar.className = 'unsaved-warning hidden';
    warnBar.id = 'em-warn-bar';
    warnBar.textContent = t('messages.unsavedChanges');

    // Filter bar
    if (config.filters && config.filters.length > 0) {
      const fb = document.createElement('div');
      fb.className = 'filter-bar';
      config.filters.forEach(f => {
        const wrap = document.createElement('div');
        wrap.className = 'filter-item';
        wrap.innerHTML = `<label>${t(f.labelKey)}</label>`;
        let input;
        if (f.type === 'select') {
          input = document.createElement('select');
          input.innerHTML = `<option value="">${t('actions.add').replace(t('actions.add'), '—')}</option>`;
          const opts = f.ref ? (this.cache[f.ref] || []) : (f.options || []);
          opts.forEach(o => {
            const op = document.createElement('option');
            op.value = f.ref ? (o.id || o.value) : o.value;
            op.textContent = f.ref ? (o[f.refLabel] || o.name || o.id) : t(o.labelKey);
            input.appendChild(op);
          });
          // Pre-populate from global year/semester state
          if (f.key === 'academic_year_id' && App.state.yearId) input.value = App.state.yearId;
          if (f.key === 'semester_id' && App.state.semesterId) input.value = App.state.semesterId;
        } else {
          input = document.createElement('input');
          input.type = 'text';
          input.placeholder = t(f.labelKey);
        }
        input.addEventListener('change', () => {
          this.filterValues[f.key] = input.value;
          this._page = 0; // reset page on filter change
          this._renderTable();
        });
        this.filterValues[f.key] = input.value;
        wrap.appendChild(input);
        fb.appendChild(wrap);
      });
      container.appendChild(fb);
    }

    // Table wrapper
    const tableWrap = document.createElement('div');
    tableWrap.className = 'table-wrapper';
    tableWrap.id = 'em-table-wrap';
    container.appendChild(tableWrap);

    // Toolbar and warning at the bottom
    container.appendChild(warnBar);
    container.appendChild(toolbar);

    this._renderTable();
  }

  _renderTable() {
    const wrap = this._container.querySelector('#em-table-wrap');
    if (!wrap) return;
    wrap.innerHTML = '';

    const { config } = this;
    const filteredRows = this._filteredRows();

    // Pagination
    const totalExisting = filteredRows.length;
    const totalPages = Math.max(1, Math.ceil(totalExisting / this._pageSize));
    if (this._page >= totalPages) this._page = Math.max(0, totalPages - 1);
    const pageStart = this._page * this._pageSize;
    const pageEnd   = pageStart + this._pageSize;
    const pagedRows  = filteredRows.slice(pageStart, pageEnd);

    // new rows always show at top (no pagination)
    const allRows = [...this.newRows, ...pagedRows];

    if (allRows.length === 0 && totalExisting === 0) {
      wrap.innerHTML = `<div class="empty-state">${t('messages.noData') || '—'}</div>`;
      return;
    }

    const table = document.createElement('table');
    table.className = 'entity-table';

    // Header
    const thead = document.createElement('thead');
    let thHtml = '<tr><th class="col-check"><input type="checkbox" id="em-check-all"></th>';
    config.columns.forEach(col => {
      thHtml += `<th>${t(col.labelKey)}</th>`;
    });
    if (config.extraActions && config.extraActions.length > 0) thHtml += '<th></th>';
    thHtml += '</tr>';
    thead.innerHTML = thHtml;
    table.appendChild(thead);

    thead.querySelector('#em-check-all')?.addEventListener('change', (e) => {
      table.querySelectorAll('.row-check').forEach(cb => { cb.checked = e.target.checked; });
      this._updateDeleteBtn();
    });

    // Body
    const tbody = document.createElement('tbody');
    allRows.forEach(row => {
      const isNew     = row._tempId !== undefined;
      const isDeleted = !isNew && this.deletedIds.has(row.id);
      const isDirty   = !isNew && this.changes.has(row.id);

      const tr = document.createElement('tr');
      tr.dataset.id  = isNew ? row._tempId : row.id;
      tr.dataset.new = isNew ? '1' : '0';
      if (isNew)     tr.classList.add('row-new');
      if (isDeleted) tr.classList.add('row-deleted');
      if (isDirty)   tr.classList.add('row-changed');

      // Checkbox
      const tdCheck = document.createElement('td');
      tdCheck.className = 'col-check';
      const cb = document.createElement('input');
      cb.type = 'checkbox';
      cb.className = 'row-check';
      cb.addEventListener('change', () => this._updateDeleteBtn());
      tdCheck.appendChild(cb);
      tr.appendChild(tdCheck);

      // Data cells
      config.columns.forEach(col => {
        const td = document.createElement('td');
        const rawVal = isNew ? (row[col.key] ?? '') : this._currentVal(row.id, col.key, row[col.key]);

        if (col.type === 'badge') {
          const key = String(rawVal);
          td.innerHTML = `<span class="badge badge-${key === 'true' || rawVal === true ? 'success' : 'neutral'}">${t('status.' + (col.values ? (col.values[key] || key) : key))}</span>`;
        } else if (col.readonly || config.noEdit) {
          td.textContent = this._displayVal(col, rawVal);
        } else {
          td.className = 'editable';
          td.dataset.field = col.key;
          td.dataset.origVal = isNew ? '' : String(row[col.key] ?? '');
          td.textContent = this._displayVal(col, rawVal);

          const isDirtyCell = !isNew && this.changes.get(row.id)?.[col.key] !== undefined;
          if (isDirtyCell) td.classList.add('cell-changed');

          td.addEventListener('click', () => this._startEdit(td, col, isNew ? null : row.id, row, isNew));
        }
        tr.appendChild(td);
      });

      // Extra actions
      if (config.extraActions && config.extraActions.length > 0) {
        const tdAct = document.createElement('td');
        tdAct.className = 'col-actions';
        config.extraActions.forEach(act => {
          const btn = document.createElement('button');
          btn.className = act.classFn ? act.classFn(row) : 'btn btn-sm btn-secondary';
          btn.textContent = act.labelFn ? act.labelFn(row) : (act.label || '');
          btn.addEventListener('click', (e) => { e.stopPropagation(); act.handler(row, this); });
          tdAct.appendChild(btn);
        });
        tr.appendChild(tdAct);
      }

      tbody.appendChild(tr);
    });
    table.appendChild(tbody);
    wrap.appendChild(table);

    // Pagination controls (only if more than one page worth of existing rows)
    if (totalExisting > this._pageSize) {
      const pg = document.createElement('div');
      pg.className = 'pagination';
      const rangeStart = pageStart + 1;
      const rangeEnd   = Math.min(pageEnd, totalExisting);
      pg.innerHTML = `
        <button class="btn btn-sm btn-secondary" id="pg-prev" ${this._page === 0 ? 'disabled' : ''}>‹</button>
        <span class="pg-info">${rangeStart}–${rangeEnd} / ${totalExisting}</span>
        <button class="btn btn-sm btn-secondary" id="pg-next" ${this._page === totalPages - 1 ? 'disabled' : ''}>›</button>`;
      wrap.appendChild(pg);

      pg.querySelector('#pg-prev').addEventListener('click', () => {
        if (this._page > 0) { this._page--; this._renderTable(); }
      });
      pg.querySelector('#pg-next').addEventListener('click', () => {
        if (this._page < totalPages - 1) { this._page++; this._renderTable(); }
      });
    }
  }

  _currentVal(id, field, fallback) {
    return this.changes.get(id)?.[field] !== undefined
      ? this.changes.get(id)[field]
      : fallback;
  }

  _displayVal(col, val) {
    if (val === null || val === undefined || val === '') return '—';
    if (col.type === 'bool') return val ? '✓' : '✗';
    if (col.type === 'date') return val ? String(val).slice(0, 10) : '—';
    if (col.type === 'select') {
      if (col.ref) {
        const list = this.cache[col.ref] || [];
        const found = list.find(o => String(o.id || o.value) === String(val));
        return found ? (found[col.refLabel] || found.name || found.id) : String(val);
      }
      if (col.options) {
        const o = col.options.find(x => x.value === String(val));
        return o ? t(o.labelKey) : String(val);
      }
    }
    return String(val);
  }

  _startEdit(td, col, id, row, isNew) {
    if (td.querySelector('input,select')) return; // already editing
    const currentVal = isNew ? (row[col.key] ?? '') : this._currentVal(id, col.key, row[col.key]);
    td.classList.add('cell-editing');
    const origText = td.textContent;
    td.textContent = '';

    let input;
    if (col.type === 'select') {
      input = document.createElement('select');
      if (col.ref) {
        const list = this.cache[col.ref] || [];
        input.innerHTML = `<option value="">—</option>`;
        list.forEach(o => {
          const opt = document.createElement('option');
          opt.value = o.id || o.value;
          opt.textContent = o[col.refLabel] || o.name || o.id;
          if (String(o.id || o.value) === String(currentVal)) opt.selected = true;
          input.appendChild(opt);
        });
      } else if (col.options) {
        input.innerHTML = `<option value="">—</option>`;
        col.options.forEach(o => {
          const opt = document.createElement('option');
          opt.value = o.value;
          opt.textContent = t(o.labelKey);
          if (o.value === String(currentVal)) opt.selected = true;
          input.appendChild(opt);
        });
      }
    } else if (col.type === 'bool') {
      input = document.createElement('select');
      input.innerHTML = `<option value="true" ${currentVal ? 'selected' : ''}>✓</option>
                         <option value="false" ${!currentVal ? 'selected' : ''}>✗</option>`;
    } else {
      input = document.createElement('input');
      input.type = col.type === 'number' ? 'number' : (col.type === 'date' ? 'date' : 'text');
      input.value = col.type === 'date' ? String(currentVal || '').slice(0, 10) : (currentVal ?? '');
    }

    input.className = 'cell-input';
    td.appendChild(input);
    input.focus();

    const commit = () => {
      let newVal = input.value;
      if (col.type === 'bool') newVal = newVal === 'true';
      if (col.type === 'number') newVal = newVal === '' ? null : Number(newVal);
      if (col.type === 'select' && col.ref) newVal = newVal === '' ? null : Number(newVal);
      td.classList.remove('cell-editing');
      td.textContent = this._displayVal(col, newVal) || '—';

      if (isNew) {
        row[col.key] = newVal;
        td.classList.add('cell-changed');
      } else {
        this._markCellChanged(id, col.key, newVal, td);
      }
    };

    input.addEventListener('blur',    commit);
    input.addEventListener('keydown', (e) => {
      if (e.key === 'Enter')  { e.preventDefault(); input.blur(); }
      if (e.key === 'Escape') { input.value = col.type === 'bool' ? String(currentVal) : (currentVal ?? ''); input.blur(); }
    });
  }

  _updateDeleteBtn() {
    const checked = this._container?.querySelectorAll('.row-check:checked').length > 0;
    if (this._delBtn) this._delBtn.disabled = !checked;
  }

  _updateToolbar() {
    const dirty = this.changes.size > 0 || this.newRows.length > 0 || this.deletedIds.size > 0;
    this._saveBtn?.classList.toggle('hidden', !dirty);
    this._discardBtn?.classList.toggle('hidden', !dirty);
    const warn = document.querySelector('.unsaved-warning');
    if (warn) warn.classList.toggle('hidden', !dirty);
  }

  _formatDatesOut(obj) {
    const res = { ...obj };
    this.config.columns.forEach(col => {
      if (col.type === 'date' && res[col.key]) {
        // If it's a date only (YYYY-MM-DD), append T00:00:00Z
        if (typeof res[col.key] === 'string' && res[col.key].length === 10) {
          res[col.key] = res[col.key] + 'T00:00:00Z';
        }
      }
    });
    return res;
  }

  // ── Delete selected ───────────────────────────────────────────────
  async _deleteSelected() {
    const checked = [...this._container.querySelectorAll('.row-check:checked')];
    if (!checked.length) return;

    const ok = await confirmDialog(t('messages.confirmDelete'));
    if (!ok) return;

    checked.forEach(cb => {
      const tr  = cb.closest('tr');
      const isNew = tr.dataset.new === '1';
      const id    = tr.dataset.id;
      if (isNew) {
        this.newRows = this.newRows.filter(r => String(r._tempId) !== id);
      } else {
        this.deletedIds.add(Number(id));
      }
    });
    this._renderTable();
    this._updateToolbar();
  }

  // ── Save ──────────────────────────────────────────────────────────
  async save() {
    try {
      // 1. Create new rows
      if (this.newRows.length > 0) {
        const toCreate = this.newRows.map(r => {
          const obj = this._formatDatesOut(r);
          delete obj._tempId;
          if (this.config.autoYearFilter && App.state.yearId) {
            obj.academic_year_id = Number(App.state.yearId);
          }
          return obj;
        });
        const created = await Api.post(this.config.apiPath, toCreate);
        if (Array.isArray(created)) {
          created.forEach(r => this.rows.push(r));
        }
        this.newRows = [];
      }

      // 2. Update changed rows
      if (this.changes.size > 0) {
        const toUpdate = [];
        this.changes.forEach((delta, id) => {
          const orig = this.originals.get(id) || {};
          const updated = this._formatDatesOut({ ...orig, ...delta, id });
          if (this.config.autoYearFilter && App.state.yearId) {
            updated.academic_year_id = Number(App.state.yearId);
          }
          toUpdate.push(updated);
        });
        await Api.put(this.config.apiPath, toUpdate);
        // merge into originals
        toUpdate.forEach(r => {
          this.originals.set(r.id, { ...r });
          const row = this.rows.find(x => x.id === r.id);
          if (row) Object.assign(row, r);
        });
        this.changes.clear();
      }

      // 3. Delete
      if (this.deletedIds.size > 0) {
        await Api.post(this.config.apiPath + '/delete', [...this.deletedIds]);
        this.rows = this.rows.filter(r => !this.deletedIds.has(r.id));
        this.deletedIds.clear();
      }

      Toast.success(t('messages.saved'));
    } catch (e) {
      Toast.error(t('messages.saveError') + ': ' + e.message);
      return;
    }

    this._renderTable();
    this._updateToolbar();
    // refresh cache for this entity so other dropdowns update
    App.refreshCache(this.config.apiPath);
  }

  discard() {
    this.changes.clear();
    this.newRows = [];
    this.deletedIds.clear();
    this._renderTable();
    this._updateToolbar();
  }

  // Called externally after activate/deactivate
  reloadRow(updatedRow) {
    const idx = this.rows.findIndex(r => r.id === updatedRow.id);
    if (idx >= 0) this.rows[idx] = updatedRow;
    this.originals.set(updatedRow.id, { ...updatedRow });
    this._renderTable();
  }
}
