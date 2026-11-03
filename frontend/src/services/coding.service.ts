import apiClient, { unwrap } from './api';
import { CodeExecutionRequest, CodeExecutionResult, SubmissionHistoryItem } from '../types/learning.types';
import { toJudgeLanguage } from '../utils/language';
import { cacheInvalidate } from '../utils/apiCache';

const MAX_POLL_ATTEMPTS = 40;
const POLL_INTERVAL_MS = 400;

function mapResult(resultData: any, submissionId: string, diagnosis?: any): CodeExecutionResult {
  const passedTests = resultData.test_cases?.passed ?? resultData.passed_test_cases ?? 0;
  const totalTests = resultData.test_cases?.total ?? resultData.total_test_cases ?? 0;
  const executionTimeMs = resultData.execution_time_ms ?? 0;
  const memoryKb = resultData.memory_used_kb ?? 0;

  let stdout = '';
  let stderr = resultData.error_message || '';
  const results = Array.isArray(resultData.results) ? resultData.results : [];
  if (results.length > 0) {
    const firstVisible = results.find((r: any) => !r.is_hidden) || results[0];
    stdout = firstVisible.actual_output || '';
    if (!stderr && firstVisible.error_message) stderr = firstVisible.error_message;
  }

  return {
    status: resultData.status || 'Runtime Error',
    runtimeMs: executionTimeMs,
    memoryMb: Math.round((memoryKb / 1024) * 10) / 10,
    stdout,
    stderr,
    passedTests,
    totalTests,
    submissionId,
    diagnosis,
    results,
  };
}

class CodingService {
  async runCode(req: CodeExecutionRequest): Promise<CodeExecutionResult> {
    return this.executeSubmission(req);
  }

  async submitCode(req: CodeExecutionRequest): Promise<CodeExecutionResult> {
    return this.executeSubmission(req);
  }

  async retrySubmission(submissionId: string): Promise<CodeExecutionResult> {
    const submitRes = await apiClient.post(`/api/v1/submissions/${submissionId}/retry`);
    const data = unwrap<any>(submitRes.data, 'data');
    const newId = data?.submission_id || submissionId;
    return this.pollUntilDone(newId);
  }

  async getProblemSubmissions(problemId: string): Promise<SubmissionHistoryItem[]> {
    try {
      const res = await apiClient.get(`/api/v1/problems/${problemId}/submissions`, {
        params: { limit: 20 },
      });
      const list = unwrap<any>(res.data, 'data', 'submissions');
      const items = list?.submissions ?? list;
      return Array.isArray(items) ? items : [];
    } catch {
      return [];
    }
  }

  async getStudentSubmissions(params?: {
    page?: number;
    limit?: number;
    problemId?: string;
    status?: string;
  }): Promise<{ submissions: SubmissionHistoryItem[]; total: number; page: number; totalPages: number }> {
    try {
      const res = await apiClient.get('/api/v1/student/submissions', {
        params: {
          page: params?.page ?? 1,
          limit: params?.limit ?? 20,
          problem_id: params?.problemId,
          status: params?.status,
        },
      });
      const data = unwrap<any>(res.data, 'data');
      const list = data?.submissions || [];
      return {
        submissions: list.map((s: any) => ({
          id: s.submission_id || s.id,
          submission_id: s.submission_id || s.id,
          problem_id: s.problem_id,
          problem_title: s.problem_title || s.title,
          status: s.status,
          language: s.language,
          execution_time_ms: s.execution_time_ms,
          memory_used_kb: s.memory_used_kb,
          submitted_at: s.submitted_at,
          source_code: s.source_code,
        })),
        total: data?.total || list.length,
        page: data?.page || 1,
        totalPages: data?.total_pages || 1,
      };
    } catch {
      return { submissions: [], total: 0, page: 1, totalPages: 1 };
    }
  }

  async getLatestSubmission(problemId: string): Promise<SubmissionHistoryItem | null> {
    try {
      const res = await apiClient.get(`/api/v1/problems/${problemId}/submissions/latest`);
      return unwrap<any>(res.data, 'data') ?? null;
    } catch {
      return null;
    }
  }

  private async executeSubmission(req: CodeExecutionRequest): Promise<CodeExecutionResult> {
    try {
      const submitRes = await apiClient.post('/api/v1/submissions', {
        problem_id: req.problemId,
        language: toJudgeLanguage(req.language),
        source_code: req.code,
      });

      const submissionId =
        submitRes.data?.data?.submission_id ||
        submitRes.data?.submission_id ||
        submitRes.data?.id;

      if (!submissionId) {
        throw new Error('Failed to obtain a submission ID from the execution engine.');
      }

      const result = await this.pollUntilDone(submissionId);
      cacheInvalidate('problems');
      cacheInvalidate('problem|');
      cacheInvalidate('roadmap');
      cacheInvalidate('analytics');
      return result;
    } catch (err: any) {
      const message = err.response?.data?.error || err.response?.data?.message || err.message;

      if (err.response?.status === 401) {
        return {
          status: 'Compilation Error',
          runtimeMs: 0,
          memoryMb: 0,
          stdout: '',
          stderr: 'Authentication required: please log in to submit and evaluate code.',
          passedTests: 0,
          totalTests: 1,
        };
      }

      return {
        status: 'Runtime Error',
        runtimeMs: 0,
        memoryMb: 0,
        stdout: '',
        stderr: message || 'Execution failed',
        passedTests: 0,
        totalTests: 1,
      };
    }
  }

  private async pollUntilDone(submissionId: string): Promise<CodeExecutionResult> {
    for (let attempt = 0; attempt < MAX_POLL_ATTEMPTS; attempt++) {
      await new Promise((res) => setTimeout(res, POLL_INTERVAL_MS));

      const statusRes = await apiClient.get(`/api/v1/submissions/${submissionId}`);
      const statusData = unwrap<any>(statusRes.data, 'data');
      const status = statusData?.status;
      if (!status || status === 'Queued' || status === 'Processing') continue;

      const resultRes = await apiClient.get(`/api/v1/submissions/${submissionId}/result`);
      const resultData = unwrap<any>(resultRes.data, 'data') || statusData;

      let diagnosis = undefined;
      if (status !== 'Accepted') {
        try {
          const diagRes = await apiClient.get(`/api/v1/submissions/${submissionId}/diagnosis`);
          diagnosis = unwrap<any>(diagRes.data, 'data');
        } catch {
          // Diagnosis may still be running
        }
      }

      return mapResult(resultData, submissionId, diagnosis);
    }

    return {
      status: 'Time Limit Exceeded',
      runtimeMs: 16000,
      memoryMb: 0,
      stdout: '',
      stderr: 'Code assessment timed out waiting for the sandbox worker.',
      passedTests: 0,
      totalTests: 1,
      submissionId,
    };
  }
}

export const codingService = new CodingService();
export default codingService;
