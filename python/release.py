import os
import shutil
import subprocess
import sys
from pathlib import Path

PROJECT_ROOT = Path(__file__).parent.parent
PACKAGE_ROOT = PROJECT_ROOT / "python" / "ogopego"
# sys.path.insert(0, str(PACKAGE_ROOT))

from ogopego import build  # type: ignore

def clean_workspace(root_dir: Path):
    """Cleans up stale build artifacts to ensure reproducible packaging."""
    dist_dir = root_dir / "dist"
    lib_dir = PACKAGE_ROOT / build.BIN_DIR_NAME

    if dist_dir.exists():
        shutil.rmtree(dist_dir)
    if lib_dir.exists():
        shutil.rmtree(lib_dir)
        
    lib_dir.mkdir(parents=True, exist_ok=True)

def build_targeted_wheels(root_dir: Path):
    """Builds individual, optimized wheels for each OS/Arch matrix target."""

    build.compile_all()

    print("\n--- Starting Targeted Wheel Packaging with uv ---")
    
    for (goos, goarch), binary_name in build.MATRIX.items():
        wheel_tag = build.get_wheel_tag(goos, goarch)
        _resolved_binary = build.get_binary_path(goos, goarch)
        
        print(f"\nPackaging wheel for target: {goos}_{goarch} ({wheel_tag})...")

        env = os.environ.copy()
        env["HATCH_BUILD_HOOKS_ENABLE"] = "false"
        
        # We instruct Hatchling to explicitly overwrite its target wheel platform tag
        # and specify exactly which binary file from our lib directory should be bundled.
        # Files not matching this pattern are skipped for this specific wheel package pass.
        env["HATCH_WHEEL_PLATFORM_TAG"] = wheel_tag
        env["HATCH_BIN_NAME"] = binary_name

        # Trigger uv build targeting only the wheel distribution channel
        cmd = [
            "uv", "build", "--wheel",
            "--force-pep517",
            "-C", f"build-data.tag=py3-none-{wheel_tag}",
            str(PROJECT_ROOT),
        ]

        try:
            # We enforce execution context from the root module folder
            subprocess.run(cmd, env=env, cwd=str(root_dir), check=True)
        except subprocess.CalledProcessError as e:
            print(f"Failed to package wheel distribution for target {wheel_tag}: {e}", file=sys.stderr)
            sys.exit(1)

    print("\nMatrix wheel build successfully complete. Check your ./dist/ folder.")

if __name__ == "__main__":
    project_root = Path(__file__).parent.resolve()
    clean_workspace(project_root)
    # build_targeted_wheels(project_root)
    build.compile_all()


