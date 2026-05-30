import os
import platform
import shutil
import subprocess
import sys
import tarfile
import zipfile
from pathlib import Path
from urllib.request import urlopen

# Centralized configuration constants
BIN_DIR_NAME = "bin"

_GO_VERSION = "1.26.3"

_GO_DOWNLOAD_URL = (
    "https://go.dev/dl/go{version}.{os}-{arch}.tar.gz"
)

_GO_DOWNLOAD_URL_WINDOWS = (
    "https://go.dev/dl/go{version}.{os}-{arch}.zip"
)

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

_GO_CACHE_DIR = Path.home() / ".cache" / "ogopego" / "go"


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


def _go_os_name(goos: str) -> str:
    return {"darwin": "darwin", "linux": "linux", "windows": "windows"}.get(goos, goos)


def _go_arch_name(goarch: str) -> str:
    return {"amd64": "amd64", "arm64": "arm64"}.get(goarch, goarch)


def _ensure_go() -> tuple[Path, str | None]:
    """Return (go_binary, goroot_or_None).  goroot is None for system Go."""
    found = shutil.which("go")
    if found:
        return Path(found), None

    goos, goarch = normalize_platform()
    os_name = _go_os_name(goos)
    arch_name = _go_arch_name(goarch)

    cached_dir = _GO_CACHE_DIR / f"go{_GO_VERSION}" / "go"
    cached_bin = cached_dir / "bin" / ("go.exe" if goos == "windows" else "go")

    if cached_bin.exists():
        return cached_bin, str(cached_dir)

    _install_go(goos, os_name, arch_name, cached_dir)

    if not cached_bin.exists():
        print(
            f"Error: Go binary not found after installation at {cached_bin}",
            file=sys.stderr,
        )
        sys.exit(1)

    return cached_bin, str(cached_dir)


def _install_go(goos: str, os_name: str, arch_name: str, target_dir: Path) -> None:
    """Download and extract Go to the cache directory."""
    if goos == "windows":
        url = _GO_DOWNLOAD_URL_WINDOWS.format(
            version=_GO_VERSION, os=os_name, arch=arch_name
        )
    else:
        url = _GO_DOWNLOAD_URL.format(
            version=_GO_VERSION, os=os_name, arch=arch_name
        )

    archive_name = url.rsplit("/", 1)[-1]
    cache_file = _GO_CACHE_DIR / archive_name

    target_dir.parent.mkdir(parents=True, exist_ok=True)

    if not cache_file.exists():
        print(f"Downloading Go {_GO_VERSION} from {url}...")
        with urlopen(url) as response:
            with open(cache_file, "wb") as f:
                shutil.copyfileobj(response, f)

    print(f"Extracting Go {_GO_VERSION} to {target_dir.parent}...")
    if goos == "windows":
        with zipfile.ZipFile(cache_file, "r") as zf:
            zf.extractall(target_dir.parent)
    else:
        with tarfile.open(cache_file, "r:gz") as tf:
            tf.extractall(target_dir.parent)


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

    go_bin, goroot = _ensure_go()

    env = os.environ.copy()
    env["GOOS"] = goos
    env["GOARCH"] = goarch
    env["CGO_ENABLED"] = "0"
    env["PATH"] = str(go_bin.parent) + os.pathsep + env.get("PATH", "")
    if goroot is not None:
        env["GOROOT"] = goroot

    # Execute the compiler from the root of the Go module context
    go_module_root = Path(__file__).parent.parent.parent.resolve()

    cmd = [str(go_bin), "build", "-o", str(output_path), "./cmd/"]

    try:
        subprocess.run(cmd, env=env, cwd=str(go_module_root), check=True)
        output_path.chmod(0o755)
    except subprocess.CalledProcessError as e:
        print(
            f"Go compilation failed for target {goos}/{goarch}: {e}",
            file=sys.stderr,
        )
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
