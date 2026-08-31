#!/usr/bin/env node

import process from "node:process";

const chunks = [];
for await (const chunk of process.stdin) chunks.push(chunk);

let event;
try {
  event = JSON.parse(Buffer.concat(chunks).toString("utf8") || "{}");
} catch (error) {
  console.error(`Auto Company guard could not parse hook input: ${error.message}`);
  process.exit(2);
}

const toolName = String(event.tool_name ?? event.toolName ?? "");
const input = event.tool_input ?? event.toolInput ?? {};
const candidatePath = String(input.file_path ?? input.path ?? "").replaceAll("\\", "/");
const lowerPath = candidatePath.toLowerCase();

const allowedSecretExamples = [".env.example", ".env.sample", ".env.template"];
const isAllowedExample = allowedSecretExamples.some((suffix) => lowerPath.endsWith(suffix));
const secretPathPatterns = [
  /(^|\/)\.env($|\.)/,
  /(^|\/)id_rsa($|\.)/,
  /(^|\/)id_ed25519($|\.)/,
  /\.pem$/,
  /\.p12$/,
  /\.pfx$/,
  /(^|\/)credentials($|\.)/,
  /(^|\/)secrets?($|\/)/,
  /(^|\/)\.aws\/credentials$/,
  /(^|\/)\.kube\/config$/,
];

if (["Read", "Write", "Edit"].includes(toolName) && candidatePath && !isAllowedExample) {
  const matched = secretPathPatterns.some((pattern) => pattern.test(lowerPath));
  if (matched) {
    console.error(
      `AUTO-COMPANY BLOCKED: access to possible secret file '${candidatePath}'. ` +
      "Use a sanitized example file or have the owner perform the operation manually."
    );
    process.exit(2);
  }
}

if (toolName === "Bash") {
  const command = String(input.command ?? "");
  const normalized = command.toLowerCase().replaceAll("\\", "/").replace(/\s+/g, " ").trim();
  const blocked = [
    { pattern: /git\s+push[^\n]*(--force|-f)(\s|$)/, reason: "force push" },
    { pattern: /git\s+push[^\n]*(?:\s|:)(?:refs\/heads\/)?(?:main|master)(?:\s|$)/, reason: "direct push to main/master" },
    { pattern: /git\s+reset\s+--hard/, reason: "destructive git reset" },
    { pattern: /git\s+clean\s+[^\n]*-[^\n]*f[^\n]*d/, reason: "destructive git clean" },
    { pattern: /gh\s+pr\s+merge/, reason: "pull request merge" },
    { pattern: /terraform\s+destroy/, reason: "infrastructure destruction" },
    { pattern: /kubectl\s+delete\s+(namespace|ns)/, reason: "namespace deletion" },
    { pattern: /helm\s+uninstall/, reason: "release uninstall" },
    { pattern: /drop\s+(database|schema)/, reason: "database/schema destruction" },
    { pattern: /npm\s+publish/, reason: "package publication" },
    { pattern: /eas\s+submit/, reason: "app-store submission" },
    { pattern: /vercel\s+[^\n]*--prod/, reason: "production deployment" },
    { pattern: /rm\s+-rf\s+(\/|~|\$home)(\s|$)/, reason: "broad filesystem deletion" },
    { pattern: /remove-item\s+[^\n]*c:\/[^\n]*-recurse[^\n]*-force/, reason: "broad Windows filesystem deletion" },
  ];
  const match = blocked.find(({ pattern }) => pattern.test(normalized));
  if (match) {
    console.error(
      `AUTO-COMPANY BLOCKED: ${match.reason}. ` +
      "This action requires explicit owner approval and should be executed manually or after temporarily disabling the guard."
    );
    process.exit(2);
  }
}

process.exit(0);
