package designpattern

import "testing"

// ============================================================
// 创建型模式测试
// ============================================================

func TestSingleton(t *testing.T) {
	SingletonDemo()
}

func TestFactory(t *testing.T) {
	FactoryDemo()
}

func TestAbstractFactory(t *testing.T) {
	AbstractFactoryDemo()
}

func TestBuilder(t *testing.T) {
	BuilderDemo()
}

func TestPrototype(t *testing.T) {
	PrototypeDemo()
}

// ============================================================
// 结构型模式测试
// ============================================================

func TestAdapter(t *testing.T) {
	AdapterDemo()
}

func TestDecorator(t *testing.T) {
	DecoratorDemo()
}

func TestProxy(t *testing.T) {
	ProxyDemo()
}

func TestComposite(t *testing.T) {
	CompositeDemo()
}

func TestFacade(t *testing.T) {
	FacadeDemo()
}

func TestBridge(t *testing.T) {
	BridgeDemo()
}

// ============================================================
// 行为型模式测试
// ============================================================

func TestStrategy(t *testing.T) {
	StrategyDemo()
}

func TestObserver(t *testing.T) {
	ObserverDemo()
}

func TestCommand(t *testing.T) {
	CommandDemo()
}

func TestChain(t *testing.T) {
	ChainDemo()
}

func TestTemplate(t *testing.T) {
	TemplateDemo()
}

func TestIterator(t *testing.T) {
	IteratorDemo()
}

func TestState(t *testing.T) {
	StateDemo()
}

func TestMemento(t *testing.T) {
	MementoDemo()
}
