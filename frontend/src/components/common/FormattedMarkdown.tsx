import React, { useState } from 'react';
import { Copy, Check, Terminal } from 'lucide-react';

interface FormattedMarkdownProps {
  content: string;
  className?: string;
}

export const FormattedMarkdown: React.FC<FormattedMarkdownProps> = ({ content, className = '' }) => {
  const [copiedIndex, setCopiedIndex] = useState<number | null>(null);

  const handleCopy = (code: string, idx: number) => {
    navigator.clipboard.writeText(code);
    setCopiedIndex(idx);
    setTimeout(() => setCopiedIndex(null), 2000);
  };

  if (!content) return null;

  // Split by code blocks: ```lang ... ```
  const codeBlockRegex = /```([a-zA-Z0-9_-]*)\n([\s\S]*?)```/g;
  const parts: Array<{ type: 'text' | 'code'; content: string; language?: string }> = [];

  let lastIndex = 0;
  let match: RegExpExecArray | null;

  while ((match = codeBlockRegex.exec(content)) !== null) {
    if (match.index > lastIndex) {
      parts.push({
        type: 'text',
        content: content.slice(lastIndex, match.index),
      });
    }
    parts.push({
      type: 'code',
      language: match[1]?.trim() || 'code',
      content: match[2]?.trimEnd() || '',
    });
    lastIndex = codeBlockRegex.lastIndex;
  }

  if (lastIndex < content.length) {
    parts.push({
      type: 'text',
      content: content.slice(lastIndex),
    });
  }

  // Helper to format inline markdown (bold, italic, inline code)
  const renderInline = (text: string): React.ReactNode[] => {
    // Split by inline code: `code`
    const inlineCodeParts = text.split(/(`[^`]+`)/g);

    return inlineCodeParts.map((part, pIdx) => {
      if (part.startsWith('`') && part.endsWith('`') && part.length >= 2) {
        return (
          <code
            key={pIdx}
            className="px-1.5 py-0.5 mx-0.5 text-[11px] font-mono bg-[#161626] text-emerald-400 border border-[#27273F] rounded-md font-medium"
          >
            {part.slice(1, -1)}
          </code>
        );
      }

      // Format bold (**text** or __text__) and italic (*text* or _text_)
      const boldParts = part.split(/(\*\*[^*]+\*\*|__[^_]+__)/g);

      return (
        <span key={pIdx}>
          {boldParts.map((bPart, bIdx) => {
            if (
              (bPart.startsWith('**') && bPart.endsWith('**')) ||
              (bPart.startsWith('__') && bPart.endsWith('__'))
            ) {
              return (
                <strong key={bIdx} className="font-semibold text-white">
                  {bPart.slice(2, -2)}
                </strong>
              );
            }

            // Italic
            const italicParts = bPart.split(/(\*[^*]+\*|_[^_]+_)/g);
            return (
              <span key={bIdx}>
                {italicParts.map((iPart, iIdx) => {
                  if (
                    (iPart.startsWith('*') && iPart.endsWith('*')) ||
                    (iPart.startsWith('_') && iPart.endsWith('_'))
                  ) {
                    return (
                      <em key={iIdx} className="italic text-[#CBD5E1]">
                        {iPart.slice(1, -1)}
                      </em>
                    );
                  }
                  return iPart;
                })}
              </span>
            );
          })}
        </span>
      );
    });
  };

  // Helper to format a block of text containing paragraphs, lists, quotes, headings
  const renderTextBlock = (text: string, blockKey: number | string) => {
    const lines = text.split('\n');
    const elements: React.ReactNode[] = [];
    let currentList: { type: 'ul' | 'ol'; items: string[] } | null = null;

    const flushList = () => {
      if (!currentList) return;
      const listKey = `list-${elements.length}`;
      if (currentList.type === 'ul') {
        elements.push(
          <ul key={listKey} className="my-2 space-y-1.5 pl-1">
            {currentList.items.map((item, idx) => (
              <li key={idx} className="flex items-start gap-2.5 text-xs text-[#E2E8F0] leading-relaxed">
                <span className="w-1.5 h-1.5 rounded-full bg-[#EF4444] mt-2 shrink-0" />
                <span className="flex-1">{renderInline(item)}</span>
              </li>
            ))}
          </ul>
        );
      } else {
        elements.push(
          <ol key={listKey} className="my-2 space-y-1.5 pl-1">
            {currentList.items.map((item, idx) => (
              <li key={idx} className="flex items-start gap-2 text-xs text-[#E2E8F0] leading-relaxed">
                <span className="px-1.5 py-0.2 rounded bg-[#1F1F30] border border-[#2E2E48] font-mono text-[10px] text-[#DC2626] shrink-0 mt-0.5">
                  {idx + 1}
                </span>
                <span className="flex-1">{renderInline(item)}</span>
              </li>
            ))}
          </ol>
        );
      }
      currentList = null;
    };

    lines.forEach((rawLine, lIdx) => {
      const line = rawLine.trimEnd();
      const trimmed = line.trim();

      if (!trimmed) {
        flushList();
        return;
      }

      // Heading 3: ###
      if (trimmed.startsWith('### ')) {
        flushList();
        elements.push(
          <h4 key={`h3-${lIdx}`} className="text-sm font-bold text-[#0F172A] mt-3.5 mb-1.5 flex items-center gap-1.5">
            <span className="w-1.5 h-3.5 bg-[#EF4444] rounded-sm" />
            {renderInline(trimmed.slice(4))}
          </h4>
        );
        return;
      }

      // Heading 2: ##
      if (trimmed.startsWith('## ')) {
        flushList();
        elements.push(
          <h3 key={`h2-${lIdx}`} className="text-[15px] font-bold text-white mt-4 mb-2 pb-1 border-b border-[#E2E8F0]">
            {renderInline(trimmed.slice(3))}
          </h3>
        );
        return;
      }

      // Heading 1: #
      if (trimmed.startsWith('# ')) {
        flushList();
        elements.push(
          <h2 key={`h1-${lIdx}`} className="text-base font-extrabold text-white mt-4 mb-2 pb-1 border-b border-[#E2E8F0]">
            {renderInline(trimmed.slice(2))}
          </h2>
        );
        return;
      }

      // Blockquote: >
      if (trimmed.startsWith('>')) {
        flushList();
        elements.push(
          <blockquote
            key={`quote-${lIdx}`}
            className="my-2 pl-3 py-1.5 border-l-2 border-[#EF4444] bg-[#12121C] rounded-r-lg text-xs text-[#CBD5E1] italic"
          >
            {renderInline(trimmed.replace(/^>\s*/, ''))}
          </blockquote>
        );
        return;
      }

      // Unordered list: - or *
      const ulMatch = trimmed.match(/^[-*]\s+(.*)$/);
      if (ulMatch) {
        if (!currentList || currentList.type !== 'ul') {
          flushList();
          currentList = { type: 'ul', items: [] };
        }
        currentList.items.push(ulMatch[1]);
        return;
      }

      // Ordered list: 1. or 2.
      const olMatch = trimmed.match(/^(\d+)\.\s+(.*)$/);
      if (olMatch) {
        if (!currentList || currentList.type !== 'ol') {
          flushList();
          currentList = { type: 'ol', items: [] };
        }
        currentList.items.push(olMatch[2]);
        return;
      }

      // Normal paragraph line
      flushList();
      elements.push(
        <p key={`p-${lIdx}`} className="text-xs text-[#E2E8F0] leading-relaxed my-1">
          {renderInline(trimmed)}
        </p>
      );
    });

    flushList();

    return (
      <div key={blockKey} className="space-y-1">
        {elements}
      </div>
    );
  };

  return (
    <div className={`space-y-3 font-sans ${className}`}>
      {parts.map((part, idx) => {
        if (part.type === 'code') {
          const isCopied = copiedIndex === idx;
          return (
            <div
              key={idx}
              className="my-3 rounded-xl overflow-hidden border border-[#232338] bg-[#0A0A10] shadow-sm font-mono text-xs"
            >
              {/* Code Header Bar */}
              <div className="flex items-center justify-between px-3.5 py-1.5 bg-[#12121D] border-b border-[#232338] text-[11px] text-[#64748B]">
                <div className="flex items-center gap-2">
                  <Terminal className="w-3.5 h-3.5 text-[#EF4444]" />
                  <span className="font-semibold uppercase tracking-wider text-[10px] text-[#DC2626]">
                    {part.language || 'code'}
                  </span>
                </div>
                <button
                  type="button"
                  onClick={() => handleCopy(part.content, idx)}
                  className="flex items-center gap-1.5 px-2 py-0.5 rounded-md hover:bg-[#1C1C2C] text-[#64748B] hover:text-[#0F172A] transition-colors cursor-pointer"
                >
                  {isCopied ? (
                    <>
                      <Check className="w-3 h-3 text-emerald-400" />
                      <span className="text-[10px] text-emerald-400 font-medium">Copied!</span>
                    </>
                  ) : (
                    <>
                      <Copy className="w-3 h-3" />
                      <span className="text-[10px]">Copy</span>
                    </>
                  )}
                </button>
              </div>

              {/* Code Content */}
              <div className="p-3.5 overflow-x-auto text-[#E2E8F0] leading-relaxed font-mono">
                <pre>{part.content}</pre>
              </div>
            </div>
          );
        }

        return renderTextBlock(part.content, idx);
      })}
    </div>
  );
};

export default FormattedMarkdown;
