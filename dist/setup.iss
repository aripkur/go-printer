; Script generated by the Inno Setup Script Wizard.
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!

[Setup]
AppName=aaaServicePrintV2
AppVersion=1.0
DefaultDirName={pf}\aaaServicePrintV2
DefaultGroupName=aaaServicePrintV2
OutputDir=output
OutputBaseFilename=aaaServicePrintV2_RUN_ADMIN
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin

[Files]
Source: "main.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "PDFtoPrinter.exe"; DestDir: "{app}"; Flags: ignoreversion


[Run]
Filename: "{app}\main.exe"; Description: "Launch aaaServicePrintV2"; Flags: nowait postinstall skipifsilent
