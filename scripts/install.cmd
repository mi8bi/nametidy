@echo off
setlocal enabledelayedexpansion

:: Script to download and install the latest NameTidy binary for Windows

echo NameTidy Installer for Windows
echo ------------------------------

:: --- Configuration ---
set "BINARY_NAME=nametidy.exe"
set "INSTALL_DIR=%USERPROFILE%\bin"
set "RELEASE_BASE_URL=https://github.com/mi8bi/NameTidy/releases/latest/download"
set "OS_NAME=windows"
set "ARCH="
set "MANUAL_MOVE_NEEDED=false"
set "TEMP_DIR="

:: --- Determine Architecture ---
echo Detecting system architecture...
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set "ARCH=amd64"
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set "ARCH=arm64"
) else (
    echo Error: Unsupported architecture: %PROCESSOR_ARCHITECTURE%. Only AMD64 and ARM64 are supported.
    goto :error_exit
)
echo Detected Architecture: %ARCH%

:: --- Fetch Latest Release Information ---
echo.
echo Fetching latest release information...
set "LATEST_RELEASE_INFO_URL=https://api.github.com/repos/mi8bi/NameTidy/releases/latest"
set "TEMP_JSON_FILE=%TEMP%\release_info_%RANDOM%.json"
set "TAG_NAME="
set "VERSION="

:: Check for curl
where curl >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: curl is required to fetch release information but not found.
    echo Please install curl (e.g., from https://curl.se/windows/) and ensure it's in your PATH.
    goto :error_exit
)

:: Fetch release info using curl
echo Fetching from %LATEST_RELEASE_INFO_URL% ...
curl -LsS -o "%TEMP_JSON_FILE%" "%LATEST_RELEASE_INFO_URL%"
if errorlevel 1 (
    echo Error: Failed to download release information from %LATEST_RELEASE_INFO_URL%.
    if exist "%TEMP_JSON_FILE%" del "%TEMP_JSON_FILE%"
    goto :error_exit
)

if not exist "%TEMP_JSON_FILE%" (
    echo Error: Release information file was not created at "%TEMP_JSON_FILE%".
    goto :error_exit
)

:: Parse tag_name from JSON. This is fragile; assumes "tag_name": "vX.Y.Z" format.
FOR /F "tokens=2 delims=:," %%g IN ('findstr /C:"\"tag_name\":" "%TEMP_JSON_FILE%"') DO (
    FOR /F "tokens=1 delims= " %%h IN ("%%g") DO (
        set "TAG_NAME=%%~h"
    )
)

if defined TEMP_JSON_FILE if exist "%TEMP_JSON_FILE%" (
    del "%TEMP_JSON_FILE%"
)

if not defined TAG_NAME (
    echo Error: Could not parse tag_name from release information.
    echo The format of the release JSON may have changed or the file was empty.
    goto :error_exit
)
:: Remove potential surrounding quotes from TAG_NAME
set "TAG_NAME=%TAG_NAME:"=%"
echo Latest release tag: %TAG_NAME%

:: Remove 'v' prefix from tag to get version. Example: v0.1.0 -> 0.1.0
if "%TAG_NAME:~0,1%"=="v" (
    set "VERSION=%TAG_NAME:~1%"
) else (
    set "VERSION=%TAG_NAME%"
)

if not defined VERSION (
    echo Error: Could not extract version from tag '%TAG_NAME%'.
    goto :error_exit
)
echo Detected version: %VERSION%

:: Construct the download URL
set "RELEASE_DOWNLOAD_URL_BASE=https://github.com/mi8bi/NameTidy/releases/download"
set "ASSET_NAME=NameTidy_%VERSION%_%OS_NAME%_%ARCH%.zip"
set "DOWNLOAD_URL=%RELEASE_DOWNLOAD_URL_BASE%/%TAG_NAME%/%ASSET_NAME%"

echo Download URL: %DOWNLOAD_URL%

:: --- Temporary Download Path ---
:: Create a temporary directory for downloads
set "TEMP_DIR=%TEMP%\NameTidy_Install_%RANDOM%"
mkdir "%TEMP_DIR%"
if not exist "%TEMP_DIR%\" (
    echo Error: Failed to create temporary directory: %TEMP_DIR%
    goto :error_exit
)
echo Temporary directory created: %TEMP_DIR%
set "ARCHIVE_PATH=%TEMP_DIR%\%ASSET_NAME%"
set "EXTRACTED_BINARY_PATH=%TEMP_DIR%\%BINARY_NAME%"

:: --- Download Logic ---
echo.
echo Downloading NameTidy release asset...

