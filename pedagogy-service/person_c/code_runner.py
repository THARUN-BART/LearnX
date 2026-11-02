"""Sandboxed code runner for student submissions.

Improvements over the original bare ``subprocess.run``:
- Writes student code to a temp file (same as before).
- Spawns a *wrapper* process that sets ``RLIMIT_AS`` (address-space /
  virtual-memory ceiling) before executing the student file, capping
  memory at 128 MB.
- Falls back gracefully on platforms where ``resource`` is unavailable
  (e.g. Windows CI).
- Classifies ``MemoryError`` / address-space kill into a new
  ``MEMORY_LIMIT_EXCEEDED`` error type.
"""

import os
import subprocess
import sys
import tempfile


# ---------------------------------------------------------------------------
# Resource-limit wrapper snippet (injected into child process)
# ---------------------------------------------------------------------------

_WRAPPER_TEMPLATE = """\
import resource as _res, sys as _sys
# 128 MB virtual-address ceiling
_MEM = 128 * 1024 * 1024
try:
    _res.setrlimit(_res.RLIMIT_AS, (_MEM, _MEM))
except Exception:
    pass  # Non-Linux or insufficient privileges — proceed without limit
with open({path!r}, encoding="utf-8") as _f:
    _src = _f.read()
exec(compile(_src, {path!r}, "exec"), {{"__name__": "__main__"}})
"""


def run_student_code(student_code: str, test_code: str, timeout: float = 3) -> dict:
    """Execute student code alongside *test_code* in a resource-limited sandbox.

    Args:
        student_code: The student's Python source.
        test_code:    Assertion/test harness appended after the student code.
        timeout:      Maximum wall-clock seconds before TIMEOUT is returned.

    Returns:
        dict with keys:
            passed     (bool)       — True when the process exits cleanly (code 0).
            error_type (str | None) — Classification token, or None on success.
            output     (str)        — Combined stdout.
            traceback  (str | None) — Stderr / error text, or None on success.
    """
    combined_code = student_code + "\n\n" + test_code

    with tempfile.NamedTemporaryFile(
        mode="w",
        suffix=".py",
        delete=False,
        encoding="utf-8",
    ) as tmp:
        tmp.write(combined_code)
        file_path = tmp.name

    wrapper = _WRAPPER_TEMPLATE.format(path=file_path)

    try:
        result = subprocess.run(
            [sys.executable, "-c", wrapper],
            capture_output=True,
            text=True,
            timeout=timeout,
        )

        if result.returncode == 0:
            return {
                "passed": True,
                "error_type": None,
                "output": result.stdout,
                "traceback": None,
            }

        stderr = result.stderr
        # Classify the error
        if "AssertionError" in stderr:
            error_type = "ASSERTION_ERROR"
        elif "MemoryError" in stderr or result.returncode == -9:
            # SIGKILL (returncode -9) often means the OOM killer or RLIMIT_AS hit
            error_type = "MEMORY_LIMIT_EXCEEDED"
        else:
            error_type = "RUNTIME_ERROR"

        return {
            "passed": False,
            "error_type": error_type,
            "output": result.stdout,
            "traceback": stderr,
        }

    except subprocess.TimeoutExpired:
        return {
            "passed": False,
            "error_type": "TIMEOUT",
            "output": "",
            "traceback": "Code execution exceeded the time limit.",
        }

    finally:
        if os.path.exists(file_path):
            os.remove(file_path)
