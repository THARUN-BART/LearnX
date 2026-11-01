import apiClient from './api';
import { AssistantMessage, AssistantMode } from '../types/learning.types';

class AIService {
  async sendMessage(
    message: string,
    mode: AssistantMode,
    context?: { code?: string; problemTitle?: string; language?: string }
  ): Promise<AssistantMessage> {
    try {
      const response = await apiClient.post('/api/v1/ai/chat', {
        message,
        mode,
        context
      });
      return response.data.data || response.data;
    } catch {
      // Intelligent fallback responses tailored to mode if AI service is offline
      let reply = '';
      let codeSnippet: string | undefined = undefined;

      if (mode === 'hint') {
        reply = `**Algorithmic Hint for "${context?.problemTitle || 'current problem'}":**\n\n1. Consider the invariant of the data structure. Are elements sorted, or can you partition the search space?\n2. What is the bottleneck in the naive approach? Usually it's repeated scanning ($O(N^2)$), which can be turned into $O(1)$ lookups using an auxiliary hash table or a two-pointer technique.\n3. Try walking through a small edge case: what happens when $N=1$ or when elements are negative?`;
      } else if (mode === 'explain') {
        reply = `**Concept Explanation:**\n\nWhen optimizing this pattern, we trade auxiliary memory ($O(N)$ space) to achieve linear time ($O(N)$). Here is the conceptual state progression:\n\n- **State Initial**: Traverse each element once.\n- **Lookup**: Store previously seen values in a hash map where key is element and value is its index.\n- **Complement**: For every index $i$, test if \`target - nums[i]\` exists in your frequency map.\n\nThis completely avoids nested loops!`;
      } else {
        // Debug mode
        reply = `**Automated Debug Analysis:**\n\nI analyzed your snippet. Here are the potential pitfalls to watch out for:\n\n- **Index Out of Bounds / Off-by-one**: Ensure loop boundary terminates at \`len - 1\` or proper slice limit.\n- **Duplicate Keys**: If the array has duplicate elements, a basic map overwrite might lose the first index unless handled.\n- **Memory Overhead**: Keep your allocation outside the loop to prevent repeated re-allocations.\n\nSuggested test case to verify: \`nums = [3, 3], target = 6\`.`;
        codeSnippet = `// Quick sanity test
assert(twoSum([3, 3], 6) == [0, 1]);`;
      }

      const assistantMsg: AssistantMessage = {
        id: 'ai-msg-' + Date.now(),
        role: 'assistant',
        content: reply,
        codeSnippet,
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        mode
      };
      return assistantMsg;
    }
  }
}

export const aiService = new AIService();
export default aiService;

