/**
 * Proposal service for backend API integration (fetching and submitting proposals/SOWs)
 */

export interface Proposal {
  id: string;
  title: string;
  client: string;
  created_at: string;
  status: string;
  content: string;
  // Add more fields as needed based on backend response
}

// Fetch all proposals
export async function fetchProposals(): Promise<Proposal[]> {
  const response = await fetch("/api/proposals", {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch proposals: ${response.statusText}`);
  }

  return response.json();
}

// Submit a new proposal
export interface SubmitProposalRequest {
  title: string;
  client: string;
  content: string;
  // Add more fields as needed
}

export async function submitProposal(data: SubmitProposalRequest): Promise<Proposal> {
  const response = await fetch("/api/proposals", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw new Error(`Failed to submit proposal: ${response.statusText}`);
  }

  return response.json();
}