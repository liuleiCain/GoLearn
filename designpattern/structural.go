package designpattern

import "fmt"

// ============================================================
// 结构型模式（Structural Patterns）
// ============================================================

// ------------------------------------------------------------
// 1. 适配器模式（Adapter）
// ------------------------------------------------------------

// MediaPlayer 媒体播放器接口
type MediaPlayer interface {
	Play(audioType string, fileName string)
}

// AdvancedMediaPlayer 高级媒体播放器接口
type AdvancedMediaPlayer interface {
	PlayVlc(fileName string)
	PlayMp4(fileName string)
}

// VlcPlayer VLC播放器
type VlcPlayer struct{}

func (p *VlcPlayer) PlayVlc(fileName string) {
	fmt.Println("播放VLC文件:", fileName)
}

func (p *VlcPlayer) PlayMp4(fileName string) {}

// Mp4Player MP4播放器
type Mp4Player struct{}

func (p *Mp4Player) PlayVlc(fileName string) {}

func (p *Mp4Player) PlayMp4(fileName string) {
	fmt.Println("播放MP4文件:", fileName)
}

// MediaAdapter 媒体适配器
type MediaAdapter struct {
	advancedPlayer AdvancedMediaPlayer
}

// NewMediaAdapter 创建媒体适配器
func NewMediaAdapter(audioType string) *MediaAdapter {
	var player AdvancedMediaPlayer
	switch audioType {
	case "vlc":
		player = &VlcPlayer{}
	case "mp4":
		player = &Mp4Player{}
	}
	return &MediaAdapter{advancedPlayer: player}
}

// Play 播放
func (a *MediaAdapter) Play(audioType string, fileName string) {
	switch audioType {
	case "vlc":
		a.advancedPlayer.PlayVlc(fileName)
	case "mp4":
		a.advancedPlayer.PlayMp4(fileName)
	}
}

// AudioPlayer 音频播放器
type AudioPlayer struct{}

func (p *AudioPlayer) Play(audioType string, fileName string) {
	if audioType == "mp3" {
		fmt.Println("播放MP3文件:", fileName)
	} else if audioType == "vlc" || audioType == "mp4" {
		adapter := NewMediaAdapter(audioType)
		adapter.Play(audioType, fileName)
	} else {
		fmt.Println("不支持的格式:", audioType)
	}
}

// AdapterDemo 演示适配器模式
func AdapterDemo() {
	fmt.Println("=== 适配器模式 ===")

	player := &AudioPlayer{}
	player.Play("mp3", "song.mp3")
	player.Play("vlc", "movie.vlc")
	player.Play("mp4", "video.mp4")
	player.Play("avi", "video.avi")
}

// ------------------------------------------------------------
// 2. 装饰器模式（Decorator）
// ------------------------------------------------------------

// Coffee 咖啡接口
type Coffee interface {
	GetCost() int
	GetDescription() string
}

// SimpleCoffee 简单咖啡
type SimpleCoffee struct{}

func (c *SimpleCoffee) GetCost() int {
	return 10
}

func (c *SimpleCoffee) GetDescription() string {
	return "简单咖啡"
}

// CoffeeDecorator 咖啡装饰器
type CoffeeDecorator struct {
	coffee Coffee
}

func (d *CoffeeDecorator) GetCost() int {
	return d.coffee.GetCost()
}

func (d *CoffeeDecorator) GetDescription() string {
	return d.coffee.GetDescription()
}

// MilkDecorator 牛奶装饰器
type MilkDecorator struct {
	CoffeeDecorator
}

func NewMilkDecorator(coffee Coffee) *MilkDecorator {
	return &MilkDecorator{CoffeeDecorator{coffee: coffee}}
}

func (d *MilkDecorator) GetCost() int {
	return d.coffee.GetCost() + 5
}

func (d *MilkDecorator) GetDescription() string {
	return d.coffee.GetDescription() + " + 牛奶"
}

// SugarDecorator 糖装饰器
type SugarDecorator struct {
	CoffeeDecorator
}

func NewSugarDecorator(coffee Coffee) *SugarDecorator {
	return &SugarDecorator{CoffeeDecorator{coffee: coffee}}
}

func (d *SugarDecorator) GetCost() int {
	return d.coffee.GetCost() + 2
}

func (d *SugarDecorator) GetDescription() string {
	return d.coffee.GetDescription() + " + 糖"
}

// DecoratorDemo 演示装饰器模式
func DecoratorDemo() {
	fmt.Println("\n=== 装饰器模式 ===")

	coffee := &SimpleCoffee{}
	fmt.Printf("%s: %d元\n", coffee.GetDescription(), coffee.GetCost())

	coffeeWithMilk := NewMilkDecorator(coffee)
	fmt.Printf("%s: %d元\n", coffeeWithMilk.GetDescription(), coffeeWithMilk.GetCost())

	coffeeWithMilkAndSugar := NewSugarDecorator(coffeeWithMilk)
	fmt.Printf("%s: %d元\n", coffeeWithMilkAndSugar.GetDescription(), coffeeWithMilkAndSugar.GetCost())
}

// ------------------------------------------------------------
// 3. 代理模式（Proxy）
// ------------------------------------------------------------

// Image 图像接口
type Image interface {
	Display()
}

// RealImage 真实图像
type RealImage struct {
	fileName string
}

