"""Diagnosis engine for LearnX Pedagogy Service.

Architecture (three layers, cheapest first):
  Layer 1 — Traceback rule registry  (deterministic, O(1))
  Layer 2 — AST scanner              (structural, no execution)
  Layer 3 — Gemini API               (LLM, only when L1+L2 fail)
"""

import json

from google.genai import types

try:
    from person_c.config import get_gemini_client, get_model_name
    from person_c.domain_loader import load_domain_data
    from person_c.traceback_parser import parse_traceback
    from person_c.ast_scanner import scan_ast
except ImportError:  # pragma: no cover
    from config import get_gemini_client, get_model_name
    from domain_loader import load_domain_data
    from traceback_parser import parse_traceback
    from ast_scanner import scan_ast


domain_data = load_domain_data()
_client = None
_model_name = None


def _get_client():
    global _client
    if _client is None:
        _client = get_gemini_client()
    return _client


def _get_model_name():
    global _model_name
    if _model_name is None:
        _model_name = get_model_name()
    return _model_name


def _lookup_misconception(mid: str) -> dict | None:
    """Return the domain entry for *mid*, or None if not found."""
    return next(
        (item for item in domain_data["misconceptions"] if item["id"] == mid),
        None,
    )


def _build_result(mid: str, skill: str, confidence: float, evidence: str) -> dict:
    """Build a normalised diagnosis dict from a known misconception ID."""
    entry = _lookup_misconception(mid)
    return {
        "skill": skill,
        "misconception_id": mid,
        "misconception": entry["name"] if entry else mid,
        "confidence": confidence,
        "evidence": evidence,
    }


# ---------------------------------------------------------------------------
# Layer 1 — Traceback rule registry
# ---------------------------------------------------------------------------
# Each rule is a dict:
#   id       — misconception ID to return on match
#   match    — callable(parsed: dict, raw_tb: str, error_type: str) → bool
#   confidence — float 0–1
#   evidence_fn — optional callable(parsed, raw_tb) → str
# ---------------------------------------------------------------------------

def _tb(raw_tb: str, msg: str) -> bool:
    return msg in raw_tb


_TRACEBACK_RULES: list[dict] = [
    # ------------------------------------------------------------------ MIS-04
    {
        "id": "MIS-04",
        "match": lambda p, tb, et: et == "TIMEOUT",
        "confidence": 0.98,
        "evidence_fn": lambda p, tb: "Code execution timed out — likely an infinite loop.",
    },
    # ------------------------------------------------------------------ MIS-07
    {
        "id": "MIS-07",
        "match": lambda p, tb, et: (
            p["exception_type"] == "AttributeError"
            and "'int' object has no attribute" in tb
        ),
        "confidence": 0.95,
        "evidence_fn": lambda p, tb: p["exception_msg"] or tb,
    },
    # ------------------------------------------------------------------ MIS-02
    # More specific check BEFORE the generic NoneType rule so MIS-08 can
    # intercept single-node cases first.
    # ------------------------------------------------------------------ MIS-08
    {
        "id": "MIS-08",
        "match": lambda p, tb, et: (
            p["exception_type"] == "AttributeError"
            and "NoneType" in tb
            and ("head.next" in tb or ".next.next" in tb)
        ),
        "confidence": 0.88,
        "evidence_fn": lambda p, tb: p["exception_msg"] or tb,
    },
    # ------------------------------------------------------------------ MIS-02
    {
        "id": "MIS-02",
        "match": lambda p, tb, et: (
            "NoneType" in tb
            and "has no attribute" in tb
        ),
        "confidence": 0.92,
        "evidence_fn": lambda p, tb: p["exception_msg"] or tb,
    },
    # ------------------------------------------------------------------ ARR-MIS-01
    {
        "id": "ARR-MIS-01",
        "match": lambda p, tb, et: (
            p["exception_type"] == "IndexError"
            and "list index out of range" in tb
        ),
        "confidence": 0.90,
        "evidence_fn": lambda p, tb: "IndexError: list index out of range",
    },
    # ------------------------------------------------------------------ SQ-MIS-01
    {
        "id": "SQ-MIS-01",
        "match": lambda p, tb, et: (
            p["exception_type"] == "IndexError"
            and "pop from empty" in tb
        ),
        "confidence": 0.95,
        "evidence_fn": lambda p, tb: p["exception_msg"] or tb,
    },
    # ------------------------------------------------------------------ REC-MIS-01
    {
        "id": "REC-MIS-01",
        "match": lambda p, tb, et: "RecursionError" in (p["exception_type"] or ""),
        "confidence": 0.97,
        "evidence_fn": lambda p, tb: "RecursionError: maximum recursion depth exceeded",
    },
    # ------------------------------------------------------------------ BS-MIS-03
    {
        "id": "BS-MIS-03",
        "match": lambda p, tb, et: "ZeroDivisionError" in (p["exception_type"] or ""),
        "confidence": 0.92,
        "evidence_fn": lambda p, tb: p["exception_msg"] or tb,
    },
    # ------------------------------------------------------------------ SORT-MIS-01
    {
        "id": "SORT-MIS-01",
        "match": lambda p, tb, et: (
            p["exception_type"] in ("TypeError", "AttributeError")
            and "comparison" in tb
        ),
        "confidence": 0.80,
        "evidence_fn": lambda p, tb: p["exception_msg"] or tb,
    },
]


