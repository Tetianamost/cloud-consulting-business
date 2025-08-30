/**
 * Report service for fetching reports from backend
 */

// Temporary local type until project type is available
export interface ReportWithInquiry {
  id: string;
  title: string;
  created_at: string;
  status: string;
  type: string;
  inquiry?: {
    name?: string;
    company?: string;
    // Add more fields as needed
  };
  inquiry_id?: string;
  generated_by?: string;
  content: string;
  // Add more fields as needed based on backend response
}

import apiService from "./api";

// Fetch all reports
export async function fetchReports(): Promise<ReportWithInquiry[]> {
  // apiService.listReports returns AdminReportsResponse, extract .data
  const res = await apiService.listReports();
  return res.data;
}
/**
 * Download a report as PDF or HTML
 * @param inquiryId The inquiry ID of the report
 * @param format 'pdf' or 'html'
 * @returns Blob of the downloaded file
 */
export async function downloadReport(
  inquiryId: string | undefined,
  format: 'pdf' | 'html'
): Promise<Blob> {
  if (!inquiryId) {
    throw new Error('Missing inquiry ID for report download');
  }
  return await apiService.downloadReport(inquiryId, format);
}