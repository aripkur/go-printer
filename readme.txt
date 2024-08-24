Spool: send print files directly to a Windows printer
=====================================================

Spool is a small Win32 console utility that sends the contents of a file
to the spooler of a Windows printer, bypassing the printer driver. The purpose
is to send a raw print file, such as one produced by the "print to file"
functionality of Windows and many DOS programs, to a printer.

In our documentation generation system, based on TeX, the source documents are
processed to a "DVI" file (device-independent) which then passes through a DVI
printer driver to generate the codes for the specific printer. When the printer
is directly connected to the PC, the DVI printer driver can output to that
printer directly, but for printing to a network printer, it must send its
output to a file, for later spooling to the printer. This "Spool" utility
performs this last step.


License
=======

Spool 1.3 is copyrighted software that is free for personal and commercial use.
You may use it and distribute it without limitations. You may however not remove
or conceal the copyright. There are no guarantees or warranties whatsoever; use
it at your own risk.


Usage
=====

Starting the command "spool" without any command-line arguments shows the usage,
and a list with the available printers. The default printer is also identified.
Below is typical output, on a laptop that has one local (virtual) printer
(PDFCreator), a second printer that is connected to another PC in the network
(called "<$smallcaps "server">") and a third that is directly connected to the
network. This third printer is set as the default.

        Spool 1.3
        (c) Copyright 2008-2020, CompuPhase, Netherlands

        USAGE: spool.exe filename [printer name] [document title]

        Available printers:
                  PDFCreator
                  \\SERVER\Brother
                * \\http://192.168.0.88:631\BRN0ACD3B

        The printer marked with a * is the default printer

Running "spool" with a filename, sends that file to the "default" printer. To
print to any other printer, add the printer name to the command line as the
second argument to spool. The printer name may optionally include the server
path. Given the list of printers in the above example, the names
"\\server\brother" and "brother" would both select the same printer.

If a filename or the printer name contains space characters, the entire name
must be between double quotes.

The optional document title is the name that the print queue will show. If it
contains space characters, it must be between double quotes.