def diagnose_error(runner_result: dict, skill: str = "SKL-02") -> dict:
    """Layer-1 rule-based diagnosis from execution result only.

    Iterates the rule registry in priority order and returns the first
    matching result.  Falls back to ``MC_UNKNOWN`` (or ``None`` for a
    clean pass) when no rule fires.
    """
    error_type = runner_result.get("error_type") or ""
    raw_tb = runner_result.get("traceback") or ""
    parsed = parse_traceback(raw_tb)

    for rule in _TRACEBACK_RULES:
        try:
            if rule["match"](parsed, raw_tb, error_type):
                mid = rule["id"]
                # Only return a named misconception if it exists in the KB.
                # New IDs (ARR-MIS-01 etc.) are always valid once KB is extended.
                evidence = (
                    rule["evidence_fn"](parsed, raw_tb)
                    if "evidence_fn" in rule
                    else raw_tb or error_type
                )
                return _build_result(mid, skill, rule["confidence"], evidence)
        except Exception:  # pragma: no cover — defensive
            continue

    # No rule matched
    if error_type in ("ASSERTION_ERROR", "RUNTIME_ERROR"):
        misconception_id = "MC_UNKNOWN"
    else:
        misconception_id = None

    return {
        "skill": skill,
        "misconception_id": misconception_id,
        "misconception": "Unknown misconception",
        "confidence": 0.40,
        "evidence": raw_tb or "Incorrect output.",
    }


# ---------------------------------------------------------------------------
# Layer 3 — Gemini API fallback
# ---------------------------------------------------------------------------

