let backends = [];
let requests = [];
let rrIndex = 0;

// Canvas
const canvas = document.getElementById('simulation');
const ctx = canvas.getContext('2d');
ctx.font = '14px Arial';

const METRICS_URL = 'http://localhost:8080/metrics';
const LOAD_BALANCER_X = 50;
const LOAD_BALANCER_Y = canvas.height / 2;

// ---------------- BACKEND POSITIONING ----------------
function setBackendPositions() {
  const n = backends.length;
  backends.forEach((b, i) => {
    b.x = 600;
    b.y = (canvas.height / (n + 1)) * (i + 1) - 25;
    b.vizX = b.x + 60;
    b.vizY = b.y + 25;
  });
}

// ---------------- DRAW LOOP ----------------
function draw() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // Load Balancer
  ctx.fillStyle = '#3b82f6';
  ctx.fillRect(20, LOAD_BALANCER_Y - 30, 80, 60);
  ctx.fillStyle = '#fff';
  ctx.fillText('LB', 50, LOAD_BALANCER_Y + 5);

  // Backends
  backends.forEach(b => {
    ctx.fillStyle = b.Alive ? '#22c55e' : '#ef4444';
    ctx.fillRect(b.x, b.y, 140, 60);

    ctx.fillStyle = '#fff';
    ctx.fillText(b.Address, b.x + 5, b.y + 18);
    ctx.fillText(`Conn: ${b.ActiveConns}`, b.x + 5, b.y + 34);
    ctx.fillText(`Latency: ${b.Latency}µs`, b.x + 5, b.y + 50);
    ctx.fillText(`Errors: ${b.ErrorCount}`, b.x + 90, b.y + 34);
  });

  // Requests
  requests.forEach(r => {
    ctx.fillStyle = '#facc15';
    ctx.beginPath();
    ctx.arc(r.x, r.y, 6, 0, Math.PI * 2);
    ctx.fill();

    r.x += 6;
    r.y += (r.targetY - r.y) * 0.1;
  });

  // Cleanup
  requests = requests.filter(r => {
    if (r.x >= r.targetX) {
      r.backend.ActiveConns = Math.max(0, r.backend.ActiveConns - 1);
      return false;
    }
    return true;
  });

  requestAnimationFrame(draw);
}

// ---------------- LOAD BALANCER LOGIC ----------------
function pickBackend() {
  const alive = backends.filter(b => b.Alive);
  if (!alive.length) return null;

  const algo = document.getElementById('algo').value;

  if (algo === 'roundrobin') {
    const b = alive[rrIndex % alive.length];
    rrIndex++;
    return b;
  }

  if (algo === 'leastconnections') {
    return alive.reduce((a, b) =>
      a.ActiveConns <= b.ActiveConns ? a : b
    );
  }

  return alive[Math.floor(Math.random() * alive.length)];
}

// ---------------- REQUEST SIMULATION ----------------
function sendRequest() {
  const backend = pickBackend();
  if (!backend) return;

  backend.ActiveConns++;

  requests.push({
    x: LOAD_BALANCER_X + 80,
    y: LOAD_BALANCER_Y,
    targetX: backend.vizX,
    targetY: backend.vizY,
    backend
  });
}

// ---------------- METRICS FETCH ----------------
async function fetchMetrics() {
  try {
    const res = await fetch(METRICS_URL);
    const data = await res.json();

    // First-time init
    if (!backends.length) {
      backends = data.map(b => ({ ...b }));
      setBackendPositions();
      return;
    }

    // Update existing backends by Address
    data.forEach(metricBackend => {
      const local = backends.find(
        b => b.Address === metricBackend.Address
      );
      if (!local) return;

      local.Alive = metricBackend.Alive;
      local.Latency = metricBackend.Latency;
      local.ErrorCount = metricBackend.ErrorCount;
    });

    // Simulate traffic
    if (Math.random() > 0.4) sendRequest();

  } catch (err) {
    console.error('Metrics fetch failed:', err);
  }
}

// ---------------- START ----------------
setInterval(fetchMetrics, 1000);
draw();
