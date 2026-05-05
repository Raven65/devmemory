// DevMemory Web UI
const API = '';

// --- View Switching ---
document.querySelectorAll('.nav-btn').forEach(btn => {
  btn.addEventListener('click', () => {
    document.querySelectorAll('.nav-btn').forEach(b => b.classList.remove('active'));
    document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));
    btn.classList.add('active');
    document.getElementById('view-' + btn.dataset.view).classList.add('active');
    if (btn.dataset.view === 'today') loadToday();
  });
});

// --- Capture ---
document.getElementById('capture-submit').addEventListener('click', async () => {
  const content = document.getElementById('capture-content').value.trim();
  if (!content) {
    showMsg('capture-msg', 'Content is required', 'error');
    return;
  }

  const entry = {
    content,
    type: document.getElementById('capture-type').value || undefined,
    title: document.getElementById('capture-title').value.trim() || undefined,
    project: document.getElementById('capture-project').value.trim() || undefined,
    tags: parseTags(document.getElementById('capture-tags').value),
  };

  try {
    const res = await fetch(API + '/api/entries', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(entry),
    });
    if (!res.ok) throw new Error(await res.text());
    const created = await res.json();
    showMsg('capture-msg', `Created ${created.type} entry: ${created.id.slice(0, 8)}`, 'success');
    // Clear form
    document.getElementById('capture-content').value = '';
    document.getElementById('capture-title').value = '';
    document.getElementById('capture-project').value = '';
    document.getElementById('capture-tags').value = '';
    document.getElementById('capture-type').value = '';
  } catch (e) {
    showMsg('capture-msg', 'Error: ' + e.message, 'error');
  }
});

function parseTags(s) {
  if (!s.trim()) return undefined;
  return s.split(',').map(t => t.trim()).filter(t => t);
}

// --- Search ---
document.getElementById('search-btn').addEventListener('click', doSearch);
document.getElementById('search-input').addEventListener('keydown', e => {
  if (e.key === 'Enter') doSearch();
});

async function doSearch() {
  const q = document.getElementById('search-input').value.trim();
  if (!q) return;

  const type = document.getElementById('search-type').value;
  let url = API + '/api/search?q=' + encodeURIComponent(q);
  if (type) url += '&type=' + type;

  try {
    const res = await fetch(url);
    if (!res.ok) throw new Error(await res.text());
    const results = await res.json();
    renderSearchResults(results);
  } catch (e) {
    document.getElementById('search-results').innerHTML =
      `<div class="empty-state">Error: ${e.message}</div>`;
  }
}

function renderSearchResults(results) {
  const container = document.getElementById('search-results');
  if (results.length === 0) {
    container.innerHTML = '<div class="empty-state">No results found</div>';
    return;
  }
  container.innerHTML = results.map(r => entryCard(r.entry, r.score)).join('');
  bindCardClicks(container);
}

// --- Today ---
async function loadToday() {
  try {
    const res = await fetch(API + '/api/today');
    if (!res.ok) throw new Error(await res.text());
    const data = await res.json();
    document.getElementById('today-date').textContent = data.date;
    renderTodayEntries(data.entries);
  } catch (e) {
    document.getElementById('today-entries').innerHTML =
      `<div class="empty-state">Error: ${e.message}</div>`;
  }
}

document.getElementById('today-refresh').addEventListener('click', loadToday);

function renderTodayEntries(entries) {
  const container = document.getElementById('today-entries');
  if (entries.length === 0) {
    container.innerHTML = '<div class="empty-state">(no entries yet)</div>';
    return;
  }
  container.innerHTML = entries.map(e => entryCard(e)).join('');
  bindCardClicks(container);
}

// --- Entry Card ---
function entryCard(entry, score) {
  const time = new Date(entry.created_at).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
  const title = escapeHtml(entry.title || truncate(entry.content, 80));
  const tags = (entry.tags || []).map(t => `<span class="tag">${escapeHtml(t)}</span>`).join('');
  const danger = entry.dangerous ? '<span class="badge-danger">DANGEROUS</span>' : '';
  const fav = entry.favorite ? '<span class="badge-fav">&#9733;</span>' : '';
  const scoreHtml = score ? `<span class="entry-score">${score.toFixed(1)}</span>` : '';

  return `
    <div class="entry-card" data-id="${entry.id}">
      <div class="entry-header">
        <span class="entry-type">${entry.type}${danger}${fav}</span>
        <span class="entry-time">${time}</span>
      </div>
      <div class="entry-title">${title}${scoreHtml}</div>
      ${entry.project ? `<div class="entry-project">${escapeHtml(entry.project)}</div>` : ''}
      ${tags ? `<div class="entry-tags">${tags}</div>` : ''}
    </div>`;
}

