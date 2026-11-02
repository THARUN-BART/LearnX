const ALIAS: Record<string, string> = {
  python: 'python',
  python3: 'python',
  py: 'python',
  cpp: 'cpp',
  'c++': 'cpp',
  cplusplus: 'cpp',
  c: 'c',
  java: 'java',
  go: 'go',
  golang: 'go',
};

export function toJudgeLanguage(lang: string | undefined): string {
  const key = (lang || 'python').trim().toLowerCase();
  return ALIAS[key] || 'python';
}

export function toMonacoLanguage(lang: string | undefined): string {
  const mapped = toJudgeLanguage(lang);
  return mapped === 'cpp' ? 'cpp' : mapped;
}

export function starterCodeKey(lang: string | undefined): string {
  return toJudgeLanguage(lang);
}
