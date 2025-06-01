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
    // Step 1: Read input file from the user
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Enter file content (or type 'EOF' to finish input):")
    var inputLines []string
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            log.Fatalf("error reading input: %v", err)
        }
        if strings.TrimSpace(line) == "EOF" {
            break
        }
        inputLines = append(inputLines, line)
    }

    // Step 2: Extract MAC addresses from input file
    macRegex := regexp.MustCompile(`([0-9a-fA-F]{4}\.[0-9a-fA-F]{4}\.[0-9a-fA-F]{4})|([0-9a-fA-F]{2}[:-][0-9a-fA-F]{2}[:-][0-9a-fA-F]{2}[:-][0-9a-fA-F]{2}[:-][0-9a-fA-F]{2}[:-][0-9a-fA-F]{2})`)
    macAddresses := make(map[string]bool)
    for _, line := range inputLines {
        matches := macRegex.FindAllString(line, -1)
        for _, mac := range matches {
            normalizedMac := normalizeMAC(mac)
            macAddresses[normalizedMac[:6]] = true
        }
    }

    // Step 3: Match MAC addresses with vendor information
    vendorInfo := parseVendorData(vendorData)
    for i, line := range inputLines {
        for mac := range macAddresses {
            if vendor, found := vendorInfo[mac]; found {
                inputLines[i] = fmt.Sprintf("%s  Vendor: %s", strings.TrimSpace(line), vendor)
                break
            }
        }
    }

    // Step 4: Print the original file with matched vendor information
    fmt.Println("Output with Vendor Information:")
    for _, line := range inputLines {
        fmt.Println(line)
    }
}

func normalizeMAC(mac string) string {
    // Normalize MAC address to a common format (no separators for base 16)
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
        if strings.Contains(line, "(hex)") {
            fields := strings.Fields(line)
            if len(fields) > 1 {
                macPrefix := strings.ReplaceAll(strings.ToLower(fields[0]), "-", "")
                vendor := strings.Join(fields[2:], " ")
                vendorInfo[macPrefix] = vendor
            }
        }
        if strings.Contains(line, "(base 16)") {
            fields := strings.Fields(line)
            if len(fields) > 1 {
                macPrefix := strings.ReplaceAll(strings.ToLower(fields[0]), ".", "")
                vendor := strings.Join(fields[2:], " ")
                vendorInfo[macPrefix] = vendor
            }
        }
    }
    return vendorInfo
}