function bindCardClicks(container) {
  container.querySelectorAll('.entry-card').forEach(card => {
    card.addEventListener('click', () => showModal(card.dataset.id));
  });
}

// --- Modal ---
let currentEntry = null;

document.getElementById('modal-close').addEventListener('click', hideModal);
document.getElementById('modal').addEventListener('click', e => {
  if (e.target.id === 'modal') hideModal();
});

async function showModal(id) {
  try {
    const res = await fetch(API + '/api/entries/' + id);
    if (!res.ok) throw new Error('not found');
    currentEntry = await res.json();
    renderModal(currentEntry);
    document.getElementById('modal').classList.remove('hidden');
  } catch (e) {
    alert('Error loading entry: ' + e.message);
  }
}

function renderModal(e) {
  const lastUsed = e.last_used_at ? new Date(e.last_used_at).toLocaleString() : 'Never';
  document.getElementById('modal-title').textContent =
    `${e.type} — ${e.id.slice(0, 8)}${e.dangerous ? ' <span class="badge-danger">DANGEROUS</span>' : ''}${e.favorite ? ' <span class="badge-fav">&#9733;</span>' : ''}`;

  document.getElementById('modal-body').innerHTML = `
    ${e.title ? `<div class="field"><div class="field-label">Title</div><div class="field-value">${escapeHtml(e.title)}</div></div>` : ''}
    <div class="field"><div class="field-label">Content</div><div class="field-value content">${escapeHtml(e.content)}</div></div>
    ${e.summary ? `<div class="field"><div class="field-label">Summary</div><div class="field-value">${escapeHtml(e.summary)}</div></div>` : ''}
    ${e.project ? `<div class="field"><div class="field-label">Project</div><div class="field-value">${escapeHtml(e.project)}</div></div>` : ''}
    ${e.tags && e.tags.length ? `<div class="field"><div class="field-label">Tags</div><div class="field-value">${e.tags.map(t => `<span class="tag">${escapeHtml(t)}</span>`).join(' ')}</div></div>` : ''}
    <div class="field"><div class="field-label">Stats</div><div class="field-value">Used ${e.use_count || 0} times &middot; Last used: ${lastUsed}</div></div>
    <div class="field"><div class="field-label">Created</div><div class="field-value">${new Date(e.created_at).toLocaleString()}</div></div>
    <div class="field"><div class="field-label">Updated</div><div class="field-value">${new Date(e.updated_at).toLocaleString()}</div></div>`;

  document.getElementById('modal-toggle-fav').textContent = e.favorite ? 'Unfavorite' : 'Favorite';
  document.getElementById('modal-archive').textContent = e.archived ? 'Unarchive' : 'Archive';
  document.getElementById('modal-actions-view').classList.remove('hidden');
  document.getElementById('modal-actions-edit').classList.add('hidden');
}

function hideModal() {
  document.getElementById('modal').classList.add('hidden');
  currentEntry = null;
}

// Modal actions
document.getElementById('modal-copy').addEventListener('click', async () => {
  if (!currentEntry) return;
  try {
    await fetch(API + '/api/entries/' + currentEntry.id + '/copy', { method: 'POST' });
    await navigator.clipboard.writeText(currentEntry.content);
    alert('Copied to clipboard!');
  } catch (e) {
    alert('Copy failed: ' + e.message);
  }
});

document.getElementById('modal-edit').addEventListener('click', () => {
  if (!currentEntry) return;
  renderEditModal(currentEntry);
});

document.getElementById('modal-edit-cancel').addEventListener('click', () => {
  if (!currentEntry) return;
  renderModal(currentEntry);
  document.getElementById('modal-actions-edit').classList.add('hidden');
  document.getElementById('modal-actions-view').classList.remove('hidden');
});

