import sys
from pathlib import Path
from hatchling.build.hooks.plugin.interface import BuildHookInterface

# 1. Dynamically append the package subdirectory to the python search path.
# This ensures hatch_build.py can import build.py cleanly during an sdist build.
PACKAGE_ROOT = Path(__file__).parent / "python" / "ogopego"
sys.path.insert(0, str(PACKAGE_ROOT))

import build  # type: ignore


class GoMatrixBuildHook(BuildHookInterface):
    def initialize(self, version, build_data):
        # 2. Extract the target platform configuration tag provided by the installer
        platform_tag = build_data.get("infer_tag", {}).get("platform", "").lower()

        # 3. Use build.py's native system parser as a solid baseline
        goos, goarch = build.normalize_platform()

        # 4. Refine the baseline choices if the installer specifies a specific target tag
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

        # 5. Delegate compilation entirely to the single source of truth module
        # This writes the compiled binary directly into your chosen layout ('lib/')
        build.compile_binary(goos=goos, goarch=goarch)
