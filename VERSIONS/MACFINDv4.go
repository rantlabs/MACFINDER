package main

import (
    "bufio"
    _ "embed"
    "fmt"
    "log"
    "os"
    "regexp"
    "strings"
)

//go:embed oui_v2.txt
var vendorData string

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run main.go <input_file>")
    }

    inputFile := os.Args[1]
    file, err := os.Open(inputFile)
    if err != nil {
        log.Fatalf("error opening input file: %v", err)
    }
    defer file.Close()

    var inputLines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        inputLines = append(inputLines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("error reading input file: %v", err)
    }

    // Step 2: Extract MAC addresses from input file
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

    // Step 4: Print the original file with matched vendor information
    fmt.Println("Output with Vendor Information:")
    for _, line := range inputLines {
        // Append vendor info if available
        if updatedLine, found := findUpdatedLine(line, macAddresses, macRegex); found {
            fmt.Println(updatedLine)
        } else {
            fmt.Println(line)
        }
    }
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
        if strings.Contains(line, "(hex)") || strings.Contains(line, "(base 16)") {
            fields := strings.Fields(line)
            if len(fields) > 1 {
                macPrefix := strings.ReplaceAll(strings.ToLower(fields[0]), "-", "")
                macPrefix = strings.ReplaceAll(macPrefix, ".", "")
                vendor := strings.Join(fields[2:], " ")
                vendorInfo[macPrefix] = vendor
            }
        }
    }
    return vendorInfo
}

