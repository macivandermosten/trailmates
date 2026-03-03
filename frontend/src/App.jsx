export default function App() {
  return (
    <div style={{
      fontFamily: "system-ui, sans-serif",
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "center",
      minHeight: "100vh",
      margin: 0,
      background: "#0f172a",
      color: "#f8fafc",
    }}>
      <h1 style={{ fontSize: "3rem", marginBottom: "0.5rem" }}>TrailMates</h1>
      <p style={{ color: "#94a3b8", fontSize: "1.2rem", marginBottom: "2rem" }}>
        Plan smarter Europe trips. Meet compatible travelers.
      </p>
      <p style={{ color: "#64748b", fontSize: "0.9rem" }}>
        Infrastructure ready — app coming soon.
      </p>
      <a
        href="/api/health"
        style={{ marginTop: "1.5rem", color: "#38bdf8", fontSize: "0.85rem" }}
      >
        API health check →
      </a>
    </div>
  );
}
