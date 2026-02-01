import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { Activity, Server, Shield, Plus } from 'lucide-react';
import { useEffect, useState } from 'react';
import { getServices, getRoutes } from './api';
import './index.css';

function ServicesList() {
  const [services, setServices] = useState<any[]>([]);

  useEffect(() => {
    getServices().then(setServices).catch(err => console.error("Failed to fetch services", err));
  }, []);

  return (
    <div className="card">
      <h2 className="mb-4">Services</h2>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Protocol</th>
            <th>Host</th>
            <th>Port</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {services.map(service => (
            <tr key={service.id}>
              <td>{service.name}</td>
              <td>{service.protocol}</td>
              <td>{service.host}</td>
              <td>{service.port}</td>
              <td><span className="status-badge status-active">Enabled</span></td>
            </tr>
          ))}
          {services.length === 0 && <tr><td colSpan={5}>No services found. Run setup_kong.sh</td></tr>}
        </tbody>
      </table>
    </div>
  );
}

function Dashboard() {
  const [services, setServices] = useState<any[]>([]);
  const [routes, setRoutes] = useState<any[]>([]);

  useEffect(() => {
    // Mock data for visual fallback
    const mockServices = [
      { id: '1', name: 'service-public', protocol: 'http', host: 'service-public', port: 8080 },
      { id: '2', name: 'service-private-1', protocol: 'http', host: 'service-private-1', port: 8080 },
    ];
    setServices(mockServices); // Show mock initially

    getServices().then(data => {
      if (data && data.length > 0) setServices(data);
    }).catch(err => console.error(err));

    getRoutes().then(setRoutes).catch(err => console.error(err));
  }, []);

  return (
    <div className="container">
      <div className="grid">
        <div className="card">
          <h2><Server size={20} style={{ marginRight: '8px', verticalAlign: 'bottom' }} /> Total Services</h2>
          <p className="text-3xl font-bold text-blue-600">{services.length}</p>
        </div>
        <div className="card">
          <h2><Activity size={20} style={{ marginRight: '8px', verticalAlign: 'bottom' }} /> Total Routes</h2>
          <p className="text-3xl font-bold text-green-600">{routes.length}</p>
        </div>
        <div className="card">
          <h2><Shield size={20} style={{ marginRight: '8px', verticalAlign: 'bottom' }} /> Security</h2>
          <p className="text-sm">OIDC & Key-Auth Enabled</p>
        </div>
      </div>

      <div style={{ marginTop: '2rem' }}>
        <ServicesList />
      </div>

      <div className="card" style={{ marginTop: '2rem' }}>
        <h3>Quick Actions</h3>
        <p>Use the provided <code>setup_kong.sh</code> script to reset configuration.</p>
      </div>
    </div>
  );
}

function App() {
  return (
    <Router>
      <div className="app">
        <header className="header">
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <Shield size={32} color="var(--primary-color)" />
            <h1 style={{ marginLeft: '1rem' }}>Kong Admin Dashboard</h1>
          </div>
          <nav className="nav">
            <Link to="/">Dashboard</Link>
            <Link to="/services">Services</Link>
          </nav>
        </header>
        <main>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/services" element={<ServicesList />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
