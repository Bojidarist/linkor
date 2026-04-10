const API_BASE = '/admin/api';
let adminKey = '';
let links = [];
let editingId = null;

function init() {
    const params = new URLSearchParams(window.location.search);
    adminKey = params.get('key') || '';
    if (!adminKey) {
        document.getElementById('app').innerHTML =
            '<div class="empty-state"><h2>Unauthorized</h2><p>Missing admin key parameter.</p></div>';
        return;
    }
    loadLinks();
    document.getElementById('btn-create').addEventListener('click', () => openModal());
    document.getElementById('modal-overlay').addEventListener('click', (e) => {
        if (e.target === e.currentTarget) closeModal();
    });
    document.getElementById('link-form').addEventListener('submit', handleSubmit);
    document.getElementById('btn-cancel').addEventListener('click', closeModal);
}

async function apiFetch(path, options = {}) {
    const headers = { 'X-Admin-Key': adminKey, 'Content-Type': 'application/json', ...options.headers };
    const res = await fetch(API_BASE + path, { ...options, headers });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Request failed');
    return data;
}

async function loadLinks() {
    try {
        links = await apiFetch('/links');
        render();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

function render() {
    const totalClicks = links.reduce((s, l) => s + l.clicks, 0);
    const totalUnique = links.reduce((s, l) => s + l.unique_clicks, 0);
    document.getElementById('stat-links').textContent = links.length;
    document.getElementById('stat-clicks').textContent = totalClicks;
    document.getElementById('stat-unique').textContent = totalUnique;

    const tbody = document.getElementById('links-tbody');
    if (links.length === 0) {
        tbody.innerHTML = `<tr><td colspan="6" class="empty-state"><p>No links yet. Create one to get started.</p></td></tr>`;
        return;
    }

    tbody.innerHTML = links.map(link => `
        <tr>
            <td>${esc(link.name)}</td>
            <td class="url-cell">/${esc(link.short_url)}</td>
            <td class="target-url-cell" title="${esc(link.target_url)}">${esc(link.target_url)}</td>
            <td class="click-count">${link.clicks}</td>
            <td class="click-count">${link.unique_clicks}</td>
            <td class="actions">
                <button class="btn btn-ghost btn-sm" onclick="openModal(${link.id})">Edit</button>
                <button class="btn btn-danger btn-sm" onclick="deleteLink(${link.id})">Delete</button>
            </td>
        </tr>
    `).join('');
}

function openModal(id = null) {
    editingId = id;
    const modal = document.getElementById('modal-overlay');
    const title = document.getElementById('modal-title');
    const form = document.getElementById('link-form');
    form.reset();

    if (id !== null) {
        const link = links.find(l => l.id === id);
        if (!link) return;
        title.textContent = 'Edit Link';
        document.getElementById('field-name').value = link.name;
        document.getElementById('field-short-url').value = link.short_url;
        document.getElementById('field-target-url').value = link.target_url;
    } else {
        title.textContent = 'Create Link';
    }
    modal.classList.add('active');
    document.getElementById('field-name').focus();
}

function closeModal() {
    document.getElementById('modal-overlay').classList.remove('active');
    editingId = null;
}

async function handleSubmit(e) {
    e.preventDefault();
    const payload = {
        name: document.getElementById('field-name').value,
        short_url: document.getElementById('field-short-url').value,
        target_url: document.getElementById('field-target-url').value,
    };

    try {
        if (editingId !== null) {
            await apiFetch(`/links/${editingId}`, { method: 'PUT', body: JSON.stringify(payload) });
            showToast('Link updated', 'success');
        } else {
            await apiFetch('/links', { method: 'POST', body: JSON.stringify(payload) });
            showToast('Link created', 'success');
        }
        closeModal();
        await loadLinks();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function deleteLink(id) {
    if (!confirm('Delete this link? This action cannot be undone.')) return;
    try {
        await apiFetch(`/links/${id}`, { method: 'DELETE' });
        showToast('Link deleted', 'success');
        await loadLinks();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

function showToast(message, type) {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = `toast ${type} visible`;
    setTimeout(() => toast.classList.remove('visible'), 3000);
}

function esc(str) {
    const el = document.createElement('span');
    el.textContent = str;
    return el.innerHTML;
}

document.addEventListener('DOMContentLoaded', init);
