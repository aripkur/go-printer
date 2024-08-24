/* Send a document as raw bytes to a Windows printer
 * The purpose of this utility is to allow a console application (especially
 * a DOS program) to send printer output to a network printer. The DOS program
 * must be reconfigured to send the data to a print file; this utility will
 * then send the print file to the printer.
 *
 * Build with OpenWatcom C/C++:
 *  wcl386 -wx -d2 -q spool.c
 */

#include <windows.h>
#include <stdio.h>
#include <assert.h>


char **PrinterNames = NULL;     // Array to store printer names
DWORD NumPrinters = 0;


BOOL FindPrinters(void)
{
  // Cycle though the list of available printers.
  // Allocate PrinterNames of the appropriate size, store the printer names.
  BYTE *pbuf;
  DWORD dwSize, index;

  // Find out how much bytes is needed to hold PRINTER_INFO_5 array
  EnumPrinters(PRINTER_ENUM_LOCAL, NULL, 5, NULL,0, &dwSize, &NumPrinters);
  pbuf = malloc(dwSize*sizeof(BYTE));
  assert(pbuf != NULL);

  // Read in the printer array
  EnumPrinters(PRINTER_ENUM_LOCAL, NULL, 5, pbuf, dwSize, &dwSize, &NumPrinters);
  if (NumPrinters > 0) {
    PRINTER_INFO_5 *pPrnInfo = (PRINTER_INFO_5 *)pbuf;
    PrinterNames = malloc(NumPrinters * sizeof(char*));
    assert(PrinterNames != NULL);
    for (index=0; index < NumPrinters; index++,pPrnInfo++) {
      PrinterNames[index] = strdup(pPrnInfo->pPrinterName);
      assert(PrinterNames[index] != NULL);
    } /* for */
  } /* if */
  free(pbuf);

  return (NumPrinters > 0);
}


// From Q246772
//
// We are explicitly linking to GetDefaultPrinter, because linking
// implicitly on Windows 95/98 or NT4 results in a runtime error.
// This block specifies which text version you explicitly link to.
#ifdef UNICODE
  #define GETDEFAULTPRINTER "GetDefaultPrinterW"
#else
  #define GETDEFAULTPRINTER "GetDefaultPrinterA"
#endif

// Size of internal buffer used to hold "printername,drivername,portname"
// string. You may have to increase this for huge strings.
#define MAXBUFFERSIZE 250

