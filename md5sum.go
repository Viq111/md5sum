package main

import (
    "bufio"
    "crypto/md5"
    "encoding/hex"
    "errors"
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

const chunkSize = 1024*1024

func usage() {
    fmt.Fprintf(os.Stderr, "usage: %s file\n", os.Args[0])
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr, "positional arguments:\n")
    fmt.Fprintf(os.Stderr, "file\t\tfile to compute md5 from\n")
}

func isFile(filename string) bool {
    fi, err := os.Stat(filename)
    if err != nil {
        panic(err)
    }
    if fi.Mode().IsRegular() {
        return true
    }
    return false
}

func filesInDir(dirname string) []string {
    // Return a list of all files in the directory, this is recursive
    var files []string
    walkFn := func(path string, info os.FileInfo, err error) error {
        if isFile(path) {
            files = append(files, path)
        }
        return nil
    }
    filepath.Walk(dirname, walkFn)
    return files
}

func Md5SumFile(filename string) string {
    // Open File
    inFile, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer inFile.Close()
    // Read chunks
    buf := make([]byte, chunkSize)
    hash := md5.New()
    for {
        nbRead, err := inFile.Read(buf)
        if err != nil && err != io.EOF {
            panic(err)
        }
        if nbRead == 0 {
            break
        }
        // Hash Chunk
        hash.Write(buf[:nbRead])
    }
    // Return HexSum
    return hex.EncodeToString(hash.Sum(nil)[:hash.Size()])
}

func VerifyMD5Sum(filenames []string, sums []string) ([]bool, error) {
    result := []bool{}
    if len(filenames) != len(sums) {
        return result, errors.New("filenames and sums sizes missmatch")
    }
    for i := range filenames {
        current_hash := Md5SumFile(filenames[i])
        previous_hash := sums[i]
        if (current_hash == previous_hash) {
            result = append(result, true)
        } else {
            result = append(result, false)
        }
    }
    return result, nil
}

func ParseVerifyFile(filename string) ([]string, []string, error) {
    files := []string{}
    hashes := []string{}
    file, err := os.Open(filename)
    if err != nil {
        return files, hashes, err
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        sliced := strings.Split(scanner.Text(), " ")
        if len(sliced) != 2 {
            return files, hashes, errors.New("Malformated line")
        }
        hashes = append(hashes, sliced[0])
        files = append(files, sliced[1])
    }
    err2 := scanner.Err()
    if err != nil {
        return files, hashes, err2
    }
    return files, hashes, nil
}

func main() {
    // Init command line
    var recursive bool
    var check bool
    var hide_ok bool
    var status_only bool
    flag.BoolVar(&recursive, "recursive", false, "Recurse a directory")
    flag.BoolVar(&check, "check", false, "read MD5 sums from file and check them")
    flag.BoolVar(&hide_ok, "quiet", false, "don't print OK for each successfully verified file")
    flag.BoolVar(&status_only, "status", false, "don't output anything, status code shows success")
    flag.Usage = usage
    flag.Parse()
    if len(flag.Args()) < 1 {
        flag.Usage()
        os.Exit(1)
    }
    inFilename := flag.Arg(0)

    // Get all files
    files := []string{}
    if recursive {
        files = filesInDir(inFilename)
    } else {
        if !isFile(inFilename) {
            fmt.Printf("Path is a directory, please use recursive\n")
            return
        }
        files = append(files, inFilename)
    }

    // Compute MD5 or Check
    if check {
        files, sums, err := ParseVerifyFile(inFilename)
        if err != nil {
            panic(err)
        }
        matching, err := VerifyMD5Sum(files, sums)
        if err != nil {
            panic(err)
        }
        something_failed := false
        for i := range files {
            if matching[i] {
                if (!hide_ok) && (!status_only) {
                    fmt.Printf("%s OK\n", files[i])
                }
            } else {
                if (!status_only) {
                    fmt.Fprintf(os.Stderr, "%s Failed\n", files[i])
                }
                something_failed = true
            }
        }
        if something_failed {
            os.Exit(1)
        }
    } else {
        for _, file := range files {
            hash := Md5SumFile(file)
            fmt.Printf("%s %s\n", hash, file)
        }
    }
}
