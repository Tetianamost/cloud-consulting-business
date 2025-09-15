import React, { useState, useEffect } from "react";
import { fetchProposals, submitProposal, Proposal, SubmitProposalRequest } from "../../services/proposalService";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { AlertCircle, RefreshCw } from "lucide-react";

const ProposalsPage: React.FC = () => {
  const [proposals, setProposals] = useState<Proposal[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Submission state
  const [title, setTitle] = useState("");
  const [client, setClient] = useState("");
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  useEffect(() => {
    loadProposals();
  }, []);

  const loadProposals = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchProposals();
      setProposals(data);
    } catch (err: any) {
      setError(err.message || "Failed to load proposals");
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setSubmitError(null);
    try {
      const req: SubmitProposalRequest = { title, client, content };
      await submitProposal(req);
      setTitle("");
      setClient("");
      setContent("");
      await loadProposals();
    } catch (err: any) {
      setSubmitError(err.message || "Failed to submit proposal");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <RefreshCw className="w-8 h-8 animate-spin mx-auto mb-4 text-blue-600" />
          <p className="text-gray-600">Loading proposals...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <AlertCircle className="w-8 h-8 mx-auto mb-4 text-red-600" />
          <p className="text-red-600 mb-4">{error}</p>
          <Button onClick={loadProposals} variant="outline">
            <RefreshCw className="w-4 h-4 mr-2" />
            Retry
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Proposals</h1>
          <p className="text-gray-600">Manage and generate client proposals/SOWs</p>
        </div>
        <Button onClick={loadProposals} variant="outline">
          <RefreshCw className="w-4 h-4 mr-2" />
          Refresh
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Submit New Proposal</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              placeholder="Title"
              value={title}
              onChange={e => setTitle(e.target.value)}
              required
            />
            <Input
              placeholder="Client"
              value={client}
              onChange={e => setClient(e.target.value)}
              required
            />
            <Input
              placeholder="Content"
              value={content}
              onChange={e => setContent(e.target.value)}
              required
            />
            <Button type="submit" disabled={submitting}>
              {submitting ? "Submitting..." : "Submit Proposal"}
            </Button>
            {submitError && (
              <p className="text-red-600 text-sm">{submitError}</p>
            )}
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Existing Proposals</CardTitle>
        </CardHeader>
        <CardContent>
          {proposals.length === 0 ? (
            <p className="text-gray-600">No proposals found.</p>
          ) : (
            <ul className="divide-y">
              {proposals.map((proposal) => (
                <li key={proposal.id} className="py-3">
                  <div className="font-semibold">{proposal.title}</div>
                  <div className="text-sm text-gray-500">{proposal.client} &bull; {new Date(proposal.created_at).toLocaleString()}</div>
                  <div className="text-sm">Status: {proposal.status}</div>
                  <div className="text-sm mt-1">{proposal.content.substring(0, 200)}...</div>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default ProposalsPage;