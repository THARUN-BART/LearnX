"""AST-based structural scanner for common student misconceptions.

Analyses the abstract syntax tree of student-submitted Python code and
returns a list of potential misconception matches *without executing any
code*.  This provides a second layer of diagnosis between the fast
traceback-rule engine and the slower Gemini API fallback.

Each finding is a dict::

    {
        "misconception_id": str,   # e.g. "MIS-01"
        "confidence":       float, # 0.0 – 1.0
        "line":             int,   # 1-based source line of the suspicious node
        "detail":           str,   # human-readable explanation
    }
"""

import ast
from typing import Any


# ---------------------------------------------------------------------------
# Internal helpers
# ---------------------------------------------------------------------------

def _collect_names_in(node: ast.AST) -> set[str]:
    """Return the set of Name node ids that appear anywhere inside *node*."""
    return {n.id for n in ast.walk(node) if isinstance(n, ast.Name)}


def _is_advancement(stmt: ast.stmt, ptr_names: set[str]) -> bool:
    """Return True if *stmt* is an assignment like ``curr = curr.next``."""
    if not isinstance(stmt, ast.Assign):
        return False
    # LHS must be a simple Name that is one of the pointer names
    if len(stmt.targets) != 1 or not isinstance(stmt.targets[0], ast.Name):
        return False
    lhs_name = stmt.targets[0].id
    if lhs_name not in ptr_names:
        return False
    # RHS must be an Attribute access whose value is a Name
    if not isinstance(stmt.value, ast.Attribute):
        return False
    return True


# ---------------------------------------------------------------------------
# Individual checkers
# ---------------------------------------------------------------------------

def _check_mis01_head_reassign_in_loop(tree: ast.Module) -> list[dict]:
    """MIS-01 — Losing Head Pointer.

    Detects ``head = …`` (or any assignment to a variable named ``head``)
    inside the body of a ``while`` or ``for`` loop.  This is the classic
    mistake where the student moves the head pointer while traversing and
    loses the list.
    """
    findings: list[dict] = []

    class _Visitor(ast.NodeVisitor):
        def _scan_loop_body(self, body: list[ast.stmt]) -> None:
            for stmt in body:
                if isinstance(stmt, ast.Assign):
                    for target in stmt.targets:
                        if isinstance(target, ast.Name) and target.id == "head":
                            findings.append({
                                "misconception_id": "MIS-01",
                                "confidence": 0.82,
                                "line": stmt.lineno,
                                "detail": (
                                    "Variable 'head' is reassigned inside a loop body. "
                                    "This moves the head pointer and loses the list reference."
                                ),
                            })

        def visit_While(self, node: ast.While) -> None:
            self._scan_loop_body(node.body)
            self.generic_visit(node)

        def visit_For(self, node: ast.For) -> None:
            self._scan_loop_body(node.body)
            self.generic_visit(node)

    _Visitor().visit(tree)
    return findings


def _check_mis04_no_advancement_in_while(tree: ast.Module) -> list[dict]:
    """MIS-04 (early warning) — Potential Infinite Loop.

    Detects ``while`` loops whose body contains a loop-pointer variable
    (``curr``, ``current``, ``node``, ``ptr``) but *no* advancement
    statement that changes that pointer (``curr = curr.next`` etc.).
    """
    findings: list[dict] = []
    PTR_NAMES = {"curr", "current", "node", "ptr", "temp"}

    class _Visitor(ast.NodeVisitor):
        def visit_While(self, node: ast.While) -> None:
            used_ptrs = _collect_names_in(node) & PTR_NAMES
            if not used_ptrs:
                self.generic_visit(node)
                return

            # Check whether *any* statement in the immediate body is an
            # advancement of one of the pointer variables.
            has_advancement = any(
                _is_advancement(stmt, used_ptrs) for stmt in ast.walk(node)
                if isinstance(stmt, ast.Assign)
            )

            # Also allow augmented assignments like ``ptr = ptr.next``
            # already covered by _is_advancement; but we need to make sure
            # we don't false-positive on ``while True:`` with ``break``.
            has_break = any(isinstance(n, ast.Break) for n in ast.walk(node))

            if not has_advancement and not has_break:
                findings.append({
                    "misconception_id": "MIS-04",
                    "confidence": 0.70,
                    "line": node.lineno,
                    "detail": (
                        f"While loop at line {node.lineno} uses pointer variable(s) "
                        f"{used_ptrs} but has no advancement statement (e.g. curr = curr.next). "
                        "This may cause an infinite loop."
                    ),
                })
            self.generic_visit(node)

    _Visitor().visit(tree)
    return findings


