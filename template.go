package main

import "html/template"

var html = template.Must(template.New("chat_room").Parse(`
<html>
<head>
    <title>{{.roomid}}</title>
    <link rel="stylesheet" type="text/css" href="http://meyerweb.com/eric/tools/css/reset/reset.css"/>
    <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.7/jquery.js"></script>
    <script src="http://malsup.github.com/jquery.form.js"></script>
    <script type="text/javascript">
        $('#message_form').focus();
        $(document).ready(function() {
            // bind 'myForm' and provide a simple callback function
            console.log("A");
            $('#myForm').ajaxForm(function() {
                $('#message_form').val('');
                $('#message_form').focus();
            });
            if (!!window.EventSource) {
                console.log("B");
                var source = new EventSource('{{.baseUrl}}/streams/{{.roomid}}');
                source.addEventListener('message', function(e) {
                    $('#messages').append(e.data + "</br>");
                    $('html, body').animate({scrollTop:$(document).height()}, 'slow');
                }, false);
            } else {
                alert("NOT SUPPORTED");
            }
        });
    </script>
    </head>
    <body>
    <h1>Welcome to {{.roomid}} room</h1>
    <div id="messages"></div>
    <form id="myForm" action="{{.baseUrl}}/rooms/{{.roomid}}" method="post">
    User: <input id="user_form" name="user" value="{{.userid}}"></input>
    Message: <input id="message_form" name="message"></input>
    <input type="submit" value="Submit" />
    </form>
</body>
</html>
`))
