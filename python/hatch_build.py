import os
import platform
import sys
from pathlib import Path
from hatchling.builders.hooks.plugin.interface import BuildHookInterface

PACKAGE_ROOT = Path(__file__).parent
sys.path.insert(0, str(PACKAGE_ROOT))

from ogopego import build  # type: ignore


class CustomBuildHook(BuildHookInterface):
    def initialize(self, version, build_data):
        # Instruct Hatchling that this package outputs binary wheels
        build_data["pure_python"] = False

        # Check if an explicit cross-compilation wheel tag was provided by release.py
        platform_tag = os.environ.get("HATCH_WHEEL_PLATFORM_TAG", "").lower()

        if platform_tag:
            goos, goarch = build.normalize_platform()

            if "linux" in platform_tag:
                goos = "linux"
            elif "mac" in platform_tag or "darwin" in platform_tag:
                goos = "darwin"
            elif "win" in platform_tag:
                goos = "windows"

            if "x86_64" in platform_tag or "amd64" in platform_tag:
                goarch = "amd64"
            elif "arm64" in platform_tag or "aarch64" in platform_tag:
                goarch = "arm64"
        else:
            # Local fallback route: Use standard host architecture info
            goos = "windows" if os.name == "nt" else sys.platform
            if goos == "darwin":
                goos = "darwin"
            elif goos.startswith("linux"):
                goos = "linux"

            arch_raw = platform.machine().lower()
            if arch_raw in ("x86_64", "amd64"):
                goarch = "amd64"
            elif arch_raw in ("arm64", "aarch64"):
                goarch = "arm64"
            else:
                goarch = arch_raw

        build_data["tag"] = f"py3-none-{build.get_wheel_tag(goos, goarch)}"

        expected_binary = build.get_binary_path(goos, goarch)

        if not expected_binary.exists():
            build.compile_binary(goos=goos, goarch=goarch)
