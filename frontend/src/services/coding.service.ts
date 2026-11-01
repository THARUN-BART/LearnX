import apiClient from './api';
import { CodeExecutionRequest, CodeExecutionResult } from '../types/learning.types';

class CodingService {
  async runCode(req: CodeExecutionRequest): Promise<CodeExecutionResult> {
    try {
      const response = await apiClient.post('/api/v1/code/run', req);
      return response.data.data || response.data;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        // Return clear status when backend execution container is offline
        return {
          status: 'Compilation Error',
          runtimeMs: 0,
          memoryMb: 0,
          stdout: '',
          stderr: 'Backend code execution service (Go runner) is currently offline at ' + apiClient.defaults.baseURL + '.\nPlease ensure the Go backend server is running to execute real code.',
          passedTests: 0,
          totalTests: 3
        };
      }
      throw err.response?.data?.message || err.message || 'Execution failed';
    }
  }

  async submitCode(req: CodeExecutionRequest): Promise<CodeExecutionResult> {
    try {
      const response = await apiClient.post('/api/v1/code/submit', req);
      return response.data.data || response.data;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        return {
          status: 'Compilation Error',
          runtimeMs: 0,
          memoryMb: 0,
          stdout: '',
          stderr: 'Go backend code execution service unreachable.\nConnect to Go REST API to submit solutions for test-suite verification.',
          passedTests: 0,
          totalTests: 3
        };
      }
      throw err.response?.data?.message || err.message || 'Submission failed';
    }
  }
}

export const codingService = new CodingService();
export default codingService;

