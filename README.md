# What is soksan
*soksan* allows you to embed a go playground widget on your personal website (the backend can be a PHP-based or a Go server).

![Custom playground widget on your website](/images/example.jpg "Custom playground widget on your website")

# Why ?

The official Playground website doesn't set any <code>Origin</code> header and doesn't implement CORS mechanism. So it is not possible to emit a JSON request from a Javascript belonging to your domain and that targets playground domain (this configuration is called "cross domain request"). In such a context, you need a mediator server application that will receive the JSON requests of your application and reroute them to the playground service. That is why *soksan* has been developed. It contains server code and client code samples that you can directly use in your own web site.

# Usage

The go team accept this usage, but you need to warn them first. Before embedding a widget that will use the go playground service, define a unique user-agent and send an e-mail to golang-dev mailing list first. See : http://blog.golang.org/playground#TOC_7

## Backend

### PHP backend

Install the files on your web server :
* "php" folder
* ".htaccess" file
* optionally the frontend sample files

configure the address of the playground and the user agent name you chose in php/config.php file
* <code>HOST_PLAYGROUND</code> : playground service (e.g. 'http://play.golang.org/' or your own instance).
* <code>USER_AGENT</code> : user agent (e.g. 'soksan') **YOU NEED TO CHANGE THAT VALUE**.
* <code>SAMPLES_PATH</code> : Source path of go code samples stored on your server '../gocode/' (for <code>run</code> endpoint).

You can use the .htaccess file provided at the root of this repository so as to define rewrite rules or enable gzip compression

### Golang backend

See the Go example provided in main.go which is a standalone web server relying on soksan library. First include the lib :

<code>import "github.com/bbalet/soksan/soksan"</code>

It will add the following handlers to your application :
* <code>/fmt</code> Format a go code using gofmt
* <code>/compile</code> Compile and execute a go code
* <code>/run</code> Variant of <code>compile</code> but with a file stored on the mediator server

You need to initialize the library :
* <code>soksan.HostPlayGround</code> : playground service (e.g. 'http://play.golang.org/' or your own instance).
* <code>soksan.UserAgent</code> : user agent (e.g. 'soksan') **YOU NEED TO CHANGE THAT VALUE**.
* <code>soksan.SamplePath</code> : Source path of go code samples stored on your server '../gocode/' (for <code>run</code> endpoint).

It is your responsability to start the web server as it is illustrated into main.go example. If you want to play with this code example, compile it and run it with these commands (you may modify config/config.json to suit your needs first) :
<code>
go build .<br />
soksan run
</code>

## Frontend

### Examples

*soksan* comes with some frontend examples :

* *go-playground.html* how to use the genuine playground.js file available on play.golang.org
* *custom-playground.html* how to use a custom playground javascript to mimic the playground
* *custom-playground-modal.html* how to use my custom playground javascript to run a go code sample in a modal form
* *file-playground.html* how to display a portion of code but compile a file stored on the server side

You can see a live example on this french website :
http://decouvrir-golang.net/generalites/des-nouvelles-du-front-annee-2014.html

(please submit your own websites)

### Usage

* Create the editor. If you are using the genuine playground.js, it must be a TEXTAREA. Otherwise, any container.
* Create an ouput DIV
* Create an Execute button
* Optionally, create a format button
* Insert a Javascript code that initializes the components

<code>
	$(function() {<br />
		//Example with a code editor, output container, run button and a format button<br />
		initEditor($('#code'), $('#output'), $('#run'), $('#fmt'));<br />
	});
</code>

If you want to display a portion of code but compile a file stored on the server side, add a <code>data-file</code> attribute to the code editor and pass it to the inti function :

<code>
	$(function() {<br />
		//Example with a code editor, output container, run button and a data-file attribute<br />
		initEditor($('#code'), $('#output'), $('#run'), null, true);<br />
	});
</code>

