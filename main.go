package main

import (
    "fmt"
    "zlib"
)

func main() {
    fmt.Println("Hello")
    const HeaderLen int= 28

    compressedIn :=[]uint8  {0x39, 0x3, 0xdc, 0x2, 0x9a, 0x0,
            0x0, 0x0, 0x21, 0x17, 0x8, 0x5d, 0x1,
            0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,0x0,
            0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
            0x1, 0x4, 0x28, 0xeb, 0x84, 0x4e, 0xe8,
            0x47, 0xc2, 0x81, 0x26, 0x94, 0x1b, 0x30,
            0x84, 0x51, 0xbb, 0x11, 0x78, 0x34, 0x41,
            0xcb, 0x87, 0x22, 0xfe, 0xc5, 0x0, 0x0,
            0x0, 0x1, 0x2a, 0xfd, 0xab, 0xef, 0x81,
            0x10, 0xf1, 0x8, 0xc9, 0x6, 0xde, 0x7e,
            0x90, 0xea, 0xa1 }
    decompressedIn := []uint8 { 0x39, 0x3, 0xdc, 0x2, 0x9a, 0x0,
            0x0, 0x0, 0x21, 0x17, 0x8, 0x5d, 0x1,
            0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,0x0,
            0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
            0x15, 0x10, 0x39, 0x3, 0x8f, 0xc2, 0xfb, 0xc2,
            0x0, 0x0, 0x78, 0xc2, 0xcd, 0xc, 0x89, 0x43,
            0x0, 0x0, 0x2, 0x0, 0x5d, 0x0, 0x0, 0x0,
            0x69, 0x59, 0x9a, 0x67, 0x0, 0x0, 0x0, 0x0,
            0x0, 0x0, 0x0, 0x0, 0xf1, 0x8, 0xc9, 0x6,
            0xde, 0x7e, 0x90, 0xea, 0xa1 }

    compressed := make( []int8, len( compressedIn ))
    for i, v := range compressedIn {
        compressed[i] = int8( v)
    }

    decompressed := make( []int8, len(decompressedIn ))
    for i,v := range decompressedIn {
        decompressed[i] = int8( v)
    }

    var zLibObj zlib.Zlib
    zlib.Zlib_init(&zLibObj)
    // no ideas why len-header - 5 - 16 times 8
    tailDec := zlib.Decompress(&zLibObj, compressed[HeaderLen:], (len(compressed) - HeaderLen - 5 - 16) * 8 )

    decompressionPassed := true
    fmt.Println("Decompressed Original len =", len( decompressed), " (w/o header ", len( decompressed) - HeaderLen, ") and len of decompressed resulted ", len( tailDec ))
    for i :=0; i < len( tailDec ) && i + HeaderLen < len(decompressed) ; i++ {
        if tailDec[i] != decompressed[i + HeaderLen] {
            fmt.Println("mismatch in Decompressed out at subposition ", i, " : expected ",
                decompressed[i + HeaderLen], "   but got", tailDec[i] )
            decompressionPassed = false
        }
    }
    if decompressionPassed {
        fmt.Println("Check of decompression passed")
    } else {
        fmt.Println("Failed check of decompression" )
    }

    tailEnc := zlib.Compress(&zLibObj, decompressed[HeaderLen:], len(decompressed) - HeaderLen )
    compressionPassed := true
    fmt.Println("Compressed Original len =", len( compressed), " w/o header is ", len( compressed) - HeaderLen, "  and len of decompressed resulted ", len( tailEnc ))
    for i :=0; i < len( tailEnc ) && i + HeaderLen < len( compressed); i++ {
        if tailEnc[i] != compressed[i + HeaderLen] {
            fmt.Println("mismatch in Decompressed out at subposition ", i, " : expected ",
                compressed[i + HeaderLen], "   but got", tailEnc[i] )
            compressionPassed = false
        }
    }
    if compressionPassed {
        fmt.Println("Check of compression passed")
    } else {
        fmt.Println("Failed check of compression" )
    }


    fmt.Println("Done")
}

