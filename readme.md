# REST API UNTUK PRINT PDF LANGSUNG KE PRINTER

#### INSTALL SERVICE
klik kanan "aaaPDFprint_RUN_ADMIN.exe" klik run as administrator, klik yes, klik install kemudian centang launch aaaPDFprint

#### ENDPOINT
PORT: 8888

mendapatkan list printer yang ada di komputer

```curl --location 'http://172.16.41.100:8888/'```

mencetak pdf ke printer, jika parameter printer tidak di isi maka akan print ke default printer
```
curl --location 'http://localhost:8888/print' \
--form 'file=@"/C:/Users/edprs/Downloads/SURAT KETERANGAN LAHIR - 006965.pdf"' \
--form 'printer="Microsoft Print to PDF"'
```

#### UNINSTALL SERVICE 
masuk folder instal "C:\Program Files (x86)\aaaPDFprint" jalankan unins000.exe
kemudian buka cmd administrator jalankan sc delete "aaaPDFprint"

#### DEV
jalankan terminal "go run main.go i"