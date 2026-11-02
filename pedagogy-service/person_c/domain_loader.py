import json
from pathlib import Path


def load_domain_data(path=None):
    """Load the domain knowledge base and validate its expected top-level fields.

    The JSON file may contain a ``new_misconceptions`` key with extended domain
    entries (arrays, binary search, stacks, etc.).  These are merged into the
    main ``misconceptions`` list at load time so the rest of the system sees a
    single unified list.
    """
    if path is None:
        base_dir = Path(__file__).resolve().parent
        candidates = [
            base_dir / "domain_knowledge_base.json",
            base_dir / "domain_knowledge_base.md"
        ]
        for candidate in candidates:
            if candidate.exists():
                path = candidate
                break
        else:
            raise FileNotFoundError(
                "No domain knowledge file found in the project directory."
            )

    with open(path, "r", encoding="utf-8") as file:
        data = json.load(file)

    required_fields = {"skills", "misconceptions", "questions"}
    missing_fields = required_fields - set(data.keys())
    if missing_fields:
        raise ValueError(
            f"Domain knowledge file is missing required top-level fields: {sorted(missing_fields)}"
        )

    # Merge extended misconceptions (new domains) into the main list.
    # Using a separate top-level key keeps existing MIS-* entries untouched.
    if "new_misconceptions" in data:
        existing_ids = {m["id"] for m in data["misconceptions"]}
        for entry in data["new_misconceptions"]:
            if entry["id"] not in existing_ids:
                data["misconceptions"].append(entry)

    return data
