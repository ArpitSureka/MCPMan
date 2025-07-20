import subprocess
import tempfile
import os
import shutil
import sys


# List of packages to be downloaded and bundled
PACKAGES_TO_BUNDLE = ["mcp-server-fetch", "mcp-server-git", "mcp-server-time"]

# Output directory for the executables
OUTPUT_DIRECTORY = "scripts/PythonMCPExecutables"

# Temporary build directory for PyInstaller
BUILD_DIRECTORY_TEMP = "temp_build"

def run_command(command, check=True, capture_output=False, text=False):
    """
    Helper function to run a shell command.
    Prints the command and its output.
    Exits if the command fails and check is True.
    """
    print(f"\nRunning command: {' '.join(command)}")
    try:
        process = subprocess.run(
            command,
            check=check,
            capture_output=capture_output,
            text=text,
            # Set a timeout to prevent indefinite hanging
            timeout=300 # 5 minutes, adjust as needed
        )
        if capture_output:
            if process.stdout:
                print("Output:\n", process.stdout)
            if process.stderr:
                print("Error output:\n", process.stderr)
        return process
    except subprocess.CalledProcessError as e:
        print(f"Error executing command: {' '.join(command)}")
        print(f"Return code: {e.returncode}")
        if e.stdout:
            print(f"Stdout: {e.stdout}")
        if e.stderr:
            print(f"Stderr: {e.stderr}")
        if check:
            sys.exit(f"Command failed: {' '.join(command)}")
        return None
    except FileNotFoundError:
        print(f"Error: Command '{command[0]}' not found. Is it installed and in your PATH?")
        if check:
            sys.exit(f"Command not found: {command[0]}")
        return None
    except subprocess.TimeoutExpired:
        print(f"Error: Command {' '.join(command)} timed out.")
        if check:
            sys.exit(f"Command timed out: {' '.join(command)}")
        return None


def check_and_install_pyinstaller():
    """Checks if PyInstaller is installed, and if not, offers to install it."""
    print("Checking for PyInstaller...")
    try:
        # Try to import PyInstaller to check if it's available
        import PyInstaller
        print("PyInstaller is already installed.")
    except ImportError:
        print("PyInstaller is not found.")
        choice = input("Do you want to try installing PyInstaller now? (yes/no): ").lower()
        if choice == 'yes':
            print("Installing PyInstaller...")
            run_command([sys.executable, "-m", "pip", "install", "pyinstaller"])
            print("PyInstaller installation attempted. Please re-run the script if it was successful.")
            sys.exit(0) # Exit so user can re-run with PyInstaller available
        else:
            print("PyInstaller is required to create executables. Please install it manually (`pip install pyinstaller`) and try again.")
            sys.exit(1)

def find_script_path(script_name):
    """
    Tries to find the path of an installed script.
    Relies on the script being in the system's PATH, which pip usually handles.
    """
    print(f"Searching for script: {script_name}")
    script_path = shutil.which(script_name)
    if script_path:
        print(f"Found script at: {script_path}")
        return script_path
    else:
        # Fallback for Windows: check with .exe extension if not already there
        if sys.platform == "win32" and not script_name.endswith(".exe"):
            script_path_exe = shutil.which(script_name + ".exe")
            if script_path_exe:
                print(f"Found script at: {script_path_exe}")
                return script_path_exe
        
        # Fallback for scripts installed in user's local bin
        if sys.platform != "win32":
            local_bin_path = os.path.expanduser("~/.local/bin")
            potential_path = os.path.join(local_bin_path, script_name)
            if os.path.exists(potential_path) and os.access(potential_path, os.X_OK):
                 print(f"Found script at: {potential_path} (in user's local bin)")
                 return potential_path

        print(f"Warning: Script '{script_name}' not found in PATH.")
        print("This might happen if the package doesn't install a script with this exact name,")
        print("or if the Python environment's Scripts/bin directory is not in your PATH.")
        return None

