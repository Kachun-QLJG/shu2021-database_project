function login(){
	let formData = new FormData();
	formData.append("username", document.getElementById("username").value);
	formData.append("password", document.getElementById("password").value);
	formData.append("ver_code", document.getElementById("ver_code").value);
	let config = {
		headers: {"Content-Type": "multipart/form-data"}
	};
	axios.post("/login", formData, config).then(res => {
		var response = res.data;
		if(response.status === "失败")
		{
			if(response.data === '验证码'){
				alert("验证码错误！");
				document.getElementById("wrongPass").innerHTML = '请输入正确的验证码！';
				document.getElementById("wrongPass").style.display = '';
				document.getElementById('ver_pic').click();
			}else{
				alert("用户名或密码错误！");
				document.getElementById("wrongPass").innerHTML = '请输入正确的账号和密码！';
				document.getElementById("wrongPass").style.display = '';
				document.getElementById('ver_pic').click();
			}
		}
		else{
			if(response.data === '已登录') {
				alert("已登录，请勿重复登录！");
			}
			window.open("/index", "_self");
		}
	})
}

async function logout() {
	await axios({
		method: 'post',
		url: '/logout'
	});
	window.open("/index", "_self");
}
function register(){
	let registerData = new FormData();
	registerData.append("phone", document.getElementById("phone").value);
	registerData.append("pswd", document.getElementById("pswd").value);
	registerData.append("ver", document.getElementById("ver").value);
	let configRegister = {
		headers: {"Content-Type": "multipart/form-data"}
	};
	// let loginData = new FormData();
	// loginData.append("username", document.getElementById("phone").value);
	// loginData.append("password", document.getElementById("pswd").value);
	// loginData.append("ver_code", document.getElementById("ver").value);
	// let configLogin = {
	// 	headers: {"Content-Type": "multipart/form-data"}
	// };
	axios.post("/register", registerData, configRegister).then(res => {
		var response = res.data;
		if(response.status === "失败")
		{
			if(response.data === '验证码'){
				alert("验证码错误！");
				document.getElementById("wrongPass").innerHTML = '验证码错误！';
			}else{
				alert(response.data);
				document.getElementById("wrongPass").innerHTML = response.data;
			}
			document.getElementById("wrongPass").style.display = '';
			document.getElementById('ver_pic').click();
		}
		else{
			alert("注册成功！");
			window.open("/index", "_self");
			// axios.post("/login", loginData, configLogin).then(res => {
			// 	var response = res.data;
			// 	alert(response.status);
			// 	if (response.status === "失败") {
			// 		if (response.data === '验证码') {
			// 			alert("验证码错误！");
			// 			document.getElementById("wrongPass").innerHTML = '请输入正确的验证码！';
			// 			document.getElementById("wrongPass").style.display = '';
			// 			document.getElementById('ver_pic').click();
			// 		} else {
			// 			alert("用户名或密码错误！");
			// 			document.getElementById("wrongPass").innerHTML = '请输入正确的账号和密码！';
			// 			document.getElementById("wrongPass").style.display = '';
			// 			document.getElementById('ver_pic').click();
			// 		}
			// 	} else {
			// 		if (response.data === '已登录') {
			// 			alert("已登录，请勿重复登录！");
			// 		}
			// 		window.open("/index", "_self");
			// 	}
			// })
		}
	})


}

