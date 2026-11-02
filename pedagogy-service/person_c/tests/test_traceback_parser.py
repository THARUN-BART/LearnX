"""Tests for person_c.traceback_parser."""

import pytest

try:
    from person_c.traceback_parser import parse_traceback
except ImportError:
    from traceback_parser import parse_traceback


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

_NONE_TYPE_TB = """\
Traceback (most recent call last):
  File "sol.py", line 12, in traverse
    print(curr.val)
AttributeError: 'NoneType' object has no attribute 'val'
"""

_INDEX_ERROR_TB = """\
Traceback (most recent call last):
  File "sol.py", line 8, in find_max
    if arr[i] > arr[i + 1]:
IndexError: list index out of range
"""

_RECURSION_TB = """\
Traceback (most recent call last):
  File "sol.py", line 4, in factorial
    return n * factorial(n)
  [Previous line repeated 996 more times]
RecursionError: maximum recursion depth exceeded
"""

_ZERO_DIV_TB = """\
Traceback (most recent call last):
  File "sol.py", line 7, in divide
    return a / b
ZeroDivisionError: division by zero
"""

_EMPTY = ""
_NONE = None


# ---------------------------------------------------------------------------
# Tests
# ---------------------------------------------------------------------------

def test_parses_attribute_error_none_type():
    result = parse_traceback(_NONE_TYPE_TB)
    assert result["exception_type"] == "AttributeError"
    assert "NoneType" in result["exception_msg"]
    assert result["failing_lineno"] == 12
    assert result["raw"] == _NONE_TYPE_TB


def test_parses_index_error():
    result = parse_traceback(_INDEX_ERROR_TB)
    assert result["exception_type"] == "IndexError"
    assert result["error_token"] == "INDEX_ERROR"
    assert "list index out of range" in result["exception_msg"]
    assert result["failing_lineno"] == 8


def test_parses_recursion_error():
    result = parse_traceback(_RECURSION_TB)
    assert result["exception_type"] == "RecursionError"
    assert result["error_token"] == "RECURSION_ERROR"


def test_parses_zero_division_error():
    result = parse_traceback(_ZERO_DIV_TB)
    assert result["exception_type"] == "ZeroDivisionError"
    assert result["error_token"] == "ZERO_DIVISION_ERROR"
    assert result["failing_lineno"] == 7


def test_empty_string_returns_none_fields():
    result = parse_traceback(_EMPTY)
    assert result["exception_type"] is None
    assert result["error_token"] is None
    assert result["exception_msg"] is None
    assert result["failing_lineno"] is None
    assert result["raw"] == _EMPTY


def test_none_input_returns_none_fields():
    result = parse_traceback(_NONE)
    assert result["exception_type"] is None
    assert result["raw"] == ""


def test_no_traceback_just_exception():
    raw = "ValueError: invalid literal for int() with base 10: 'x'"
    result = parse_traceback(raw)
    assert result["exception_type"] == "ValueError"
    assert "invalid literal" in result["exception_msg"]
    assert result["failing_lineno"] is None
