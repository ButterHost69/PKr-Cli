CLI_OUTPUT="C:\Users\Lappy\OneDrive\Desktop\Cringe xD\Dim Future\Go\Picker-Pal\PKr-Cli\PKr-Cli.exe"

TEST_DEST="C:\Users\Lappy\OneDrive\Desktop\Cringe xD\Dim Future\Go\Picker-Pal\PKr-Test\"
TEST_MOIT="C:\Users\Lappy\OneDrive\Desktop\Cringe xD\Dim Future\Go\Picker-Pal\PKr-Test\Moit"
TEST_PALAS="C:\Users\Lappy\OneDrive\Desktop\Cringe xD\Dim Future\Go\Picker-Pal\PKr-Test\Palas"

build2test:clean build copy done

done:
	@echo $(TEST) is built

build:
	@cls
	@echo Building the PKr-Cli file ...
	@go build -o PKr-Cli.exe

copy:
	@echo Copying the executable to the destination ...

	@copy $(CLI_OUTPUT) $(TEST_DEST)
	@copy $(CLI_OUTPUT) $(TEST_MOIT)
	@copy $(CLI_OUTPUT) $(TEST_PALAS)

	@del $(CLI_OUTPUT)

clean:
	@cls
	@echo Cleaning up ...

	@del $(TEST_DEST)\PKr-Cli.exe || exit 0
	@del $(TEST_MOIT)\PKr-Cli.exe || exit 0
	@del $(TEST_PALAS)\PKr-Cli.exe || exit 0

grpc-out:
	protoc ./proto/*.proto --go_out=. --go-grpc_out=.

get-new-kcp:
	go get github.com/ButterHost69/kcp-go@latest

upgrade-base:
	go get "github.com/ButterHost69/PKr-Base"@latest
