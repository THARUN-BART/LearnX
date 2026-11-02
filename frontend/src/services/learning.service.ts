import apiClient, { unwrap } from './api';
import { Problem, TestCase } from '../types/learning.types';
import { cacheKey, cached, TTL } from '../utils/apiCache';

function generateDefaultStarterCode(): Record<string, string> {
  return {
    python: `import sys

def solution():
    data = sys.stdin.read().split()
    print(" ".join(data))

if __name__ == '__main__':
    solution()
`,
    cpp: `#include <iostream>
using namespace std;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    int n;
    while (cin >> n) {
        cout << n << "\\n";
    }
    return 0;
}
`,
    c: `#include <stdio.h>
int main() {
    int n;
    while (scanf("%d", &n) != EOF) printf("%d\\n", n);
    return 0;
}
`,
    java: `import java.util.Scanner;
public class Main {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        while (scanner.hasNextInt()) {
            System.out.println(scanner.nextInt());
        }
        scanner.close();
    }
}
`,
    go: `package main
import (
    "bufio"
    "fmt"
    "os"
)
func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }
}
`,
  };
}

function mapSampleCases(sampleCases: any[]): TestCase[] {
  return sampleCases.map((tc: any, idx: number) => ({
    id: tc.id ?? `tc-${idx + 1}`,
    input: tc.input_data ?? tc.input ?? '',
    expectedOutput: tc.expected_output ?? tc.expectedOutput ?? '',
    actualOutput: '',
    passed: false,
  }));
}

export function mapRawProblemToProblem(raw: any, sampleCases?: any[]): Problem {
  const tests = Array.isArray(sampleCases) && sampleCases.length
    ? mapSampleCases(sampleCases)
    : raw.testCases || [];

  const examples =
    raw.examples ||
    tests.slice(0, 3).map((tc: TestCase) => ({
      input: tc.input,
      output: tc.expectedOutput,
    }));

  return {
    id: raw.id,
    title: raw.title || 'Coding Problem',
    slug: raw.slug || '',
    difficulty: raw.difficulty || 'Unspecified',
    accuracyRate: raw.accuracyRate || (raw.difficulty === 'Easy' ? 88 : raw.difficulty === 'Medium' ? 72 : 54),
    problem_order: raw.problem_order || 1,
    external_platform: raw.external_platform || 'takeuforward',
    external_problem_url: raw.external_problem_url || '',
    source_sheet: raw.source_sheet || 'strivers-a2z-dsa-sheet',
    source_problem_id: raw.source_problem_id || '',
    tags: Array.isArray(raw.tags) ? raw.tags : ['DSA'],
    progress_status: raw.progress_status || 'Not Started',
    subtopic_id: raw.subtopic_id,
    topic_id: raw.topic_id,
    topicId: raw.topic_id,
    summary: raw.summary || `Practice problem from Striver's A2Z DSA Sheet: ${raw.title}`,
    description:
      raw.description ||
      raw.summary ||
      `Solve "${raw.title}". Open the official TakeUForward statement for the full problem, then implement a solution that reads stdin and writes stdout.`,
    constraints: raw.constraints || [
      'Follow the official problem constraints from the source sheet.',
      'Time Limit: 2.0s',
      'Memory Limit: 256MB',
    ],
    examples: examples.length
      ? examples
      : [{ input: 'See official statement', output: 'See official statement' }],
    hints: raw.hints || [
      'Start from the brute-force idea, then look for a linear or n-log-n pattern.',
      'Check empty, single-element, and already-sorted inputs.',
    ],
    starterCode: raw.starterCode || generateDefaultStarterCode(),
    testCases: tests.length ? tests : [{ id: 'tc-sample-1', input: '', expectedOutput: '' }],
    related: Array.isArray(raw.related) ? raw.related : undefined,
  };
}

export interface ProblemListParams {
  topicId?: string;
  subtopicId?: string;
  difficulty?: string;
  status?: string;
  q?: string;
  page?: number;
  limit?: number;
}

class LearningService {
  async getProblems(topicIdOrParams?: string | ProblemListParams): Promise<Problem[]> {
    const params: ProblemListParams =
      typeof topicIdOrParams === 'string' ? { topicId: topicIdOrParams } : topicIdOrParams || {};
    const limit = params.limit ?? 40;
    const page = params.page ?? 1;
    const key = cacheKey([
      'problems',
      params.topicId,
      params.subtopicId,
      params.difficulty,
      params.status,
      params.q,
      page,
      limit,
    ]);

    return cached(key, TTL.short, async () => {
      try {
        const response = await apiClient.get('/api/v1/problems', {
          params: {
            topic_id: params.topicId,
            subtopic_id: params.subtopicId,
            difficulty: params.difficulty,
            status: params.status,
            q: params.q,
            page,
            limit,
            sort: 'problem_order',
          },
        });
        const rawProblems = unwrap<any[]>(response.data, 'problems');
        if (!Array.isArray(rawProblems)) return [];
        return rawProblems.map((p) => mapRawProblemToProblem(p));
      } catch {
        return [];
      }
    });
  }

  async searchProblems(q: string, limit = 20): Promise<Problem[]> {
    if (!q.trim()) return [];
    try {
      const response = await apiClient.get('/api/v1/problems/search', {
        params: { q: q.trim(), limit },
      });
      const rawProblems = unwrap<any[]>(response.data, 'problems');
      if (!Array.isArray(rawProblems)) return [];
      return rawProblems.map((p) => mapRawProblemToProblem(p));
    } catch {
      return [];
    }
  }

  async getProblemsBySubtopic(subtopicId: string): Promise<Problem[]> {
    return cached(cacheKey(['subtopic-problems', subtopicId]), TTL.medium, async () => {
      try {
        const response = await apiClient.get(`/api/v1/subtopics/${subtopicId}/problems`, {
          params: { limit: 80 },
        });
        const rawProblems = unwrap<any[]>(response.data, 'problems');
        if (!Array.isArray(rawProblems)) return [];
        return rawProblems.map((p) => mapRawProblemToProblem(p));
      } catch {
        return this.getProblems({ subtopicId, limit: 80 });
      }
    });
  }

  async getProblemById(problemId: string): Promise<Problem | null> {
    return cached(cacheKey(['problem', problemId]), TTL.short, async () => {
      try {
        const [response, tcRes] = await Promise.allSettled([
          apiClient.get(`/api/v1/problems/${problemId}`),
          apiClient.get(`/api/v1/problems/${problemId}/testcases/sample`),
        ]);

        if (response.status !== 'fulfilled') return null;
        const raw = unwrap<any>(response.value.data, 'problem');
        if (!raw || !raw.id) return null;

        const related = response.value.data?.related;
        const samples =
          tcRes.status === 'fulfilled' ? unwrap<any[]>(tcRes.value.data, 'data', 'test_cases') : [];

        const problem = mapRawProblemToProblem(
          { ...raw, related },
          Array.isArray(samples) ? samples : []
        );
        return problem;
      } catch {
        return null;
      }
    });
  }
}

export const learningService = new LearningService();
export default learningService;
