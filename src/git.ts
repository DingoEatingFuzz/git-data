import git from 'simple-git';

interface GitOptions {
  dir: string;
}

export default class Git {
  repoUrl: string;
  dir: string;

  constructor(repoUrl: string, options: GitOptions) {
    this.repoUrl = repoUrl;
    this.dir = options.dir;
  }

  async init() {
    // Clone the repo
  }
}
