import { createHash, randomBytes } from 'node:crypto';
import { spawn } from 'node:child_process';
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = dirname(dirname(fileURLToPath(import.meta.url)));
const backend = join(root, 'apps', 'backend');
const envPath = join(backend, '.env');
const version = '4.48';
const expectedHash = 'fe90c04c0620ad1a1c756f86cd5e1443773f56a446688077a4bb5ec04c3cc874';
const binDir = join(root, '.local', 'seaweedfs', version);
const binary = join(binDir, 'weed.exe');
const archive = join(binDir, 'windows_amd64.zip');
const dataDir = join(root, '.local', 'storage-data');
const prefix = 'TASK_MANAGER_STORAGE.';

function run(command, args, options = {}) {
  return new Promise((resolve, reject) => {
    const child = spawn(command, args, { cwd: root, stdio: 'inherit', windowsHide: true, ...options });
    child.on('error', reject);
    child.on('exit', (code, signal) => code === 0 ? resolve() : reject(new Error(`${command} exited with ${signal ?? code}`)));
  });
}

function readStorage() {
  const text = existsSync(envPath) ? readFileSync(envPath, 'utf8') : '';
  const values = {};
  for (const line of text.split(/\r?\n/)) {
    const match = line.match(/^\s*TASK_MANAGER_STORAGE\.([A-Z_]+)\s*=\s*(.*?)\s*$/);
    if (!match) continue;
    let value = match[2];
    if (value.startsWith('"')) value = JSON.parse(value);
    else if (value.startsWith("'") && value.endsWith("'")) value = value.slice(1, -1);
    else value = value.split(/\s+#/)[0].trim();
    values[match[1]] = value;
  }
  return { text, values };
}

async function setup() {
  if (process.platform !== 'win32' || process.arch !== 'x64') throw new Error('This installer supports Windows x64. See docs/local-storage.md for other systems.');
  mkdirSync(binDir, { recursive: true });
  if (!existsSync(binary)) {
    if (!existsSync(archive) || createHash('sha256').update(readFileSync(archive)).digest('hex') !== expectedHash) {
      await run('curl.exe', ['--fail', '--location', '--connect-timeout', '20', '--max-time', '300', '--output', archive,
        `https://github.com/seaweedfs/seaweedfs/releases/download/${version}/windows_amd64.zip`]);
    }
    if (createHash('sha256').update(readFileSync(archive)).digest('hex') !== expectedHash) throw new Error('SeaweedFS checksum mismatch; download was not installed.');
    await run('powershell.exe', ['-NoProfile', '-NonInteractive', '-Command',
      'Expand-Archive -LiteralPath $env:TASK_STORAGE_ARCHIVE -DestinationPath $env:TASK_STORAGE_BIN -Force'],
      { env: { ...process.env, TASK_STORAGE_ARCHIVE: archive, TASK_STORAGE_BIN: binDir } });
  }
  const { text, values } = readStorage();
  const defaults = {
    ENDPOINT: 'http://127.0.0.1:8333', REGION: 'us-east-1', BUCKET: 'task-manager',
    ACCESS_KEY: `local-${randomBytes(12).toString('hex')}`,
    SECRET_KEY: randomBytes(32).toString('hex'), USE_PATH_STYLE: 'true',
  };
  const missing = Object.entries(defaults).filter(([key]) => !(key in values));
  if (missing.length) {
    writeFileSync(envPath, text + '\n# Local object storage (bun storage:start)\n' +
      missing.map(([key, value]) => `${prefix}${key}=${JSON.stringify(value)}`).join('\n') + '\n');
  }
  mkdirSync(dataDir, { recursive: true });
  console.log('Local storage installed. Storage settings saved in apps/backend/.env; existing settings preserved.');
  console.log('Run bun storage:start, then bun storage:check in a second terminal.');
}

async function start() {
  if (!existsSync(binary)) throw new Error('Run bun storage:setup first.');
  const { values } = readStorage();
  for (const key of ['ENDPOINT', 'REGION', 'BUCKET', 'ACCESS_KEY', 'SECRET_KEY', 'USE_PATH_STYLE']) {
    if (!values[key]) throw new Error(`Missing ${prefix}${key}; run bun storage:setup or check apps/backend/.env.`);
  }
  if (!['http://127.0.0.1:8333', 'http://localhost:8333'].includes(values.ENDPOINT.replace(/\/$/, '')) || values.USE_PATH_STYLE !== 'true') {
    throw new Error('storage:start is for local SeaweedFS at port 8333 with USE_PATH_STYLE=true. Your configured storage endpoint was preserved.');
  }
  mkdirSync(dataDir, { recursive: true });
  console.log(`Starting local storage at ${values.ENDPOINT}; bucket: ${values.BUCKET}. Ctrl+C stops it; files are retained.`);
  // Pass only storage credentials to SeaweedFS, never Clerk or other .env secrets.
  await run(binary, ['mini', `-dir=${dataDir}`, '-ip=127.0.0.1', '-ip.bind=127.0.0.1',
    '-admin.ui=false', '-webdav=false', '-s3.port.iceberg=0', '-s3.port.lance=0',
    '-s3.allowedOrigins=http://localhost:3000', '-s3.autoCreateBucket=false', '-s3.allowDeleteBucketNotEmpty=false'],
    { env: { ...process.env, AWS_ACCESS_KEY_ID: values.ACCESS_KEY, AWS_SECRET_ACCESS_KEY: values.SECRET_KEY, S3_BUCKET: values.BUCKET } });
}

try {
  switch (process.argv[2]) {
    case 'setup': await setup(); break;
    case 'start': await start(); break;
    case 'check': await run('go', ['run', './cmd/storage-check'], { cwd: backend }); break;
    default: throw new Error('Use bun storage:setup, bun storage:start, or bun storage:check.');
  }
} catch (error) {
  console.error(error.message);
  process.exitCode = 1;
}
