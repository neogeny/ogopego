import sys
from pathlib import Path

# Ensure the python/ directory is on sys.path so ogopego is importable
_pkg_root = Path(__file__).resolve().parent.parent
if str(_pkg_root) not in sys.path:
    sys.path.insert(0, str(_pkg_root))
