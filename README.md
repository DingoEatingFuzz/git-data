# git.json

Turn a git or GitHub repo into a collection of json files.

These files include:

  1. Git log
  2. Issue metadata
  3. PR metadata

## Why would I want this?

Git is an invaluable data source but it's trapped in its built-for-purpose blobs and trees structure. Getting all that data in standardized formats makes it much simpler to integrate git and GitHub as data sources to feed into analytic/BI systems.

## Usage

You can use git.json as a standalone tool:

```shellsession
$ npx git.json https://github.com/hashicorp/vagrant
generating files...
```

Or as a package in your own JS/TS scripts:

```js
import GJ from 'git.json'

const repo = await GJ.init('https://github.com/hashicorp/vagrant');
const log = await repo.log();

const jan1 = new Date('01/01/2024');
const yearToDate = log.filter(commit => commit.date > jan1);
```
