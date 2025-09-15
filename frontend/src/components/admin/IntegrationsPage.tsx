import React, { useEffect, useState } from "react";
import { Integration, fetchIntegrations } from "../../services/integrationsService";

const IntegrationsPage: React.FC = () => {
  const [integrations, setIntegrations] = useState<Integration[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    fetchIntegrations()
      .then((res) => {
        setIntegrations(res.data);
        setError(null);
      })
      .catch((err) => {
        setError(err.message || "Failed to load integrations");
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <h1>Integrations</h1>
      {loading && <p>Loading integrations...</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}
      {!loading && !error && (
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Status</th>
              <th>Description</th>
              <th>Last Sync</th>
            </tr>
          </thead>
          <tbody>
            {integrations.map((integration) => (
              <tr key={integration.id}>
                <td>{integration.name}</td>
                <td>{integration.type}</td>
                <td>{integration.status}</td>
                <td>{integration.description || "-"}</td>
                <td>{integration.lastSync ? new Date(integration.lastSync).toLocaleString() : "-"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default IntegrationsPage;