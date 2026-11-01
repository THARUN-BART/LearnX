import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Sparkles,
  ArrowRight,
  Code2,
  Cpu,
  Brain,
  Terminal,
  Zap,
  CheckCircle2,
  TrendingUp,
  ShieldAlert,
  Layers,
  Flame,
  Award,
  GitBranch,
  Bot,
  Activity,
  Check,
  X,
  Compass,
  BarChart3,
  Server,
  Workflow,
  Search,
  ChevronRight
} from 'lucide-react';
import Button from '../../components/ui/Button';

export const LandingPage: React.FC = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'editor' | 'ai' | 'roadmap' | 'stats'>('editor');
  const [activeHowStep, setActiveHowStep] = useState<number>(1);

  return (
    <div className="relative overflow-hidden bg-[#07070B] text-[#F4F4F6]">
      {/* Background technical grid */}
      <div className="absolute inset-0 bg-tech-grid opacity-25 pointer-events-none" />

      {/* Ambient glowing radial lights */}
      <div className="absolute top-20 left-1/2 -translate-x-1/2 w-[650px] h-[350px] bg-gradient-to-b from-[#8B5CF6]/15 via-[#6366F1]/5 to-transparent blur-[140px] pointer-events-none -z-10" />

      {/* HERO SECTION */}
      <section className="pt-24 pb-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto text-center relative z-10">
        <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-[#12121B] border border-[#1C1C28] text-xs font-mono text-[#A78BFA] mb-8 shadow-inner">
          <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
          <span>LEARNX V2.4 • COGNITIVE DSA COMPILER</span>
        </div>

        <h1 className="text-4xl sm:text-6xl lg:text-7xl font-bold tracking-tight text-[#F4F4F6] max-w-4xl mx-auto leading-[1.08] font-mono">
          Don't just learn to code.{' '}
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#A78BFA] via-[#C084FC] to-[#818CF8]">
            Learn to think like a developer.
          </span>
        </h1>

        <p className="mt-6 text-base sm:text-lg text-[#8E8EA0] max-w-2xl mx-auto leading-relaxed">
          An intelligent learning environment that adapts to your skills, analyzes your mistakes and helps you become a better problem solver.
        </p>

        <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button
            size="lg"
            variant="primary"
            onClick={() => navigate('/signup')}
            rightIcon={<ArrowRight className="w-5 h-5" />}
            className="w-full sm:w-auto px-8 font-mono text-sm"
          >
            Start Learning
          </Button>
          <Button
            size="lg"
            variant="secondary"
            onClick={() => {
              const el = document.getElementById('how-it-works');
              el?.scrollIntoView({ behavior: 'smooth' });
            }}
            className="w-full sm:w-auto px-8 font-mono text-sm"
          >
            Explore LearnX
          </Button>
        </div>

        {/* INTERACTIVE PRODUCT VISUALIZATION */}
        <div className="mt-16 sm:mt-24 max-w-5xl mx-auto text-left">
          {/* Preview Navigation Tabs */}
          <div className="flex items-center justify-between bg-[#0D0D14] border border-[#1C1C28] rounded-t-2xl px-4 py-3">
            <div className="flex items-center gap-2">
              <div className="w-2.5 h-2.5 rounded-full bg-red-500/80" />
              <div className="w-2.5 h-2.5 rounded-full bg-amber-500/80" />
              <div className="w-2.5 h-2.5 rounded-full bg-emerald-500/80" />
              <span className="ml-2 text-xs font-mono text-[#5A5A6E]">learnx-interactive-cockpit</span>
            </div>

            <div className="flex items-center gap-1 bg-[#12121B] p-1 rounded-xl border border-[#1C1C28] text-xs">
              <button
                onClick={() => setActiveTab('editor')}
                className={`px-3 py-1 rounded-lg font-mono transition-colors cursor-pointer ${
                  activeTab === 'editor' ? 'bg-[#8B5CF6] text-white' : 'text-[#8E8EA0] hover:text-white'
                }`}
              >
                Monaco IDE
              </button>
              <button
                onClick={() => setActiveTab('ai')}
                className={`px-3 py-1 rounded-lg font-mono transition-colors cursor-pointer ${
                  activeTab === 'ai' ? 'bg-[#8B5CF6] text-white' : 'text-[#8E8EA0] hover:text-white'
                }`}
              >
                AI Mentor
              </button>
              <button
                onClick={() => setActiveTab('roadmap')}
                className={`px-3 py-1 rounded-lg font-mono transition-colors cursor-pointer ${
                  activeTab === 'roadmap' ? 'bg-[#8B5CF6] text-white' : 'text-[#8E8EA0] hover:text-white'
                }`}
              >
                Curriculum
              </button>
              <button
                onClick={() => setActiveTab('stats')}
                className={`px-3 py-1 rounded-lg font-mono transition-colors cursor-pointer ${
                  activeTab === 'stats' ? 'bg-[#8B5CF6] text-white' : 'text-[#8E8EA0] hover:text-white'
                }`}
              >
                Telemetry
              </button>
            </div>
          </div>

          {/* Tab Screen Content */}
          <div className="bg-[#0D0D14] border-x border-b border-[#1C1C28] rounded-b-2xl p-6 shadow-2xl shadow-black min-h-[380px]">
            {activeTab === 'editor' && (
              <div className="grid grid-cols-1 md:grid-cols-12 gap-6">
                <div className="md:col-span-4 space-y-3 font-mono text-xs border-r border-[#1C1C28] pr-4">
                  <div className="text-[10px] uppercase font-mono px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 w-fit">
                    Array // Two Pointers
                  </div>
                  <h3 className="text-base font-bold text-[#F4F4F6]">3Sum Triplet Target</h3>
                  <p className="text-[#8E8EA0] leading-relaxed font-sans">
                    Given an integer array nums, return all triplets [nums[i], nums[j], nums[k]] such that i != j != k and their sum is 0.
                  </p>
                  <div className="p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-1">
                    <span className="text-[#5A5A6E]">Runtime Constraints:</span>
                    <div className="text-emerald-400">3 &lt;= nums.length &lt;= 3000</div>
                    <div className="text-[#818CF8]">Target Complexity: O(N^2)</div>
                  </div>
                </div>

                <div className="md:col-span-8 flex flex-col justify-between font-mono text-xs bg-[#09090E] p-4 rounded-xl border border-[#1C1C28]">
                  <pre className="text-emerald-300 leading-relaxed overflow-x-auto">
{`class Solution:
    def threeSum(self, nums: list[int]) -> list[list[int]]:
        nums.sort()
        res = []
        for i, a in enumerate(nums):
            if i > 0 and a == nums[i - 1]:
                continue
            l, r = i + 1, len(nums) - 1
            while l < r:
                three_sum = a + nums[l] + nums[r]
                if three_sum > 0:
                    r -= 1
                elif three_sum < 0:
                    l += 1
                else:
                    res.append([a, nums[l], nums[r]])
                    l += 1
                    while nums[l] == nums[l - 1] and l < r:
                        l += 1
        return res`}
                  </pre>
                  <div className="mt-4 pt-3 border-t border-[#1C1C28] flex items-center justify-between text-[#8E8EA0]">
                    <span className="text-emerald-400 font-semibold flex items-center gap-1.5">
                      <CheckCircle2 className="w-4 h-4" /> 32 / 32 Test Cases Verified (0.04s)
                    </span>
                    <span className="text-[11px] text-[#5A5A6E]">Memory: 18.2 MB (Top 96%)</span>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'ai' && (
              <div className="space-y-4 max-w-2xl mx-auto font-mono">
                <div className="flex items-start gap-3">
                  <div className="w-8 h-8 rounded-xl bg-[#8B5CF6]/20 border border-[#8B5CF6]/40 flex items-center justify-center text-[#A78BFA] shrink-0">
                    <Bot className="w-4 h-4" />
                  </div>
                  <div className="bg-[#12121B] border border-[#1C1C28] p-4 rounded-2xl text-xs space-y-2">
                    <div className="font-semibold text-[#A78BFA] flex items-center gap-2">
                      <span>LearnX Adaptive Diagnostic</span>
                      <span className="text-[10px] text-[#5A5A6E]">Hint Mode</span>
                    </div>
                    <p className="text-[#F4F4F6] leading-relaxed font-sans">
                      Sorting the array takes $O(N \log N)$. Once you fix the first integer, the subproblem reduces to a standard <strong>Sorted Two Sum</strong>.
                    </p>
                    <p className="text-[#8E8EA0] font-sans">
                      To avoid duplicate triplets in $O(1)$ space, increment the left pointer and decrement the right pointer while adjacent values match.
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-2 justify-end">
                  <div className="bg-[#8B5CF6]/20 border border-[#8B5CF6]/40 text-[#F4F4F6] px-4 py-2 rounded-xl text-xs font-mono">
                    Does advancing past duplicates alter the linear time traversal?
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'roadmap' && (
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 font-mono">
                <div className="p-4 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-2">
                  <span className="text-[10px] px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                    Completed
                  </span>
                  <h4 className="font-bold text-sm text-[#F4F4F6]">Arrays & Two Pointers</h4>
                  <p className="text-xs text-[#8E8EA0]">18 / 18 problems solved</p>
                  <div className="w-full bg-[#181824] h-1.5 rounded-full overflow-hidden">
                    <div className="bg-emerald-500 h-full w-full" />
                  </div>
                </div>

                <div className="p-4 bg-[#12121B] rounded-xl border border-[#8B5CF6]/40 glow-purple space-y-2">
                  <span className="text-[10px] px-2 py-0.5 rounded bg-[#8B5CF6]/20 text-[#A78BFA] border border-[#8B5CF6]/30">
                    Active Pillar
                  </span>
                  <h4 className="font-bold text-sm text-[#F4F4F6]">Binary Search & Answers</h4>
                  <p className="text-xs text-[#8E8EA0]">14 / 22 problems solved (63%)</p>
                  <div className="w-full bg-[#181824] h-1.5 rounded-full overflow-hidden">
                    <div className="bg-[#8B5CF6] h-full w-[63%]" />
                  </div>
                </div>

                <div className="p-4 bg-[#12121B] rounded-xl border border-[#1C1C28] opacity-60 space-y-2">
                  <span className="text-[10px] px-2 py-0.5 rounded bg-[#181824] text-[#8E8EA0] border border-[#2B2B3C]">
                    Locked
                  </span>
                  <h4 className="font-bold text-sm text-[#F4F4F6]">Dynamic Programming</h4>
                  <p className="text-xs text-[#8E8EA0]">Requires tree recursion mastery</p>
                  <div className="w-full bg-[#181824] h-1.5 rounded-full overflow-hidden">
                    <div className="bg-[#2B2B3C] h-full w-0" />
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'stats' && (
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 text-center font-mono">
                <div className="p-5 bg-[#12121B] rounded-xl border border-[#1C1C28]">
                  <div className="text-3xl font-bold text-emerald-400">79.4%</div>
                  <div className="text-xs text-[#8E8EA0] mt-1">First-Try Pass Velocity</div>
                </div>
                <div className="p-5 bg-[#12121B] rounded-xl border border-[#1C1C28]">
                  <div className="text-3xl font-bold text-[#A78BFA]">14 Days</div>
                  <div className="text-xs text-[#8E8EA0] mt-1">Unbroken Solving Consistency</div>
                </div>
                <div className="p-5 bg-[#12121B] rounded-xl border border-[#1C1C28]">
                  <div className="text-2xl font-bold text-amber-400">Edge Cases</div>
                  <div className="text-xs text-[#8E8EA0] mt-1">Top Addressed Anti-Pattern</div>
                </div>
              </div>
            )}
          </div>
        </div>
      </section>

      {/* SECTION 1: THE LEARNING PROBLEM */}
      <section id="problem" className="py-24 border-t border-[#1C1C28] px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto scroll-mt-20">
        <div className="max-w-2xl mx-auto text-center space-y-3 mb-16">
          <span className="tech-pill text-red-400 border-red-500/20">The Core Bottleneck</span>
          <h2 className="text-3xl sm:text-4xl font-bold text-[#F4F4F6] font-mono">
            Traditional coding platforms measure answers, not understanding.
          </h2>
          <p className="text-sm text-[#8E8EA0]">
            Students grind hundreds of questions by memorizing boilerplates, only to freeze when interviewers tweak a single constraint.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 font-mono text-xs">
          <div className="p-6 rounded-2xl bg-[#0D0D14] border border-red-500/20 space-y-3">
            <ShieldAlert className="w-6 h-6 text-red-400" />
            <h3 className="font-bold text-sm text-[#F4F4F6]">The Solution Leaking Trap</h3>
            <p className="text-[#8E8EA0] leading-relaxed font-sans">
              Opening the discussion tab tricks the brain into passive comprehension without reinforcing actual neural pattern recall.
            </p>
          </div>

          <div className="p-6 rounded-2xl bg-[#0D0D14] border border-amber-500/20 space-y-3">
            <Layers className="w-6 h-6 text-amber-400" />
            <h3 className="font-bold text-sm text-[#F4F4F6]">Rigid, Non-Adaptive Curricula</h3>
            <p className="text-[#8E8EA0] leading-relaxed font-sans">
              Every learner gets an identical static list regardless of whether their weakness is recursion, pointer mechanics, or time complexity analysis.
            </p>
          </div>

          <div className="p-6 rounded-2xl bg-[#0D0D14] border border-[#8B5CF6]/20 space-y-3">
            <Brain className="w-6 h-6 text-[#A78BFA]" />
            <h3 className="font-bold text-sm text-[#F4F4F6]">Zero Error Forensics</h3>
            <p className="text-[#8E8EA0] leading-relaxed font-sans">
              Platforms spit out "Wrong Answer on testcase 48/52" and abandon you without explaining why your mental model diverged.
            </p>
          </div>
        </div>
      </section>

      {/* SECTION 2: OUR APPROACH */}
      <section id="approach" className="py-24 border-t border-[#1C1C28] bg-[#0A0A10] px-4 sm:px-6 lg:px-8 scroll-mt-20">
        <div className="max-w-7xl mx-auto space-y-16">
          <div className="max-w-3xl mx-auto text-center space-y-3">
            <span className="tech-pill text-[#A78BFA] border-[#8B5CF6]/30">The LearnX Approach</span>
            <h2 className="text-3xl sm:text-4xl font-bold font-mono text-[#F4F4F6]">
              A Cognitive Feedback Loop Designed for Technical Precision
            </h2>
            <p className="text-sm text-[#8E8EA0] font-sans">
              Rather than spoon-feeding solutions, LearnX operates like a senior principal engineer pairing with you through progressive cues, telemetry diagnostics, and mistake prevention.
            </p>
          </div>

          {/* Architecture Flow Diagram */}
          <div className="p-6 sm:p-8 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] relative">
            <div className="grid grid-cols-1 lg:grid-cols-4 gap-6 font-mono text-xs">
              {/* Step 1 */}
              <div className="p-5 rounded-xl bg-[#12121B] border border-[#1C1C28] space-y-2.5 relative">
                <div className="w-8 h-8 rounded-lg bg-[#8B5CF6]/20 text-[#A78BFA] flex items-center justify-center font-bold">
                  01
                </div>
                <h4 className="text-sm font-bold text-[#F4F4F6]">Static AST & Logic Parsing</h4>
                <p className="text-[#8E8EA0] font-sans text-xs">
                  Inspects loop invariants, nested iterations, recursive base bounds, and memory allocations in real-time.
                </p>
              </div>

              {/* Step 2 */}
              <div className="p-5 rounded-xl bg-[#12121B] border border-[#1C1C28] space-y-2.5">
                <div className="w-8 h-8 rounded-lg bg-[#6366F1]/20 text-[#818CF8] flex items-center justify-center font-bold">
                  02
                </div>
                <h4 className="text-sm font-bold text-[#F4F4F6]">Cognitive Gap Detection</h4>
                <p className="text-[#8E8EA0] font-sans text-xs">
                  Classifies whether a failure was caused by asymptotic time complexity, pointer underflow, or missing edge cases.
                </p>
              </div>

              {/* Step 3 */}
              <div className="p-5 rounded-xl bg-[#12121B] border border-[#1C1C28] space-y-2.5">
                <div className="w-8 h-8 rounded-lg bg-amber-500/20 text-amber-400 flex items-center justify-center font-bold">
                  03
                </div>
                <h4 className="text-sm font-bold text-[#F4F4F6]">Socratic Hint Calibration</h4>
                <p className="text-[#8E8EA0] font-sans text-xs">
                  Generates progressive prompts (Hint, Explain, Debug) that guide you toward the invariant without spoiling the answer.
                </p>
              </div>

              {/* Step 4 */}
              <div className="p-5 rounded-xl bg-[#12121B] border border-[#1C1C28] space-y-2.5">
                <div className="w-8 h-8 rounded-lg bg-emerald-500/20 text-emerald-400 flex items-center justify-center font-bold">
                  04
                </div>
                <h4 className="text-sm font-bold text-[#F4F4F6]">Curriculum Weighting</h4>
                <p className="text-[#8E8EA0] font-sans text-xs">
                  Updates your student profile in Supabase PostgreSQL, queuing targeted drill problems to eradicate the mistake.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* SECTION 3: HOW IT WORKS */}
      <section id="how-it-works" className="py-24 border-t border-[#1C1C28] px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto scroll-mt-20">
        <div className="max-w-2xl mx-auto text-center space-y-3 mb-16">
          <span className="tech-pill text-cyan-400 border-cyan-500/20">Operational Workflow</span>
          <h2 className="text-3xl sm:text-4xl font-bold text-[#F4F4F6] font-mono">
            How LearnX Transforms Your Coding Velocity
          </h2>
          <p className="text-sm text-[#8E8EA0]">
            From university placement exams to FAANG technical screening rounds, follow a disciplined three-stage progression.
          </p>
        </div>

        {/* Step Selector & Details */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center font-mono">
          <div className="lg:col-span-5 space-y-3">
            {[
              {
                step: 1,
                title: 'Phase 1: Academic & Target Calibration',
                desc: 'Set your university, course, preferred language (Python, Go, Java, C++, C), and DSA proficiency tier.'
              },
              {
                step: 2,
                title: 'Phase 2: 16-Pillar Milestone Progression',
                desc: 'Tackle the comprehensive roadmap spanning Programming Fundamentals, Arrays, Recursion, Trees, and Dynamic Programming.'
              },
              {
                step: 3,
                title: 'Phase 3: Automated Anti-Pattern Forensics',
                desc: 'Inspect classified mistake retrospectives with line-by-line diffs to build interview intuition that lasts.'
              }
            ].map((s) => (
              <div
                key={s.step}
                onClick={() => setActiveHowStep(s.step)}
                className={`p-5 rounded-2xl border cursor-pointer transition-all ${
                  activeHowStep === s.step
                    ? 'bg-[#12121B] border-[#8B5CF6]/50 shadow-lg shadow-[#8B5CF6]/10'
                    : 'bg-[#0D0D14] border-[#1C1C28] hover:border-[#2B2B3C]'
                }`}
              >
                <div className="flex items-center gap-3">
                  <span
                    className={`w-7 h-7 rounded-lg flex items-center justify-center text-xs font-bold ${
                      activeHowStep === s.step
                        ? 'bg-[#8B5CF6] text-white'
                        : 'bg-[#181824] text-[#8E8EA0]'
                    }`}
                  >
                    0{s.step}
                  </span>
                  <h4 className="text-sm font-bold text-[#F4F4F6]">{s.title}</h4>
                </div>
                <p className="text-xs text-[#8E8EA0] mt-2 font-sans pl-10 leading-relaxed">
                  {s.desc}
                </p>
              </div>
            ))}
          </div>

          <div className="lg:col-span-7 p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
            <div className="flex items-center justify-between border-b border-[#1C1C28] pb-3 text-xs text-[#5A5A6E]">
              <div className="flex items-center gap-2">
                <Workflow className="w-4 h-4 text-[#8B5CF6]" />
                <span>learnx-pipeline-stage-0{activeHowStep}.json</span>
              </div>
              <span className="text-emerald-400">ACTIVE PIPELINE</span>
            </div>

            {activeHowStep === 1 && (
              <div className="space-y-3 text-xs">
                <div className="p-3.5 bg-[#12121B] rounded-xl border border-[#1C1C28]">
                  <span className="text-[#A78BFA] block font-bold mb-1">// Calibration Schema</span>
                  <div className="text-[#8E8EA0]">Current Academic Year: <strong className="text-white">3rd Year</strong></div>
                  <div className="text-[#8E8EA0]">Preferred Language: <strong className="text-white">Python / Go</strong></div>
                  <div className="text-[#8E8EA0]">DSA Baseline Tier: <strong className="text-white">Intermediate</strong></div>
                </div>
                <p className="text-xs text-[#8E8EA0] font-sans leading-relaxed">
                  LearnX calibrates problem recommendations against your placement deadline, ensuring you avoid wasting time on problems below your level.
                </p>
              </div>
            )}

            {activeHowStep === 2 && (
              <div className="space-y-3 text-xs">
                <div className="p-3.5 bg-[#12121B] rounded-xl border border-[#1C1C28]">
                  <span className="text-emerald-400 block font-bold mb-1">// Active Milestone</span>
                  <div className="text-[#8E8EA0]">Pillar: <strong className="text-white">Searching & Sorting // Binary Search</strong></div>
                  <div className="text-[#8E8EA0]">Core Concept: <strong className="text-white">Search on Monotonic Answer Spaces</strong></div>
                  <div className="text-[#8E8EA0]">Verified Submissions: <strong className="text-emerald-400">14 / 22 Passed</strong></div>
                </div>
                <p className="text-xs text-[#8E8EA0] font-sans leading-relaxed">
                  Work directly in our embedded Monaco editor. Run test cases against our backend execution sandbox with instant stdout/stderr telemetry.
                </p>
              </div>
            )}

            {activeHowStep === 3 && (
              <div className="space-y-3 text-xs">
                <div className="p-3.5 bg-[#12121B] rounded-xl border border-red-500/20">
                  <span className="text-red-400 block font-bold mb-1">// Mistake Forensic Ledger</span>
                  <div className="text-[#8E8EA0]">Anti-Pattern: <strong className="text-amber-400">Off-by-one upper bound slip</strong></div>
                  <div className="text-[#8E8EA0]">Failed Testcase: <strong className="text-white">nums=[3,3], target=6</strong></div>
                  <div className="text-[#8E8EA0]">Recovery Drill: <strong className="text-emerald-400">Two Sum Frequency Map</strong></div>
                </div>
                <p className="text-xs text-[#8E8EA0] font-sans leading-relaxed">
                  Mistakes are not dead-ends. LearnX flags identical failure modes across your submissions and prescribes targeted re-practice drills.
                </p>
              </div>
            )}
          </div>
        </div>
      </section>

      {/* SECTION 4: PRODUCT SHOWCASE */}
      <section id="showcase" className="py-24 border-t border-[#1C1C28] bg-[#0A0A10] px-4 sm:px-6 lg:px-8 scroll-mt-20">
        <div className="max-w-7xl mx-auto space-y-16">
          <div className="max-w-3xl mx-auto text-center space-y-3">
            <span className="tech-pill text-[#818CF8] border-[#6366F1]/30">Product Showcase</span>
            <h2 className="text-3xl sm:text-4xl font-bold font-mono text-[#F4F4F6]">
              Engineered Like a High-End Developer Tool
            </h2>
            <p className="text-sm text-[#8E8EA0] font-sans">
              Built on battle-tested developer primitives: Monaco editor, Raycast-style command palettes, Go REST execution kernels, and PostgreSQL session stores.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 font-mono text-xs">
            <div className="p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-3 hover:border-[#8B5CF6]/50 transition-all">
              <Code2 className="w-6 h-6 text-[#8B5CF6]" />
              <h4 className="font-bold text-sm text-[#F4F4F6]">Monaco IDE Engine</h4>
              <p className="text-[#8E8EA0] font-sans leading-relaxed">
                The identical editing core that powers VS Code. Full syntax highlighting, keybindings, and multi-language compilation (C++, Java, Python, C, Go).
              </p>
            </div>

            <div className="p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-3 hover:border-[#6366F1]/50 transition-all">
              <Bot className="w-6 h-6 text-[#6366F1]" />
              <h4 className="font-bold text-sm text-[#F4F4F6]">Socratic Copilot</h4>
              <p className="text-[#8E8EA0] font-sans leading-relaxed">
                Choose Hint Mode, Explain Mode, or Debug Mode depending on how deep you want guidance. Features one-click code copy and prompt history.
              </p>
            </div>

            <div className="p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-3 hover:border-emerald-500/50 transition-all">
              <GitBranch className="w-6 h-6 text-emerald-400" />
              <h4 className="font-bold text-sm text-[#F4F4F6]">16-Pillar Curriculum</h4>
              <p className="text-[#8E8EA0] font-sans leading-relaxed">
                Structured progression from fundamental memory models up to advanced backtracking, graph shortest-paths, and 2D dynamic programming.
              </p>
            </div>

            <div className="p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-3 hover:border-amber-500/50 transition-all">
              <BarChart3 className="w-6 h-6 text-amber-400" />
              <h4 className="font-bold text-sm text-[#F4F4F6]">Precision Analytics</h4>
              <p className="text-[#8E8EA0] font-sans leading-relaxed">
                Quantitative tracking of solving volume, solving speed percentiles, 20-week heatmap consistency, and topic competency dials.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* BENCHMARK COMPARISON TABLE */}
      <section className="py-20 border-t border-[#1C1C28] bg-[#07070B] px-4 sm:px-6 lg:px-8">
        <div className="max-w-5xl mx-auto space-y-12">
          <div className="text-center space-y-2">
            <span className="tech-pill text-[#A78BFA] border-[#8B5CF6]/30">Architectural Difference</span>
            <h2 className="text-3xl font-bold font-mono text-[#F4F4F6]">
              How LearnX Compares
            </h2>
          </div>

          <div className="overflow-x-auto rounded-2xl border border-[#1C1C28] bg-[#0D0D14]">
            <table className="w-full text-left text-xs font-mono">
              <thead className="bg-[#12121B] border-b border-[#1C1C28] text-[#8E8EA0]">
                <tr>
                  <th className="p-4">Capability</th>
                  <th className="p-4 text-[#8B5CF6] font-bold">LearnX Adaptive Engine</th>
                  <th className="p-4">Generic Problem Platforms</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#1C1C28] text-[#D1D1DE]">
                <tr>
                  <td className="p-4 font-semibold text-[#F4F4F6]">Hint Methodology</td>
                  <td className="p-4 text-emerald-400 flex items-center gap-2">
                    <Check className="w-4 h-4" /> Socratic progressive cues without spoiling solution
                  </td>
                  <td className="p-4 text-[#5A5A6E] flex items-center gap-2">
                    <X className="w-4 h-4 text-red-400" /> Complete solution dumps or zero hints
                  </td>
                </tr>
                <tr>
                  <td className="p-4 font-semibold text-[#F4F4F6]">Mistake Categorization</td>
                  <td className="p-4 text-emerald-400 flex items-center gap-2">
                    <Check className="w-4 h-4" /> Classifies edge cases, complexity slips, & logic regressions
                  </td>
                  <td className="p-4 text-[#5A5A6E] flex items-center gap-2">
                    <X className="w-4 h-4 text-red-400" /> Generic "Wrong Answer" boolean flag
                  </td>
                </tr>
                <tr>
                  <td className="p-4 font-semibold text-[#F4F4F6]">Curriculum Calibration</td>
                  <td className="p-4 text-emerald-400 flex items-center gap-2">
                    <Check className="w-4 h-4" /> Dynamic based on error frequency & academic target
                  </td>
                  <td className="p-4 text-[#5A5A6E] flex items-center gap-2">
                    <X className="w-4 h-4 text-red-400" /> Static linear question lists
                  </td>
                </tr>
                <tr>
                  <td className="p-4 font-semibold text-[#F4F4F6]">Execution Sandbox</td>
                  <td className="p-4 text-emerald-400 flex items-center gap-2">
                    <Check className="w-4 h-4" /> Go, Python, Java, C++, C with instant execution
                  </td>
                  <td className="p-4 text-[#5A5A6E] flex items-center gap-2">
                    <Check className="w-4 h-4 text-emerald-400" /> Standard web compilers
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      {/* FINAL CTA */}
      <section className="py-24 border-t border-[#1C1C28] px-4 sm:px-6 lg:px-8 text-center relative">
        <div className="max-w-3xl mx-auto space-y-6">
          <h2 className="text-3xl sm:text-5xl font-bold tracking-tight text-[#F4F4F6] font-mono">
            Elevate your algorithmic problem solving.
          </h2>
          <p className="text-sm text-[#8E8EA0]">
            Join aspiring engineers mastering data structures and technical interviews with LearnX.
          </p>
          <div className="pt-4 flex flex-col sm:flex-row items-center justify-center gap-4">
            <Button
              size="lg"
              variant="primary"
              onClick={() => navigate('/signup')}
              rightIcon={<ArrowRight className="w-5 h-5" />}
              className="w-full sm:w-auto px-8 font-mono text-sm"
            >
              Start Free Today
            </Button>
            <Button
              size="lg"
              variant="secondary"
              onClick={() => navigate('/login')}
              className="w-full sm:w-auto px-8 font-mono text-sm"
            >
              Sign In to Account
            </Button>
          </div>
        </div>
      </section>

      {/* FOOTER */}
      <footer className="border-t border-[#1C1C28] bg-[#0A0A10] py-10 px-4 sm:px-6 lg:px-8 text-xs font-mono text-[#5A5A6E]">
        <div className="max-w-7xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 rounded-lg bg-gradient-to-tr from-[#8B5CF6] to-[#6366F1] flex items-center justify-center">
              <Sparkles className="w-3 h-3 text-white" />
            </div>
            <span className="font-bold text-sm text-[#F4F4F6]">LearnX</span>
            <span className="text-[#5A5A6E]">| Adaptive Developer Tool</span>
          </div>

          <div className="flex items-center gap-6 text-[#8E8EA0]">
            <a href="#problem" className="hover:text-white transition-colors">Problem</a>
            <a href="#approach" className="hover:text-white transition-colors">Our Approach</a>
            <a href="#how-it-works" className="hover:text-white transition-colors">How It Works</a>
            <a href="#showcase" className="hover:text-white transition-colors">Product Showcase</a>
          </div>

          <div>
            © {new Date().getFullYear()} LearnX Inc. Built for high-velocity software engineers.
          </div>
        </div>
      </footer>
    </div>
  );
};

export default LandingPage;
