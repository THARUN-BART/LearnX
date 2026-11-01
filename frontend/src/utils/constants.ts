import { Topic, Problem } from '../types/learning.types';

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
export const TOKEN_KEY = 'learnx_access_token';
export const REFRESH_TOKEN_KEY = 'learnx_refresh_token';
export const USER_KEY = 'learnx_auth_user';
export const PROFILE_KEY = 'learnx_student_profile';

export const DSA_TOPICS: Topic[] = [
  {
    id: 'prog-fundamentals',
    title: 'Programming Fundamentals',
    slug: 'programming-fundamentals',
    description: 'Master variables, memory layouts, pointers, scope, conditional flow, and bitwise operations.',
    difficulty: 'Easy',
    estimatedHours: 8,
    totalProblems: 12,
    completedProblems: 12,
    status: 'completed',
    prerequisites: ['Basic math', 'Logic'],
    subtopics: ['Variables & Memory', 'Control Structures', 'Functions & Call Stack', 'Bitwise Operators'],
    explanation: 'Programming fundamentals lay the rock-solid bedrock for writing efficient algorithms. Understanding how the memory stack and heap operate allows you to reason about cache locality and pointer arithmetic.'
  },
  {
    id: 'complexity-analysis',
    title: 'Complexity Analysis',
    slug: 'complexity-analysis',
    description: 'Asymptotic notation, Big-O, Big-Theta, Big-Omega, amortized analysis, and recurrence relations.',
    difficulty: 'Easy',
    estimatedHours: 6,
    totalProblems: 8,
    completedProblems: 8,
    status: 'completed',
    prerequisites: ['Programming Fundamentals'],
    subtopics: ['Big-O Notation', 'Space Complexity', 'Amortized Analysis', 'Master Theorem'],
    explanation: 'Algorithm efficiency is measured using asymptotic analysis. You will learn how to analyze loops, nested iterations, recursive call stacks, and space overheads.'
  },
  {
    id: 'arrays',
    title: 'Arrays',
    slug: 'arrays',
    description: 'Contiguous memory structures, two pointers, sliding window, prefix sums, and Kadane algorithm.',
    difficulty: 'Easy',
    estimatedHours: 14,
    totalProblems: 24,
    completedProblems: 18,
    status: 'in_progress',
    prerequisites: ['Complexity Analysis'],
    subtopics: ['Dynamic Arrays', 'Two Pointers', 'Sliding Window', 'Prefix Sums', 'Matrix Traversal'],
    explanation: 'Arrays provide O(1) random access but require O(N) insertion/deletion. Advanced patterns like two-pointer traversal and sliding window transform O(N^2) brute force solutions into optimal O(N) routines.'
  },
  {
    id: 'strings',
    title: 'Strings',
    slug: 'strings',
    description: 'String manipulation, anagrams, palindrome verification, rolling hash, and KMP pattern matching.',
    difficulty: 'Easy',
    estimatedHours: 12,
    totalProblems: 18,
    completedProblems: 12,
    status: 'in_progress',
    prerequisites: ['Arrays'],
    subtopics: ['Immutable vs Mutable Strings', 'Anagram Detection', 'Palindrome Properties', 'Pattern Matching (KMP)'],
    explanation: 'Strings are arrays of characters with unique algorithmic challenges such as immutability in Java/Python and character encoding. Efficient string parsing is vital in interview challenges.'
  },
  {
    id: 'searching-sorting',
    title: 'Searching and Sorting',
    slug: 'searching-and-sorting',
    description: 'Binary search variations, search on answer space, QuickSort, MergeSort, and stability.',
    difficulty: 'Medium',
    estimatedHours: 16,
    totalProblems: 22,
    completedProblems: 14,
    status: 'in_progress',
    prerequisites: ['Arrays'],
    subtopics: ['Binary Search on Rotated Arrays', 'Binary Search on Answers', 'QuickSort & Partitioning', 'MergeSort & Count Inversions'],
    explanation: 'Sorting establishes order to unlock logarithmic searching. Binary search is not just for sorted arrays; it is a fundamental problem-solving technique for monotonic search spaces.'
  },
  {
    id: 'recursion',
    title: 'Recursion',
    slug: 'recursion',
    description: 'Call stacks, base cases, tree recursions, subset generation, and divide-and-conquer paradigms.',
    difficulty: 'Medium',
    estimatedHours: 12,
    totalProblems: 16,
    completedProblems: 7,
    status: 'in_progress',
    prerequisites: ['Programming Fundamentals'],
    subtopics: ['Stack Frames & Recursion Limits', 'Tree Recursion', 'Subsets & Permutations', 'Divide and Conquer'],
    explanation: 'Recursion breaks complex problems into self-similar sub-problems. Mastering base cases and unwinding state transitions is crucial for trees, graphs, and dynamic programming.'
  },
  {
    id: 'linked-lists',
    title: 'Linked Lists',
    slug: 'linked-lists',
    description: 'Singly, doubly, and circular linked lists, Fast & Slow pointer (Floyd cycle detection), and reordering.',
    difficulty: 'Medium',
    estimatedHours: 10,
    totalProblems: 16,
    completedProblems: 9,
    status: 'in_progress',
    prerequisites: ['Programming Fundamentals'],
    subtopics: ['Singly & Doubly Linked Lists', 'Reversal in Groups of K', 'Floyd Cycle Detection', 'Merge K Sorted Lists'],
    explanation: 'Linked lists teach explicit pointer manipulation without contiguous storage. Pointers can be tricky with edge cases like empty lists, single nodes, and cycle detection.'
  },
  {
    id: 'stacks-queues',
    title: 'Stacks and Queues',
    slug: 'stacks-and-queues',
    description: 'LIFO & FIFO primitives, Monotonic Stack, Next Greater Element, and circular deque design.',
    difficulty: 'Medium',
    estimatedHours: 12,
    totalProblems: 18,
    completedProblems: 8,
    status: 'in_progress',
    prerequisites: ['Arrays', 'Linked Lists'],
    subtopics: ['Monotonic Stack Technique', 'Largest Rectangle in Histogram', 'Sliding Window Maximum (Deque)', 'Min Stack Design'],
    explanation: 'Monotonic stacks and queues are some of the most powerful tools in competitive programming for maintaining running maximums/minimums in amortized O(1) time.'
  },
  {
    id: 'hashing',
    title: 'Hashing',
    slug: 'hashing',
    description: 'Hash tables, collision handling, open addressing, rolling hashes, and sub-array frequency maps.',
    difficulty: 'Easy',
    estimatedHours: 8,
    totalProblems: 15,
    completedProblems: 6,
    status: 'in_progress',
    prerequisites: ['Arrays'],
    subtopics: ['Hash Function Design', 'Collision Resolution (Chaining vs Open Addressing)', 'Subarray Sum Equals K', 'LRU Cache Design'],
    explanation: 'Hash maps trade memory for O(1) expected time operations. They form the backbone of database indices, caches, and fast lookup algorithms.'
  },
  {
    id: 'trees',
    title: 'Trees',
    slug: 'trees',
    description: 'Binary trees, traversals (in-order, pre-order, post-order, level-order), LCA, and maximum depth.',
    difficulty: 'Medium',
    estimatedHours: 18,
    totalProblems: 24,
    completedProblems: 5,
    status: 'in_progress',
    prerequisites: ['Recursion', 'Stacks-Queues'],
    subtopics: ['Tree Traversals (DFS & BFS)', 'Lowest Common Ancestor', 'Diameter & Height', 'Serialization & Deserialization'],
    explanation: 'Hierarchical data structures represent file systems, DOM trees, and network routing. Tree algorithms heavily reinforce recursive reasoning.'
  },
  {
    id: 'bst',
    title: 'Binary Search Trees',
    slug: 'binary-search-trees',
    description: 'BST validation, inorder property, AVL self-balancing trees, red-black trees, and range queries.',
    difficulty: 'Medium',
    estimatedHours: 12,
    totalProblems: 14,
    completedProblems: 2,
    status: 'in_progress',
    prerequisites: ['Trees'],
    subtopics: ['BST Invariants', 'Successor & Predecessor', 'K-th Smallest Element', 'Self-Balancing Concepts (AVL)'],
    explanation: 'Binary Search Trees maintain dynamic sorted order, enabling logarithmic search, insertion, and deletion.'
  },
  {
    id: 'heaps',
    title: 'Heaps & Priority Queues',
    slug: 'heaps',
    description: 'Min-heaps, max-heaps, heapify, top-K frequent elements, and median from data stream.',
    difficulty: 'Medium',
    estimatedHours: 10,
    totalProblems: 14,
    completedProblems: 1,
    status: 'in_progress',
    prerequisites: ['Arrays', 'Trees'],
    subtopics: ['Binary Heap Representation', 'Heapify O(N)', 'Top K Elements', 'Continuous Median Stream'],
    explanation: 'Heaps guarantee O(1) access to extremes (min or max) with logarithmic updates, powering Dijkstra algorithm and real-time scheduling.'
  },
  {
    id: 'graphs',
    title: 'Graphs',
    slug: 'graphs',
    description: 'Adjacency lists, BFS, DFS, Dijkstra, Bellman-Ford, Floyd-Warshall, Topological Sort, and Disjoint Set.',
    difficulty: 'Hard',
    estimatedHours: 24,
    totalProblems: 30,
    completedProblems: 0,
    status: 'locked',
    prerequisites: ['Trees', 'Stacks-Queues', 'Hashing'],
    subtopics: ['BFS/DFS Traversals', 'Cycle Detection in Directed/Undirected', 'Topological Sort (Kahn Algorithm)', 'Shortest Path (Dijkstra)', 'Union-Find (DSU)'],
    explanation: 'Graphs model real-world networks from social graphs to transportation grids. Graph traversal and shortest-path algorithms are gold-standard technical interview subjects.'
  },
  {
    id: 'greedy',
    title: 'Greedy Algorithms',
    slug: 'greedy-algorithms',
    description: 'Locally optimal choices, interval scheduling, fractional knapsack, and Huffman coding.',
    difficulty: 'Medium',
    estimatedHours: 12,
    totalProblems: 16,
    completedProblems: 0,
    status: 'locked',
    prerequisites: ['Searching and Sorting'],
    subtopics: ['Greedy Choice Property', 'Interval Merging & Scheduling', 'Gas Station / Jump Game', 'Minimum Spanning Trees'],
    explanation: 'Greedy algorithms construct solutions piece-by-piece, always picking the immediately obvious benefit. Proving correctness requires mathematical induction or exchange arguments.'
  },
  {
    id: 'dp',
    title: 'Dynamic Programming',
    slug: 'dynamic-programming',
    description: 'Overlapping subproblems, optimal substructure, 1D & 2D memoization, tabulation, knapsack, and LCS.',
    difficulty: 'Hard',
    estimatedHours: 32,
    totalProblems: 36,
    completedProblems: 0,
    status: 'locked',
    prerequisites: ['Recursion'],
    subtopics: ['Memoization vs Tabulation', '0/1 Knapsack & Variations', 'Longest Common Subsequence (LCS)', 'Longest Increasing Subsequence (LIS)', 'Matrix Chain DP'],
    explanation: 'Dynamic Programming solves complex problems by breaking them into simpler overlapping subproblems and caching their answers to eliminate redundant computation.'
  },
  {
    id: 'backtracking',
    title: 'Backtracking',
    slug: 'backtracking',
    description: 'State-space trees, N-Queens problem, Sudoku solver, word search, and pruning techniques.',
    difficulty: 'Hard',
    estimatedHours: 16,
    totalProblems: 18,
    completedProblems: 0,
    status: 'locked',
    prerequisites: ['Recursion'],
    subtopics: ['State Space Exploration', 'N-Queens Formulation', 'Sudoku Constraints & Pruning', 'Hamiltonian Paths'],
    explanation: 'Backtracking systematically searches all configurations in a combinatorial search space while abandoning ("backtracking") branches that violate constraints early.'
  }
];