function changePswd(){
	let formData = new FormData();
	formData.append("oldpswd", document.getElementById("oldpswd").value);
	formData.append("pswd", document.getElementById("pswd").value);
	let config = {
		headers: {"Content-Type": "multipart/form-data"}
	};
	axios.post("/change_password", formData, config).then(res => {
		var response = res.data;
		if(response.status === "失败")
		{
			alert(response.data);
		}
		else{
			alert("密码更改成功！请重新登录！");
			window.open("/index", "_self");
		}
	})
}
/*function displayChangeStatus() {
	document.getElementById("register").style.display="none";
	var text = document.getElementById("text");
	axios({
		method : 'get',
		url: '/check_status'
	})
		.then(function(response1) {
			text.innerHTML = text.innerHTML + "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;工作状态：" + response1.data;
			var body = document.getElementById("con");
			var div = document.getElementById("workStatus");
			body.appendChild(div);
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
			if (response1.data === "正常") {
				option1.setAttribute("selected", true);
			}
			if (response1.data === "休假") {
				option2.setAttribute("selected", true);
			}
			if (response1.data === "离职") {
				option3.setAttribute("selected", true);
			}

			var button = document.createElement("button");
			div.appendChild(button);
			button.onclick = function () {
				changeStatus();
			};
			button.style.cssText = "width: 50px; height:20px";
			button.innerHTML = "更改";
		})
}
*/
function changeStatus(){		//https://blog.csdn.net/weixin_41949511/article/details/93630346
	let formData = new FormData();
	var select = document.getElementById("StatusChange");
	formData.append("status", select.value);
	let config = {
	   headers: {"Content-Type": "multipart/form-data"}
	};
	axios.post("/change_status", formData, config).then(res => {
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
		url: '/check_notification'
	})
		.then(function(response) {
			var last = JSON.stringify(response.data);
			if(last !== "null"){
				var notification = window.eval(response.data);
				var title1 =　notification.title;
				var content1 ="　　" + notification.content;
				var time1 = notification.time;
				var body = document.getElementById("body");
				var div = document.createElement("div");
				body.appendChild(div);
				div.style.cssText = "width:35%;position: fixed; top: 10%; left: 36%; text-align: center; background-color: white;"
				div.id = "alert";

				var modal= document.createElement("div");
				div.appendChild(modal);
				modal.className="modal-dialog modal-content"
				modal.style.cssText="width:100%;"

				var modalHeader=document.createElement("div");
				modal.appendChild(modalHeader);
				modalHeader.className="modal-header"
				var title_1=document.createElement("h3");
				modalHeader.appendChild(title_1)
				title_1.innerHTML = title1;
				title_1.className="modal-title"

				var modalBody=document.createElement("h4");
				modal.appendChild(modalBody);
				modalBody.innerHTML = content1;
				modalBody.className="modal-body"
				modalBody.style.cssText="text-align: left;"

				var time=document.createElement("h4");
				modal.appendChild(time);
				time.innerHTML = time1;
				time.className="modal-body"
				time.style.cssText="text-align: right;"

				var foot=document.createElement("form");
				foot.className ="modal-footer";
				modal.appendChild(foot);
				var button = document.createElement("button");
				foot.appendChild(button);
				button.onclick= function () { read(); };
				button.className = "btn btn-primary";
				button.type="button";
				button.innerHTML = "已读";
				var button1 = document.createElement("button");
				foot.appendChild(button1);
				button1.onclick= function () { closeNotification(); };
				button1.className = "btn btn-default";
				button1.type="button";
				button1.innerHTML = "暂时忽略";
			}
		});
}

function showNotificationOrNot() {
	axios({
		method: 'get',
		url: '/check_notification'
	})
		.then(function (response) {
			var last = response.data;
			if(last !== null){
				document.getElementById("notification").style.display = "";
			}
			else{
				document.getElementById("notification").style.display = "none";
			}
		})
}

function read(){
	axios({
		method: 'post',
		url: '/read'
	});
	closeNotification(1);
	checkNotification();
	showNotificationOrNot();
}
function closeNotification(){
	var father = document.getElementById("body");
	var child = document.getElementById("alert");
	father.removeChild(child);
}
/*
function checkUser(){
	document.getElementById("logout").style.display="none";
	axios({
		method : 'get',
		url: '/check_group'
	})
		.then(function(response){
			if(response.data === "维修员"||response.data ==="业务员"||response.data ==="普通用户"){
				document.getElementById("register").style.display="none";
				document.getElementById("login").style.display="none";
				document.getElementById("logout").style.display="";
			}
		})
}
*/
function searchForParts(){
	delProjectResults();
	if (document.getElementById("detail1").value === "") return;
	axios({
		method: 'get',
		url: '/search_for_parts?text=' + document.getElementById("detail1").value + '&type=A'
	})
		.then(function (response) {
			{
				console.log(response.data);
				var div = document.getElementById("con");//寻找定位的区块
				var table = document.createElement("div");
				table.setAttribute("name", "table");
				table.className = "table-responsive";
				div.appendChild(table);
				var ordertable = document.createElement("table");//生成表格ordertable
				ordertable.id = "ordertable";
				ordertable.className = "table  table-bordered table-hover";
				table.appendChild(ordertable);
				var tbody = document.createElement("tbody");//不知道这是啥
				ordertable.appendChild(tbody);
				var tablehead = document.createElement("tr");//生成表格头tablehead
				tablehead.id = "tablehead";
				tbody.appendChild(tablehead);
				var head_1 = document.createElement("th");
				head_1.innerHTML = "维修项目名";
				tablehead.appendChild(head_1);
				var head_2 = document.createElement("th");
				head_2.innerHTML = "维修项目编号";
				tablehead.appendChild(head_2);

				for (var i = 0; response.data[i].Name !== ""; ++i) {
					var temp = document.createElement("tr");
					var t1 = document.createElement("td");
					t1.innerHTML = JSON.stringify(response.data[i].Name);
					temp.appendChild(t1);
					var t2 = document.createElement("td");
					t2.innerHTML = JSON.stringify(response.data[i].Id);
					temp.appendChild(t2);
					tbody.appendChild(temp);
				}
			}
		})
}

