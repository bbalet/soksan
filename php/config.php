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

// Go Playground, your own google app engine or your own host
const HOST_PLAYGROUND = 'http://play.golang.org/';

// If you use play.golang.org, please use a unique user agent
// and contact golang-dev@googlegroups.com first
// See: http://blog.golang.org/playground#TOC_7
const USER_AGENT = 'soksan';

// Location of go code sample files
const SAMPLES_PATH = '../gocode/';
