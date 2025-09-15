// MeetingPrepPage.tsx
import React, { useEffect, useState } from 'react';
import {
  fetchPreMeetingBriefings,
  fetchQuestionBanks,
  fetchCompetitiveAnalysis,
  fetchFollowUpActionItems,
  PreMeetingBriefing,
  QuestionBank,
  CompetitiveAnalysis,
  FollowUpActionItem,
} from '../../services/meetingPrepService';

const DEMO_MEETING_ID = 'demo-meeting-123'; // Replace with actual meetingId as needed

const MeetingPrepPage: React.FC = () => {
  const [briefings, setBriefings] = useState<PreMeetingBriefing[]>([]);
  const [questionBanks, setQuestionBanks] = useState<QuestionBank[]>([]);
  const [competitiveAnalysis, setCompetitiveAnalysis] = useState<CompetitiveAnalysis[]>([]);
  const [followUpItems, setFollowUpItems] = useState<FollowUpActionItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    Promise.all([
      fetchPreMeetingBriefings(DEMO_MEETING_ID),
      fetchQuestionBanks(DEMO_MEETING_ID),
      fetchCompetitiveAnalysis(DEMO_MEETING_ID),
      fetchFollowUpActionItems(DEMO_MEETING_ID),
    ])
      .then(([briefings, questionBanks, competitiveAnalysis, followUpItems]) => {
        setBriefings(briefings);
        setQuestionBanks(questionBanks);
        setCompetitiveAnalysis(competitiveAnalysis);
        setFollowUpItems(followUpItems);
        setLoading(false);
      })
      .catch((err) => {
        setError('Failed to fetch meeting preparation data.');
        setLoading(false);
      });
  }, []);

  if (loading) return <div>Loading meeting preparation data...</div>;
  if (error) return <div>{error}</div>;

  // Placeholder: Replace with actual UI rendering as needed
  return (
    <div>
      <h1>Meeting Preparation</h1>
      <section>
        <h2>Pre-Meeting Briefings</h2>
        <pre>{JSON.stringify(briefings, null, 2)}</pre>
      </section>
      <section>
        <h2>Question Banks</h2>
        <pre>{JSON.stringify(questionBanks, null, 2)}</pre>
      </section>
      <section>
        <h2>Competitive Analysis</h2>
        <pre>{JSON.stringify(competitiveAnalysis, null, 2)}</pre>
      </section>
      <section>
        <h2>Follow-Up Action Items</h2>
        <pre>{JSON.stringify(followUpItems, null, 2)}</pre>
      </section>
    </div>
  );
};

export default MeetingPrepPage;