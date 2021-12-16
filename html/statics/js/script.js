function displayChangeStatus() {
	document.getElementById("register").style.display="none";
	var text = document.getElementById("text");
	axios({
		method : 'get',
		url: '/checkStatus'
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

function checkGroup(){
	axios({
		method : 'get',
		url: '/checkGroup'
	})
	.then(function(response){
		if(response.data === "维修员"){
			displayChangeStatus()
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

function searchForProjects(){
	delProjectResults();
	if (document.getElementById("txt").value === "") return;
	axios({
		method: 'get',
		url: '/search_for_projects?text=' + document.getElementById("txt").value + '&type=A'
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

function delProjectResults(){
	try {
		var div = document.getElementById("con");
		var res = document.getElementsByName("table");
		for(var i=0; i<res.length; i++)
			div.removeChild(res[i]);
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
		searchForProjects();
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

function changeinfo(){
	let formData = new FormData();
	formData.append("name", "郑宇");
	formData.append("property", "个人");
	formData.append("contact_person", "郑宇");
	let config = {
	   headers: {"Content-Type": "multipart/form-data"}
	};
	axios.post("/change_userinfo", formData, config).then(res => {
				 alert(res.data);
				 location.reload();
			})
}

function addVehicle(){
	let formData = new FormData();
	formData.append("number", document.getElementById("number").value);
	formData.append("license_number", document.getElementById("license_number").value);
	formData.append("color", document.getElementById("color").value);
	formData.append("model", document.getElementById("model").value);
	formData.append("type", document.getElementById("type").value);
	let config = {
		headers: {"Content-Type": "multipart/form-data"}
	};
	axios.post("/addVehicle", formData, config).then(res => {
		alert(res.data);
		location.reload();
	})
}
