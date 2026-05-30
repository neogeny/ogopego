"""Python wrapper around the ogopego Go CLI binary.

Provides compile(), parse(), and parse_file() with a TatSu-compatible
signature. The Go binary is bundled in the wheel and invoked via
subprocess.
"""

import hashlib
import json
import subprocess
import sys
import tempfile
from pathlib import Path

from .build import get_binary_path


def ogo(*args: str, **kwargs) -> subprocess.CompletedProcess:
    """Executes the embedded Go binary with the provided arguments.

    Forwards stdin, stdout, and stderr to the parent process by default
    unless overridden via kwargs (e.g., capture_output=True).
    """
    binary_path = get_binary_path()

    if not binary_path.exists():
        raise FileNotFoundError(
            f"Native ogopego binary component is missing at: {binary_path}"
        )

    # Convert all arguments to strings and build the execution command list
    cmd = [str(binary_path)] + [str(arg) for arg in args]

    # Provide safe terminal defaults while allowing complete subprocess parameter overrides
    kwargs.setdefault("stdout", None)
    kwargs.setdefault("stderr", None)
    kwargs.setdefault("stdin", None)

    try:
        return subprocess.run(cmd, check=True, **kwargs)
    except subprocess.CalledProcessError as e:
        # Transparently pass through the exit code if the Go executable fails internally
        sys.exit(e.returncode)


__cache: dict[str, dict] = {}

_UNSUPPORTED_DEFAULTS: dict[str, object] = {
    "config": None,
    "basetype": None,
    "semantics": None,
    "asmodel": False,
    "builderconfig": None,
    "synthok": True,
    "typedefs": None,
    "constructors": None,
    "start": None,
    "name": None,
    "filename": None,
}


class OgoError(Exception):
    """Raised when the ogopego Go binary exits with a non-zero status."""

    def __init__(self, returncode: int, stderr: str = ""):
        self.returncode = returncode
        self.stderr = stderr
        super().__init__(f"ogo failed (exit {returncode}): {stderr}")


def _hasha(text: str) -> str:
    return hashlib.sha256(str(text).encode("utf-8")).hexdigest()


def _binary_path() -> Path:
    p = get_binary_path()
    if not p.exists():
        raise FileNotFoundError(
            f"Native ogopego binary component is missing at: {p}"
        )
    return p


def _ogo_capture(*args: str) -> str:
    binary = _binary_path()
    cmd = [str(binary), *args]
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        raise OgoError(result.returncode, result.stderr)
    if not result.stdout.strip():
        raise OgoError(result.returncode, result.stderr)
    return result.stdout


def _check_unsupported(**kwargs: object) -> None:
    bad = {}
    for name, value in kwargs.items():
        if name in _UNSUPPORTED_DEFAULTS and value is not _UNSUPPORTED_DEFAULTS[name]:
            bad[name] = value
    if bad:
        items = ", ".join(f"{k}={v!r}" for k, v in bad.items())
        raise ValueError(f"ogopego does not support: {items}")


def _build_cli_args(
    subcommand: str,
    sub_flags: list[str],
    positional: list[str],
    trace: bool = False,
    color: str = "auto",
) -> list[str]:
    args: list[str] = []
    if trace:
        args.append("--trace")
    # Force plain output when auto-detecting, since stdout is always captured
    args.append("--color")
    args.append("never" if color == "auto" else color)
    args.append(subcommand)
    args.extend(sub_flags)
    args.extend(positional)
    return args