def create_executable(script_path, package_name):
    """
    Creates a standalone executable using PyInstaller.
    """
    if not script_path:
        print(f"Skipping executable creation for {package_name} as its script was not found.")
        return False

    executable_name = package_name # Use package name as executable name
    print(f"\nCreating executable for: {package_name} (from {script_path})")
    print(f"Output executable name will be: {executable_name}")

    # PyInstaller command
    # --onefile: Create a single executable file
    # --name: Name of the executable
    # --distpath: Directory to put the final executable
    # --workpath: Directory to put temporary build files
    # --specpath: Directory to put the .spec file
    # --clean: Clean PyInstaller cache and remove temporary files before building
    # --noconfirm: Overwrite output directory without asking
    pyinstaller_command = [
        sys.executable, "-m", "PyInstaller",
        "--onefile",
        "--name", executable_name,
        "--distpath", OUTPUT_DIRECTORY,
        "--workpath", BUILD_DIRECTORY_TEMP,
        "--specpath", os.getcwd(), # Put .spec file in current dir
        "--clean",
        "--noconfirm", # Suppress y/n questions from PyInstaller
        script_path
    ]

    try:
        run_command(pyinstaller_command)
        print(f"Successfully created executable for {package_name} in '{OUTPUT_DIRECTORY}' directory.")
        
        # Clean up the .spec file for this package
        spec_file = f"{executable_name}.spec"
        if os.path.exists(spec_file):
            try:
                os.remove(spec_file)
                print(f"Removed spec file: {spec_file}")
            except OSError as e:
                print(f"Warning: Could not remove spec file {spec_file}: {e}")
        
        return True
    except Exception as e:
        print(f"Error creating executable for {package_name}: {e}")
        return False

def main():
    """Main function to orchestrate the process."""
    print("--- MCP Server Executable Bundler ---")

    check_and_install_pyinstaller()

    if not os.path.exists(OUTPUT_DIRECTORY):
        print(f"Creating output directory: {OUTPUT_DIRECTORY}")
        os.makedirs(OUTPUT_DIRECTORY)
    
    # Clean up temporary build directory from previous runs if it exists
    if os.path.exists(BUILD_DIRECTORY_TEMP):
        print(f"Removing old temporary build directory: {BUILD_DIRECTORY_TEMP}")
        shutil.rmtree(BUILD_DIRECTORY_TEMP, ignore_errors=True)


    successful_bundles = []
    failed_bundles = []

    for package_name in PACKAGES_TO_BUNDLE:
        print(f"\n--- Processing package: {package_name} ---")

        # Step 1: Install/Upgrade the package using pip
        print(f"Installing/upgrading package: {package_name}...")
        pip_command = [sys.executable, "-m", "pip", "install", "--upgrade", package_name]
        install_process = run_command(pip_command, check=False) # Don't exit immediately if pip fails

        if install_process and install_process.returncode != 0:
            print(f"Failed to install package {package_name}. Skipping bundling for this package.")
            failed_bundles.append(package_name)
            continue
        
        print(f"Package {package_name} installed/updated successfully.")

        # Step 2: Find the script path
        # We assume the script name is the same as the package name
        script_to_bundle = find_script_path(package_name)

        if not script_to_bundle:
            print(f"Could not find an executable script for {package_name}. It might not install a CLI tool with that name, or it's not in PATH.")
            failed_bundles.append(package_name)
            continue
            
        # Step 3: Create the executable
        if create_executable(script_to_bundle, package_name):
            successful_bundles.append(package_name)
        else:
            failed_bundles.append(package_name)
        
        # Clean up temporary build directory after each PyInstaller run to save space
        if os.path.exists(BUILD_DIRECTORY_TEMP):
            shutil.rmtree(BUILD_DIRECTORY_TEMP, ignore_errors=True)


    print("\n--- Bundling Summary ---")
    if successful_bundles:
        print("Successfully created executables for:")
        for pkg in successful_bundles:
            exe_path = os.path.join(OUTPUT_DIRECTORY, pkg)
            if sys.platform == "win32" and not exe_path.endswith(".exe"):
                exe_path += ".exe"
            print(f"  - {pkg} (as {exe_path})")
    if failed_bundles:
        print("\nFailed to create executables or find scripts for:")
        for pkg in failed_bundles:
            print(f"  - {pkg}")
    
    if not successful_bundles and not failed_bundles:
        print("No packages were processed.")
    elif not successful_bundles and failed_bundles:
        print("\nNo executables were successfully created.")
    else:
        print(f"\nAll created executables are located in the '{OUTPUT_DIRECTORY}' directory.")

    # Final cleanup of temporary build directory if it somehow still exists
    if os.path.exists(BUILD_DIRECTORY_TEMP):
        shutil.rmtree(BUILD_DIRECTORY_TEMP, ignore_errors=True)
    print("-------------------------")

if __name__ == "__main__":
    main()
