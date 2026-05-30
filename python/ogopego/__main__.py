import sys
from ._ogopego import ogo


def main():
    try:
        # sys.argv[1:] passes all command line flags down into the driver, ignoring the calling script path
        ogo(*sys.argv[1:])
    except FileNotFoundError as e:
        print(f"Initialization Error: {e}", file=sys.stderr)
        sys.exit(1)
    except KeyboardInterrupt:
        sys.exit(130)  # Standard Linux exit code for SIGINT


if __name__ == "__main__":
    main()
