export interface ConceptMemory {
  id: string;
  student_id: string;
  topic_id?: string;
  subtopic_id?: string;
  concept_key: string;
  concept_name: string;
  topic_name?: string;
  learning_state: 'discovered' | 'practicing' | 'reinforcing' | 'mastered' | 'struggling';
  estimated_mastery: number;
  confidence_level: 'low' | 'moderate' | 'high';
  total_attempts: number;
  independent_solves: number;
  assisted_solves: number;
  failed_attempts: number;
  hint_requests: number;
  mistake_count: number;
  corrected_mistakes: number;
  first_practiced_at: string;
  last_practiced_at: string;
  last_successful_at?: string;
  next_review_at?: string;
}

export interface MistakeMemoryItem {
  id: string;
  student_id: string;
  category: string;
  subtype: string;
  concept_key: string;
  occurrence_count: number;
  corrected_count: number;
  first_occurrence_at: string;
  last_occurrence_at: string;
  last_corrected_at?: string;
  status: 'active' | 'improving' | 'resolved';
  supporting_diagnosis_ids: string[];
}

export interface ConceptStrength {
  concept_key: string;
  concept_name: string;
  topic_name: string;
  independent_solves: number;
  estimated_mastery: number;
  confidence_level: string;
  reason: string;
}

export interface ConceptWeakness {
  concept_key: string;
  concept_name: string;
  topic_name: string;
  mistake_category: string;
  mistake_subtype: string;
  occurrence_count: number;
  corrected_count: number;
  status: string;
  reason: string;
}

export interface LearningEventItem {
  id: string;
  student_id: string;
  event_type: string;
  problem_id?: string;
  submission_id?: string;
  diagnosis_id?: string;
  topic_id?: string;
  subtopic_id?: string;
  concept_key?: string;
  event_data: Record<string, any>;
  occurred_at: string;
}

export interface TimelinePoint {
  date: string;
  independent_solves: number;
  assisted_solves: number;
  mistakes_detected: number;
  mistakes_corrected: number;
  concepts_practised: number;
}

export interface LearningMemorySummary {
  total_concepts_practised: number;
  independent_solves: number;
  assisted_solves: number;
  total_mistakes: number;
  corrected_mistakes: number;
  active_learning_days: number;
  mastered_concepts: number;
  struggling_concepts: number;
  recent_concepts: ConceptMemory[];
  active_weaknesses: MistakeMemoryItem[];
  key_strengths: ConceptStrength[];
}

export interface StudentFacingNarrative {
  headline: string;
  narrative: string;
  top_strengths: string[];
  areas_for_focus: string[];
  next_action_prompt: string;
}

export interface ConceptDetailResponse {
  concept: ConceptMemory;
  mistakes: MistakeMemoryItem[];
  timeline: LearningEventItem[];
  next_steps: string[];
}

export interface LearningSnapshot {
  id: string;
  student_id: string;
  snapshot_date: string;
  total_concepts_practised: number;
  independent_solves: number;
  assisted_solves: number;
  mistake_count: number;
  corrected_mistake_count: number;
  active_learning_days: number;
  summary?: Record<string, any>;
  created_at?: string;
}
