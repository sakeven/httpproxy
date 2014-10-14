<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	<title>Http Proxy Web</title>
	 <link href="/static/css/main.css" rel="stylesheet" type="text/css" />
	 <script src="/static/js/jquery.min.js" type="text/javascript"></script>
    <script src="/static/js/action.js" type="text/javascript"></script>
</head>
<body>
 	<div class="container">
 		<div id="pageHeader">
		<h1>Http Proxy Web</h1> 		
 		</div>

    <div id="navibar" class="span-3 last">
        <ul>
          <li>{{if eq .Nav "home"}}<span>主页</span>{{else}}<a href="/">主页</a>{{end}}</li>
          <li>{{if eq .Nav "user"}}<span>用户</span>{{else}}<a href="/user/list/detail">用户</a>{{end}}</li>
          <li>{{if eq .Nav "setting"}}<span>设置</span>{{else}}<a href="/setting/list">设置</a>{{end}}</li>
        </ul>
      </div>

      <div id="body" class="span-22">
	 {{template "content" .}}      	
      </div>
	</div>
</body>
</html>