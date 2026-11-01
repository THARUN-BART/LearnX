import { Topic, Problem } from '../types/learning.types';

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
export const TOKEN_KEY = 'learnx_access_token';
export const REFRESH_TOKEN_KEY = 'learnx_refresh_token';
export const USER_KEY = 'learnx_auth_user';
export const PROFILE_KEY = 'learnx_student_profile';

export const DEFAULT_ROADMAP_ID = 'c3e1a6b0-2c4f-4d9a-9b11-4f0c8e1a2b33';

// Official TakeUForward Striver's A2Z DSA Sheet 19 Sections
export const DSA_TOPICS: Topic[] = [
  {
    id: 'top-1',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 1,
    title: 'Beginner Problems',
    slug: 'beginner-problems',
    description: 'Patterns, language basics, I/O handling, conditionals, loops, and basic time complexity intuition.',
    difficulty: 'Easy',
    estimatedHours: 10,
    totalProblems: 24,
    problem_count: 24,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Basic math', 'Logic'],
    subtopics: ['Patterns', 'Basic Math', 'Basic Recursion', 'Basic Hashing'],
    explanation: 'Beginner Problems form the baseline of algorithmic problem-solving: nested loops for rectangular and pyramid patterns, division/modulo for integer extraction, and basic recursion.'
  },
  {
    id: 'top-2',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 2,
    title: 'Sorting',
    slug: 'sorting',
    description: 'Selection sort, Bubble sort, Insertion sort, Merge sort, Quick sort, and Recursive sorting invariants.',
    difficulty: 'Easy',
    estimatedHours: 8,
    totalProblems: 7,
    problem_count: 7,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Beginner Problems'],
    subtopics: ['Selection Sort', 'Bubble Sort', 'Insertion Sort', 'Merge Sort', 'Quick Sort'],
    explanation: 'Sorting algorithms demonstrate fundamental algorithmic design paradigms: quadratic brute-force passes versus divide-and-conquer routines with logarithmic call stacks.'
  },
  {
    id: 'top-3',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 3,
    title: 'Arrays',
    slug: 'arrays',
    description: 'Contiguous memory structures, two pointers, sliding window, prefix sums, Kadane algorithm, and Dutch National Flag.',
    difficulty: 'Medium',
    estimatedHours: 24,
    totalProblems: 40,
    problem_count: 40,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Sorting'],
    subtopics: ['Easy Arrays', 'Medium Arrays (Kadane, 2Sum)', 'Hard Arrays (3Sum, 4Sum, Merge Intervals)'],
    explanation: 'Arrays teach optimal linear-time traversals, in-place partitions, and prefix hash lookups to replace nested iterations.'
  },
  {
    id: 'top-4',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 4,
    title: 'Binary Search',
    slug: 'binary-search',
    description: 'Search in sorted arrays, rotated sorted arrays, search on answer spaces, and capacity allocation models.',
    difficulty: 'Medium',
    estimatedHours: 18,
    totalProblems: 32,
    problem_count: 32,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Arrays'],
    subtopics: ['BS on 1D Arrays', 'BS on Answers (Koko, Book Allocation)', 'BS on 2D Matrices'],
    explanation: 'Binary Search is a monotonic reduction technique. When the condition f(x) is monotonic, the optimal answer can be searched in logarithmic iterations.'
  },
  {
    id: 'top-5',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 5,
    title: 'Strings (Advanced Algo)',
    slug: 'strings',
    description: 'Anagrams, palindrome verification, rolling hash, Roman numerals, Rabin-Karp, and KMP pattern matching.',
    difficulty: 'Medium',
    estimatedHours: 14,
    totalProblems: 15,
    problem_count: 15,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Arrays'],
    subtopics: ['Basic String Ops', 'Medium String Problems', 'Advanced Pattern Search'],
    explanation: 'Strings test character encoding, immutable memory considerations, frequency counting, and string hashing.'
  },
  {
    id: 'top-6',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 6,
    title: 'Linked-List',
    slug: 'linked-list',
    description: 'Singly, doubly, and circular linked lists, Fast & Slow pointer (Floyd cycle detection), and reversal in groups of K.',
    difficulty: 'Medium',
    estimatedHours: 16,
    totalProblems: 31,
    problem_count: 31,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Arrays'],
    subtopics: ['1D Linked List', 'Doubly Linked List', 'Medium LL Problems', 'Hard LL Problems'],
    explanation: 'Linked list manipulation trains rigorous pointer pointer arithmetic and edge case guards for null references and cycles.'
  },
  {
    id: 'top-7',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 7,
    title: 'Recursion',
    slug: 'recursion',
    description: 'Call stacks, base cases, tree recursions, subset generation, and divide-and-conquer paradigms.',
    difficulty: 'Medium',
    estimatedHours: 14,
    totalProblems: 19,
    problem_count: 19,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Beginner Problems'],
    subtopics: ['Subsequence Patterns', 'Combination Sum', 'Subset Generation'],
    explanation: 'Recursion systematic state exploration creates the foundation for backtracking, trees, graphs, and dynamic programming.'
  },
  {
    id: 'top-8',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 8,
    title: 'Bit Manipulation',
    slug: 'bit-manipulation',
    description: 'Bitwise AND, OR, XOR, shifts, count set bits, single number variations, and power set generator.',
    difficulty: 'Easy',
    estimatedHours: 8,
    totalProblems: 18,
    problem_count: 18,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Beginner Problems'],
    subtopics: ['Learn Bit Magic', 'Interview Problems', 'Advanced Maths & Bits'],
    explanation: 'Hardware bit-level arithmetic operates in single CPU cycles, powering XOR cancellation tricks and set masks.'
  },
  {
    id: 'top-9',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 9,
    title: 'Stack / Queues',
    slug: 'stack-queues',
    description: 'LIFO & FIFO primitives, Monotonic Stack, Next Greater Element, and circular deque design.',
    difficulty: 'Medium',
    estimatedHours: 16,
    totalProblems: 30,
    problem_count: 30,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Arrays', 'Linked-List'],
    subtopics: ['Learning Stack & Queue', 'Prefix/Infix/Postfix', 'Monotonic Stack / Queue Problems'],
    explanation: 'Monotonic stacks maintain running maximums or minimums in amortized O(1) time per element.'
  },
  {
    id: 'top-10',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 10,
    title: 'Sliding Window / 2 Pointer',
    slug: 'sliding-window-2-pointer',
    description: 'Dynamic and fixed size windows, longest substring without repeating characters, and minimum window substring.',
    difficulty: 'Medium',
    estimatedHours: 14,
    totalProblems: 12,
    problem_count: 12,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Arrays', 'Stack / Queues'],
    subtopics: ['Constant Window', 'Longest Substring Patterns', 'Counting Subarrays'],
    explanation: 'Sliding windows maintain a contiguous state range without re-evaluating elements from scratch, cutting quadratic costs to linear.'
  },
  {
    id: 'top-11',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 11,
    title: 'Heaps',
    slug: 'heaps',
    description: 'Min-heaps, max-heaps, heapify, top-K frequent elements, and median from data stream.',
    difficulty: 'Medium',
    estimatedHours: 12,
    totalProblems: 17,
    problem_count: 17,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Arrays'],
    subtopics: ['Learning Heaps', 'Medium Problems (Top K)', 'Hard Problems (Median Stream)'],
    explanation: 'Priority queues guarantee O(1) peek with O(log N) updates, essential for Dijkstra and greedy scheduling.'
  },
  {
    id: 'top-12',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 12,
    title: 'Greedy Algorithms',
    slug: 'greedy-algorithms',
    description: 'Locally optimal choices, interval scheduling, fractional knapsack, jump game, and N meetings in one room.',
    difficulty: 'Medium',
    estimatedHours: 12,
    totalProblems: 15,
    problem_count: 15,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Sorting'],
    subtopics: ['Easy Greedy', 'Medium/Hard Greedy'],
    explanation: 'Greedy strategies construct solutions step by step, choosing the best local option under mathematically proven conditions.'
  },
  {
    id: 'top-13',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 13,
    title: 'Binary Trees',
    slug: 'binary-trees',
    description: 'Tree traversals (DFS & BFS), diameter, height, lowest common ancestor, view traversals, and serialization.',
    difficulty: 'Medium',
    estimatedHours: 20,
    totalProblems: 36,
    problem_count: 36,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Recursion'],
    subtopics: ['Traversals', 'Medium Problems', 'Hard Problems'],
    explanation: 'Hierarchical tree branching mirrors real-world file trees and ASTs, reinforcing depth-first and breadth-first recursion.'
  },
  {
    id: 'top-14',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 14,
    title: 'Binary Search Trees',
    slug: 'binary-search-trees',
    description: 'BST validation, inorder property, successor, predecessor, floor, ceil, and two-sum in BST.',
    difficulty: 'Medium',
    estimatedHours: 12,
    totalProblems: 16,
    problem_count: 16,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Binary Trees'],
    subtopics: ['Concepts', 'Practice Problems'],
    explanation: 'BST invariants maintain dynamically sorted keys, achieving logarithmic search, insertion, and deletion.'
  },
  {
    id: 'top-15',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 15,
    title: 'Graphs',
    slug: 'graphs',
    description: 'Adjacency lists, BFS, DFS, Dijkstra, Bellman-Ford, Floyd-Warshall, Topological Sort, and Disjoint Set (DSU).',
    difficulty: 'Hard',
    estimatedHours: 28,
    totalProblems: 54,
    problem_count: 54,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Binary Trees', 'Stack / Queues'],
    subtopics: ['Learning Graph', 'BFS / DFS Problems', 'Topo Sort', 'Shortest Paths', 'MST / Disjoint Set'],
    explanation: 'Graphs model general relational networks. Traversal, cycle detection, and shortest paths are foundational computer science problems.'
  },
  {
    id: 'top-16',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 16,
    title: 'Dynamic Programming',
    slug: 'dynamic-programming',
    description: 'Overlapping subproblems, 1D & 2D memoization, tabulation, 0/1 knapsack, LCS, LIS, and partition DP.',
    difficulty: 'Hard',
    estimatedHours: 36,
    totalProblems: 56,
    problem_count: 56,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Recursion'],
    subtopics: ['1D DP', '2D / Grid DP', 'DP on Subsequences', 'DP on Strings', 'DP on Stocks', 'DP on LIS'],
    explanation: 'Dynamic Programming caches answers to overlapping sub-problems, trading memory for polynomial time complexity.'
  },
  {
    id: 'top-17',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 17,
    title: 'Tries',
    slug: 'tries',
    description: 'Prefix trees, insert, search, startsWith, maximum XOR of two numbers in an array, and bitwise tries.',
    difficulty: 'Medium',
    estimatedHours: 10,
    totalProblems: 7,
    problem_count: 7,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Binary Trees', 'Bit Manipulation'],
    subtopics: ['Theory', 'Problems'],
    explanation: 'Tries provide prefix matching in O(L) time where L is word length, powering search engine autocompletes.'
  },
  {
    id: 'top-18',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 18,
    title: 'Hashing',
    slug: 'hashing',
    description: 'Hash tables, collision handling, open addressing, rolling hashes, and sub-array frequency maps.',
    difficulty: 'Easy',
    estimatedHours: 8,
    totalProblems: 15,
    problem_count: 15,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Arrays'],
    subtopics: ['Hashing Basics', 'Frequency Mapping', 'Subarray Sum Equals K'],
    explanation: 'Hash maps trade memory for O(1) expected operations, enabling rapid frequency lookups and duplicate detection.'
  },
  {
    id: 'top-19',
    roadmap_id: DEFAULT_ROADMAP_ID,
    topic_order: 19,
    title: 'Maths',
    slug: 'maths',
    description: 'GCD/HCF Euclidean algorithm, Sieve of Eratosthenes, prime factorization, modular exponentiation, and combinatorics.',
    difficulty: 'Medium',
    estimatedHours: 10,
    totalProblems: 12,
    problem_count: 12,
    completedProblems: 0,
    status: 'in_progress',
    prerequisites: ['Beginner Problems'],
    subtopics: ['Prime Numbers', 'Modular Arithmetic', 'Combinatorics'],
    explanation: 'Number theory fundamentals provide fast primality testing and modular arithmetic needed in competitive programming.'
  }
];