export const SAMPLE_PROBLEMS: Problem[] = [
  {
    id: 'prob-two-sum',
    title: 'Two Sum',
    slug: 'two-sum',
    topicId: 'arrays',
    difficulty: 'Easy',
    accuracyRate: 88,
    description: `Given an array of integers \`nums\` and an integer \`target\`, return indices of the two numbers such that they add up to \`target\`.

You may assume that each input would have exactly one solution, and you may not use the same element twice.

You can return the answer in any order.`,
    constraints: [
      '2 <= nums.length <= 10^4',
      '-10^9 <= nums[i] <= 10^9',
      '-10^9 <= target <= 10^9',
      'Only one valid answer exists.'
    ],
    examples: [
      {
        input: 'nums = [2,7,11,15], target = 9',
        output: '[0,1]',
        explanation: 'Because nums[0] + nums[1] == 9, we return [0, 1].'
      },
      {
        input: 'nums = [3,2,4], target = 6',
        output: '[1,2]',
        explanation: 'Because nums[1] + nums[2] == 6, we return [1, 2].'
      }
    ],
    hints: [
      'A brute force approach would search all pairs in O(n^2) time. Can we do better with auxiliary memory?',
      'Can you use a Hash Map to check if the complement (target - current_number) has already been seen in O(1) time?'
    ],
    starterCode: {
      python: `class Solution:
    def twoSum(self, nums: list[int], target: int) -> list[int]:
        # Write your adaptive solution here
        seen = {}
        for i, num in enumerate(nums):
            complement = target - num
            if complement in seen:
                return [seen[complement], i]
            seen[num] = i
        return []
`,
      cpp: `#include <vector>
#include <unordered_map>
using namespace std;

class Solution {
public:
    vector<int> twoSum(vector<int>& nums, int target) {
        // Write your solution here
        unordered_map<int, int> seen;
        for (int i = 0; i < nums.size(); i++) {
            int complement = target - nums[i];
            if (seen.find(complement) != seen.end()) {
                return {seen[complement], i};
            }
            seen[nums[i]] = i;
        }
        return {};
    }
};`,
      java: `import java.util.HashMap;
import java.util.Map;

class Solution {
    public int[] twoSum(int[] nums, int target) {
        // Write your solution here
        Map<Integer, Integer> seen = new HashMap<>();
        for (int i = 0; i < nums.length; i++) {
            int complement = target - nums[i];
            if (seen.containsKey(complement)) {
                return new int[] { seen.get(complement), i };
            }
            seen.put(nums[i], i);
        }
        return new int[] {};
    }
}`,
      c: `/**
 * Note: The returned array must be malloced, assume caller calls free().
 */
int* twoSum(int* nums, int numsSize, int target, int* returnSize) {
    // Write C implementation
    *returnSize = 2;
    int* result = (int*)malloc(2 * sizeof(int));
    for (int i = 0; i < numsSize; i++) {
        for (int j = i + 1; j < numsSize; j++) {
            if (nums[i] + nums[j] == target) {
                result[0] = i;
                result[1] = j;
                return result;
            }
        }
    }
    return result;
}`,
      go: `package main

func twoSum(nums []int, target int) []int {
    seen := make(map[int]int)
    for i, num := range nums {
        complement := target - num
        if idx, ok := seen[complement]; ok {
            return []int{idx, i}
        }
        seen[num] = i
    }
    return nil
}`
    },
    testCases: [
      { id: 'tc-1', input: 'nums = [2,7,11,15], target = 9', expectedOutput: '[0,1]' },
      { id: 'tc-2', input: 'nums = [3,2,4], target = 6', expectedOutput: '[1,2]' },
      { id: 'tc-3', input: 'nums = [3,3], target = 6', expectedOutput: '[0,1]' }
    ]
  },
  {
    id: 'prob-max-subarray',
    title: 'Maximum Subarray (Kadane)',
    slug: 'maximum-subarray',
    topicId: 'arrays',
    difficulty: 'Medium',
    accuracyRate: 72,
    description: `Given an integer array \`nums\`, find the subarray with the largest sum, and return its sum.

A subarray is a contiguous non-empty sequence of elements within an array.`,
    constraints: [
      '1 <= nums.length <= 10^5',
      '-10^4 <= nums[i] <= 10^4'
    ],
    examples: [
      {
        input: 'nums = [-2,1,-3,4,-1,2,1,-5,4]',
        output: '6',
        explanation: 'The subarray [4,-1,2,1] has the largest sum 6.'
      },
      {
        input: 'nums = [1]',
        output: '1',
        explanation: 'The subarray [1] has the largest sum 1.'
      }
    ],
    hints: [
      'If you have already solved the problem for nums[0..i-1], how does nums[i] change the solution?',
      'Consider Kadane algorithm: at each index, you decide whether to append to the existing subarray or start fresh from current index.'
    ],
    starterCode: {
      python: `class Solution:
    def maxSubArray(self, nums: list[int]) -> int:
        # Implement Kadane's Algorithm
        max_sum = nums[0]
        curr_sum = 0
        for x in nums:
            curr_sum = max(x, curr_sum + x)
            max_sum = max(max_sum, curr_sum)
        return max_sum
`,
      cpp: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int maxSubArray(vector<int>& nums) {
        int maxSum = nums[0];
        int currSum = 0;
        for (int x : nums) {
            currSum = max(x, currSum + x);
            maxSum = max(maxSum, currSum);
        }
        return maxSum;
    }
};`,
      java: `class Solution {
    public int maxSubArray(int[] nums) {
        int maxSum = nums[0];
        int currSum = 0;
        for (int x : nums) {
            currSum = Math.max(x, currSum + x);
            maxSum = Math.max(maxSum, currSum);
        }
        return maxSum;
    }
}`,
      c: `int maxSubArray(int* nums, int numsSize) {
    int maxSum = nums[0];
    int currSum = 0;
    for (int i = 0; i < numsSize; i++) {
        currSum = (nums[i] > currSum + nums[i]) ? nums[i] : currSum + nums[i];
        if (currSum > maxSum) maxSum = currSum;
    }
    return maxSum;
}`,
      go: `package main

func maxSubArray(nums []int) int {
    maxSum := nums[0]
    currSum := 0
    for _, x := range nums {
        if x > currSum+x {
            currSum = x
        } else {
            currSum += x
        }
        if currSum > maxSum {
            maxSum = currSum
        }
    }
    return maxSum
}`
    },
    testCases: [
      { id: 'tc-1', input: 'nums = [-2,1,-3,4,-1,2,1,-5,4]', expectedOutput: '6' },
      { id: 'tc-2', input: 'nums = [1]', expectedOutput: '1' },
      { id: 'tc-3', input: 'nums = [5,4,-1,7,8]', expectedOutput: '23' }
    ]
  },
  {
    id: 'prob-valid-parentheses',
    title: 'Valid Parentheses',
    slug: 'valid-parentheses',
    topicId: 'stacks-queues',
    difficulty: 'Easy',
    accuracyRate: 85,
    description: `Given a string \`s\` containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.

An input string is valid if:
1. Open brackets must be closed by the same type of brackets.
2. Open brackets must be closed in the correct order.
3. Every close bracket has a corresponding open bracket of the same type.`,
    constraints: [
      '1 <= s.length <= 10^4',
      's consists of parentheses only \x27()[]{}\x27'
    ],
    examples: [
      { input: 's = "()"', output: 'true' },
      { input: 's = "()[]{}"', output: 'true' },
      { input: 's = "(]"', output: 'false' }
    ],
    hints: [
      'Use a stack to track open brackets. Whenever you see a closing bracket, check if it matches the top of stack.'
    ],
    starterCode: {
      python: `class Solution:
    def isValid(self, s: str) -> bool:
        stack = []
        mapping = {")": "(", "}": "{", "]": "["}
        for char in s:
            if char in mapping:
                top = stack.pop() if stack else '#'
                if mapping[char] != top:
                    return False
            else:
                stack.append(char)
        return not stack
`,
      cpp: `#include <string>
#include <stack>
using namespace std;

class Solution {
public:
    bool isValid(string s) {
        stack<char> st;
        for (char c : s) {
            if (c == '(' || c == '{' || c == '[') st.push(c);
            else {
                if (st.empty()) return false;
                if (c == ')' && st.top() != '(') return false;
                if (c == '}' && st.top() != '{') return false;
                if (c == ']' && st.top() != '[') return false;
                st.pop();
            }
        }
        return st.empty();
    }
};`,
      java: `import java.util.Stack;

class Solution {
    public boolean isValid(String s) {
        Stack<Character> stack = new Stack<>();
        for (char c : s.toCharArray()) {
            if (c == '(' || c == '{' || c == '[') stack.push(c);
            else {
                if (stack.isEmpty()) return false;
                char top = stack.pop();
                if (c == ')' && top != '(') return false;
                if (c == '}' && top != '{') return false;
                if (c == ']' && top != '[') return false;
            }
        }
        return stack.isEmpty();
    }
}`,
      c: `bool isValid(char* s) {
    int len = strlen(s);
    char* stack = (char*)malloc(len);
    int top = -1;
    for (int i = 0; i < len; i++) {
        char c = s[i];
        if (c == '(' || c == '{' || c == '[') {
            stack[++top] = c;
        } else {
            if (top < 0) { free(stack); return false; }
            char open = stack[top--];
            if (c == ')' && open != '(') { free(stack); return false; }
            if (c == '}' && open != '{') { free(stack); return false; }
            if (c == ']' && open != '[') { free(stack); return false; }
        }
    }
    bool valid = (top == -1);
    free(stack);
    return valid;
}`,
      go: `package main

func isValid(s string) bool {
    var stack []rune
    pairs := map[rune]rune{')': '(', '}': '{', ']': '['}
    for _, ch := range s {
        if open, exists := pairs[ch]; exists {
            if len(stack) == 0 || stack[len(stack)-1] != open {
                return false
            }
            stack = stack[:len(stack)-1]
        } else {
            stack = append(stack, ch)
        }
    }
    return len(stack) == 0
}`
    },
    testCases: [
      { id: 'tc-1', input: 's = "()"', expectedOutput: 'true' },
      { id: 'tc-2', input: 's = "()[]{}"', expectedOutput: 'true' },
      { id: 'tc-3', input: 's = "(]"', expectedOutput: 'false' }
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

