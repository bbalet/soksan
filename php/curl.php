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

/**
 * Wrapper to send the POST parameters of the page
 * @uses  postPlaygroundRequest()
 * @param string $endpoint
 * @return string JSON message sent back by the playground
 */
function sendPlaygroundRequest($endpoint) {
	return postPlaygroundRequest($endpoint, $_POST);
}

/**
 * Make a curl POST request to the configured go playground service and 
 * @param string $source
 * @param string $params
 * @return string JSON message sent back by the playground
 */
function postPlaygroundRequest($endpoint, $params) {
// Setup cURL
$ch = curl_init(HOST_PLAYGROUND . '/' . $endpoint);
curl_setopt_array($ch, array(
    CURLOPT_POST => TRUE,
    CURLOPT_RETURNTRANSFER => TRUE,
    CURLOPT_POSTFIELDS => $params,
	CURLOPT_USERAGENT => USER_AGENT
));

// Send the request
$response = curl_exec($ch);

// Check for errors
if($response === FALSE){
	header($_SERVER['SERVER_PROTOCOL'] . 'Internal Server Error', true, 500);
    die(curl_error($ch));
}

// return the response as received
header("Content-type: application/json; charset: utf-8");
return $response;
}