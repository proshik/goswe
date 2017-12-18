@echo off

set mypath=%cd%

if exist goswe.exe (
    del goswe.exe
) else (
    ECHO File not found for delete
)

FOR /F "tokens=* USEBACKQ" %%F IN (`go build`) DO (
    SET var=%%F
)

if exist goswe.exe (
    start %mypath%\goswe.exe
) else (
    ECHO %var%
)


