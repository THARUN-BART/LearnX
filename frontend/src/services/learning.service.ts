import apiClient from './api';
import { Problem, TestCase } from '../types/learning.types';

function generateDefaultStarterCode(title: string): Record<string, string> {
  const safeName = title.replace(/[^a-zA-Z0-9]/g, '');
  return {
    python: `import sys

def solution():
    # Read standard input
    input_data = sys.stdin.read().split()
    if not input_data:
        return
    # TODO: Implement solution logic
    print(" ".join(input_data))

if __name__ == '__main__':
    solution()
`,
    cpp: `#include <iostream>
#include <vector>
#include <string>
#include <algorithm>

using namespace std;

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    
    // Read input and solve
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
    while (scanf("%d", &n) != EOF) {
        printf("%d\\n", n);
    }
    return 0;
}
`,
    java: `import java.util.Scanner;

public class Main {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        while (scanner.hasNextInt()) {
            int n = scanner.nextInt();
            System.out.println(n);
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
`
  };
}

function mapRawProblemToProblem(raw: any): Problem {
  return {
    id: raw.id,
    title: raw.title || 'Coding Problem',
    slug: raw.slug || '',
    difficulty: raw.difficulty || 'Easy',
    accuracyRate: raw.accuracyRate || (raw.difficulty === 'Easy' ? 88 : raw.difficulty === 'Medium' ? 72 : 54),
    problem_order: raw.problem_order || 1,
    external_platform: raw.external_platform || 'takeuforward',
    external_problem_url: raw.external_problem_url || '',
    source_sheet: raw.source_sheet || "strivers-a2z-dsa-sheet",
    source_problem_id: raw.source_problem_id || '',
    tags: Array.isArray(raw.tags) ? raw.tags : ['DSA', 'Striver A2Z'],
    progress_status: raw.progress_status || 'Not Started',
    subtopic_id: raw.subtopic_id,
    topic_id: raw.topic_id,
    summary: raw.summary || `Practice problem from Striver's A2Z DSA Sheet: ${raw.title}`,
    description: raw.description || raw.summary || `Solve the problem "${raw.title}". Read the input from standard input and print the result to standard output.`,
    constraints: raw.constraints || [
      '1 <= N <= 10^5',
      '-10^9 <= Value <= 10^9',
      'Time Limit: 2.0s',
      'Memory Limit: 256MB'
    ],
    examples: raw.examples || [
      {
        input: 'Sample Input',
        output: 'Expected Output',
        explanation: 'Follow standard algorithmic approach.'
      }
    ],
    hints: raw.hints || [
      'Consider optimal time complexity before coding.',
      'Check edge cases with empty input or extreme bounds.'
    ],
    starterCode: raw.starterCode || generateDefaultStarterCode(raw.title || 'Solution'),
    testCases: raw.testCases || [
      {
        id: 'tc-sample-1',
        input: '1 2 3',
        expectedOutput: '1 2 3'
      }
    ]
  };
}

class LearningService {
  /**
   * Fetches the list of problems, optionally filtered by topic ID or slug.
   * Returns a clean, non-null Problem[] array.
   */
  async getProblems(topicId?: string): Promise<Problem[]> {
    try {
      const response = await apiClient.get('/api/v1/problems', {
        params: topicId ? { topic_id: topicId, limit: 150 } : { limit: 100 }
      });

      // Extract array safely
      const rawProblems =
        response.data.problems ||
        response.data.data ||
        (Array.isArray(response.data) ? response.data : []);

      if (!Array.isArray(rawProblems) || rawProblems.length === 0) {
        return [];
      }

      return rawProblems.map(mapRawProblemToProblem);
    } catch {
      return [];
    }
  }

  /**
   * Fetches problem details by ID or slug, along with its sample test cases.
   */
  async getProblemById(problemId: string): Promise<Problem | null> {
    try {
      const response = await apiClient.get(`/api/v1/problems/${problemId}`);
      const raw = response.data.problem || response.data.data || response.data;
      if (!raw || !raw.id) {
        return null;
      }

      const problem = mapRawProblemToProblem(raw);

      // Try fetching sample test cases registered in the backend
      try {
        const tcRes = await apiClient.get(`/api/v1/problems/${raw.id}/testcases/sample`);
        const sampleCases = tcRes.data.data || tcRes.data.test_cases || [];
        if (Array.isArray(sampleCases) && sampleCases.length > 0) {
          problem.testCases = sampleCases.map((tc: any) => ({
            id: tc.id,
            input: tc.input_data,
            expectedOutput: tc.expected_output,
            actualOutput: '',
            passed: false
          }));
        }
      } catch {
        // Sample test cases not found or unauthenticated
      }

      return problem;
    } catch {
      return null;
    }
  }
}

export const learningService = new LearningService();
export default learningService;
