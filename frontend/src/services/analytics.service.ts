import apiClient from './api';
import { AnalyticsSummary, MistakeRecord } from '../types/learning.types';

export const MOCK_MISTAKES: MistakeRecord[] = [
  {
    id: 'mst-1',
    problemId: 'prob-two-sum',
    problemTitle: 'Two Sum',
    topic: 'Arrays & Hashing',
    category: 'edge_cases',
    categoryLabel: 'Edge Case Neglect',
    frequency: 3,
    date: '2026-10-28',
    failedCodeSnippet: `seen = {}
for i, num in enumerate(nums):
    seen[num] = i  # Overwrote before checking complement!
    if target - num in seen:
        return [seen[target - num], i]`,
    explanation: 'Overwriting the dictionary key with the current index before testing for the complement caused identical elements to pair with themselves.',
    whyFailed: 'Failed on testcase nums=[3,3], target=6. The code returned [0, 0] instead of [0, 1].',
    suggestedCorrection: 'Check `if target - num in seen` before inserting `num` into the hash table.',
    recommendedTopicId: 'hashing',
    recommendedTopicTitle: 'Hashing & Frequency Maps'
  },
  {
    id: 'mst-2',
    problemId: 'prob-max-subarray',
    problemTitle: 'Maximum Subarray',
    topic: 'Arrays',
    category: 'time_complexity',
    categoryLabel: 'Time Limit Exceeded (O(N^3))',
    frequency: 2,
    date: '2026-10-25',
    failedCodeSnippet: `max_sum = -float('inf')
for i in range(len(nums)):
    for j in range(i, len(nums)):
        max_sum = max(max_sum, sum(nums[i:j+1]))`,
    explanation: 'Nested loops calculating slice sums resulted in O(N^3) time complexity, timing out on 10^5 elements.',
    whyFailed: 'Exceeded 2000ms limit on large arrays.',
    suggestedCorrection: 'Use Kadanes single-pass algorithm to achieve optimal O(N) time and O(1) space.',
    recommendedTopicId: 'arrays',
    recommendedTopicTitle: 'Prefix Sums & Dynamic Kadane'
  },
  {
    id: 'mst-3',
    problemId: 'prob-valid-parentheses',
    problemTitle: 'Valid Parentheses',
    topic: 'Stacks and Queues',
    category: 'logical',
    categoryLabel: 'Logical Stack Underflow',
    frequency: 4,
    date: '2026-10-22',
    failedCodeSnippet: `stack = []
for c in s:
    if c == ')':
        top = stack.pop() # Error when stack is empty!
        if top != '(': return False`,
    explanation: 'Popping from an empty stack when input starts with a closing bracket raises an IndexError / panic.',
    whyFailed: 'Crashed on string s = "]" or ")(',
    suggestedCorrection: 'Always guard stack operations: `if not stack: return False` before invoking `.pop()`.',
    recommendedTopicId: 'stacks-queues',
    recommendedTopicTitle: 'Monotonic Stacks & Invariants'
  },
  {
    id: 'mst-4',
    problemId: 'prob-binary-search',
    problemTitle: 'Search in Rotated Sorted Array',
    topic: 'Searching and Sorting',
    category: 'algorithm_selection',
    categoryLabel: 'Incorrect Pivot Logic',
    frequency: 3,
    date: '2026-10-18',
    failedCodeSnippet: `mid = (low + high) // 2
if nums[mid] == target: return mid
if nums[low] <= nums[mid]:
    # Incomplete boundary condition
    if nums[low] <= target < nums[mid]: high = mid - 1
    else: low = mid + 1`,
    explanation: 'Boundary inequality missed strict monotonic half checking when duplicates or non-strict bounds occurred.',
    whyFailed: 'Infinite loop when low == high - 1.',
    suggestedCorrection: 'Verify whether the left half or right half is strictly sorted before testing inclusion.',
    recommendedTopicId: 'searching-sorting',
    recommendedTopicTitle: 'Binary Search On Answer Spaces'
  }
];

export const MOCK_ANALYTICS: AnalyticsSummary = {
  totalSolved: 64,
  easySolved: 32,
  mediumSolved: 26,
  hardSolved: 6,
  streakDays: 14,
  hoursSpent: 48,
  accuracyRate: 79.4,
  weeklyActivity: [
    { day: 'Mon', solved: 4, hours: 2.5 },
    { day: 'Tue', solved: 7, hours: 3.8 },
    { day: 'Wed', solved: 5, hours: 2.1 },
    { day: 'Thu', solved: 9, hours: 4.5 },
    { day: 'Fri', solved: 6, hours: 3.0 },
    { day: 'Sat', solved: 11, hours: 5.2 },
    { day: 'Sun', solved: 8, hours: 4.0 }
  ],
  topicProficiency: [
    { topic: 'Arrays', score: 85, total: 100 },
    { topic: 'Strings', score: 78, total: 100 },
    { topic: 'Sorting', score: 72, total: 100 },
    { topic: 'Linked Lists', score: 68, total: 100 },
    { topic: 'Stacks & Queues', score: 82, total: 100 },
    { topic: 'Hashing', score: 90, total: 100 },
    { topic: 'Trees', score: 55, total: 100 },
    { topic: 'DP', score: 35, total: 100 }
  ],
  difficultyDistribution: [
    { name: 'Easy', count: 32, color: '#22C55E' },
    { name: 'Medium', count: 26, color: '#F59E0B' },
    { name: 'Hard', count: 6, color: '#EF4444' }
  ],
  monthlyProgress: [
    { month: 'Jun', solved: 12 },
    { month: 'Jul', solved: 22 },
    { month: 'Aug', solved: 35 },
    { month: 'Sep', solved: 48 },
    { month: 'Oct', solved: 64 }
  ]
};

class AnalyticsService {
  async getAnalytics(): Promise<AnalyticsSummary> {
    try {
      const response = await apiClient.get('/api/v1/analytics/dashboard');
      return response.data.data || response.data;
    } catch {
      return MOCK_ANALYTICS;
    }
  }

  async getMistakes(): Promise<MistakeRecord[]> {
    try {
      const response = await apiClient.get('/api/v1/analytics/mistakes');
      return response.data.data || response.data;
    } catch {
      return MOCK_MISTAKES;
    }
  }
}

export const analyticsService = new AnalyticsService();
export default analyticsService;

