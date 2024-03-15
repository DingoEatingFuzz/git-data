import { createHash } from 'crypto';
import { mkdir, exists } from 'fs/promises';
import { tmpdir } from 'os';
import * as path from 'path';
import Git from './git';
import GitHub from './vcs/github';

interface Args {
  positionals: string[],
  values: { [key:string]: string },
  ghPat: string,
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

  console.log(`Cloning ${args.positionals[0]}...`);
  await repo.init();
  if (repo.initialized) {
    const sh = await repo.git!.show();
    console.log(sh);
  }

  // Authenticate github api
  const parts = githubParts(args.positionals[0]);
  const gh = parts ? new GitHub(parts.owner, parts.repo, { auth: args.ghPat }) : null;

  if (gh) {
    console.log('Fetching all issues');
    console.log(await gh.allIssues());
  }

  // Run scripts
}

function githubParts(repo: string): { owner: string, repo: string } | null {
  const matches = repo.match(/github.com\/(.+?)\/(.+?)(\/|\.git)?$/);
  if (!matches) {
    console.error('Could not parse GitHub URL');
    return null;
  }

  console.log(matches);

  return {
    owner: matches[1] ?? '',
    repo: matches[2] ?? '',
  }
}


async function makeTmpDir(args: Args): Promise<string | undefined> {
  try {
    const sha = checksum(JSON.stringify(args));
    const tmpPath = path.join(tmpdir(), 'git.json', sha);

    await mkdir(tmpPath, { recursive: true });

    return tmpPath;
  } catch (err) {
    console.log('Could not create tmp dir', err);
  }
}
