{{define "content"}}
<h1 class="compact">设置列表</h1>
<form accept-charset="UTF-8" id="new_config">
	<div id="field">
	<label for="port">端口</label><span> [Read only]</span>
	<br />
	<input type="text" id="port" name="port" value="{{.Port}}" readonly="" size="30"/>
	</div>
	<div id="field">
	<label for="auth">是否开启反向代理</label><span> [true/false]</span>
	<br />
	<input type="text" pattern="false|true" id="reverse" name="reverse" value="{{.Reverse}}" size="30" />
	</div>
	<div id="field">
	<label for="auth">反向代理目标地址</label><span> eg:"127.0.0.1:8090"</span>
	<br />
	<input type="text" id="proxy_pass" name="proxy_pass" value="{{.ProxyPass}}" size="30" />
	</div>
	<div id="field">
	<label for="auth">是否开启认证</label><span> [true/false]</span>
	<br />
	<input type="text" pattern="false|true" id="auth" name="auth" value="{{.Auth}}" size="30" />
	</div>
	<div id="field">
	<label for="cache">是否开启缓存</label><span> [true/false]</span>
	<br />
	<input type="text" pattern="false|true" id="cache" name="cache" value="{{.Cache}}" size="30" />
	</div>
	<div id="field">
	<label for="cachetimeout">缓存更新时间</label><span> [minutes]</span>
	<br />
	<input type="text" pattern="[0-9]+" id="cachetimeout" name="cachetimeout" value="{{.CacheTimeout}}" size="30" />
	</div>
	<div id="field">
	<label for="gfwlist">网站过滤列表</label><span> [关键词英文分号分割]</span>
	<br />
	<input type="text" id="gfwlist" name="gfwlist" {{if eq (len .GFWList) 0}}{{else}}value="{{range .GFWList}}{{.}};{{end}}"{{end}} size="30" />
	</div>
	<div id="field">
	<label for="log">日志</label><span> [1 调试模式/0 普通模式]</span>
	<br />
	<input type="text" pattern="[0-1]" id="log" name="log" value="{{.Log}}" size="30" readonly />
	</div>
	<div class="actions"><input type="submit" value="设置" /></div>
</form>
</script>
<script type="text/javascript">
	$('#new_config').submit( function(e) {
		$.ajax({
			type:'POST',
			url:'/setting/set',
			data:$(this).serialize(),
			error: function(response) {
				alert("failed");
            },
			success: function() {
				window.location.reload()
			}
		});
	});
	</script>
{{end}}