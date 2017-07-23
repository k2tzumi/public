<?php
# Copyright (c) 2015, phpfmt and its authors
# All rights reserved.
#
# Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
#
# 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
#
# 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
#
# 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

function showHelp(array $argv, bool $enableCache, bool $inPhar) {
	echo 'Usage: ' . $argv[0] . ' [-hv] [-o=FILENAME] [--config=FILENAME] ' . ($enableCache ? '[--cache[=FILENAME]] ' : '') . '[options] <target>', PHP_EOL;
	$options = [
		'--cache[=FILENAME]' => 'cache file. Default: ',
		'--config=FILENAME' => 'configuration file. Default: .phpfmt.ini',
		'--dry-run' => 'Runs the formatter without atually changing files; returns exit code 1 if changes would have been applied',
		'--ignore=PATTERN-1,PATTERN-N,...' => 'ignore file names whose names contain any PATTERN-N',
		'--lint-before' => 'lint files before pretty printing (PHP must be declared in %PATH%/$PATH)',
		'--no-backup' => 'no backup file (original.php~)',
		'-h, --help' => 'this help message',
		'-o=file' => 'output the formatted code to "file"',
		'-o=-' => 'output the formatted code to standard output',
		'-v' => 'verbose',
	];
	if ($inPhar) {
		$options['--selfupdate[=channel]'] = 'self-update fmt.phar from Github. options: lts, alpha (default: lts)';
		$options['--version'] = 'version';
	}
	$options['--cache[=FILENAME]'] .= (Cacher::DEFAULT_CACHE_FILENAME);
	if (!$enableCache) {
		unset($options['--cache[=FILENAME]']);
	}
	ksort($options);
	$maxLen = max(array_map(function ($v) {
		return strlen($v);
	}, array_keys($options)));
	foreach ($options as $k => $v) {
		echo '  ', str_pad($k, $maxLen), '  ', $v, PHP_EOL;
	}

	echo PHP_EOL, 'If <target> is "-", it reads from stdin', PHP_EOL;
	echo PHP_EOL, 'phpfmt is free for educational and personal use only.', PHP_EOL, 'Thanks for using it fairly.', PHP_EOL, 'Check license at http://phpfmt.github.io/fmt/', PHP_EOL, PHP_EOL;
}

$getoptLongOptions = [
	'cache::',
	'config:',
	'dry-run',
	'help',
	'ignore:',
	'lint-before',
	'no-backup',
	'oracleDB::',
	'profile:',
];
if ($inPhar) {
	$getoptLongOptions[] = 'selfupdate::';
	$getoptLongOptions[] = 'version';
}
if (!$enableCache) {
	unset($getoptLongOptions['cache::']);
}
$opts = getopt(
	'ihvo:',
	$getoptLongOptions
);