export const SAMPLE_PROBLEMS: Problem[] = [
  {
    id: '9c1e2a44-1111-4aaa-8bbb-0123456789ab',
    title: 'Pattern 1',
    slug: 'pattern-1',
    topic_id: 'top-1',
    topicId: 'top-1',
    difficulty: 'Easy',
    accuracyRate: 94,
    problem_order: 1,
    external_platform: 'takeuforward',
    external_problem_url: 'https://takeuforward.org/practice/dsa/pattern-1?category=patterns&source=strivers-a2z-dsa-sheet',
    source_sheet: 'strivers-a2z-dsa-sheet',
    source_problem_id: 'a2z-401',
    tags: ['basic', 'Beginner Problems', 'Patterns'],
    progress_status: 'Not Started',
    summary: "Practice problem from Striver's A2Z DSA Sheet (Beginner Problems — Patterns).",
    description: `Given an integer \`N\`, print an \`N x N\` grid of asterisks \`*\`.

Input Format: A single integer \`N\`.
Output Format: Grid of \`N\` rows with \`N\` stars each.`,
    constraints: ['1 <= N <= 100'],
    examples: [
      { input: 'N = 3', output: '***\n***\n***' },
      { input: 'N = 4', output: '****\n****\n****\n****' }
    ],
    hints: ['Use two nested loops: outer loop runs N times for rows, inner loop runs N times for columns.'],
    starterCode: {
      python: `def printPattern(n: int) -> None:\n    for i in range(n):\n        print('*' * n)\n`,
      cpp: `#include <iostream>\nusing namespace std;\n\nvoid printPattern(int n) {\n    for (int i = 0; i < n; i++) {\n        for (int j = 0; j < n; j++) cout << "*";\n        cout << "\\n";\n    }\n}`,
      java: `class Solution {\n    public static void printPattern(int n) {\n        for (int i = 0; i < n; i++) {\n            for (int j = 0; j < n; j++) System.out.print("*");\n            System.out.println();\n        }\n    }\n}`,
      c: `#include <stdio.h>\nvoid printPattern(int n) {\n    for(int i=0; i<n; i++) {\n        for(int j=0; j<n; j++) printf("*");\n        printf("\\n");\n    }\n}`,
      go: `package main\nimport "fmt"\n\nfunc printPattern(n int) {\n    for i := 0; i < n; i++ {\n        for j := 0; j < n; j++ { fmt.Print("*") }\n        fmt.Println()\n    }\n}`
    },
    testCases: [
      { id: 'tc-1', input: 'N = 3', expectedOutput: '***\n***\n***' },
      { id: 'tc-2', input: 'N = 4', expectedOutput: '****\n****\n****\n****' }
    ]
  },
  {
    id: 'prob-two-sum',
    title: 'Two Sum',
    slug: 'two-sum',
    topic_id: 'top-3',
    topicId: 'top-3',
    difficulty: 'Easy',
    accuracyRate: 88,
    problem_order: 1,
    external_platform: 'takeuforward',
    external_problem_url: 'https://takeuforward.org/data-structure/two-sum-check-if-a-pair-with-given-sum-exists-in-array/',
    source_sheet: 'strivers-a2z-dsa-sheet',
    source_problem_id: 'a2z-402',
    tags: ['Arrays', 'Hashing', 'Medium Arrays'],
    progress_status: 'Not Started',
    summary: "Two Sum check if a pair with given sum exists in array.",
    description: `Given an array of integers \`nums\` and an integer \`target\`, return indices of the two numbers such that they add up to \`target\`.

You may assume that each input would have exactly one solution, and you may not use the same element twice.`,
    constraints: [
      '2 <= nums.length <= 10^4',
      '-10^9 <= nums[i] <= 10^9',
      'Only one valid answer exists.'
    ],
    examples: [
      { input: 'nums = [2,7,11,15], target = 9', output: '[0,1]' },
      { input: 'nums = [3,2,4], target = 6', output: '[1,2]' }
    ],
    hints: ['Can you use a Hash Map to check if target - current_number has been seen in O(1) time?'],
    starterCode: {
      python: `class Solution:\n    def twoSum(self, nums: list[int], target: int) -> list[int]:\n        seen = {}\n        for i, num in enumerate(nums):\n            comp = target - num\n            if comp in seen:\n                return [seen[comp], i]\n            seen[num] = i\n        return []\n`,
      cpp: `#include <vector>\n#include <unordered_map>\nusing namespace std;\n\nclass Solution {\npublic:\n    vector<int> twoSum(vector<int>& nums, int target) {\n        unordered_map<int, int> seen;\n        for (int i = 0; i < nums.size(); i++) {\n            int comp = target - nums[i];\n            if (seen.count(comp)) return {seen[comp], i};\n            seen[nums[i]] = i;\n        }\n        return {};\n    }\n};`,
      java: `import java.util.HashMap;\nimport java.util.Map;\n\nclass Solution {\n    public int[] twoSum(int[] nums, int target) {\n        Map<Integer, Integer> seen = new HashMap<>();\n        for (int i = 0; i < nums.length; i++) {\n            int comp = target - nums[i];\n            if (seen.containsKey(comp)) return new int[] { seen.get(comp), i };\n            seen.put(nums[i], i);\n        }\n        return new int[] {};\n    }\n}`,
      c: `int* twoSum(int* nums, int numsSize, int target, int* returnSize) {\n    *returnSize = 2;\n    int* res = (int*)malloc(2 * sizeof(int));\n    for (int i=0; i<numsSize; i++) {\n        for (int j=i+1; j<numsSize; j++) {\n            if (nums[i] + nums[j] == target) {\n                res[0] = i; res[1] = j;\n                return res;\n            }\n        }\n    }\n    return res;\n}`,
      go: `package main\n\nfunc twoSum(nums []int, target int) []int {\n    seen := make(map[int]int)\n    for i, num := range nums {\n        if idx, ok := seen[target-num]; ok {\n            return []int{idx, i}\n        }\n        seen[num] = i\n    }\n    return nil\n}`
    },
    testCases: [
      { id: 'tc-1', input: 'nums = [2,7,11,15], target = 9', expectedOutput: '[0,1]' },
      { id: 'tc-2', input: 'nums = [3,2,4], target = 6', expectedOutput: '[1,2]' }
    ]
  },
  {
    id: 'prob-max-subarray',
    title: 'Maximum Subarray (Kadane)',
    slug: 'maximum-subarray',
    topic_id: 'top-3',
    topicId: 'top-3',
    difficulty: 'Medium',
    accuracyRate: 72,
    problem_order: 2,
    external_platform: 'takeuforward',
    external_problem_url: 'https://takeuforward.org/data-structure/kadanes-algorithm-maximum-subarray-sum-in-an-array/',
    source_sheet: 'strivers-a2z-dsa-sheet',
    source_problem_id: 'a2z-403',
    tags: ['Arrays', 'Kadane', 'Medium Arrays'],
    progress_status: 'Attempted',
    summary: "Kadane's Algorithm : Maximum Subarray Sum in an Array.",
    description: `Given an integer array \`nums\`, find the contiguous subarray with the largest sum and return its sum.`,
    constraints: ['1 <= nums.length <= 10^5', '-10^4 <= nums[i] <= 10^4'],
    examples: [
      { input: 'nums = [-2,1,-3,4,-1,2,1,-5,4]', output: '6' },
      { input: 'nums = [1]', output: '1' }
    ],
    hints: ['Keep a running sum. If running sum falls below 0, reset it to 0.'],
    starterCode: {
      python: `class Solution:\n    def maxSubArray(self, nums: list[int]) -> int:\n        max_sum = nums[0]\n        curr = 0\n        for x in nums:\n            curr = max(x, curr + x)\n            max_sum = max(max_sum, curr)\n        return max_sum\n`,
      cpp: `#include <vector>\n#include <algorithm>\nusing namespace std;\n\nclass Solution {\npublic:\n    int maxSubArray(vector<int>& nums) {\n        int maxSum = nums[0], curr = 0;\n        for (int x : nums) {\n            curr = max(x, curr + x);\n            maxSum = max(maxSum, curr);\n        }\n        return maxSum;\n    }\n};`,
      java: `class Solution {\n    public int maxSubArray(int[] nums) {\n        int maxSum = nums[0], curr = 0;\n        for (int x : nums) {\n            curr = Math.max(x, curr + x);\n            maxSum = Math.max(maxSum, curr);\n        }\n        return maxSum;\n    }\n}`,
      c: `int maxSubArray(int* nums, int numsSize) {\n    int maxSum = nums[0], curr = 0;\n    for(int i=0; i<numsSize; i++) {\n        curr = (nums[i] > curr + nums[i]) ? nums[i] : (curr + nums[i]);\n        if (curr > maxSum) maxSum = curr;\n    }\n    return maxSum;\n}`,
      go: `package main\n\nfunc maxSubArray(nums []int) int {\n    maxSum := nums[0]\n    curr := 0\n    for _, x := range nums {\n        if x > curr+x { curr = x } else { curr += x }\n        if curr > maxSum { maxSum = curr }\n    }\n    return maxSum\n}`
    },
    testCases: [
      { id: 'tc-1', input: 'nums = [-2,1,-3,4,-1,2,1,-5,4]', expectedOutput: '6' },
      { id: 'tc-2', input: 'nums = [1]', expectedOutput: '1' }
    ]
  }
];

