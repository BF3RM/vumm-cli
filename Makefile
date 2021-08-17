WIN_X86_PATH := dist/vumm_windows_386
WIN_X64_PATH := dist/vumm_windows_amd64

VERSION ?= 0.0.1

.PHONY: msi
msi:
ifeq ($(OS),Windows_NT)
	@candle -nologo -arch x86 -dVersion=$(VERSION) -dPath=$(WIN_X86_PATH)/vumm.exe packaging/vumm.wxs -out $(WIN_X86_PATH)/vumm.wixobj
	@candle -nologo -arch x64 -dVersion=$(VERSION) -dPath=$(WIN_X64_PATH)/vumm.exe packaging/vumm.wxs -out $(WIN_X64_PATH)/vumm.wixobj

	@light -nologo $(WIN_X86_PATH)/vumm.wixobj -o $(WIN_X86_PATH)/vumm.msi -spdb
	@light -nologo $(WIN_X64_PATH)/vumm.wixobj -o $(WIN_X64_PATH)/vumm.msi -spdb
else
	@wixl -a x86 -D Version=$(VERSION) -D Path=$(WIN_X86_PATH)/vumm.exe -o $(WIN_X86_PATH)/vumm.msi packaging/vumm.wxs
	@wixl -a x64 -D Version=$(VERSION) -D Path=$(WIN_X64_PATH)/vumm.exe -o $(WIN_X64_PATH)/vumm.msi packaging/vumm.wxs
endif