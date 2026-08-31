#!/usr/bin/env node

import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import path from "node:path";

const here = path.dirname(fileURLToPath(import.meta.url));
const guard = path.join(here, "guard-tool-use.mjs");

function run(payload) {
  return spawnSync(process.execPath, [guard], {
    input: JSON.stringify(payload),
    encoding: "utf8",
  });
}

function allowed(payload) {
  const result = run(payload);
  assert.equal(result.status, 0, `expected allowed, stderr=${result.stderr}`);
}

function blocked(payload) {
  const result = run(payload);
  assert.equal(result.status, 2, `expected blocked, stderr=${result.stderr}`);
  assert.match(result.stderr, /AUTO-COMPANY BLOCKED/);
}

allowed({ tool_name: "Read", tool_input: { file_path: "src/app.ts" } });
allowed({ tool_name: "Read", tool_input: { file_path: ".env.example" } });
blocked({ tool_name: "Read", tool_input: { file_path: ".env" } });
blocked({ tool_name: "Read", tool_input: { file_path: "C:\\Users\\aleka\\project\\.env.production" } });
blocked({ tool_name: "Write", tool_input: { file_path: "certs/private.pem" } });

allowed({ tool_name: "Bash", tool_input: { command: "git push origin feature/x" } });
blocked({ tool_name: "Bash", tool_input: { command: "git push --force origin feature/x" } });
blocked({ tool_name: "Bash", tool_input: { command: "git push origin main" } });
blocked({ tool_name: "Bash", tool_input: { command: "git push origin HEAD:main" } });
blocked({ tool_name: "Bash", tool_input: { command: "gh pr merge 42 --squash" } });
blocked({ tool_name: "Bash", tool_input: { command: "terraform destroy -auto-approve" } });
blocked({ tool_name: "Bash", tool_input: { command: "Remove-Item C:\\ -Recurse -Force" } });

console.log("guard-tool-use: all tests passed");
