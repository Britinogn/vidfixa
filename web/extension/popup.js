const CONFIG = {
  API_BASE: "https://vidfixa-api.onrender.com/api",
  CREATE_PATH: "/downloads/",
  STATUS_PATH: (id) => `/downloads/${encodeURIComponent(id)}`,
  FILE_TICKET_PATH: (id) => `/downloads/${encodeURIComponent(id)}/file-ticket`,
  FREE_MONTHLY_LIMIT: 15,
};

const urlInput = document.getElementById("url");
const goBtn = document.getElementById("go");
const statusEl = document.getElementById("status");
const resultEl = document.getElementById("result");

function setStatus(msg, isError = false) {
  statusEl.textContent = msg;
  statusEl.className = isError ? "error" : "";
}

chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
  if (tabs[0] && tabs[0].url) urlInput.value = tabs[0].url;
});

async function authHeaders() {
  const { token } = await chrome.storage.local.get("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
}

async function api(path, options = {}) {
  const res = await fetch(CONFIG.API_BASE + path, {
    ...options,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(await authHeaders()),
      ...(options.headers || {}),
    },
  });
  const text = await res.text();
  if (!text.trim()) {
    throw new Error(
      `VidFixa returned an empty response (HTTP ${res.status}). The API must return the created download job with its ID.`,
    );
  }

  let data = null;
  try {
    data = JSON.parse(text);
  } catch {
    if (res.ok) {
      throw new Error(`VidFixa returned a non-JSON response (HTTP ${res.status}).`);
    }
    data = { message: text.trim() };
  }

  if (!res.ok) {
    if (res.status === 429) {
      throw new Error(`Monthly download limit reached (${CONFIG.FREE_MONTHLY_LIMIT} downloads on the Free plan).`);
    }
    if (res.status === 400) {
      throw new Error("This link is invalid or unsupported. Use an Instagram, Facebook, X, LinkedIn, or TikTok video link.");
    }
    if (res.status >= 500) {
      throw new Error("VidFixa is having a server problem. Please try again shortly.");
    }
    throw new Error(data?.error || data?.message || `Request failed (${res.status}).`);
  }
  if (data === null || typeof data !== "object" || Array.isArray(data)) {
    throw new Error(`VidFixa returned an invalid response (HTTP ${res.status}). Please try again.`);
  }
  return data.data && typeof data.data === "object" && !Array.isArray(data.data)
    ? data.data
    : data;
}

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

function isSupportedLink(value) {
  let parsed;
  try {
    parsed = new URL(value);
  } catch {
    return false;
  }
  if (parsed.protocol !== "https:" && parsed.protocol !== "http:") return false;

  const hostname = parsed.hostname.toLowerCase();
  return [
    "instagram.com",
    "facebook.com",
    "fb.watch",
    "x.com",
    "twitter.com",
    "linkedin.com",
    "lnkd.in",
    "tiktok.com",
    "vt.tiktok.com",
  ].some((domain) => hostname === domain || hostname.endsWith(`.${domain}`));
}

async function startDownload() {
  const url = urlInput.value.trim();
  if (!isSupportedLink(url)) {
    resultEl.replaceChildren();
    setStatus("Enter a valid Instagram, Facebook, X, LinkedIn, or TikTok video link.", true);
    return;
  }

  goBtn.disabled = true;
  resultEl.replaceChildren();
  setStatus("Starting download...");

  try {
    const slowRequestTimer = window.setTimeout(() => {
      setStatus("VidFixa is taking a while to wake up. The first request can take up to a minute...");
    }, 5000);
    let job;
    try {
      job = await api(CONFIG.CREATE_PATH, {
        method: "POST",
        body: JSON.stringify({ url }),
      });
    } finally {
      window.clearTimeout(slowRequestTimer);
    }

    const id = job.id;
    if (typeof id !== "string" || !id) {
      throw new Error("VidFixa did not return a download job ID.");
    }

    for (let i = 0; i < 100; i++) {
      const status = await api(CONFIG.STATUS_PATH(id));
      const state = status.status;

      if (state === "completed") {
        const ticket = await api(CONFIG.FILE_TICKET_PATH(id));
        if (typeof ticket.path !== "string" || !ticket.path.startsWith("/downloads/")) {
          throw new Error("VidFixa did not return a valid video link.");
        }

        const link = document.createElement("a");
        link.href = CONFIG.API_BASE + ticket.path;
        link.textContent = "Save video";
        link.download = "";
        setStatus("Done!");
        resultEl.appendChild(link);
        return;
      }
      if (state === "failed") {
        throw new Error(status.error || "VidFixa could not download this video. Please try another link.");
      }
      if (state === "cancelled") {
        throw new Error("This download was cancelled. Please try again.");
      }
      if (state !== "queued" && state !== "processing") {
        throw new Error("VidFixa returned an unknown download status. Please try again.");
      }

      setStatus(`Working... (${state})`);
      await sleep(3000);
    }
    throw new Error("Timed out waiting for the download.");
  } catch (err) {
    const message = err instanceof Error ? err.message : "Something went wrong. Please try again.";
    setStatus(message, true);
  } finally {
    goBtn.disabled = false;
  }
}

goBtn.addEventListener("click", startDownload);