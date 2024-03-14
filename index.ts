#!/usr/bin/env node

import { parseArgs } from 'util';
import { main } from './src/main';

const args = parseArgs({
  allowPositionals: true,
  options: {},
});

await main(args);