if (isset($opts['selfupdate'])) {
	selfupdate($argv, $inPhar, $opts['selfupdate']);
}
if (isset($opts['version'])) {
	if ($inPhar) {
		echo $argv[0], ' ', VERSION, PHP_EOL;
	}
	exit(0);
}
if (isset($opts['config'])) {
	$argv = extractFromArgv($argv, 'config');

	if ('scan' == $opts['config']) {
		$cfgfn = getcwd() . DIRECTORY_SEPARATOR . '.phpfmt.ini';
		$lastcfgfn = '';
		fwrite(STDERR, 'Scanning for configuration file...');
		while (!is_file($cfgfn) && $lastcfgfn != $cfgfn) {
			$lastcfgfn = $cfgfn;
			$cfgfn = dirname(dirname($cfgfn)) . DIRECTORY_SEPARATOR . '.phpfmt.ini';
		}
		$opts['config'] = $cfgfn;
		if (file_exists($opts['config']) && is_file($opts['config'])) {
			fwrite(STDERR, $opts['config']);
			$iniOpts = parse_ini_file($opts['config'], true);
			if (!empty($iniOpts)) {
				$opts += $iniOpts;
			}
		}
		fwrite(STDERR, PHP_EOL);
	} else {
		if (!file_exists($opts['config']) || !is_file($opts['config'])) {
			fwrite(STDERR, 'Custom configuration not file found' . PHP_EOL);
			exit(255);
		}
		$iniOpts = parse_ini_file($opts['config'], true);
		if (!empty($iniOpts)) {
			$opts += $iniOpts;
		}
	}
} elseif (file_exists(getcwd() . DIRECTORY_SEPARATOR . '.phpfmt.ini') && is_file(getcwd() . DIRECTORY_SEPARATOR . '.phpfmt.ini')) {
	fwrite(STDERR, 'Configuration file found' . PHP_EOL);
	$iniOpts = parse_ini_file(getcwd() . DIRECTORY_SEPARATOR . '.phpfmt.ini', true);
	if (isset($opts['profile'])) {
		$argv = extractFromArgv($argv, 'profile');
		$profile = &$iniOpts[$opts['profile']];
		if (isset($profile)) {
			$iniOpts = $profile;
		}
	}
	$opts = array_merge($iniOpts, $opts);
}
if (isset($opts['h']) || isset($opts['help'])) {
	showHelp($argv, $enableCache, $inPhar);
	exit(0);
}

$cache = new CacheDummy();
$cache_fn = '';
if ($enableCache && isset($opts['cache'])) {
	$argv = extractFromArgv($argv, 'cache');
	$cache_fn = $opts['cache'];
	$cache = new Cache($cache_fn);
	fwrite(STDERR, 'Using cache ...' . PHP_EOL);
}

$backup = true;
if (isset($opts['no-backup'])) {
	$argv = extractFromArgv($argv, 'no-backup');
	$backup = false;
}

$dryRun = false;
if (isset($opts['dry-run'])) {
	$argv = extractFromArgv($argv, 'dry-run');
	$dryRun = true;
}

$ignore_list = null;
if (isset($opts['ignore'])) {
	$argv = extractFromArgv($argv, 'ignore');
	$ignore_list = array_map(function ($v) {
		return trim($v);
	}, explode(',', $opts['ignore']));
}

$lintBefore = false;
if (isset($opts['lint-before'])) {
	$argv = extractFromArgv($argv, 'lint-before');
	$lintBefore = true;
}

$fmt = new CodeFormatter();

if (isset($opts['oracleDB'])) {
	$argv = extractFromArgv($argv, 'oracleDB');

	if ('scan' == $opts['oracleDB']) {
		$oracle = getcwd() . DIRECTORY_SEPARATOR . 'oracle.sqlite';
		$lastoracle = '';
		while (!is_file($oracle) && $lastoracle != $oracle) {
			$lastoracle = $oracle;
			$oracle = dirname(dirname($oracle)) . DIRECTORY_SEPARATOR . 'oracle.sqlite';
		}
		$opts['oracleDB'] = $oracle;
		fwrite(STDERR, PHP_EOL);
	}

	if (file_exists($opts['oracleDB']) && is_file($opts['oracleDB'])) {
		$fmt->enablePass('AutoImportPass', $opts['oracleDB']);
	}
}

if (isset($opts['v'])) {
	$argv = extractFromArgvShort($argv, 'v');
	fwrite(STDERR, 'Used passes: ' . implode(', ', $fmt->getPassesNames()) . PHP_EOL);
}

