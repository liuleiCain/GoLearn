# receiver

#### 介绍
通过receiver实现对象方法，将方法绑定到对象
1. 方法名的首字母是否大写决定了该方法是不是导出方法
2. 方法定义与类型定义要放在同一个包中
3. 每个方法只能有一个receiver参数，不支持多个或变长的参数
4. receiver参数的基类型本身不能是指针类型或接口类型

#### receiver类型
1.  当receiver参数的类型为T时，选择值类型的receiver
2.  当receiver参数的类型为*T时，选择指针类型的receiver
