# Xnopyt's Bulk nHentai Downloader
A Tool for bulk downloads from nHentai, written in Go
<br />
<br />
## Usage
Run the file from a terminal (or command prompt in windows) and follow the on screen instructions.
<br />
<br />
### Advanced Usage
The downloader supports 2 flags: <br />
`--verbose`, this will run the downloader in verbose mode <br />
`--threads 15`, this will define the max number of goroutines that will be using to download
<br />
<br />
## Downloading
Binary downloads for Windows(x86 and x64), linux(i386 and amd64) and MacOS can be found in the releases tab.<br />
<br /> 
<br />
## Building from source
Assuming that you already have the golang and git packages for your system, you can build using the following:<br />
```bash
go get github.com/Xnopyt/nhentai-go
git clone https://github.com/Xnopyt/nhentai-dl.git
cd nhentai-dl
go build -o hentai
```
<br />
In the last line "hentai" will be the file name and should be changed to include the extention for your system.