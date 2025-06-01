package main

import (
    "bufio"
    _ "embed"
    "flag"
    "fmt"
    "log"
    "os"
    "regexp"
    "strings"
)

//go:embed oui_v2.txt
var vendorData string

func main() {
    var inputFile string
    var macAddress string
    var outputFile string

    flag.StringVar(&inputFile, "file", "", "Input file containing MAC addresses")
    flag.StringVar(&macAddress, "mac", "", "Single MAC address to look up")
    flag.StringVar(&outputFile, "output", "", "Output file to save results")
    help := flag.Bool("help", false, "Display help")

    flag.Parse()

    // Check if data is being piped through stdin
    stat, err := os.Stdin.Stat()
    if err != nil {
        log.Fatalf("error getting stdin stat: %v", err)
    }
    isPipedInput := (stat.Mode() & os.ModeCharDevice) == 0

    if *help || (inputFile == "" && macAddress == "" && !isPipedInput) {
        flag.Usage()
        return
    }

    var inputLines []string
    if inputFile != "" {
        file, err := os.Open(inputFile)
        if err != nil {
            log.Fatalf("Error opening input file: %v", err)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            inputLines = append(inputLines, scanner.Text())
        }
        if err := scanner.Err(); err != nil {
            log.Fatalf("Error reading input file: %v", err)
        }
    } else if macAddress != "" {
        inputLines = append(inputLines, macAddress)
    } else {
        // Reading from stdin
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            inputLines = append(inputLines, scanner.Text())
        }
        if err := scanner.Err(); err != nil {
            log.Fatalf("Error reading from stdin: %v", err)
        }
    }

    // Step 2: Extract MAC addresses from input
    macRegex := regexp.MustCompile(`([0-9a-fA-F]{4}\.[0-9a-fA-F]{4}\.[0-9a-fA-F]{4})|([0-9a-fA-F]{2}[:-][0-9a-fA-F]{2}[:-][0-9a-fA-F]{2}[:-][0-9a-fA-F]{2}[:-][0-9a-fA-F]{2}[:-][0-9a-fA-F]{2})`)
    macAddresses := make(map[string]string)
    for _, line := range inputLines {
        matches := macRegex.FindAllString(line, -1)
        for _, mac := range matches {
            normalizedMac := normalizeMAC(mac)
            macAddresses[normalizedMac] = line
        }
    }

    // Step 3: Match MAC addresses with vendor information
    vendorInfo := parseVendorData(vendorData)
    for mac, line := range macAddresses {
        prefix := mac[:6]
        if vendor, found := vendorInfo[prefix]; found {
            updatedLine := fmt.Sprintf("%s  Vendor: %s", strings.TrimSpace(line), vendor)
            macAddresses[mac] = updatedLine
        }
    }

    // Step 4: Prepare to write output
    var writer *bufio.Writer
    if outputFile != "" {
        output, err := os.Create(outputFile)
        if err != nil {
            log.Fatalf("Error creating output file: %v", err)
        }
        defer output.Close()
        writer = bufio.NewWriter(output)
    } else {
        writer = bufio.NewWriter(os.Stdout)
    }

    fmt.Fprintln(writer, "Output with Vendor Information:")
    for _, line := range inputLines {
        if updatedLine, found := findUpdatedLine(line, macAddresses, macRegex); found {
            fmt.Fprintln(writer, updatedLine)
        } else {
            fmt.Fprintln(writer, line)
        }
    }
    writer.Flush()
}

func findUpdatedLine(line string, macAddresses map[string]string, macRegex *regexp.Regexp) (string, bool) {
    matches := macRegex.FindAllString(line, -1)
    for _, mac := range matches {
        normalizedMac := normalizeMAC(mac)
        if updatedLine, found := macAddresses[normalizedMac]; found {
            return updatedLine, true
        }
    }
    return "", false
}

func normalizeMAC(mac string) string {
    // Normalize MAC address to a common format (no separators, lowercase)
    mac = strings.ToLower(mac)
    if strings.Contains(mac, ".") {
        return strings.ReplaceAll(mac, ".", "")
    }
    return strings.ReplaceAll(strings.ReplaceAll(mac, ":", ""), "-", "")
}

func parseVendorData(data string) map[string]string {
    vendorInfo := make(map[string]string)
    scanner := bufio.NewScanner(strings.NewReader(data))
    for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Fields(line)
        if len(fields) >= 2 {
            macPrefix := strings.ReplaceAll(strings.ToLower(fields[0]), "-", "")
            macPrefix = strings.ReplaceAll(macPrefix, ".", "")
            vendor := strings.Join(fields[1:], " ")
            vendorInfo[macPrefix] = vendor
        }
    }
    return vendorInfo
}
