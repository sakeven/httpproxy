{{define "content"}}
<h1 class="compact">用户列表</h1>
<table class="userlist">
		<thead>
		<tr>
		    <th class="header">用户名</th>
		    <th class="header">密码</th>
		    <th class="header">操作</th>
		</tr>
		</thead>
		<tbody>
		{{with .User}}
		{{range $user,$passwd :=.}}
		<tr>
			<td>{{$user}}</td>
			<td><input type="text" id="passwd{{$user}}" name="passwd" value="{{$passwd}}" required /></td>
			<td><select name="action" data-id="{{$user}}">
			<option selected value="">无</option>
			<option value="delete">删除</option>
			<option value="modify">修改</option>
			</select></td>
		</tr>
		{{end}}
		{{end}}
		<tr>
		<form accept-charset="UTF-8" id="new_user">
			<td><input type="text" name="user" required /></td>
			<td><input type="text" name="passwd" required /></td>
			<td><div class="actions"><input type="submit" value="增加" /></div></td>
		</form>
		</tr>
	</tbody>
</table>
<script type="text/javascript">
$('select').on('change', function() {
	var user = $(this).data('id')
  	var ret = confirm($(this).val()+' the user '+user+'?');
  	if (ret == true){
  		if ($(this).val() == 'delete'){
  		$.ajax({
    	type:'POST',
      	url:'/user/delete/'+user,
      	data:$(this).serialize(),
      	error: function() {
        	alert('failed!');
      	},
      	success: function() {
        	window.location.reload();
      	}
    	});
  	} else if ( $(this).val() == 'modify') {
  		$.ajax({
    	type:'POST',
      	url:'/user/modify/'+user,
      	data:$('input#passwd'+user).serialize(),
      	error: function() {
        	alert('failed!');
      	},
      	success: function() {
        	window.location.reload()
      	}
    	});
  	};
  }
});
</script>
<script type="text/javascript">
	$(document).ready(function(){
		$(".userlist tr:even").addClass("even");
	});
</script>
<script type="text/javascript">
	$('#new_user').submit( function(e) {
		$.ajax({
			type:'POST',
			url:'/user/add/new',
			data:$(this).serialize(),
			error: function(response) {
				alert('failed')
			},
			success: function() {
				window.location.reload()
			}
		});
	});
	</script>
{{end}}