function getProjectTime(){
	axios({
		method: 'get',
		url: '/get_project_time?project=' + document.getElementById("project1").value + '&type=' +'A'
	})
		.then(function (response) {
			time = response.data.Time;
			document.getElementById("predict_finish_time1").value = time;
		})
}
function SsearchForProjects(){
 	if (document.getElementById("project1").value === "") return;
 	axios({
 		method: 'get',
 		url: '/search_for_projects?text=' + document.getElementById("project1").value + '&type=' +'A'
 	})
 		.then(function (response) {
 			{
				 var div = document.createElement("div");
				 div.id ="projects";
				 div.setAttribute("style", "width:150px; height:90px; position: absolute; top: 20px; left: 0; background-color: white; border: 1px solid black; cursor: pointer;");
				 document.getElementById("project_input").appendChild(div);
				 //table.setAttribute("style", "display: inline; width: 90px; position: absolute");
				 /*list.setAttribute("name", "list");
				 list.className = "table-responsive";
				 list.id = "li";
 				 div.appendChild(list);
				 var ordertable = document.createElement("table");//生成表格ordertable
 				 ordertable.id = "ordertable";
 				 ordertable.className = "table  table-bordered table-hover";
 				 table.appendChild(ordertable);
 				 var tbody = document.createElement("tbody");//不知道这是啥
 				 ordertable.appendChild(tbody);
				 var tablehead = document.createElement("tr");//生成表格头tablehead
 				 tablehead.id = "tablehead";
 				 tbody.appendChild(tablehead);
 				 var head_1 = document.createElement("th");
 				 head_1.innerHTML = "维修项目名";
 				 tablehead.appendChild(head_1);
 				 var head_2 = document.createElement("th");
 				 head_2.innerHTML = "维修项目编号";
 				 tablehead.appendChild(head_2);
				*/
 				for (var i = 1; response.data[i].Name !== ""; ++i) {
					 if(i === 6){
						 break;
					 }
					// var tr = document.createElement("tr");
					 //table.appendChild(tr);
					// tr.setAttribute("style","text-align: center; border: 1px solid #000");
					 var temp = document.createElement("div");
					 temp.class="select";
					 temp.setAttribute("onclick","fill(this);delProjectResults()");
					 temp.setAttribute("style","overflow:hidden;text-overflow:ellipsis;white-space:nowrap");
					 //var t1 = document.createElement("td");
					 //temp.setAttribute("value",JSON.stringify(response.data[i].Name).replace("\"","").replace("\"",""));
 					 temp.innerHTML = JSON.stringify(response.data[i].Name).replace("\"","").replace("\"","");
					 //temp.appendChild(t1);
 					 //var t2 = document.createElement("td");
 					 //t2.innerHTML = JSON.stringify(response.data[i].Id);
 					 //temp.appendChild(t2);
 					 //alert(temp);
 					 div.appendChild(temp);
 				}
 			}
 		})
 }

function delProjectResults(){
	try {
		var div = document.getElementById("projects");
			document.getElementById("project_input").removeChild(div);
	}catch (error)
	{}
}

var flagParts = 0;
var flagProjects = 0;
var timerParts;
var timerProjects;

function openFlagParts () {
	timerParts = setTimeout(function(){
		flagParts = 1;
		searchForParts();
		clearTimeout(timerParts);
		flagParts=0;
	}, 200);
}
function openFlagProjects () {
	    timerProjects = setTimeout(function(){
		flagProjects = 1;
		SsearchForProjects();
		clearTimeout(timerProjects);
		flagProjects=0;
	}, 200);
}
function closeFlagParts() {
	clearTimeout(timerParts);// 取消定时器
	flagParts = 0;
}
function closeFlagProjects() {
	clearTimeout(timerProjects);// 取消定时器
	flagProjects = 0;
}

