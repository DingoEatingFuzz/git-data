import { readdir } from 'fs/promises';
import git, { type SimpleGit } from 'simple-git';

interface GitOptions {
  dir: string;
}

export default class Git {
  repoUrl: string;
  dir: string;
  git?: SimpleGit;
  initialized: boolean = false;

  constructor(repoUrl: string, options: GitOptions) {
    this.repoUrl = repoUrl;
    this.dir = options.dir;
  }

  async init() {
    try {
      const files = await readdir(this.dir);
      if (files.length === 0) {
        // Directory exists but it is empty
        await git().clone(this.repoUrl, this.dir);
      }
      this.git = git(this.dir);
      this.initialized = true;
    } catch (err) {
      console.log('Could not clone repo', err);
    }
  }
}
