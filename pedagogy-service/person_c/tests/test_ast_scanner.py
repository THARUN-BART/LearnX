"""Tests for person_c.ast_scanner."""

import pytest

try:
    from person_c.ast_scanner import scan_ast
except ImportError:
    from ast_scanner import scan_ast


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def _find(findings, misconception_id):
    return next((f for f in findings if f["misconception_id"] == misconception_id), None)


# ---------------------------------------------------------------------------
# MIS-01 — Losing Head Pointer
# ---------------------------------------------------------------------------

_MIS01_CODE = """\
def traverse(head):
    while head:
        print(head.val)
        head = head.next  # BUG: should be curr = head, then curr = curr.next
"""

_CLEAN_TRAVERSAL = """\
def traverse(head):
    curr = head
    while curr:
        print(curr.val)
        curr = curr.next
"""


def test_mis01_detected_when_head_reassigned_in_while():
    findings = scan_ast(_MIS01_CODE)
    hit = _find(findings, "MIS-01")
    assert hit is not None
    assert hit["confidence"] >= 0.70
    assert hit["line"] is not None


def test_mis01_not_triggered_on_clean_traversal():
    findings = scan_ast(_CLEAN_TRAVERSAL)
    hit = _find(findings, "MIS-01")
    assert hit is None


# ---------------------------------------------------------------------------
# MIS-04 — Potential infinite loop (no advancement)
# ---------------------------------------------------------------------------

_MIS04_CODE = """\
def find(head, val):
    curr = head
    while curr:
        if curr.val == val:
            return curr
        # BUG: forgot curr = curr.next
"""

_MIS04_CLEAN = """\
def find(head, val):
    curr = head
    while curr:
        if curr.val == val:
            return curr
        curr = curr.next
    return None
"""


def test_mis04_detected_when_no_advancement():
    findings = scan_ast(_MIS04_CODE)
    hit = _find(findings, "MIS-04")
    assert hit is not None
    assert hit["confidence"] >= 0.60


def test_mis04_not_triggered_on_clean_loop():
    findings = scan_ast(_MIS04_CLEAN)
    hit = _find(findings, "MIS-04")
    assert hit is None


# ---------------------------------------------------------------------------
# ARR-MIS-01 — Off-by-one in range
# ---------------------------------------------------------------------------

_ARR_MIS01_LEN_CODE = """\
def sum_array(arr):
    total = 0
    for i in range(len(arr) + 1):  # BUG: should be range(len(arr))
        total += arr[i]
    return total
"""

_ARR_MIS01_N_CODE = """\
def print_all(n, arr):
    for i in range(n + 1):  # BUG
        print(arr[i])
"""

_ARR_CLEAN = """\
def sum_array(arr):
    return sum(arr)
"""


def test_arr_mis01_detected_with_len_plus_one():
    findings = scan_ast(_ARR_MIS01_LEN_CODE)
    hit = _find(findings, "ARR-MIS-01")
    assert hit is not None
    assert hit["confidence"] >= 0.85  # len() case gets high confidence


def test_arr_mis01_detected_with_n_plus_one():
    findings = scan_ast(_ARR_MIS01_N_CODE)
    hit = _find(findings, "ARR-MIS-01")
    assert hit is not None


def test_arr_mis01_not_triggered_on_clean_code():
    findings = scan_ast(_ARR_CLEAN)
    hit = _find(findings, "ARR-MIS-01")
    assert hit is None


# ---------------------------------------------------------------------------
# Empty / unparseable input
# ---------------------------------------------------------------------------

def test_empty_code_returns_empty_list():
    assert scan_ast("") == []


def test_whitespace_only_returns_empty_list():
    assert scan_ast("   \n\n  ") == []


def test_syntax_error_returns_empty_list():
    assert scan_ast("def broken(:\n    pass") == []


# ---------------------------------------------------------------------------
# Deduplication
# ---------------------------------------------------------------------------

_DOUBLE_MIS01 = """\
def bad(head):
    while head:
        head = head.next   # line 3 — MIS-01
    while head:
        head = head.next   # line 5 — MIS-01 again
"""


def test_findings_deduplicated_by_misconception_id():
    findings = scan_ast(_DOUBLE_MIS01)
    ids = [f["misconception_id"] for f in findings]
    # Each ID should appear at most once
    assert len(ids) == len(set(ids))
