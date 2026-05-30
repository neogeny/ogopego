import platform
import subprocess
from pathlib import Path

def get_binary_path() -> Path:
    base_dir = Path(__file__).parent / "binaries"
    system = platform.system().lower() # 'darwin', 'linux', 'windows'
    machine = platform.machine().lower() # 'x86_64', 'arm64', 'amd64'

    # Normalize architecture names
    arch = "amd64" if machine in ["x86_64", "amd64"] else "arm64" if "arm" in machine else machine
    
    ext = ".exe" if system == "windows" else ""
    binary_name = f"ogopego_{system}_{arch}{ext}"
    
    binary_path = base_dir / binary_name
    if not binary_path.exists():
        raise RuntimeError(f"Unsupported platform platform configuration: {system}_{arch}")
        
    return binary_path

def start_go_submodule():
    binary = get_binary_path()
    
    # Example: Fire up the Go binary via a background subprocess
    # Your Go binary can listen on a dynamic port or process stdin/stdout lines
    process = subprocess.Popen(
        [str(binary), "--port", "8080"],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )
    return process
