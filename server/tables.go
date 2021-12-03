package main

type User struct { //客户表
	Number        string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //客户编号
	Password      string `gorm:"type:varchar(255);not null"`                        //密码
	Name          string `gorm:"type:varchar(100)"`                                 //客户名称
	Property      string `gorm:"type:varchar(4)"`                                   //性质
	DiscountRate  int    `gorm:"type:int(2)"`                                       //折扣率
	ContactPerson string `gorm:"type:varchar(10)"`                                  //联系人
	ContactTel    string `gorm:"type:varchar(20);not null;unique_index"`            //联系电话
}

type Salesman struct { //业务员表
	Number   string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //工号
	Name     string `gorm:"type:varchar(20);not null"`                         //姓名
	Password string `gorm:"type:varchar(255);not null"`                        //密码
}

type Repairman struct { //维修员表
	Number          string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //工号
	Name            string `gorm:"type:varchar(20);not null"`                         //姓名
	Password        string `gorm:"type:varchar(255);not null"`                        //密码
	Type            string `gorm:"type:varchar(10);not null"`                         //工种
	CurrentWorkHour int    `gorm:"type:int(4);not null"`                              //当前工时
	Status          string `gorm:"type:varchar(10);not null"`                         //工人状态
}

type TypeOverview struct { //工种总览表
	ProjectNumber   string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //维修项目编号
	ProjectName     string `gorm:"type:varchar(100);not null"`                        //维修项目名称
	Type            string `gorm:"type:varchar(10);not null;primary_key"`             //工种
	ProjectSpelling string `gorm:"type:varchar(50);not null;primary_key"`             //项目名称首字母
}

type PartsOverview struct { //零件总览表
	PartsNumber string  `gorm:"type:varchar(8);unique_index;not null;primary_key"` //零件编号
	PartsName   string  `gorm:"type:varchar(50);not null"`                         //零件名称
	Unit        string  `gorm:"type:varchar(6);not null"`                          //计量单位
	PartsCost   float64 `gorm:"type:double(8,2);not null"`                         //零件价格
}

type Vehicle struct { //车辆表
	Number        string `gorm:"type:varchar(17);not null;primary_key"` //车架号
	LicenseNumber string `gorm:"type:varchar(10);not null"`             //车牌号
	UserID        string `gorm:"type:varchar(8);not null;primary_key"`  //客户编号
	Color         string `gorm:"type:varchar(10);not null"`             //车辆颜色
	Model         string `gorm:"type:varchar(40);not null"`             //车型
	Type          string `gorm:"type:varchar(10);not null"`             //车辆类别
	Time          string `gorm:"type:varchar(20);not null"`             //绑定时间
}

type Notification struct { //通知表
	Number  string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //通知编号
	UserID  string `gorm:"type:varchar(8);not null"`                          //用户号
	Title   string `gorm:"type:varchar(50);not null"`                         //通知标题
	Content string `gorm:"type:varchar(255);not null"`                        //通知内容
	Status  string `gorm:"type:varchar(4);not null"`                          //接收状态
}

type Attorney struct { //委托书表
	Number            string  `gorm:"type:varchar(11);not null;unique_index;primary_key"` //订单号
	UserID            string  `gorm:"type:varchar(8);not null"`                           //客户编号
	VehicleNumber     string  `gorm:"type:varchar(17);not null"`                          //车架号
	RepairType        string  `gorm:"type:varchar(4)"`                                    //维修类型
	Classification    string  `gorm:"type:varchar(4)"`                                    //作业分类
	PayMethod         string  `gorm:"type:varchar(4)"`                                    //结算方式
	StartTime         string  `gorm:"type:varchar(20);not null"`                          //进场时间
	SalesmanID        string  `gorm:"type:varchar(8)"`                                    //业务员编号
	PredictFinishTime string  `gorm:"type:varchar(20)"`                                   //预计完工时间时间
	ActualFinishTime  string  `gorm:"type:varchar(20)"`                                   //实际完工时间时间
	RoughProblem      string  `gorm:"type:varchar(255);not null"`                         //粗略故障描述
	SpecificProblem   string  `gorm:"type:varchar(255);not null"`                         //详细故障描述
	Progress          string  `gorm:"type:varchar(10);not null"`                          //进展
	TotalCost         float64 `gorm:"type:double(6,2);not null"`                          //总价
	StartPetrol       float64 `gorm:"type:double(5,2);not null"`                          //进厂油量
	StartMile         float64 `gorm:"type:double(8,2);not null"`                          //进厂里程
	EndPetrol         float64 `gorm:"type:double(5,2)"`                                   //出厂油量
	EndMile           float64 `gorm:"type:double(8,2)"`                                   //出厂里程
	OutRange          string  `gorm:"type:varchar(255)"`                                  //非维修范围
}

type Arrangement struct { //派工单表
	OrderNumber     string `gorm:"type:varchar(11);not null;primary_key"` //订单号
	ProjectNumber   string `gorm:"type:varchar(8);not null;primary_key"`  //维修项目编号
	PredictTime     int    `gorm:"type:int(3);not null"`                  //预计工时
	ActualTIme      int    `gorm:"type:int(3);not null"`                  //实际工时
	RepairmanNumber string `gorm:"type:varchar(8);not null;primary_key"`  //维修工工号
	PartsNumber     string `gorm:"type:varchar(8);not null"`              //零件号
	PartsCount      int    `gorm:"type:int(2);not null"`                  //零件数量
	Progress        string `gorm:"type:varchar(6);not null"`              //进展
}

type AuthSession struct { //登录表
	TimeHash  string `gorm:"type:varchar(64);not null;unique_index;primary_key"` //时间戳的哈希值
	LastVisit string `gorm:"type:varchar(30);not null"`                          //最后一次访问的时间戳（精确到秒）
	Username  string `gorm:"type:varchar(20);not null"`                          //当前session对应的用户信息
}
