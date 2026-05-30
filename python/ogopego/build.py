import os
import platform
import subprocess
import sys
from pathlib import Path

# Centralized configuration constants
BIN_DIR_NAME = "bin"

MATRIX = {
    ("darwin", "amd64"): "ogopego_darwin_amd64",
    ("darwin", "arm64"): "ogopego_darwin_arm64",
    ("linux", "amd64"): "ogopego_linux_amd64",
    ("linux", "arm64"): "ogopego_linux_arm64",
    ("windows", "amd64"): "ogopego_windows_amd64.exe",
}

WHEEL_TAG_MATRIX = {
    ("darwin", "amd64"): "macosx_10_12_x86_64",
    ("darwin", "arm64"): "macosx_11_0_arm64",
    ("linux", "amd64"): "manylinux_2_28_x86_64",
    ("linux", "arm64"): "manylinux_2_28_aarch64",
    ("windows", "amd64"): "win_amd64",
}

def get_wheel_tag(goos: str, goarch: str) -> str:
    """Returns the official PEP 425 platform tag for a given Go target."""
    return WHEEL_TAG_MATRIX.get((goos, goarch), "any")


def normalize_platform(goos: str = None, goarch: str = None) -> tuple[str, str]:
    """Normalizes host tokens into standard Go GOOS/GOARCH tokens."""
    if not goos:
        system = platform.system().lower()
        if "darwin" in system:
            goos = "darwin"
        elif "win" in system:
            goos = "windows"
        else:
            goos = "linux"

    if not goarch:
        machine = platform.machine().lower()
        if machine in ["x86_64", "amd64"]:
            goarch = "amd64"
        elif "arm64" in machine or "aarch64" in machine:
            goarch = "arm64"
        else:
            goarch = "amd64"

    return goos, goarch

def get_binary_name(goos: str = None, goarch: str = None) -> str:
    """Returns the standardized filename of the binary for a platform configuration."""
    goos, goarch = normalize_platform(goos, goarch)
    binary_name = MATRIX.get((goos, goarch))

    if not binary_name:
        ext = ".exe" if goos == "windows" else ""
        binary_name = f"ogopego_{goos}_{goarch}{ext}"

    return binary_name

def get_binary_path(goos: str = None, goarch: str = None) -> Path:
    """Resolves the absolute path to the binary artifact based on package location."""
    package_dir = Path(__file__).parent.resolve()
    lib_dir = package_dir / BIN_DIR_NAME
    return lib_dir / get_binary_name(goos, goarch)

def compile_binary(goos: str, goarch: str) -> Path:
    """Compiles a single statically linked Go binary directly into the owned lib target."""
    goos, goarch = normalize_platform(goos, goarch)
    output_path = get_binary_path(goos, goarch)

    # Ensure the target directory exists before running the compiler
    output_path.parent.mkdir(parents=True, exist_ok=True)

    print(f"Compiling Go binary for target: {goos}/{goarch} -> {output_path}")

    env = os.environ.copy()
    env["GOOS"] = goos
    env["GOARCH"] = goarch
    env["CGO_ENABLED"] = "0"

    # Crucial: Execute the compiler from the root of the Go module context
    # Assumes go.mod is two directories above python/ogopego/build.py
    go_module_root = Path(__file__).parent.parent.parent.resolve()

    cmd = ["go", "build", "-o", str(output_path), "./api"]

    try:
        subprocess.run(cmd, env=env, cwd=str(go_module_root), check=True)
    except FileNotFoundError:
        print("Error: The 'go' compiler executable could not be found.", file=sys.stderr)
        sys.exit(1)
    except subprocess.CalledProcessError as e:
        print(f"Go compilation failed for target {goos}/{goarch}: {e}", file=sys.stderr)
        sys.exit(1)

    return output_path

def compile_all():
    """Loops through the entire matrix configuration to build all binary artifacts."""
    print(f"Initializing matrix build for {len(MATRIX)} targets...")
    for goos, goarch in MATRIX.keys():
        compile_binary(goos, goarch)

if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] == "--all":
        compile_all()
    else:
        compile_binary(goos=None, goarch=None)
