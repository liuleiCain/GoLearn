package designpattern

import "fmt"

// ============================================================
// 创建型模式（Creational Patterns）
// ============================================================

// ------------------------------------------------------------
// 1. 单例模式（Singleton）
// ------------------------------------------------------------

// Singleton 单例结构体
type Singleton struct {
	name string
}

var (
	singletonInstance *Singleton
	singletonOnce     Once
)

// Once 简化的sync.Once实现
type Once struct {
	done bool
}

// GetSingleton 获取单例实例
func GetSingleton() *Singleton {
	if !singletonOnce.done {
		singletonInstance = &Singleton{name: "singleton"}
		singletonOnce.done = true
	}
	return singletonInstance
}

// SingletonDemo 演示单例模式
func SingletonDemo() {
	fmt.Println("=== 单例模式 ===")

	s1 := GetSingleton()
	s2 := GetSingleton()

	fmt.Printf("s1 == s2: %v\n", s1 == s2)
	fmt.Printf("s1地址: %p, s2地址: %p\n", s1, s2)
}

// ------------------------------------------------------------
// 2. 工厂模式（Factory）
// ------------------------------------------------------------

// Shape 形状接口
type Shape interface {
	Draw()
}

// Circle 圆形
type Circle struct{}

func (c *Circle) Draw() {
	fmt.Println("绘制圆形")
}

// Rectangle 矩形
type Rectangle struct{}

func (r *Rectangle) Draw() {
	fmt.Println("绘制矩形")
}

// Triangle 三角形
type Triangle struct{}

func (t *Triangle) Draw() {
	fmt.Println("绘制三角形")
}

// ShapeFactory 形状工厂
func ShapeFactory(shapeType string) Shape {
	switch shapeType {
	case "circle":
		return &Circle{}
	case "rectangle":
		return &Rectangle{}
	case "triangle":
		return &Triangle{}
	default:
		return nil
	}
}

// FactoryDemo 演示工厂模式
func FactoryDemo() {
	fmt.Println("\n=== 工厂模式 ===")

	circle := ShapeFactory("circle")
	circle.Draw()

	rectangle := ShapeFactory("rectangle")
	rectangle.Draw()

	triangle := ShapeFactory("triangle")
	triangle.Draw()
}

// ------------------------------------------------------------
// 3. 抽象工厂模式（Abstract Factory）
// ------------------------------------------------------------

// Button 按钮接口
type Button interface {
	Render()
}

// WindowsButton Windows按钮
type WindowsButton struct{}

func (b *WindowsButton) Render() {
	fmt.Println("渲染Windows风格按钮")
}

// MacOSButton MacOS按钮
type MacOSButton struct{}

func (b *MacOSButton) Render() {
	fmt.Println("渲染MacOS风格按钮")
}

// Checkbox 复选框接口
type Checkbox interface {
	Check()
}

// WindowsCheckbox Windows复选框
type WindowsCheckbox struct{}

func (c *WindowsCheckbox) Check() {
	fmt.Println("Windows复选框选中")
}

// MacOSCheckbox MacOS复选框
type MacOSCheckbox struct{}

func (c *MacOSCheckbox) Check() {
	fmt.Println("MacOS复选框选中")
}

// GUIFactory GUI工厂接口
type GUIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

// WindowsFactory Windows工厂
type WindowsFactory struct{}

func (f *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

func (f *WindowsFactory) CreateCheckbox() Checkbox {
	return &WindowsCheckbox{}
}

// MacOSFactory MacOS工厂
type MacOSFactory struct{}

func (f *MacOSFactory) CreateButton() Button {
	return &MacOSButton{}
}

func (f *MacOSFactory) CreateCheckbox() Checkbox {
	return &MacOSCheckbox{}
}

// AbstractFactoryDemo 演示抽象工厂模式
func AbstractFactoryDemo() {
	fmt.Println("\n=== 抽象工厂模式 ===")

	var factory GUIFactory

	// Windows风格
	factory = &WindowsFactory{}
	button := factory.CreateButton()
	checkbox := factory.CreateCheckbox()
	button.Render()
	checkbox.Check()

	// MacOS风格
	factory = &MacOSFactory{}
	button = factory.CreateButton()
	checkbox = factory.CreateCheckbox()
	button.Render()
	checkbox.Check()
}

// ------------------------------------------------------------
// 4. 建造者模式（Builder）
// ------------------------------------------------------------

// Computer 电脑
type Computer struct {
	CPU    string
	Memory string
	Disk   string
	GPU    string
}

// ComputerBuilder 电脑建造者
type ComputerBuilder struct {
	computer *Computer
}

// NewComputerBuilder 创建建造者
func NewComputerBuilder() *ComputerBuilder {
	return &ComputerBuilder{computer: &Computer{}}
}

// SetCPU 设置CPU
func (b *ComputerBuilder) SetCPU(cpu string) *ComputerBuilder {
	b.computer.CPU = cpu
	return b
}

// SetMemory 设置内存
func (b *ComputerBuilder) SetMemory(memory string) *ComputerBuilder {
	b.computer.Memory = memory
	return b
}

// SetDisk 设置硬盘
func (b *ComputerBuilder) SetDisk(disk string) *ComputerBuilder {
	b.computer.Disk = disk
	return b
}

// SetGPU 设置显卡
func (b *ComputerBuilder) SetGPU(gpu string) *ComputerBuilder {
	b.computer.GPU = gpu
	return b
}

// Build 构建电脑
func (b *ComputerBuilder) Build() *Computer {
	return b.computer
}

// BuilderDemo 演示建造者模式
func BuilderDemo() {
	fmt.Println("\n=== 建造者模式 ===")

	// 链式调用构建
	computer := NewComputerBuilder().
		SetCPU("Intel i7").
		SetMemory("32GB").
		SetDisk("1TB SSD").
		SetGPU("RTX 4090").
		Build()

	fmt.Printf("电脑配置: CPU=%s, Memory=%s, Disk=%s, GPU=%s\n",
		computer.CPU, computer.Memory, computer.Disk, computer.GPU)
}

// ------------------------------------------------------------
// 5. 原型模式（Prototype）
// ------------------------------------------------------------

// Prototype 原型接口
type Prototype interface {
	Clone() Prototype
}

// Document 文档
type Document struct {
	Title   string
	Content string
	Author  string
}

// Clone 克隆文档
func (d *Document) Clone() Prototype {
	return &Document{
		Title:   d.Title + " (副本)",
		Content: d.Content,
		Author:  d.Author,
	}
}

// PrototypeDemo 演示原型模式
func PrototypeDemo() {
	fmt.Println("\n=== 原型模式 ===")

	original := &Document{
		Title:   "设计模式文档",
		Content: "这是一份设计模式文档",
		Author:  "张三",
	}

	clone := original.Clone().(*Document)

	fmt.Printf("原始文档: %+v\n", original)
	fmt.Printf("克隆文档: %+v\n", clone)
}
