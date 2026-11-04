import React, { useEffect, useState } from 'react';
import {
  Bot,
  Send,
  Sparkles,
  HelpCircle,
  BookOpen,
  Bug,
  Copy,
  Check,
  Terminal,
  User,
  Plus,
  Cpu,
  Layers
} from 'lucide-react';
import { AIConversation, AssistantMessage, AssistantMode } from '../../types/learning.types';
import aiService from '../../services/ai.service';
import Button from '../../components/ui/Button';
import FormattedMarkdown from '../../components/common/FormattedMarkdown';
import { toast } from 'sonner';

export const AssistantPage: React.FC = () => {
  const [mode, setMode] = useState<AssistantMode>('hint');
  const [messages, setMessages] = useState<AssistantMessage[]>([
    {
      id: 'init-msg',
      role: 'assistant',
      content:
        'Hello! I am your LearnX Adaptive Coding Mentor. Choose a mode below and ask for hints, algorithmic derivations, or code debugging.',
      timestamp: 'Now',
      mode: 'hint'
    }
  ]);
  const [inputText, setInputText] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [conversationId, setConversationId] = useState<string | undefined>();
  const [conversations, setConversations] = useState<AIConversation[]>([]);

  useEffect(() => {
    aiService.listConversations().then(setConversations).catch(() => undefined);
  }, []);

  const suggestedPrompts = [
    { text: 'Explain the difference between Dynamic Programming and Memoization.', mode: 'explain' as AssistantMode },
    { text: 'Give me a hint for solving 3Sum without O(N^3) time.', mode: 'hint' as AssistantMode },
    { text: 'Why does my two-pointer loop cause an IndexError on single-element arrays?', mode: 'debug' as AssistantMode },
    { text: 'How do Monotonic Stacks compute Next Greater Element in linear time?', mode: 'explain' as AssistantMode }
  ];

  const handleSend = async (textToSend?: string, overrideMode?: AssistantMode) => {
    const text = textToSend || inputText;
    if (!text.trim()) return;

    const currentMode = overrideMode || mode;

    const userMsg: AssistantMessage = {
      id: 'usr-' + Date.now(),
      role: 'user',
      content: text,
      timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      mode: currentMode
    };

    setMessages((prev) => [...prev, userMsg]);
    setInputText('');
    setIsLoading(true);

    try {
      const response = await aiService.sendMessage(text, currentMode, {
        conversationId,
      });
      if (response.conversationId) {
        setConversationId(response.conversationId);
      }
      setMessages((prev) => [...prev, response]);
    } catch {
      toast.error('Failed to receive AI response');
    } finally {
      setIsLoading(false);
    }
  };

  const loadConversation = async (id: string) => {
    setConversationId(id);
    const history = await aiService.getConversationMessages(id);
    if (history.length) {
      setMessages(history);
    }
  };

  const handleCopyCode = (snippet: string, id: string) => {
    navigator.clipboard.writeText(snippet);
    setCopiedId(id);
    toast.success('Code copied!');
    setTimeout(() => setCopiedId(null), 2000);
  };

  const handleNewConversation = () => {
    setConversationId(undefined);
    setMessages([
      {
        id: 'new-init',
        role: 'assistant',
        content: 'New session initialized. How can I assist your coding today?',
        timestamp: 'Now',
        mode
      }
    ]);
  };

  return (
    <div className="h-[calc(100vh-6.5rem)] flex flex-col gap-3 max-w-5xl mx-auto font-mono">
      {/* Header Toolbar: Mode selector & New Chat */}
      <div className="flex flex-wrap items-center justify-between gap-3 bg-white border border-[#E2E8F0] px-4 py-2.5 rounded-xl shrink-0">
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-1.5 p-1 bg-[#F1F5F9] rounded-xl border border-[#E2E8F0]">
            <button
              onClick={() => setMode('hint')}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-mono transition-all cursor-pointer ${
                mode === 'hint'
                  ? 'bg-[#EF4444] text-white shadow-sm'
                  : 'text-[#64748B] hover:text-[#0F172A]'
              }`}
            >
              <HelpCircle className="w-3.5 h-3.5" />
              <span>Hint Mode</span>
              <kbd className="keycap hidden sm:inline">1</kbd>
            </button>
            <button
              onClick={() => setMode('explain')}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-mono transition-all cursor-pointer ${
                mode === 'explain'
                  ? 'bg-[#DC2626] text-white shadow-sm'
                  : 'text-[#64748B] hover:text-[#0F172A]'
              }`}
            >
              <BookOpen className="w-3.5 h-3.5" />
              <span>Explain Mode</span>
              <kbd className="keycap hidden sm:inline">2</kbd>
            </button>
            <button
              onClick={() => setMode('debug')}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-mono transition-all cursor-pointer ${
                mode === 'debug'
                  ? 'bg-amber-600 text-white shadow-sm'
                  : 'text-[#64748B] hover:text-[#0F172A]'
              }`}
            >
              <Bug className="w-3.5 h-3.5" />
              <span>Debug Mode</span>
              <kbd className="keycap hidden sm:inline">3</kbd>
            </button>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {conversations.length > 0 && (
            <select
              value={conversationId || ''}
              onChange={(e) => {
                if (e.target.value) loadConversation(e.target.value);
              }}
              className="bg-[#F1F5F9] text-[#0F172A] text-[11px] font-mono px-2 py-1.5 rounded-lg border border-[#E2E8F0] max-w-[180px]"
            >
              <option value="">Past sessions</option>
              {conversations.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.title || 'Conversation'}
                </option>
              ))}
            </select>
          )}
          <Button
            variant="secondary"
            size="sm"
            onClick={handleNewConversation}
            leftIcon={<Plus className="w-3.5 h-3.5" />}
            className="text-xs font-mono h-8"
          >
            New Session
          </Button>
        </div>
      </div>

      {/* Messages Scroll Area */}
      <div className="flex-1 bg-white border border-[#E2E8F0] rounded-xl p-4 sm:p-6 overflow-y-auto space-y-6">
        {messages.map((msg) => {
          const isUser = msg.role === 'user';
          return (
            <div
              key={msg.id}
              className={`flex gap-3 max-w-3xl ${isUser ? 'ml-auto justify-end' : ''}`}
            >
              {!isUser && (
                <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-[#EF4444] to-[#B91C1C] flex items-center justify-center text-white shrink-0 mt-0.5">
                  <Bot className="w-4 h-4" />
                </div>
              )}

              <div
                className={`p-4 rounded-2xl text-xs space-y-3 leading-relaxed ${
                  isUser
                    ? 'bg-[#EF4444] text-white rounded-tr-none'
                    : 'bg-[#F1F5F9] text-[#0F172A] border border-[#E2E8F0] rounded-tl-none font-sans'
                }`}
              >
                {!isUser && msg.mode && (
                  <div className="flex items-center gap-2 pb-1.5 border-b border-[#E2E8F0] text-[10px] font-mono text-[#DC2626] uppercase">
                    <span>LearnX Copilot // {msg.mode}</span>
                  </div>
                )}

                {isUser ? (
                  <div className="whitespace-pre-wrap">{msg.content}</div>
                ) : (
                  <FormattedMarkdown content={msg.content} />
                )}

                {msg.codeSnippet && (
                  <div className="mt-3 bg-[#F8FAFC] rounded-xl border border-[#E2E8F0] p-3.5 font-mono text-xs relative text-emerald-300">
                    <button
                      onClick={() => handleCopyCode(msg.codeSnippet!, msg.id)}
                      className="absolute top-2.5 right-2.5 p-1 rounded bg-[#F1F5F9] text-[#64748B] hover:text-[#0F172A] border border-[#E2E8F0] cursor-pointer"
                    >
                      {copiedId === msg.id ? (
                        <Check className="w-3 h-3 text-emerald-400" />
                      ) : (
                        <Copy className="w-3 h-3" />
                      )}
                    </button>
                    <pre className="overflow-x-auto">{msg.codeSnippet}</pre>
                  </div>
                )}

                <span className="text-[10px] opacity-60 block text-right font-mono">
                  {msg.timestamp}
                </span>
              </div>

              {isUser && (
                <div className="w-8 h-8 rounded-xl bg-[#F1F5F9] flex items-center justify-center text-[#64748B] shrink-0 mt-0.5 font-mono text-xs font-bold">
                  U
                </div>
              )}
            </div>
          );
        })}

        {isLoading && (
          <div className="flex gap-3 max-w-xl font-mono">
            <div className="w-8 h-8 rounded-xl bg-[#EF4444]/20 border border-[#EF4444]/40 flex items-center justify-center text-[#DC2626]">
              <Sparkles className="w-4 h-4 animate-spin" />
            </div>
            <div className="p-4 bg-[#F1F5F9] border border-[#E2E8F0] rounded-2xl text-xs text-[#64748B] flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-[#EF4444] animate-ping" />
              <span>Synthesizing cognitive advice...</span>
            </div>
          </div>
        )}
      </div>

      {/* Suggested Prompts Pill Row (when few messages) */}
      {messages.length <= 2 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 shrink-0">
          {suggestedPrompts.map((p, idx) => (
            <button
              key={idx}
              onClick={() => handleSend(p.text, p.mode)}
              className="p-2.5 bg-white hover:bg-[#F1F5F9] border border-[#E2E8F0] rounded-xl text-left text-[11px] text-[#64748B] hover:text-[#0F172A] transition-all truncate flex items-center gap-2 cursor-pointer font-sans"
            >
              <Sparkles className="w-3 h-3 text-[#EF4444] shrink-0" />
              <span className="truncate">{p.text}</span>
            </button>
          ))}
        </div>
      )}

      {/* Message Input Box */}
      <div className="bg-white border border-[#E2E8F0] rounded-xl p-2.5 flex items-center gap-2 shrink-0">
        <textarea
          rows={1}
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
              e.preventDefault();
              handleSend();
            }
          }}
          placeholder={`Formulate prompt in ${mode} mode... (Enter to send)`}
          className="flex-1 bg-transparent px-3 py-1.5 text-xs text-[#0F172A] placeholder-[#94A3B8] outline-none resize-none font-mono"
        />
        <Button
          variant="primary"
          size="sm"
          onClick={() => handleSend()}
          isLoading={isLoading}
          rightIcon={<Send className="w-3 h-3" />}
          className="text-xs font-mono h-8"
        >
          Send
        </Button>
      </div>
    </div>
  );
};

export default AssistantPage;
