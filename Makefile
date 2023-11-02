APP_NAME=chip8-app

build:
	go build -o ${APP_NAME}

run: build
	./${APP_NAME} -rom "test_opcode.ch8"
