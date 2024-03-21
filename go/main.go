package main

// TODO
// 1. Port the git.ts file to golang
// 2. Port vcs/github to golang
// 3. Devise a way to run scripts in parallel (fan out? semaphor?)
// 4. Use bubbletea to make things pretty. They will probably look something like:
//
// Git
//
// Commits  |||||||||||||||||-------------- 50%
// Authors  |||||||||||||||||||||||||||---- 80%
// Tags     ------------------------------- 0%
// Branches ------------------------------- 0%
//
// GitHub
//
// Issues      |||||||--------------------- 10%
// PRs         ||-------------------------- 3%
// Discussions ||||------------------------ 7%
//
// 5. Interactive script selectors
// 6. Other VCS APIs

func main() {
	// TODO these should come from an args parser
	repo := git{repoUrl: "https://github.com/hashicorp/vagrant.git", dir: "/tmp/foo"}
	// TODO how to make this async?
	repo.Init()
}
