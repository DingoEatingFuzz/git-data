#!/usr/bin/env node

import { parseArgs } from 'util';
import { main } from './src/main';

const args = parseArgs({
  allowPositionals: true,
  options: {},
});

const ghPat = process.env['GH_PAT'];
if (!ghPat) {
  console.error('Expected GH_PAT in env.');
  process.exit(1);
}

await main({ ...args, ghPat });
