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
				modalBody.innerHTML = content1.replace(/\n/g,"<br>　　");
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

function searchForParts(){

	delPartsResults();
	if (document.getElementById("part").value === "") return;
	axios({
		method: 'get',
		url: '/search_for_parts?text=' + document.getElementById("part").value
	})
		.then(function (response) {
				var div = document.createElement("div");
				div.id ="parts";
				div.setAttribute("style", "width:150px; height:auto; position: absolute; top: 35px; left: 90px; background-color: white; border: 1px solid black; cursor: pointer;");
				document.getElementById("part_input").appendChild(div);
				for (var i = 0; response.data[i].Name !== ""; ++i) {
					if(i === 20){
						break;
					}
					var temp = document.createElement("div");
					temp.setAttribute("class", "select");
					temp.setAttribute("id",response.data[i].Id);
					temp.setAttribute("onclick","fill(this)");
					temp.setAttribute("style","overflow:hidden;text-overflow:ellipsis;white-space:nowrap");
					temp.innerHTML = JSON.stringify(response.data[i].Name).replace("\"","").replace("\"","");
					div.appendChild(temp);
				}
			})
		}

function getProjectTime(){
	axios({
		method: 'get',
		url: '/get_project_time?project=' + document.getElementById("project").value + '&type=' + document.getElementById("order_number").class
	})
		.then(function (response) {
			let time = response.data.Time;
			if(time === undefined){
				time = 0;
			}
			document.getElementById("predict_finish_time").value = time;
		})
}
function SsearchForProjects(){
	delProjectResults();
 	if (document.getElementById("project").value === ""){
		return;
	}
 	axios({
 		method: 'get',
 		url: '/search_for_projects?text=' + document.getElementById("project").value + '&type=' +document.getElementById("order_number").class
 	})
 		.then(function (response) {
 			{
				 if(response.data.length === 0){
					 return;
				 }
				 var div = document.createElement("div");
				 div.id ="projects";
				 div.setAttribute("style", "width:150px; height:auto; position: absolute; top: 20px; left: 0; background-color: white; border: 1px solid black; cursor: pointer;");
				 document.getElementById("project_input").appendChild(div);
				 var len = response.data.length;
 				 for (var i = 0; i < len; ++i) {
					  if(i === 20){
						  break;
					  }
					  var temp = document.createElement("div");
					  temp.setAttribute("class","select");
					  temp.setAttribute("id",response.data[i].Id);
					  temp.setAttribute("onclick","fill(this);delProjectResults()");
					  temp.setAttribute("style","overflow:hidden;text-overflow:ellipsis;white-space:nowrap;");
 					  temp.innerHTML = JSON.stringify(response.data[i].Name).replace("\"","").replace("\"","");
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

function delPartsResults(){
	try {
		var div = document.getElementById("parts");
		document.getElementById("part_input").removeChild(div);
	}catch (error)
	{}
}

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

function checkSubmit(param){
	var temp = document.getElementsByClassName(param);
	var num = temp.length;
	for(var i = 0; i < num; i++){
		if(temp[i].value === ""){
			alert("请把表单填写完整！");
			return false;
		}
	}
	return true;
}

function loadPDF(param)
{
	try {document.getElementById("pdf_div").removeChild(document.getElementById("show_pdf"));
	}
	catch(error){}
	var temp = param.parentNode.parentNode.children[0].innerText;
	var pdf = document.createElement("embed");
	var pdfDiv = document.getElementById("pdf_div");
	pdfDiv.appendChild(pdf);
	pdf.id = "show_pdf";
	pdf.setAttribute("height","840px");
	pdf.setAttribute("width","90%");
	pdf.setAttribute("src","/show_pdf?attorney_no="+temp);
}

function downloadPDF(param)
{
	var temp = param.parentNode.parentNode.children[0].innerText;
	window.open("/download_pdf?attorney_no="+temp,"_self");
}