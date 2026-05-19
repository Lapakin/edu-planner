console.log("Welcome to Edu Planner!");

let accessToken = null;
let isRefreshing = false;
let refreshSubscribers = [];

function subscribeTokenRefresh(cb) {
    refreshSubscribers.push(cb);
}

function onRefreshed(token) {
    refreshSubscribers.map(cb => cb(token));
    refreshSubscribers = [];
}

async function apiCall(url, options = {}) {
    if (!options.headers) {
        options.headers = {};
    }

    if (accessToken) {
        options.headers['Authorization'] = `Bearer ${accessToken}`;
    }

    // Format dates before sending in custom apiCall
    if (options.body && typeof options.body === 'string') {
        try {
            let parsed = JSON.parse(options.body);
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
            parsed = formatDates(parsed);
            options.body = JSON.stringify(parsed);
        } catch (e) {
            // Ignore parsing errors, pass original body
        }
    }

    let response = await fetch(url, options);

    if (response.status === 401) {
        if (!isRefreshing) {
            isRefreshing = true;
            try {
                // UI proxy endpoint that forwards POST to user-management and uses httpOnly cookie
                const refreshResponse = await fetch('/api/auth/refresh', { method: 'POST' });
                if (refreshResponse.ok) {
                    const data = await refreshResponse.json();
                    accessToken = data.access_token;
                    isRefreshing = false;
                    onRefreshed(accessToken);
                } else {
                    accessToken = null;
                    isRefreshing = false;
                    showLogin();
                    return Promise.reject("Refresh failed");
                }
            } catch (err) {
                isRefreshing = false;
                showLogin();
                return Promise.reject(err);
            }
        }

        // Wait for the token
        const retryOriginalRequest = new Promise((resolve) => {
            subscribeTokenRefresh((token) => {
                options.headers['Authorization'] = `Bearer ${token}`;
                resolve(fetch(url, options));
            });
        });
        return retryOriginalRequest;
    }

    return response;
}

async function login() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        if (response.ok) {
            const data = await response.json();
            accessToken = data.access_token; // Refresh token goes to httpOnly cookie automatically
            showApp();
            loadAcademicYears();
        } else {
            alert("Login failed");
        }
    } catch (error) {
        console.error("Login error", error);
    }
}

async function logout() {
    await fetch('/api/auth/logout', { method: 'POST' });
    accessToken = null;
    showLogin();
}

function showLogin() {
    document.getElementById('login-container').classList.remove('hidden');
    document.getElementById('app').classList.add('hidden');
    document.getElementById('main-header').classList.add('hidden');
}

function showApp() {
    document.getElementById('login-container').classList.add('hidden');
    document.getElementById('app').classList.remove('hidden');
    document.getElementById('main-header').classList.remove('hidden');
}

function switchTab(tab) {
    if (tab === 'home') {
        document.getElementById('home-tab').classList.remove('hidden');
        document.getElementById('reference-tab').classList.add('hidden');
    } else if (tab === 'reference') {
        document.getElementById('home-tab').classList.add('hidden');
        document.getElementById('reference-tab').classList.remove('hidden');
    }
}

async function loadAcademicYears() {
    try {
        const res = await apiCall('/api/syllabus/academic-years');
        if (res.ok) {
            const data = await res.json();
            const select = document.getElementById('academic-year-select');
            const list = document.getElementById('academic-years-list');

            // Clear existing
            while (select.options.length > 1) { select.remove(1); }
            list.innerHTML = "";

            if(data.years) {
                data.years.forEach(yr => {
                    const opt = document.createElement('option');
                    opt.value = yr.id;
                    opt.textContent = yr.name;
                    select.appendChild(opt);

                    const li = document.createElement('li');
                    li.textContent = yr.name;
                    list.appendChild(li);
                });
            }
        }
    } catch (e) {
        console.error("Error loading years", e);
    }
}

async function createAcademicYear() {
    const val = document.getElementById('new-academic-year').value;
    if (!val) return;

    // Use default dates if missing in UI, or backend will fail
    const currentYear = new Date().getFullYear();
    const start_date = `${currentYear}-09-01T00:00:00Z`;
    const end_date = `${currentYear + 1}-06-30T00:00:00Z`;

    try {
        const res = await apiCall('/api/syllabus/academic-years', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify([{ name: val, start_date, end_date }])
        });

        if (res.ok) {
            document.getElementById('new-academic-year').value = '';
            loadAcademicYears();
        } else {
            alert('Failed to create academic year');
        }
    } catch (e) {
        console.error("Error creating year", e);
    }
}
