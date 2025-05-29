# defer

#### 介绍
defer常见的使用常见有资源释放、异常拦截、函数返回值处理

#### defer的执行顺序
1.  defer执行顺序为后进先出，是因为底层使用栈
2.  defer在return之前执行，若出现panic，defer执行完成后执行panic，panic之后的defer将不被执行
