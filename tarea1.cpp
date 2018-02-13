#include <stdio.h>

int main(){
    FILE *archivoFuente, *archivoSalida; // archivos fuente y salida
    char c; // caracter a leer
    char diagonal = '/';
    char estrella = '*'; 
    char espacio = ' '; 
    char salto = '\n';

    archivoSalida = fopen("salida", "w");
    archivoFuente = fopen("tarea1.cpp", "r");
    /* 
        esto es un comentario lol
    */
    while( (c = getc(archivoFuente)) != EOF ){
        if(c == espacio || c == salto){
            putc(espacio, archivoSalida);
            while((c = getc(archivoFuente)) == espacio || c == salto);
        }

        if (c == diagonal){
            c = getc(archivoFuente);
            if ( c == diagonal){
                // while hasta que encuentre el salto de linea
                while( (c = getc(archivoFuente)) != salto);
            }else if(c == estrella){
                bool bandera = true;
                c = getc(archivoFuente);
                // while hasta encontrar '*/'
                while( bandera ){
                    while( c != estrella){
                        c = getc(archivoFuente);
                    }
                    if( (c = getc(archivoFuente)) == diagonal ) bandera = false;
                }
            }
        }else{
            putc(c, archivoSalida);
        }
    }
    return 0;
}