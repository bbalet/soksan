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
	
	function initEditor(editor, output, btnRun, btnFormat, file) {
	
		//Handle default parameters
		if(typeof(btnFormat)==='undefined') btnFormat = null;
		if(typeof(file)==='undefined') file = false;
	
		//Set the height of the editor with the length of code
		if (editor.prop("tagName") == "TEXTAREA") {
			editor.height(0);
			editor.height($('#code').prop('scrollHeight'));
			//Automatically update the height of the editor on any change
			editor.bind('input propertychange', function() {
				editor.height( 0 );
				editor.height( $('#code').prop('scrollHeight') );
			});
		}
		
		btnRun.click(function() {
			output.css('visibility', 'visible');
			output.html("<img src='/images/ajax-loader.gif'>&nbsp;<i>Compiling...</i><br />");
			output.trigger('show');
			
			var mydata;
			var myurl = '/compile';
			if (file) {
				mydata = {'version': 2, 'file': editor.data('filename')};
				myurl = '/run';
			}
			else {
				var code = "";
				if (editor.prop("tagName") == "TEXTAREA")
					code = editor.val();
				else
					code = editor.text();
				mydata = {'version': 2, 'body': code};
			}
		
			$.ajax(myurl, {
				type: 'POST',
				data: mydata,
				dataType: 'json',
				success: function(data) {
					if (data.Errors) {
						output.html("<span class='text-error'>" + data.Errors.replace("\n","<br />") + "</span>");
					}
					else {
						output.html('');
						playback(output, data.Events);
					}
				},
				error: function() {
					output.html("<strong>Communication error</strong><br />");
				}
			});
		});
		
		if (btnFormat) {
			btnFormat.click(function() {
				output.css('visibility', 'visible');
				output.html("<img src='/images/ajax-loader.gif'>&nbsp;<i>Formatting...</i><br />");
				output.trigger('show');
			
				$.ajax("/fmt", {
				data: {"body":editor.text()},
				type: "POST",
				dataType: "json",
				success: function(data) {
				  if (data.Error) {
					output.html("<strong>Communication error</strong><br />");
					var text = "";
					for (var i = 0; i < data.Errors.length; i++) {
						text += data.Errors[i].Message + "<br />";
					}
					output.html(text);
				  } else {
					editor.val(data.Body);
					editor.trigger('propertychange');
					output.html('Formatted');
				  }
				}
			  });
			});
		}		
	}

//------------------------------------------------------------------------------
//Most of the code below this line has been stolen from http://play.golang.org/
	function playback(output, events) {
		var timeout;
		write(output, {Kind: 'start'});
		function next() {
			if (events.length === 0) {
				write(output, {Kind: 'end'});
				return;
			}
			var e = events.shift();
			if (e.Delay === 0) {
				write(output, {Kind: 'stdout', Body: e.Message});
				next();
				return;
			}
			timeout = setTimeout(function() {
				write(output, {Kind: 'stdout', Body: e.Message});
				next();
			}, e.Delay / 1000000);
		}
		next();
		return {
			Stop: function() {
				clearTimeout(timeout);
			}
		}
	}
	
	function write(output, event) {
		if (event.Kind == 'start') {
			output.innerHTML = '';
			return;
		}

		var cl = 'system';
		if (event.Kind == 'stdout' || event.Kind == 'stderr')
			cl = event.Kind;

		var m = event.Body;
		if (event.Kind == 'end') 
			m = '\nProgram exited' + (m?(': '+m):'.');

		if (m.indexOf('IMAGE:') === 0) {
			var url = 'data:image/png;base64,' + m.substr(6);
			var img = document.createElement('img');
			img.src = url;
			output.appendChild(img);
			return;
		}

		// ^L clears the screen.
		var s = m.split('\x0c');
		if (s.length > 1) {
			output.innerHTML = '';
			m = s.pop();
		}

		m = m.replace(/&/g, '&amp;');
		m = m.replace(/</g, '&lt;');
		m = m.replace(/>/g, '&gt;');

		var needScroll = (output.scrollTop + output.offsetHeight) == output.scrollHeight;

		var span = document.createElement('span');
		span.className = cl;
		span.innerHTML = m;
		output.append(span);

		if (needScroll)
			output.scrollTop = output.scrollHeight - output.offsetHeight;
	}
	