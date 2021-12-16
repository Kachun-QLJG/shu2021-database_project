document.writeln("<!--上侧区块-->");
document.writeln("<div class=\'container-fluid\'>");
document.writeln("    <div class=\'row\'>");
document.writeln("        <!--左下导航栏-->");
document.writeln("        <div class=\'hidden-xs  col-md-2\'>");
document.writeln("            <div class=\'list-group\'>");
document.writeln("                <a class=\'list-group-item list-group-item-success\'>");
document.writeln("                    功能");
document.writeln("                </a>");
document.writeln("                <a class=\'list-group-item\' href=\'/change_password\'>更改密码</a>");
document.writeln("                <a class=\'list-group-item\' href=\'/r_check_orders\' style=\'cursor:pointer\'>订单详情</a>");
document.writeln("                <a class=\'list-group-item\' data-toggle=\'modal\' data-target=\'#workStatus\' style=\'cursor:pointer\'>设置状态</a>");
document.writeln("            </div>");
document.writeln("        </div>");
document.writeln("        <!--右下导航栏-->");
document.writeln("        <div class=\'col-xs-10\' id=\'con\'>");
document.writeln("            <div class=\'container-fluid\'>");
document.writeln("                <div id=\'test_div\'>");
document.writeln("                </div>");
document.writeln("            </div>");
document.writeln("        </div>");
document.writeln("    </div>");
document.writeln("</div>");
document.writeln("");
document.writeln("<!-- 更改个人状态 -->");
document.writeln("<div class=\'modal fade\' id=\'workStatus\' tabindex=\'-1\' role=\'dialog\' aria-labelledby=\'myModalLabel\' aria-hidden=\'true\'>");
document.writeln("    <div class=\'modal-dialog\'>");
document.writeln("        <div class=\'modal-content\'>");
document.writeln("            <div class=\'modal-header\'>");
document.writeln("                <button type=\'button\' class=\'close\' data-dismiss=\'modal\' aria-hidden=\'true\'>&times;</button>");
document.writeln("                <h3 class=\'modal-title\' id=\'changeWorkStatus\'>更改状态</h3>");
document.writeln("            </div>");
document.writeln("            <div class=\'modal-body\'>");
document.writeln("                <img src=\'/statics/user/井号.png\' width=\'20px\'>");
document.writeln("                <label style=\'font-size:15px;width: 80px;\'>更改状态：</label>");
document.writeln("                <select name=\'status\'>");
document.writeln("                    <option value=\'正常\' selected>正常</option>");
document.writeln("                    <option value=\'休假\'>休假</option>");
document.writeln("                    <option value=\'离职\'>离职</option>");
document.writeln("                </select>");
document.writeln("            </div>");
document.writeln("            <form class=\'modal-footer\'>");
document.writeln("                <button type=\'button\' class=\'btn btn-primary\' onclick=\'changeStatus()\' id=\'confirmStatusChange\'>确认更改</button>");
document.writeln("                <button type=\'button\' class=\'btn btn-default\' data-dismiss=\'modal\'>取消</button>");
document.writeln("            </form>");
document.writeln("        </div>");
document.writeln("    </div>");
document.writeln("</div>");