/*----------------------------------------------------------------*/
/* DPGetDefaultPrinter                                            */
/*                                                                */
/* Parameters:                                                    */
/*   pPrinterName: Buffer alloc'd by caller to hold printer name. */
/*   pdwBufferSize: On input, ptr to size of pPrinterName.        */
/*          On output, min required size of pPrinterName.         */
/*                                                                */
/* NOTE: You must include enough space for the NULL terminator!   */
/*                                                                */
/* Returns: TRUE for success, FALSE for failure.                  */
/*----------------------------------------------------------------*/
BOOL DPGetDefaultPrinter(LPTSTR pPrinterName, LPDWORD pdwBufferSize)
{
  BOOL bFlag;
  OSVERSIONINFO osv;
  TCHAR cBuffer[MAXBUFFERSIZE];
  PRINTER_INFO_2 *ppi2 = NULL;
  DWORD dwNeeded = 0;
  DWORD dwReturned = 0;
  HMODULE hWinSpool = NULL;
  PROC fnGetDefaultPrinter = NULL;

  // What version of Windows are you running?
  osv.dwOSVersionInfoSize = sizeof(OSVERSIONINFO);
  GetVersionEx(&osv);

  // If Windows 95 or 98, use EnumPrinters.
  if (osv.dwPlatformId == VER_PLATFORM_WIN32_WINDOWS)
  {
    // The first EnumPrinters() tells you how big our buffer must
    // be to hold ALL of PRINTER_INFO_2. Note that this will
    // typically return FALSE. This only means that the buffer (the 4th
    // parameter) was not filled in. You do not want it filled in here.
    SetLastError(0);
    bFlag = EnumPrinters(PRINTER_ENUM_DEFAULT, NULL, 2, NULL, 0, &dwNeeded, &dwReturned);
    {
      if ((GetLastError() != ERROR_INSUFFICIENT_BUFFER) || (dwNeeded == 0))
        return FALSE;
    }

    // Allocate enough space for PRINTER_INFO_2.
    ppi2 = (PRINTER_INFO_2 *)GlobalAlloc(GPTR, dwNeeded);
    if (!ppi2)
      return FALSE;

    // The second EnumPrinters() will fill in all the current information.
    bFlag = EnumPrinters(PRINTER_ENUM_DEFAULT, NULL, 2, (LPBYTE)ppi2, dwNeeded, &dwNeeded, &dwReturned);
    if (!bFlag)
    {
      GlobalFree(ppi2);
      return FALSE;
    }

    // If specified buffer is too small, set required size and fail.
    if ((DWORD)lstrlen(ppi2->pPrinterName) >= *pdwBufferSize)
    {
      *pdwBufferSize = (DWORD)lstrlen(ppi2->pPrinterName) + 1;
      GlobalFree(ppi2);
      return FALSE;
    }

    // Copy printer name into passed-in buffer.
    lstrcpy(pPrinterName, ppi2->pPrinterName);

    // Set buffer size parameter to minimum required buffer size.
    *pdwBufferSize = (DWORD)lstrlen(ppi2->pPrinterName) + 1;
  }

  // If Windows NT, use the GetDefaultPrinter API for Windows 2000,
  // or GetProfileString for version 4.0 and earlier.
  else if (osv.dwPlatformId == VER_PLATFORM_WIN32_NT)
  {
    if (osv.dwMajorVersion >= 5) // Windows 2000 or later (use explicit call)
    {
      hWinSpool = LoadLibrary("winspool.drv");
      if (!hWinSpool)
        return FALSE;
      fnGetDefaultPrinter = GetProcAddress(hWinSpool, GETDEFAULTPRINTER);
      if (!fnGetDefaultPrinter)
      {
        FreeLibrary(hWinSpool);
        return FALSE;
      }

      bFlag = fnGetDefaultPrinter(pPrinterName, pdwBufferSize);
      FreeLibrary(hWinSpool);
      if (!bFlag)
        return FALSE;
    }

    else // NT4.0 or earlier
    {
      // Retrieve the default string from Win.ini (the registry).
      // String will be in form "printername,drivername,portname".
      if (GetProfileString("windows", "device", ",,,", cBuffer, MAXBUFFERSIZE) == 0)
        return FALSE;

      // Printer name precedes first "," character.
      strtok(cBuffer, ",");

      // If specified buffer is too small, set required size and fail.
      if ((DWORD)lstrlen(cBuffer) >= *pdwBufferSize)
      {
        *pdwBufferSize = (DWORD)lstrlen(cBuffer) + 1;
        return FALSE;
      }

      // Copy printer name into passed-in buffer.
      lstrcpy(pPrinterName, cBuffer);

      // Set buffer size parameter to minimum required buffer size.
      *pdwBufferSize = (DWORD)lstrlen(cBuffer) + 1;
    }
  }

  // Clean up.
  if (ppi2)
    GlobalFree(ppi2);

  return TRUE;
}
#undef MAXBUFFERSIZE
#undef GETDEFAULTPRINTER


static const char *skippath(const char *name)
{
  const char *ptr=strrchr(name,'\\');
  if (ptr==NULL)
    ptr=strchr(name,':');
  if (ptr==NULL)
    return name;
  return ptr+1;
}