:: Check for curl
where curl >nul 2>nul
if %errorlevel% equ 0 (
    echo Found curl. Attempting download...
    curl -LSsf -o "%ARCHIVE_PATH%" "%DOWNLOAD_URL%"
    if errorlevel 1 (
        echo Error: curl download failed. Check URL or network connection.
        echo Asset might not be available for %OS_NAME%/%ARCH%.
        goto :error_exit
    )
    echo Download successful using curl.
    goto :extract_logic
)

:: Check for bitsadmin
where bitsadmin >nul 2>nul
if %errorlevel% equ 0 (
    echo Found bitsadmin. Attempting download...
    bitsadmin /transfer NameTidyDownloadJob /download /priority NORMAL "%DOWNLOAD_URL%" "%ARCHIVE_PATH%"
    if errorlevel 1 (
        echo Error: bitsadmin download failed. Check URL or network connection.
        echo Asset might not be available for %OS_NAME%/%ARCH%.
        goto :error_exit
    )
    echo Download successful using bitsadmin.
    goto :extract_logic
)

echo Error: Neither curl nor bitsadmin found.
echo Please install curl (recommended: https://curl.se/windows/) or download the file manually:
echo %DOWNLOAD_URL%
echo And place %ASSET_NAME% in "%TEMP_DIR%"
pause
if not exist "%ARCHIVE_PATH%" (
    echo Manual download not completed or file not found at expected location.
    goto :error_exit
)
echo Assuming manual download completed.

:extract_logic
:: --- Extraction Logic ---
echo.
echo Extracting archive from "%ARCHIVE_PATH%" into "%TEMP_DIR%"...
set "EXTRACTION_ATTEMPTED=false"
set "EXTRACTION_SUCCESSFUL=false"

:: Check for tar (bsdtar, included in modern Windows)
if not %EXTRACTION_SUCCESSFUL%==true (
    where tar >nul 2>nul
    if %errorlevel% equ 0 (
        set "EXTRACTION_ATTEMPTED=true"
        echo Found tar. Attempting extraction...
        tar -xf "%ARCHIVE_PATH%" -C "%TEMP_DIR%" >nul 2>nul
        if errorlevel 1 (
            echo Warning: tar extraction may have failed or reported errors (errorlevel %errorlevel%).
            echo Will proceed to search for the binary anyway.
            :: Don't set EXTRACTION_SUCCESSFUL to true here if tar reports an error,
            :: but allow search to proceed as tar might still have extracted some files.
        ) else (
            echo Extraction with tar seems complete.
            set "EXTRACTION_SUCCESSFUL=true"
        )
        goto :search_for_binary_after_extraction
    )
)

:: Check for PowerShell Expand-Archive
if not %EXTRACTION_SUCCESSFUL%==true (
    where powershell >nul 2>nul
    if %errorlevel% equ 0 (
        set "EXTRACTION_ATTEMPTED=true"
        echo Found PowerShell. Attempting extraction...
        powershell -NoProfile -ExecutionPolicy Bypass -Command "try { Expand-Archive -Path '%ARCHIVE_PATH%' -DestinationPath '%TEMP_DIR%' -Force } catch { Write-Error $_; exit 1 }"
        if errorlevel 1 (
            echo Error: PowerShell Expand-Archive failed. Archive might be corrupt.
            :: Unlike tar, if PowerShell fails, it's usually a more definite failure.
            goto :error_exit
        )
        echo Extraction successful using PowerShell.
        set "EXTRACTION_SUCCESSFUL=true"
        goto :search_for_binary_after_extraction
    )
)

:search_for_binary_after_extraction
if not "%EXTRACTION_ATTEMPTED%"=="true" (
    echo Error: Neither tar nor PowerShell Expand-Archive found. Cannot automate extraction.
    echo Please manually extract %BINARY_NAME% from "%ARCHIVE_PATH%"
    echo into folder: "%TEMP_DIR%"
    pause
    echo Resuming to search for binary after manual extraction attempt...
    :: After pause, script will proceed to search. User must have placed files in TEMP_DIR.
)

:: --- Search for the binary ---
echo.
echo Searching for the binary in "%TEMP_DIR%"...
set "FOUND_BINARY_PATH="
:: %BINARY_NAME% is "nametidy.exe"
set "POSSIBLE_BINARY_NAMES=%BINARY_NAME% NameTidy.exe"

FOR %%N IN (%POSSIBLE_BINARY_NAMES%) DO (
    IF not defined FOUND_BINARY_PATH (
        echo Trying to find %%N...
        FOR /F "delims=" %%F IN ('dir /s /b "%TEMP_DIR%\%%N" 2^>nul') DO (
            IF not defined FOUND_BINARY_PATH (
                echo Found binary at: "%%F"
                set "FOUND_BINARY_PATH=%%F"
            )
        )
    )
)

IF not defined FOUND_BINARY_PATH (
    echo Error: Could not find %BINARY_NAME% or NameTidy.exe in the extracted files at "%TEMP_DIR%".
    echo Listing contents of "%TEMP_DIR%" (recursive):
    dir "%TEMP_DIR%" /s /b /A
    goto :error_exit
)

:: Update EXTRACTED_BINARY_PATH to the actual found path
set "EXTRACTED_BINARY_PATH=%FOUND_BINARY_PATH%"
echo Binary to be installed: %EXTRACTED_BINARY_PATH%
:: Note: The original EXTRACTED_BINARY_PATH was %TEMP_DIR%\%BINARY_NAME%. This update is crucial.
:: The FINAL_INSTALL_PATH is %INSTALL_DIR%\%BINARY_NAME%.
:: If "NameTidy.exe" is found, it will be moved and renamed to "nametidy.exe" in the install directory. This is desired.

goto :install_logic :: Proceed to installation


:install_logic
:: --- Installation Logic ---
echo.
echo Installing %BINARY_NAME%...

:: Create installation directory if it doesn't exist
if not exist "%INSTALL_DIR%\" (
    echo Creating installation directory: "%INSTALL_DIR%"
    mkdir "%INSTALL_DIR%"
    if errorlevel 1 (
        echo Error: Failed to create installation directory: "%INSTALL_DIR%". Check permissions.
        goto :error_exit
    )
)

