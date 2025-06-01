# MACFINDER
#### A simple very fast MAC address vendor lookup command line application. The vendor database is embedded in executable. MACFINDER is written in GO and can vendor lookup over a hundred thousand mac address per second. MACFINDER takes any input and translates any mac address found on a line and places the vendor name of the mac address at the end of the line.

#### MACFINDER can take the output of any text file such as a GatherDB containing the output of commands such as "show mac address-table" and "show ip arp" and add the vendor for every device. One simple example of MACFINDER's utility is finding all Netgear devices on any size campus. I would do this with the following command:

```
cat GatherDB | grep "show mac address-table" | macfinder | grep -i netgear
```
#### For command line options type macfinder or macfinder -help
```
macfinder 
Usage of macfinder:
  -file string
    	Input file containing MAC addresses
  -help
    	Display help
  -mac string
    	Single MAC address to look up
  -output string
    	Output file to save results
```

#### Input methods:

```
Redirect a mac address of any format (DOT, DASH, or COLON) into macfinder

Example 1: Directing mac address data direclty into macfinder
echo "841b.5e91.7e2d" | macfinder
Output with Vendor Information:
841b.5e91.7e2d  Vendor: NETGEAR

Example 2: Inputing data via unix pipes and redirects.
File (file.txt) with the following mac address:
1. 9cc9.ebd2.3b9d
2. 9418.65e9.f240
3. 5407.7d1c.a2f9
4. 6cb0.ce22.b7c3

macfinder < file.txt
Output with Vendor Information:
1. 9cc9.ebd2.3b9d  Vendor: NETGEAR
2. 9418.65e9.f240  Vendor: NETGEAR
3. 5407.7d1c.a2f9  Vendor: NETGEAR
4. 6cb0.ce22.b7c3  Vendor: NETGEAR

cat file.txt | macfinder
Output with Vendor Information:
1. 9cc9.ebd2.3b9d  Vendor: NETGEAR
2. 9418.65e9.f240  Vendor: NETGEAR
3. 5407.7d1c.a2f9  Vendor: NETGEAR
4. 6cb0.ce22.b7c3  Vendor: NETGEAR

Example 3: Using the command line application flags
macfinder -file file.txt
Output with Vendor Information:
1. 9cc9.ebd2.3b9d  Vendor: NETGEAR
2. 9418.65e9.f240  Vendor: NETGEAR
3. 5407.7d1c.a2f9  Vendor: NETGEAR
4. 6cb0.ce22.b7c3  Vendor: NETGEAR

macfinder -mac "841b.5e91.7e2d"            
Output with Vendor Information:
841b.5e91.7e2d  Vendor: NETGEAR

macfinder -mac 841b.5e91.7e2d 
Output with Vendor Information:
841b.5e91.7e2d  Vendor: NETGEAR
```
#### Executables binaries available for download
```
macfinder-apl - apple mac
macfinder-lnx32 - Linux 32 bit
macfinder-lnx64 - Linux 64 bit
macfinder-rsbpi - Raspberry Pie
macfinder-win32.exe - Windows 32 bit
macfinder-win64.exe - Windows 64 bit
```
#### Further instructions to compile the Go code are in notes. Make sure the oui_v2.txt file/TextDB is in the same directory when compiling. The Go code is in the file named MACFIND_FINALv9.go.






