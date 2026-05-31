import subprocess
from pathlib import Path
from hatchling.build.hooks.plugin.interface import BuildHookInterface

class CustomBuildHook(BuildHookInterface):
    PLUGIN_NAME = "custom"

    def initialize(self, version, build_data):
        go_src_dir = Path(self.root) / "go_src"
        
        print("\n--- Hatch Build Hook: Invoking locked project gopy tool ---")
        
        # 'go tool gopy' uses the specific, patched version in go.mod
        subprocess.run(
            ["go", "tool", "gopy", "pkg", "-vm=python3", "./geometry"],
            cwd=str(go_src_dir),
            check=True
        )
        
        # Then proceed with mapping files into your wheel artifacts...
