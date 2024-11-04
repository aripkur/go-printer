; Script generated by the Inno Setup Script Wizard.
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!

[Setup]
AppName=aaaPDFprint
AppVersion=1.0
AppPublisher=arip.rsuislamboyolali
AppPublisherURL=https://rsuislamboyolali.co.id/
DefaultDirName={pf}\aaaPDFprint
DefaultGroupName=aaaServicePrintV2
OutputDir=output
OutputBaseFilename=aaaPDFprint_RUN_ADMIN
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin

[Files]
Source: "main.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "PDFtoPrinter.exe"; DestDir: "{app}"; Flags: ignoreversion


[Run]
Filename: "{app}\main.exe"; Description: "Launch aaaPDFprint"; Flags: nowait postinstall skipifsilent
