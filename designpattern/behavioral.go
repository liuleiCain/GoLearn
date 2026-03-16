package designpattern

import "fmt"

// ============================================================
// 行为型模式（Behavioral Patterns）
// ============================================================

// ------------------------------------------------------------
// 1. 策略模式（Strategy）
// ------------------------------------------------------------

// PaymentStrategy 支付策略接口
type PaymentStrategy interface {
	Pay(amount float64)
}

// CreditCardStrategy 信用卡支付策略
type CreditCardStrategy struct {
	name   string
	card   string
	cvv    string
	date   string
}

func NewCreditCardStrategy(name, card, cvv, date string) *CreditCardStrategy {
	return &CreditCardStrategy{name: name, card: card, cvv: cvv, date: date}
}

func (s *CreditCardStrategy) Pay(amount float64) {
	fmt.Printf("使用信用卡支付 %.2f 元\n", amount)
}

// AlipayStrategy 支付宝支付策略
type AlipayStrategy struct {
	email string
}

func NewAlipayStrategy(email string) *AlipayStrategy {
	return &AlipayStrategy{email: email}
}

func (s *AlipayStrategy) Pay(amount float64) {
	fmt.Printf("使用支付宝支付 %.2f 元\n", amount)
}

// WechatPayStrategy 微信支付策略
type WechatPayStrategy struct {
	phone string
}

func NewWechatPayStrategy(phone string) *WechatPayStrategy {
	return &WechatPayStrategy{phone: phone}
}

func (s *WechatPayStrategy) Pay(amount float64) {
	fmt.Printf("使用微信支付 %.2f 元\n", amount)
}

// ShoppingCart 购物车
type ShoppingCart struct {
	items   []string
	strategy PaymentStrategy
}

func (c *ShoppingCart) AddItem(item string) {
	c.items = append(c.items, item)
}

func (c *ShoppingCart) SetStrategy(strategy PaymentStrategy) {
	c.strategy = strategy
}

func (c *ShoppingCart) Checkout() {
	amount := float64(len(c.items) * 100)
	c.strategy.Pay(amount)
}

// StrategyDemo 演示策略模式
func StrategyDemo() {
	fmt.Println("=== 策略模式 ===")

	cart := &ShoppingCart{}
	cart.AddItem("商品1")
	cart.AddItem("商品2")

	cart.SetStrategy(NewCreditCardStrategy("张三", "1234567890123456", "123", "12/25"))
	cart.Checkout()

	cart.SetStrategy(NewAlipayStrategy("zhangsan@example.com"))
	cart.Checkout()

	cart.SetStrategy(NewWechatPayStrategy("13800138000"))
	cart.Checkout()
}

// ------------------------------------------------------------
// 2. 观察者模式（Observer）
// ------------------------------------------------------------

// Observer 观察者接口
type Observer interface {
	Update(message string)
}

// Subject 主题
type Subject struct {
	observers []Observer
}

