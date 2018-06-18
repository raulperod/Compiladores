new Vue({
    el: '#app',
    data: {
        code: `// codigo de ejemplo
package main

import "fmt"
        
func main(){
    fmt.Println("Hola Mundo")
}`      ,
        message: ''
    },
    methods: {
        ejecutar(){
            fetch('/parser', {
                method: 'POST',
                body: JSON.stringify(this.code),
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                }
            })
            .then(res => res.json())
            .then(data => {
                this.message = data.msg
            })
        }
    }
})