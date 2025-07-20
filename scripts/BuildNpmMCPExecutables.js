#!/usr/bin/env node
const fs = require('fs').promises;
const path = require('path');
const { exec } = require('child_process');
const util = require('util');

// Promisify the exec function to use it with async/await
const execPromise = util.promisify(exec);

// --- Configuration ---
const PACKAGES_TO_BUNDLE = [
  '@modelcontextprotocol/server-filesystem',
  '@modelcontextprotocol/server-memory',
  '@modelcontextprotocol/server-sequential-thinking'
];

// Output directory for the final executables
const OUTPUT_DIRECTORY = 'scripts/NpmMCPExecutables';

// Specify target platforms for pkg. You can customize this.
// Format: node<version>-<platform>-<arch>
// See `pkg --help` for all options.
// const PKG_TARGETS = 'node18-linux-x64,node18-win-x64,node18-macos-x64';
const PKG_TARGETS = 'node18-macos-x64';
// --- End Configuration ---


/**
 * Helper function to run a shell command and stream its output.
 * @param {string} command The command to execute.
 * @returns {Promise<{stdout: string, stderr: string}>}
 */
async function runCommand(command) {
  console.log(`\n$ ${command}`);
  try {
    const { stdout, stderr } = await execPromise(command);
    if (stdout) console.log(stdout);
    if (stderr) console.error(stderr);
    return { stdout, stderr };
  } catch (error) {
    console.error(`Error executing command: ${command}`);
    console.error(`Error: ${error.message}`);
    // Re-throw the error to be caught by the calling function
    throw error;
  }
}

/**
 * Checks if pkg is installed globally.
 */
async function checkPkgInstallation() {
  console.log('Checking for `pkg` installation...');
  try {
    await runCommand('pkg --version');
    console.log('`pkg` is installed.');
    return true;
  } catch (error) {
    console.error('`pkg` command not found.');
    console.error('Please install it globally to continue:');
    console.error('npm install -g pkg');
    return false;
  }
}

/**
 * Main function to orchestrate the bundling process.
 */
async function main() {
  console.log('--- NPM Executable Bundler ---');

  if (!(await checkPkgInstallation())) {
    process.exit(1);
  }

  // Create output directory if it doesn't exist
  try {
    await fs.mkdir(OUTPUT_DIRECTORY, { recursive: true });
  } catch (error) {
    console.error(`Failed to create output directory: ${OUTPUT_DIRECTORY}`);
    process.exit(1);
  }

  const successfulBundles = [];
  const failedBundles = [];

  for (const pkgName of PACKAGES_TO_BUNDLE) {
    console.log(`\n--- Processing package: ${pkgName} ---`);
    try {
      // Step 1: Install the package to get its files
      console.log(`Installing ${pkgName} to determine its entry point...`);
      await runCommand(`npm install ${pkgName}`);

      // Step 2: Find the package's entry point from its package.json
      const pkgJsonPath = path.resolve(process.cwd(), 'node_modules', pkgName, 'package.json');
      const pkgJsonContent = await fs.readFile(pkgJsonPath, 'utf8');
      const pkgJson = JSON.parse(pkgJsonContent);

      if (!pkgJson.bin) {
        throw new Error(`Package ${pkgName} does not have a 'bin' field in its package.json. Cannot determine entry point.`);
      }

      // The 'bin' field can be a string or an object. We'll handle both.
      const binEntries = typeof pkgJson.bin === 'string' ? { [pkgJson.name]: pkgJson.bin } : pkgJson.bin;
      const executableNames = Object.keys(binEntries);
      
      if (executableNames.length === 0) {
          throw new Error(`The 'bin' field for ${pkgName} is empty.`);
      }

      // We'll process the first binary found. For packages with multiple binaries, you might need to adjust this.
      const firstExecName = executableNames[0];
      const entryScriptRelativePath = binEntries[firstExecName];
      const entryScriptFullPath = path.resolve(path.dirname(pkgJsonPath), entryScriptRelativePath);

      console.log(`Found entry point for '${firstExecName}': ${entryScriptFullPath}`);

      // Step 3: Run pkg to create the executable
      // The output path will be inside our output directory.
      // We use path.join to handle platform-specific separators.
      const outputPath = path.join(OUTPUT_DIRECTORY, firstExecName);
      
      console.log(`Bundling with pkg... This might take a while.`);
      const pkgCommand = `pkg --targets ${PKG_TARGETS} --output ${outputPath} ${entryScriptFullPath}`;
      await runCommand(pkgCommand);

      console.log(`Successfully created executables for ${pkgName}.`);
      successfulBundles.push(pkgName);

    } catch (error) {
      console.error(`\nFailed to bundle package ${pkgName}.`);
      console.error(error.message);
      failedBundles.push(pkgName);
    }
  }

  // Final Summary
  console.log('\n--- Bundling Summary ---');
  if (successfulBundles.length > 0) {
    console.log('✅ Successfully created executables for:');
    successfulBundles.forEach(pkg => console.log(`  - ${pkg}`));
    console.log(`\nFind them in the '${OUTPUT_DIRECTORY}' directory.`);
    console.log(`Platforms targeted: ${PKG_TARGETS}`);
  }
  if (failedBundles.length > 0) {
    console.log('\n❌ Failed to create executables for:');
    failedBundles.forEach(pkg => console.log(`  - ${pkg}`));
  }
  console.log('--------------------------');
}

main().catch(error => {
  console.error('\nAn unexpected error occurred:', error);
  process.exit(1);
});
