import { useState } from 'react';

function App() {
  const [token, setToken] = useState<string>('');
  const [msg, setMsg] = useState<string>('');

  const login = async () => {
    const res = await fetch('http://api:8080/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: 'harryadmin', password: 'StrongPass!2025!' }),
    });
    const data = await res.json();
    setToken(data.access_token);
    setMsg('Login successful – token received');
  };

  return (
    <div style={{ padding: '2rem', fontFamily: 'sans-serif' }}>
      <h1>Orchestrator UI (minimal)</h1>
      <button onClick={login}>Login (get JWT)</button>
      {msg && <p>{msg}</p>}
      {token && (
        <>
          <p>Token: {token.substring(0, 12)}…</p>
          <p>You can now call the API (e.g., <code>curl -H "Authorization: Bearer {token}" http://api:8080/health</code>).</p>
        </>
      )}
    </div>
  );
}

export default App;