def _check_mis03_wrong_pointer_order(tree: ast.Module) -> list[dict]:
    """MIS-03 — Wrong Pointer Reassignment Order.

    Detects the pattern where ``curr.next`` is overwritten *before*
    ``temp`` (or ``next_node``) is saved to it — the canonical pointer
    ordering bug in linked-list operations.
    """
    findings: list[dict] = []
    SAVE_NAMES = {"temp", "next_node", "saved", "nxt", "next"}

    class _FuncVisitor(ast.NodeVisitor):
        def visit_FunctionDef(self, node: ast.FunctionDef) -> None:
            self._check_body(node.body, node)
            self.generic_visit(node)

        visit_AsyncFunctionDef = visit_FunctionDef

        def _check_body(self, stmts: list[ast.stmt], _parent: Any) -> None:
            for i, stmt in enumerate(stmts):
                if not isinstance(stmt, ast.Assign):
                    continue
                # LHS is an Attribute like ``curr.next``
                if not (
                    len(stmt.targets) == 1
                    and isinstance(stmt.targets[0], ast.Attribute)
                    and stmt.targets[0].attr == "next"
                ):
                    continue
                # This statement overwrites X.next — check whether any
                # *later* statement in the same block saves X.next first.
                for later in stmts[i + 1:]:
                    if not isinstance(later, ast.Assign):
                        continue
                    if len(later.targets) != 1:
                        continue
                    target = later.targets[0]
                    if (
                        isinstance(target, ast.Name)
                        and target.id in SAVE_NAMES
                        and isinstance(later.value, ast.Attribute)
                        and later.value.attr == "next"
                    ):
                        findings.append({
                            "misconception_id": "MIS-03",
                            "confidence": 0.75,
                            "line": stmt.lineno,
                            "detail": (
                                f"'curr.next' is overwritten at line {stmt.lineno} "
                                f"but the next-pointer is only saved to '{target.id}' "
                                f"at line {later.lineno}. Save 'curr.next' BEFORE "
                                "overwriting it."
                            ),
                        })
                        break

    _FuncVisitor().visit(tree)
    return findings


def _check_arr_mis01_range_off_by_one(tree: ast.Module) -> list[dict]:
    """ARR-MIS-01 — Array Index Out of Bounds (static check).

    Detects ``range(len(x) + 1)`` or ``range(n + 1)`` patterns which
    iterate one step past the last valid index.
    """
    findings: list[dict] = []

    class _Visitor(ast.NodeVisitor):
        def visit_Call(self, node: ast.Call) -> None:
            # We're looking for range(something + 1)
            if not (
                isinstance(node.func, ast.Name) and node.func.id == "range"
                and node.args
            ):
                self.generic_visit(node)
                return

            arg = node.args[0]
            # Match: range(X + 1)
            if (
                isinstance(arg, ast.BinOp)
                and isinstance(arg.op, ast.Add)
                and isinstance(arg.right, ast.Constant)
                and arg.right.value == 1
            ):
                # Extra confidence if X is a len() call
                is_len = (
                    isinstance(arg.left, ast.Call)
                    and isinstance(arg.left.func, ast.Name)
                    and arg.left.func.id == "len"
                )
                findings.append({
                    "misconception_id": "ARR-MIS-01",
                    "confidence": 0.88 if is_len else 0.65,
                    "line": node.lineno,
                    "detail": (
                        f"range() at line {node.lineno} iterates up to "
                        f"{'len(…)' if is_len else 'n'} + 1, which is one past "
                        "the last valid index. Use range(len(arr)) or range(n)."
                    ),
                })
            self.generic_visit(node)

    _Visitor().visit(tree)
    return findings


