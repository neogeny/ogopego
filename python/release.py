import os
import shutil
import subprocess
import sys
from pathlib import Path

# Ensure the parent of this script (python/) is on sys.path so
# "from ogopego import build" resolves correctly whether run as
# "python python/release.py" or "python -m python.release".
_SCRIPTS_DIR = Path(__file__).parent.resolve()
sys.path.insert(0, str(_SCRIPTS_DIR))

from ogopego import build  # type: ignore

PROJECT_ROOT = _SCRIPTS_DIR.parent
PACKAGE_ROOT = PROJECT_ROOT / "python" / "ogopego"


def clean_workspace():
    """Cleans up stale build artifacts to ensure reproducible packaging."""
    dist_dir = PROJECT_ROOT / "dist"
    lib_dir = PACKAGE_ROOT / build.BIN_DIR_NAME

    if dist_dir.exists():
        shutil.rmtree(dist_dir)
    if lib_dir.exists():
        shutil.rmtree(lib_dir)

    lib_dir.mkdir(parents=True, exist_ok=True)


def build_targeted_wheels():
    """Builds individual, optimized wheels for each OS/Arch matrix target."""

    build.compile_all()

    print("\n--- Starting Targeted Wheel Packaging with uv ---")

    for (goos, goarch), binary_name in build.MATRIX.items():
        wheel_tag = build.get_wheel_tag(goos, goarch)

        print(f"\nPackaging wheel for target: {goos}_{goarch} ({wheel_tag})...")

        env = os.environ.copy()
        env["HATCH_BUILD_HOOKS_ENABLE"] = "false"
        env["HATCH_WHEEL_PLATFORM_TAG"] = wheel_tag
        env["HATCH_BIN_NAME"] = binary_name

        cmd = [
            "uv", "build", "--wheel",
            "--force-pep517",
            "-C", f"build-data.tag=py3-none-{wheel_tag}",
            str(PROJECT_ROOT),
        ]

        try:
            subprocess.run(cmd, env=env, cwd=str(PROJECT_ROOT), check=True)
        except subprocess.CalledProcessError as e:
            print(
                f"Failed to package wheel distribution for target {wheel_tag}: {e}",
                file=sys.stderr,
            )
            sys.exit(1)

    print("\nMatrix wheel build successfully complete. Check your ./dist/ folder.")


if __name__ == "__main__":
    clean_workspace()
    build.compile_all()


