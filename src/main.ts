import { createHash } from 'crypto';
import { mkdir, exists } from 'fs/promises';
import { tmpdir } from 'os';
import * as path from 'path';
import Git from './git';

interface Args {
  positionals: string[],
  values: { [key:string]: string },
}

const checksum = (digest: string): string => {
  return createHash('sha256')
    .update(digest)
    .digest()
    .toString('hex')
}

export async function main(args: Args) {
  // Make tmp dir
  const dir = await makeTmpDir(args);
  if (!dir) {
    process.exit(1);
  }

  // Clone the repo
  const repo = new Git(args.positionals[0], { dir })

  console.log('Cloning repo...');
  await repo.init();
  if (repo.initialized) {
    const sh = await repo.git!.show();
    console.log(sh);
  }

  // Authenticate github api
  // Run scripts
}


async function makeTmpDir(args: Args): Promise<string | undefined> {
  try {
    const sha = checksum(JSON.stringify(args));
    const tmpPath = path.join(tmpdir(), 'git.json', sha);

    if (!exists(tmpPath)) {
      await mkdir(tmpPath, { recursive: true });
    }

    return tmpPath;
  } catch (err) {
    console.log('Could not create tmp dir', err);
  }
}