def _check_mis07_int_as_node(tree: ast.Module) -> list[dict]:
    """MIS-07 — Value vs Node Confusion (static check).

    Detects attempts to access ``.next`` or ``.val`` on what is likely
    an integer: ``node = some_list[i]`` followed by ``node.next``, or
    ``curr = curr.val`` followed by ``curr.next``.

    This is a lightweight heuristic; it checks if ``.next`` is accessed
    on a variable that was assigned from a subscript or from ``.val``.
    """
    findings: list[dict] = []

    class _Visitor(ast.NodeVisitor):
        def visit_FunctionDef(self, node: ast.FunctionDef) -> None:
            # Build a rough map: variable → "int_likely" | "node_likely"
            int_likely: set[str] = set()

            for stmt in ast.walk(node):
                if not isinstance(stmt, ast.Assign):
                    continue
                if len(stmt.targets) != 1 or not isinstance(stmt.targets[0], ast.Name):
                    continue
                lhs = stmt.targets[0].id
                rhs = stmt.value
                # var = something[i]  → likely an integer value, not a node
                if isinstance(rhs, ast.Subscript):
                    int_likely.add(lhs)
                # var = other.val → definitely an integer
                if isinstance(rhs, ast.Attribute) and rhs.attr == "val":
                    int_likely.add(lhs)

            # Now look for .next access on int_likely variables
            for expr_node in ast.walk(node):
                if (
                    isinstance(expr_node, ast.Attribute)
                    and expr_node.attr == "next"
                    and isinstance(expr_node.value, ast.Name)
                    and expr_node.value.id in int_likely
                ):
                    findings.append({
                        "misconception_id": "MIS-07",
                        "confidence": 0.72,
                        "line": expr_node.lineno,
                        "detail": (
                            f"Variable '{expr_node.value.id}' at line {expr_node.lineno} "
                            "appears to hold an integer value (assigned from a subscript "
                            "or .val), but '.next' is accessed on it as if it were a node."
                        ),
                    })
            self.generic_visit(node)

        visit_AsyncFunctionDef = visit_FunctionDef

    _Visitor().visit(tree)
    return findings


# ---------------------------------------------------------------------------
# Public API
# ---------------------------------------------------------------------------

_CHECKERS = [
    _check_mis01_head_reassign_in_loop,
    _check_mis04_no_advancement_in_while,
    _check_mis03_wrong_pointer_order,
    _check_arr_mis01_range_off_by_one,
    _check_mis07_int_as_node,
]


def scan_ast(student_code: str) -> list[dict]:
    """Scan *student_code* for structural misconception patterns.

    Args:
        student_code: Raw Python source code as a string.

    Returns:
        A list of findings, each with keys
        ``misconception_id``, ``confidence``, ``line``, ``detail``.
        Returns an empty list if the code cannot be parsed or no
        patterns are detected.
    """
    if not student_code or not student_code.strip():
        return []

    try:
        tree = ast.parse(student_code)
    except SyntaxError:
        return []

    findings: list[dict] = []
    for checker in _CHECKERS:
        try:
            findings.extend(checker(tree))
        except Exception:  # pragma: no cover — defensive
            pass

    # Deduplicate by misconception_id, keeping the highest confidence
    seen: dict[str, dict] = {}
    for f in findings:
        mid = f["misconception_id"]
        if mid not in seen or f["confidence"] > seen[mid]["confidence"]:
            seen[mid] = f

    return list(seen.values())
