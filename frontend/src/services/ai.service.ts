import apiClient from './api';
import {
  AssistantMessage,
  AssistantMode,
  ProgressiveHint,
  DryRunResult,
  ComplexityResult,
  ExplainMistakeResult
} from '../types/learning.types';

class AIService {
  /**
   * Request a progressive scaffolding hint (levels 1-5).
   */
  async getProgressiveHint(
    problemId: string,
    level: number,
    submissionId?: string
  ): Promise<ProgressiveHint> {
    const res = await apiClient.post('/api/v1/ai/hint', {
      problem_id: problemId,
      level,
      submission_id: submissionId
    });
    return res.data?.data || res.data;
  }

  /**
   * Retrieve all previously unlocked hints for a problem.
   */
  async getUnlockedHints(problemId: string): Promise<ProgressiveHint[]> {
    const res = await apiClient.get(`/api/v1/problems/${problemId}/hints`);
    return res.data?.data || res.data || [];
  }

  /**
   * Plain-English explanation of a diagnosed submission mistake.
   */
  async explainMistake(submissionId: string): Promise<ExplainMistakeResult> {
    const res = await apiClient.post('/api/v1/ai/explain-mistake', {
      submission_id: submissionId
    });
    return res.data?.data || res.data;
  }

  /**
   * Comprehensive section-by-section code walkthrough.
   */
  async getWalkthrough(
    problemId: string,
    sourceCode?: string,
    language?: string
  ): Promise<any> {
    const res = await apiClient.post('/api/v1/ai/walkthrough', {
      problem_id: problemId,
      source_code: sourceCode,
      language: language || 'cpp'
    });
    return res.data?.data || res.data;
  }

  /**
   * Step-by-step dry-run execution trace for visualization.
   */
  async getDryRun(
    problemId: string,
    sourceCode?: string,
    language?: string,
    customInput?: string
  ): Promise<DryRunResult> {
    const res = await apiClient.post('/api/v1/ai/dry-run', {
      problem_id: problemId,
      source_code: sourceCode,
      language: language || 'cpp',
      custom_input: customInput
    });
    return res.data?.data || res.data;
  }

  /**
   * Analyze time & space complexity against optimal bounds.
   */
  async getComplexity(
    problemId: string,
    sourceCode?: string,
    language?: string
  ): Promise<ComplexityResult> {
    const res = await apiClient.post('/api/v1/ai/complexity', {
      problem_id: problemId,
      source_code: sourceCode,
      language: language || 'cpp'
    });
    return res.data?.data || res.data;
  }

  /**
   * Conversational tutor query with context memory.
   */
  async askTutor(
    question: string,
    problemId?: string,
    conversationId?: string
  ): Promise<{ conversationId: string; answer: string }> {
    const res = await apiClient.post('/api/v1/ai/ask', {
      question,
      problem_id: problemId,
      conversation_id: conversationId
    });
    const d = res.data?.data || res.data;
    return {
      conversationId: d.conversation_id,
      answer: d.answer
    };
  }

  /**
   * Submit student feedback on diagnosis accuracy.
   */
  async submitMistakeFeedback(
    diagnosisId: string,
    feedbackType: 'helpful' | 'unhelpful' | 'inaccurate' | 'too_vague',
    studentFeedback?: string
  ): Promise<void> {
    await apiClient.post(`/api/v1/student/mistakes/${diagnosisId}/feedback`, {
      feedback_type: feedbackType,
      student_feedback: studentFeedback || ''
    });
  }

  /**
   * Conversational interface for the Assistant page.
   */
  async sendMessage(
    message: string,
    mode: AssistantMode,
    context?: { code?: string; problemTitle?: string; language?: string; problemId?: string; conversationId?: string }
  ): Promise<AssistantMessage> {
    try {
      if (mode === 'explain' && context?.code) {
        const wt = await this.getWalkthrough(context.problemId || 'general', context.code, context.language);
        const reply = wt.overall_summary || 'Here is the breakdown of your code.';
        return {
          id: 'ai-msg-' + Date.now(),
          role: 'assistant',
          content: reply,
          timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
          mode
        };
      }

      // Default Socratic tutor query
      const { answer } = await this.askTutor(message, context?.problemId, context?.conversationId);
      return {
        id: 'ai-msg-' + Date.now(),
        role: 'assistant',
        content: answer,
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        mode
      };
    } catch {
      // Offline fallback
      let reply = '';
      if (mode === 'hint') {
        reply = `**Algorithmic Hint for "${context?.problemTitle || 'current problem'}":**\n\n1. What is the bottleneck in the naive approach? Usually it's repeated scanning ($O(N^2)$), which can be turned into $O(1)$ lookups using an auxiliary hash table or a two-pointer technique.\n2. Invariant: Maintain previously visited elements in a hash set to detect complements immediately.`;
      } else if (mode === 'explain') {
        reply = `**Concept Explanation:**\n\nWhen optimizing this pattern, we trade auxiliary memory ($O(N)$ space) to achieve linear time ($O(N)$). Hash map lookups operate in expected $O(1)$ amortized time.`;
      } else {
        reply = `**Automated Debug Analysis:**\n\nEnsure array indices stay within $[0, N-1]$. Check loop boundaries and single-element inputs carefully.`;
      }

      return {
        id: 'ai-msg-' + Date.now(),
        role: 'assistant',
        content: reply,
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        mode
      };
    }
  }
}

export const aiService = new AIService();
export default aiService;
