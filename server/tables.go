package main

type User struct { //客户表
	Number        string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //客户编号
	Password      string `gorm:"type:varchar(255);not null"`                        //密码
	Name          string `gorm:"type:varchar(100)"`                                 //客户名称
	Property      string `gorm:"type:varchar(4)"`                                   //性质
	DiscountRate  int    `gorm:"type:int(3)"`                                       //折扣率
	ContactPerson string `gorm:"type:varchar(10)"`                                  //联系人
	ContactTel    string `gorm:"type:varchar(20);not null;unique_index"`            //联系电话
}

type Salesman struct { //业务员表
	Number   string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //工号
	Name     string `gorm:"type:varchar(20);not null"`                         //姓名
	Password string `gorm:"type:varchar(255);not null"`                        //密码
}

type Repairman struct { //维修员表
	Number          string  `gorm:"type:varchar(8);unique_index;not null;primary_key"` //工号
	Name            string  `gorm:"type:varchar(20);not null"`                         //姓名
	Password        string  `gorm:"type:varchar(255);not null"`                        //密码
	Type            string  `gorm:"type:varchar(10);not null"`                         //工种
	CurrentWorkHour float64 `gorm:"type:double(5,1);not null"`                         //当前工时
	Status          string  `gorm:"type:varchar(10);not null"`                         //工人状态
}

type TimeOverview struct { //工时总览表
	ProjectNumber   string  `gorm:"type:varchar(10);unique_index;not null;primary_key"` //维修项目编号
	ProjectName     string  `gorm:"type:varchar(50);not null"`                          //维修项目名称
	ProjectSpelling string  `gorm:"type:varchar(25);not null"`                          //项目名称首字母
	TimeA           float64 `gorm:"type:double(5,1)"`                                   //工时A
	TimeB           float64 `gorm:"type:double(5,1)"`                                   //工时B
	TimeC           float64 `gorm:"type:double(5,1)"`                                   //工时C
	TimeD           float64 `gorm:"type:double(5,1)"`                                   //工时D
	TimeE           float64 `gorm:"type:double(5,1)"`                                   //工时E
	Remark          string  `gorm:"type:varchar(20)"`                                   //备注
}

type PartsOverview struct { //零件总览表
	PartsNumber   string  `gorm:"type:varchar(25);unique_index;not null;primary_key"` //零件编号
	PartsName     string  `gorm:"type:varchar(70);not null"`                          //零件名称
	PartsSpelling string  `gorm:"type:varchar(35);not null"`                          //零件首字母
	PartsCost     float64 `gorm:"type:double(8,2);not null"`                          //零件价格
}

type Vehicle struct { //车辆表
	Number        string `gorm:"type:varchar(17);not null;primary_key"` //车架号
	LicenseNumber string `gorm:"type:varchar(10);not null"`             //车牌号
	UserID        string `gorm:"type:varchar(8);not null;"`             //客户编号
	Color         string `gorm:"type:varchar(10);not null"`             //车辆颜色
	Model         string `gorm:"type:varchar(40);not null"`             //车型
	Type          string `gorm:"type:varchar(10);not null"`             //车辆类别
	Time          string `gorm:"type:varchar(20);not null"`             //绑定时间
}

type Notification struct { //通知表
	Number  string `gorm:"type:varchar(8);unique_index;not null;primary_key"` //通知编号
	UserID  string `gorm:"type:varchar(11);not null"`                         //用户号
	Title   string `gorm:"type:varchar(50);not null"`                         //通知标题
	Content string `gorm:"type:varchar(255);not null"`                        //通知内容
	Status  string `gorm:"type:varchar(4);not null"`                          //接收状态
	Time    string `gorm:"type:varchar(20);not null"`                         //绑定时间
}

type Attorney struct { //委托书表
	Number            string  `gorm:"type:varchar(11);not null;unique_index;primary_key"` //订单号
	UserID            string  `gorm:"type:varchar(8);not null"`                           //客户编号
	VehicleNumber     string  `gorm:"type:varchar(17);not null"`                          //车架号（要填）
	RepairType        string  `gorm:"type:varchar(4)"`                                    //维修类型
	Classification    string  `gorm:"type:varchar(4)"`                                    //作业分类
	PayMethod         string  `gorm:"type:varchar(4);not null"`                           //结算方式（要填）
	StartTime         string  `gorm:"type:varchar(20);not null"`                          //进场时间（要填）
	SalesmanID        string  `gorm:"type:varchar(8)"`                                    //业务员编号
	PredictFinishTime string  `gorm:"type:varchar(20)"`                                   //预计完工时间时间
	ActualFinishTime  string  `gorm:"type:varchar(20)"`                                   //实际完工时间时间
	RoughProblem      string  `gorm:"type:varchar(255);not null"`                         //粗略故障描述（要填）
	SpecificProblem   string  `gorm:"type:varchar(255);"`                                 //详细故障描述
	Progress          string  `gorm:"type:varchar(10);not null"`                          //进展
	StartPetrol       float64 `gorm:"type:double(5,2);not null"`                          //进厂油量（要填）
	StartMile         float64 `gorm:"type:double(8,2);not null"`                          //进厂里程（要填）
	EndPetrol         float64 `gorm:"type:double(5,2)"`                                   //出厂油量
	EndMile           float64 `gorm:"type:double(8,2)"`                                   //出厂里程
	OutRange          string  `gorm:"type:varchar(255)"`                                  //非维修范围
}

type Arrangement struct { //派工单表
	OrderNumber     string `gorm:"type:varchar(11);not null;primary_key"` //订单号
	ProjectNumber   string `gorm:"type:varchar(10);not null;primary_key"` //维修项目编号
	RepairmanNumber string `gorm:"type:varchar(8);not null;primary_key"`  //维修工工号
	Progress        string `gorm:"type:varchar(6);not null"`              //进展
}

type RepairParts struct { //维修零件表
	OrderNumber   string `gorm:"type:varchar(11);not null;primary_key"` //订单号
	ProjectNumber string `gorm:"type:varchar(10);not null;primary_key"` //维修项目编号
	PartsNumber   string `gorm:"type:varchar(25);primary_key"`          //零件号
	PartsCount    int    `gorm:"type:int(2)"`                           //零件数量
}

type Remark struct { //登录表
	RemarkNumber string `gorm:"type:varchar(20);not null;unique_index;primary_key"` //备注编号
	Content      string `gorm:"type:varchar(255);not null"`                         //备注信息
}

type AuthSession struct { //登录表
	TimeHash  string `gorm:"type:varchar(64);not null;unique_index;primary_key"` //时间戳的哈希值
	LastVisit string `gorm:"type:varchar(30);not null"`                          //最后一次访问的时间戳（精确到秒）
	Username  string `gorm:"type:varchar(20);not null"`                          //当前session对应的用户信息
}
