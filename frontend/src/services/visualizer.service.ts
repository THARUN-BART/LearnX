import apiClient, { unwrap } from './api';
import {
  AlgorithmMetadata,
  ExecuteResponse,
  ExplainStepResponse,
  PracticeQuestion,
  VisualizerSession,
} from '../types/visualizer.types';

class VisualizerService {
  async getAlgorithms(): Promise<AlgorithmMetadata[]> {
    try {
      const res = await apiClient.get('/api/v1/visualizer/algorithms');
      return unwrap<AlgorithmMetadata[]>(res.data, 'algorithms') || [];
    } catch {
      return [];
    }
  }

  async getAlgorithmById(id: string): Promise<AlgorithmMetadata | null> {
    try {
      const res = await apiClient.get(`/api/v1/visualizer/algorithms/${id}`);
      return unwrap<AlgorithmMetadata>(res.data, 'algorithm') || null;
    } catch {
      return null;
    }
  }

  async executeAlgorithm(
    algorithmId: string,
    customInput?: Record<string, any>
  ): Promise<ExecuteResponse | null> {
    try {
      const res = await apiClient.post('/api/v1/visualizer/execute', {
        algorithm_id: algorithmId,
        custom_input: customInput || {},
      });
      return unwrap<ExecuteResponse>(res.data, 'data') || null;
    } catch {
      return null;
    }
  }

  async explainStep(params: {
    algorithmId: string;
    stepNumber: number;
    totalSteps: number;
    operation: string;
    explanation: string;
    currentState: Record<string, any>;
    userQuestion?: string;
  }): Promise<ExplainStepResponse | null> {
    try {
      const res = await apiClient.post('/api/v1/visualizer/explain-step', {
        algorithm_id: params.algorithmId,
        step_number: params.stepNumber,
        total_steps: params.totalSteps,
        operation: params.operation,
        explanation: params.explanation,
        current_state: params.currentState,
        user_question: params.userQuestion,
      });
      return unwrap<ExplainStepResponse>(res.data, 'data') || null;
    } catch {
      return null;
    }
  }

  async getPracticeQuestions(algorithmId: string): Promise<PracticeQuestion[]> {
    try {
      const res = await apiClient.get(`/api/v1/visualizer/practice-questions?algorithm_id=${algorithmId}`);
      return unwrap<PracticeQuestion[]>(res.data, 'questions') || [];
    } catch {
      return [];
    }
  }

  async createSession(params: {
    algorithmId: string;
    category: string;
    inputConfig: Record<string, any>;
    totalSteps: number;
  }): Promise<VisualizerSession | null> {
    try {
      const res = await apiClient.post('/api/v1/visualizer/sessions', {
        algorithm_id: params.algorithmId,
        category: params.category,
        input_config: params.inputConfig,
        total_steps: params.totalSteps,
      });
      return unwrap<VisualizerSession>(res.data, 'session') || null;
    } catch {
      return null;
    }
  }

  async updateSession(
    sessionId: string,
    params: { stepsCompleted: number; timeSpentSeconds: number; completed: boolean }
  ): Promise<boolean> {
    try {
      await apiClient.patch(`/api/v1/visualizer/sessions/${sessionId}`, {
        steps_completed: params.stepsCompleted,
        time_spent_seconds: params.timeSpentSeconds,
        completed: params.completed,
      });
      return true;
    } catch {
      return false;
    }
  }

  async listSessions(): Promise<VisualizerSession[]> {
    try {
      const res = await apiClient.get('/api/v1/visualizer/sessions');
      return unwrap<VisualizerSession[]>(res.data, 'sessions') || [];
    } catch {
      return [];
    }
  }
}

export const visualizerService = new VisualizerService();
export default visualizerService;

