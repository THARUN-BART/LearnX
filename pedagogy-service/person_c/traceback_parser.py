"""Structured parser for Python exception tracebacks.

Extracts structured fields from raw stderr output so that diagnosis
rules can pattern-match on typed data instead of raw strings.
"""

import re


# Maps common Python exception names to short error-type tokens.
_EXCEPTION_TYPE_MAP = {
    "AssertionError": "ASSERTION_ERROR",
    "IndexError": "INDEX_ERROR",
    "AttributeError": "ATTRIBUTE_ERROR",
    "RecursionError": "RECURSION_ERROR",
    "ZeroDivisionError": "ZERO_DIVISION_ERROR",
    "MemoryError": "MEMORY_ERROR",
    "ValueError": "VALUE_ERROR",
    "TypeError": "TYPE_ERROR",
    "KeyError": "KEY_ERROR",
    "StopIteration": "STOP_ITERATION",
    "OverflowError": "OVERFLOW_ERROR",
}

_EXCEPTION_LINE_RE = re.compile(r"^(\w+(?:\.\w+)*(?:Error|Exception|Warning|Interrupt)):\s*(.*)")
_FILE_LINE_RE = re.compile(r'^\s*File ".*?", line (\d+)')


def parse_traceback(raw: str) -> dict:
    """Parse a raw Python traceback string into structured fields.

    Args:
        raw: Raw stderr output from a Python process.

    Returns:
        dict with keys:
            exception_type (str | None): Python exception class name, e.g. ``"IndexError"``.
            error_token    (str | None): Short token, e.g. ``"INDEX_ERROR"``, or the raw
                                        exception_type when not in the known map.
            exception_msg  (str | None): The message after the colon, e.g.
                                        ``"list index out of range"``.
            failing_lineno (int | None): Last ``File … line N`` number seen in the traceback.
            raw            (str):        The original input, unchanged.
    """
    lines = (raw or "").strip().splitlines()

    exception_type: str | None = None
    exception_msg: str | None = None
    failing_lineno: int | None = None

    for line in lines:
        m = _FILE_LINE_RE.match(line)
        if m:
            failing_lineno = int(m.group(1))

        m = _EXCEPTION_LINE_RE.match(line)
        if m:
            exception_type = m.group(1)
            exception_msg = m.group(2).strip()

    error_token: str | None = None
    if exception_type:
        # Try exact match first, then suffix-based lookup.
        error_token = _EXCEPTION_TYPE_MAP.get(exception_type)
        if error_token is None:
            for suffix, token in _EXCEPTION_TYPE_MAP.items():
                if exception_type.endswith(suffix):
                    error_token = token
                    break
        if error_token is None:
            error_token = exception_type.upper()

    return {
        "exception_type": exception_type,
        "error_token": error_token,
        "exception_msg": exception_msg,
        "failing_lineno": failing_lineno,
        "raw": raw or "",
    }