if (isset($opts['i'])) {
	echo 'php.tools fmt.php interactive mode.', PHP_EOL;
	echo 'no <?php is necessary', PHP_EOL;
	echo 'type a lone "." to finish input.', PHP_EOL;
	echo 'type "quit" to finish.', PHP_EOL;
	while (true) {
		$str = '';
		do {
			$line = readline('> ');
			$str .= $line;
		} while (!('.' == $line || 'quit' == $line));
		if ('quit' == $line) {
			exit(0);
		}
		readline_add_history(substr($str, 0, -1));
		echo $fmt->formatCode('<?php ' . substr($str, 0, -1)), PHP_EOL;
	}
} elseif (isset($opts['o'])) {
	$argv = extractFromArgvShort($argv, 'o');
	if ('-' == $opts['o'] && '-' == $argv[1]) {
		echo $fmt->formatCode(file_get_contents('php://stdin'));
		exit(0);
	}
	if ($inPhar) {
		if (!file_exists($argv[1])) {
			$argv[1] = getcwd() . DIRECTORY_SEPARATOR . $argv[1];
		}
	}
	if (!is_file($argv[1])) {
		fwrite(STDERR, 'File not found: ' . $argv[1] . PHP_EOL);
		exit(255);
	}
	if ('-' == $opts['o']) {
		echo $fmt->formatCode(file_get_contents($argv[1]));
		exit(0);
	}
	$argv = array_values($argv);
	file_put_contents($opts['o'], $fmt->formatCode(file_get_contents($argv[1])));
} elseif (isset($argv[1])) {
	if ('-' == $argv[1]) {
		echo $fmt->formatCode(file_get_contents('php://stdin'));
		exit(0);
	}
	$fileNotFound = false;
	$start = microtime(true);
	fwrite(STDERR, 'Formatting ...' . PHP_EOL);
	$missingFiles = [];
	$fileCount = 0;

	$cacheHitCount = 0;
	$workers = 4;

	$hasFnSeparator = false;

	// Used with dry-run to flag if any files would have been changed
	$filesChanged = false;

	for ($j = 1; $j < $argc; ++$j) {
		$arg = &$argv[$j];
		if (!isset($arg)) {
			continue;
		}
		if ('--' == $arg) {
			$hasFnSeparator = true;
			continue;
		}
		if ($inPhar && !file_exists($arg)) {
			$arg = getcwd() . DIRECTORY_SEPARATOR . $arg;
		}
		if (is_file($arg)) {
			$file = $arg;
			if ($lintBefore && !lint($file)) {
				fwrite(STDERR, 'Error lint:' . $file . PHP_EOL);
				continue;
			}
			++$fileCount;
			fwrite(STDERR, '.');
			$fileContents = file_get_contents($file);
			$formattedCode = $fmt->formatCode($fileContents);
			if ($dryRun) {
				if ($fileContents !== $formattedCode) {
					$filesChanged = true;
				}
			} else {
				file_put_contents($file . '-tmp', $formattedCode);
				$oldchmod = fileperms($file);
				rename($file . '-tmp', $file);
				chmod($file, $oldchmod);
			}
		} elseif (is_dir($arg)) {
			fwrite(STDERR, $arg . PHP_EOL);

			$target_dir = $arg;
			$dir = new RecursiveDirectoryIterator($target_dir);
			$it = new RecursiveIteratorIterator($dir);
			$files = new RegexIterator($it, '/^.+\.php$/i', RecursiveRegexIterator::GET_MATCH);

			if ($concurrent) {
				$chn = make_channel();
				$chn_done = make_channel();
				if ($concurrent) {
					fwrite(STDERR, 'Starting ' . $workers . ' workers ...' . PHP_EOL);
				}
				for ($i = 0; $i < $workers; ++$i) {
					cofunc(function ($fmt, $backup, $cache_fn, $chn, $chn_done, $lintBefore, $dryRun) {
						$cache = new Cache($cache_fn);
						$cacheHitCount = 0;
						$cache_miss_count = 0;
						$filesChanged = false;
						while (true) {
							$msg = $chn->out();
							if (null === $msg) {
								break;
							}
							$target_dir = $msg['target_dir'];
							$file = $msg['file'];
							if (empty($file)) {
								continue;
							}
							if ($lintBefore && !lint($file)) {
								fwrite(STDERR, 'Error lint:' . $file . PHP_EOL);
								continue;
							}

							$content = $cache->is_changed($target_dir, $file);
							if (false === $content) {
								++$cacheHitCount;
								continue;
							}

							++$cache_miss_count;
							$fmtCode = $fmt->formatCode($content);
							if (null !== $cache) {
								$cache->upsert($target_dir, $file, $fmtCode);
							}
							if ($dryRun) {
								if ($fmtCode !== $content) {
									$filesChanged = true;
								}
							} else {
								file_put_contents($file . '-tmp', $fmtCode);
								$oldchmod = fileperms($file);
								$backup && rename($file, $file . '~');
								rename($file . '-tmp', $file);
								chmod($file, $oldchmod);
							}
						}
						$chn_done->in([$cacheHitCount, $cache_miss_count, $filesChanged]);
					}, $fmt, $backup, $cache_fn, $chn, $chn_done, $lintBefore, $dryRun);
				}
			}

			foreach ($files as $file) {
				$file = $file[0];
				if (null !== $ignore_list) {
					foreach ($ignore_list as $pattern) {
						if (false !== strpos($file, $pattern)) {
							continue 2;
						}
					}
				}

				fwrite(STDERR, '.');

				++$fileCount;
				if ($concurrent) {
					$chn->in([
						'target_dir' => $target_dir,
						'file' => $file,
					]);
				} else {
					if (0 == ($fileCount % 20)) {
						fwrite(STDERR, ' ' . $fileCount . PHP_EOL);
					}
					$content = $cache->is_changed($target_dir, $file);
					if (false === $content) {
						++$fileCount;
						++$cacheHitCount;
						continue;
					}
					if ($lintBefore && !lint($file)) {
						fwrite(STDERR, 'Error lint:' . $file . PHP_EOL);
						continue;
					}
					$fmtCode = $fmt->formatCode($content);
					fwrite(STDERR, '.');
					if (null !== $cache) {
						$cache->upsert($target_dir, $file, $fmtCode);
					}
					if ($dryRun) {
						if ($fmtCode !== $content) {
							$filesChanged = true;
						}
					} else {
						file_put_contents($file . '-tmp', $fmtCode);
						$oldchmod = fileperms($file);
						$backup && rename($file, $file . '~');
						rename($file . '-tmp', $file);
						chmod($file, $oldchmod);
					}
				}
			}
			if ($concurrent) {
				for ($i = 0; $i < $workers; ++$i) {
					$chn->in(null);
				}
				for ($i = 0; $i < $workers; ++$i) {
					list($cache_hit, $cache_miss, $filesChanged) = $chn_done->out();
					$cacheHitCount += $cache_hit;
				}
				$chn_done->close();
				$chn->close();
			}
			fwrite(STDERR, PHP_EOL);

			continue;
		} elseif (
			!is_file($arg) &&
			('--' != substr($arg, 0, 2) || $hasFnSeparator)
		) {
			$fileNotFound = true;
			$missingFiles[] = $arg;
			fwrite(STDERR, '!');
		}
		if (0 == ($fileCount % 20)) {
			fwrite(STDERR, ' ' . $fileCount . PHP_EOL);
		}
	}
	fwrite(STDERR, PHP_EOL);
	if (null !== $cache) {
		fwrite(STDERR, ' ' . $cacheHitCount . ' files untouched (cache hit)' . PHP_EOL);
	}
	fwrite(STDERR, ' ' . $fileCount . ' files total' . PHP_EOL);
	fwrite(STDERR, 'Took ' . round(microtime(true) - $start, 2) . 's' . PHP_EOL);
	if (sizeof($missingFiles)) {
		fwrite(STDERR, 'Files not found: ' . PHP_EOL);
		foreach ($missingFiles as $file) {
			fwrite(STDERR, "\t - " . $file . PHP_EOL);
		}
	}
	if ($dryRun && $filesChanged) {
		exit(1);
	}
	if ($fileNotFound) {
		exit(255);
	}
} else {
	showHelp($argv, $enableCache, $inPhar);
}
exit(0);