def compile(
    grammar: str,
    name: str | None = None,
    *,
    filename=None,
    trace: bool = False,
    color: str = "auto",
) -> dict:
    """Compile a PEG grammar string and return the compiled grammar as a dict.

    The result is cached by grammar content hash so repeated calls with
    the same grammar string return the cached result.

    Args:
        grammar: PEG grammar source text.
        name: Grammar name (unused, raises ValueError if set).
        filename: Source filename hint (unused, raises ValueError if set).
        trace: Enable trace output from the Go binary.
        color: Color mode for CLI output ("auto", "never", "always").

    Returns:
        A dict representing the compiled grammar (JSON-serializable).

    Raises:
        OgoError: If the Go binary fails.
        ValueError: If an unsupported argument is provided.
    """
    _check_unsupported(name=name, filename=filename)

    key = _hasha(grammar)
    if key in __cache:
        return __cache[key]

    with tempfile.NamedTemporaryFile(
        mode="w", suffix=".ebnf", delete=False
    ) as f:
        f.write(grammar)
        tmpfile = f.name

    try:
        cli_args = _build_cli_args(
            "grammar", ["--json"], [tmpfile],
            trace=trace, color=color,
        )
        stdout = _ogo_capture(*cli_args)
        result = json.loads(stdout)
        __cache[key] = result
        return result
    finally:
        Path(tmpfile).unlink(missing_ok=True)


def parse(
    grammar: str,
    text: str,
    /,
    *,
    start=None,
    name=None,
    filename=None,
    trace: bool = False,
    color: str = "auto",
):
    """Parse input text against a PEG grammar and return the parse tree.

    Compiles the grammar first (with caching), then runs the parser on
    the input text via the Go binary.

    Args:
        grammar: PEG grammar source text.
        text: Input text to parse.
        start: Start rule name (unused, raises ValueError if set).
        name: Grammar name hint (unused, raises ValueError if set).
        filename: Source filename hint (unused, raises ValueError if set).
        trace: Enable trace output from the Go binary.
        color: Color mode for CLI output ("auto", "never", "always").

    Returns:
        A dict representing the parse tree (JSON-serializable).

    Raises:
        OgoError: If compilation or parsing fails.
        ValueError: If an unsupported argument is provided.
    """
    _check_unsupported(start=start, name=name, filename=filename)

    grammar_dict = compile(grammar, trace=trace, color=color)

    with tempfile.NamedTemporaryFile(
        mode="w", suffix=".json", delete=False
    ) as f:
        json.dump(grammar_dict, f)
        grammar_file = f.name

    with tempfile.NamedTemporaryFile(
        mode="w", suffix=".txt", delete=False
    ) as f:
        f.write(text)
        input_file = f.name

    try:
        cli_args = _build_cli_args(
            "run", ["--json"], [grammar_file, input_file],
            trace=trace, color=color,
        )
        stdout = _ogo_capture(*cli_args)
        return json.loads(stdout)
    finally:
        Path(grammar_file).unlink(missing_ok=True)
        Path(input_file).unlink(missing_ok=True)


def parse_file(
    grammar: str,
    path: str,
    /,
    *,
    start=None,
    name=None,
    filename=None,
    trace: bool = False,
    color: str = "auto",
):
    """Parse a file against a PEG grammar and return the parse tree.

    Compiles the grammar first (with caching), then runs the parser on
    the file at *path* via the Go binary.

    Args:
        grammar: PEG grammar source text.
        path: Path to the input file to parse.
        start: Start rule name (unused, raises ValueError if set).
        name: Grammar name hint (unused, raises ValueError if set).
        filename: Source filename hint (unused, raises ValueError if set).
        trace: Enable trace output from the Go binary.
        color: Color mode for CLI output ("auto", "never", "always").

    Returns:
        A dict representing the parse tree (JSON-serializable).

    Raises:
        OgoError: If compilation or parsing fails.
        ValueError: If an unsupported argument is provided.
    """
    _check_unsupported(start=start, name=name, filename=filename)

    grammar_dict = compile(grammar, trace=trace, color=color)

    with tempfile.NamedTemporaryFile(
        mode="w", suffix=".json", delete=False
    ) as f:
        json.dump(grammar_dict, f)
        grammar_file = f.name

    try:
        cli_args = _build_cli_args(
            "run", ["--json"], [grammar_file, path],
            trace=trace, color=color,
        )
        stdout = _ogo_capture(*cli_args)
        return json.loads(stdout)
    finally:
        Path(grammar_file).unlink(missing_ok=True)