document.getElementById('modal-edit-save').addEventListener('click', async () => {
  if (!currentEntry) return;
  const fields = {
    title: document.getElementById('edit-title').value.trim(),
    content: document.getElementById('edit-content').value.trim(),
    type: document.getElementById('edit-type').value,
    project: document.getElementById('edit-project').value.trim(),
    tags: parseTags(document.getElementById('edit-tags').value),
  };
  if (!fields.content) {
    alert('Content cannot be empty');
    return;
  }
  try {
    await updateEntry(currentEntry.id, fields);
    showModal(currentEntry.id);
  } catch (e) {
    alert('Save failed: ' + e.message);
  }
});

function renderEditModal(e) {
  document.getElementById('modal-title').textContent =
    `Edit — ${e.type} ${e.id.slice(0, 8)}`;

  document.getElementById('modal-body').innerHTML = `
    <div class="field">
      <div class="field-label">Title</div>
      <input type="text" id="edit-title" value="${escapeAttr(e.title || '')}">
    </div>
    <div class="field">
      <div class="field-label">Content</div>
      <textarea id="edit-content" rows="4">${escapeHtml(e.content)}</textarea>
    </div>
    <div class="field">
      <div class="field-label">Type</div>
      <select id="edit-type">
        <option value="command"${e.type==='command'?' selected':''}>Command</option>
        <option value="url"${e.type==='url'?' selected':''}>URL</option>
        <option value="snippet"${e.type==='snippet'?' selected':''}>Snippet</option>
        <option value="prompt"${e.type==='prompt'?' selected':''}>Prompt</option>
        <option value="note"${e.type==='note'?' selected':''}>Note</option>
        <option value="issue"${e.type==='issue'?' selected':''}>Issue</option>
        <option value="journal"${e.type==='journal'?' selected':''}>Journal</option>
        <option value="business"${e.type==='business'?' selected':''}>Business</option>
        <option value="task"${e.type==='task'?' selected':''}>Task</option>
        <option value="file"${e.type==='file'?' selected':''}>File</option>
        <option value="folder"${e.type==='folder'?' selected':''}>Folder</option>
      </select>
    </div>
    <div class="field">
      <div class="field-label">Project</div>
      <input type="text" id="edit-project" value="${escapeAttr(e.project || '')}">
    </div>
    <div class="field">
      <div class="field-label">Tags</div>
      <input type="text" id="edit-tags" value="${escapeAttr((e.tags||[]).join(', '))}">
    </div>`;

  document.getElementById('modal-actions-view').classList.add('hidden');
  document.getElementById('modal-actions-edit').classList.remove('hidden');
}

function escapeAttr(s) {
  if (!s) return '';
  return s.replace(/&/g,'&amp;').replace(/"/g,'&quot;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

document.getElementById('modal-toggle-fav').addEventListener('click', async () => {
  if (!currentEntry) return;
  await updateEntry(currentEntry.id, { favorite: !currentEntry.favorite });
  showModal(currentEntry.id);
});

document.getElementById('modal-archive').addEventListener('click', async () => {
  if (!currentEntry) return;
  await updateEntry(currentEntry.id, { archived: !currentEntry.archived });
  showModal(currentEntry.id);
});

document.getElementById('modal-delete').addEventListener('click', async () => {
  if (!currentEntry) return;
  if (!confirm('Delete this entry?')) return;
  await fetch(API + '/api/entries/' + currentEntry.id, { method: 'DELETE' });
  hideModal();
  // Refresh current view
  const activeView = document.querySelector('.nav-btn.active').dataset.view;
  if (activeView === 'today') loadToday();
  else if (activeView === 'search') doSearch();
});

async function updateEntry(id, fields) {
  await fetch(API + '/api/entries/' + id, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(fields),
  });
}

// --- Helpers ---
function showMsg(id, text, type) {
  const el = document.getElementById(id);
  el.textContent = text;
  el.className = 'msg ' + type;
  setTimeout(() => { el.textContent = ''; el.className = 'msg'; }, 3000);
}

function escapeHtml(s) {
  if (!s) return '';
  const div = document.createElement('div');
  div.textContent = s;
  return div.innerHTML;
}

function truncate(s, n) {
  if (!s) return '';
  const first = s.split('\n')[0];
  return first.length > n ? first.slice(0, n) + '...' : first;
}
