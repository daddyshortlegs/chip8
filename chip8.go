package main

import (
	"chip8"
)

func main() {
	chip8Display := Chip8Display{}
	defer chip8Display.shutdown()
	chip8Display.startUp()

	random := chip8.NewRandom()

	vm := chip8.NewVM(&chip8Display, random)

	//dat, _ := ioutil.ReadFile("test_opcode.ch8")
	//check(err)

	//vm.Load(dat)

	vm.Load(testOpcode())
	vm.Run()
	//Chip8Display.ClearScreen()
}

const x0 = 0x08
const x1 = 0x09
const x2 = 0x0A
const y = 0x0B

func testOpcode() []byte {
	a := chip8.NewAssembler()

	// programs start at 0x200
	const start = 0x200
	address := start

	a.Jump(0x248)
	address += 2

	imageok, address := data(address, a, []byte{0xEA, 0xAC, 0xAA, 0xEA})
	imagefalse, address := data(address, a, []byte{0xCE, 0xAA, 0xAA, 0xAE})
	_, address = data(address, a, []byte{0xE0, 0xA0, 0xA0, 0xE0})
	_, address = data(address, a, []byte{0xC0, 0x40, 0x40, 0xE0})
	_, address = data(address, a, []byte{0xE0, 0x20, 0xC0, 0xE0})
	im3, address := data(address, a, []byte{0xE0, 0x60, 0x20, 0xE0})
	im4, address := data(address, a, []byte{0xA0, 0xE0, 0x20, 0x20})
	im5, address := data(address, a, []byte{0x60, 0x40, 0x20, 0x40})
	_, address = data(address, a, []byte{0xE0, 0x80, 0xE0, 0xE0})
	im7, address := data(address, a, []byte{0xE0, 0x20, 0x20, 0x20})
	_, address = data(address, a, []byte{0xE0, 0xE0, 0xA0, 0xE0})
	im9, address := data(address, a, []byte{0xE0, 0xE0, 0x20, 0xE0})
	imA, address := data(address, a, []byte{0x40, 0xA0, 0xE0, 0xA0})
	_, address = data(address, a, []byte{0xE0, 0xC0, 0x80, 0xE0})
	_, address = data(address, a, []byte{0xE0, 0x80, 0xC0, 0x80})
	imX, address := data(address, a, []byte{0xA0, 0x40, 0xA0, 0xA0})

	// testAX
	// 0x242
	println("testAX = ", address)
	testAX := address
	a.SetIndexRegister(uint16(imageok))
	address += 2
	a.Display(x2, y, 4)
	address += 2
	a.Return()
	address += 2

	// MAIN
	// 0x248
	println("main = ", address)
	a.SetRegister(x0, 1)
	a.SetRegister(x1, 5)
	a.SetRegister(x2, 10)
	a.SetRegister(y, 1)
	a.SetRegister(0x05, 42)
	a.SetRegister(0x06, 43)

	drawop(a, im3, imX)

	a.SetIndexRegister(uint16(imageok))
	a.SkipIfEqual(0x06, 43)
	a.SetIndexRegister(uint16(imagefalse))
	a.Display(x2, y, 4)

	//test 4x
	a.SetRegister(y, 6)
	drawop(a, im4, imX)
	a.SetIndexRegister(uint16(imagefalse))
	a.SkipIfNotEqual(0x05, 42)
	a.SetIndexRegister(uint16(imageok))
	a.Display(x2, y, 4)

	//test 5x
	a.SetRegister(y, 11)
	drawop(a, im5, imX)
	a.SetIndexRegister(uint16(imagefalse))
	a.SkipIfRegistersEqual(0x05, 0x06)
	a.SetIndexRegister(uint16(imageok))
	a.Display(x2, y, 4)

	//test 7x
	a.SetRegister(y, 16)
	drawop(a, im7, imX)
	a.SetIndexRegister(uint16(imagefalse))
	a.AddToRegister(0x06, 255)
	a.SkipIfNotEqual(0x06, 42)
	a.SetIndexRegister(uint16(imageok))
	a.Display(x2, y, 4)

	//test 9x
	a.SetRegister(y, 21)
	drawop(a, im9, imX)
	a.SetIndexRegister(uint16(imagefalse))
	a.SkipIfRegistersNotEqual(0x05, 0x06)
	a.SetIndexRegister(uint16(imageok))
	a.Display(x2, y, 4)

	// test AX
	a.SetRegister(y, 26)
	drawop(a, imA, imX)
	a.Sub(uint16(testAX))

	//test 0E
	//a.SetRegister(x0, 2)
	//a.SetRegister(x1, 2)
	//a.SetRegister(x2, 32)
	//a.SetRegister(y, 10)
	//
	//drawop(a, im0, imE)
	//// drawop im0 imE
	//a.SetIndexRegister(uint16(im0))
	//a.Display(x0, y, 4)
	//a.SetIndexRegister(uint16(imE))
	//a.Display(x1, y, 4)

	a.GetKey(3)
	a.GetKey(3)

	return a.Assemble()
}

func drawop(a *chip8.Assembler, im3 int, imX int) {
	a.SetIndexRegister(uint16(im3))
	a.Display(byte(x0), byte(y), 4)
	a.SetIndexRegister(uint16(imX))
	a.Display(byte(x1), byte(y), 4)
}

func data(address int, a *chip8.Assembler, bytes []byte) (int, int) {
	im3 := address
	a.Data(bytes)
	address += len(bytes)
	return im3, address
}

func andysProgram() []byte {
	a := chip8.NewAssembler()

	a.ClearScreen()
	//drawChar(a, 0, 0, 0x1)
	//drawChar(a, 5, 0, 0x2)
	//drawChar(a, 10, 0, 0x3)
	//drawChar(a, 15, 0, 0x4)
	//drawChar(a, 20, 0, 0x5)
	//drawChar(a, 25, 0, 0xA)
	drawChar(a, 0xB, 30, 0)
	drawChar(a, 0xC, 35, 0)
	drawChar(a, 0xD, 40, 0)

	a.GetKey(3)
	a.ClearScreen()

	a.GetKey(4) // This will not wait

	return a.Assemble()
}

func drawChar(a *chip8.Assembler, char byte, xPosition int, yPosition int) {
	a.SetRegister(0, char)
	a.FontChar(0)

	a.SetRegister(1, byte(xPosition))
	a.SetRegister(2, byte(yPosition))
	a.Display(1, 2, byte(5))
}