def diagnose_with_gemini(student_code: str, runner_result: dict, skill: str = "SKL-02") -> dict:
    """Call the Gemini API to diagnose a submission when rules and AST fail."""
    catalogue = [
        {
            "id": item["id"],
            "name": item["name"],
            "root_cause": item["root_cause"],
            "deterministic_trigger": item["deterministic_trigger"],
        }
        for item in domain_data["misconceptions"]
    ]

    prompt = f"""
    You are an expert Python and Data Structures tutor.

    Identify the most likely misconception in the student's code.

    Student skill: {skill}

    Student code:
    {student_code}

    Execution result:
    {json.dumps(runner_result)}

    OFFICIAL MISCONCEPTION CATALOGUE:
    {json.dumps(catalogue, indent=2)}

    Instructions:
1. Analyse the student's code and execution result carefully.
2. Compare the mistake with every official misconception.
3. Select an ID ONLY when the root cause clearly matches.
4. Do not select a misconception just because it sounds similar.
5. MIS-06 is specifically about incorrect position indexing
   during insertion or deletion, NOT skipping nodes during traversal.
6. If no official misconception accurately matches, return MC_UNKNOWN.
7. Never invent an ID.
8. Do not provide the complete corrected solution.

If the student uses current.next.next during normal
linked-list traversal, check whether any official
misconception explicitly covers skipping nodes.
If none does, return MC_UNKNOWN.
    Return JSON with these fields:
    {{
        "skill": "{skill}",
        "misconception_id": "MIS-01",
        "misconception": "Misconception name",
        "confidence": 0.0,
        "evidence": "Evidence from the student's code"
    }}
    """

    try:
        response = _get_client().models.generate_content(
            model=_get_model_name(),
            contents=prompt,
            config=types.GenerateContentConfig(
                response_mime_type="application/json"
            )
        )

        diagnosis = json.loads(response.text)

        allowed_ids = {item["id"] for item in domain_data["misconceptions"]}

        required_keys = [
            "skill",
            "misconception_id",
            "misconception",
            "confidence",
            "evidence"
        ]

        for key in required_keys:
            if key not in diagnosis:
                raise ValueError(f"Missing field: {key}")

        if diagnosis["misconception_id"] not in allowed_ids | {"MC_UNKNOWN"}:
            raise ValueError("Gemini returned an invalid misconception ID")

        confidence = float(diagnosis["confidence"])
        if not 0 <= confidence <= 1:
            raise ValueError("Invalid confidence value")

        diagnosis["confidence"] = confidence

        if diagnosis["misconception_id"] in allowed_ids:
            official_item = next(
                item for item in domain_data["misconceptions"]
                if item["id"] == diagnosis["misconception_id"]
            )
            diagnosis["misconception"] = official_item["name"]

        return diagnosis

    except Exception as exc:
        return {
            "skill": skill,
            "misconception_id": "MC_UNKNOWN",
            "misconception": "Diagnosis unavailable",
            "confidence": 0.0,
            "evidence": str(exc)
        }


# ---------------------------------------------------------------------------
# Orchestrator — three-layer pipeline
# ---------------------------------------------------------------------------

def diagnose_student_submission(
    student_code: str,
    runner_result: dict,
    skill: str = "SKL-02",
) -> dict:
    """Run all diagnosis layers and return the best result.

    Pipeline:
      1. Rule-based traceback matching  (instant)
      2. AST structural scanning        (fast, no network)
      3. Gemini API                     (slow, only when needed)
    """
    # --- Layer 1: Rule-based ---
    rule_result = diagnose_error(runner_result, skill)
    if rule_result["misconception_id"] not in (None, "MC_UNKNOWN"):
        rule_result["diagnosis_source"] = "rule_based"
        return rule_result

    # Clean pass — no error, no wrong answer to diagnose
    if runner_result.get("passed"):
        rule_result["diagnosis_source"] = "rule_based"
        return rule_result

    # --- Layer 2: AST scanner ---
    ast_findings = scan_ast(student_code)
    if ast_findings:
        best = max(ast_findings, key=lambda x: x["confidence"])
        if best["confidence"] >= 0.70:
            entry = _lookup_misconception(best["misconception_id"])
            return {
                "skill": skill,
                "misconception_id": best["misconception_id"],
                "misconception": entry["name"] if entry else best["misconception_id"],
                "confidence": best["confidence"],
                "evidence": best["detail"],
                "diagnosis_source": "ast_scanner",
                "ast_line": best["line"],
            }

    # --- Layer 3: Gemini ---
    gemini_result = diagnose_with_gemini(student_code, runner_result, skill)
    gemini_result["diagnosis_source"] = "gemini"
    return gemini_result
