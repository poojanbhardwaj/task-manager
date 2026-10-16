import { existsSync, realpathSync, rmSync } from "node:fs";
import { dirname, isAbsolute, relative, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const projectRoot = realpathSync(resolve(dirname(fileURLToPath(import.meta.url)), ".."));
const workspaces = [
  "apps/frontend",
  "packages/zod",
  "packages/openapi",
  "packages/emails",
];
const scope = process.argv[2] ?? "all";

if (process.argv.length > 3 || (scope !== "all" && !workspaces.includes(scope))) {
  throw new Error("Choose all or a configured frontend/package workspace.");
}

const scopes = scope === "all" ? [...workspaces, "."] : [scope];
const artifacts = ["dist", "tsconfig.tsbuildinfo", ".turbo", "node_modules"];
const targets = [];

function assertInsideProject(path) {
  const child = relative(projectRoot, path);
  if (!child || child === ".." || child.startsWith("../") ||
      child.startsWith("..\\") || isAbsolute(child)) {
    throw new Error("Refusing to clean outside the project: " + path);
  }
}

// Validate every workspace before deleting anything. A linked workspace must
// not point outside this project. rmSync removes artifact links themselves.
for (const workspace of scopes) {
  const directory = resolve(projectRoot, workspace);
  if (!existsSync(directory)) continue;
  if (workspace !== ".") assertInsideProject(realpathSync(directory));
  for (const artifact of artifacts) {
    const target = resolve(directory, artifact);
    assertInsideProject(target);
    targets.push(target);
  }
}

for (const target of targets) {
  rmSync(target, { recursive: true, force: true, maxRetries: 3, retryDelay: 100 });
}

console.log("Cleaned generated files and dependencies for " + scope + ".");