export const COLLEGES = [
  'Indian Institute of Technology (IIT) Bombay',
  'Indian Institute of Technology (IIT) Delhi',
  'National Institute of Technology (NIT) Trichy',
  'BITS Pilani',
  'Delhi Technological University (DTU)',
  'Vellore Institute of Technology (VIT)',
  'Stanford University',
  'Carnegie Mellon University',
  'Other / Self-Taught'
];

export const COURSES = [
  'B.Tech / B.E. Computer Science',
  'B.Tech Information Technology',
  'B.Tech Electronics & Communication',
  'B.Sc Computer Science / Data Science',
  'BCA / MCA',
  'M.Tech Computer Science',
  'Non-CS Degree / Career Switcher'
];

export const DEPARTMENTS = [
  'Computer Science & Engineering',
  'Information Technology',
  'Artificial Intelligence & Data Science',
  'Electrical & Computer Engineering',
  'Software Systems',
  'Mathematics & Computing',
  'Other'
];

export const ACADEMIC_YEARS = [
  '1st Year',
  '2nd Year',
  '3rd Year',
  '4th Year'
];

export const LEARNING_GOALS = [
  { id: 'Interview Preparation', title: 'Interview Preparation', desc: 'Targeting product companies & campus placements' },
  { id: 'DSA', title: 'DSA Core Foundations', desc: 'Mastering algorithm patterns and data structure fundamentals' },
  { id: 'Competitive Programming', title: 'Competitive Programming', desc: 'Speed problem-solving for Codeforces and contest rounds' }
];

export const DSA_LEVELS = [
  { id: 'Beginner', title: 'Beginner', desc: 'Starting DSA fundamentals, loops, arrays, and basic conditionals' },
  { id: 'Intermediate', title: 'Intermediate', desc: 'Comfortable with recursion, two-pointers, stacks, and binary search' },
  { id: 'Advanced', title: 'Advanced', desc: 'Optimizing graph algorithms, dynamic programming, and hard edge cases' }
];

export const PROGRAMMING_LANGUAGES = [
  { id: 'Python', name: 'Python', version: 'Python 3.12', ext: 'py' },
  { id: 'Go', name: 'Go', version: 'Go 1.23', ext: 'go' },
  { id: 'Java', name: 'Java', version: 'OpenJDK 21', ext: 'java' },
  { id: 'C++', name: 'C++', version: 'GCC C++20', ext: 'cpp' },
  { id: 'C', name: 'C', version: 'Clang C17', ext: 'c' }
];
