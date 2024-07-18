@echo off
set "file_path=%~1"

:: Change to the directory where the batch script is located
cd /d "%~dp0"

:: Execute the Python script within the src directory with the full path to the dragged file
src\cutsheet-traveller.exe "%file_path%"

pause