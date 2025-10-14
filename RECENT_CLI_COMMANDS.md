# Recent CLI Commands Used for Project Tasks

This document lists the terminal commands used in the last 5 major tasks. Some of these are excellent candidates for custom shortcuts, scripts, or integration as tools in Neovim and VS Code.

---

## 1. Check Current Git Status and Branch

```sh
cd /home/alex/SHDProj/SSTorytime && git status
```

- Shows current branch, untracked files, and sync status with remote.

---

## 2. Compare Changes Between Branches

```sh
cd /home/alex/SHDProj/SSTorytime && git diff --name-status main..text2n4l
```

- Lists all files added, modified, or deleted between `main` and `text2n4l` branches.

---

## 3. Run Integration Tests in n4l

```sh
cd /home/alex/SHDProj/SSTorytime/n4l/tests && timeout 30 ../run_tests 2>/dev/null | tail -5
```

- Runs the integration test suite for the n4l package and shows the last 5 lines of output.

---

## 4. Rename a File in the Project Root

```sh
mv PROJECT_SUMMARY_FOR_DR_BURGESS.md PROJECT_SUMMARY.md
```

- Renames the summary document to remove personal references.

---

## 5. Check for Uncommitted Changes

```sh
git status
```

- Shows the current state of the working directory and staging area.

---

## Suggestions for Shortcuts/Tools

- `git diff --name-status main..text2n4l` as a custom command or keybinding for quick branch comparison
- Integration test runner as a VS Code/Neovim task
- File renaming and project cleanup scripts
- Automated test output tailing for fast feedback