int main(int argc, char *argv[])
{
  DWORD index;
  DWORD Count;  // Count of bytes sent to printer
  DWORD Job;
  DOC_INFO_1 DocAttrib;
  HANDLE hPrinter;
  FILE *fp;
  unsigned char buffer[512];
  int numrd;
  const char *name;

  if (!FindPrinters()) {
    printf("ERROR: no printers available\n");
    return 1;
  } /* if */

  if (argc < 2 || argc > 4) {
    printf("Spool 1.3\n(c) Copyright 2008-2020, CompuPhase, Netherlands\n\n");
    printf("USAGE: %s filename [printer name] [document title]\n\n",skippath(argv[0]));
    /* get name of the default printer */
    Count = sizeof buffer;
    DPGetDefaultPrinter(buffer, &Count);
    printf("Available printers:\n");
    for (index=0; index<NumPrinters; index++) {
      if (stricmp(buffer,PrinterNames[index])==0)
        printf("\t* ");
      else
        printf("\t  ");
      printf("%s\n",PrinterNames[index]);
    } /* for */
    printf("\nThe printer marked with a * is the default printer\n");
    return 1;
  } /* if */

  /* find the selected printer */
  if (argc > 2) {
    name = argv[2];
  } else {
    /* get name of the default printer, if none is supplied */
    Count = sizeof buffer;
    DPGetDefaultPrinter(buffer, &Count);
    name = (const char *)buffer;
  } /* if */
  for (index=0; index<NumPrinters; index++)
    if (stricmp(PrinterNames[index],name)==0)
      break;
  if (index>=NumPrinters) {
    /* try finding the printer after ignoring the server path */
    for (index=0; index<NumPrinters; index++)
      if (stricmp(skippath(PrinterNames[index]),name)==0)
        break;
  } /* if */
  if (index>=NumPrinters) {
    printf("ERROR: printer %s is not found\n",name);
    printf("Run %s without parameters to see a list of available printers\n",skippath(argv[0]));
    return 1;
  } /* if */

  /* open the document */
  fp = fopen(argv[1],"rb");
  if (fp == NULL) {
    printf("ERROR: cannot open print file %s\n",argv[1]);
    return 1;
  } /* if */

  if (!OpenPrinter(PrinterNames[index], &hPrinter, NULL)) {
    fclose(fp);
    printf("ERROR: failed to open printer %s\n",PrinterNames[index]);
    return 1;
  } /* if */
  assert(hPrinter!=NULL);

  DocAttrib.pDocName = (argc > 3) ? argv[3] : (char*)skippath(argv[1]);
  DocAttrib.pOutputFile = NULL;
  DocAttrib.pDatatype= "RAW"; // Must be "RAW" or else StarDocPrinter() will return Invalid Datatype

  if ((Job = StartDocPrinter(hPrinter, 1, (LPBYTE)&DocAttrib)) == 0) {
    ClosePrinter(hPrinter);
    fclose(fp);
    printf("ERROR: failed to start print job\n");
    return 1;
  } /* if */

  if (!StartPagePrinter(hPrinter)) {  // entire document is a single virtual "page"
    EndDocPrinter(hPrinter);
    ClosePrinter(hPrinter);
    fclose(fp);
    printf("ERROR: failed to start new (virtual) page\n");
    return 1;
  } /* if */

  while ((numrd=fread(buffer,1,sizeof buffer,fp)) > 0) {
    if (!WritePrinter(hPrinter,buffer, numrd, &Count) || Count != numrd) {
      EndPagePrinter(hPrinter);
      EndDocPrinter(hPrinter);
      ClosePrinter(hPrinter);
      fclose(fp);
      printf("ERROR: printer does not accept (all) data\n");
      return 1;
    } /* if */
  } /* while */

  EndPagePrinter(hPrinter);   //End of virtual Page
  EndDocPrinter(hPrinter);
  ClosePrinter(hPrinter);
  fclose(fp);

  assert(PrinterNames != NULL);
  for (index=0; index<NumPrinters; index++) {
    assert(PrinterNames[index] != NULL);
    free(PrinterNames[index]);
  } /* if */
  free(PrinterNames);

  return 0;
}
