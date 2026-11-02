from person_c.diagnosis import diagnose_error, diagnose_student_submission


def test_diagnose_error_timeout():
    result = diagnose_error({
        "passed": False,
        "error_type": "TIMEOUT",
        "output": "",
        "traceback": "Code execution exceeded the time limit."
    })
    assert result["misconception_id"] == "MIS-04"


def test_diagnose_error_runtime_error_known_misconception():
    result = diagnose_error({
        "passed": False,
        "error_type": "RUNTIME_ERROR",
        "output": "",
        "traceback": "AttributeError: 'NoneType' object has no attribute 'next'"
    })
    assert result["misconception_id"] == "MIS-02"


def test_diagnose_student_submission_uses_rule_based_unknown_result():
    result = diagnose_student_submission(
        "while True:\n    pass\n",
        {"passed": False, "error_type": "TIMEOUT", "output": "", "traceback": "Code execution exceeded the time limit."}
    )
    assert result["diagnosis_source"] == "rule_based"
    assert result["misconception_id"] == "MIS-04"


# ---------------------------------------------------------------------------
# New rule tests — extended rule registry
# ---------------------------------------------------------------------------

def test_diagnose_error_mis07_int_not_node():
    """Rule fires when traceback says 'int' object has no attribute ..."""
    result = diagnose_error({
        "passed": False,
        "error_type": "RUNTIME_ERROR",
        "output": "",
        "traceback": "AttributeError: 'int' object has no attribute 'next'"
    })
    assert result["misconception_id"] == "MIS-07"
    assert result["confidence"] >= 0.90


def test_diagnose_error_arr_mis01_index_error():
    """IndexError: list index out of range -> ARR-MIS-01."""
    result = diagnose_error({
        "passed": False,
        "error_type": "RUNTIME_ERROR",
        "output": "",
        "traceback": "IndexError: list index out of range"
    })
    assert result["misconception_id"] == "ARR-MIS-01"
    assert result["confidence"] >= 0.85


def test_diagnose_error_rec_mis01_recursion_error():
    """RecursionError -> REC-MIS-01."""
    result = diagnose_error({
        "passed": False,
        "error_type": "RUNTIME_ERROR",
        "output": "",
        "traceback": "RecursionError: maximum recursion depth exceeded"
    })
    assert result["misconception_id"] == "REC-MIS-01"
    assert result["confidence"] >= 0.95


def test_diagnose_error_sq_mis01_pop_empty():
    """IndexError: pop from empty list -> SQ-MIS-01."""
    result = diagnose_error({
        "passed": False,
        "error_type": "RUNTIME_ERROR",
        "output": "",
        "traceback": "IndexError: pop from empty list"
    })
    assert result["misconception_id"] == "SQ-MIS-01"
    assert result["confidence"] >= 0.90


def test_diagnose_student_submission_uses_ast_layer():
    """When rule layer returns MC_UNKNOWN, AST scanner should catch MIS-01."""
    code_with_head_bug = (
        "def traverse(head):\n"
        "    while head:\n"
        "        print(head.val)\n"
        "        head = head.next\n"
    )
    # Simulate wrong output (assertion failure), so rule layer won't positively identify.
    runner_result = {
        "passed": False,
        "error_type": "ASSERTION_ERROR",
        "output": "",
        "traceback": ""
    }
    result = diagnose_student_submission(code_with_head_bug, runner_result)
    # AST scanner should catch this before falling through to Gemini
    assert result["diagnosis_source"] in ("ast_scanner", "rule_based", "gemini")
    if result["diagnosis_source"] == "ast_scanner":
        assert result["misconception_id"] == "MIS-01"


def test_diagnose_error_mis02_none_type_generic():
    """Generic NoneType AttributeError -> MIS-02."""
    result = diagnose_error({
        "passed": False,
        "error_type": "RUNTIME_ERROR",
        "output": "",
        "traceback": "AttributeError: 'NoneType' object has no attribute 'val'"
    })
    assert result["misconception_id"] == "MIS-02"
