BASE_OUTPUT=/home/anorak/Desktop/Projects/PKr-cli/PKr-cli/PKr-base/PKr-base.exe
CLI_OUTPUT=/home/anorak/Desktop/Projects/PKr-cli/PKr-cli/PKr-cli.exe

TEST_DEST=/home/anorak/Desktop/Projects/PKr-test-runs/$(TEST)
TEST=Test-2


#  WIN PATHs
WIN_BASE_OUTPUT=E:\Projects\PKR-Re\Re_PKr-Base\PKr-base.exe
WIN_CLI_OUTPUT=E:\Projects\PKR-Re\Re_PKr-cli\PKr-cli.exe
WIN_TEST_DEST=E:\Projects\PKR-Re\Tests\$(TEST)
WIN_TEST=Test-2


build2test: clean build copy

build: clean
	@echo Building the Go file...
	go build
	cd Re_PKr-base &&  go build

copy: build
	@echo Copying the executable to the destination...
	@copy $(BASE_OUTPUT) $(TEST_DEST)
	@copy $(CLI_OUTPUT) $(TEST_DEST)

clean:
	@echo Cleaning up...
	del $(CLI_OUTPUT)
	del $(BASE_OUTPUT)
	del $(TEST_DEST)\PKr-base.exe
	del $(TEST_DEST)\PKr-cli.exe


#  -------------------- New --------------------
buildall: clean_all build_cli build_base move_exec
# buildall: build_cli build_base move_exec

build_cli:
	@echo Building the Go CLI...
	go build

build_base:
	@echo Building the Go Base...
	cd ../Re_Pkr-Base &&  go build

move_exec:
	copy $(WIN_CLI_OUTPUT) $(WIN_TEST_DEST)\PKr-Cli.exe
	copy $(WIN_BASE_OUTPUT) $(WIN_TEST_DEST)\PKr-Base.exe

clean_all:
	@echo Cleaning up...
	del $(WIN_CLI_OUTPUT)
	del $(WIN_BASE_OUTPUT)
	del $(WIN_TEST_DEST)\PKr-base.exe
	del $(WIN_TEST_DEST)\PKr-cli.exe

.PHONY: build2test build copy clean buildall build_cli build_base move_exec