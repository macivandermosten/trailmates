import { useState, useEffect } from "react";

const API = "/api";

async function apiFetch(path, options = {}) {
  const token = localStorage.getItem("token");
  const headers = { "Content-Type": "application/json", ...options.headers };
  if (token) headers["Authorization"] = `Bearer ${token}`;
  const res = await fetch(`${API}${path}`, { ...options, headers });
  if (res.status === 204) return null;
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || "Request failed");
  return data;
}

// ---- Login Page ----
function LoginPage({ onLogin }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isRegister, setIsRegister] = useState(false);
  const [displayName, setDisplayName] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    try {
      const endpoint = isRegister ? "/register" : "/login";
      const body = isRegister
        ? { email, password, display_name: displayName }
        : { email, password };
      const data = await apiFetch(endpoint, {
        method: "POST",
        body: JSON.stringify(body),
      });
      localStorage.setItem("token", data.token);
      onLogin(data.user);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div style={styles.page}>
      <div style={styles.card}>
        <h1 style={{ fontSize: "2rem", marginBottom: "0.25rem" }}>TrailMates</h1>
        <p style={{ color: "#94a3b8", marginBottom: "1.5rem" }}>
          Plan smarter Europe trips. Meet compatible travelers.
        </p>
        <h2 style={{ fontSize: "1.25rem", marginBottom: "1rem" }}>
          {isRegister ? "Create Account" : "Log In"}
        </h2>
        <form onSubmit={handleSubmit} style={styles.form}>
          {isRegister && (
            <input
              style={styles.input}
              placeholder="Display Name"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
              required
            />
          )}
          <input
            style={styles.input}
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            style={styles.input}
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {error && <p style={{ color: "#ef4444", fontSize: "0.875rem" }}>{error}</p>}
          <button style={styles.button} type="submit">
            {isRegister ? "Sign Up" : "Log In"}
          </button>
        </form>
        <p style={{ marginTop: "1rem", fontSize: "0.875rem", color: "#94a3b8" }}>
          {isRegister ? "Already have an account?" : "Don't have an account?"}{" "}
          <span
            style={{ color: "#38bdf8", cursor: "pointer" }}
            onClick={() => { setIsRegister(!isRegister); setError(""); }}
          >
            {isRegister ? "Log in" : "Sign up"}
          </span>
        </p>
      </div>
    </div>
  );
}

