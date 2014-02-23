<?php
/*  soksan allows you to interact with a go playground 
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

// Make a custom request by filling body variable with the content of a designated file
array_push($_POST, 'version', '2');
array_push($_POST, 'body', file_get_contents(SAMPLES_PATH . $_POST['file']));
unset($_POST['file']);	//Don't bothe Go playground with extra parameters

//echo var_dump($_POST);

// Send a compile request
echo sendPlaygroundRequest('compile');
