/**
 * LearnX Profile Cache
 *
 * Two-tier cache for the student profile:
 *  - Tier 1: In-memory (per session, fastest reads, no serialization cost)
 *  - Tier 2: localStorage (survives page refresh, cold-start restore)
 *
 * Cache entries are invalidated after TTL_MS milliseconds. Any write to the
 * profile (update, onboarding submit) must call `profileCache.set(profile)`
 * so both tiers stay in sync.
 */

import { StudentProfile } from '../types/profile.types';
import { PROFILE_KEY } from './constants';

/** How long a cached profile is considered fresh (5 minutes). */
const TTL_MS = 5 * 60 * 1000;

interface CacheEntry {
  data: StudentProfile;
  /** Epoch ms when this entry expires */
  expiresAt: number;
}

// ---------------------------------------------------------------------------
// In-memory store (module-level singleton)
// ---------------------------------------------------------------------------
let memoryCache: CacheEntry | null = null;

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
function isExpired(entry: CacheEntry): boolean {
  return Date.now() > entry.expiresAt;
}

function makeEntry(data: StudentProfile): CacheEntry {
  return { data, expiresAt: Date.now() + TTL_MS };
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

/** Return the cached profile if it exists and has not expired. */
function get(): StudentProfile | null {
  // 1. Check memory cache
  if (memoryCache && !isExpired(memoryCache)) {
    return memoryCache.data;
  }

  // 2. Fall through to localStorage
  try {
    const raw = localStorage.getItem(PROFILE_KEY);
    if (!raw) return null;

    const entry: CacheEntry = JSON.parse(raw);

    // Migrate plain objects written by older code (no expiresAt)
    if (!entry.expiresAt || isExpired(entry)) {
      // If the stored value is the raw profile object (old format), still return it
      // so the caller can decide whether to re-fetch or use stale data.
      const maybeProfile = entry as unknown as StudentProfile;
      if (maybeProfile.id && maybeProfile.email) {
        return maybeProfile; // stale but usable for initial render
      }
      return null;
    }

    // Warm the memory cache from localStorage so next call is instant
    memoryCache = entry;
    return entry.data;
  } catch {
    return null;
  }
}

/** Store a profile in both tiers and reset the TTL. */
function set(profile: StudentProfile): void {
  const entry = makeEntry(profile);
  memoryCache = entry;

  try {
    localStorage.setItem(PROFILE_KEY, JSON.stringify(entry));
  } catch (e) {
    // localStorage quota exceeded – silently degrade to memory-only
    console.warn('[LearnX Cache] localStorage write failed:', e);
  }
}

/** Merge a partial update into the cached profile (both tiers). */
function patch(updates: Partial<StudentProfile>): StudentProfile | null {
  const current = get();
  if (!current) return null;
  const merged = { ...current, ...updates, updated_at: new Date().toISOString() };
  set(merged);
  return merged;
}

/** Remove the profile from both tiers (e.g. on logout). */
function clear(): void {
  memoryCache = null;
  try {
    localStorage.removeItem(PROFILE_KEY);
  } catch {}
}

/** Check whether a fresh (non-expired) cache entry exists. */
function isFresh(): boolean {
  if (memoryCache && !isExpired(memoryCache)) return true;
  // Check localStorage TTL without warming memory cache
  try {
    const raw = localStorage.getItem(PROFILE_KEY);
    if (!raw) return false;
    const entry: CacheEntry = JSON.parse(raw);
    return Boolean(entry.expiresAt && !isExpired(entry));
  } catch {
    return false;
  }
}

export const profileCache = { get, set, patch, clear, isFresh };
export default profileCache;

