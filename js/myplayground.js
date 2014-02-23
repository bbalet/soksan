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
			output.html("<img src='/images/ajax-loader.gif'>&nbsp;<i>En cours d'exécution</i><br />");
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
						var text = "<span class='text-info'>";
						for (var i = 0; i < data.Events.length; i++) {
							text += data.Events[i].Message + "<br />";
						}
						text += "</span>";
						text += "<span class='muted'>Fin du programme.</span>";
						output.html(text);
					}
				},
				error: function() {
					output.html("<strong>Problème de communication avec le serveur</strong><br />");
				}
			});
		});
		
		if (btnFormat) {
			btnFormat.click(function() {
				$.ajax("/fmt", {
				data: {"body":editor.text()},
				type: "POST",
				dataType: "json",
				success: function(data) {
				  if (data.Error) {
					output.html("<strong>Problème de communication avec le serveur</strong><br />");
					var text = "";
					for (var i = 0; i < data.Errors.length; i++) {
						text += data.Errors[i].Message + "<br />";
					}
					output.html(text);
				  } else {
					editor.val(data.Body);
					editor.trigger('propertychange');
				  }
				}
			  });
			});
		}		
	}
	