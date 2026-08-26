const sampleList = document.querySelector('#sample-list');
const sampleMessage = document.querySelector('#sample-message');
const form = document.querySelector('#sample-form');

function setText(id, value) { document.querySelector(id).textContent = value; }

async function getJSON(path) {
  const response = await fetch(path);
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

function renderSamples(samples) {
  sampleList.replaceChildren();
  if (!samples.length) {
    const empty = document.createElement('div');
    empty.className = 'empty';
    empty.textContent = '还没有样品，先从左侧登记一块。';
    sampleList.append(empty);
    return;
  }
  samples.slice().reverse().forEach((sample) => {
    const item = document.createElement('article');
    item.className = 'sample';
    item.innerHTML = `<div><strong>${sample.id}</strong><br><small>${sample.location} · ${sample.material}</small></div><span class="tag">${sample.status}</span>`;
    sampleList.append(item);
  });
}

async function refresh() {
  sampleMessage.textContent = '读取中';
  try {
    const [samples, batches, metrics] = await Promise.all([getJSON('/api/samples'), getJSON('/api/batches'), getJSON('/api/metrics')]);
    renderSamples(samples);
    setText('#sample-count', samples.length);
    setText('#batch-count', batches.length);
    setText('#observation-count', metrics['observations.recorded'] || 0);
    setText('#report-count', metrics['reports.generated'] || 0);
    sampleMessage.textContent = '已同步';
  } catch (error) {
    sampleMessage.textContent = '服务不可用';
  }
}

form.addEventListener('submit', async (event) => {
  event.preventDefault();
  const data = Object.fromEntries(new FormData(form));
  data.depth_mm = Number(data.depth_mm);
  data.collection_date = new Date(data.collection_date).toISOString();
  const response = await fetch('/api/samples', { method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(data) });
  if (!response.ok) { alert('样品保存失败'); return; }
  form.reset();
  await refresh();
});

document.querySelector('#refresh').addEventListener('click', refresh);
refresh();