set "FINAL_INSTALL_PATH=%INSTALL_DIR%\%BINARY_NAME%"

:: Check if binary already exists and prompt for overwrite
if exist "%FINAL_INSTALL_PATH%" (
    echo Warning: %BINARY_NAME% already exists at "%FINAL_INSTALL_PATH%".
    choice /C YN /M "Do you want to overwrite it?"
    if errorlevel 2 (
        echo Overwrite cancelled by user. Exiting.
        goto :user_exit_graceful
    )
    echo Proceeding with overwrite...
)

:: Move the binary
echo Moving %BINARY_NAME% to "%FINAL_INSTALL_PATH%"...
move /Y "%EXTRACTED_BINARY_PATH%" "%FINAL_INSTALL_PATH%" >nul
if errorlevel 1 (
    echo Error: Failed to move %BINARY_NAME% to "%FINAL_INSTALL_PATH%". Check permissions.
    echo The extracted binary is still available at "%EXTRACTED_BINARY_PATH%"
    set "MANUAL_MOVE_NEEDED=true"
    goto :path_logic_or_skip
)
echo %BINARY_NAME% successfully installed to "%FINAL_INSTALL_PATH%".

:: --- PATH Update ---
:path_logic_or_skip
if "%MANUAL_MOVE_NEEDED%"=="true" (
    echo Skipping PATH update as installation was not fully automatic.
    goto :final_message
)

echo.
echo Attempting to add "%INSTALL_DIR%" to your User PATH...
:: Important: setx modifies the persistent (registry) PATH, not the current session's PATH.
:: %PATH% here refers to the current session's PATH.
:: This command appends INSTALL_DIR to the user's PATH. If it's already there, it will be duplicated.
:: A more sophisticated script might check for existence first via reg query.
setx PATH "%PATH%;%INSTALL_DIR%"
if errorlevel 1 (
    echo Warning: Failed to automatically add "%INSTALL_DIR%" to PATH using setx.
    echo This can happen if the PATH is too long or due to permissions.
    echo You may need to add it manually through System Properties (Environment Variables).
) else (
    echo "%INSTALL_DIR%" has been scheduled to be added to your User PATH.
    echo This change will take effect in new command prompts.
    echo You might need to restart your current prompt, or log out and log back in.
)

:final_message
echo.
if "%MANUAL_MOVE_NEEDED%"=="true" (
    echo Installation requires manual intervention.
    echo Please manually move %BINARY_NAME% from "%EXTRACTED_BINARY_PATH%"
    echo to "%INSTALL_DIR%"
    echo Then, ensure "%INSTALL_DIR%" is in your PATH.
    echo The temporary files (including the extracted binary at "%EXTRACTED_BINARY_PATH%") will be kept.
    echo Please clean up "%TEMP_DIR%" manually after moving the binary.
    set "TEMP_DIR=" :: Prevent cleanup routine from deleting it
    goto :eof
)

echo Installation complete!
echo Please open a new command prompt to use %BINARY_NAME%.
call :cleanup_and_exit

:user_exit_graceful
echo Exiting script.
call :cleanup_and_exit

:error_exit
echo An error occurred. Installation aborted.
call :cleanup_and_exit

:cleanup_and_exit
if defined TEMP_DIR if exist "%TEMP_DIR%\" (
    echo Cleaning up temporary directory: "%TEMP_DIR%"
    rmdir /S /Q "%TEMP_DIR%"
    set "TEMP_DIR="
)
goto :eof

endlocal
