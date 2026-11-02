import { aiClient, unwrap } from './api';
import {
  AssistantMessage,
  AssistantMode,
  ProgressiveHint,
  DryRunResult,
  ComplexityResult,
  ExplainMistakeResult,
  PedagogyAnalysis,
  AIConversation,
} from '../types/learning.types';
import { toJudgeLanguage } from '../utils/language';

class AIService {
  async getProgressiveHint(
    problemId: string,
    level: number,
    submissionId?: string
  ): Promise<ProgressiveHint> {
    const res = await aiClient.post('/api/v1/ai/hint', {
      problem_id: problemId,
      level,
      submission_id: submissionId,
    });
    return unwrap<ProgressiveHint>(res.data, 'data');
  }

  async getUnlockedHints(problemId: string): Promise<ProgressiveHint[]> {
    const res = await aiClient.get(`/api/v1/problems/${problemId}/hints`);
    const list = unwrap<any>(res.data, 'data');
    return Array.isArray(list) ? list : [];
  }

  async explainMistake(submissionId: string): Promise<ExplainMistakeResult> {
    const res = await aiClient.post('/api/v1/ai/explain-mistake', {
      submission_id: submissionId,
    });
    return unwrap<ExplainMistakeResult>(res.data, 'data');
  }

  async getWalkthrough(problemId: string, sourceCode?: string, language?: string): Promise<any> {
    const res = await aiClient.post('/api/v1/ai/walkthrough', {
      problem_id: problemId,
      source_code: sourceCode,
      language: toJudgeLanguage(language),
    });
    return unwrap<any>(res.data, 'data');
  }

  async getDryRun(
    problemId: string,
    sourceCode?: string,
    language?: string,
    customInput?: string
  ): Promise<DryRunResult> {
    const res = await aiClient.post('/api/v1/ai/dry-run', {
      problem_id: problemId,
      source_code: sourceCode,
      language: toJudgeLanguage(language),
      custom_input: customInput,
      test_case_input: customInput,
    });
    return unwrap<DryRunResult>(res.data, 'data');
  }

  async getComplexity(
    problemId: string,
    sourceCode?: string,
    language?: string
  ): Promise<ComplexityResult> {
    const res = await aiClient.post('/api/v1/ai/complexity', {
      problem_id: problemId,
      source_code: sourceCode,
      language: toJudgeLanguage(language),
    });
    return unwrap<ComplexityResult>(res.data, 'data');
  }

  async askTutor(
    question: string,
    problemId?: string,
    conversationId?: string
  ): Promise<{ conversationId: string; answer: string }> {
    const res = await aiClient.post('/api/v1/ai/ask', {
      question,
      problem_id: problemId,
      conversation_id: conversationId,
    });
    const d = unwrap<any>(res.data, 'data');
    return {
      conversationId: d.conversation_id,
      answer: d.answer,
    };
  }

  async listConversations(): Promise<AIConversation[]> {
    try {
      const res = await aiClient.get('/api/v1/ai/conversations', { params: { limit: 20 } });
      const list = unwrap<any[]>(res.data, 'data');
      return Array.isArray(list) ? list : [];
    } catch {
      return [];
    }
  }

  async getConversationMessages(conversationId: string): Promise<AssistantMessage[]> {
    try {
      const res = await aiClient.get(`/api/v1/ai/conversations/${conversationId}`);
      const payload = unwrap<any>(res.data, 'data');
      const list = Array.isArray(payload) ? payload : payload?.messages;
      const messages = Array.isArray(list) ? list : [];
      return messages.map((m: any) => ({
        id: m.id || String(m.created_at),
        role: m.role === 'user' ? 'user' : 'assistant',
        content: m.message ?? m.content ?? '',
        timestamp: m.created_at
          ? new Date(m.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
          : '',
      }));
    } catch {
      return [];
    }
  }

  async submitMistakeFeedback(
    diagnosisId: string,
    feedbackType: 'helpful' | 'unhelpful' | 'inaccurate' | 'too_vague',
    studentFeedback?: string
  ): Promise<void> {
    await aiClient.post(`/api/v1/student/mistakes/${diagnosisId}/feedback`, {
      feedback_type: feedbackType,
      student_feedback: studentFeedback || '',
    });
  }

  async analyzePedagogy(
    studentCode: string,
    testCode: string,
    hintLevel: string = 'L1'
  ): Promise<PedagogyAnalysis> {
    const res = await aiClient.post('/api/v1/ai/pedagogy/analyze', {
      student_code: studentCode,
      test_code: testCode,
      hint_level: hintLevel,
      timeout: 3.0,
    });
    return unwrap<PedagogyAnalysis>(res.data, 'data');
  }

  async sendMessage(
    message: string,
    mode: AssistantMode,
    context?: {
      code?: string;
      problemTitle?: string;
      language?: string;
      problemId?: string;
      conversationId?: string;
    }
  ): Promise<AssistantMessage & { conversationId?: string }> {
    if (mode === 'explain' && context?.code) {
      const wt = await this.getWalkthrough(
        context.problemId || 'general',
        context.code,
        context.language
      );
      const reply = wt.overall_summary || 'Here is the breakdown of your code.';
      const extras = [
        ...(wt.key_strengths || []).map((s: string) => `• Strength: ${s}`),
        ...(wt.areas_to_improve || []).map((s: string) => `• Improve: ${s}`),
      ].join('\n');
      return {
        id: 'ai-msg-' + Date.now(),
        role: 'assistant',
        content: extras ? `${reply}\n\n${extras}` : reply,
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        mode,
      };
    }

    const prompt =
      mode === 'hint'
        ? `Give a Socratic hint (no full solution) for: ${message}`
        : mode === 'debug'
        ? `Debug this DSA issue without dumping a full rewrite: ${message}`
        : message;

    const { answer, conversationId } = await this.askTutor(
      prompt,
      context?.problemId,
      context?.conversationId
    );
    return {
      id: 'ai-msg-' + Date.now(),
      role: 'assistant',
      content: answer,
      timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      mode,
      conversationId,
    };
  }
}

export const aiService = new AIService();
export default aiService;
