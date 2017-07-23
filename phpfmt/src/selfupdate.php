<?php
# Copyright (c) 2016, phpfmt and its authors
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

function selfupdate(array $argv, bool $inPhar, string $channel = 'lts') {
	$opts = [
		'http' => [
			'method' => 'GET',
			'header' => "User-agent: phpfmt fmt.phar selfupdate\r\n",
		],
	];

	$context = stream_context_create($opts);

	// current release
	$channels = json_decode(file_get_contents('https://raw.githubusercontent.com/phpfmt/releases/master/releases.json', false, $context), true);
	if (!isset($channels[$channel])) {
		fwrite(STDERR, 'channel not found: ' . $channel . PHP_EOL);
		exit(1);
	}

	$version = $channels[$channel];
	$downloadURL = 'https://github.com/phpfmt/releases/raw/master/releases/' . $channel . '/' . $version . '/fmt.phar';

	$pharFile = file_get_contents($downloadURL);
	$pharSHA1 = file_get_contents($downloadURL . '.sha1');

	if ($inPhar && !file_exists($argv[0])) {
		$argv[0] = dirname(Phar::running(false)) . DIRECTORY_SEPARATOR . $argv[0];
	}
	if (sha1_file($argv[0]) != $pharSHA1) {
		copy($argv[0], $argv[0] . '~');
		file_put_contents($argv[0], $pharFile);
		chmod($argv[0], 0777 & ~umask());
		fwrite(STDERR, 'Updated successfully' . PHP_EOL);
		exit(0);
	}
	fwrite(STDERR, 'Up-to-date!' . PHP_EOL);
	exit(0);
}