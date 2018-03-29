grupo_de_sentencias	    : 	sentencia
						| 	sentencia grupo_de_sentencias 			
						;
						
sentencia 	: 	asignacion 
            | 	sentencia_if 								
            |   for						
            |	funcion				
            ;		
                                                        
asignacion	:	asignacion_simple 
            | 	asignacion_multiple 
            ;
			
asignacion_simple   :   asignacion_simple_izquierda   asignacion_simple_centro   asignacion_simple_derecha
                    ;
					
asignacion_simple_izquierda     :   T_VAR   T_IDENT
                                |   T_IDENT
                                ;

asignacion_simple_centro   :   tipo_dato
                           |   EPSILON
                           ;
                    
tipo_dato   :   T_INT_R
            |   T_FLOAT_R
            |   T_STRING_R
            ;
                
asignacion_simple_derecha   :   T_EQ   datos
                            |   T_COLON_EQ datos
                            ;

datos   :   T_IDENT 
        |   T_INT_V   
        |   T_FLOAT_V 
        |   T_STRING_V 
        ;

asignacion_multiple :   FALTA
                    ;

sentencia_if    :   if grupo_de_elseif else
                ;

if  :   T_IF   condicion   T_CBKT_LEFT   grupo_de_sentencias   T_CBKT_RIGHT
    ;

grupo_de_elseif :   elseif
                |   elseif grupo_de_elseif
                |   EPSILON
                ;

elseif  :   T_ELSE   if
        ;

else    :   T_ELSE   T_CBKT_LEFT   grupo_de_sentencias   T_CBKT_RIGHT   
        |   EPSILON
        ;

condicion   :	datos condicion_derecha
            ;
				

condicion_derecha   :	simbolo_comparador datos
                    |   EPSILON
                    ;		
							
simbolo_comparador	:	T_GREATER_THAN 
                    |   T_GREATER_EQ
					| 	T_LESS_THAN 
					| 	T_LESS_EQ 
					| 	T_EQ_EQ
					|	T_NOT_EQ
                    ;

for :   T_FOR   for_centro  for_derecha
    ;

for_centro  :   condicion   
            |   asignacion_simple   T_SEMICOLON   condicion   T_SEMICOLON   asignacion_simple 
            ;

for_derecha :  T_CBKT_LEFT   grupo_de_sentencias   T_CBKT_RIGHT
            ;