function searchForParts(){
	delPartsResults();
	if(document.getElementById("txt").value === "") return;
	axios({
		method: 'get',
		url: '/search_for_parts?text='+document.getElementById("txt").value
	})
		.then(function(response) {
			{
				console.log(response.data);
				var div = document.getElementById("con");//寻找定位的区块
				var table = document.createElement("div");
				table.setAttribute("name", "table");
				table.className="table-responsive";
				div.appendChild(table);
				var ordertable = document.createElement("table");//生成表格ordertable
				ordertable.id = "ordertable";
				ordertable.className="table  table-bordered table-hover";
				table.appendChild(ordertable);
				var tbody = document.createElement("tbody");//不知道这是啥
				ordertable.appendChild(tbody);
				var tablehead = document.createElement("tr");//生成表格头tablehead
				tablehead.id = "tablehead";
				tbody.appendChild(tablehead);
				var head_1 = document.createElement("th");
				head_1.innerHTML="维修项目名";
				tablehead.appendChild(head_1);
				var head_2 = document.createElement("th");
				head_2.innerHTML="维修项目编号";
				tablehead.appendChild(head_2);

				for(var i=0; response.data[i].Name!==""; ++i) {
					var temp = document.createElement("tr");
					var t1 = document.createElement("td");
					t1.innerHTML = JSON.stringify(response.data[i].Name);
					temp.appendChild(t1);
					var t2 = document.createElement("td");
					t2.innerHTML = JSON.stringify(response.data[i].Id);
					temp.appendChild(t2);
					tbody.appendChild(temp);
				}
			}
		})
}

function delPartsResults(){
	try {
		var div = document.getElementById("con");
		var res = document.getElementsByName("table");
		for(var i=0; i<res.length; i++)
			div.removeChild(res[i]);
	}catch (error)
	{}
}
/*
function delProjectsResults(){
	try {
		var div = document.getElementById("project_input");
		var res = document.getElementById("project1");
		div.removeChild(res);
	}catch (error)
	{}
}
*/
/*
function getinfo(){
	axios({
		method : 'get',
		url: '/userinfo'
	})
	.then(function(response){
		var data = response.data
		for(var k in data ){//遍历packJson 对象的每个key/value对,k为key
			var info = document.getElementById("info");
			info.innerHTML = info.innerHTML + "<br>" + k +": " + data[k];
		}
	})
}
*/
function changeinfo(){
	let formData = new FormData();
	formData.append("name", document.getElementById("name").value);
	formData.append("property", document.getElementById("property").value);
	formData.append("contact_person", document.getElementById("contact_person").value);
	let config = {
	   headers: {"Content-Type": "multipart/form-data"}
	};
	axios.post("/change_userinfo", formData, config).then(res => {
		alert(res.data);
		location.reload();
	})
	document.getElementById("confirmChangeInfo").innerHTML="更改成功";
	window.setTimeout("window.location.reload()", 2000);
}

function addVehicle(){
	let formData = new FormData();
	formData.append("number", document.getElementById("number").value);
	formData.append("license_number", document.getElementById("license_number").value.toUpperCase());
	formData.append("color", document.getElementById("color").value);
	formData.append("model", document.getElementById("model").value);
	formData.append("type", document.getElementById("type").value);
	let config = {
		headers: {"Content-Type": "multipart/form-data"}
	};
	axios.post("/add_vehicle", formData, config).then(res => {
		var response = res.data;
		if(response.status === "失败")
		{
			alert(response.data);
		}
		else{
			alert("添加成功！");
		}
		location.reload(true);
	})
	// document.getElementById("confirmAddVehicle").innerHTML="新增成功";
	// window.setTimeout("window.location.reload()", 2000);
}

function checkVehicle() {
	document.getElementById("confirmAddVehicle").innerHTML="提交中...";
	axios.get('/check_vehicle',{
		params:{
			number: document.getElementById("number").value
		}
	})
		.then(function(response) {
			var last = response.data;
			if(last === 1){
				var r = confirm("该车辆已被绑定，是否覆盖？");
				if(r===true){
					addVehicle();
				}
			}
			else{
				addVehicle();
			}
		});
}
function getUserName(){
	axios({
		method : 'get',
		url: "/get_username"
	})
		.then(function(response){
			var text = document.getElementById("text");
			if(response.data==null){
				text.innerHTML ="用户名：未登录";

			}
			else{
				text.innerHTML ="用户名：" + response.data;
				document.getElementById("logout").style.display="";
			}
		})
}
function getUsername(){
	axios({
		method : 'get',
		url: "/get_username"
	})
		.then(function(response){
			var text1 = document.getElementById("text1");
			if(response.data===null){
				text1.innerHTML ="用户名：未登录";
			}
			else{
				text1.innerHTML ="当前用户名：" + response.data;
			}
		})
}

/*function getusername(){
	axios({
		method : 'get',
		url: "/get_username"
	})
		.then(function(response){
			var text2 = document.getElementById("navbar").getElementsByTagName("li");
			if(response.data===null){
				text2[0].innerHTML ="用户名：未登录";
			}
			else{
				text2[0].innerHTML ="用户名：" + response.data;
			}
		})
}*/