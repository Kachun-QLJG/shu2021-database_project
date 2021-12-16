document.writeln("<div class=\'navbar-header\' id=\'head1\'>");
document.writeln("    <a class=\'navbar-brand\'>502汽修厂管理系统</a>");
document.writeln("</div>");
document.writeln("");
document.writeln("<!--右上区块-->");
document.writeln("<div class=\'collapse navbar-collapse\' id=\'bs-example-navbar-collapse-1\'>");
document.writeln("    <ul id=\'navbar\' class=\'nav navbar-nav navbar-right\' >");
document.writeln("        <li><a id=\'text\'>用户名：{{.username}}&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;用户组别：{{.group}}</a></li>");
document.writeln("        <li><a id=\'register\' href=\'/register\' title=\'注册\' >前往注册</a></li>");
document.writeln("        <li><a id=\'login\'href=\'/login\' title=\'登录\' >前往登录</a></li>");
document.writeln("        <li><a id=\'logout\' title=\'退出登录\' data-toggle=\'modal\' data-target=\'#log_out\'>退出登录</a></li>");
document.writeln("        <li><a href=\'/change_password\'id=\'changePassword\' title=\'更改密码\'>更改密码</a></li>");
document.writeln("    </ul>");
document.writeln("</div>");
document.writeln("");
document.writeln("<!--登出模态框（Modal）-->");
document.writeln("<div class=\'modal fade\' id=\'log_out\' tabindex=\'-1\' role=\'dialog\' aria-labelledby=\'myModalLabel\' aria-hidden=\'true\'>");
document.writeln("    <div class=\'modal-dialog\'>");
document.writeln("        <div class=\'modal-content\'>");
document.writeln("            <div class=\'modal-header\'>");
document.writeln("                <button type=\'button\' class=\'close\' data-dismiss=\'modal\' aria-hidden=\'true\'>&times;</button>");
document.writeln("                <h3 class=\'modal-title\' id=\'myModalLabel\'>您确定要登出当前账号吗？</h3>");
document.writeln("                <h4 class=\'modal-body\'>当前登录状态：{{.username}}&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;用户组别：{{.group}}</h4>");
document.writeln("            </div>");
document.writeln("            <form class=\'modal-footer\'>");
document.writeln("                <button type=\'submit\' value=\'退出登录\' class=\'btn btn-primary\' onclick=\'logout()\'>退出登录</button>");
document.writeln("                <button type=\'button\' class=\'btn btn-default\' data-dismiss=\'modal\'>取消</button>");
document.writeln("            </form>");
document.writeln("        </div>");
document.writeln("    </div>");
document.writeln("</div>");