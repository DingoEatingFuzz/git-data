import { Octokit } from "octokit";

interface GraphqlOptions {
  owner?: string;
  repo?: string;
}

interface GitHubOptions {
  auth: string;
}

interface PaginationOptions extends GraphqlOptions {
  num?: number;
}

export default class GitHub {
  owner: string;
  repo: string;
  octokit: Octokit;

  constructor(owner: string, repo: string, options: GitHubOptions) {
    this.owner = owner;
    this.repo = repo;
    this.octokit = new Octokit({ auth: options.auth });
  }

  async allIssues(options: PaginationOptions = {}) {
    const iter = await this.graphql(`
      query allIssues($owner: String!, $repo: String!, $num: Int = 50, $cursor: String) {
        repository(owner: $owner, name: $repo) {
          issues(first: $num, after: $cursor) {
            edges {
              node {
                title,
                author {
                  login,
                },
                participants(first: 100) {
                  nodes {
                    login
                  },
                  totalCount,
                }
                createdAt,
                closedAt,
                closed,
                comments {
                  totalCount
                },
                reactions {
                  totalCount,
                },
                locked,
              }
            }
            pageInfo {
              hasNextPage
              endCursor
            }
          }
        }
      }`,
      options
    );

    let collect = [];
    for await (const res of iter) {
      const issues = res.repository.issues.edges;
      collect.push(...issues);
      console.log(`Fetching issues ${collect.length}`);
    }

    return collect;
  }

  async allPRs(options: PaginationOptions = {}) {
    const iter = await this.graphql(`
      query allPRs($owner: String!, $repo: String!, $num: Int = 50, $cursor: String) {
        repository(owner: $owner, name: $repo) {
          issues(first: $num, after: $cursor) {
            edges {
              node {
                title,
                author {
                  login,
                },
                participants(first: 100) {
                  nodes {
                    login
                  },
                  totalCount,
                }
                createdAt,
                closedAt,
                closed,
                comments {
                  totalCount
                },
                reactions {
                  totalCount,
                },
                locked,
              }
            }
            pageInfo {
              hasNextPage
              endCursor
            }
          }
        }
      }`,
      options
    );

    let collect = [];
    for await (const res of iter) {
      const prs = res.repository.pullRequests.edges;
      collect.push(...prs);
      console.log(`Fetching pull requests ${collect.length}`);
    }

    return collect;
  }

  async graphql(query: string, options: GraphqlOptions = {}) {
    const { owner, repo } = this;
    return this.octokit.graphql.paginate.iterator(
      query,
      Object.assign(options, { owner, repo })
    );
  }
}