// ---- Cities Browser ----
function CitiesPage() {
  const [cities, setCities] = useState([]);
  const [selected, setSelected] = useState(null);
  const [attractions, setAttractions] = useState([]);

  useEffect(() => {
    apiFetch("/cities").then(setCities).catch(() => {});
  }, []);

  const selectCity = async (city) => {
    setSelected(city);
    const data = await apiFetch(`/cities/${city.id}`);
    setAttractions(data.attractions || []);
  };

  if (selected) {
    return (
      <div>
        <button style={styles.backBtn} onClick={() => setSelected(null)}>
          ← Back to Cities
        </button>
        <h2 style={{ fontSize: "1.5rem" }}>{selected.name}, {selected.country}</h2>
        <p style={{ color: "#94a3b8", marginBottom: "1rem" }}>{selected.description}</p>
        <div style={{ display: "flex", gap: "0.5rem", marginBottom: "1rem" }}>
          <span style={styles.badge}>{selected.cost_level} cost</span>
          <span style={styles.badge}>{selected.region}</span>
        </div>
        <h3 style={{ marginBottom: "0.75rem" }}>Top Attractions ({attractions.length})</h3>
        {attractions.map((a) => (
          <div key={a.id} style={styles.listItem}>
            <div style={{ fontWeight: 600 }}>{a.name}</div>
            <div style={{ color: "#94a3b8", fontSize: "0.875rem" }}>{a.description}</div>
            <div style={{ display: "flex", gap: "0.5rem", marginTop: "0.25rem" }}>
              <span style={styles.smallBadge}>{a.category}</span>
              <span style={styles.smallBadge}>{a.estimated_hours}h</span>
              <span style={styles.smallBadge}>{a.cost_level}</span>
            </div>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div>
      <h2 style={{ fontSize: "1.5rem", marginBottom: "1rem" }}>Explore Cities ({cities.length})</h2>
      <div style={styles.grid}>
        {cities.map((c) => (
          <div key={c.id} style={styles.cityCard} onClick={() => selectCity(c)}>
            <div style={{ fontWeight: 600, fontSize: "1.1rem" }}>{c.name}</div>
            <div style={{ color: "#94a3b8", fontSize: "0.875rem" }}>{c.country}</div>
            <div style={{ display: "flex", gap: "0.5rem", marginTop: "0.5rem" }}>
              <span style={styles.smallBadge}>{c.cost_level}</span>
              <span style={styles.smallBadge}>{c.region}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// ---- Trips Dashboard ----
function TripsPage() {
  const [trips, setTrips] = useState([]);
  const [selectedTrip, setSelectedTrip] = useState(null);
  const [route, setRoute] = useState(null);

  useEffect(() => {
    apiFetch("/trips").then(setTrips).catch(() => {});
  }, []);

  const viewTrip = async (trip) => {
    const detail = await apiFetch(`/trips/${trip.id}`);
    setSelectedTrip(detail);
    try {
      const r = await apiFetch(`/trips/${trip.id}/route`);
      setRoute(r);
    } catch { setRoute(null); }
  };

  if (selectedTrip) {
    return (
      <div>
        <button style={styles.backBtn} onClick={() => { setSelectedTrip(null); setRoute(null); }}>
          ← Back to Trips
        </button>
        <h2 style={{ fontSize: "1.5rem" }}>{selectedTrip.name}</h2>
        <div style={{ display: "flex", gap: "0.5rem", marginBottom: "1rem" }}>
          <span style={styles.badge}>{selectedTrip.status}</span>
          <span style={styles.badge}>{selectedTrip.budget_style}</span>
          {selectedTrip.start_date && (
            <span style={styles.badge}>{selectedTrip.start_date} → {selectedTrip.end_date}</span>
          )}
        </div>
        <h3>Cities ({selectedTrip.cities.length})</h3>
        {selectedTrip.cities.map((c, i) => (
          <div key={c.id} style={styles.listItem}>
            <span style={{ fontWeight: 600 }}>#{i + 1} {c.city_name}, {c.country}</span>
            {c.arrival_date && (
              <span style={{ color: "#94a3b8", marginLeft: "0.5rem", fontSize: "0.875rem" }}>
                {c.arrival_date} → {c.departure_date}
              </span>
            )}
          </div>
        ))}
        {route && route.hops && route.hops.length > 0 && (
          <>
            <h3 style={{ marginTop: "1rem" }}>Optimized Route</h3>
            {route.hops.map((h, i) => (
              <div key={i} style={styles.listItem}>
                {h.from} → {h.to}: {h.distance_km} km
                <span style={{ ...styles.smallBadge, marginLeft: "0.5rem" }}>{h.cost_level}</span>
              </div>
            ))}
          </>
        )}
      </div>
    );
  }

  return (
    <div>
      <h2 style={{ fontSize: "1.5rem", marginBottom: "1rem" }}>My Trips ({trips.length})</h2>
      {trips.length === 0 && <p style={{ color: "#94a3b8" }}>No trips yet.</p>}
      {trips.map((t) => (
        <div key={t.id} style={styles.cityCard} onClick={() => viewTrip(t)}>
          <div style={{ fontWeight: 600 }}>{t.name}</div>
          <div style={{ display: "flex", gap: "0.5rem", marginTop: "0.5rem" }}>
            <span style={styles.smallBadge}>{t.status}</span>
            <span style={styles.smallBadge}>{t.budget_style}</span>
            {t.start_date && <span style={styles.smallBadge}>{t.start_date}</span>}
          </div>
        </div>
      ))}
    </div>
  );
}

// ---- Profile Page ----
function ProfilePage() {
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    apiFetch("/profile").then(setProfile).catch(() => {});
  }, []);

  if (!profile) return <p>Loading...</p>;

  return (
    <div>
      <h2 style={{ fontSize: "1.5rem", marginBottom: "1rem" }}>Profile</h2>
      <div style={styles.card}>
        <div style={{ fontSize: "1.25rem", fontWeight: 600 }}>{profile.display_name}</div>
        <p style={{ color: "#94a3b8", margin: "0.5rem 0" }}>{profile.bio}</p>
        <div style={{ display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
          <span style={styles.badge}>{profile.travel_style}</span>
          {profile.interests?.map((i) => (
            <span key={i} style={styles.smallBadge}>{i}</span>
          ))}
        </div>
      </div>
    </div>
  );
}

// ---- Connections Page ----
function ConnectionsPage() {
  const [connections, setConnections] = useState([]);

  useEffect(() => {
    apiFetch("/connections").then(setConnections).catch(() => {});
  }, []);

  return (
    <div>
      <h2 style={{ fontSize: "1.5rem", marginBottom: "1rem" }}>Connections ({connections.length})</h2>
      {connections.length === 0 && <p style={{ color: "#94a3b8" }}>No connections yet.</p>}
      {connections.map((c) => (
        <div key={c.id} style={styles.listItem}>
          <div>
            <span style={{ fontWeight: 600 }}>{c.requester_name}</span>
            <span style={{ color: "#94a3b8" }}> → </span>
            <span style={{ fontWeight: 600 }}>{c.recipient_name}</span>
          </div>
          <div style={{ display: "flex", gap: "0.5rem", marginTop: "0.25rem" }}>
            <span style={styles.smallBadge}>{c.status}</span>
            <span style={styles.smallBadge}>{c.trip_name}</span>
          </div>
          {c.message && <p style={{ color: "#94a3b8", fontSize: "0.875rem", marginTop: "0.25rem" }}>{c.message}</p>}
        </div>
      ))}
    </div>
  );
}

// ---- Main App ----
export default function App() {
  const [user, setUser] = useState(null);
  const [page, setPage] = useState("cities");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      apiFetch("/auth/verify")
        .then((data) => setUser(data.user))
        .catch(() => localStorage.removeItem("token"))
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
    setPage("cities");
  };

  if (loading) return <div style={styles.page}><p>Loading...</p></div>;
  if (!user) return <LoginPage onLogin={(u) => { setUser(u); setPage("cities"); }} />;

  const navItems = [
    { key: "cities", label: "Cities" },
    { key: "trips", label: "My Trips" },
    { key: "connections", label: "Connections" },
    { key: "profile", label: "Profile" },
  ];

  return (
    <div style={{ minHeight: "100vh", background: "#0f172a", color: "#f8fafc" }}>
      <nav style={styles.nav}>
        <span style={{ fontWeight: 700, fontSize: "1.25rem" }}>TrailMates</span>
        <div style={{ display: "flex", gap: "1rem", alignItems: "center" }}>
          {navItems.map((item) => (
            <span
              key={item.key}
              onClick={() => setPage(item.key)}
              style={{
                cursor: "pointer",
                color: page === item.key ? "#38bdf8" : "#94a3b8",
                fontWeight: page === item.key ? 600 : 400,
              }}
            >
              {item.label}
            </span>
          ))}
          <span style={{ color: "#64748b" }}>|</span>
          <span style={{ color: "#94a3b8", fontSize: "0.875rem" }}>{user.display_name}</span>
          <span onClick={logout} style={{ cursor: "pointer", color: "#ef4444", fontSize: "0.875rem" }}>
            Logout
          </span>
        </div>
      </nav>
      <main style={{ maxWidth: 900, margin: "0 auto", padding: "2rem 1rem" }}>
        {page === "cities" && <CitiesPage />}
        {page === "trips" && <TripsPage />}
        {page === "connections" && <ConnectionsPage />}
        {page === "profile" && <ProfilePage />}
      </main>
    </div>
  );
}

// ---- Styles ----
const styles = {
  page: {
    fontFamily: "system-ui, sans-serif",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    minHeight: "100vh",
    background: "#0f172a",
    color: "#f8fafc",
  },
  card: {
    background: "#1e293b",
    borderRadius: "0.75rem",
    padding: "2rem",
    maxWidth: 420,
    width: "100%",
  },
  form: {
    display: "flex",
    flexDirection: "column",
    gap: "0.75rem",
  },
  input: {
    padding: "0.625rem 0.75rem",
    borderRadius: "0.5rem",
    border: "1px solid #334155",
    background: "#0f172a",
    color: "#f8fafc",
    fontSize: "1rem",
  },
  button: {
    padding: "0.625rem",
    borderRadius: "0.5rem",
    border: "none",
    background: "#2563eb",
    color: "#fff",
    fontSize: "1rem",
    fontWeight: 600,
    cursor: "pointer",
  },
  nav: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    padding: "1rem 2rem",
    borderBottom: "1px solid #1e293b",
    background: "#0f172a",
  },
  grid: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fill, minmax(220px, 1fr))",
    gap: "1rem",
  },
  cityCard: {
    background: "#1e293b",
    borderRadius: "0.75rem",
    padding: "1rem",
    cursor: "pointer",
    transition: "background 0.15s",
  },
  listItem: {
    background: "#1e293b",
    borderRadius: "0.5rem",
    padding: "0.75rem 1rem",
    marginBottom: "0.5rem",
  },
  badge: {
    background: "#334155",
    padding: "0.25rem 0.625rem",
    borderRadius: "1rem",
    fontSize: "0.8rem",
    color: "#cbd5e1",
  },
  smallBadge: {
    background: "#1e293b",
    border: "1px solid #334155",
    padding: "0.125rem 0.5rem",
    borderRadius: "0.75rem",
    fontSize: "0.75rem",
    color: "#94a3b8",
  },
  backBtn: {
    background: "none",
    border: "none",
    color: "#38bdf8",
    cursor: "pointer",
    fontSize: "0.9rem",
    padding: 0,
    marginBottom: "1rem",
  },
};
