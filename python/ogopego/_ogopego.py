import subprocess
import sys
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
