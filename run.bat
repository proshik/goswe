@echo off

set mypath=%cd%

if exist gotrew.exe (
    del gotrew.exe
) else (
    ECHO File not found for delete
)

FOR /F "tokens=* USEBACKQ" %%F IN (`go build`) DO (
    SET var=%%F
)

if exist gotrew.exe (
    start %mypath%\gotrew.exe %*
) else (
    ECHO %var%
)