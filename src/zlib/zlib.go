package zlib

type Zlib struct {
    enc []uint32
    jump []Zlib
    val int8
}

func Zlib_init(  obj *Zlib) {
    // compress array -> vector of uint32
    // decompress array -> vector dec
    // populate_jump_table ???
    obj.enc = make( []uint32, len(compress) / 4)
    for i:=0; i < len( obj.enc); i++ {
        obj.enc[i] = uint32(compress[ i*4 + 0]) + uint32(compress[ i*4 +1 ]<<8 ) +
            uint32(compress[i*4 + 2]<<16) + uint32(compress[i*4 + 3]<<24)
    }
    dec := make( []uint32, len(decompress)/4)
    for i:=0; i < len(dec); i++ {
        dec[i] = uint32(decompress[ i*4 + 0]) + uint32(decompress[ i*4 +1 ]<<8 ) +
             uint32(decompress[i*4 + 2]<<16) + uint32(decompress[i*4 + 3]<<24)
    }
    // populate
    obj.jump = make( []Zlib, len( dec))
    base := dec[0] - 4
    for i:=0; i < len(dec); i++ {
        if dec[i] > 0xFF {
            obj.jump[i].jump = obj.jump[ ((dec[i] - base )/ 4) : ]
        } else {
            obj.jump[i].val = int8(dec[i])
            obj.jump[i].jump = nil
            //asserts?
        }
    }
    return
}

func zlib_compressed_size( sz uint32) uint32 {
    return (sz + 7) / 8
}

func jmpBit( table []int8, i uint) int8 {
    return (table[i/8] >> ( i & 7)) & 1
}

func compressSub(fourBytes [4]int8, read uint32, elem uint32, out []int8 ) (outA []int8, state bool ){
    state = true
    if zlib_compressed_size( elem ) > 4 {
        state = false
        return
    }
    outA = make( []int8, len( out))
    copy( outA, out)
    var i uint32
    for i= 0; i < elem; i++ {
        var shift uint8 = uint8( (read + i) &7 )
        var v uint32 = (read + i) / 8
        inv_mask := int8( uint8(0xFF ^ (1 << shift)))
        if int(v) > len( outA) {
            if int(v) < cap( outA) {
                outA = outA[:v]
            } else {
                t := make([]int8, v)
                copy( t, outA)
                outA = t
            }
        }
        outA[v] = (inv_mask & outA[v]) + (jmpBit(fourBytes[:], uint(i)) << shift)
    }
    return
}

func Compress( zlibObj * Zlib, in []int8, in_sz int ) ( out []int8) {
    out = make( []int8, 1)
    var read uint32 = 0
    var i int
    for i=0; i < in_sz; i++ {
        elem := zlibObj.enc[ uint32( int32(in[i]) + 0x180) ]    // cast to int32 so that no overflow when + x180; cast to unsigned afterwards to use as subscript index
        index := uint32( int32(in[i]) + 0x80 )
        // assert(index < zlib.enc.size())
        v := zlibObj.enc[index]
        var fourB [4]int8
        fourB[0] = int8( uint8(v & 0xFF))
        fourB[1] = int8( uint8((v>>8) & 0xFF))
        fourB[2] = int8( uint8((v>>16) & 0xFF))
        fourB[3] = int8( uint8((v>>24) & 0xFF))
        //var state bool
        //out, state = compressSub( fourB, read, elem, out[1:] )
        out, _ = compressSub( fourB, read, elem, out[1:] )
        read  += elem
    }
    out[0] = 1
    return
}

func Decompress(ZlibObj *Zlib, in []int8, inSize int) (out []int8) {
    out = make( []int8, 0)
    if ( in[0] != 1 ) {
        return
    }

    data := in[1:]
    jmp := ZlibObj.jump[0]
    var i int
    for i =1; i < inSize; i++ {
        jmp = jmp.jump[ jmpBit( data, uint(i)) ]
        // assert
        if  jmp.jump[0].jump != nil || jmp.jump[1].jump != nil  {
            continue }

        out = append(out, jmp.jump[3].val )
        jmp = ZlibObj.jump[0]
    }
    return
}