func NewRealImage(fileName string) *RealImage {
	img := &RealImage{fileName: fileName}
	img.loadFromDisk()
	return img
}

func (i *RealImage) loadFromDisk() {
	fmt.Println("从磁盘加载图像:", i.fileName)
}

func (i *RealImage) Display() {
	fmt.Println("显示图像:", i.fileName)
}

// ImageProxy 图像代理
type ImageProxy struct {
	realImage *RealImage
	fileName  string
}

func NewImageProxy(fileName string) *ImageProxy {
	return &ImageProxy{fileName: fileName}
}

func (p *ImageProxy) Display() {
	if p.realImage == nil {
		p.realImage = NewRealImage(p.fileName)
	}
	p.realImage.Display()
}

// ProxyDemo 演示代理模式
func ProxyDemo() {
	fmt.Println("\n=== 代理模式 ===")

	image := NewImageProxy("photo.jpg")

	fmt.Println("图像创建完成，但尚未加载")
	image.Display()
	fmt.Println("再次显示图像（无需重新加载）")
	image.Display()
}

// ------------------------------------------------------------
// 4. 组合模式（Composite）
// ------------------------------------------------------------

// FileSystemNode 文件系统节点接口
type FileSystemNode interface {
	Display(indent string)
}

// File 文件
type File struct {
	name string
}

func (f *File) Display(indent string) {
	fmt.Println(indent + "- " + f.name)
}

// Directory 目录
type Directory struct {
	name     string
	children []FileSystemNode
}

func (d *Directory) Add(node FileSystemNode) {
	d.children = append(d.children, node)
}

func (d *Directory) Display(indent string) {
	fmt.Println(indent + "+ " + d.name)
	for _, child := range d.children {
		child.Display(indent + "  ")
	}
}

// CompositeDemo 演示组合模式
func CompositeDemo() {
	fmt.Println("\n=== 组合模式 ===")

	root := &Directory{name: "root"}
	src := &Directory{name: "src"}
	bin := &Directory{name: "bin"}

	root.Add(src)
	root.Add(bin)

	src.Add(&File{name: "main.go"})
	src.Add(&File{name: "utils.go"})
	bin.Add(&File{name: "app.exe"})

	root.Display("")
}

// ------------------------------------------------------------
// 5. 外观模式（Facade）
// ------------------------------------------------------------

// CPU CPU
type CPU struct{}

func (c *CPU) Freeze() {
	fmt.Println("CPU冻结")
}

func (c *CPU) Jump(position int) {
	fmt.Printf("CPU跳转到位置 %d\n", position)
}

func (c *CPU) Execute() {
	fmt.Println("CPU执行")
}

// Memory 内存
type Memory struct{}

func (m *Memory) Load(position int, data string) {
	fmt.Printf("内存加载数据到位置 %d: %s\n", position, data)
}

// HardDrive 硬盘
type HardDrive struct{}

func (h *HardDrive) Read(lba int, size int) string {
	fmt.Printf("硬盘读取: LBA=%d, 大小=%d\n", lba, size)
	return "boot_data"
}

// ComputerFacade 电脑外观
type ComputerFacade struct {
	cpu       *CPU
	memory    *Memory
	hardDrive *HardDrive
}

func NewComputerFacade() *ComputerFacade {
	return &ComputerFacade{
		cpu:       &CPU{},
		memory:    &Memory{},
		hardDrive: &HardDrive{},
	}
}

func (f *ComputerFacade) Start() {
	fmt.Println("启动电脑...")
	f.cpu.Freeze()
	f.memory.Load(0, f.hardDrive.Read(0, 1024))
	f.cpu.Jump(0)
	f.cpu.Execute()
	fmt.Println("电脑启动完成")
}

// FacadeDemo 演示外观模式
func FacadeDemo() {
	fmt.Println("\n=== 外观模式 ===")

	computer := NewComputerFacade()
	computer.Start()
}

// ------------------------------------------------------------
// 6. 桥接模式（Bridge）
// ------------------------------------------------------------

// Renderer 渲染器接口
type Renderer interface {
	RenderCircle(radius float64)
}

// VectorRenderer 矢量渲染器
type VectorRenderer struct{}

func (r *VectorRenderer) RenderCircle(radius float64) {
	fmt.Printf("矢量渲染圆形，半径: %.1f\n", radius)
}

// RasterRenderer 光栅渲染器
type RasterRenderer struct{}

func (r *RasterRenderer) RenderCircle(radius float64) {
	fmt.Printf("光栅渲染圆形，半径: %.1f\n", radius)
}

// BridgeShape 桥接模式形状基类
type BridgeShape struct {
	renderer Renderer
}

func (s *BridgeShape) Draw() {}

// CircleShape 圆形
type CircleShape struct {
	BridgeShape
	radius float64
}

func NewCircleShape(renderer Renderer, radius float64) *CircleShape {
	return &CircleShape{
		BridgeShape: BridgeShape{renderer: renderer},
		radius:      radius,
	}
}

func (c *CircleShape) Draw() {
	c.renderer.RenderCircle(c.radius)
}

// BridgeDemo 演示桥接模式
func BridgeDemo() {
	fmt.Println("\n=== 桥接模式 ===")

	vector := &VectorRenderer{}
	raster := &RasterRenderer{}

	circle1 := NewCircleShape(vector, 5.0)
	circle1.Draw()

	circle2 := NewCircleShape(raster, 10.0)
	circle2.Draw()
}
