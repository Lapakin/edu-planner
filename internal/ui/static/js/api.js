'use strict';

const Api = (() => {
  let _token = localStorage.getItem('token') || null;

  function setToken(token) {
    _token = token;
    if (token) localStorage.setItem('token', token);
    else localStorage.removeItem('token');
  }

  function getToken() { return _token; }

  function parseJWT(token) {
    try {
      const base64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
      return JSON.parse(atob(base64));
    } catch { return null; }
  }

  function isTokenValid() {
    if (!_token) return false;
    const claims = parseJWT(_token);
    if (!claims) return false;
    return claims.exp * 1000 > Date.now();
  }

  async function request(method, url, body, opts = {}) {
    const headers = { 'Content-Type': 'application/json' };
    if (_token) headers['Authorization'] = `Bearer ${_token}`;
    const fetchOpts = { method, headers };

    // Format dates before sending
    if (body !== undefined && body !== null) {
      const formatDates = (obj) => {
        if (Array.isArray(obj)) return obj.map(formatDates);
        if (typeof obj !== 'object' || obj === null) return obj;
        const out = { ...obj };
        for (const [k, v] of Object.entries(out)) {
          if (typeof v === 'string' && v.match(/^\d{4}-\d{2}-\d{2}$/)) {
            out[k] = v + 'T00:00:00Z';
          } else if (typeof v === 'object') {
            out[k] = formatDates(v);
          }
        }
        return out;
      };
      body = formatDates(body);
      fetchOpts.body = JSON.stringify(body);
    }

    let resp;
    try {
      resp = await fetch(url, fetchOpts);
    } catch (e) {
      throw new Error(t('messages.loadError'));
    }

    if (resp.status === 401 && !opts.isLogin) {
      setToken(null);
      window.dispatchEvent(new CustomEvent('auth:logout'));
      throw new Error('unauthorized');
    }

    if (!resp.ok) {
      let msg = `HTTP ${resp.status}`;
      try {
        const err = await resp.json();
        msg = err.error || err.message || msg;
      } catch { /* ignore */ }
      throw new Error(msg);
    }

    if (resp.status === 204 || resp.headers.get('content-length') === '0') return null;
    const ct = resp.headers.get('content-type') || '';
    if (ct.includes('application/json')) return resp.json();
    return null;
  }

  return {
    setToken,
    getToken,
    parseJWT,
    isTokenValid,
    get:    (url)          => request('GET',    url),
    post:   (url, body)    => request('POST',   url, body),
    put:    (url, body)    => request('PUT',    url, body),
    delete: (url, body)    => request('DELETE', url, body),
    login:  (email, pass)  => request('POST', '/api/auth/login',
                                { email, password: pass }, { isLogin: true }),
  };
})();
