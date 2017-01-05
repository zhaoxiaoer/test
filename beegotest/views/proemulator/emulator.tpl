<form action={{.url}} method="post">
请求 url:<br>
<input type="text" name="posturl" value={{.posturl}} style="width:600px;height:30px"><br>
请求参数:<br>
<input type="text" name="postbody" value="key=1111" style="width:600px;height:30px"><br>
<input type="submit" value="测试">
</form><br>
{{if .body}}
<p>{{.body}}</p>
{{end}}