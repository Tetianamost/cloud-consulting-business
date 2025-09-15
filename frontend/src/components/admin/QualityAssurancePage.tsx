// QualityAssurancePage.tsx
import React, { useEffect, useState } from "react";
import {
  fetchRecommendationAccuracy,
  fetchPeerReviews,
  fetchClientOutcomes,
  fetchContinuousImprovementItems,
  RecommendationAccuracy,
  PeerReview,
  ClientOutcome,
  ContinuousImprovementItem,
} from "../../services/qualityAssuranceService";

const QualityAssurancePage: React.FC = () => {
  const [recommendationAccuracy, setRecommendationAccuracy] = useState<RecommendationAccuracy[] | null>(null);
  const [peerReviews, setPeerReviews] = useState<PeerReview[] | null>(null);
  const [clientOutcomes, setClientOutcomes] = useState<ClientOutcome[] | null>(null);
  const [improvementItems, setImprovementItems] = useState<ContinuousImprovementItem[] | null>(null);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    setError(null);

    Promise.all([
      fetchRecommendationAccuracy().catch(e => { setError(e.message); return []; }),
      fetchPeerReviews().catch(e => { setError(e.message); return []; }),
      fetchClientOutcomes().catch(e => { setError(e.message); return []; }),
      fetchContinuousImprovementItems().catch(e => { setError(e.message); return []; }),
    ]).then(([recAcc, peerRev, clientOut, improv]) => {
      setRecommendationAccuracy(recAcc);
      setPeerReviews(peerRev);
      setClientOutcomes(clientOut);
      setImprovementItems(improv);
      setLoading(false);
    });
  }, []);

  if (loading) return <div>Loading Quality Assurance data...</div>;
  if (error) return <div>Error loading QA data: {error}</div>;

  return (
    <div>
      <h1>Quality Assurance</h1>
      <section>
        <h2>Recommendation Accuracy</h2>
        <div>
          {recommendationAccuracy && recommendationAccuracy.length > 0
            ? <div>{recommendationAccuracy.length} records loaded</div>
            : <div>No data</div>}
        </div>
      </section>
      <section>
        <h2>Peer Reviews</h2>
        <div>
          {peerReviews && peerReviews.length > 0
            ? <div>{peerReviews.length} records loaded</div>
            : <div>No data</div>}
        </div>
      </section>
      <section>
        <h2>Client Outcomes</h2>
        <div>
          {clientOutcomes && clientOutcomes.length > 0
            ? <div>{clientOutcomes.length} records loaded</div>
            : <div>No data</div>}
        </div>
      </section>
      <section>
        <h2>Continuous Improvement</h2>
        <div>
          {improvementItems && improvementItems.length > 0
            ? <div>{improvementItems.length} items loaded</div>
            : <div>No data</div>}
        </div>
      </section>
    </div>
  );
};

export default QualityAssurancePage;