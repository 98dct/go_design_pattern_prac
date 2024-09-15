package adapter

import (
	"fmt"
	"testing"
)

/**
常规适配器模式：
被适配的对象需要适配器进行适配，
适配器类中需要包含被适配的类，对外提供一个可以被接受的方法
*/

// 稳定提供5v电压
type PhoneCharge interface {
	Output5v()
}

// 华为充电器
type HuaweiCharger struct {
}

func NewHuaweiCharger() *HuaweiCharger {
	return &HuaweiCharger{}
}

func (h *HuaweiCharger) Output5v() {
	fmt.Println("华为手机充电器 输出5v电压 ...")
}

// 小米充电器
type XiaomiCharger struct {
}

func NewXiaomiCharger() *XiaomiCharger {
	return &XiaomiCharger{}
}

func (x *XiaomiCharger) Output5v() {
	fmt.Println("小米充电器 输出电压5v...")
}

// 苹果充电器
type AppleCharger struct {
}

func NewAppleCharge() *AppleCharger {
	return &AppleCharger{}
}

func (apple *AppleCharger) Output28v() {
	fmt.Println("苹果充电器 输出电压28v。。。")
}

type AppleChargerAdapter struct {
	adaptee *AppleCharger
}

func NewAppleChargerAdapter(adaptee *AppleCharger) *AppleChargerAdapter {
	return &AppleChargerAdapter{adaptee: adaptee}
}

func (appleChargerAdapter *AppleChargerAdapter) Output5v() {
	appleChargerAdapter.adaptee.Output28v()
	fmt.Println("适配器将输出电压调整为5v。。。")
}

type Phone interface {
	charge(phoneCharger PhoneCharge)
}

// 华为手机
type HuaweiPhone struct {
}

func NewHuaweiPhone() Phone {
	return &HuaweiPhone{}
}

func (huaweiPhone *HuaweiPhone) charge(phoneCharge PhoneCharge) {
	fmt.Println("华为手机开始充电")
	phoneCharge.Output5v()
}

func TestAdapter(t *testing.T) {
	huaweiPhone := NewHuaweiPhone()
	// 使用华为手机充电器充电
	huaweiCharger := NewHuaweiCharger()
	huaweiPhone.charge(huaweiCharger)

	//使用适配器转换后的华为充电器充电
	appleCharge := NewAppleCharge()
	appleChargerAdapter := NewAppleChargerAdapter(appleCharge)
	huaweiPhone.charge(appleChargerAdapter)
}
