#include <stdio.h>
#include <stdlib.h>

int main( int argc, char * argv[])
{
    FILE * fin;
    int lineC =0;
    unsigned char Char;
    int rv;
    if( argc < 2 ) { fprintf(stderr, "input file name should be provided. \nExiting\n"); exit( EXIT_FAILURE); }
    fin = fopen( argv[1], "rb");
    if( NULL == fin) { fprintf( stderr, "failed to open input file\n"); exit( EXIT_FAILURE); }

    while( ! feof( fin)) {
        rv = fread(&Char, 1, 1, fin);
        if( rv < 1 ) continue;
        printf("0x%02hhx,", Char);
        lineC++;
        if( lineC == 16 ) {
            printf("\n");
            lineC =0;
        } else {
            printf(" ");
        }
    }
    fclose( fin);
};
