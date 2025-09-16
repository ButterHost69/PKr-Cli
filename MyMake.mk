ROOT_DIR=E:\Projects\Picker-Pal
CLI_OUTPUT=$(ROOT_DIR)\PKr-Cli\PKr-Cli.exe
TEST_DEST=$(ROOT_DIR)\PKr-Test
TEST_MOIT=$(TEST_DEST)\Moit
TEST_PALAS=$(TEST_DEST)\Palas

build2test:clean build copy

build:
	@cls
	@echo Building the PKr-Cli File ...
	@go build -o PKr-Cli.exe

copy:
	@echo Copying the Executable to Test Destination ...

	@copy "$(CLI_OUTPUT)" "$(TEST_DEST)"
	@copy "$(CLI_OUTPUT)" "$(TEST_MOIT)"
	@copy "$(CLI_OUTPUT)" "$(TEST_PALAS)"
	
	@del "$(CLI_OUTPUT)"

clean:
	@cls
	@echo Cleaning Up ...

	@del "$(TEST_DEST)\PKr-Cli.exe" || exit 0
	@del "$(TEST_MOIT)\PKr-Cli.exe" || exit 0
	@del "$(TEST_PALAS)\PKr-Cli.exe" || exit 0
