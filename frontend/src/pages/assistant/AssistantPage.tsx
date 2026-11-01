import React, { useState } from 'react';
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
import { AssistantMessage, AssistantMode } from '../../types/learning.types';
import aiService from '../../services/ai.service';
import Button from '../../components/ui/Button';
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
      const response = await aiService.sendMessage(text, currentMode);
      setMessages((prev) => [...prev, response]);
    } catch {
      toast.error('Failed to receive AI response');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopyCode = (snippet: string, id: string) => {
    navigator.clipboard.writeText(snippet);
    setCopiedId(id);
    toast.success('Code copied!');
    setTimeout(() => setCopiedId(null), 2000);
  };

  const handleNewConversation = () => {
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
      <div className="flex flex-wrap items-center justify-between gap-3 bg-[#0D0D14] border border-[#1C1C28] px-4 py-2.5 rounded-xl shrink-0">
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-1.5 p-1 bg-[#12121B] rounded-xl border border-[#1C1C28]">
            <button
              onClick={() => setMode('hint')}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-mono transition-all cursor-pointer ${
                mode === 'hint'
                  ? 'bg-[#8B5CF6] text-white shadow-sm'
                  : 'text-[#8E8EA0] hover:text-white'
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
                  ? 'bg-[#6366F1] text-white shadow-sm'
                  : 'text-[#8E8EA0] hover:text-white'
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
                  : 'text-[#8E8EA0] hover:text-white'
              }`}
            >
              <Bug className="w-3.5 h-3.5" />
              <span>Debug Mode</span>
              <kbd className="keycap hidden sm:inline">3</kbd>
            </button>
          </div>
        </div>

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

      {/* Messages Scroll Area */}
      <div className="flex-1 bg-[#0D0D14] border border-[#1C1C28] rounded-xl p-4 sm:p-6 overflow-y-auto space-y-6">
        {messages.map((msg) => {
          const isUser = msg.role === 'user';
          return (
            <div
              key={msg.id}
              className={`flex gap-3 max-w-3xl ${isUser ? 'ml-auto justify-end' : ''}`}
            >
              {!isUser && (
                <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-[#8B5CF6] to-[#6366F1] flex items-center justify-center text-white shrink-0 mt-0.5">
                  <Bot className="w-4 h-4" />
                </div>
              )}

              <div
                className={`p-4 rounded-2xl text-xs space-y-3 leading-relaxed ${
                  isUser
                    ? 'bg-[#8B5CF6] text-white rounded-tr-none'
                    : 'bg-[#12121B] text-[#F4F4F6] border border-[#1C1C28] rounded-tl-none font-sans'
                }`}
              >
                {!isUser && msg.mode && (
                  <div className="flex items-center gap-2 pb-1.5 border-b border-[#1C1C28] text-[10px] font-mono text-[#A78BFA] uppercase">
                    <span>LearnX Copilot // {msg.mode}</span>
                  </div>
                )}

                <div className="whitespace-pre-wrap">{msg.content}</div>

                {msg.codeSnippet && (
                  <div className="mt-3 bg-[#08080C] rounded-xl border border-[#1C1C28] p-3.5 font-mono text-xs relative text-emerald-300">
                    <button
                      onClick={() => handleCopyCode(msg.codeSnippet!, msg.id)}
                      className="absolute top-2.5 right-2.5 p-1 rounded bg-[#12121B] text-[#8E8EA0] hover:text-white border border-[#1C1C28] cursor-pointer"
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
                <div className="w-8 h-8 rounded-xl bg-[#1C1C28] flex items-center justify-center text-[#8E8EA0] shrink-0 mt-0.5 font-mono text-xs font-bold">
                  U
                </div>
              )}
            </div>
          );
        })}

        {isLoading && (
          <div className="flex gap-3 max-w-xl font-mono">
            <div className="w-8 h-8 rounded-xl bg-[#8B5CF6]/20 border border-[#8B5CF6]/40 flex items-center justify-center text-[#A78BFA]">
              <Sparkles className="w-4 h-4 animate-spin" />
            </div>
            <div className="p-4 bg-[#12121B] border border-[#1C1C28] rounded-2xl text-xs text-[#8E8EA0] flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-[#8B5CF6] animate-ping" />
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
              className="p-2.5 bg-[#0D0D14] hover:bg-[#12121B] border border-[#1C1C28] rounded-xl text-left text-[11px] text-[#8E8EA0] hover:text-[#F4F4F6] transition-all truncate flex items-center gap-2 cursor-pointer font-sans"
            >
              <Sparkles className="w-3 h-3 text-[#8B5CF6] shrink-0" />
              <span className="truncate">{p.text}</span>
            </button>
          ))}
        </div>
      )}

      {/* Message Input Box */}
      <div className="bg-[#0D0D14] border border-[#1C1C28] rounded-xl p-2.5 flex items-center gap-2 shrink-0">
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
          className="flex-1 bg-transparent px-3 py-1.5 text-xs text-[#F4F4F6] placeholder-[#5A5A6E] outline-none resize-none font-mono"
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
