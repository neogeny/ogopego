import os
import platform
import shutil
import subprocess
import sys
import tarfile
import urllib.request
import zipfile
from pathlib import Path

# Explicit version tag to prevent dynamic breaking changes
GO_VERSION = "1.26.3"

# Release lookup map pointing back to official download storage mirrors
DOWNLOAD_MATRIX = {
    ("darwin", "amd64"): f"go{GO_VERSION}.darwin-amd64.tar.gz",
    ("darwin", "arm64"): f"go{GO_VERSION}.darwin-arm64.tar.gz",
    ("linux", "amd64"): f"go{GO_VERSION}.linux-amd64.tar.gz",
    ("linux", "arm64"): f"go{GO_VERSION}.linux-arm64.tar.gz",
    ("windows", "amd64"): f"go{GO_VERSION}.windows-amd64.zip",
}

def get_current_tokens() -> tuple[str, str]:
    """Extracts standardized host platform tokens."""
    system = platform.system().lower()
    goos = "darwin" if "darwin" in system else "windows" if "win" in system else "linux"
    
    machine = platform.machine().lower()
    goarch = "amd64" if machine in ["x86_64", "amd64"] else "arm64" if "arm" in machine else "amd64"
    return goos, goarch

def find_system_go() -> str | None:
    """Checks if the Go toolchain is already present in the user's default PATH."""
    return shutil.which("go")

def ensure_go_installed() -> Path:
    """Verifies Go existence or performs a localized sandbox installation."""
    # 1. Immediate path resolution bypass if already present
    system_go = find_system_go()
    if system_go:
        return Path(system_go)

    # 2. Configure target installation folder parameters
    package_root = Path(__file__).parent.resolve()
    go_sandbox_root = package_root / ".go_cache"
    go_binary_path = go_sandbox_root / "go" / "bin" / ("go.exe" if platform.system().lower() == "windows" else "go")

    if go_binary_path.exists():
        return go_binary_path

    print(f"Go toolchain not found on host. Initiating localized sandbox install of Go v{GO_VERSION}...")
    goos, goarch = get_current_tokens()
    archive_name = DOWNLOAD_MATRIX.get((goos, goarch))

    if not archive_name:
        raise RuntimeError(f"Unsupported bootstrapping platform matrix target: {goos}/{goarch}")

    url = f"https://go.dev/dl/{archive_name}"
    download_target = go_sandbox_root / archive_name
    
    go_sandbox_root.mkdir(parents=True, exist_ok=True)

    # 3. Stream binary distribution package from official mirrors
    print(f"Downloading toolchain from: {url}")
    try:
        urllib.request.urlretrieve(url, download_target)
    except Exception as e:
        raise RuntimeError(f"Failed to fetch Go distribution from mirror: {e}")

    # 4. Extract based on type specifications without external dependencies
    print("Extracting package payloads...")
    if archive_name.endswith(".zip"):
        with zipfile.ZipFile(download_target, 'r') as zip_ref:
            zip_ref.extractall(go_sandbox_root)
    else:
        with tarfile.open(download_target, "r:gz") as tar_ref:
            tar_ref.extractall(go_sandbox_root)

    # Clean up installation file archive
    download_target.unlink()

    if not go_binary_path.exists():
        raise RuntimeError("Extraction completed but binary path resolution validation failed.")

    print(f"Localized Go toolchain successfully established at: {go_binary_path}")
    return go_binary_path

if __name__ == "__main__":
    try:
        resolved_path = ensure_go_installed()
        print(f"SUCCESS: Go executable target available at -> {resolved_path}")
    except Exception as e:
        print(f"BOOTSTRAP FAILURE: {e}", file=sys.stderr)
        sys.exit(1)
