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


_CFG_KEYS = frozenset({"trace", "color"})


def compile(
    grammar: str,
    name: str | None = None,
    *,
    config=None,
    filename=None,
    basetype=None,
    semantics=None,
    asmodel=False,
    builderconfig=None,
    synthok=True,
    typedefs=None,
    constructors=None,
    **settings,
) -> dict:
    _check_unsupported(
        name=name,
        config=config,
        filename=filename,
        basetype=basetype,
        semantics=semantics,
        asmodel=asmodel,
        builderconfig=builderconfig,
        synthok=synthok,
        typedefs=typedefs,
        constructors=constructors,
    )
    unknown = {k for k in settings if k not in _CFG_KEYS}
    if unknown:
        _raise_unknown("compile", settings)
    trace = settings.get("trace", False)
    color = settings.get("color", "auto")

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
    config=None,
    start=None,
    name=None,
    filename=None,
    semantics=None,
    asmodel=False,
    builderconfig=None,
    basetype=None,
    synthok=True,
    typedefs=None,
    constructors=None,
    **settings,
):
    _check_unsupported(
        config=config,
        start=start,
        name=name,
        filename=filename,
        semantics=semantics,
        asmodel=asmodel,
        builderconfig=builderconfig,
        basetype=basetype,
        synthok=synthok,
        typedefs=typedefs,
        constructors=constructors,
    )
    unknown = {k for k in settings if k not in _CFG_KEYS}
    if unknown:
        _raise_unknown("parse", settings)
    trace = settings.get("trace", False)
    color = settings.get("color", "auto")

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
    config=None,
    start=None,
    name=None,
    filename=None,
    semantics=None,
    asmodel=False,
    builderconfig=None,
    basetype=None,
    synthok=True,
    typedefs=None,
    constructors=None,
    **settings,
):
    _check_unsupported(
        config=config,
        start=start,
        name=name,
        filename=filename,
        semantics=semantics,
        asmodel=asmodel,
        builderconfig=builderconfig,
        basetype=basetype,
        synthok=synthok,
        typedefs=typedefs,
        constructors=constructors,
    )
    unknown = {k for k in settings if k not in _CFG_KEYS}
    if unknown:
        _raise_unknown("parse_file", settings)
    trace = settings.get("trace", False)
    color = settings.get("color", "auto")

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


def _raise_unknown(caller: str, settings: dict) -> None:
    items = ", ".join(
        f"{k}={v!r}" for k, v in settings.items() if k not in _CFG_KEYS
    )
    raise ValueError(f"ogopego {caller} does not support: {items}")
