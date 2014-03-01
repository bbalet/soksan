<?php
/*  soksan allows you to embed a go playground in your website
    Copyright (C) 2014 Benjamin BALET

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.*/
	
include 'config.php';
include 'curl.php';

$filename = pathJoin(realpath(dirname(__FILE__)), pathClean(SAMPLES_PATH), pathClean($_POST['file']));

//Test if file exists
if (!file_exists($filename)) {
	header($_SERVER['SERVER_PROTOCOL'] . 'Internal Server Error', true, 500);
    echo "file " . $filename . " doesn't exist";
}
else {
	// Make a custom request by filling body variable with the content of a designated file
	$data = array('version' => '2', 'body' => file_get_contents($filename));

	// Send a compile request
	echo postPlaygroundRequest('compile', $data);
}

 /**
 * This function cleans a path by uniformizing the path separator
 * which can vary depending on the OS system
 * @param string $path The path to normalize
 * @return string Sanitized path
 */
function pathClean($path) {
	$sanitizedPath = str_replace("\\", DIRECTORY_SEPARATOR, $path);
	$sanitizedPath = str_replace("/", DIRECTORY_SEPARATOR, $sanitizedPath);
	return $sanitizedPath;
}

 /**
 * This function takes a variable amount of strings and joins
 * them together so that they form a valid file path.
 * @param array $pieces - The pieces of the file path
 * @returns string The final file path
 */
function pathJoin() {
	$pieces = array_filter(func_get_args(), function($value) {
		return $value;
	});
	return pathNormalize(implode(DIRECTORY_SEPARATOR, $pieces));
}

/**
 * This function takes a valid file path and nomalizes it into
 * the simplest form possible.
 * @param string $path The path to normalize
 * @returns string The normailized path
 */
function pathNormalize($path) {
	if (!strlen($path)) {
		return ".";
	}

	$isAbsolute    = $path[0];
	$trailingSlash = $path[strlen($path) - 1];

	$up     = 0;
	$pieces = array_values(array_filter(explode(DIRECTORY_SEPARATOR, $path), function($n) {
		return !!$n;
	}));
	for ($i = count($pieces) - 1; $i >= 0; $i--) {
		$last = $pieces[$i];
		if ($last == ".") {
			array_splice($pieces, $i, 1);
		} else if ($last == "..") {
			array_splice($pieces, $i, 1);
			$up++;
		} else if ($up) {
			array_splice($pieces, $i, 1);
			$up--;
		}
	}

	$path = implode(DIRECTORY_SEPARATOR, $pieces);

	if (!$path && !$isAbsolute) {
		$path = ".";
	}

	if ($path && $trailingSlash == DIRECTORY_SEPARATOR) {
		$path .= DIRECTORY_SEPARATOR;
	}

	return ($isAbsolute == DIRECTORY_SEPARATOR ? DIRECTORY_SEPARATOR : "") . $path;
}