func (s *Subject) Attach(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *Subject) Detach(observer Observer) {
	for i, o := range s.observers {
		if o == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *Subject) Notify(message string) {
	for _, observer := range s.observers {
		observer.Update(message)
	}
}

// UserObserver 用户观察者
type UserObserver struct {
	name string
}

func (o *UserObserver) Update(message string) {
	fmt.Printf("%s 收到消息: %s\n", o.name, message)
}

// ObserverDemo 演示观察者模式
func ObserverDemo() {
	fmt.Println("\n=== 观察者模式 ===")

	subject := &Subject{}

	user1 := &UserObserver{name: "用户1"}
	user2 := &UserObserver{name: "用户2"}
	user3 := &UserObserver{name: "用户3"}

	subject.Attach(user1)
	subject.Attach(user2)
	subject.Attach(user3)

	subject.Notify("系统维护通知")

	subject.Detach(user2)
	subject.Notify("维护完成")
}

// ------------------------------------------------------------
// 3. 命令模式（Command）
// ------------------------------------------------------------

// Command 命令接口
type Command interface {
	Execute()
	Undo()
}

// Light 灯
type Light struct {
	isOn bool
}

func (l *Light) On() {
	l.isOn = true
	fmt.Println("灯已打开")
}

func (l *Light) Off() {
	l.isOn = false
	fmt.Println("灯已关闭")
}

// LightOnCommand 开灯命令
type LightOnCommand struct {
	light *Light
}

func NewLightOnCommand(light *Light) *LightOnCommand {
	return &LightOnCommand{light: light}
}

func (c *LightOnCommand) Execute() {
	c.light.On()
}

func (c *LightOnCommand) Undo() {
	c.light.Off()
}

// LightOffCommand 关灯命令
type LightOffCommand struct {
	light *Light
}

func NewLightOffCommand(light *Light) *LightOffCommand {
	return &LightOffCommand{light: light}
}

func (c *LightOffCommand) Execute() {
	c.light.Off()
}

func (c *LightOffCommand) Undo() {
	c.light.On()
}

// RemoteControl 遥控器
type RemoteControl struct {
	command Command
}

func (r *RemoteControl) SetCommand(command Command) {
	r.command = command
}

func (r *RemoteControl) PressButton() {
	r.command.Execute()
}

func (r *RemoteControl) PressUndo() {
	r.command.Undo()
}

// CommandDemo 演示命令模式
func CommandDemo() {
	fmt.Println("\n=== 命令模式 ===")

	light := &Light{}
	remote := &RemoteControl{}

	remote.SetCommand(NewLightOnCommand(light))
	remote.PressButton()
	remote.PressUndo()

	remote.SetCommand(NewLightOffCommand(light))
	remote.PressButton()
	remote.PressUndo()
}

// ------------------------------------------------------------
// 4. 责任链模式（Chain of Responsibility）
// ------------------------------------------------------------

// Handler 处理器接口
type Handler interface {
	SetNext(handler Handler)
	Handle(request string) string
}

// AbstractHandler 抽象处理器
type AbstractHandler struct {
	next Handler
}

func (h *AbstractHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *AbstractHandler) HandleNext(request string) string {
	if h.next != nil {
		return h.next.Handle(request)
	}
	return ""
}

// ManagerHandler 经理处理器
type ManagerHandler struct {
	AbstractHandler
}

func (h *ManagerHandler) Handle(request string) string {
	if request == "请假1天" {
		return "经理批准: " + request
	}
	return h.HandleNext(request)
}

// DirectorHandler 总监处理器
type DirectorHandler struct {
	AbstractHandler
}

func (h *DirectorHandler) Handle(request string) string {
	if request == "请假3天" {
		return "总监批准: " + request
	}
	return h.HandleNext(request)
}

// CEOHandler CEO处理器
type CEOHandler struct {
	AbstractHandler
}

func (h *CEOHandler) Handle(request string) string {
	if request == "请假7天" {
		return "CEO批准: " + request
	}
	return "无人能批准: " + request
}

// ChainDemo 演示责任链模式
func ChainDemo() {
	fmt.Println("\n=== 责任链模式 ===")

	manager := &ManagerHandler{}
	director := &DirectorHandler{}
	ceo := &CEOHandler{}

	manager.SetNext(director)
	director.SetNext(ceo)

	fmt.Println(manager.Handle("请假1天"))
	fmt.Println(manager.Handle("请假3天"))
	fmt.Println(manager.Handle("请假7天"))
	fmt.Println(manager.Handle("请假30天"))
}

// ------------------------------------------------------------
// 5. 模板方法模式（Template Method）
// ------------------------------------------------------------

// DataMiner 数据挖掘器
type DataMiner interface {
	OpenFile()
	ExtractData()
	ParseData()
	AnalyzeData()
	SendReport()
	CloseFile()
}

// DataMinerTemplate 数据挖掘模板
type DataMinerTemplate struct{}

func (t *DataMinerTemplate) Mine(miner DataMiner) {
	miner.OpenFile()
	miner.ExtractData()
	miner.ParseData()
	miner.AnalyzeData()
	miner.SendReport()
	miner.CloseFile()
}

// PDFMiner PDF挖掘器
type PDFMiner struct {
	DataMinerTemplate
}

func (m *PDFMiner) OpenFile()   { fmt.Println("打开PDF文件") }
func (m *PDFMiner) ExtractData() { fmt.Println("提取PDF数据") }
func (m *PDFMiner) ParseData()   { fmt.Println("解析PDF数据") }
func (m *PDFMiner) AnalyzeData() { fmt.Println("分析PDF数据") }
func (m *PDFMiner) SendReport()  { fmt.Println("发送PDF报告") }
func (m *PDFMiner) CloseFile()   { fmt.Println("关闭PDF文件") }

// CSVMiner CSV挖掘器
type CSVMiner struct {
	DataMinerTemplate
}

func (m *CSVMiner) OpenFile()   { fmt.Println("打开CSV文件") }
func (m *CSVMiner) ExtractData() { fmt.Println("提取CSV数据") }
func (m *CSVMiner) ParseData()   { fmt.Println("解析CSV数据") }
func (m *CSVMiner) AnalyzeData() { fmt.Println("分析CSV数据") }
func (m *CSVMiner) SendReport()  { fmt.Println("发送CSV报告") }
func (m *CSVMiner) CloseFile()   { fmt.Println("关闭CSV文件") }

// TemplateDemo 演示模板方法模式
func TemplateDemo() {
	fmt.Println("\n=== 模板方法模式 ===")

	template := &DataMinerTemplate{}

	pdfMiner := &PDFMiner{}
	template.Mine(pdfMiner)

	fmt.Println()

	csvMiner := &CSVMiner{}
	template.Mine(csvMiner)
}

// ------------------------------------------------------------
// 6. 迭代器模式（Iterator）
// ------------------------------------------------------------

// Iterator 迭代器接口
type Iterator interface {
	HasNext() bool
	Next() interface{}
}

// StringIterator 字符串迭代器
type StringIterator struct {
	items []string
	index int
}

func NewStringIterator(items []string) *StringIterator {
	return &StringIterator{items: items}
}

func (i *StringIterator) HasNext() bool {
	return i.index < len(i.items)
}

func (i *StringIterator) Next() interface{} {
	if i.HasNext() {
		item := i.items[i.index]
		i.index++
		return item
	}
	return nil
}

// StringCollection 字符串集合
type StringCollection struct {
	items []string
}

func NewStringCollection(items []string) *StringCollection {
	return &StringCollection{items: items}
}

func (c *StringCollection) CreateIterator() Iterator {
	return NewStringIterator(c.items)
}

// IteratorDemo 演示迭代器模式
func IteratorDemo() {
	fmt.Println("\n=== 迭代器模式 ===")

	collection := NewStringCollection([]string{"A", "B", "C", "D"})
	iterator := collection.CreateIterator()

	for iterator.HasNext() {
		fmt.Println(iterator.Next())
	}
}

// ------------------------------------------------------------
// 7. 状态模式（State）
// ------------------------------------------------------------

// State 状态接口
type State interface {
	Handle()
}

// Order 订单
type Order struct {
	state State
}

func (o *Order) SetState(state State) {
	o.state = state
}

func (o *Order) Process() {
	o.state.Handle()
}

// NewOrderState 新订单状态
type NewOrderState struct{}

func (s *NewOrderState) Handle() {
	fmt.Println("处理新订单")
}

// PaidOrderState 已支付状态
type PaidOrderState struct{}

func (s *PaidOrderState) Handle() {
	fmt.Println("订单已支付，准备发货")
}

// ShippedOrderState 已发货状态
type ShippedOrderState struct{}

func (s *ShippedOrderState) Handle() {
	fmt.Println("订单已发货，等待收货")
}

// CompletedOrderState 已完成状态
type CompletedOrderState struct{}

func (s *CompletedOrderState) Handle() {
	fmt.Println("订单已完成")
}

// StateDemo 演示状态模式
func StateDemo() {
	fmt.Println("\n=== 状态模式 ===")

	order := &Order{}

	order.SetState(&NewOrderState{})
	order.Process()

	order.SetState(&PaidOrderState{})
	order.Process()

	order.SetState(&ShippedOrderState{})
	order.Process()

	order.SetState(&CompletedOrderState{})
	order.Process()
}

// ------------------------------------------------------------
// 8. 备忘录模式（Memento）
// ------------------------------------------------------------

// Memento 备忘录
type Memento struct {
	state string
}

func (m *Memento) GetState() string {
	return m.state
}

// Originator 原发器
type Originator struct {
	state string
}

func (o *Originator) SetState(state string) {
	o.state = state
}

func (o *Originator) Save() *Memento {
	return &Memento{state: o.state}
}

func (o *Originator) Restore(memento *Memento) {
	o.state = memento.GetState()
}

// Caretaker 管理者
type Caretaker struct {
	mementos []*Memento
}

func (c *Caretaker) Add(memento *Memento) {
	c.mementos = append(c.mementos, memento)
}

func (c *Caretaker) Get(index int) *Memento {
	return c.mementos[index]
}

// MementoDemo 演示备忘录模式
func MementoDemo() {
	fmt.Println("\n=== 备忘录模式 ===")

	originator := &Originator{}
	caretaker := &Caretaker{}

	originator.SetState("状态1")
	fmt.Println("当前状态:", originator.state)
	caretaker.Add(originator.Save())

	originator.SetState("状态2")
	fmt.Println("当前状态:", originator.state)
	caretaker.Add(originator.Save())

	originator.SetState("状态3")
	fmt.Println("当前状态:", originator.state)

	originator.Restore(caretaker.Get(0))
	fmt.Println("恢复后状态:", originator.state)
}
