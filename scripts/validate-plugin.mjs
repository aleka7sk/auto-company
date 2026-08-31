#!/usr/bin/env node

import assert from "node:assert/strict";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");

function read(relative) {
  return fs.readFileSync(path.join(root, relative), "utf8");
}

function parseJSON(relative) {
  try {
    return JSON.parse(read(relative));
  } catch (error) {
    throw new Error(`${relative}: invalid JSON: ${error.message}`);
  }
}

function frontmatter(relative) {
  const text = read(relative);
  assert.ok(text.startsWith("---\n"), `${relative}: missing YAML frontmatter`);
  const end = text.indexOf("\n---\n", 4);
  assert.ok(end > 4, `${relative}: unterminated YAML frontmatter`);
  const header = text.slice(4, end);
  const name = header.match(/^name:\s*([^\n]+)$/m)?.[1]?.trim();
  const description = header.match(/^description:\s*([^\n]+)$/m)?.[1]?.trim();
  assert.ok(name, `${relative}: missing name`);
  assert.ok(description, `${relative}: missing description`);
  return { name, description, text };
}

const plugin = parseJSON(".claude-plugin/plugin.json");
assert.equal(plugin.name, "auto-company");
assert.match(plugin.version, /^\d+\.\d+\.\d+$/);
assert.equal(plugin.license, "MIT");

const marketplace = parseJSON(".claude-plugin/marketplace.json");
assert.equal(marketplace.name, "auto-company-marketplace");
assert.ok(Array.isArray(marketplace.plugins));
const marketPlugin = marketplace.plugins.find((item) => item.name === "auto-company");
assert.ok(marketPlugin, "marketplace does not expose auto-company");
assert.equal(marketPlugin.source, "./");

const hooks = parseJSON("hooks/hooks.json");
const preTool = hooks?.hooks?.PreToolUse;
assert.ok(Array.isArray(preTool) && preTool.length > 0, "PreToolUse guard is missing");
const hookCommands = preTool.flatMap((entry) => entry.hooks ?? []).map((entry) => entry.command ?? "");
assert.ok(hookCommands.some((command) => command.includes("guard-tool-use.mjs")));
assert.ok(hookCommands.every((command) => !/curl|wget|invoke-webrequest|npm\s+install/i.test(command)), "hooks must not install or fetch code");

const skillRoot = path.join(root, "skills");
const skillDirs = fs.readdirSync(skillRoot, { withFileTypes: true }).filter((entry) => entry.isDirectory());
assert.ok(skillDirs.length >= 7, "expected at least seven skills");
const skillNames = new Set();
for (const dir of skillDirs) {
  const relative = path.posix.join("skills", dir.name, "SKILL.md");
  assert.ok(fs.existsSync(path.join(root, relative)), `${relative}: missing`);
  const meta = frontmatter(relative);
  assert.ok(!skillNames.has(meta.name), `${relative}: duplicate skill name ${meta.name}`);
  skillNames.add(meta.name);
}
assert.ok(skillNames.has("create-saas"), "orchestrator skill is missing");

const agentRoot = path.join(root, "agents");
const agentFiles = fs.readdirSync(agentRoot).filter((name) => name.endsWith(".md"));
assert.ok(agentFiles.length >= 4, "expected at least four specialist agents");
for (const name of agentFiles) {
  const meta = frontmatter(path.posix.join("agents", name));
  if (meta.name !== "release-reviewer") {
    assert.ok(!/^tools:.*\b(Bash|Write|Edit)\b/m.test(meta.text), `${name}: research/critique agent must be read-only`);
  }
}

const registry = parseJSON("internal/assets/registry/integrations.json");
assert.equal(registry.schemaVersion, "1.0");
assert.ok(Array.isArray(registry.integrations) && registry.integrations.length >= 7);
const ids = new Set();
for (const item of registry.integrations) {
  for (const field of ["id", "name", "repository", "license", "maturity", "role"]) {
    assert.ok(item[field], `registry integration missing ${field}`);
  }
  assert.ok(!ids.has(item.id), `duplicate integration id ${item.id}`);
  ids.add(item.id);
  assert.ok(item.install && typeof item.install === "object", `${item.id}: install guidance missing`);
}

console.log(`plugin validation passed: ${skillDirs.length} skills, ${agentFiles.length} agents, ${registry.integrations.length} integrations`);
