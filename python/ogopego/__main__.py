import sys
from .build import get_binary_path
from .main import start_backend_service

def main():
    # Ask the builder toolchain exactly where the local or host binary path is
    binary_path = get_binary_path()

    if not binary_path.exists():
        print(f"Error: Native binary component is missing at target location: {binary_path}", file=sys.stderr)
        sys.exit(1)

    print(f"Starting ogopego backend using: {binary_path.name}...")

    process = start_backend_service(binary_path)

    try:
        process.wait()
    except KeyboardInterrupt:
        process.terminate()
        print("\nBackend service stopped.")

if __name__ == "__main__":
    main()
