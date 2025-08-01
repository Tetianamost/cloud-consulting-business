// Service for Meeting Preparation API integration

export interface PreMeetingBriefing {
  id: string;
  meetingId: string;
  title: string;
  summary: string;
  recommendations: string[];
  createdAt: string;
}

export interface QuestionBank {
  id: string;
  meetingId: string;
  questions: string[];
  category: string;
  createdAt: string;
}

export interface CompetitiveAnalysis {
  id: string;
  meetingId: string;
  competitors: string[];
  strengths: string[];
  weaknesses: string[];
  opportunities: string[];
  threats: string[];
  summary: string;
  createdAt: string;
}

export interface FollowUpActionItem {
  id: string;
  meetingId: string;
  description: string;
  assignedTo: string;
  dueDate: string;
  status: "pending" | "completed" | "in_progress";
  createdAt: string;
}

// Use apiService for admin endpoints
import apiService from "./api";

export async function fetchPreMeetingBriefings(meetingId: string): Promise<PreMeetingBriefing[]> {
  const query = `?meetingId=${encodeURIComponent(meetingId)}`;
  const res = await apiService.getPreMeetingBriefings(query);
  return res.data as PreMeetingBriefing[];
}

// Use apiService for admin endpoints
// Use apiService for admin endpoints
export async function fetchQuestionBanks(meetingId: string): Promise<QuestionBank[]> {
  const query = `?meetingId=${encodeURIComponent(meetingId)}`;
  const res = await apiService.getQuestionBanks(query);
  return res.data as QuestionBank[];
}

// Use apiService for admin endpoints
// Use apiService for admin endpoints
export async function fetchCompetitiveAnalysis(meetingId: string): Promise<CompetitiveAnalysis[]> {
  const query = `?meetingId=${encodeURIComponent(meetingId)}`;
  const res = await apiService.getCompetitiveAnalysis(query);
  return res.data as CompetitiveAnalysis[];
}

// Use apiService for admin endpoints
// Use apiService for admin endpoints
export async function fetchFollowUpActionItems(meetingId: string): Promise<FollowUpActionItem[]> {
  const query = `?meetingId=${encodeURIComponent(meetingId)}`;
  const res = await apiService.getFollowUpActionItems(query);
  return res.data as FollowUpActionItem[];
}