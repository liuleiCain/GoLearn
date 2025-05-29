package learn_receiver

import (
	"fmt"
	"time"
)

type Person struct {
	Name string
}

func (p *Person) PointerSetName(name string) {
	p.Name = name
}

func (p Person) SetName(name string) {
	p.Name = name
}

func (p *Person) printName() {
	fmt.Println("Name=", p.Name)
}

func ReceiverTypeCompare() {
	p := Person{
		Name: "张三",
	}
	fmt.Println("Name=", p.Name)
	p.SetName("李四")
	fmt.Println("Name=", p.Name)
	p.PointerSetName("王五")
	fmt.Println("Name=", p.Name)
}

func ReceiverListTypeCompare() {
	data1 := []*Person{
		{Name: "张三"},
		{Name: "李四"},
		{Name: "王五"},
	}
	for _, v := range data1 {
		go v.printName()
	}
	data2 := []Person{
		{Name: "张三"},
		{Name: "李四"},
		{Name: "王五"},
	}
	for _, v := range data2 {
		go v.printName()
	}

	time.Sleep(time.Second)
}
