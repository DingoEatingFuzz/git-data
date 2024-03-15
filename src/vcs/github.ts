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
    const res = await this.graphql(`
      query allIssues($owner: String!, $repo: String!, $num: Int = 10, $cursor: String) {
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

    return res;
  }

  async allPRs(options: PaginationOptions = {}) {
    const res = await this.graphql(`
      query allPRs($owner: String!, $repo: String!, $num: Int = 10, $cursor: String) {
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

    return res;
  }

  async graphql(query: string, options: GraphqlOptions = {}) {
    const { owner, repo } = this;
    return await this.octokit.graphql.paginate(
      query,
      Object.assign(options, { owner, repo })
    );
  }
}
