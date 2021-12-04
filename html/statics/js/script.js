function checkGroup(){
	axios({
		method : 'get',
		url: '/checkGroup'
	})
	.then(function(response){
		if(response.data === "维修员"){
			document.getElementById("register").style.display="none";
			var text = document.getElementById("text");
			axios({
				method : 'get',
				url: '/checkStatus'
			})
				.then(function(response1){
					text.innerHTML = text.innerHTML + "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;工作状态：" + response1.data;
				var body = document.getElementById("body");
				var div = document.createElement("div");
				body.appendChild(div);
				div.id = "status";
				var h1 = document.createElement("h1");
				div.appendChild(h1);
				h1.style.cssText = "display: inline-block;";
				h1.innerHTML = "更改工作状态：";
				var select = document.createElement("select");
				select.id = "select";
				div.appendChild(select);
				select.name = "status";

				var option1 = document.createElement("option");
				select.appendChild(option1);
				option1.value = "正常";
				option1.innerHTML = "正常";
				var option2 = document.createElement("option");
				select.appendChild(option2);
				option2.value = "休假";
				option2.innerHTML = "休假";
				var option3 = document.createElement("option");
				select.appendChild(option3);
				option3.value = "离职";
				option3.innerHTML = "离职";
				console.log(response1.data);
				console.log(response1.data);
				if(response1.data === "正常"){
					option1.setAttribute("selected", true);
				}
				if(response1.data === "休假"){
					option2.setAttribute("selected", true);
				}
				if(response1.data === "离职"){
					option3.setAttribute("selected", true);
				}

				var button = document.createElement("button");
				div.appendChild(button);
				button.onclick= function () { changeStatus(); };
				button.style.cssText = "width: 50px; height:20px";
				button.innerHTML = "更改";

				})
		}
	})
}
function changeStatus(){		//https://blog.csdn.net/weixin_41949511/article/details/93630346
	axios({
		method: 'post',
		url: '/changeStatus'
	});
	let formData = new FormData();
	var select = document.getElementById("select");
	formData.append("status", select.value);
	let config = {
	   headers: {"Content-Type": "multipart/form-data"}
	};
	axios.post("/changeStatus", formData, config).then(res => {
			 alert(res.data);
			 location.reload();
		}).catch(error => {
			 alert(error.response.data);
			 location.reload();
		});
}

function checkNotification(){
	axios({
	  method: 'get',
	  url: '/checkNotification'
	  })
	  .then(function(response) {
		var last = JSON.stringify(response.data);
		if(last !== "null"){
			var body = document.getElementById("body");
			var div = document.createElement("div");
			body.appendChild(div);
			div.style.cssText = "width: 50%; height: 50%;position: absolute; top: 25%; left: 25%; text-align: center; background-color: white;";
			div.id = "alert";
			var title = document.createElement("h1");
			div.appendChild(title);
			title.innerHTML = last;
			var button = document.createElement("button");
			div.appendChild(button);
			button.onclick= function () { read(); };
			button.style.cssText = "width: 50px; height:30px";
			button.innerHTML = "已读";
			var button1 = document.createElement("button");
			div.appendChild(button1);
			button1.onclick= function () { closeNotification(); };
			button1.style.cssText = "width: 100px; height:30px";
			button1.innerHTML = "暂时忽略";
		}
	  });
}
function read(){
	axios({
		method: 'post',
		url: '/read'
	});
	closeNotification();
	checkNotification();
}
function closeNotification(){
	var father = document.getElementById("body");
	var child = document.getElementById("alert");
	father.removeChild(child);
}

function checkUser(){
	document.getElementById("logout").style.display="none";
	document.getElementById("changePassword").style.display="none";
	axios({
		method : 'get',
		url: '/checkGroup'
	})
		.then(function(response){
			if(response.data === "维修员"||response.data ==="业务员"||response.data ==="普通用户"){
				document.getElementById("register").style.display="none";
				document.getElementById("login").style.display="none";
				document.getElementById("logout").style.display="";
				document.getElementById("changePassword").style.display="";
			}
		})
}

function test(){
	if(document.getElementById("txt").value === "") return;
	axios({
		method: 'get',
		url: '/test?text='+document.getElementById("txt").value
	})
		.then(function(response) {
			try {
				var text = document.getElementById("res");
				text.innerHTML = JSON.stringify(response.data);
			}catch (error)
			{
				var div = document.getElementById("test_div");
				var result_div = document.createElement("div");
				div.appendChild(result_div);
				result_div.id = "result";
				result_div.style.cssText = "width: 200px; height: 200px;position: relative; text-align: left; background-color: red; margin-left:10%";
				var text = document.createElement("p");
				result_div.appendChild(text);
				text.id = "res";
				text.innerHTML = JSON.stringify(response.data);
			}
		})
}

function deltest(){
	try {
		var div = document.getElementById("test_div");
		var res = document.getElementById("result");
		div.removeChild(res);
	}catch (error)
	{

	}
}